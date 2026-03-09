// Package enterpriseconnection provides the Enterprise Connections Backend API.
package enterpriseconnection

import (
	"context"
	"net/http"
	"net/url"

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

type CreateParams struct {
	clerk.APIParams
	Name           *string `json:"name,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Protocol       *string `json:"protocol,omitempty"`
	Provider       *string `json:"provider,omitempty"`
	Domains        *[]string               `json:"domains,omitempty"`
	IdpEntityID   *string                 `json:"idp_entity_id,omitempty"`
	IdpSsoURL     *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate   *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL   *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata      *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping *AttributeMappingParams `json:"attribute_mapping,omitempty"`
	ForceAuthn       *bool                   `json:"force_authn,omitempty"`
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

type UpdateParams struct {
	clerk.APIParams
	Name                             *string                 `json:"name,omitempty"`
	OrganizationID                    *string                 `json:"organization_id,omitempty"`
	Domains                          *[]string               `json:"domains,omitempty"`
	IdpEntityID                      *string                 `json:"idp_entity_id,omitempty"`
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
	ForceAuthn                       *bool                   `json:"force_authn,omitempty"`
	CustomAttributes                 *[]clerk.CustomAttribute `json:"custom_attributes,omitempty"`
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

// List returns a list of enterprise connections.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.EnterpriseConnectionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.EnterpriseConnectionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
