// Package apikey provides the API Keys API.
package apikey

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/api_keys"

// Client is used to invoke the API Keys API.
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
	Type                   *string         `json:"type,omitempty"`
	Name                   *string         `json:"name,omitempty"`
	Description            *string         `json:"description,omitempty"`
	Subject                *string         `json:"subject,omitempty"`
	Claims                 json.RawMessage `json:"claims,omitempty"`
	Scopes                 []string        `json:"scopes,omitempty"`
	CreatedBy              *string         `json:"created_by,omitempty"`
	SecondsUntilExpiration *int64          `json:"seconds_until_expiration,omitempty"`
}

// Create creates a new API key.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.APIKeyWithSecret, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.APIKeyWithSecret{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Get returns an API key
func (c *Client) Get(ctx context.Context, apiKeyID string) (*clerk.APIKey, error) {
	path, err := clerk.JoinPath(path, apiKeyID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.APIKey{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	Type           *string `json:"type,omitempty"`
	Subject        *string `json:"subject,omitempty"`
	IncludeInvalid *string `json:"include_invalid,omitempty"`
}

func (params *ListParams) ToQuery() url.Values {
	q := url.Values{}
	if params.Type != nil {
		q.Set("type", *params.Type)
	}
	if params.Subject != nil {
		q.Set("subject", *params.Subject)
	}
	if params.IncludeInvalid != nil {
		q.Set("include_invalid", *params.IncludeInvalid)
	}
	// Add list params
	if params.Limit != nil {
		q.Set("limit", strconv.FormatInt(*params.Limit, 10))
	}
	if params.Offset != nil {
		q.Set("offset", strconv.FormatInt(*params.Offset, 10))
	}
	return q
}

// List returns a list of API keys.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.APIKeyList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &clerk.APIKeyList{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// GetSecret retrieves the secret for an API key.
func (c *Client) GetSecret(ctx context.Context, apiKeyID string) (*APIKeySecret, error) {
	path, err := clerk.JoinPath(path, apiKeyID, "secret")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &APIKeySecret{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type APIKeySecret struct {
	clerk.APIResource
	Secret string `json:"secret"`
}

type UpdateParams struct {
	clerk.APIParams
	Claims                 json.RawMessage `json:"claims,omitempty"`
	Scopes                 []string        `json:"scopes,omitempty"`
	Description            *string         `json:"description,omitempty"`
	Subject                *string         `json:"subject,omitempty"`
	SecondsUntilExpiration *int64          `json:"seconds_until_expiration,omitempty"`
}

// Update updates an API key.
func (c *Client) Update(ctx context.Context, apiKeyID string, params *UpdateParams) (*clerk.APIKey, error) {
	path, err := clerk.JoinPath(path, apiKeyID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	resource := &clerk.APIKey{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Delete deletes an API key
func (c *Client) Delete(ctx context.Context, apiKeyID string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, apiKeyID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type RevokeParams struct {
	clerk.APIParams
	RevocationReason *string `json:"revocation_reason,omitempty"`
}

// Revoke revokes an API key.
func (c *Client) Revoke(ctx context.Context, apiKeyID string, params *RevokeParams) (*clerk.APIKey, error) {
	path, err := clerk.JoinPath(path, apiKeyID, "revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.APIKey{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type VerifyParams struct {
	clerk.APIParams
	Secret string `json:"secret"`
}

// Verify verifies an API key.
func (c *Client) Verify(ctx context.Context, params *VerifyParams) (*clerk.APIKey, error) {
	path, err := clerk.JoinPath(path, "verify")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.APIKey{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
