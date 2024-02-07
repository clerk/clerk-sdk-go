// Package blocklistidentifier provides the Blocklist Identifiers API.
package blocklistidentifier

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/blocklist_identifiers"

// Client is used to invoke the Blocklist Identifiers API.
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
}

// Create adds a new identifier to the blocklist.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.BlocklistIdentifier, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	identifier := &clerk.BlocklistIdentifier{}
	err := c.Backend.Call(ctx, req, identifier)
	return identifier, err
}

// Delete removes an identifier from the blocklist.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
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

// List returns all the identifiers in the blocklist.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.BlocklistIdentifierList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	list := &clerk.BlocklistIdentifierList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
