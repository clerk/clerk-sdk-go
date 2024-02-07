// Package emailaddress provides the Email Addresses API.
package emailaddress

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/email_addresses"

// Client is used to invoke the Email Addresses API.
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
	UserID       *string `json:"user_id,omitempty"`
	EmailAddress *string `json:"email_address,omitempty"`
	Verified     *bool   `json:"verified,omitempty"`
	Primary      *bool   `json:"primary,omitempty"`
}

// Create creates a new email address.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.EmailAddress, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	emailAddress := &clerk.EmailAddress{}
	err := c.Backend.Call(ctx, req, emailAddress)
	return emailAddress, err
}

// Get retrieves an email address.
func (c *Client) Get(ctx context.Context, id string) (*clerk.EmailAddress, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	emailAddress := &clerk.EmailAddress{}
	err = c.Backend.Call(ctx, req, emailAddress)
	return emailAddress, err
}

type UpdateParams struct {
	clerk.APIParams
	Verified *bool `json:"verified,omitempty"`
	Primary  *bool `json:"primary,omitempty"`
}

// Update updates the email address specified by id.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.EmailAddress, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	emailAddress := &clerk.EmailAddress{}
	err = c.Backend.Call(ctx, req, emailAddress)
	return emailAddress, err
}

// Delete deletes an email address.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	emailAddress := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, emailAddress)
	return emailAddress, err
}
