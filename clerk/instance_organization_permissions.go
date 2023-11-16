package clerk

import (
	"net/http"
)

type ListInstanceOrganizationPermissionsParams struct {
	Limit   *int
	Offset  *int
	Query   *string
	OrderBy *string
}

func (s *InstanceService) ListOrganizationPermissions(params ListInstanceOrganizationPermissionsParams) (*PermissionsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationPermissionsUrl)

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)

	if params.Query != nil {
		query.Set("query", *params.Query)
	}
	if params.OrderBy != nil {
		query.Set("order_by", *params.OrderBy)
	}

	req.URL.RawQuery = query.Encode()

	response := &PermissionsResponse{}
	_, err := s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
