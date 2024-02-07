package clerk

import "encoding/json"

type EmailAddress struct {
	APIResource
	ID           string                  `json:"id"`
	Object       string                  `json:"object"`
	EmailAddress string                  `json:"email_address"`
	Reserved     bool                    `json:"reserved"`
	Verification *Verification           `json:"verification"`
	LinkedTo     []*LinkedIdentification `json:"linked_to"`
}

type Verification struct {
	Status                          string          `json:"status"`
	Strategy                        string          `json:"strategy"`
	Attempts                        *int64          `json:"attempts"`
	ExpireAt                        *int64          `json:"expire_at"`
	VerifiedAtClient                string          `json:"verified_at_client,omitempty"`
	Nonce                           *string         `json:"nonce,omitempty"`
	ExternalVerificationRedirectURL *string         `json:"external_verification_redirect_url,omitempty"`
	Error                           json.RawMessage `json:"error,omitempty"`
}

type LinkedIdentification struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
