// Package organizationdomain provides the Organization Domains API.
package organizationdomain

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organization Domains API.
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
	Name           *string `json:"name,omitempty"`
	EnrollmentMode *string `json:"enrollment_mode,omitempty"`
	Verified       *bool   `json:"verified,omitempty"`
}

// Create adds a new domain to the organization.
func (c *Client) Create(ctx context.Context, organizationID string, params *CreateParams) (*clerk.OrganizationDomain, error) {
	path, err := clerk.JoinPath(path, organizationID, "/domains")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	domain := &clerk.OrganizationDomain{}
	err = c.Backend.Call(ctx, req, domain)
	return domain, err
}

type UpdateParams struct {
	clerk.APIParams
	OrganizationID string  `json:"-"`
	DomainID       string  `json:"-"`
	EnrollmentMode *string `json:"enrollment_mode,omitempty"`
	Verified       *bool   `json:"verified,omitempty"`
}

// Update updates an organization domain.
func (c *Client) Update(ctx context.Context, params *UpdateParams) (*clerk.OrganizationDomain, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/domains", params.DomainID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	domain := &clerk.OrganizationDomain{}
	err = c.Backend.Call(ctx, req, domain)
	return domain, err
}

type DeleteParams struct {
	clerk.APIParams
	OrganizationID string `json:"-"`
	DomainID       string `json:"-"`
}

// Delete removes a domain from an organization.
func (c *Client) Delete(ctx context.Context, params *DeleteParams) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/domains", params.DomainID)
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
	Verified        *bool     `json:"verified,omitempty"`
	EnrollmentModes *[]string `json:"enrollment_mode,omitempty"`
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()

	if params.Verified != nil {
		q.Set("verified", strconv.FormatBool(*params.Verified))
	}

	if params.EnrollmentModes != nil && len(*params.EnrollmentModes) > 0 {
		q["enrollment_mode"] = *params.EnrollmentModes
	}

	return q
}

// List returns a list of organization domains.
func (c *Client) List(ctx context.Context, organizationID string, params *ListParams) (*clerk.OrganizationDomainList, error) {
	path, err := clerk.JoinPath(path, organizationID, "/domains")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	domains := &clerk.OrganizationDomainList{}
	err = c.Backend.Call(ctx, req, domains)
	return domains, err
}

type ListAllFromInstanceParams struct {
	clerk.APIParams
	clerk.ListParams
	Verified       *bool     `json:"verified,omitempty"`
	EnrollmentMode *string   `json:"enrollment_mode,omitempty"`
	OrganizationID *string   `json:"organization_id,omitempty"`
	OrderBy        *string   `json:"order_by,omitempty"`
	Query          *string   `json:"query,omitempty"`
	Domains        *[]string `json:"domains,omitempty"`
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (params *ListAllFromInstanceParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()

	if params.Verified != nil {
		q.Set("verified", strconv.FormatBool(*params.Verified))
	}

	if params.EnrollmentMode != nil {
		q.Set("enrollment_mode", *params.EnrollmentMode)
	}

	if params.OrganizationID != nil {
		q.Set("organization_id", *params.OrganizationID)
	}

	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}

	if params.Query != nil {
		q.Set("query", *params.Query)
	}

	if params.Domains != nil && len(*params.Domains) > 0 {
		q["domains"] = *params.Domains
	}

	return q
}

// ListFromInstance returns a list of organization domains from all organizations in the instance.
func (c *Client) ListAllFromInstance(ctx context.Context, params *ListAllFromInstanceParams) (*clerk.OrganizationDomainList, error) {
	path, err := clerk.JoinPath("/organization_domains")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	domains := &clerk.OrganizationDomainList{}
	err = c.Backend.Call(ctx, req, domains)
	return domains, err
}
