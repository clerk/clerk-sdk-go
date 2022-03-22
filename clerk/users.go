package clerk

import (
	"fmt"
	"strconv"
)

type UsersService service

type User struct {
	ID                    string         `json:"id"`
	Object                string         `json:"object"`
	Username              *string        `json:"username"`
	FirstName             *string        `json:"first_name"`
	LastName              *string        `json:"last_name"`
	Gender                *string        `json:"gender"`
	Birthday              *string        `json:"birthday"`
	ProfileImageURL       string         `json:"profile_image_url"`
	PrimaryEmailAddressID *string        `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  *string        `json:"primary_phone_number_id"`
	PasswordEnabled       bool           `json:"password_enabled"`
	TwoFactorEnabled      bool           `json:"two_factor_enabled"`
	EmailAddresses        []EmailAddress `json:"email_addresses"`
	PhoneNumbers          []PhoneNumber  `json:"phone_numbers"`
	ExternalAccounts      []interface{}  `json:"external_accounts"`
	PublicMetadata        interface{}    `json:"public_metadata"`
	PrivateMetadata       interface{}    `json:"private_metadata"`
	CreatedAt             int64          `json:"created_at"`
	UpdatedAt             int64          `json:"updated_at"`
}

type EmailAddress struct {
	ID           string               `json:"id"`
	Object       string               `json:"object"`
	EmailAddress string               `json:"email_address"`
	Verification interface{}          `json:"verification"`
	LinkedTo     []IdentificationLink `json:"linked_to"`
}

type PhoneNumber struct {
	ID                      string               `json:"id"`
	Object                  string               `json:"object"`
	PhoneNumber             string               `json:"phone_number"`
	ReservedForSecondFactor bool                 `json:"reserved_for_second_factor"`
	Verification            interface{}          `json:"verification"`
	LinkedTo                []IdentificationLink `json:"linked_to"`
}

type IdentificationLink struct {
	IdentType string `json:"type"`
	IdentID   string `json:"id"`
}

type ListAllUsersParams struct {
	Limit          *int
	Offset         *int
	EmailAddresses []string
	PhoneNumbers   []string
	Web3Wallets    []string
	Usernames      []string
	UserIDs        []string
	Query          *string
	OrderBy        *string
}

func (s *UsersService) ListAll(params ListAllUsersParams) ([]User, error) {
	req, _ := s.client.NewRequest("GET", UsersUrl)

	query := req.URL.Query()
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", strconv.Itoa(*params.Offset))
	}
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
	FirstName             *string     `json:"first_name,omitempty"`
	LastName              *string     `json:"last_name,omitempty"`
	PrimaryEmailAddressID *string     `json:"primary_email_address_id,omitempty"`
	PrimaryPhoneNumberID  *string     `json:"primary_phone_number_id,omitempty"`
	ProfileImage          *string     `json:"profile_image,omitempty"`
	Password              *string     `json:"password,omitempty"`
	PublicMetadata        interface{} `json:"public_metadata,omitempty"`
	PrivateMetadata       interface{} `json:"private_metadata,omitempty"`
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
