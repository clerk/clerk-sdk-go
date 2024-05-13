package clerk

type InstanceRestrictions struct {
	APIResource
	Object                      string `json:"object"`
	Allowlist                   bool   `json:"allowlist"`
	Blocklist                   bool   `json:"blocklist"`
	BlockEmailSubaddresses      bool   `json:"block_email_subaddresses"`
	BlockDisposableEmailDomains bool   `json:"block_disposable_email_domains"`
	IgnoreDotsForGmailAddresses bool   `json:"ignore_dots_for_gmail_addresses"`
}

type OrganizationSettings struct {
	APIResource
	Object                 string   `json:"object"`
	Enabled                bool     `json:"enabled"`
	MaxAllowedMemberships  int64    `json:"max_allowed_memberships"`
	MaxAllowedRoles        int64    `json:"max_allowed_roles"`
	MaxAllowedPermissions  int64    `json:"max_allowed_permissions"`
	CreatorRole            string   `json:"creator_role"`
	AdminDeleteEnabled     bool     `json:"admin_delete_enabled"`
	DomainsEnabled         bool     `json:"domains_enabled"`
	DomainsEnrollmentModes []string `json:"domains_enrollment_modes"`
	DomainsDefaultRole     string   `json:"domains_default_role"`
}
