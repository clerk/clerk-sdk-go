// Package directorygrouprolemapping provides the Directory Group Role Mappings API.
//
// Directories are an experimental feature, not enabled for all instances.
package directorygrouprolemapping

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/directories"

// Client is used to invoke the Directory Group Role Mappings API.
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

// List returns all group role mappings for a directory.
func (c *Client) List(ctx context.Context, directoryID string) (*clerk.DirectoryGroupRoleMappingList, error) {
	path, err := clerk.JoinPath(path, directoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.DirectoryGroupRoleMappingList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

// ListGroups returns directory groups for a directory with cursor pagination.
func (c *Client) ListGroups(ctx context.Context, directoryID string, params *ListGroupsParams) (*clerk.DirectoryGroupList, error) {
	path, err := clerk.JoinPath(path, directoryID, "groups")
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
	DirectoryGroupID string `json:"directory_group_id"`
	RoleID           string `json:"role_id"`
	Precedence       *int   `json:"precedence,omitempty"`
}

// Create creates a new group role mapping.
func (c *Client) Create(ctx context.Context, directoryID string, params *CreateParams) (*clerk.DirectoryGroupRoleMapping, error) {
	path, err := clerk.JoinPath(path, directoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.DirectoryGroupRoleMapping{}
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
func (c *Client) BulkUpdate(ctx context.Context, directoryID string, params *BulkUpdateParams) (*clerk.DirectoryGroupRoleMappingList, error) {
	path, err := clerk.JoinPath(path, directoryID, "group_role_mappings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	list := &clerk.DirectoryGroupRoleMappingList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

// Delete deletes a group role mapping.
func (c *Client) Delete(ctx context.Context, directoryID, mappingID string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, directoryID, "group_role_mappings", mappingID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
