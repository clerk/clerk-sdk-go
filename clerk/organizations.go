package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type OrganizationsService service

type Organization struct {
	Object                string          `json:"object"`
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Slug                  *string         `json:"slug"`
	LogoURL               *string         `json:"logo_url"`
	MembersCount          *int            `json:"members_count,omitempty"`
	MaxAllowedMemberships int             `json:"max_allowed_memberships"`
	PublicMetadata        json.RawMessage `json:"public_metadata"`
	PrivateMetadata       json.RawMessage `json:"private_metadata,omitempty"`
	CreatedAt             int64           `json:"created_at"`
	UpdatedAt             int64           `json:"updated_at"`
}

type CreateOrganizationParams struct {
	Name                  string          `json:"name"`
	Slug                  *string         `json:"slug,omitempty"`
	CreatedBy             string          `json:"created_by"`
	MaxAllowedMemberships *int            `json:"max_allowed_memberships,omitempty"`
	PublicMetadata        json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata       json.RawMessage `json:"private_metadata,omitempty"`
}

func (s *OrganizationsService) Create(params CreateOrganizationParams) (*Organization, error) {
	req, _ := s.client.NewRequest(http.MethodPost, OrganizationsUrl, &params)

	var organization Organization
	_, err := s.client.Do(req, &organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

type UpdateOrganizationParams struct {
	Name                  *string `json:"name,omitempty"`
	MaxAllowedMemberships *int    `json:"max_allowed_memberships,omitempty"`
}

func (s *OrganizationsService) Update(organizationID string, params UpdateOrganizationParams) (*Organization, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", OrganizationsUrl, organizationID), &params)

	var organization Organization
	_, err := s.client.Do(req, &organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

type UpdateMetadataParams struct {
	PublicMetadata  json.RawMessage `json:"public_metadata"`
	PrivateMetadata json.RawMessage `json:"private_metadata"`
	OrganizationID  string          `json:"-"`
}

func (s *OrganizationsService) UpdateMetadata(params UpdateMetadataParams) (*Organization, error) {

	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s/metadata", OrganizationsUrl, params.OrganizationID), &params)

	var organization Organization
	_, err := s.client.Do(req, &organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (s *OrganizationsService) Delete(organizationID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", OrganizationsUrl, organizationID))

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}

func (s *OrganizationsService) Read(organizationIDOrSlug string) (*Organization, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", OrganizationsUrl, organizationIDOrSlug))
	if err != nil {
		return nil, err
	}

	var organization Organization
	_, err = s.client.Do(req, &organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

type OrganizationsResponse struct {
	Data       []Organization `json:"data"`
	TotalCount int64          `json:"total_count"`
}

type ListAllOrganizationsParams struct {
	Limit               *int
	Offset              *int
	IncludeMembersCount bool
}

func (s *OrganizationsService) ListAll(params ListAllOrganizationsParams) (*OrganizationsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationsUrl)

	query := req.URL.Query()
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.IncludeMembersCount {
		query.Set("include_members_count", strconv.FormatBool(params.IncludeMembersCount))
	}
	req.URL.RawQuery = query.Encode()

	var organizationsResponse *OrganizationsResponse
	_, err := s.client.Do(req, &organizationsResponse)
	if err != nil {
		return nil, err
	}
	return organizationsResponse, nil
}

type OrganizationInvitation struct {
	Object         string          `json:"object"`
	ID             string          `json:"id"`
	EmailAddress   string          `json:"email_address"`
	OrganizationID string          `json:"organization_id"`
	PublicMetadata json.RawMessage `json:"public_metadata"`
	Role           string          `json:"role"`
	Status         string          `json:"status"`
	CreatedAt      int64           `json:"created_at"`
	UpdatedAt      int64           `json:"updated_at"`
}

type CreateOrganizationInvitationParams struct {
	EmailAddress   string          `json:"email_address"`
	InviterUserID  string          `json:"inviter_user_id"`
	OrganizationID string          `json:"organization_id"`
	PublicMetadata json.RawMessage `json:"public_metadata,omitempty"`
	RedirectURL    string          `json:"redirect_url,omitempty"`
	Role           string          `json:"role"`
}

func (s *OrganizationsService) CreateInvitation(params CreateOrganizationInvitationParams) (*OrganizationInvitation, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", OrganizationsUrl, params.OrganizationID, InvitationsURL)
	req, _ := s.client.NewRequest(http.MethodPost, endpoint, &params)
	var organizationInvitation OrganizationInvitation
	_, err := s.client.Do(req, &organizationInvitation)
	return &organizationInvitation, err
}
