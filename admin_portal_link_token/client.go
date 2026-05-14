// Package admin_portal_link_token provides the Admin Portal Link Tokens API.
//
// Admin Portal Link Tokens are single-use, URL-borne credentials issued
// against an origin-level admin portal link. They bootstrap an IT contact
// into a Clerk-hosted admin portal session (e.g. for SSO setup), and sit
// alongside api_keys and m2m_tokens as a peer opaque-token primitive.
package admin_portal_link_token

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/admin_portal_link_tokens"

// Client is used to invoke the Admin Portal Link Tokens API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// CreateParams configures a new admin portal link token. All fields are
// optional — the parent admin_portal_link id is minted server-side and
// returned on the response.
type CreateParams struct {
	clerk.APIParams
	// OrganizationID scopes the token to an existing organization. Omit for
	// flows that bootstrap before any org exists (e.g. first-time setup).
	OrganizationID *string `json:"organization_id,omitempty"`
	// ITContactID is an opaque reference to the IT contact this link is
	// associated with. Surfaced on the token's claims.
	ITContactID *string `json:"it_contact_id,omitempty"`
	// Scopes attached to the token. The verify path (admin portal frontend)
	// decides what each scope grants.
	Scopes []string `json:"scopes,omitempty"`
	// SecondsUntilExpiration overrides the default 1h TTL. Capped at 24h server-side.
	SecondsUntilExpiration *int64 `json:"seconds_until_expiration,omitempty"`
}

// Create mints a new admin portal link token. The returned token is shown
// once and must be embedded in the URL delivered to the IT contact.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.AdminPortalLinkTokenWithToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.AdminPortalLinkTokenWithToken{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// RevokeParams configures a revoke request.
type RevokeParams struct {
	clerk.APIParams
	RevocationReason *string `json:"revocation_reason,omitempty"`
}

// Revoke revokes an admin portal link token.
func (c *Client) Revoke(ctx context.Context, adminPortalLinkTokenID string, params *RevokeParams) (*clerk.AdminPortalLinkToken, error) {
	path, err := clerk.JoinPath(path, adminPortalLinkTokenID, "revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.AdminPortalLinkToken{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
