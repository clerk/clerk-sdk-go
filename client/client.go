// Package client provides the Client API.
package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/clients"

// Client is used to invoke the Client API.
// This is an API client for interacting with Clerk Client resources.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// Get retrieves the client specified by ID.
func (c *Client) Get(ctx context.Context, id string) (*clerk.Client, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	client := &clerk.Client{}
	err = c.Backend.Call(ctx, req, client)
	return client, err
}

type VerifyParams struct {
	clerk.APIParams
	Token *string `json:"token,omitempty"`
}

// Verify verifies the Client in the provided JWT.
func (c *Client) Verify(ctx context.Context, params *VerifyParams) (*clerk.Client, error) {
	path, err := clerk.JoinPath(path, "/verify")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	client := &clerk.Client{}
	err = c.Backend.Call(ctx, req, client)
	return client, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
}

func (params *ListParams) ToQuery() url.Values {
	return params.ListParams.ToQuery()
}

// List returns a list of all the clients.
//
// Deprecated: The operation is deprecated and will be removed in
// future versions.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.ClientList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.ClientList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
