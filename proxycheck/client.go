// Package proxycheck provides the Proxy Checks API.
//
// Deprecated: The Proxy Checks API is deprecated and will be
// removed in future versions.
package proxycheck

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/proxy_checks"

// Client is used to invoke the Proxy Checks API.
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
	DomainID *string `json:"domain_id,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
}

// Create creates a proxy check.
//
// Deprecated: The operation is deprecated and will be removed in
// future versions.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.ProxyCheck, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.ProxyCheck{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}
