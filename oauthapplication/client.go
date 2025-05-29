// Package oauthapplication provides the OAuth applications API.
package oauthapplication

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/oauth_applications"

type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// Get fetches a single OAuth application by its ID.
func (c *Client) Get(ctx context.Context, id string) (*clerk.OAuthApplication, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.OAuthApplication{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrderBy *string `json:"order_by,omitempty"`
}

func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()

	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}

	return q
}

// List retrieves all OAuth applications.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OAuthApplicationList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OAuthApplicationList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type CreateParams struct {
	clerk.APIParams
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	Scopes       *string  `json:"scopes,omitempty"`
	Public       *bool    `json:"public,omitempty"`

	// Deprecated: Use RedirectURIs instead
	CallbackURL *string `json:"callback_url,omitempty"`
}

// Create creates a new OAuth application with the given parameters.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.OAuthApplication, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	authApplication := &clerk.OAuthApplication{}
	err := c.Backend.Call(ctx, req, authApplication)
	return authApplication, err
}

type UpdateParams struct {
	clerk.APIParams
	Name         *string  `json:"name,omitempty"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	Scopes       *string  `json:"scopes,omitempty"`
	Public       *bool    `json:"public,omitempty"`

	// Deprecated: Use RedirectURIs instead
	CallbackURL *string `json:"callback_url,omitempty"`
}

// Update updates an existing OAuth application.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.OAuthApplication, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	authApplication := &clerk.OAuthApplication{}
	err = c.Backend.Call(ctx, req, authApplication)
	return authApplication, err
}

// Delete deletes the given OAuth application
func (c *Client) DeleteOAuthApplication(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	authApplication := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, authApplication)
	return authApplication, err
}

// RotateClientSecret rotates the OAuth application's client secret
func (c *Client) RotateClientSecret(ctx context.Context, id string) (*clerk.OAuthApplication, error) {
	path, err := clerk.JoinPath(path, id, "rotate_secret")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	authApplication := &clerk.OAuthApplication{}
	err = c.Backend.Call(ctx, req, authApplication)
	return authApplication, err
}
