// Package instancesettings provides the Instance Settings API.
package instancesettings

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/instance"

// Client is used to invoke the Instance Settings API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type UpdateParams struct {
	clerk.APIParams
	// TestMode can be used to toggle test mode for this instance.
	// Defaults to true for development instances.
	TestMode *bool `json:"test_mode,omitempty"`

	// HIBP is used to configure whether Clerk should use the
	// "Have I Been Pawned" service to check passwords against
	// known security breaches.
	// By default, this is enabled in all instances.
	HIBP *bool `json:"hibp,omitempty"`

	// EnhancedEmailDeliverability controls how Clerk delivers emails.
	// Specifically, when set to true, if the instance is a production
	// instance, OTP verification emails are sent by the Clerk's shared
	// domain via Postmark.
	EnhancedEmailDeliverability *bool `json:"enhanced_email_deliverability,omitempty"`

	// SupportEmail is the contact email address that will be displayed
	// on the frontend, in case your instance users need support.
	// If the empty string is provided, the support email that is currently
	// configured in the instance will be removed.
	SupportEmail *string `json:"support_email,omitempty"`

	// ClerkJSVersion allows you to request a specific Clerk JS version on the Clerk Hosted Account pages.
	// If an empty string is provided, the stored version will be removed.
	// If an explicit version is not set, the Clerk JS version will be automatically be resolved.
	ClerkJSVersion *string `json:"clerk_js_version,omitempty"`

	// CookielessDev can be used to enable the new mode in which no third-party
	// cookies are used in development instances. Make sure to also enable the
	// setting in Clerk.js
	//
	// Deprecated: Use URLBasedSessionSyncing instead
	CookielessDev *bool `json:"cookieless_dev,omitempty"`

	// URLBasedSessionSyncing can be used to enable the new mode in which no third-party
	// cookies are used in development instances. Make sure to also enable the
	// setting in Clerk.js
	URLBasedSessionSyncing *bool `json:"url_based_session_syncing,omitempty"`

	// URL that is going to be used in development instances in order to create custom redirects
	// and fix the third-party cookies issues.
	DevelopmentOrigin *string `json:"development_origin,omitempty"`
}

// Update updates the instance's settings.
func (c *Client) Update(ctx context.Context, params *UpdateParams) error {
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	err := c.Backend.Call(ctx, req, &clerk.APIResource{})
	return err
}

type UpdateRestrictionsParams struct {
	clerk.APIParams
	Allowlist                   *bool `json:"allowlist,omitempty"`
	Blocklist                   *bool `json:"blocklist,omitempty"`
	BlockEmailSubaddresses      *bool `json:"block_email_subaddresses,omitempty"`
	BlockDisposableEmailDomains *bool `json:"block_disposable_email_domains,omitempty"`
	IgnoreDotsForGmailAddresses *bool `json:"ignore_dots_for_gmail_addresses,omitempty"`
}

// UpdateRestrictions updates the restriction settings of the instance.
func (c *Client) UpdateRestrictions(ctx context.Context, params *UpdateRestrictionsParams) (*clerk.InstanceRestrictions, error) {
	path, err := clerk.JoinPath(path, "/restrictions")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	instanceRestrictions := &clerk.InstanceRestrictions{}
	err = c.Backend.Call(ctx, req, instanceRestrictions)
	return instanceRestrictions, err
}

type UpdateOrganizationSettingsParams struct {
	clerk.APIParams
	Enabled                *bool     `json:"enabled,omitempty"`
	MaxAllowedMemberships  *int64    `json:"max_allowed_memberships,omitempty"`
	CreatorRoleID          *string   `json:"creator_role_id,omitempty"`
	AdminDeleteEnabled     *bool     `json:"admin_delete_enabled,omitempty"`
	DomainsEnabled         *bool     `json:"domains_enabled,omitempty"`
	DomainsEnrollmentModes *[]string `json:"domains_enrollment_modes,omitempty"`
	DomainsDefaultRoleID   *string   `json:"domains_default_role_id,omitempty"`
}

// UpdateOrganizationSettings updates the organization settings of the instance.
func (c *Client) UpdateOrganizationSettings(ctx context.Context, params *UpdateOrganizationSettingsParams) (*clerk.OrganizationSettings, error) {
	path, err := clerk.JoinPath(path, "/organization_settings")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	orgSettings := &clerk.OrganizationSettings{}
	err = c.Backend.Call(ctx, req, orgSettings)
	return orgSettings, err
}
