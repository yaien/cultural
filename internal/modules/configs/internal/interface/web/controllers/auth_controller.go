package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/auth"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"golang.org/x/oauth2"
)

type AuthController struct {
	accounts *auth.Accounts
	store    sessions.Store
	config   *oauth2.Config
}

type AuthState struct {
	Redirect string
}

func NewAuthController(accounts *auth.Accounts, store sessions.Store, config *oauth2.Config) *AuthController {
	return &AuthController{accounts, store, config}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	state := AuthState{
		Redirect: r.URL.Query().Get("redirect"),
	}

	bs, err := json.Marshal(state)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to marshal state: %w", err))
		return
	}

	s := base64.RawURLEncoding.EncodeToString(bs)

	url := c.config.AuthCodeURL(s)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {

	session, _ := c.store.Get(r, middlewares.SessionKey)
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)

	if err := session.Save(r, w); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to clear session: %w", err))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := c.config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to exchange token: %w", err))
		return
	}

	client := c.config.Client(ctx, token)

	// Call Google's userinfo endpoint to get profile information.
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to get user info: %w", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		WriteHTMLErr(w, fmt.Errorf("unexpected status from userinfo endpoint: %s", resp.Status))
		return
	}

	// Decode the JSON response.
	var body struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to decode user info: %w", err))
		return
	}

	account := &auth.Account{
		Provider:     "google",
		ID:           body.ID,
		Email:        body.Email,
		Name:         body.Name,
		AvatarUrl:    body.Picture,
		AcccessToken: token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		LastUsedAt:   time.Now(),
	}

	user, err := c.accounts.Sync(ctx, account)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to sync account: %w", err))
		return
	}

	session, _ := c.store.Get(r, middlewares.SessionKey)
	session.Values[middlewares.UserIDKey] = user.ID.Hex()
	if err := session.Save(r, w); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to save session: %w", err))
		return
	}

	var state AuthState
	bs, err := base64.RawURLEncoding.DecodeString(r.URL.Query().Get("state"))
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to decode state: %w", err))
		return
	}

	if err := json.Unmarshal(bs, &state); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to unmarshal state: %w", err))
		return
	}

	redirect := "/"
	if state.Redirect != "" {
		redirect = state.Redirect
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
