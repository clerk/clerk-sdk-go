package clerk

// AdminPortalLinkToken represents an admin portal link token response without
// the secret. Returned by every endpoint except Create.
type AdminPortalLinkToken struct {
	APIResource
	Object            string   `json:"object"`
	ID                string   `json:"id"`
	AdminPortalLinkID string   `json:"admin_portal_link_id"`
	InstanceID        string   `json:"instance_id"`
	OrganizationID    *string  `json:"organization_id"`
	ITContactID       *string  `json:"it_contact_id"`
	Scopes            []string `json:"scopes"`
	Revoked           bool     `json:"revoked"`
	RevocationReason  *string  `json:"revocation_reason"`
	Expired           bool     `json:"expired"`
	Expiration        *int64   `json:"expiration"`
	CreatedAt         int64    `json:"created_at"`
	UpdatedAt         int64    `json:"updated_at"`
}

// AdminPortalLinkTokenWithToken represents an admin portal link token response
// that includes the single-use token. Returned only by Create — the token is
// shown once and must be embedded in the URL delivered to the IT contact.
type AdminPortalLinkTokenWithToken struct {
	AdminPortalLinkToken
	Token string `json:"token"`
}
