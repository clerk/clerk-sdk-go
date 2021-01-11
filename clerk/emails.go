package clerk

import "errors"

type EmailService service

type Email struct {
	FromEmailName  string  `json:"from_email_name"`
	Subject        string  `json:"subject"`
	Body           string  `json:"body"`
	ToEmailAddress *string `json:"to_email_address,omitempty"`
	EmailAddressID *string `json:"email_address_id,omitempty"`
}

type EmailResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Status string `json:"status,omitempty"`
	Email
}

// Sends an email.
// Please note that one of ToEmailAddress or EmailAddressID must be supplied.
func (s *EmailService) Create(email Email) (*EmailResponse, error) {
	if email.ToEmailAddress == nil && email.EmailAddressID == nil {
		return nil, errors.New("one of ToEmailAddress or EmailAddressID must be supplied")
	}

	req, _ := s.client.NewRequest("POST", "emails", &email)

	var emailResponse EmailResponse
	_, err := s.client.Do(req, &emailResponse)
	if err != nil {
		return nil, err
	}
	return &emailResponse, nil
}
