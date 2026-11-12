package clerk

type AccessPolicyDomain struct {
	Name string `json:"name"`
}

type AccessPolicy struct {
	APIResource
	Object  string `json:"object"`
	ID      string `json:"id"`
	Default bool   `json:"default"`
	// OrganizationID is set while an organization's sso access policy signs
	// in through this policy.
	OrganizationID *string               `json:"organization_id"`
	Domains        []*AccessPolicyDomain `json:"domains"`
	Status         string                `json:"status"`
	// Authentication is nil on the default policy.
	Authentication *OrganizationAccessPolicyAuthentication `json:"authentication"`
	CreatedAt      int64                                   `json:"created_at"`
	UpdatedAt      int64                                   `json:"updated_at"`
}

type AccessPolicyList struct {
	APIResource
	AccessPolicies []*AccessPolicy `json:"data"`
	TotalCount     int64           `json:"total_count"`
}
