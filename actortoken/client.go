// Package actortoken provides the Actor Tokens API.
package actortoken

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/actor_tokens"

// Client is used to invoke the Actor Tokens API.
type Client struct {
	Backend clerk.Backend
}

type ClientConfig struct {
	clerk.BackendConfig
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type CreateParams struct {
	clerk.APIParams
	UserID                      *string         `json:"user_id,omitempty"`
	Actor                       json.RawMessage `json:"actor,omitempty"`
	ExpiresInSeconds            *int64          `json:"expires_in_seconds,omitempty"`
	SessionMaxDurationInSeconds *int64          `json:"session_max_duration_in_seconds,omitempty"`
}

// Create creates a new actor token.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.ActorToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	token := &clerk.ActorToken{}
	err := c.Backend.Call(ctx, req, token)
	return token, err
}

// Revoke revokes a pending actor token.
func (c *Client) Revoke(ctx context.Context, id string) (*clerk.ActorToken, error) {
	token := &clerk.ActorToken{}
	path, err := clerk.JoinPath(path, id, "revoke")
	if err != nil {
		return token, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	err = c.Backend.Call(ctx, req, token)
	return token, err
}
