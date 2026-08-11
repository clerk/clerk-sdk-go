package clerk

// DirectoryGroupRoleMapping represents a mapping from a directory group to an
// organization role.
type DirectoryGroupRoleMapping struct {
	APIResource
	Object                    string            `json:"object"`
	ID                        string            `json:"id"`
	DirectoryID               string            `json:"directory_id"`
	DirectoryGroupID          string            `json:"directory_group_id"`
	DirectoryGroupDisplayName string            `json:"directory_group_display_name"`
	Role                      *OrganizationRole `json:"role,omitempty"`
	Precedence                int               `json:"precedence"`
	CreatedAt                 int64             `json:"created_at"`
	UpdatedAt                 int64             `json:"updated_at"`
}

// DirectoryGroupRoleMappingList represents a list of directory group role mappings.
type DirectoryGroupRoleMappingList struct {
	APIResource
	Data       []*DirectoryGroupRoleMapping `json:"data"`
	TotalCount int64                        `json:"total_count"`
}

// DirectoryGroup represents a group from the directory provider.
type DirectoryGroup struct {
	APIResource
	Object      string `json:"object"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	UpdatedAt   int64  `json:"updated_at"`
}

// DirectoryGroupList represents a paginated list of directory groups.
type DirectoryGroupList struct {
	APIResource
	Data   []*DirectoryGroup `json:"data"`
	Cursor *PaginationCursor `json:"cursor"`
}
