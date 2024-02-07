// Package allowlistidentifier provides the Allowlist Identifiers API.
package allowlistidentifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/allowlist_identifiers"

// Client is used to invoke the Allowlist Identifiers API.
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
	Identifier *string `json:"identifier,omitempty"`
	Notify     *bool   `json:"notify,omitempty"`
}

// Create adds a new identifier to the allowlist.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.AllowlistIdentifier, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	identifier := &clerk.AllowlistIdentifier{}
	err := c.Backend.Call(ctx, req, identifier)
	return identifier, err
}

// Delete removes an identifier from the allowlist.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := url.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	identifier := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, identifier)
	return identifier, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns all the identifiers in the allowlist.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.AllowlistIdentifierList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	list := &clerk.AllowlistIdentifierList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
