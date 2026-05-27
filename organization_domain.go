package clerk

type OrganizationDomainVerification struct {
	Status     string `json:"status"`
	Strategy   string `json:"strategy"`
	Attempts   int64  `json:"attempts"`
	ExpireAt   *int64 `json:"expire_at"`
	VerifiedAt *int64 `json:"verified_at"`
}

// OrganizationDomainOwnershipVerification describes how ownership of the
// underlying DNS domain was proven (e.g. via a DNS TXT challenge). The TXT
// record fields are populated only while a DNS challenge is pending.
type OrganizationDomainOwnershipVerification struct {
	OrganizationDomainVerification
	TXTRecordName  *string `json:"txt_record_name"`
	TXTRecordValue *string `json:"txt_record_value"`
}

type OrganizationDomain struct {
	APIResource
	Object                  string                                   `json:"object"`
	ID                      string                                   `json:"id"`
	OrganizationID          string                                   `json:"organization_id"`
	Name                    string                                   `json:"name"`
	EnrollmentMode          string                                   `json:"enrollment_mode"`
	AffiliationEmailAddress *string                                  `json:"affiliation_email_address"`
	AffiliationVerification *OrganizationDomainVerification          `json:"affiliation_verification"`
	OwnershipVerification   *OrganizationDomainOwnershipVerification `json:"ownership_verification"`
	// Deprecated: use AffiliationVerification instead.
	Verification            *OrganizationDomainVerification `json:"verification"`
	TotalPendingInvitations int64                           `json:"total_pending_invitations"`
	TotalPendingSuggestions int64                           `json:"total_pending_suggestions"`
	CreatedAt               int64                           `json:"created_at"`
	UpdatedAt               int64                           `json:"updated_at"`
	PublicOrganizationData  *PublicOrganizationData         `json:"public_organization_data,omitempty"`
}

type OrganizationDomainList struct {
	APIResource
	OrganizationDomains []*OrganizationDomain `json:"data"`
	TotalCount          int64                 `json:"total_count"`
}
