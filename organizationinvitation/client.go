// Package organizationinvitation provides the Organization Invitations API.
package organizationinvitation

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organization Invitations API.
type Client struct {
	Backend clerk.Backend
}

type ClientConfig struct {
	clerk.BackendConfig
}

func NewClient(config *ClientConfig) *Client {
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
