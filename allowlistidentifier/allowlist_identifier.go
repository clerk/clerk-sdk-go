// Package allowlistidentifier provides the Allowlist Identifiers API.
package allowlistidentifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

const path = "/allowlist_identifiers"

type CreateParams struct {
	clerk.APIParams
	Identifier *string `json:"identifier,omitempty"`
	Notify     *bool   `json:"notify,omitempty"`
}

// Create adds a new identifier to the allowlist.
func Create(ctx context.Context, params *CreateParams) (*clerk.AllowlistIdentifier, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	identifier := &clerk.AllowlistIdentifier{}
	err := clerk.GetBackend().Call(ctx, req, identifier)
	return identifier, err
}

// Delete removes an identifier from the allowlist.
func Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := url.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	identifier := &clerk.DeletedResource{}
	err = clerk.GetBackend().Call(ctx, req, identifier)
	return identifier, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns all the identifiers in the allowlist.
func List(ctx context.Context, params *ListParams) (*clerk.AllowlistIdentifierList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	list := &clerk.AllowlistIdentifierList{}
	err := clerk.GetBackend().Call(ctx, req, list)
	return list, err
}
