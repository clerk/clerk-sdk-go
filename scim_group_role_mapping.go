package clerk

// SCIMGroupRoleMappingRole is the former name of the role carried by a group
// role mapping.
//
// Deprecated: use [OrganizationRole]. The old standalone struct declared
// Permissions as []string, but the API serializes permissions as objects, so
// it could not decode a real response. This alias resolves to the correct
// shape.
type SCIMGroupRoleMappingRole = OrganizationRole

// SCIMGroupRoleMapping is the former name of [DirectoryGroupRoleMapping].
//
// Deprecated: use [DirectoryGroupRoleMapping]. This is a distinct type rather
// than an alias because it decodes the legacy `scim_*` JSON fields. The API
// emits both spellings with the same values, so either type decodes any
// response.
type SCIMGroupRoleMapping struct {
	APIResource
	Object               string            `json:"object"`
	ID                   string            `json:"id"`
	SCIMDirectoryID      string            `json:"scim_directory_id"`
	SCIMGroupID          string            `json:"scim_group_id"`
	SCIMGroupDisplayName string            `json:"scim_group_display_name"`
	Role                 *OrganizationRole `json:"role,omitempty"`
	Precedence           int               `json:"precedence"`
	CreatedAt            int64             `json:"created_at"`
	UpdatedAt            int64             `json:"updated_at"`
}

// SCIMGroupRoleMappingList is the former name of
// [DirectoryGroupRoleMappingList].
//
// Deprecated: use [DirectoryGroupRoleMappingList].
type SCIMGroupRoleMappingList struct {
	APIResource
	Data       []*SCIMGroupRoleMapping `json:"data"`
	TotalCount int64                   `json:"total_count"`
}

// SCIMGroup is the former name of [DirectoryGroup].
//
// Deprecated: use [DirectoryGroup].
type SCIMGroup = DirectoryGroup

// SCIMGroupList is the former name of [DirectoryGroupList].
//
// Deprecated: use [DirectoryGroupList].
type SCIMGroupList = DirectoryGroupList
