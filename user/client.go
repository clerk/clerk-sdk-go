// Package user provides the Users API.
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/users"

// Client is used to invoke the Users API.
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
	EmailAddresses          *[]string        `json:"email_address,omitempty"`
	PhoneNumbers            *[]string        `json:"phone_number,omitempty"`
	Web3Wallets             *[]string        `json:"web3_wallet,omitempty"`
	Username                *string          `json:"username,omitempty"`
	Password                *string          `json:"password,omitempty"`
	FirstName               *string          `json:"first_name,omitempty"`
	LastName                *string          `json:"last_name,omitempty"`
	ExternalID              *string          `json:"external_id,omitempty"`
	UnsafeMetadata          *json.RawMessage `json:"unsafe_metadata,omitempty"`
	PublicMetadata          *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata         *json.RawMessage `json:"private_metadata,omitempty"`
	PasswordDigest          *string          `json:"password_digest,omitempty"`
	PasswordHasher          *string          `json:"password_hasher,omitempty"`
	SkipPasswordRequirement *bool            `json:"skip_password_requirement,omitempty"`
	SkipPasswordChecks      *bool            `json:"skip_password_checks,omitempty"`
	TOTPSecret              *string          `json:"totp_secret,omitempty"`
	BackupCodes             *[]string        `json:"backup_codes,omitempty"`
	// Specified in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
}

// Create creates a new user.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.User, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.User{}
	err := c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Get retrieves details about the user.
func (c *Client) Get(ctx context.Context, id string) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type UpdateParams struct {
	clerk.APIParams
	FirstName                        *string          `json:"first_name,omitempty"`
	LastName                         *string          `json:"last_name,omitempty"`
	PrimaryEmailAddressID            *string          `json:"primary_email_address_id,omitempty"`
	NotifyPrimaryEmailAddressChanged *bool            `json:"notify_primary_email_address_changed,omitempty"`
	PrimaryPhoneNumberID             *string          `json:"primary_phone_number_id,omitempty"`
	PrimaryWeb3WalletID              *string          `json:"primary_web3_wallet_id,omitempty"`
	Username                         *string          `json:"username,omitempty"`
	ProfileImageID                   *string          `json:"profile_image_id,omitempty"`
	ProfileImage                     *string          `json:"profile_image,omitempty"`
	Password                         *string          `json:"password,omitempty"`
	PasswordDigest                   *string          `json:"password_digest,omitempty"`
	PasswordHasher                   *string          `json:"password_hasher,omitempty"`
	SkipPasswordChecks               *bool            `json:"skip_password_checks,omitempty"`
	SignOutOfOtherSessions           *bool            `json:"sign_out_of_other_sessions,omitempty"`
	ExternalID                       *string          `json:"external_id,omitempty"`
	PublicMetadata                   *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata                  *json.RawMessage `json:"private_metadata,omitempty"`
	UnsafeMetadata                   *json.RawMessage `json:"unsafe_metadata,omitempty"`
	TOTPSecret                       *string          `json:"totp_secret,omitempty"`
	BackupCodes                      *[]string        `json:"backup_codes,omitempty"`
	DeleteSelfEnabled                *bool            `json:"delete_self_enabled,omitempty"`
	CreateOrganizationEnabled        *bool            `json:"create_organization_enabled,omitempty"`
	// Specified in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
}

// Update updates a user.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type UpdateMetadataParams struct {
	clerk.APIParams
	PublicMetadata  *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata *json.RawMessage `json:"private_metadata,omitempty"`
	UnsafeMetadata  *json.RawMessage `json:"unsafe_metadata,omitempty"`
}

// UpdateMetadata updates the user's metadata by merging the
// provided values with the existing ones.
func (c *Client) UpdateMetadata(ctx context.Context, id string, params *UpdateMetadataParams) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id, "/metadata")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Delete deletes a user.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrderBy           *string  `json:"order_by,omitempty"`
	Query             *string  `json:"query,omitempty"`
	EmailAddresses    []string `json:"email_address,omitempty"`
	ExternalIDs       []string `json:"external_id,omitempty"`
	PhoneNumbers      []string `json:"phone_number,omitempty"`
	Web3Wallets       []string `json:"web3_wallet,omitempty"`
	Usernames         []string `json:"username,omitempty"`
	UserIDs           []string `json:"user_id,omitempty"`
	LastActiveAtSince *int64   `json:"last_active_at_since,omitempty"`
}

// ToQuery returns url.Values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.OrderBy != nil {
		q.Add("order_by", *params.OrderBy)
	}
	if params.Query != nil {
		q.Add("query", *params.Query)
	}
	for _, v := range params.EmailAddresses {
		q.Add("email_address", v)
	}
	for _, v := range params.ExternalIDs {
		q.Add("external_id", v)
	}
	for _, v := range params.PhoneNumbers {
		q.Add("phone_number", v)
	}
	for _, v := range params.Web3Wallets {
		q.Add("web3_wallet", v)
	}
	for _, v := range params.Usernames {
		q.Add("username", v)
	}
	for _, v := range params.UserIDs {
		q.Add("user_id", v)
	}
	if params.LastActiveAtSince != nil {
		q.Add("last_active_at_since", strconv.FormatInt(*params.LastActiveAtSince, 10))
	}
	return q
}

// List returns a list of users.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.UserList, error) {
	// The Clerk API returns the results of GET /v1/users as an
	// array. In order to build the final response that includes
	// the total count, we need to make two API calls.
	// GET /v1/users retrieves the actual results
	// GET /v1/users/count retrieves the total count
	// The response is then synthesized from the individual responses.

	// GET /v1/users
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	data := &userList{}
	err := c.Backend.Call(ctx, req, data)
	if err != nil {
		return nil, err
	}

	// GET /v1/users/count
	totalCount, err := c.Count(ctx, params)
	if err != nil {
		return nil, err
	}

	users := []*clerk.User(*data)
	return &clerk.UserList{
		Users:      users,
		TotalCount: totalCount.TotalCount,
	}, nil
}

// Count returns the total count of users satisfying the parameters.
func (c *Client) Count(ctx context.Context, params *ListParams) (*TotalCount, error) {
	path, err := clerk.JoinPath(path, "/count")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &TotalCount{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Custom type needed in order to store the GET /v1/users results
// array.
type userList []*clerk.User

// Read implements the clerk.ResponseReader interface.
// The implementation is empty, meaning that we'll lose
// the raw response from the server.
func (_ *userList) Read(res *clerk.APIResponse) {
	// no-op
}

// Response schema for GET /v1/users/count
type TotalCount struct {
	clerk.APIResource
	Object     string `json:"object"`
	TotalCount int64  `json:"total_count"`
}

type ListOAuthAccessTokensParams struct {
	clerk.APIParams
	ID       string `json:"-"`
	Provider string `json:"-"`
}

// ListOAuthAccessTokens retrieves a list of the user's access
// tokens for a specific OAuth provider.
func (c *Client) ListOAuthAccessTokens(ctx context.Context, params *ListOAuthAccessTokensParams) (*clerk.OAuthAccessTokenList, error) {
	path, err := clerk.JoinPath(path, params.ID, "/oauth_access_tokens", fmt.Sprintf("%s?paginated=true", params.Provider))
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OAuthAccessTokenList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}

type DeleteMFAParams struct {
	clerk.APIParams
	ID string `json:"-"`
}

// DeleteMFA disables a user's multi-factor authentication methods.
func (c *Client) DeleteMFA(ctx context.Context, params *DeleteMFAParams) (*MultifactorAuthentication, error) {
	path, err := clerk.JoinPath(path, params.ID, "/mfa")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &MultifactorAuthentication{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type MultifactorAuthentication struct {
	clerk.APIResource
	UserID string `json:"user_id"`
}

// Ban marks the user as banned.
func (c *Client) Ban(ctx context.Context, id string) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id, "/ban")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Unban removes the ban for a user.
func (c *Client) Unban(ctx context.Context, id string) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id, "/unban")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Lock marks the user as locked.
func (c *Client) Lock(ctx context.Context, id string) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id, "/lock")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// Unlock removes the lock for a user.
func (c *Client) Unlock(ctx context.Context, id string) (*clerk.User, error) {
	path, err := clerk.JoinPath(path, id, "/unlock")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.User{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListOrganizationMembershipsParams struct {
	clerk.APIParams
	clerk.ListParams
	ID string `json:"-"`
}

// ToQuery returns url.Values from the params.
func (params *ListOrganizationMembershipsParams) ToQuery() url.Values {
	return params.ListParams.ToQuery()
}

// ListOrganizationMemberships lists all the user's organization memberships.
func (c *Client) ListOrganizationMemberships(ctx context.Context, params *ListOrganizationMembershipsParams) (*clerk.OrganizationMembershipList, error) {
	path, err := clerk.JoinPath(path, params.ID, "/organization_memberships")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.OrganizationMembershipList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}
