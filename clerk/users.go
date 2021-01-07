package clerk

import (
	"fmt"
	"time"
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
	Metadata              interface{}    `json:"metadata"`
	PrivateMetadata       interface{}    `json:"private_metadata,omitempty"`
	CreatedAt             *time.Time     `json:"created_at,omitempty"`
	UpdatedAt             *time.Time     `json:"updated_at,omitempty"`
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

func (s *UsersService) ListAll() ([]User, error) {
	usersUrl := "users"
	req, _ := s.client.NewRequest("GET", usersUrl)

	var users []User
	_, err := s.client.Do(req, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UsersService) Read(userId string) (*User, error) {
	userUrl := fmt.Sprintf("users/%v", userId)
	req, _ := s.client.NewRequest("GET", userUrl)

	var user User
	_, err := s.client.Do(req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func (s *UsersService) Delete(userId string) (*DeleteResponse, error) {
	userUrl := fmt.Sprintf("users/%v", userId)
	req, _ := s.client.NewRequest("DELETE", userUrl)

	var delResponse DeleteResponse
	_, err := s.client.Do(req, &delResponse)
	if err != nil {
		return nil, err
	}
	return &delResponse, nil
}
