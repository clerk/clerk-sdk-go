// Package redirecturl provides the Redirect URLs API.
package redirecturl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/redirect_urls"

// Client is used to invoke the Redirect URLs API.
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
	URL *string `json:"url,omitempty"`
}

// Create creates a new redirect url.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.RedirectURL, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	redirectURL := &clerk.RedirectURL{}
	err := c.Backend.Call(ctx, req, redirectURL)
	return redirectURL, err
}

// Get retrieves details for a redirect url by ID.
func (c *Client) Get(ctx context.Context, id string) (*clerk.RedirectURL, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	redirectURL := &clerk.RedirectURL{}
	err = c.Backend.Call(ctx, req, redirectURL)
	return redirectURL, err
}

// Delete deletes a redirect url.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	redirectURL := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, redirectURL)
	return redirectURL, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns a list of redirect urls.
func (c *Client) List(ctx context.Context, _ *ListParams) (*clerk.RedirectURLList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	list := &clerk.RedirectURLList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
