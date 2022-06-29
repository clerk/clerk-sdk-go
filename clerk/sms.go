package clerk

import (
	"encoding/json"
)

type SMSService service

type SMSMessage struct {
	Message       string `json:"message"`
	PhoneNumberID string `json:"phone_number_id"`
}

type SMSMessageResponse struct {
	Object           string          `json:"object"`
	ID               string          `json:"id"`
	FromPhoneNumber  string          `json:"from_phone_number"`
	ToPhoneNumber    *string         `json:"to_phone_number,omitempty"`
	Status           string          `json:"status"`
	DeliveredByClerk bool            `json:"delivered_by_clerk"`
	Data             json.RawMessage `json:"data"`
	SMSMessage
}

func (s *SMSService) Create(message SMSMessage) (*SMSMessageResponse, error) {
	req, _ := s.client.NewRequest("POST", SMSUrl, &message)

	var smsResponse SMSMessageResponse
	_, err := s.client.Do(req, &smsResponse)
	if err != nil {
		return nil, err
	}
	return &smsResponse, nil
}
