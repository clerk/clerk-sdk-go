// Package accesspolicy provides the Access Policies API.
package accesspolicy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/access_policies"

// Client is used to invoke the Access Policies API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type AuthenticationParams struct {
	Connections        *[]string `json:"connections,omitempty"`
	Directories        *[]string `json:"directories,omitempty"`
	Provisioning       *string   `json:"provisioning,omitempty"`
	ConnectionRequired *bool     `json:"connection_required,omitempty"`
}

type CreateParams struct {
	clerk.APIParams
	Domains        *[]string             `json:"domains,omitempty"`
	Authentication *AuthenticationParams `json:"authentication,omitempty"`
	// Status defaults to draft.
	Status *string `json:"status,omitempty"`
}

// Create adds a new access policy to the instance.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.AccessPolicy, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	policy := &clerk.AccessPolicy{}
	err := c.Backend.Call(ctx, req, policy)
	return policy, err
}

// Get retrieves an instance access policy.
func (c *Client) Get(ctx context.Context, id string) (*clerk.AccessPolicy, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	policy := &clerk.AccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}

type UpdateParams struct {
	clerk.APIParams
	// Domains replaces the policy's domains when set.
	Domains        *[]string             `json:"domains,omitempty"`
	Authentication *AuthenticationParams `json:"authentication,omitempty"`
	Status         *string               `json:"status,omitempty"`
}

// Update updates an instance access policy.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.AccessPolicy, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	policy := &clerk.AccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}

// Delete removes an access policy from the instance.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	res := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, res)
	return res, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (params *ListParams) ToQuery() url.Values {
	return params.ListParams.ToQuery()
}

// List returns a list of the instance's access policies.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.AccessPolicyList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	policies := &clerk.AccessPolicyList{}
	err := c.Backend.Call(ctx, req, policies)
	return policies, err
}
