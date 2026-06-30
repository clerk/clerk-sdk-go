// Package samlconnection provides the SAML Connections API.
//
// Deprecated: Use the enterpriseconnection package and the Enterprise Connections
// Backend API (https://clerk.com/docs/reference/backend-api/tag/enterprise-connections) instead.
package samlconnection

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/saml_connections"

// Client is used to invoke the SAML Connections API.
//
// Deprecated: Use enterpriseconnection.Client instead.
type Client struct {
	Backend clerk.Backend
}

// NewClient returns a Client for the SAML Connections API.
//
// Deprecated: Use enterpriseconnection.NewClient instead.
func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type AttributeMappingParams struct {
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type CreateParams struct {
	clerk.APIParams
	Name           *string `json:"name,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	// Deprecated: Use `domains` instead.
	Domain                 *string                 `json:"domain,omitempty"`
	Domains                *[]string               `json:"domains,omitempty"`
	Provider               *string                 `json:"provider,omitempty"`
	IdpEntityID            *string                 `json:"idp_entity_id,omitempty"`
	IdpSsoURL              *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate         *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL         *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata            *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping       *AttributeMappingParams `json:"attribute_mapping,omitempty"`
	ForceAuthn             *bool                   `json:"force_authn,omitempty"`
	DisableJITProvisioning *bool                   `json:"disable_jit_provisioning,omitempty"`
	// CustomAttributes is an Experimental feature, not available for all customers.
	CustomAttributes *[]clerk.CustomAttribute `json:"custom_attributes,omitempty"`
}

// Create creates a new SAML Connection.
//
// Deprecated: Use enterpriseconnection.Create instead.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.SAMLConnection, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	connection := &clerk.SAMLConnection{}
	err := c.Backend.Call(ctx, req, connection)
	return connection, err
}

// Get returns details about a SAML Connection.
//
// Deprecated: Use enterpriseconnection.Get instead.
func (c *Client) Get(ctx context.Context, id string) (*clerk.SAMLConnection, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	connection := &clerk.SAMLConnection{}
	err = c.Backend.Call(ctx, req, connection)
	return connection, err
}

type UpdateParams struct {
	clerk.APIParams
	Name *string `json:"name,omitempty"`
	// Deprecated: Use `domains` instead.
	Domain *string `json:"domain,omitempty"`
	// Domains represents the complete, desired set of domains.
	//   - If nil or an empty slice is provided, no changes will be made.
	//   - Otherwise, the provided slice is treated as the full target list:
	//     • Any existing domains not in this list will be removed.
	//     • Any domains in this list not already present will be added.
	//
	// For example, if the current domains are ["b", "c"] and you set:
	//     Domains: []string{"a", "c", "d"}
	// then:
	//     - "a" and "d" will be added
	//     - "b" will be removed
	//     - "c" will remain unchanged
	Domains     *[]string `json:"domains,omitempty"`
	IdpEntityID *string   `json:"idp_entity_id,omitempty"`
	// OrganizationID is a nullable optional field.
	//   - If nil or unset, no action will be taken.
	//   - If an empty value (""), the organization_id will be unset.
	//   - If a valid ID is provided, the organization_id will be updated.
	OrganizationID                   *string                 `json:"organization_id,omitempty"`
	IdpSsoURL                        *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate                   *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL                   *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata                      *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping                 *AttributeMappingParams `json:"attribute_mapping,omitempty"`
	Active                           *bool                   `json:"active,omitempty"`
	SyncUserAttributes               *bool                   `json:"sync_user_attributes,omitempty"`
	AllowSubdomains                  *bool                   `json:"allow_subdomains,omitempty"`
	AllowIdpInitiated                *bool                   `json:"allow_idp_initiated,omitempty"`
	DisableAdditionalIdentifications *bool                   `json:"disable_additional_identifications,omitempty"`
	AllowOrganizationAccountLinking  *bool                   `json:"allow_organization_account_linking,omitempty"`
	ForceAuthn                       *bool                   `json:"force_authn,omitempty"`
	DisableJITProvisioning           *bool                   `json:"disable_jit_provisioning,omitempty"`
	ConsentVerifiedDomainsDeletion   *bool                   `json:"consent_verified_domains_deletion,omitempty"`
	Authenticatable                  *bool                   `json:"authenticatable,omitempty"`
	// CustomAttributes is an Experimental feature, not available for all customers.
	CustomAttributes *[]clerk.CustomAttribute `json:"custom_attributes,omitempty"`
}

// Update updates the SAML Connection specified by id.
//
// Deprecated: Use enterpriseconnection.Update instead.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.SAMLConnection, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	connection := &clerk.SAMLConnection{}
	err = c.Backend.Call(ctx, req, connection)
	return connection, err
}

// Delete deletes a SAML Connection.
//
// Deprecated: Use enterpriseconnection.Delete instead.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	template := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, template)
	return template, err
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

// List returns a list of SAML Connections.
//
// Deprecated: Use enterpriseconnection.List instead.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.SAMLConnectionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.SAMLConnectionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
