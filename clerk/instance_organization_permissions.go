package clerk

import (
	"net/http"
)

type ListInstanceOrganizationPermissionsParams struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

func (s *InstanceService) ListOrganizationPermissions(params ListInstanceOrganizationPermissionsParams) (*PermissionsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationPermissionsUrl)

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)
	req.URL.RawQuery = query.Encode()

	response := &PermissionsResponse{}
	_, err := s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
