package clerk

import (
	"fmt"
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

type CreateInstanceOrganizationPermissionParams struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

func (s *InstanceService) CreateOrganizationPermission(params CreateInstanceOrganizationPermissionParams) (*Permission, error) {
	req, _ := s.client.NewRequest(http.MethodPost, OrganizationPermissionsUrl)

	var orgPermission Permission
	_, err := s.client.Do(req, &orgPermission)
	if err != nil {
		return nil, err
	}
	return &orgPermission, nil
}

func (s *InstanceService) ReadOrganizationPermission(orgPermissionID string) (*Permission, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", OrganizationPermissionsUrl, orgPermissionID))
	if err != nil {
		return nil, err
	}

	var orgPermission Permission
	_, err = s.client.Do(req, &orgPermission)
	if err != nil {
		return nil, err
	}
	return &orgPermission, nil
}

type UpdateInstanceOrganizationPermissionParams struct {
	Name        *string `json:"name,omitempty"`
	Key         *string `json:"key,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *InstanceService) UpdateOrganizationPermission(orgPermissionID string, params UpdateInstanceOrganizationPermissionParams) (*Permission, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", OrganizationPermissionsUrl, orgPermissionID), &params)

	var orgPermission Permission
	_, err := s.client.Do(req, &orgPermission)
	if err != nil {
		return nil, err
	}
	return &orgPermission, nil
}

func (s *InstanceService) DeleteOrganizationPermission(orgPermissionID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", OrganizationPermissionsUrl, orgPermissionID))

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}
