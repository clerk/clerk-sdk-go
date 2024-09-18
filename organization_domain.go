package clerk

type OrganizationDomainVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
	Attempts int64  `json:"attempts"`
	ExpireAt *int64 `json:"expire_at"`
}

type OrganizationDomain struct {
	APIResource
	Object                  string                          `json:"object"`
	ID                      string                          `json:"id"`
	OrganizationID          string                          `json:"organization_id"`
	Name                    string                          `json:"name"`
	EnrollmentMode          string                          `json:"enrollment_mode"`
	AffiliationEmailAddress *string                         `json:"affiliation_email_address"`
	Verification            *OrganizationDomainVerification `json:"verification"`
	TotalPendingInvitations int64                           `json:"total_pending_invitations"`
	TotalPendingSuggestions int64                           `json:"total_pending_suggestions"`
	CreatedAt               int64                           `json:"created_at"`
	UpdatedAt               int64                           `json:"updated_at"`
}

type OrganizationDomainList struct {
	APIResource
	OrganizationDomains []*OrganizationDomain `json:"data"`
	TotalCount          int64                 `json:"total_count"`
}
