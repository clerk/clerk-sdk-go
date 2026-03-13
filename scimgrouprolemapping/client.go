// Package scimgrouprolemapping provides the SCIM Group Role Mappings API.
package scimgrouprolemapping

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/scim_directories"

// Client is used to invoke the SCIM Group Role Mappings API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// ListGroupsParams are the parameters for listing groups.
type ListGroupsParams struct {
	clerk.APIParams
	Cursor *string `json:"cursor,omitempty"`
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (p *ListGroupsParams) ToQuery() url.Values {
	q := url.Values{}
	if p != nil && p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	return q
}

// List returns all group role mappings for a SCIM directory.
func (c *Client) List(ctx context.Context, scimDirectoryID string) (*clerk.SCIMGroupRoleMappingList, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.SCIMGroupRoleMappingList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

// ListGroups returns SCIM groups for a directory with cursor pagination.
func (c *Client) ListGroups(ctx context.Context, scimDirectoryID string, params *ListGroupsParams) (*clerk.SCIMGroupList, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "groups")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.SCIMGroupList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

type CreateParams struct {
	clerk.APIParams
	SCIMGroupID string `json:"scim_group_id"`
	RoleID      string `json:"role_id"`
	Precedence  *int   `json:"precedence,omitempty"`
}

// Create creates a new group role mapping.
func (c *Client) Create(ctx context.Context, scimDirectoryID string, params *CreateParams) (*clerk.SCIMGroupRoleMapping, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.SCIMGroupRoleMapping{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type MappingUpdate struct {
	ID     string  `json:"id"`
	RoleID *string `json:"role_id,omitempty"`
}

type BulkUpdateParams struct {
	clerk.APIParams
	Mappings []MappingUpdate `json:"mappings"`
}

// BulkUpdate updates multiple group role mappings at once.
// The array position determines precedence (1-indexed).
// All mappings in the directory must be included.
func (c *Client) BulkUpdate(ctx context.Context, scimDirectoryID string, params *BulkUpdateParams) (*clerk.SCIMGroupRoleMappingList, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	list := &clerk.SCIMGroupRoleMappingList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

// Delete deletes a group role mapping.
func (c *Client) Delete(ctx context.Context, scimDirectoryID, mappingID string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings", mappingID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
