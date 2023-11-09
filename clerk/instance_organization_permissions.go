package clerk

import (
	"net/http"
)

type OrganizationPermission struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type OrganizationPermissionsResponse struct {
	Data       []OrganizationPermission `json:"data"`
	TotalCount int64                    `json:"total_count"`
}

type ListInstanceOrganizationPermissionsParams struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

func (s *InstanceService) ListOrganizationPermissions(params ListInstanceOrganizationPermissionsParams) (*OrganizationPermissionsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationPermissionsUrl)

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)
	req.URL.RawQuery = query.Encode()

	response := &OrganizationPermissionsResponse{}
	_, err := s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
