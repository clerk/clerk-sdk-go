package clerk

import "encoding/json"

// M2MToken represents a machine-to-machine token response without the secret field.
// This is used for most endpoints except Create.
type M2MToken struct {
	APIResource
	Object           string          `json:"object"`
	ID               string          `json:"id"`
	Subject          string          `json:"subject"`
	Claims           json.RawMessage `json:"claims"`
	Scopes           []string        `json:"scopes"`
	Revoked          bool            `json:"revoked"`
	RevocationReason *string         `json:"revocation_reason"`
	Expired          bool            `json:"expired"`
	Expiration       *int64          `json:"expiration"`
	LastUsedAt       *int64          `json:"last_used_at"`
	CreatedAt        int64           `json:"created_at"`
	UpdatedAt        int64           `json:"updated_at"`
}

// M2MTokenWithToken represents a machine-to-machine token response that includes the token field.
// This is only used for the Create endpoint.
type M2MTokenWithToken struct {
	M2MToken
	Token string `json:"token"`
}

// M2MTokenList represents a list of machine-to-machine tokens.
type M2MTokenList struct {
	APIResource
	M2MTokens  []*M2MToken `json:"m2m_tokens"`
	TotalCount int64       `json:"total_count"`
}
