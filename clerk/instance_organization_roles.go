package clerk

import (
	"fmt"
	"net/http"
	"strconv"
)

type InsOrgRole struct {
	Object      string             `json:"object"`
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Key         string             `json:"key"`
	Description string             `json:"description"`
	Permissions []InsOrgPermission `json:"permissions"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
}

type InsOrgRolesResponse struct {
	Data       []InsOrgRole `json:"data"`
	TotalCount int64        `json:"total_count"`
}

// TODO: move this to a separate file once custom permissions endpoints are done
type InsOrgPermission struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateInsOrgRoleParams struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
}

func (s *InstanceService) CreateOrganizationRole(params CreateInsOrgRoleParams) (*InsOrgRole, error) {
	req, _ := s.client.NewRequest(http.MethodPost, OrganizationRolesUrl, &params)

	var orgRole InsOrgRole
	_, err := s.client.Do(req, &orgRole)
	if err != nil {
		return nil, err
	}
	return &orgRole, nil
}

type ListInsOrgRoleParams struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

func (s *InstanceService) ListOrganizationRole(params ListInsOrgRoleParams) (*InsOrgRolesResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationRolesUrl)

	query := req.URL.Query()
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", strconv.Itoa(*params.Offset))
	}
	req.URL.RawQuery = query.Encode()

	var orgRolesResponse *InsOrgRolesResponse
	_, err := s.client.Do(req, &orgRolesResponse)
	if err != nil {
		return nil, err
	}
	return orgRolesResponse, nil
}

func (s *InstanceService) ReadOrganizationRole(orgRoleID string) (*InsOrgRole, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", OrganizationRolesUrl, orgRoleID))
	if err != nil {
		return nil, err
	}

	var orgRole InsOrgRole
	_, err = s.client.Do(req, &orgRole)
	if err != nil {
		return nil, err
	}
	return &orgRole, nil
}

type UpdateInsOrgRoleParams struct {
	Name        *string   `json:"name,omitempty"`
	Key         *string   `json:"key,omitempty"`
	Description *string   `json:"description,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
}

func (s *InstanceService) UpdateOrganizationRole(orgRoleID string, params UpdateInsOrgRoleParams) (*InsOrgRole, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", OrganizationRolesUrl, orgRoleID), &params)

	var orgRole InsOrgRole
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
