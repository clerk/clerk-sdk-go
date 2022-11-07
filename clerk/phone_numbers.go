package clerk

import "fmt"

type PhoneNumbersService service

type PhoneNumber struct {
	ID                      string               `json:"id"`
	Object                  string               `json:"object"`
	PhoneNumber             string               `json:"phone_number"`
	ReservedForSecondFactor bool                 `json:"reserved_for_second_factor"`
	DefaultSecondFactor     bool                 `json:"default_second_factor"`
	Reserved                bool                 `json:"reserved"`
	Verification            *Verification        `json:"verification"`
	LinkedTo                []IdentificationLink `json:"linked_to"`
	BackupCodes             []string             `json:"backup_codes"`
}

type CreatePhoneNumberParams struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Verified    *bool  `json:"verified"`
	Primary     *bool  `json:"primary"`
}

type UpdatePhoneNumberParams struct {
	Verified *bool `json:"verified"`
	Primary  *bool `json:"primary"`
}

func (s *PhoneNumbersService) Create(params CreatePhoneNumberParams) (*PhoneNumber, error) {
	req, _ := s.client.NewRequest("POST", PhoneNumbersURL, &params)

	var response PhoneNumber
	_, err := s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *PhoneNumbersService) Read(phoneNumberID string) (*PhoneNumber, error) {
	phoneNumberURL := fmt.Sprintf("%s/%s", PhoneNumbersURL, phoneNumberID)
	req, _ := s.client.NewRequest("GET", phoneNumberURL)

	var phoneNumber PhoneNumber
	_, err := s.client.Do(req, &phoneNumber)
	if err != nil {
		return nil, err
	}
	return &phoneNumber, nil
}

func (s *PhoneNumbersService) Update(phoneNumberID string, params UpdatePhoneNumberParams) (*PhoneNumber, error) {
	phoneNumberURL := fmt.Sprintf("%s/%s", PhoneNumbersURL, phoneNumberID)
	req, _ := s.client.NewRequest("PATCH", phoneNumberURL, &params)

	var phoneNumber PhoneNumber
	_, err := s.client.Do(req, &phoneNumber)
	if err != nil {
		return nil, err
	}
	return &phoneNumber, nil
}

func (s *PhoneNumbersService) Delete(phoneNumberID string) (*DeleteResponse, error) {
	phoneNumberURL := fmt.Sprintf("%s/%s", PhoneNumbersURL, phoneNumberID)
	req, _ := s.client.NewRequest("DELETE", phoneNumberURL)

	var delResponse DeleteResponse
	if _, err := s.client.Do(req, &delResponse); err != nil {
		return nil, err
	}
	return &delResponse, nil
}
