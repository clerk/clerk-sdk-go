package clerk

import "encoding/json"

type Email struct {
	APIResource
	ID               string          `json:"id"`
	Object           string          `json:"object"`
	Slug             *string         `json:"slug,omitempty"`
	FromEmailName    string          `json:"from_email_name"`
	ReplyToEmailName *string         `json:"reply_to_email_name,omitempty"`
	ToEmailAddress   string          `json:"to_email_address,omitempty"`
	EmailAddressID   *string         `json:"email_address_id,omitempty"`
	UserID           *string         `json:"user_id,omitempty"`
	Subject          string          `json:"subject,omitempty"`
	Body             string          `json:"body,omitempty"`
	BodyPlain        *string         `json:"body_plain,omitempty"`
	Status           string          `json:"status,omitempty"`
	Data             json.RawMessage `json:"data,omitempty"`
	DeliveredByClerk bool            `json:"delivered_by_clerk"`
}
