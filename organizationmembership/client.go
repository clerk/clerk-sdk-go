// Package organizationmembership provides the Organization Memberships API.
package organizationmembership

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/organizations"

// Client is used to invoke the Organization Memberships API.
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
	UserID         *string `json:"user_id,omitempty"`
	Role           *string `json:"role,omitempty"`
	OrganizationID string  `json:"-"`
}

// Create adds a new member to the organization.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.OrganizationMembership, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/memberships")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	membership := &clerk.OrganizationMembership{}
	err = c.Backend.Call(ctx, req, membership)
	return membership, err
}

type UpdateParams struct {
	clerk.APIParams
	Role           *string `json:"role,omitempty"`
	OrganizationID string  `json:"-"`
	UserID         string  `json:"-"`
}

// Update updates an organization membership.
func (c *Client) Update(ctx context.Context, params *UpdateParams) (*clerk.OrganizationMembership, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/memberships", params.UserID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	membership := &clerk.OrganizationMembership{}
	err = c.Backend.Call(ctx, req, membership)
	return membership, err
}

type DeleteParams struct {
	clerk.APIParams
	OrganizationID string `json:"-"`
	UserID         string `json:"-"`
}

// Delete removes a member from an organization.
func (c *Client) Delete(ctx context.Context, params *DeleteParams) (*clerk.OrganizationMembership, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/memberships", params.UserID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	membership := &clerk.OrganizationMembership{}
	err = c.Backend.Call(ctx, req, membership)
	return membership, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrderBy        *string  `json:"order_by,omitempty"`
	Query          *string  `json:"query,omitempty"`
	Roles          []string `json:"role,omitempty"`
	UserIDs        []string `json:"user_id,omitempty"`
	EmailAddresses []string `json:"email_address,omitempty"`
	PhoneNumbers   []string `json:"phone_number,omitempty"`
	Usernames      []string `json:"username,omitempty"`
	Web3Wallets    []string `json:"web3_wallet,omitempty"`
	OrganizationID string   `json:"-"`
}

// ToQuery returns the parameters as url.Values so they can be used
// in a URL query string.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	if params.Roles != nil {
		q["role"] = params.Roles
	}
	if params.UserIDs != nil {
		q["user_id"] = params.UserIDs
	}
	if params.EmailAddresses != nil {
		q["email_address"] = params.EmailAddresses
	}
	if params.PhoneNumbers != nil {
		q["phone_number"] = params.PhoneNumbers
	}
	if params.Usernames != nil {
		q["username"] = params.Usernames
	}
	if params.Web3Wallets != nil {
		q["web3_wallet"] = params.Web3Wallets
	}
	return q
}

// List returns a list of organization memberships.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.OrganizationMembershipList, error) {
	path, err := clerk.JoinPath(path, params.OrganizationID, "/memberships")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OrganizationMembershipList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}
