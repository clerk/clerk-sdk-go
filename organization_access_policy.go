package clerk

type OrganizationAccessPolicyDomain struct {
	Name                string `json:"name"`
	AffiliationVerified bool   `json:"affiliation_verified"`
	OwnershipVerified   bool   `json:"ownership_verified"`
}

type OrganizationAccessPolicyConnection struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type OrganizationAccessPolicyDirectory struct {
	ID string `json:"id"`
}

// OrganizationAccessPolicyAuthentication is set on sso policies only.
type OrganizationAccessPolicyAuthentication struct {
	Connections        []*OrganizationAccessPolicyConnection `json:"connections"`
	Directories        []*OrganizationAccessPolicyDirectory  `json:"directories"`
	Provisioning       string                                `json:"provisioning"`
	ConnectionRequired bool                                  `json:"connection_required"`
}

type OrganizationAccessPolicy struct {
	APIResource
	Object         string                            `json:"object"`
	ID             string                            `json:"id"`
	OrganizationID string                            `json:"organization_id"`
	Default        bool                              `json:"default"`
	Domains        []*OrganizationAccessPolicyDomain `json:"domains"`
	EnrollmentMode string                            `json:"enrollment_mode"`
	// Status is draft or active on sso policies, and nil on the others.
	Status         *string                                 `json:"status"`
	Authentication *OrganizationAccessPolicyAuthentication `json:"authentication"`
	CreatedAt      int64                                   `json:"created_at"`
	UpdatedAt      int64                                   `json:"updated_at"`
}

type OrganizationAccessPolicyList struct {
	APIResource
	OrganizationAccessPolicies []*OrganizationAccessPolicy `json:"data"`
	TotalCount                 int64                       `json:"total_count"`
}
