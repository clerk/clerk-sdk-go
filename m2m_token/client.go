// Package m2m_token provides the M2M Tokens API.
package m2m_token

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/m2m_tokens"

// Client is used to invoke the M2M Tokens API.
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
	SecondsUntilExpiration *int64           `json:"seconds_until_expiration,omitempty"`
	Claims                 *json.RawMessage `json:"claims,omitempty"`
	TokenFormat            *string          `json:"token_format,omitempty"`
}

// Create creates a new M2M token.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.M2MTokenWithToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.M2MTokenWithToken{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	Subject *string `json:"subject,omitempty"`
	Revoked *bool   `json:"revoked,omitempty"`
	Expired *bool   `json:"expired,omitempty"`
}

func (params *ListParams) ToQuery() url.Values {
	q := url.Values{}
	if params.Subject != nil {
		q.Set("subject", *params.Subject)
	}
	if params.Revoked != nil {
		q.Set("revoked", strconv.FormatBool(*params.Revoked))
	}
	if params.Expired != nil {
		q.Set("expired", strconv.FormatBool(*params.Expired))
	}
	// Add list params
	if params.Limit != nil {
		q.Set("limit", strconv.FormatInt(*params.Limit, 10))
	}
	if params.Offset != nil {
		q.Set("offset", strconv.FormatInt(*params.Offset, 10))
	}
	return q
}

// List returns a list of M2M tokens.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.M2MTokenList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &clerk.M2MTokenList{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

type RevokeParams struct {
	clerk.APIParams
	RevocationReason *string `json:"revocation_reason,omitempty"`
}

// Revoke revokes an M2M token.
func (c *Client) Revoke(ctx context.Context, tokenID string, params *RevokeParams) (*clerk.M2MToken, error) {
	path, err := clerk.JoinPath(path, tokenID, "revoke")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.M2MToken{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type VerifyParams struct {
	clerk.APIParams
	Token string `json:"token"`
}

// Verify verifies an M2M token.
func (c *Client) Verify(ctx context.Context, params *VerifyParams) (*clerk.M2MToken, error) {
	path, err := clerk.JoinPath(path, "verify")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.M2MToken{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
