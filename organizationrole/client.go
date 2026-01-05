// Package organizationrole provides the Organization Roles API.
package organizationrole

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organization_roles"

// Client is used to invoke the Organization Roles API.
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
	Name                    *string   `json:"name,omitempty"`
	Key                     *string   `json:"key,omitempty"`
	Description             *string   `json:"description,omitempty"`
	Permissions             *[]string `json:"permissions,omitempty"`
	IsManaged               *bool     `json:"is_managed,omitempty"`
	IncludeInInitialRoleSet *bool     `json:"include_in_initial_role_set,omitempty"`
}

// Create creates a new organization role.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.OrganizationRole, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	role := &clerk.OrganizationRole{}
	err := c.Backend.Call(ctx, req, role)
	return role, err
}

// Get retrieves details for an organization role.
func (c *Client) Get(ctx context.Context, roleID string) (*clerk.OrganizationRole, error) {
	path, err := clerk.JoinPath(path, roleID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	role := &clerk.OrganizationRole{}
	err = c.Backend.Call(ctx, req, role)
	return role, err
}

type UpdateParams struct {
	clerk.APIParams
	Name        *string   `json:"name,omitempty"`
	Key         *string   `json:"key,omitempty"`
	Description *string   `json:"description,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
}

// Update updates an organization role.
func (c *Client) Update(ctx context.Context, roleID string, params *UpdateParams) (*clerk.OrganizationRole, error) {
	path, err := clerk.JoinPath(path, roleID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	role := &clerk.OrganizationRole{}
	err = c.Backend.Call(ctx, req, role)
	return role, err
}

// Delete removes an organization role.
func (c *Client) Delete(ctx context.Context, roleID string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, roleID)
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

// List returns a list of organization roles.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OrganizationRoleList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OrganizationRoleList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

// AssignPermission assigns a permission to an organization role.
func (c *Client) AssignPermission(ctx context.Context, roleID, permissionID string) (*clerk.OrganizationRole, error) {
	path, err := clerk.JoinPath(path, roleID, "/permissions", permissionID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	role := &clerk.OrganizationRole{}
	err = c.Backend.Call(ctx, req, role)
	return role, err
}

// RemovePermission removes a permission from an organization role.
func (c *Client) RemovePermission(ctx context.Context, roleID, permissionID string) (*clerk.OrganizationRole, error) {
	path, err := clerk.JoinPath(path, roleID, "/permissions", permissionID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	role := &clerk.OrganizationRole{}
	err = c.Backend.Call(ctx, req, role)
	return role, err
}
