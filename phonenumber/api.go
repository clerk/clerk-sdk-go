// Code generated by "gen"; DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package phonenumber

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
)

// Create creates a new phone number.
func Create(ctx context.Context, params *CreateParams) (*clerk.PhoneNumber, error) {
	return getClient().Create(ctx, params)
}

// Get retrieves a phone number.
func Get(ctx context.Context, id string) (*clerk.PhoneNumber, error) {
	return getClient().Get(ctx, id)
}

// Update updates the phone number specified by id.
func Update(ctx context.Context, id string, params *UpdateParams) (*clerk.PhoneNumber, error) {
	return getClient().Update(ctx, id, params)
}

// Delete deletes a phone number.
func Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	return getClient().Delete(ctx, id)
}

func getClient() *Client {
	return &Client{
		Backend: clerk.GetBackend(),
	}
}
