package clerk

type SAMLConnection struct {
	APIResource
	ID                 string                         `json:"id"`
	Object             string                         `json:"object"`
	Name               string                         `json:"name"`
	Domain             string                         `json:"domain"`
	IdpEntityID        *string                        `json:"idp_entity_id"`
	IdpSsoURL          *string                        `json:"idp_sso_url"`
	IdpCertificate     *string                        `json:"idp_certificate"`
	IdpMetadataURL     *string                        `json:"idp_metadata_url"`
	IdpMetadata        *string                        `json:"idp_metadata"`
	AcsURL             string                         `json:"acs_url"`
	SPEntityID         string                         `json:"sp_entity_id"`
	SPMetadataURL      string                         `json:"sp_metadata_url"`
	AttributeMapping   SAMLConnectionAttributeMapping `json:"attribute_mapping"`
	Active             bool                           `json:"active"`
	Provider           string                         `json:"provider"`
	UserCount          int64                          `json:"user_count"`
	SyncUserAttributes bool                           `json:"sync_user_attributes"`
	AllowSubdomains    bool                           `json:"allow_subdomains"`
	AllowIdpInitiated  bool                           `json:"allow_idp_initiated"`
	CreatedAt          int64                          `json:"created_at"`
	UpdatedAt          int64                          `json:"updated_at"`
}

type SAMLConnectionAttributeMapping struct {
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type SAMLConnectionList struct {
	APIResource
	SAMLConnections []*SAMLConnection `json:"data"`
	TotalCount      int64             `json:"total_count"`
}
