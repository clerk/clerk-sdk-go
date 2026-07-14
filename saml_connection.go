package clerk

type CustomAttribute struct {
	Name        *string `json:"name"`
	Key         *string `json:"key"`
	SSOPath     *string `json:"sso_path"`
	SCIMPath    *string `json:"scim_path"`
	MultiValued *bool   `json:"multi_valued,omitempty"`
}

type SAMLConnection struct {
	APIResource
	ID     string `json:"id"`
	Object string `json:"object"`
	Name   string `json:"name"`
	// Deprecated: Use `domains` instead.
	Domain                           string                         `json:"domain"`
	Domains                          []string                       `json:"domains"`
	IdpEntityID                      *string                        `json:"idp_entity_id"`
	OrganizationID                   *string                        `json:"organization_id"`
	EnterpriseConnectionID           *string                        `json:"enterprise_connection_id"`
	IdpSsoURL                        *string                        `json:"idp_sso_url"`
	IdpCertificate                   *string                        `json:"idp_certificate"`
	IdpCertificateIssuedAt           *int64                         `json:"idp_certificate_issued_at"`
	IdpCertificateExpiresAt          *int64                         `json:"idp_certificate_expires_at"`
	IdpMetadataURL                   *string                        `json:"idp_metadata_url"`
	IdpMetadata                      *string                        `json:"idp_metadata"`
	AcsURL                           string                         `json:"acs_url"`
	SPEntityID                       string                         `json:"sp_entity_id"`
	SPMetadataURL                    string                         `json:"sp_metadata_url"`
	AttributeMapping                 SAMLConnectionAttributeMapping `json:"attribute_mapping"`
	Active                           bool                           `json:"active"`
	Provider                         string                         `json:"provider"`
	UserCount                        int64                          `json:"user_count"`
	SyncUserAttributes               bool                           `json:"sync_user_attributes"`
	AllowSubdomains                  bool                           `json:"allow_subdomains"`
	AllowIdpInitiated                bool                           `json:"allow_idp_initiated"`
	DisableAdditionalIdentifications bool                           `json:"disable_additional_identifications"`
	ForceAuthn                       bool                           `json:"force_authn"`
	DisableJITProvisioning           *bool                          `json:"disable_jit_provisioning,omitempty"`
	AllowOrganizationAccountLinking  bool                           `json:"allow_organization_account_linking"`
	CustomAttributes                 *[]CustomAttribute             `json:"custom_attributes,omitempty"`
	Authenticatable                  bool                           `json:"authenticatable"`
	LoginHint                        *LoginHint                     `json:"login_hint,omitempty"`
	CreatedAt                        int64                          `json:"created_at"`
	UpdatedAt                        int64                          `json:"updated_at"`
}

type SAMLConnectionAttributeMapping struct {
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// LoginHint describes the login hint a SAML connection sends to the identity
// provider when initiating a sign-in.
type LoginHint struct {
	Mode   string  `json:"mode"`
	Source *string `json:"source,omitempty"`
}

type SAMLConnectionList struct {
	APIResource
	SAMLConnections []*SAMLConnection `json:"data"`
	TotalCount      int64             `json:"total_count"`
}
