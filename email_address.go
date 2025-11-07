package clerk

import "encoding/json"

type EmailAddress struct {
	APIResource
	ID                   string                  `json:"id"`
	Object               string                  `json:"object"`
	EmailAddress         string                  `json:"email_address"`
	Reserved             bool                    `json:"reserved"`
	Verification         *Verification           `json:"verification"`
	LinkedTo             []*LinkedIdentification `json:"linked_to"`
	MatchesSSOConnection bool                    `json:"matches_sso_connection"`
	CreatedAt            int64                   `json:"created_at"`
	UpdatedAt            int64                   `json:"updated_at"`
}

type Verification struct {
	Object                          string          `json:"object"`
	Status                          string          `json:"status"`
	Strategy                        string          `json:"strategy"`
	Channel                         *string         `json:"channel,omitempty"`
	Attempts                        *int64          `json:"attempts"`
	ExpireAt                        *int64          `json:"expire_at"`
	VerifiedAtClient                string          `json:"verified_at_client,omitempty"`
	Nonce                           *string         `json:"nonce,omitempty"`
	Message                         *string         `json:"message,omitempty"`
	ExternalVerificationRedirectURL *string         `json:"external_verification_redirect_url,omitempty"`
	Error                           json.RawMessage `json:"error,omitempty"`
}

type LinkedIdentification struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
