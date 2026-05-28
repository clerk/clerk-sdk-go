// Package phonenumber provides the Phone Numbers API.
package phonenumber

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
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
	Verified                *bool `json:"verified,omitempty"`
	Primary                 *bool `json:"primary,omitempty"`
	ReservedForSecondFactor *bool `json:"reserved_for_second_factor,omitempty"`
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

type ReplaceForUserParams struct {
	clerk.APIParams
	UserID      string  `json:"-"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// ReplaceForUser replaces all of the user's phone numbers with a single
// verified, primary phone number. The new phone number is created with the
// admin verification strategy and is not reserved for second factor. Any
// existing phone numbers are deleted; replacing a phone number that is reserved
// for second factor disables the user's MFA.
func (c *Client) ReplaceForUser(ctx context.Context, params *ReplaceForUserParams) (*clerk.PhoneNumber, error) {
	path, err := clerk.JoinPath("/users", params.UserID, "/phone_number")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPut, path)
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
