package clerk

import (
	"fmt"
	"net/http"
)

type CreateInstanceOrganizationRoleParams struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
}

func (s *InstanceService) CreateOrganizationRole(params CreateInstanceOrganizationRoleParams) (*Role, error) {
	req, _ := s.client.NewRequest(http.MethodPost, OrganizationRolesUrl, &params)

	var orgRole Role
	_, err := s.client.Do(req, &orgRole)
	if err != nil {
		return nil, err
	}
	return &orgRole, nil
}

type ListInstanceOrganizationRoleParams struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

func (s *InstanceService) ListOrganizationRole(params ListInstanceOrganizationRoleParams) (*RolesResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationRolesUrl)

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)
	req.URL.RawQuery = query.Encode()

	var orgRolesResponse *RolesResponse
	_, err := s.client.Do(req, &orgRolesResponse)
	if err != nil {
		return nil, err
	}
	return orgRolesResponse, nil
}

func (s *InstanceService) ReadOrganizationRole(orgRoleID string) (*Role, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", OrganizationRolesUrl, orgRoleID))
	if err != nil {
		return nil, err
	}

	var orgRole Role
	_, err = s.client.Do(req, &orgRole)
	if err != nil {
		return nil, err
	}
	return &orgRole, nil
}

type UpdateInstanceOrganizationRoleParams struct {
	Name        *string   `json:"name,omitempty"`
	Key         *string   `json:"key,omitempty"`
	Description *string   `json:"description,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
}

func (s *InstanceService) UpdateOrganizationRole(orgRoleID string, params UpdateInstanceOrganizationRoleParams) (*Role, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", OrganizationRolesUrl, orgRoleID), &params)

	var orgRole Role
	_, err := s.client.Do(req, &orgRole)
	if err != nil {
		return nil, err
	}
	return &orgRole, nil
}

func (s *InstanceService) DeleteOrganizationRole(orgRoleID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", OrganizationRolesUrl, orgRoleID))

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}
