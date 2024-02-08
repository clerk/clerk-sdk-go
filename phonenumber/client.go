// Package phonenumber provides the Phone Numbers API.
//
// https://clerk.com/docs/reference/backend-api/tag/Phone-Numbers
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

	UserID      *string `json:"user_id,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Verified    *bool   `json:"verified,omitempty"`
	Primary     *bool   `json:"primary,omitempty"`
}

// Create creates a new phone number.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.PhoneNumber, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	template := &clerk.PhoneNumber{}
	err := c.Backend.Call(ctx, req, template)
	return template, err
}

// Get returns details about a phone number.
func (c *Client) Get(ctx context.Context, id string) (*clerk.PhoneNumber, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	template := &clerk.PhoneNumber{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
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
	template := &clerk.PhoneNumber{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
}

// Delete deletes a phone number.
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
