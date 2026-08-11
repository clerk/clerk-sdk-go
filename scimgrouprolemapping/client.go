// Package scimgrouprolemapping provides the SCIM Group Role Mappings API.
//
// Deprecated: use package directorygrouprolemapping. This package is retained
// so that code written against the SCIM-era names keeps compiling. It
// continues to call the legacy /scim_directories routes, which stay mounted
// alongside /directories and return the same resources.
package scimgrouprolemapping

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

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
	StartingAfter *string `json:"starting_after,omitempty"`
	EndingBefore  *string `json:"ending_before,omitempty"`
	Limit         *int    `json:"limit,omitempty"`
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (p *ListGroupsParams) ToQuery() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.StartingAfter != nil {
		q.Set("starting_after", *p.StartingAfter)
	}
	if p.EndingBefore != nil {
		q.Set("ending_before", *p.EndingBefore)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// List returns all group role mappings for a SCIM directory.
func (c *Client) List(ctx context.Context, scimDirectoryID string) (*clerk.SCIMGroupRoleMappingList, error) { //nolint:staticcheck // SA1019: this deprecated package returns the SCIM-named list on purpose. Its elements are SCIMGroupRoleMapping, which exposes the legacy scim_* field names that existing callers read.
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.SCIMGroupRoleMappingList{} //nolint:staticcheck // SA1019: must match this method's deprecated return type.
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

// ListGroups returns SCIM groups for a directory with cursor pagination.
func (c *Client) ListGroups(ctx context.Context, scimDirectoryID string, params *ListGroupsParams) (*clerk.DirectoryGroupList, error) {
	path, err := clerk.JoinPath(path, scimDirectoryID, "groups")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.DirectoryGroupList{}
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
func (c *Client) Create(ctx context.Context, scimDirectoryID string, params *CreateParams) (*clerk.SCIMGroupRoleMapping, error) { //nolint:staticcheck // SA1019: this deprecated package returns the SCIM-named type on purpose. It decodes the legacy scim_* JSON fields that existing callers read, so swapping the type here would break them.
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.SCIMGroupRoleMapping{} //nolint:staticcheck // SA1019: must match this method's deprecated return type.
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
func (c *Client) BulkUpdate(ctx context.Context, scimDirectoryID string, params *BulkUpdateParams) (*clerk.SCIMGroupRoleMappingList, error) { //nolint:staticcheck // SA1019: this deprecated package returns the SCIM-named list on purpose. Its elements are SCIMGroupRoleMapping, which exposes the legacy scim_* field names that existing callers read.
	path, err := clerk.JoinPath(path, scimDirectoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	list := &clerk.SCIMGroupRoleMappingList{} //nolint:staticcheck // SA1019: must match this method's deprecated return type.
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
