package clerk

import "encoding/json"

type OAuthAccessToken struct {
	Object         string          `json:"object"`
	Token          string          `json:"token"`
	Provider       string          `json:"provider"`
	PublicMetadata json.RawMessage `json:"public_metadata"`
	Label          *string         `json:"label"`
	// Only set in OAuth 2.0 tokens
	Scopes []string `json:"scopes,omitempty"`
	// Only set in OAuth 1.0 tokens
	TokenSecret *string `json:"token_secret,omitempty"`
}

type OAuthAccessTokenList struct {
	APIResource
	OAuthAccessTokens []*OAuthAccessToken `json:"data"`
	TotalCount        int64               `json:"total_count"`
}
