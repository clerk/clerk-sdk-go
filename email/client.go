// Package email provides the Email API.
package email

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/email"

// Client is used to invoke the Email API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// Mailbox is an email address, with an optional display name.
type Mailbox struct {
	// Name is the display name for the mailbox. It is currently
	// accepted but not yet rendered by the server.
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
}

type SendParams struct {
	clerk.APIParams
	To      Mailbox  `json:"to"`
	From    Mailbox  `json:"from"`
	ReplyTo *Mailbox `json:"reply_to,omitempty"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

// Send sends a transactional email.
func (c *Client) Send(ctx context.Context, params *SendParams) (*clerk.Email, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.Email{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}
