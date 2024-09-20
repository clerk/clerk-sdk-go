package invitation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/invitations"

// Client is used to invoke the Invitations API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrderBy  *string  `json:"order_by,omitempty"`
	Query    *string  `json:"query,omitempty"`
	Statuses []string `json:"status,omitempty"`
}

// ToQuery returns query string values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	for _, status := range params.Statuses {
		q.Add("status", status)
	}
	return q
}

// List returns all invitations.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.InvitationList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, fmt.Sprintf("%s?paginated=true", path))
	req.SetParams(params)
	list := &clerk.InvitationList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type CreateParams struct {
	clerk.APIParams
	EmailAddress   string           `json:"email_address"`
	PublicMetadata *json.RawMessage `json:"public_metadata,omitempty"`
	RedirectURL    *string          `json:"redirect_url,omitempty"`
	Notify         *bool            `json:"notify,omitempty"`
	IgnoreExisting *bool            `json:"ignore_existing,omitempty"`
	ExpiresInDays  *int64           `json:"expires_in_days,omitempty"`
}

// Create adds a new identifier to the allowlist.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.Invitation, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	invitation := &clerk.Invitation{}
	err := c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

// Revoke revokes a pending invitation.
func (c *Client) Revoke(ctx context.Context, id string) (*clerk.Invitation, error) {
	path, err := clerk.JoinPath(path, id, "revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	invitation := &clerk.Invitation{}
	err = c.Backend.Call(ctx, req, invitation)
	return invitation, err
}
