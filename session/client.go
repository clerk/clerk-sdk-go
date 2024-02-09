// Package session provides the Sessions API.
package session

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/sessions"

// Client is used to invoke the Sessions API.
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

// Get retrieves details for a session.
func (c *Client) Get(ctx context.Context, id string) (*clerk.Session, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	session := &clerk.Session{}
	err = c.Backend.Call(ctx, req, session)
	return session, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	ClientID *string `json:"client_id,omitempty"`
	UserID   *string `json:"user_id,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// ToQuery returns the params as url.Values.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.ClientID != nil {
		q.Add("client_id", *params.ClientID)
	}
	if params.UserID != nil {
		q.Add("user_id", *params.UserID)
	}
	if params.Status != nil {
		q.Add("status", *params.Status)
	}
	return q
}

// List returns a list of sessions.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.SessionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	req.SetParams(params)
	list := &clerk.SessionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type RevokeParams struct {
	ID string `json:"id"`
}

// Revoke marks the session as revoked.
func (c *Client) Revoke(ctx context.Context, params *RevokeParams) (*clerk.Session, error) {
	path, err := clerk.JoinPath(path, params.ID, "/revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	session := &clerk.Session{}
	err = c.Backend.Call(ctx, req, session)
	return session, err
}

type VerifyParams struct {
	ID    string  `json:"-"`
	Token *string `json:"token,omitempty"`
}

// Verify verifies the session.
//
// Deprecated: The operation is deprecated and will be removed in future versions.
// It is recommended to switch to networkless verification using short-lived
// session tokens instead.
// See https://clerk.com/docs/backend-requests/resources/session-tokens
func (c *Client) Verify(ctx context.Context, params *VerifyParams) (*clerk.Session, error) {
	path, err := clerk.JoinPath(path, params.ID, "/verify")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	session := &clerk.Session{}
	err = c.Backend.Call(ctx, req, session)
	return session, err
}
