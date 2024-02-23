// Package jwttemplate provides the JWT Templates API.
package jwttemplate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/jwt_templates"

// Client is used to invoke the JWT Templates API.
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
	Name             *string         `json:"name,omitempty"`
	Claims           json.RawMessage `json:"claims,omitempty"`
	Lifetime         *int64          `json:"lifetime,omitempty"`
	AllowedClockSkew *int64          `json:"allowed_clock_skew,omitempty"`
	CustomSigningKey *bool           `json:"custom_signing_key,omitempty"`
	SigningKey       *string         `json:"signing_key,omitempty"`
	SigningAlgorithm *string         `json:"signing_algorithm,omitempty"`
}

// Create creates a new JWT template.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.JWTTemplate, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	template := &clerk.JWTTemplate{}
	err := c.Backend.Call(ctx, req, template)
	return template, err
}

// Get returns details about a JWT template.
func (c *Client) Get(ctx context.Context, id string) (*clerk.JWTTemplate, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	template := &clerk.JWTTemplate{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
}

type UpdateParams struct {
	clerk.APIParams
	Name             *string         `json:"name,omitempty"`
	Claims           json.RawMessage `json:"claims,omitempty"`
	Lifetime         *int64          `json:"lifetime,omitempty"`
	AllowedClockSkew *int64          `json:"allowed_clock_skew,omitempty"`
	CustomSigningKey *bool           `json:"custom_signing_key,omitempty"`
	SigningKey       *string         `json:"signing_key,omitempty"`
	SigningAlgorithm *string         `json:"signing_algorithm,omitempty"`
}

// Update updates the JWT template specified by id.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.JWTTemplate, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	template := &clerk.JWTTemplate{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
}

// Delete deletes a JWT template.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	template := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns a list of JWT templates.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.JWTTemplateList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	req.SetParams(params)
	list := &clerk.JWTTemplateList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
