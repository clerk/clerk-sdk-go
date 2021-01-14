package clerk

type EmailService service

type Email struct {
	FromEmailName  string `json:"from_email_name"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
	EmailAddressID string `json:"email_address_id"`
}

type EmailResponse struct {
	ID             string  `json:"id"`
	Object         string  `json:"object"`
	Status         string  `json:"status,omitempty"`
	ToEmailAddress *string `json:"to_email_address,omitempty"`
	Email
}

// Sends an email.
// Please note that one of ToEmailAddress or EmailAddressID must be supplied.
func (s *EmailService) Create(email Email) (*EmailResponse, error) {
	req, _ := s.client.NewRequest("POST", "emails", &email)

	var emailResponse EmailResponse
	_, err := s.client.Do(req, &emailResponse)
	if err != nil {
		return nil, err
	}
	return &emailResponse, nil
}
