package clerk

import "fmt"

type EmailAddressesService service

type EmailAddress struct {
	ID           string               `json:"id"`
	Object       string               `json:"object"`
	EmailAddress string               `json:"email_address"`
	Reserved     bool                 `json:"reserved"`
	Verification *Verification        `json:"verification"`
	LinkedTo     []IdentificationLink `json:"linked_to"`
}

type CreateEmailAddressParams struct {
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	Verified     *bool  `json:"verified"`
	Primary      *bool  `json:"primary"`
}

type UpdateEmailAddressParams struct {
	Verified *bool `json:"verified"`
	Primary  *bool `json:"primary"`
}

func (s *EmailAddressesService) Create(params CreateEmailAddressParams) (*EmailAddress, error) {
	req, _ := s.client.NewRequest("POST", EmailAddressesURL, &params)

	var emailAddress EmailAddress
	_, err := s.client.Do(req, &emailAddress)
	if err != nil {
		return nil, err
	}
	return &emailAddress, nil
}

func (s *EmailAddressesService) Read(emailAddressID string) (*EmailAddress, error) {
	emailAddressURL := fmt.Sprintf("%s/%s", EmailAddressesURL, emailAddressID)
	req, _ := s.client.NewRequest("GET", emailAddressURL)

	var emailAddress EmailAddress
	_, err := s.client.Do(req, &emailAddress)
	if err != nil {
		return nil, err
	}
	return &emailAddress, nil
}

func (s *EmailAddressesService) Update(emailAddressID string, params UpdateEmailAddressParams) (*EmailAddress, error) {
	emailAddressURL := fmt.Sprintf("%s/%s", EmailAddressesURL, emailAddressID)
	req, _ := s.client.NewRequest("PATCH", emailAddressURL, &params)

	var emailAddress EmailAddress
	_, err := s.client.Do(req, &emailAddress)
	if err != nil {
		return nil, err
	}
	return &emailAddress, nil
}

func (s *EmailAddressesService) Delete(emailAddressID string) (*DeleteResponse, error) {
	emailAddressURL := fmt.Sprintf("%s/%s", EmailAddressesURL, emailAddressID)
	req, _ := s.client.NewRequest("DELETE", emailAddressURL)

	var delResponse DeleteResponse
	if _, err := s.client.Do(req, &delResponse); err != nil {
		return nil, err
	}
	return &delResponse, nil
}
