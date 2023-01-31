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
	Name                  *string         `json:"name,omitempty"`
	MaxAllowedMemberships *int            `json:"max_allowed_memberships,omitempty"`
	PublicMetadata        json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata       json.RawMessage `json:"private_metadata,omitempty"`
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

type UpdateOrganizationMetadataParams struct {
	PublicMetadata  json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata json.RawMessage `json:"private_metadata,omitempty"`
}

func (s *OrganizationsService) UpdateMetadata(organizationID string, params UpdateOrganizationMetadataParams) (*Organization, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s/metadata", OrganizationsUrl, organizationID), &params)

	var organization Organization
	_, err := s.client.Do(req, &organization)
	return &organization, err
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

type ListOrganizationMembershipsParams struct {
	OrganizationID string
	Limit          *int
	Offset         *int
	Roles          []string `json:"role"`
	UserIDs        []string `json:"user_id"`
	EmailAddresses []string `json:"email_address"`
	PhoneNumbers   []string `json:"phone_number"`
	Usernames      []string `json:"username"`
	Web3Wallets    []string `json:"web3_wallet"`
	OrderBy        *string  `json:"order_by"`
	Query          *string  `json:"query"`
}

type ΟrganizationMembershipPublicUserData struct {
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	ProfileImageURL string  `json:"profile_image_url"`
	Identifier      string  `json:"identifier"`
	UserID          string  `json:"user_id"`
}

type OrganizationMembership struct {
	Object          string          `json:"object"`
	ID              string          `json:"id"`
	PublicMetadata  json.RawMessage `json:"public_metadata"`
	PrivateMetadata json.RawMessage `json:"private_metadata"`
	Role            string          `json:"role"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`

	Organization   *Organization                         `json:"organization"`
	PublicUserData *ΟrganizationMembershipPublicUserData `json:"public_user_data"`
}

type ListOrganizationMembershipsResponse struct {
	Data       []OrganizationMembership `json:"data"`
	TotalCount int64                    `json:"total_count"`
}

func (s *OrganizationsService) addMembersSearchParamsToRequest(r *http.Request, params ListOrganizationMembershipsParams) {
	query := r.URL.Query()
	for _, email := range params.EmailAddresses {
		query.Add("email_address", email)
	}

	for _, phone := range params.PhoneNumbers {
		query.Add("phone_number", phone)
	}

	for _, web3Wallet := range params.Web3Wallets {
		query.Add("web3_wallet", web3Wallet)
	}

	for _, username := range params.Usernames {
		query.Add("username", username)
	}

	for _, userID := range params.UserIDs {
		query.Add("user_id", userID)
	}

	for _, role := range params.Roles {
		query.Add("role", role)
	}

	if params.Query != nil {
		query.Add("query", *params.Query)
	}
	r.URL.RawQuery = query.Encode()
}

func (s *OrganizationsService) ListMemberships(params ListOrganizationMembershipsParams) (*ListOrganizationMembershipsResponse, error) {
	endpoint := fmt.Sprintf("%s/%s/memberships", OrganizationsUrl, params.OrganizationID)
	req, _ := s.client.NewRequest(http.MethodGet, endpoint)

	s.addMembersSearchParamsToRequest(req, ListOrganizationMembershipsParams{
		EmailAddresses: params.EmailAddresses,
		PhoneNumbers:   params.PhoneNumbers,
		Web3Wallets:    params.Web3Wallets,
		Usernames:      params.Usernames,
		UserIDs:        params.UserIDs,
		Roles:          params.Roles,
		Query:          params.Query,
		OrderBy:        params.OrderBy,
	})

	query := req.URL.Query()
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", strconv.Itoa(*params.Offset))
	}
	req.URL.RawQuery = query.Encode()

	var membershipsResponse *ListOrganizationMembershipsResponse
	_, err := s.client.Do(req, &membershipsResponse)
	if err != nil {
		return nil, err
	}
	return membershipsResponse, nil
}

type CreateOrganizationMembershipParams struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (s *OrganizationsService) CreateMembership(organizationID string, params CreateOrganizationMembershipParams) (*OrganizationMembership, error) {
	req, _ := s.client.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/memberships", OrganizationsUrl, organizationID), &params)

	var organizationMembership OrganizationMembership
	_, err := s.client.Do(req, &organizationMembership)
	if err != nil {
		return nil, err
	}
	return &organizationMembership, nil
}

type UpdateOrganizationMembershipParams struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (s *OrganizationsService) UpdateMembership(organizationID string, params UpdateOrganizationMembershipParams) (*OrganizationMembership, error) {
	req, _ := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s/memberships/%s", OrganizationsUrl, organizationID, params.UserID), &params)

	var organizationMembership OrganizationMembership
	_, err := s.client.Do(req, &organizationMembership)
	if err != nil {
		return nil, err
	}
	return &organizationMembership, nil
}

func (s *OrganizationsService) DeleteMembership(organizationID, userID string) (*OrganizationMembership, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/memberships/%s", OrganizationsUrl, organizationID, userID))

	var organizationMembership OrganizationMembership
	_, err := s.client.Do(req, &organizationMembership)
	if err != nil {
		return nil, err
	}
	return &organizationMembership, nil
}
