package clerk

// SCIMDirectory represents a SCIM directory resource.
type SCIMDirectory struct {
	APIResource
	Object                  string  `json:"object"`
	ID                      string  `json:"id"`
	Name                    string  `json:"name"`
	EnterpriseConnectionID  *string `json:"enterprise_connection_id"`
	EndpointURL             string  `json:"endpoint_url"`
	Provider                string  `json:"provider"`
	Enabled                 bool    `json:"enabled"`
	GroupRoleMappingEnabled bool    `json:"group_role_mapping_enabled"`
	// CredentialsConfigured is only returned for pull-based directories
	// (e.g. Google Workspace): whether validated provider credentials are
	// stored for the directory.
	CredentialsConfigured *bool             `json:"credentials_configured,omitempty"`
	APIKey                *string           `json:"api_key,omitempty"`
	CreatedAt             int64             `json:"created_at"`
	UpdatedAt             int64             `json:"updated_at"`
	AttributeMapping      map[string]string `json:"attribute_mapping,omitempty"`
}

// SCIMDirectoryList represents a paginated list of SCIM directories.
type SCIMDirectoryList struct {
	APIResource
	SCIMDirectories []*SCIMDirectory `json:"data"`
	TotalCount      int64            `json:"total_count"`
}
