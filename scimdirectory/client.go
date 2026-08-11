// Package scimdirectory provides the SCIM Directories API.
// SCIM directories are an experimental features, not enabled for all instances
//
// Deprecated: use package directory. This package is retained so that code
// written against the SCIM-era names keeps compiling. It continues to call the
// legacy /scim_directories routes, which stay mounted alongside /directories
// and return the same resources.
package scimdirectory

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/scim_directories"

// Client is used to invoke the SCIM Directories API.
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
	EnterpriseConnectionID *string `json:"enterprise_connection_id,omitempty"`
	Name                   *string `json:"name,omitempty"`
	Provider               *string `json:"provider,omitempty"`
}

// Create creates a new SCIM directory.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.Directory, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.Directory{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Get returns details about a SCIM directory.
func (c *Client) Get(ctx context.Context, id string) (*clerk.Directory, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.Directory{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type UpdateParams struct {
	clerk.APIParams
	Name                    *string `json:"name,omitempty"`
	Provider                *string `json:"provider,omitempty"`
	Enabled                 *bool   `json:"enabled,omitempty"`
	GroupRoleMappingEnabled *bool   `json:"group_role_mapping_enabled,omitempty"`
	// AttributeMapping is a map of SCIM attributes to Clerk attributes.
	// The semantics of the PATCH request are as follows:
	//   - If the attribute is not present in the request, it will be left unchanged.
	//   - If the attribute is present in the request, it will be updated to the new value.
	//   - If the attribute is present in the request and the value is null, it will be removed.
	AttributeMapping *map[string]string `json:"attribute_mapping,omitempty"`
}

// Update updates a SCIM directory.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.Directory, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	resource := &clerk.Directory{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
}

func (params *ListParams) ToQuery() url.Values {
	q := url.Values{}
	if params.Limit != nil {
		q.Set("limit", strconv.FormatInt(*params.Limit, 10))
	}
	if params.Offset != nil {
		q.Set("offset", strconv.FormatInt(*params.Offset, 10))
	}
	return q
}

// List returns a paginated list of SCIM directories.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.SCIMDirectoryList, error) { //nolint:staticcheck // SA1019: this deprecated package returns the SCIM-named list on purpose. DirectoryList renames the data field from SCIMDirectories to Directories, so swapping the type here would break existing callers.
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &clerk.SCIMDirectoryList{} //nolint:staticcheck // SA1019: must match this method's deprecated return type.
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Delete deletes a SCIM directory.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// RotateAPIKey rotates the API key for a SCIM directory.
// The old API key will be valid for a grace period before expiring.
func (c *Client) RotateAPIKey(ctx context.Context, id string) (*clerk.Directory, error) {
	path, err := clerk.JoinPath(path, id, "rotate_api_key")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.Directory{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Sync enqueues a sync run for a pull-based SCIM directory. The server
// responds 202 Accepted; the run itself is asynchronous.
func (c *Client) Sync(ctx context.Context, id string) error {
	path, err := clerk.JoinPath(path, id, "sync")
	if err != nil {
		return err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	return c.Backend.Call(ctx, req, &clerk.APIResource{})
}

type CredentialsParams struct {
	clerk.APIParams
	ServiceAccountJSON *string `json:"service_account_json,omitempty"`
	SubjectEmail       *string `json:"subject_email,omitempty"`
}

// AddCredentials adds pull-provider credentials (e.g. a Google service-account
// key and subject email) for a SCIM directory. They are validated and sealed
// server-side; on success the directory is enabled.
func (c *Client) AddCredentials(ctx context.Context, id string, params *CredentialsParams) (*clerk.Directory, error) {
	path, err := clerk.JoinPath(path, id, "credentials")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.Directory{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
