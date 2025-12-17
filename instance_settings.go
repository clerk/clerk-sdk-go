package clerk

type InstanceRestrictions struct {
	APIResource
	Object                             string `json:"object"`
	Allowlist                          bool   `json:"allowlist"`
	Blocklist                          bool   `json:"blocklist"`
	AllowlistBlocklistDisabledOnSignIn bool   `json:"allowlist_blocklist_disabled_on_sign_in"`
	BlockEmailSubaddresses             bool   `json:"block_email_subaddresses"`
	BlockDisposableEmailDomains        bool   `json:"block_disposable_email_domains"`
	IgnoreDotsForGmailAddresses        bool   `json:"ignore_dots_for_gmail_addresses"`
}

type OrganizationSettings struct {
	APIResource
	Object                 string   `json:"object"`
	Enabled                bool     `json:"enabled"`
	MaxAllowedMemberships  int64    `json:"max_allowed_memberships"`
	MaxAllowedRoles        int64    `json:"max_allowed_roles"`
	MaxAllowedPermissions  int64    `json:"max_allowed_permissions"`
	MaxRoleSetsAllowed     int64    `json:"max_role_sets_allowed"`
	CreatorRole            string   `json:"creator_role"`
	AdminDeleteEnabled     bool     `json:"admin_delete_enabled"`
	DomainsEnabled         bool     `json:"domains_enabled"`
	SlugDisabled           bool     `json:"slug_disabled"`
	DomainsEnrollmentModes []string `json:"domains_enrollment_modes"`
	DomainsDefaultRole     string   `json:"domains_default_role"`
	// TODO(gabriel): remove Remove omitempty when feat is out
	OrganizationCreationDefaults *OrganizationCreationDefaults `json:"default_organization_naming,omitempty"`
	// TODO(nicolas): Remove omitempty when it's GA
	InitialRoleSetKey *string `json:"initial_role_set_key,omitempty"`
}

type OrganizationCreationDefaults struct {
	Enabled                       bool                                  `json:"enabled"`
	AutomaticOrganizationCreation AutomaticOrganizationCreationSettings `json:"automatic_organization_creation"`
	DetectFromEmailDomain         DetectFromEmailDomainSettings         `json:"detect_from_email_domain"`
	OrganizationNameTemplate      OrganizationNameTemplateSettings      `json:"organization_name_template"`
	Fallback                      FallbackSettings                      `json:"fallback"`
}

type AutomaticOrganizationCreationSettings struct {
	Enabled bool `json:"enabled"`
}

type DetectFromEmailDomainSettings struct {
	Enabled bool `json:"enabled"`
}

type OrganizationNameTemplateSettings struct {
	Enabled  bool   `json:"enabled"`
	Template string `json:"template"`
}

type FallbackSettings struct {
	Template string `json:"template"`
}
