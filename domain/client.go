// Package domain provides the Domains API.
package domain

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/domains"

// Client is used to invoke the Domains API.
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
	Name        *string `json:"name,omitempty"`
	ProxyURL    *string `json:"proxy_url,omitempty"`
	IsSatellite *bool   `json:"is_satellite,omitempty"`
}

// Create creates a new domain.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.Domain, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)

	domain := &clerk.Domain{}
	err := c.Backend.Call(ctx, req, domain)
	return domain, err
}

type UpdateParams struct {
	clerk.APIParams
	Name        *string `json:"name,omitempty"`
	ProxyURL    *string `json:"proxy_url,omitempty"`
	IsSecondary *bool   `json:"is_secondary,omitempty"`
}

// Update updates a domain's properties.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.Domain, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)

	domain := &clerk.Domain{}
	err = c.Backend.Call(ctx, req, domain)
	return domain, err
}

// Delete removes a domain.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	domain := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, domain)
	return domain, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns a list of domains.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.DomainList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.DomainList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
