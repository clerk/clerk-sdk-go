package clerk

// OrganizationPermission represents an organization permission resource.
type OrganizationPermission struct {
	APIResource
	Object      string  `json:"object"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Key         string  `json:"key"`
	Description *string `json:"description"`
	Type        string  `json:"type"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// OrganizationPermissionList represents a list of organization permissions.
type OrganizationPermissionList struct {
	APIResource
	OrganizationPermissions []*OrganizationPermission `json:"data"`
	TotalCount              int64                     `json:"total_count"`
}
