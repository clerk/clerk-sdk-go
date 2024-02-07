package clerk

import "encoding/json"

type JWTTemplate struct {
	APIResource
	Object           string          `json:"object"`
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Claims           json.RawMessage `json:"claims"`
	Lifetime         int64           `json:"lifetime"`
	AllowedClockSkew int64           `json:"allowed_clock_skew"`
	CustomSigningKey bool            `json:"custom_signing_key"`
	SigningAlgorithm string          `json:"signing_algorithm"`
	CreatedAt        int64           `json:"created_at"`
	UpdatedAt        int64           `json:"updated_at"`
}

type JWTTemplateList struct {
	APIResource
	JWTTemplates []*JWTTemplate `json:"data"`
	TotalCount   int64          `json:"total_count"`
}
