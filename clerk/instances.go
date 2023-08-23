package clerk

import (
	"net/http"
)

type InstanceService service

type UpdateInstanceParams struct {
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

func (s *InstanceService) Update(params UpdateInstanceParams) error {
	req, _ := s.client.NewRequest(http.MethodPatch, "instance", &params)

	_, err := s.client.Do(req, nil)
	return err
}

type InstanceRestrictionsResponse struct {
	Object                 string `json:"object"`
	Allowlist              bool   `json:"allowlist"`
	Blocklist              bool   `json:"blocklist"`
	BlockEmailSubaddresses bool   `json:"block_email_subaddresses"`
}

type UpdateRestrictionsParams struct {
	Allowlist              *bool `json:"allowlist,omitempty"`
	Blocklist              *bool `json:"blocklist,omitempty"`
	BlockEmailSubaddresses *bool `json:"block_email_subaddresses,omitempty"`
}

func (s *InstanceService) UpdateRestrictions(params UpdateRestrictionsParams) (*InstanceRestrictionsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, "instance/restrictions", &params)

	var instanceRestrictionsResponse InstanceRestrictionsResponse
	_, err := s.client.Do(req, &instanceRestrictionsResponse)
	if err != nil {
		return nil, err
	}
	return &instanceRestrictionsResponse, nil
}

type OrganizationSettingsResponse struct {
	Object                 string   `json:"object"`
	Enabled                bool     `json:"enabled"`
	MaxAllowedMemberships  int      `json:"max_allowed_memberships"`
	AdminDeleteEnabled     bool     `json:"admin_delete_enabled"`
	DomainsEnabled         bool     `json:"domains_enabled"`
	DomainsEnrollmentModes []string `json:"domains_enrollment_modes"`
}

type UpdateOrganizationSettingsParams struct {
	Enabled                *bool    `json:"enabled,omitempty"`
	MaxAllowedMemberships  *int     `json:"max_allowed_memberships,omitempty"`
	AdminDeleteEnabled     *bool    `json:"admin_delete_enabled,omitempty"`
	DomainsEnabled         *bool    `json:"domains_enabled,omitempty"`
	DomainsEnrollmentModes []string `json:"domains_enrollment_modes,omitempty"`
}

func (s *InstanceService) UpdateOrganizationSettings(params UpdateOrganizationSettingsParams) (*OrganizationSettingsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, "instance/organization_settings", &params)

	var organizationSettingsResponse OrganizationSettingsResponse
	_, err := s.client.Do(req, &organizationSettingsResponse)
	if err != nil {
		return nil, err
	}
	return &organizationSettingsResponse, nil
}

type UpdateHomeURLParams struct {
	HomeURL string `json:"home_url"`
}

func (s *InstanceService) UpdateHomeURL(params UpdateHomeURLParams) error {
	req, _ := s.client.NewRequest(http.MethodPost, "instance/change_domain", &params)

	_, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
