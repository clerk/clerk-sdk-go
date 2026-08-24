package clerk

// OrganizationRole represents an organization role resource.
type OrganizationRole struct {
	APIResource
	Object            string                    `json:"object"`
	ID                string                    `json:"id"`
	Name              string                    `json:"name"`
	Key               string                    `json:"key"`
	Description       *string                   `json:"description"`
	Permissions       []*OrganizationPermission `json:"permissions"`
	IsCreatorEligible bool                      `json:"is_creator_eligible"`
	CreatedAt         int64                     `json:"created_at"`
	UpdatedAt         int64                     `json:"updated_at"`
}

// OrganizationRoleList represents a list of organization roles.
type OrganizationRoleList struct {
	APIResource
	OrganizationRoles []*OrganizationRole `json:"data"`
	TotalCount        int64               `json:"total_count"`
}
