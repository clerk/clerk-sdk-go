package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UsersService service

type User struct {
	ID                            string         `json:"id"`
	Object                        string         `json:"object"`
	Username                      *string        `json:"username"`
	FirstName                     *string        `json:"first_name"`
	LastName                      *string        `json:"last_name"`
	Gender                        *string        `json:"gender"`
	Birthday                      *string        `json:"birthday"`
	ProfileImageURL               string         `json:"profile_image_url"`
	ImageURL                      *string        `json:"image_url,omitempty"`
	HasImage                      bool           `json:"has_image"`
	PrimaryEmailAddressID         *string        `json:"primary_email_address_id"`
	PrimaryPhoneNumberID          *string        `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID           *string        `json:"primary_web3_wallet_id"`
	PasswordEnabled               bool           `json:"password_enabled"`
	TwoFactorEnabled              bool           `json:"two_factor_enabled"`
	TOTPEnabled                   bool           `json:"totp_enabled"`
	BackupCodeEnabled             bool           `json:"backup_code_enabled"`
	EmailAddresses                []EmailAddress `json:"email_addresses"`
	PhoneNumbers                  []PhoneNumber  `json:"phone_numbers"`
	Web3Wallets                   []Web3Wallet   `json:"web3_wallets"`
	ExternalAccounts              []interface{}  `json:"external_accounts"`
	SAMLAccounts                  []*SAMLAccount `json:"saml_accounts"`
	PublicMetadata                interface{}    `json:"public_metadata"`
	PrivateMetadata               interface{}    `json:"private_metadata"`
	UnsafeMetadata                interface{}    `json:"unsafe_metadata"`
	LastSignInAt                  *int64         `json:"last_sign_in_at"`
	Banned                        bool           `json:"banned"`
	Locked                        bool           `json:"locked"`
	LockoutExpiresInSeconds       *int64         `json:"lockout_expires_in_seconds"`
	VerificationAttemptsRemaining *int64         `json:"verification_attempts_remaining"`
	ExternalID                    *string        `json:"external_id"`
	CreatedAt                     int64          `json:"created_at"`
	UpdatedAt                     int64          `json:"updated_at"`
	LastActiveAt                  *int64         `json:"last_active_at"`
}

type UserOAuthAccessToken struct {
	ExternalAccountID string          `json:"external_account_id"`
	Object            string          `json:"object"`
	Token             string          `json:"token"`
	Provider          string          `json:"provider"`
	PublicMetadata    json.RawMessage `json:"public_metadata"`
	Label             *string         `json:"label"`
	Scopes            []string        `json:"scopes"`
	TokenSecret       *string         `json:"token_secret"`
}

type IdentificationLink struct {
	IdentType string `json:"type"`
	IdentID   string `json:"id"`
}

type SAMLAccount struct {
	Object         string          `json:"object"`
	ID             string          `json:"id"`
	Provider       string          `json:"provider"`
	Active         bool            `json:"active"`
	EmailAddress   string          `json:"email_address"`
	FirstName      *string         `json:"first_name"`
	LastName       *string         `json:"last_name"`
	ProviderUserID *string         `json:"provider_user_id"`
	PublicMetadata json.RawMessage `json:"public_metadata"`
	Verification   *Verification   `json:"verification"`
}

type CreateUserParams struct {
	EmailAddresses          []string         `json:"email_address,omitempty"`
	PhoneNumbers            []string         `json:"phone_number,omitempty"`
	Web3Wallets             []string         `json:"web3_wallet,omitempty"`
	Username                *string          `json:"username,omitempty"`
	Password                *string          `json:"password,omitempty"`
	FirstName               *string          `json:"first_name,omitempty"`
	LastName                *string          `json:"last_name,omitempty"`
	ExternalID              *string          `json:"external_id,omitempty"`
	UnsafeMetadata          *json.RawMessage `json:"unsafe_metadata,omitempty"`
	PublicMetadata          *json.RawMessage `json:"public_metadata,omitempty"`
	PrivateMetadata         *json.RawMessage `json:"private_metadata,omitempty"`
	PasswordDigest          *string          `json:"password_digest,omitempty"`
	PasswordHasher          *string          `json:"password_hasher,omitempty"`
	SkipPasswordRequirement *bool            `json:"skip_password_requirement,omitempty"`
	SkipPasswordChecks      *bool            `json:"skip_password_checks,omitempty"`
	TOTPSecret              *string          `json:"totp_secret,omitempty"`
	BackupCodes             []string         `json:"backup_codes,omitempty"`
	// Specified in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
}

func (s *UsersService) Create(params CreateUserParams) (*User, error) {
	req, _ := s.client.NewRequest("POST", UsersUrl, &params)
	var user User
	_, err := s.client.Do(req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type ListAllUsersParams struct {
	Limit             *int
	Offset            *int
	EmailAddresses    []string
	ExternalIDs       []string
	PhoneNumbers      []string
	Web3Wallets       []string
	Usernames         []string
	UserIDs           []string
	Query             *string
	LastActiveAtSince *int64
	OrderBy           *string
}

func (s *UsersService) ListAll(params ListAllUsersParams) ([]User, error) {
	req, _ := s.client.NewRequest("GET", UsersUrl)

	s.addUserSearchParamsToRequest(req, params)

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)
	if params.OrderBy != nil {
		query.Add("order_by", *params.OrderBy)
	}
	req.URL.RawQuery = query.Encode()

	var users []User
	_, err := s.client.Do(req, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

type UserCount struct {
	Object     string `json:"object"`
	TotalCount int    `json:"total_count"`
}

func (s *UsersService) Count(params ListAllUsersParams) (*UserCount, error) {
	req, _ := s.client.NewRequest("GET", UsersCountUrl)

	s.addUserSearchParamsToRequest(req, params)

	var userCount UserCount
	_, err := s.client.Do(req, &userCount)
	if err != nil {
		return nil, err
	}
	return &userCount, nil
}

func (s *UsersService) addUserSearchParamsToRequest(r *http.Request, params ListAllUsersParams) {
	query := r.URL.Query()
	if params.EmailAddresses != nil {
		for _, email := range params.EmailAddresses {
			query.Add("email_address", email)
		}
	}
	if params.PhoneNumbers != nil {
		for _, phone := range params.PhoneNumbers {
			query.Add("phone_number", phone)
		}
	}
	if params.ExternalIDs != nil {
		for _, externalID := range params.ExternalIDs {
			query.Add("external_id", externalID)
		}
	}
	if params.Web3Wallets != nil {
		for _, web3Wallet := range params.Web3Wallets {
			query.Add("web3_wallet", web3Wallet)
		}
	}
	if params.Usernames != nil {
		for _, username := range params.Usernames {
			query.Add("username", username)
		}
	}
	if params.UserIDs != nil {
		for _, userID := range params.UserIDs {
			query.Add("user_id", userID)
		}
	}
	if params.Query != nil {
		query.Add("query", *params.Query)
	}
	if params.LastActiveAtSince != nil {
		query.Add("last_active_at_since", strconv.Itoa(int(*params.LastActiveAtSince)))
	}
	r.URL.RawQuery = query.Encode()
}

func (s *UsersService) Read(userId string) (*User, error) {
	userUrl := fmt.Sprintf("%s/%s", UsersUrl, userId)
	req, _ := s.client.NewRequest("GET", userUrl)

	var user User
	_, err := s.client.Do(req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UsersService) Delete(userId string) (*DeleteResponse, error) {
	userUrl := fmt.Sprintf("%s/%s", UsersUrl, userId)
	req, _ := s.client.NewRequest("DELETE", userUrl)

	var delResponse DeleteResponse
	if _, err := s.client.Do(req, &delResponse); err != nil {
		return nil, err
	}
	return &delResponse, nil
}

type UpdateUser struct {
	FirstName              *string     `json:"first_name,omitempty"`
	LastName               *string     `json:"last_name,omitempty"`
	PrimaryEmailAddressID  *string     `json:"primary_email_address_id,omitempty"`
	PrimaryPhoneNumberID   *string     `json:"primary_phone_number_id,omitempty"`
	PrimaryWeb3WalletID    *string     `json:"primary_web3_wallet_id,omitempty"`
	Username               *string     `json:"username,omitempty"`
	ProfileImageID         *string     `json:"profile_image_id,omitempty"`
	ProfileImage           *string     `json:"profile_image,omitempty"`
	Password               *string     `json:"password,omitempty"`
	SkipPasswordChecks     *bool       `json:"skip_password_checks,omitempty"`
	SignOutOfOtherSessions *bool       `json:"sign_out_of_other_sessions,omitempty"`
	ExternalID             *string     `json:"external_id,omitempty"`
	PublicMetadata         interface{} `json:"public_metadata,omitempty"`
	PrivateMetadata        interface{} `json:"private_metadata,omitempty"`
	UnsafeMetadata         interface{} `json:"unsafe_metadata,omitempty"`
	TOTPSecret             *string     `json:"totp_secret,omitempty"`
	BackupCodes            []string    `json:"backup_codes,omitempty"`
	// Specified in RFC3339 format
	CreatedAt *string `json:"created_at,omitempty"`
}

func (s *UsersService) Update(userId string, updateRequest *UpdateUser) (*User, error) {
	userUrl := fmt.Sprintf("%s/%s", UsersUrl, userId)
	req, _ := s.client.NewRequest("PATCH", userUrl, updateRequest)

	var updatedUser User
	_, err := s.client.Do(req, &updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

type UpdateUserMetadata struct {
	PublicMetadata  json.RawMessage `json:"public_metadata"`
	PrivateMetadata json.RawMessage `json:"private_metadata"`
	UnsafeMetadata  json.RawMessage `json:"unsafe_metadata"`
}

func (s *UsersService) UpdateMetadata(userId string, updateMetadataRequest *UpdateUserMetadata) (*User, error) {
	updateUserMetadataURL := fmt.Sprintf("%s/%s/metadata", UsersUrl, userId)
	req, _ := s.client.NewRequest(http.MethodPatch, updateUserMetadataURL, updateMetadataRequest)

	var updatedUser User
	_, err := s.client.Do(req, &updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (s *UsersService) ListOAuthAccessTokens(userID, provider string) ([]*UserOAuthAccessToken, error) {
	url := fmt.Sprintf("%s/%s/oauth_access_tokens/%s", UsersUrl, userID, provider)
	req, _ := s.client.NewRequest(http.MethodGet, url)

	response := make([]*UserOAuthAccessToken, 0)
	_, err := s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type UserDisableMFAResponse struct {
	UserID string `json:"user_id"`
}

func (s *UsersService) DisableMFA(userID string) (*UserDisableMFAResponse, error) {
	url := fmt.Sprintf("%s/%s/mfa", UsersUrl, userID)
	req, _ := s.client.NewRequest(http.MethodDelete, url)

	var response UserDisableMFAResponse
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *UsersService) Ban(userID string) (*User, error) {
	url := fmt.Sprintf("%s/%s/ban", UsersUrl, userID)
	req, _ := s.client.NewRequest(http.MethodPost, url)

	var response User
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *UsersService) Unban(userID string) (*User, error) {
	url := fmt.Sprintf("%s/%s/unban", UsersUrl, userID)
	req, _ := s.client.NewRequest(http.MethodPost, url)

	var response User
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *UsersService) Lock(userID string) (*User, error) {
	url := fmt.Sprintf("%s/%s/lock", UsersUrl, userID)
	req, _ := s.client.NewRequest(http.MethodPost, url)

	var response User
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *UsersService) Unlock(userID string) (*User, error) {
	url := fmt.Sprintf("%s/%s/unlock", UsersUrl, userID)
	req, _ := s.client.NewRequest(http.MethodPost, url)

	var response User
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ListMembershipsParams struct {
	Limit  *int
	Offset *int
	UserID string
}

func (s *UsersService) ListMemberships(params ListMembershipsParams) (*ListOrganizationMembershipsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/organization_memberships", UsersUrl, params.UserID))

	paginationParams := PaginationParams{Limit: params.Limit, Offset: params.Offset}
	query := req.URL.Query()
	addPaginationParams(query, paginationParams)
	req.URL.RawQuery = query.Encode()

	var response ListOrganizationMembershipsResponse
	if _, err := s.client.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
