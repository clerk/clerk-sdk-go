// Package email provides the Email API.
//
// Experimental: The Email API is internal and not yet public. It is subject to
// change. It is advised to pin the SDK version to avoid breaking changes. See
// https://clerk.com/docs/pinning.
package email

import (
	"context"
	"errors"
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

// Recipient is the addressee of an email. Provide exactly one of Address or
// UserID; they are mutually exclusive. When UserID is set, Clerk resolves that
// user's primary email address on the server, from the instance the secret key
// belongs to.
type Recipient struct {
	// Name is an optional display name. Only meaningful alongside Address. It is
	// currently accepted but not yet rendered by the server.
	Name string `json:"name,omitempty"`
	// Address is a literal recipient email address.
	Address string `json:"address,omitempty"`
	// UserID sends to the primary email address of the Clerk user with this ID.
	UserID string `json:"user_id,omitempty"`
}

type SendParams struct {
	clerk.APIParams
	// IdempotencyKey deduplicates retries of the same logical send. Reuse a
	// key only when recipient and content are identical. The server replay
	// window is 24 hours; keep a durable logical-send record if duplicates after
	// that window matter. Keys may contain only ASCII letters, digits,
	// underscores, and hyphens, up to 255 characters.
	IdempotencyKey string    `json:"-"`
	To             Recipient `json:"to"`
	From           Mailbox   `json:"from"`
	ReplyTo        *Mailbox  `json:"reply_to,omitempty"`
	Subject        string    `json:"subject"`
	HTML           string    `json:"html,omitempty"`
	Text           string    `json:"text,omitempty"`
}

// Send sends a transactional email.
//
// Experimental: This method calls an internal, not-yet-public endpoint and is
// subject to change. It is advised to pin the SDK version to avoid breaking
// changes. See https://clerk.com/docs/pinning.
func (c *Client) Send(ctx context.Context, params *SendParams) (*clerk.Email, error) {
	if params == nil {
		return nil, errors.New("email: send params are required")
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	if params.IdempotencyKey != "" {
		req.SetIdempotencyKey(params.IdempotencyKey)
	}
	resource := &clerk.Email{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Get returns Clerk's stored provider-acceptance state for a transactional
// email. An accepted status does not prove final delivery.
func (c *Client) Get(ctx context.Context, emailID string) (*clerk.Email, error) {
	if emailID == "" {
		return nil, errors.New("email: email ID is required")
	}
	req := clerk.NewAPIRequest(http.MethodGet, path+"/"+emailID)
	resource := &clerk.Email{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}
