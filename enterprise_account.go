package clerk

import "encoding/json"

type EnterpriseAccount struct {
	ID                     string                       `json:"id"`
	Object                 string                       `json:"object"`
	Protocol               string                       `json:"protocol"`
	Provider               string                       `json:"provider"`
	Active                 bool                         `json:"active"`
	EmailAddress           string                       `json:"email_address"`
	FirstName              *string                      `json:"first_name"`
	LastName               *string                      `json:"last_name"`
	ProviderUserID         *string                      `json:"provider_user_id"`
	LastAuthenticatedAt    *int64                       `json:"last_authenticated_at"`
	PublicMetadata         json.RawMessage              `json:"public_metadata" logger:"omit"`
	Verification           *Verification                `json:"verification"`
	EnterpriseConnection   *EnterpriseAccountConnection `json:"enterprise_connection"`
	EnterpriseConnectionID string                       `json:"enterprise_connection_id"`
}

type EnterpriseAccountConnection struct {
	// ID belongs to the underlying connection, either SAML Connection or the OAuth config
	ID                     string `json:"id"`
	EnterpriseConnectionID string `json:"enterprise_connection_id"`
	Protocol               string `json:"protocol"`
	Provider               string `json:"provider"`

	// Name is the name of this enterprise connection that we will display directly to end-users
	Name                             string   `json:"name"`
	LogoPublicURL                    *string  `json:"logo_public_url"`
	Domains                          []string `json:"domains,omitempty"`
	Active                           bool     `json:"active"`
	SyncUserAttributes               bool     `json:"sync_user_attributes"`
	DisableAdditionalIdentifications bool     `json:"disable_additional_identifications"`
	CreatedAt                        int64    `json:"created_at"`
	UpdatedAt                        int64    `json:"updated_at"`

	// SAML Connection specific fields
	AllowSubdomains   bool `json:"allow_subdomains"`
	AllowIDPInitiated bool `json:"allow_idp_initiated"`
}
