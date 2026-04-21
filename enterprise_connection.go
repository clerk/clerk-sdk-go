package clerk

// EnterpriseConnection represents an enterprise SSO connection (SAML or OAuth/OIDC)
// in the Clerk Backend API.
type EnterpriseConnection struct {
	APIResource
	ID                               string   `json:"id"`
	Object                           string   `json:"object"`
	Name                             string   `json:"name"`
	Protocol                         string   `json:"protocol"`
	Provider                         string   `json:"provider"`
	OrganizationID                   *string  `json:"organization_id,omitempty"`
	Domains                          []string `json:"domains,omitempty"`
	LogoPublicURL                    *string  `json:"logo_public_url,omitempty"`
	Active                           bool     `json:"active"`
	SyncUserAttributes               bool     `json:"sync_user_attributes"`
	DisableAdditionalIdentifications bool     `json:"disable_additional_identifications"`
	AllowSubdomains                  bool     `json:"allow_subdomains"`
	AllowIdpInitiated                bool     `json:"allow_idp_initiated"`
	ForceAuthn                       bool     `json:"force_authn"`
	AllowOrganizationAccountLinking  bool     `json:"allow_organization_account_linking"`
	Authenticatable                  bool     `json:"authenticatable"`
	CreatedAt                        int64    `json:"created_at"`
	UpdatedAt                        int64    `json:"updated_at"`

	// SAML-specific fields (when protocol is "saml")
	IdpEntityID      *string                         `json:"idp_entity_id,omitempty"`
	IdpSsoURL        *string                         `json:"idp_sso_url,omitempty"`
	AcsURL           string                          `json:"acs_url,omitempty"`
	SPEntityID       string                          `json:"sp_entity_id,omitempty"`
	SPMetadataURL    string                          `json:"sp_metadata_url,omitempty"`
	AttributeMapping *SAMLConnectionAttributeMapping `json:"attribute_mapping,omitempty"`
	CustomAttributes *[]CustomAttribute              `json:"custom_attributes,omitempty"`

	// OIDC-specific fields (when protocol is "oauth_oidc")
	ClientID         *string `json:"client_id,omitempty"`
	IssuerURL        *string `json:"issuer_url,omitempty"`
	AuthorizationURL *string `json:"authorization_url,omitempty"`
	TokenURL         *string `json:"token_url,omitempty"`
	UserInfoURL      *string `json:"user_info_url,omitempty"`
	Scopes           *string `json:"scopes,omitempty"`
}

// EnterpriseConnectionList is the response for listing enterprise connections.
type EnterpriseConnectionList struct {
	APIResource
	EnterpriseConnections []*EnterpriseConnection `json:"data"`
	TotalCount            int64                   `json:"total_count"`
}
