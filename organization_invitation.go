package clerk

import "encoding/json"

type OrganizationInvitation struct {
	APIResource
	Object          string          `json:"object"`
	ID              string          `json:"id"`
	EmailAddress    string          `json:"email_address"`
	Role            string          `json:"role"`
	OrganizationID  string          `json:"organization_id"`
	Status          string          `json:"status"`
	PublicMetadata  json.RawMessage `json:"public_metadata"`
	PrivateMetadata json.RawMessage `json:"private_metadata"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
}
