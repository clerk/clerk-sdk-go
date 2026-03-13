package clerk

// SCIMGroupRoleMappingRole represents the role in a group role mapping.
type SCIMGroupRoleMappingRole struct {
	Object            string   `json:"object"`
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Key               string   `json:"key"`
	Description       *string  `json:"description"`
	Permissions       []string `json:"permissions"`
	IsCreatorEligible bool     `json:"is_creator_eligible"`
	CreatedAt         int64    `json:"created_at"`
	UpdatedAt         int64    `json:"updated_at"`
}

// SCIMGroupRoleMapping represents a mapping from a SCIM group to an organization role.
type SCIMGroupRoleMapping struct {
	APIResource
	Object               string                    `json:"object"`
	ID                   string                    `json:"id"`
	SCIMDirectoryID      string                    `json:"scim_directory_id"`
	SCIMGroupID          string                    `json:"scim_group_id"`
	SCIMGroupDisplayName string                    `json:"scim_group_display_name"`
	Role                 *SCIMGroupRoleMappingRole `json:"role,omitempty"`
	Precedence           int                       `json:"precedence"`
	CreatedAt            int64                     `json:"created_at"`
	UpdatedAt            int64                     `json:"updated_at"`
}

// SCIMGroupRoleMappingList represents a list of SCIM group role mappings.
type SCIMGroupRoleMappingList struct {
	APIResource
	Data       []*SCIMGroupRoleMapping `json:"data"`
	TotalCount int64                   `json:"total_count"`
}

// SCIMGroup represents a SCIM group from the directory provider.
type SCIMGroup struct {
	APIResource
	Object      string `json:"object"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	UpdatedAt   int64  `json:"updated_at"`
}

// SCIMGroupList represents a paginated list of SCIM groups.
type SCIMGroupList struct {
	APIResource
	Data   []*SCIMGroup      `json:"data"`
	Cursor *PaginationCursor `json:"cursor"`
}
