package clerk

import "encoding/json"

type User struct {
	APIResource
	Object                        string                    `json:"object"`
	ID                            string                    `json:"id"`
	Username                      *string                   `json:"username"`
	FirstName                     *string                   `json:"first_name"`
	LastName                      *string                   `json:"last_name"`
	ImageURL                      *string                   `json:"image_url,omitempty"`
	HasImage                      bool                      `json:"has_image"`
	PrimaryEmailAddressID         *string                   `json:"primary_email_address_id"`
	PrimaryPhoneNumberID          *string                   `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID           *string                   `json:"primary_web3_wallet_id"`
	PasswordEnabled               bool                      `json:"password_enabled"`
	TwoFactorEnabled              bool                      `json:"two_factor_enabled"`
	TOTPEnabled                   bool                      `json:"totp_enabled"`
	BackupCodeEnabled             bool                      `json:"backup_code_enabled"`
	EmailAddresses                []*EmailAddress           `json:"email_addresses"`
	PhoneNumbers                  []*PhoneNumber            `json:"phone_numbers"`
	Web3Wallets                   []*Web3Wallet             `json:"web3_wallets"`
	Passkeys                      []*Passkey                `json:"passkeys"`
	OrganizationMemberships       []*OrganizationMembership `json:"organization_memberships,omitempty"`
	ExternalAccounts              []*ExternalAccount        `json:"external_accounts"`
	SAMLAccounts                  []*SAMLAccount            `json:"saml_accounts"`
	EnterpriseAccounts            []*EnterpriseAccount      `json:"enterprise_accounts"`
	PasswordLastUpdatedAt         *int64                    `json:"password_last_updated_at,omitempty"`
	PublicMetadata                json.RawMessage           `json:"public_metadata"`
	PrivateMetadata               json.RawMessage           `json:"private_metadata,omitempty"`
	UnsafeMetadata                json.RawMessage           `json:"unsafe_metadata,omitempty"`
	ExternalID                    *string                   `json:"external_id"`
	LastSignInAt                  *int64                    `json:"last_sign_in_at"`
	Banned                        bool                      `json:"banned"`
	Locked                        bool                      `json:"locked"`
	LockoutExpiresInSeconds       *int64                    `json:"lockout_expires_in_seconds"`
	VerificationAttemptsRemaining *int64                    `json:"verification_attempts_remaining"`
	CreatedAt                     int64                     `json:"created_at"`
	UpdatedAt                     int64                     `json:"updated_at"`
	DeleteSelfEnabled             bool                      `json:"delete_self_enabled"`
	CreateOrganizationEnabled     bool                      `json:"create_organization_enabled"`
	CreateOrganizationsLimit      *int                      `json:"create_organizations_limit,omitempty"`
	LastActiveAt                  *int64                    `json:"last_active_at"`
	MFAEnabledAt                  *int64                    `json:"mfa_enabled_at"`
	Locale                        *string                   `json:"locale"`
	MFADisabledAt                 *int64                    `json:"mfa_disabled_at"`
	LegalAcceptedAt               *int64                    `json:"legal_accepted_at"`
}

type ExternalAccount struct {
	Object           string          `json:"object"`
	ID               string          `json:"id"`
	Provider         string          `json:"provider"`
	IdentificationID string          `json:"identification_id"`
	ProviderUserID   string          `json:"provider_user_id"`
	ApprovedScopes   string          `json:"approved_scopes"`
	EmailAddress     string          `json:"email_address"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	AvatarURL        string          `json:"avatar_url"`
	ImageURL         *string         `json:"image_url,omitempty"`
	Username         *string         `json:"username"`
	PhoneNumber      *string         `json:"phone_number"`
	PublicMetadata   json.RawMessage `json:"public_metadata"`
	Label            *string         `json:"label"`
	CreatedAt        int64           `json:"created_at"`
	UpdatedAt        int64           `json:"updated_at"`
	Verification     *Verification   `json:"verification"`
}

type Web3Wallet struct {
	Object       string        `json:"object"`
	ID           string        `json:"id"`
	Web3Wallet   string        `json:"web3_wallet"`
	Verification *Verification `json:"verification"`
	CreatedAt    int64         `json:"created_at"`
	UpdatedAt    int64         `json:"updated_at"`
}

type SAMLAccount struct {
	Object                 string                 `json:"object"`
	ID                     string                 `json:"id"`
	Provider               string                 `json:"provider"`
	Active                 bool                   `json:"active"`
	EmailAddress           string                 `json:"email_address"`
	FirstName              *string                `json:"first_name"`
	LastName               *string                `json:"last_name"`
	ProviderUserID         *string                `json:"provider_user_id"`
	LastAuthenticatedAt    *int64                 `json:"last_authenticated_at"`
	PublicMetadata         json.RawMessage        `json:"public_metadata"`
	SAMLConnection         *SAMLAccountConnection `json:"saml_connection"`
	EnterpriseConnectionID string                 `json:"enterprise_connection_id"`
	Verification           *Verification          `json:"verification"`
}

type UserList struct {
	APIResource
	Users      []*User `json:"data"`
	TotalCount int64   `json:"total_count"`
}

type SAMLAccountConnection struct {
	ID                               string   `json:"id"`
	Name                             string   `json:"name"`
	Domains                          []string `json:"domains"`
	Active                           bool     `json:"active"`
	Provider                         string   `json:"provider"`
	SyncUserAttributes               bool     `json:"sync_user_attributes"`
	AllowSubdomains                  bool     `json:"allow_subdomains"`
	AllowIDPInitiated                bool     `json:"allow_idp_initiated"`
	DisableAdditionalIdentifications bool     `json:"disable_additional_identifications"`
	CreatedAt                        int64    `json:"created_at"`
	UpdatedAt                        int64    `json:"updated_at"`
}
