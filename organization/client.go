// Package organization provides the Organizations API.
package organization

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organizations API.
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
	Name                  *string          `json:"name,omitempty"`
	Slug                  *string          `json:"slug,omitempty"`
	CreatedBy             *string          `json:"created_by,omitempty"`
	MaxAllowedMemberships *int64           `json:"max_allowed_memberships,omitempty"`
	PublicMetadata        *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata       *json.RawMessage `json:"private_metadata,omitempty"`
}

// Create creates a new organization.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.Organization, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	organization := &clerk.Organization{}
	err := c.Backend.Call(ctx, req, organization)
	return organization, err
}

// Get retrieves details for an organization.
// The organization can be fetched by either the ID or its slug.
func (c *Client) Get(ctx context.Context, idOrSlug string) (*clerk.Organization, error) {
	path, err := clerk.JoinPath(path, idOrSlug)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	organization := &clerk.Organization{}
	err = c.Backend.Call(ctx, req, organization)
	return organization, err
}

type UpdateParams struct {
	clerk.APIParams
	Name                  *string          `json:"name,omitempty"`
	Slug                  *string          `json:"slug,omitempty"`
	MaxAllowedMemberships *int64           `json:"max_allowed_memberships,omitempty"`
	PublicMetadata        *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata       *json.RawMessage `json:"private_metadata,omitempty"`
	AdminDeleteEnabled    *bool            `json:"admin_delete_enabled,omitempty"`
}

// Update updates an organization.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.Organization, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	organization := &clerk.Organization{}
	err = c.Backend.Call(ctx, req, organization)
	return organization, err
}

type UpdateMetadataParams struct {
	clerk.APIParams
	PublicMetadata  *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata *json.RawMessage `json:"private_metadata,omitempty"`
}

// UpdateMetadata updates the organization's metadata by merging the
// provided values with the existing ones.
func (c *Client) UpdateMetadata(ctx context.Context, id string, params *UpdateMetadataParams) (*clerk.Organization, error) {
	path, err := clerk.JoinPath(path, id, "/metadata")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	organization := &clerk.Organization{}
	err = c.Backend.Call(ctx, req, organization)
	return organization, err
}

// Delete deletes an organization.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	organization := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, organization)
	return organization, err
}

// DeleteLogo removes the organization's logo.
func (c *Client) DeleteLogo(ctx context.Context, id string) (*clerk.Organization, error) {
	path, err := clerk.JoinPath(path, id, "/logo")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	organization := &clerk.Organization{}
	err = c.Backend.Call(ctx, req, organization)
	return organization, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	IncludeMembersCount *bool    `json:"include_members_count,omitempty"`
	OrderBy             *string  `json:"order_by,omitempty"`
	Query               *string  `json:"query,omitempty"`
	UserIDs             []string `json:"user_id,omitempty"`
}

// ToQuery returns query string values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.IncludeMembersCount != nil {
		q.Set("include_members_count", strconv.FormatBool(*params.IncludeMembersCount))
	}
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	if params.UserIDs != nil {
		q["user_id"] = params.UserIDs
	}
	return q
}

// List returns a list of organizations.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OrganizationList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OrganizationList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
