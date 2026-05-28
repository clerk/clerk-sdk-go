// Package emailaddress provides the Email Addresses API.
package emailaddress

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/email_addresses"

// Client is used to invoke the Email Addresses API.
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

type ReplaceForUserParams struct {
	clerk.APIParams
	UserID       string  `json:"-"`
	EmailAddress *string `json:"email_address,omitempty"`
}

// ReplaceForUser replaces all of the user's email addresses with a single
// verified, primary email address. The new email address is created with the
// admin verification strategy and any existing email addresses are deleted. If
// an existing email address is linked to a connected account, the request is
// rejected; remove the connected account first.
func (c *Client) ReplaceForUser(ctx context.Context, params *ReplaceForUserParams) (*clerk.EmailAddress, error) {
	path, err := clerk.JoinPath("/users", params.UserID, "/email_address")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPut, path)
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
