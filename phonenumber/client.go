// Package phonenumber provides the Phone Numbers API.
package phonenumber

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/phone_numbers"

// Client is used to invoke the Phone Numbers API.
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
	UserID      *string `json:"user_id,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Verified    *bool   `json:"verified,omitempty"`
	Primary     *bool   `json:"primary,omitempty"`
}

// Create creates a new phone number.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.PhoneNumber, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.PhoneNumber{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Get retrieves a phone number.
func (c *Client) Get(ctx context.Context, id string) (*clerk.PhoneNumber, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.PhoneNumber{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type UpdateParams struct {
	clerk.APIParams
	Verified *bool `json:"verified,omitempty"`
	Primary  *bool `json:"primary,omitempty"`
}

// Update updates the phone number specified by id.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.PhoneNumber, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	resource := &clerk.PhoneNumber{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Delete deletes a phone number.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
