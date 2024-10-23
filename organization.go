package clerk

import "encoding/json"

type Organization struct {
	APIResource
	Object                           string          `json:"object"`
	ID                               string          `json:"id"`
	Name                             string          `json:"name"`
	Slug                             string          `json:"slug"`
	ImageURL                         *string         `json:"image_url"`
	HasImage                         bool            `json:"has_image"`
	MembersCount                     *int64          `json:"members_count,omitempty"`
	HasMemberWithElevatedPermissions *bool           `json:"has_member_with_elevated_permissions,omitempty"`
	PendingInvitationsCount          *int64          `json:"pending_invitations_count,omitempty"`
	MaxAllowedMemberships            int64           `json:"max_allowed_memberships"`
	AdminDeleteEnabled               bool            `json:"admin_delete_enabled"`
	PublicMetadata                   json.RawMessage `json:"public_metadata"`
	PrivateMetadata                  json.RawMessage `json:"private_metadata"`
	CreatedBy                        string          `json:"created_by"`
	CreatedAt                        int64           `json:"created_at"`
	UpdatedAt                        int64           `json:"updated_at"`
}

type OrganizationList struct {
	APIResource
	Organizations []*Organization `json:"data"`
	TotalCount    int64           `json:"total_count"`
}
