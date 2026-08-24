package clerk

// RoleSet represents a role set resource.
type RoleSet struct {
	APIResource
	Object      string         `json:"object"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Key         string         `json:"key"`
	Description *string        `json:"description"`
	Roles       []*RoleSetItem `json:"roles"`
	// Type defines the type of role set. It can be either "initial" or "custom".
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// RoleSetItem represents a role within a role set.
type RoleSetItem struct {
	Object      string  `json:"object"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Key         string  `json:"key"`
	Description *string `json:"description"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// RoleSetList represents a list of role sets.
type RoleSetList struct {
	APIResource
	RoleSets   []*RoleSet `json:"data"`
	TotalCount int64      `json:"total_count"`
}
