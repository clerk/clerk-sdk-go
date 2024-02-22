// Package jwks provides access to the JWKS endpoint.
package jwks

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/jwks"

// Client is used to invoke the JWKS API.
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

type GetParams struct {
	clerk.APIParams
}

// Get retrieves a JSON Web Key set.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.JSONWebKeySet, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &clerk.JSONWebKeySet{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}
