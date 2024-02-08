package clerk

import "encoding/json"

type OrganizationMembership struct {
	APIResource
	Object          string                                `json:"object"`
	ID              string                                `json:"id"`
	Organization    *Organization                         `json:"organization"`
	Permissions     []string                              `json:"permissions"`
	PublicMetadata  json.RawMessage                       `json:"public_metadata"`
	PrivateMetadata json.RawMessage                       `json:"private_metadata"`
	Role            string                                `json:"role"`
	CreatedAt       int64                                 `json:"created_at"`
	UpdatedAt       int64                                 `json:"updated_at"`
	PublicUserData  *OrganizationMembershipPublicUserData `json:"public_user_data,omitempty"`
}

type OrganizationMembershipList struct {
	APIResource
	OrganizationMemberships []*OrganizationMembership `json:"data"`
	TotalCount              int64                     `json:"total_count"`
}

type OrganizationMembershipPublicUserData struct {
	UserID     string  `json:"user_id"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	ImageURL   *string `json:"image_url"`
	HasImage   bool    `json:"has_image"`
	Identifier string  `json:"identifier"`
}
