// Package roleset provides the Role Sets API.
package roleset

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/role_sets"

// Client is used to invoke the Role Sets API.
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
	// Type defines the type of role set. It can be either "initial" or "custom".
	Type *string `json:"type,omitempty"`
	// Roles are an array of role keys
	Roles          *[]string `json:"roles,omitempty"`
	DefaultRoleKey *string   `json:"default_role_key,omitempty"`
}

// Create creates a new role set.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.RoleSet, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	roleSet := &clerk.RoleSet{}
	err := c.Backend.Call(ctx, req, roleSet)
	return roleSet, err
}

// Get retrieves details for a role set.
func (c *Client) Get(ctx context.Context, roleSetKeyOrID string) (*clerk.RoleSet, error) {
	path, err := clerk.JoinPath(path, roleSetKeyOrID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	roleSet := &clerk.RoleSet{}
	err = c.Backend.Call(ctx, req, roleSet)
	return roleSet, err
}

type UpdateParams struct {
	clerk.APIParams
	Name        *string `json:"name,omitempty"`
	Key         *string `json:"key,omitempty"`
	Description *string `json:"description,omitempty"`
	// Type defines the type of role set. For update operations it can be only set to "initial".
	// There's only one "initial" role set per organization, after updating this roleset to "initial"
	// the other "initial" role sets will be updated to "custom".
	Type           *string `json:"type,omitempty"`
	DefaultRoleKey *string `json:"default_role_key,omitempty"`
}

// Update updates a role set.
func (c *Client) Update(ctx context.Context, roleSetKeyOrID string, params *UpdateParams) (*clerk.RoleSet, error) {
	path, err := clerk.JoinPath(path, roleSetKeyOrID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	roleSet := &clerk.RoleSet{}
	err = c.Backend.Call(ctx, req, roleSet)
	return roleSet, err
}

type DeleteParams struct {
	clerk.APIParams
	ReassignmentMappings *clerk.ReassignmentMappings `json:"reassignment_mappings,omitempty"`
	DestRoleSetKey       *string                     `json:"dest_role_set_key,omitempty"`
}

// Delete removes a role set.
func (c *Client) Delete(ctx context.Context, roleSetKeyOrID string, params *DeleteParams) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, roleSetKeyOrID, "/replace")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
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

// List returns a list of role sets.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.RoleSetList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.RoleSetList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type AddRolesParams struct {
	clerk.APIParams
	RoleKeys       []string `json:"role_keys,omitempty"`
	DefaultRoleKey *string  `json:"default_role_key,omitempty"`
}

// AddRoles adds roles to a role set.
func (c *Client) AddRoles(ctx context.Context, roleSetKeyOrID string, params *AddRolesParams) (*clerk.RoleSet, error) {
	path, err := clerk.JoinPath(path, roleSetKeyOrID, "/roles")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	roleSet := &clerk.RoleSet{}
	err = c.Backend.Call(ctx, req, roleSet)
	return roleSet, err
}

type RemoveRoleParams struct {
	clerk.APIParams
	RoleKey   string `json:"role_key,omitempty"`
	ToRoleKey string `json:"to_role_key,omitempty"`
}

// RemoveRoles removes roles from a role set.
func (c *Client) RemoveRole(ctx context.Context, roleSetKeyOrID string, params *RemoveRoleParams) (*clerk.RoleSet, error) {
	path, err := clerk.JoinPath(path, roleSetKeyOrID, "/roles", "/replace")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	roleSet := &clerk.RoleSet{}
	err = c.Backend.Call(ctx, req, roleSet)
	return roleSet, err
}
