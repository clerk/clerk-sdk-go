// Package scimgrouprolemapping provides the SCIM Group Role Mappings API.
package scimgrouprolemapping

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

// Client is used to invoke the SCIM Group Role Mappings API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

func path(scimDirectoryID string) string {
	return "/scim_directories/" + scimDirectoryID + "/group_role_mappings"
}

func groupsPath(scimDirectoryID string) string {
	return "/scim_directories/" + scimDirectoryID + "/groups"
}

// mappingList is a custom type to implement ResponseReader for array responses.
type mappingList []*clerk.SCIMGroupRoleMapping

func (*mappingList) Read(_ *clerk.APIResponse) {}

// groupList is a custom type to implement ResponseReader for array responses.
type groupList []*clerk.SCIMGroup

func (*groupList) Read(_ *clerk.APIResponse) {}

// List returns all group role mappings for a SCIM directory.
func (c *Client) List(ctx context.Context, scimDirectoryID string) ([]*clerk.SCIMGroupRoleMapping, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path(scimDirectoryID))
	data := &mappingList{}
	err := c.Backend.Call(ctx, req, data)
	if err != nil {
		return nil, err
	}
	return []*clerk.SCIMGroupRoleMapping(*data), nil
}

// ListGroups returns all SCIM groups for a directory.
func (c *Client) ListGroups(ctx context.Context, scimDirectoryID string) ([]*clerk.SCIMGroup, error) {
	req := clerk.NewAPIRequest(http.MethodGet, groupsPath(scimDirectoryID))
	data := &groupList{}
	err := c.Backend.Call(ctx, req, data)
	if err != nil {
		return nil, err
	}
	return []*clerk.SCIMGroup(*data), nil
}

type CreateParams struct {
	clerk.APIParams
	SCIMGroupID string `json:"scim_group_id"`
	RoleID      string `json:"role_id"`
	Precedence  *int   `json:"precedence,omitempty"`
}

// Create creates a new group role mapping.
func (c *Client) Create(ctx context.Context, scimDirectoryID string, params *CreateParams) (*clerk.SCIMGroupRoleMapping, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path(scimDirectoryID))
	req.SetParams(params)
	resource := &clerk.SCIMGroupRoleMapping{}
	err := c.Backend.Call(ctx, req, resource)
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
func (c *Client) BulkUpdate(ctx context.Context, scimDirectoryID string, params *BulkUpdateParams) ([]*clerk.SCIMGroupRoleMapping, error) {
	req := clerk.NewAPIRequest(http.MethodPatch, path(scimDirectoryID))
	req.SetParams(params)
	data := &mappingList{}
	err := c.Backend.Call(ctx, req, data)
	if err != nil {
		return nil, err
	}
	return []*clerk.SCIMGroupRoleMapping(*data), nil
}

// Delete deletes a group role mapping.
func (c *Client) Delete(ctx context.Context, scimDirectoryID, mappingID string) (*clerk.DeletedResource, error) {
	p, err := clerk.JoinPath(path(scimDirectoryID), mappingID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, p)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
