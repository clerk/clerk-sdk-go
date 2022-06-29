package clerk

import (
	"encoding/json"
)

type EmailService service

type Email struct {
	FromEmailName  string `json:"from_email_name"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
	EmailAddressID string `json:"email_address_id"`
}

type EmailResponse struct {
	ID               string          `json:"id"`
	Object           string          `json:"object"`
	Status           string          `json:"status,omitempty"`
	ToEmailAddress   *string         `json:"to_email_address,omitempty"`
	DeliveredByClerk bool            `json:"delivered_by_clerk"`
	Data             json.RawMessage `json:"data"`
	Email
}

func (s *EmailService) Create(email Email) (*EmailResponse, error) {
	req, _ := s.client.NewRequest("POST", EmailsUrl, &email)

	var emailResponse EmailResponse
	_, err := s.client.Do(req, &emailResponse)
	if err != nil {
		return nil, err
	}
	return &emailResponse, nil
}
