// Package enterpriseconnection provides the Enterprise Connections Backend API.
package enterpriseconnection

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/enterprise_connections"

// Client is used to invoke the Enterprise Connections API.
type Client struct {
	Backend clerk.Backend
}

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

// CreateParams holds parameters for creating an enterprise connection.
type CreateParams struct {
	clerk.APIParams
	// Common
	Name           *string   `json:"name,omitempty"`
	OrganizationID *string   `json:"organization_id,omitempty"`
	Protocol       *string   `json:"protocol,omitempty"` // "saml" or "oauth_oidc"
	Provider       *string   `json:"provider,omitempty"`
	Domains        *[]string `json:"domains,omitempty"`

	// SAML (protocol "saml"): https://clerk.com/docs/reference/backend-api/tag/enterprise-connections/post/enterprise_connections.body.saml
	IdpEntityID      *string                 `json:"idp_entity_id,omitempty"`
	IdpSsoURL        *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate   *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL   *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata      *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping *AttributeMappingParams `json:"attribute_mapping,omitempty"`
	ForceAuthn       *bool                   `json:"force_authn,omitempty"`

	// OIDC (protocol "oauth_oidc"): https://clerk.com/docs/reference/backend-api/tag/enterprise-connections/post/enterprise_connections.body.oidc
	ClientID         *string `json:"client_id,omitempty"`
	ClientSecret     *string `json:"client_secret,omitempty"`
	IssuerURL        *string `json:"issuer_url,omitempty"`
	AuthorizationURL *string `json:"authorization_url,omitempty"`
	TokenURL         *string `json:"token_url,omitempty"`
	UserInfoURL      *string `json:"user_info_url,omitempty"`
	Scopes           *string `json:"scopes,omitempty"`
}

// Create creates a new enterprise connection.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.EnterpriseConnection, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	conn := &clerk.EnterpriseConnection{}
	err := c.Backend.Call(ctx, req, conn)
	return conn, err
}

// Get returns an enterprise connection by ID.
func (c *Client) Get(ctx context.Context, id string) (*clerk.EnterpriseConnection, error) {
	p, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, p)
	conn := &clerk.EnterpriseConnection{}
	err = c.Backend.Call(ctx, req, conn)
	return conn, err
}

// UpdateParams holds parameters for updating an enterprise connection.
// Use SAML fields for SAML connections; use OIDC fields for OIDC connections.
type UpdateParams struct {
	clerk.APIParams
	// Common
	Name                             *string   `json:"name,omitempty"`
	OrganizationID                   *string   `json:"organization_id,omitempty"`
	Domains                          *[]string `json:"domains,omitempty"`
	Active                           *bool     `json:"active,omitempty"`
	SyncUserAttributes               *bool     `json:"sync_user_attributes,omitempty"`
	DisableAdditionalIdentifications *bool     `json:"disable_additional_identifications,omitempty"`

	// SAML: https://clerk.com/docs/reference/backend-api/tag/enterprise-connections/post/enterprise_connections.body.saml
	IdpEntityID                    *string                  `json:"idp_entity_id,omitempty"`
	IdpSsoURL                      *string                  `json:"idp_sso_url,omitempty"`
	IdpCertificate                 *string                  `json:"idp_certificate,omitempty"`
	IdpMetadataURL                 *string                  `json:"idp_metadata_url,omitempty"`
	IdpMetadata                    *string                  `json:"idp_metadata,omitempty"`
	AttributeMapping               *AttributeMappingParams  `json:"attribute_mapping,omitempty"`
	AllowSubdomains                *bool                    `json:"allow_subdomains,omitempty"`
	AllowIdpInitiated              *bool                    `json:"allow_idp_initiated,omitempty"`
	ForceAuthn                     *bool                    `json:"force_authn,omitempty"`
	ConsentVerifiedDomainsDeletion *bool                    `json:"consent_verified_domains_deletion,omitempty"`
	CustomAttributes               *[]clerk.CustomAttribute `json:"custom_attributes,omitempty"`

	// OIDC: https://clerk.com/docs/reference/backend-api/tag/enterprise-connections/post/enterprise_connections.body.oidc
	ClientID         *string `json:"client_id,omitempty"`
	ClientSecret     *string `json:"client_secret,omitempty"`
	IssuerURL        *string `json:"issuer_url,omitempty"`
	AuthorizationURL *string `json:"authorization_url,omitempty"`
	TokenURL         *string `json:"token_url,omitempty"`
	UserInfoURL      *string `json:"user_info_url,omitempty"`
	Scopes           *string `json:"scopes,omitempty"`
}

// Update updates an enterprise connection by ID.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.EnterpriseConnection, error) {
	p, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, p)
	req.SetParams(params)
	conn := &clerk.EnterpriseConnection{}
	err = c.Backend.Call(ctx, req, conn)
	return conn, err
}

// Delete deletes an enterprise connection.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	p, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, p)
	res := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, res)
	return res, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	// OrganizationID filters enterprise connections by organization ID.
	OrganizationID *string `json:"organization_id,omitempty"`
	// Active filters by active status. If true, only active connections are returned.
	// If false, only inactive connections are returned. If omitted, all connections are returned.
	Active *bool `json:"active,omitempty"`
}

// ToQuery returns query string values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.OrganizationID != nil {
		q.Set("organization_id", *params.OrganizationID)
	}
	if params.Active != nil {
		q.Set("active", strconv.FormatBool(*params.Active))
	}
	return q
}

// List returns a list of enterprise connections.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.EnterpriseConnectionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.EnterpriseConnectionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
