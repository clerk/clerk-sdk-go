package clerk

import "encoding/json"

type publicOrganizationData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ImageURL string `json:"image_url,omitempty"`
	HasImage bool   `json:"has_image"`
}

type OrganizationInvitation struct {
	APIResource
	Object                 string                  `json:"object"`
	ID                     string                  `json:"id"`
	EmailAddress           string                  `json:"email_address"`
	Role                   string                  `json:"role"`
	OrganizationID         string                  `json:"organization_id"`
	PublicOrganizationData *publicOrganizationData `json:"public_organization_data,omitempty"`
	Status                 string                  `json:"status"`
	PublicMetadata         json.RawMessage         `json:"public_metadata"`
	PrivateMetadata        json.RawMessage         `json:"private_metadata"`
	CreatedAt              int64                   `json:"created_at"`
	UpdatedAt              int64                   `json:"updated_at"`
}

type OrganizationInvitationList struct {
	APIResource
	OrganizationInvitations []*OrganizationInvitation `json:"data"`
	TotalCount              int64                     `json:"total_count"`
}
