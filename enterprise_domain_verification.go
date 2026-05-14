package clerk

type EnterpriseDomainVerificationDetails struct {
	Status       string  `json:"status"`
	Strategy     string  `json:"strategy"`
	Attempts     *int    `json:"attempts"`
	ExpireAt     *int64  `json:"expire_at"`
	DNSTxtRecord *string `json:"dns_txt_record,omitempty"`
}

type EnterpriseDomainVerification struct {
	APIResource
	Object                 string                               `json:"object"`
	ID                     string                               `json:"id"`
	Domain                 string                               `json:"domain"`
	EnterpriseConnectionID *string                              `json:"enterprise_connection_id"`
	EnterpriseDomainID     *string                              `json:"enterprise_domain_id"`
	Verified               bool                                 `json:"verified"`
	Verification           *EnterpriseDomainVerificationDetails `json:"verification"`
}
