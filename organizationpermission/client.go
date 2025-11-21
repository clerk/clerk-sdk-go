// Package organizationpermission provides the Organization Permissions API.
package organizationpermission

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organization_permissions"

// Client is used to invoke the Organization Permissions API.
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
	Name        *string `json:"name,omitempty"`
	Key         *string `json:"key,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Create creates a new organization permission.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.OrganizationPermission, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	permission := &clerk.OrganizationPermission{}
	err := c.Backend.Call(ctx, req, permission)
	return permission, err
}

// Get retrieves details for an organization permission.
func (c *Client) Get(ctx context.Context, permissionID string) (*clerk.OrganizationPermission, error) {
	path, err := clerk.JoinPath(path, permissionID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	permission := &clerk.OrganizationPermission{}
	err = c.Backend.Call(ctx, req, permission)
	return permission, err
}

type UpdateParams struct {
	clerk.APIParams
	Name        *string `json:"name,omitempty"`
	Key         *string `json:"key,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Update updates an organization permission.
func (c *Client) Update(ctx context.Context, permissionID string, params *UpdateParams) (*clerk.OrganizationPermission, error) {
	path, err := clerk.JoinPath(path, permissionID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	permission := &clerk.OrganizationPermission{}
	err = c.Backend.Call(ctx, req, permission)
	return permission, err
}

// Delete removes an organization permission.
func (c *Client) Delete(ctx context.Context, permissionID string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, permissionID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	deletedResource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, deletedResource)
	return deletedResource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams

	Query   *string `json:"query,omitempty"`
	OrderBy *string `json:"order_by,omitempty"`
}

// ToQuery returns query string values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	return q
}

// List returns a list of organization permissions.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OrganizationPermissionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OrganizationPermissionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
