// Package organizationaccesspolicy provides the Organization Access Policies API.
package organizationaccesspolicy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organization Access Policies API.
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
	EnrollmentMode *string               `json:"enrollment_mode,omitempty"`
	Authentication *AuthenticationParams `json:"authentication,omitempty"`
	// Status is taken by sso policies only, and defaults to draft.
	Status *string `json:"status,omitempty"`
}

// Create adds a new access policy to the organization.
func (c *Client) Create(ctx context.Context, organizationID string, params *CreateParams) (*clerk.OrganizationAccessPolicy, error) {
	path, err := clerk.JoinPath(path, organizationID, "/access_policies")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	policy := &clerk.OrganizationAccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}

type GetParams struct {
	OrganizationID string
	PolicyID       string
}

// Get retrieves an organization access policy.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.OrganizationAccessPolicy, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/access_policies", params.PolicyID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	policy := &clerk.OrganizationAccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}

type UpdateParams struct {
	clerk.APIParams
	OrganizationID string `json:"-"`
	PolicyID       string `json:"-"`
	// Domains replaces the policy's domains when set. An empty list clears
	// them, which is only valid on the default policy.
	Domains        *[]string             `json:"domains,omitempty"`
	EnrollmentMode *string               `json:"enrollment_mode,omitempty"`
	Authentication *AuthenticationParams `json:"authentication,omitempty"`
	Status         *string               `json:"status,omitempty"`
}

// Update updates an organization access policy.
func (c *Client) Update(ctx context.Context, params *UpdateParams) (*clerk.OrganizationAccessPolicy, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/access_policies", params.PolicyID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	policy := &clerk.OrganizationAccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}

type DeleteParams struct {
	clerk.APIParams
	OrganizationID string `json:"-"`
	PolicyID       string `json:"-"`
}

// Delete removes an access policy from an organization.
func (c *Client) Delete(ctx context.Context, params *DeleteParams) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/access_policies", params.PolicyID)
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

// List returns a list of the organization's access policies.
func (c *Client) List(ctx context.Context, organizationID string, params *ListParams) (*clerk.OrganizationAccessPolicyList, error) {
	path, err := clerk.JoinPath(path, organizationID, "/access_policies")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	policies := &clerk.OrganizationAccessPolicyList{}
	err = c.Backend.Call(ctx, req, policies)
	return policies, err
}

type VerifyParams struct {
	clerk.APIParams
	OrganizationID string `json:"-"`
	PolicyID       string `json:"-"`
	// Domain is one of the policy's domains.
	Domain *string `json:"domain,omitempty"`
	// Kind is the proof to record over the domain, affiliation or ownership.
	Kind *string `json:"kind,omitempty"`
}

// Verify marks one of the access policy's domains as verified.
func (c *Client) Verify(ctx context.Context, params *VerifyParams) (*clerk.OrganizationAccessPolicy, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/access_policies", params.PolicyID, "/verify")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	policy := &clerk.OrganizationAccessPolicy{}
	err = c.Backend.Call(ctx, req, policy)
	return policy, err
}
