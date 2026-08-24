package clerk

import "encoding/json"

// APIKey represents an API key response without the secret field.
// This is used for most endpoints except Create.
type APIKey struct {
	APIResource
	Object           string          `json:"object"`
	ID               string          `json:"id"`
	Type             string          `json:"type"`
	Subject          string          `json:"subject"`
	Name             string          `json:"name"`
	Description      *string         `json:"description"`
	Claims           json.RawMessage `json:"claims"`
	Scopes           []string        `json:"scopes"`
	Revoked          bool            `json:"revoked"`
	RevocationReason *string         `json:"revocation_reason"`
	Expired          bool            `json:"expired"`
	Expiration       *int64          `json:"expiration"`
	CreatedBy        *string         `json:"created_by"`
	LastUsedAt       *int64          `json:"last_used_at"`
	CreatedAt        int64           `json:"created_at"`
	UpdatedAt        int64           `json:"updated_at"`
}

// APIKeyWithSecret represents an API key response that includes the secret field.
// This is only used for the Create endpoint.
type APIKeyWithSecret struct {
	APIKey
	Secret string `json:"secret"`
}

// APIKeyList represents a list of API keys without secrets.
type APIKeyList struct {
	APIResource
	APIKeys    []*APIKey `json:"data"`
	TotalCount int64     `json:"total_count"`
}
