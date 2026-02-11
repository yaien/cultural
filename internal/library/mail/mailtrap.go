package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MailtrapOptions struct {
	Token   string
	Sandbox bool
	InboxID string
	Timeout time.Duration
}

func NewMailtrap(options MailtrapOptions) (*MailtrapClient, error) {
	var c MailtrapClient

	if options.Token == "" {
		return nil, fmt.Errorf("mailtrap api token is required")
	}
	c.token = options.Token

	c.endpoint = "https://send.api.mailtrap.io/api/send"
	if options.Sandbox {
		if options.InboxID == "" {
			return nil, fmt.Errorf("inbox ID is required for sandbox mode")
		}
		c.endpoint = fmt.Sprintf("https://sandbox.api.mailtrap.io/api/send/%s", options.InboxID)
	}

	var h http.Client
	if options.Timeout > 0 {
		h.Timeout = options.Timeout
	}

	c.http = &h

	return &c, nil
}

type MailtrapClient struct {
	endpoint string
	token    string
	http     *http.Client
}

type mailtrapAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailtrapPayload struct {
	From     mailtrapAddress   `json:"from"`
	To       []mailtrapAddress `json:"to"`
	Subject  string            `json:"subject"`
	HTML     string            `json:"html,omitempty"`
	Category string            `json:"category,omitempty"`
}

func (c *MailtrapClient) Send(ctx context.Context, email *Email) error {
	if c == nil {
		return fmt.Errorf("mailtrap client is nil")
	}
	if c.token == "" {
		return fmt.Errorf("mailtrap api token is required")
	}
	if email == nil {
		return fmt.Errorf("email is nil")
	}
	if email.To.Email == "" {
		return fmt.Errorf("recipient email is required")
	}
	if email.From.Email == "" {
		return fmt.Errorf("sender email is required")
	}
	if email.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	payload := c.payload(email)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal mailtrap payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create mailtrap request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mailtrap request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mailtrap request returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *MailtrapClient) payload(email *Email) mailtrapPayload {
	return mailtrapPayload{
		From: mailtrapAddress{
			Email: email.From.Email,
			Name:  email.From.Name,
		},
		To: []mailtrapAddress{
			{
				Email: email.To.Email,
				Name:  email.To.Name,
			},
		},
		Subject:  email.Subject,
		HTML:     email.Body,
		Category: email.Category,
	}
}
