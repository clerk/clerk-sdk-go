// Package enterprisedomainverification provides the Enterprise Domain
// Verification API — standalone domain-ownership challenges that exist
// independently of any enterprise connection. Once verified, a record is
// linked to an enterprise connection when the connection is created with the
// matching domain.
package enterprisedomainverification

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/enterprise_domain_verifications"

// Client is used to invoke the Enterprise Domain Verification API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// PrepareParams is the request body for starting a new domain verification.
type PrepareParams struct {
	clerk.APIParams
	Strategy                *string `json:"strategy,omitempty"`
	Domain                  *string `json:"domain,omitempty"`
	AffiliationEmailAddress *string `json:"affiliation_email_address,omitempty"`
}

// Prepare creates a new standalone domain verification.
func (c *Client) Prepare(ctx context.Context, params *PrepareParams) (*clerk.EnterpriseDomainVerification, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	res := &clerk.EnterpriseDomainVerification{}
	err := c.Backend.Call(ctx, req, res)
	return res, err
}

// AttemptParams is the request body for completing a previously prepared
// domain verification.
type AttemptParams struct {
	clerk.APIParams
	VerificationID string  `json:"-"`
	Strategy       *string `json:"strategy,omitempty"`
	Code           *string `json:"code,omitempty"`
}

// Attempt completes a previously prepared domain verification.
func (c *Client) Attempt(ctx context.Context, params *AttemptParams) (*clerk.EnterpriseDomainVerification, error) {
	p, err := clerk.JoinPath(path, params.VerificationID, "/attempt_verification")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, p)
	req.SetParams(params)
	res := &clerk.EnterpriseDomainVerification{}
	err = c.Backend.Call(ctx, req, res)
	return res, err
}
