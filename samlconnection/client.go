// Package samlconnection provides the SAML Connections API.
package samlconnection

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/saml_connections"

// Client is used to invoke the SAML Connections API.
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
	Name             *string                 `json:"name,omitempty"`
	Domain           *string                 `json:"domain,omitempty"`
	Provider         *string                 `json:"provider,omitempty"`
	IdpEntityID      *string                 `json:"idp_entity_id,omitempty"`
	IdpSsoURL        *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate   *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL   *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata      *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping *AttributeMappingParams `json:"attribute_mapping,omitempty"`
}

// Create creates a new SAML Connection.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.SAMLConnection, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	connection := &clerk.SAMLConnection{}
	err := c.Backend.Call(ctx, req, connection)
	return connection, err
}

// Get returns details about a SAML Connection.
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
	Name               *string                 `json:"name,omitempty"`
	Domain             *string                 `json:"domain,omitempty"`
	IdpEntityID        *string                 `json:"idp_entity_id,omitempty"`
	IdpSsoURL          *string                 `json:"idp_sso_url,omitempty"`
	IdpCertificate     *string                 `json:"idp_certificate,omitempty"`
	IdpMetadataURL     *string                 `json:"idp_metadata_url,omitempty"`
	IdpMetadata        *string                 `json:"idp_metadata,omitempty"`
	AttributeMapping   *AttributeMappingParams `json:"attribute_mapping,omitempty"`
	Active             *bool                   `json:"active,omitempty"`
	SyncUserAttributes *bool                   `json:"sync_user_attributes,omitempty"`
	AllowSubdomains    *bool                   `json:"allow_subdomains,omitempty"`
	AllowIdpInitiated  *bool                   `json:"allow_idp_initiated,omitempty"`
}

// Update updates the SAML Connection specified by id.
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
}

// List returns a list of SAML Connections.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.SAMLConnectionList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.SAMLConnectionList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
