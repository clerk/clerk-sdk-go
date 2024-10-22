// Package session provides the Sign In Token API.
package signintokens

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/sign_in_tokens"

// Client is used to invoke the sign-in Token API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type CreateParams struct {
	clerk.APIParams
	UserID           *string `json:"user_id,omitempty"`
	ExpiresInSeconds *int64  `json:"expires_in_seconds,omitempty"`
}

// Create creates a new sign-in token.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.SignInToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	token := &clerk.SignInToken{}
	err := c.Backend.Call(ctx, req, token)
	return token, err
}

// Revoke revokes a pending sign-in token.
func (c *Client) Revoke(ctx context.Context, id string) (*clerk.SignInToken, error) {
	token := &clerk.SignInToken{}
	path, err := clerk.JoinPath(path, id, "revoke")
	if err != nil {
		return token, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	err = c.Backend.Call(ctx, req, token)
	return token, err
}
