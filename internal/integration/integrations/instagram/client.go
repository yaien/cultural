package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	config *oauth2.Config
	token  *oauth2.Token
	base   string
}

func NewClient(token *oauth2.Token, config *oauth2.Config) *Client {
	return &Client{
		config: config,
		token:  token,
		base:   "https://graph.instagram.com",
	}
}

func (c *Client) SetToken(token *oauth2.Token) {
	c.token = token
}

func (c *Client) Token() *oauth2.Token {
	return c.token
}

func (c *Client) GetLongToken(ctx context.Context) (*oauth2.Token, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/access_token", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	query := req.URL.Query()
	query.Set("grant_type", "ig_exchange_token")
	query.Set("client_secret", c.config.ClientSecret)
	query.Set("access_token", c.token.AccessToken)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("failed clossing instagram body", "err", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		var e Error
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("error from API: %s - %s", e.Error.Type, e.Error.Message)
	}

	var token oauth2.Token

	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed decoding response: %w", err)
	}

	return &token, nil
}

func (c *Client) RefreshLongToken(ctx context.Context) (*oauth2.Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/refresh_access_token", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	query := req.URL.Query()
	query.Set("grant_type", "ig_refresh_token")
	if c.token != nil {
		query.Set("access_token", c.token.AccessToken)
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("failed clossing instagram body", "err", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		var e Error
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("error from API: %s - %s", e.Error.Type, e.Error.Message)
	}

	var token oauth2.Token
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed decoding response: %w", err)
	}

	return &token, nil
}

type User struct {
	ID                string `json:"id" bson:"id"`
	Name              string `json:"name" bson:"name"`
	Username          string `json:"username" bson:"username"`
	ProfilePictureURL string `json:"profile_picture_url" bson:"profilePictureUrl"`
}

type Error struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}
}

func (c *Client) GetUser(ctx context.Context) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/me", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	query := req.URL.Query()
	fields := "id,name,username,profile_picture_url"
	query.Set("fields", fields)
	if c.token != nil {
		query.Set("access_token", c.token.AccessToken)
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("failed clossing instagram body", "err", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		var e Error
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("error from API: %s - %s", e.Error.Type, e.Error.Message)
	}

	var me User
	if err := json.NewDecoder(res.Body).Decode(&me); err != nil {
		return nil, fmt.Errorf("failed decoding response: %w", err)
	}

	return &me, nil
}

type Child struct {
	ID        string `json:"id" bson:"id"`
	MediaURL  string `json:"media_url" bson:"mediaUrl"`
	MediaType string `json:"media_type" bson:"mediaType"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
}

type Post struct {
	ID        string `json:"id" bson:"id"`
	Caption   string `json:"caption" bson:"caption"`
	MediaURL  string `json:"media_url" bson:"mediaUrl"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
	MediaType string `json:"media_type" bson:"mediaType"`
	Permalink string `json:"permalink" bson:"permalink"`
	Children  *struct {
		Data []Child `json:"data" bson:"data"`
	} `json:"children,omitempty" bson:"children,omitempty"`
}

func (c *Client) GetPosts(ctx context.Context) ([]*Post, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/me/media", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	query := req.URL.Query()
	fields := "id,caption,media_url,timestamp,media_type,permalink,children{media_url,media_type,timestamp}"
	query.Set("fields", fields)

	if c.token != nil {
		query.Set("access_token", c.token.AccessToken)
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Warn("failed clossing instagram body", "err", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		var e Error
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("error from API: %s - %s", e.Error.Type, e.Error.Message)
	}

	var resp struct {
		Data []*Post `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed decoding response: %w", err)
	}

	return resp.Data, nil
}
