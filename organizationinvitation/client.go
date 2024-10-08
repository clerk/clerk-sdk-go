// Package organizationinvitation provides the Organization Invitations API.
package organizationinvitation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organization Invitations API.
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
	EmailAddress    *string          `json:"email_address,omitempty"`
	Role            *string          `json:"role,omitempty"`
	RedirectURL     *string          `json:"redirect_url,omitempty"`
	InviterUserID   *string          `json:"inviter_user_id,omitempty"`
	PublicMetadata  *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata *json.RawMessage `json:"private_metadata,omitempty"`
	OrganizationID  string           `json:"-"`
}

// Create creates and sends an invitation to join an organization.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.OrganizationInvitation, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/invitations")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	invitation := &clerk.OrganizationInvitation{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrganizationID string
	Statuses       *[]string
}

func (p *ListParams) ToQuery() url.Values {
	q := p.ListParams.ToQuery()

	if p.Statuses != nil && len(*p.Statuses) > 0 {
		q["status"] = *p.Statuses
	}

	return q
}

// List returns a list of organization invitations
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OrganizationInvitationList, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/invitations")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	invitation := &clerk.OrganizationInvitationList{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

type GetParams struct {
	OrganizationID string
	ID             string
}

// Get retrieves the detail for an organization invitation.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.OrganizationInvitation, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/invitations", params.ID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	invitation := &clerk.OrganizationInvitation{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

type RevokeParams struct {
	clerk.APIParams
	RequestingUserID *string `json:"requesting_user_id,omitempty"`
	OrganizationID   string  `json:"-"`
	ID               string  `json:"-"`
}

// Revoke marks the organization invitation as revoked.
func (c *Client) Revoke(ctx context.Context, params *RevokeParams) (*clerk.OrganizationInvitation, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/invitations", params.ID, "/revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	invitation := &clerk.OrganizationInvitation{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

type ListFromInstanceParams struct {
	clerk.APIParams
	clerk.ListParams
	Statuses *[]string `json:"statuses,omitempty"`
	Query    *string   `json:"query,omitempty"`
	OrderBy  *string   `json:"order_by,omitempty"`
}

func (p *ListFromInstanceParams) ToQuery() url.Values {
	q := p.ListParams.ToQuery()

	if p.Statuses != nil && len(*p.Statuses) > 0 {
		q["status"] = *p.Statuses
	}

	if p.Query != nil {
		q.Add("query", *p.Query)
	}

	if p.OrderBy != nil {
		q.Add("order_by", *p.OrderBy)
	}

	return q
}

// ListAllFromInstance lists all the organization invitations from the current instance
func (c *Client) ListFromInstance(ctx context.Context, params *ListFromInstanceParams) (*clerk.OrganizationInvitationList, error) {
	path, err := clerk.JoinPath("/organization_invitations")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	invitation := &clerk.OrganizationInvitationList{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}
