// Package svixwebhook provides the Svix Webhooks API.
package svixwebhook

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/webhooks/svix"

// Client is used to invoke the Organizations API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// Create creates a Svix app.
func (c *Client) Create(ctx context.Context) (*clerk.SvixWebhook, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.SvixWebhook{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Delete deletes the Svix app.
func (c *Client) Delete(ctx context.Context) (*clerk.SvixWebhook, error) {
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.SvixWebhook{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// RefreshURL generates a new URL for accessing Svix's dashboard.
func (c *Client) RefreshURL(ctx context.Context) (*clerk.SvixWebhook, error) {
	req := clerk.NewAPIRequest(http.MethodPost, "/webhooks/svix_url")
	resource := &clerk.SvixWebhook{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}
