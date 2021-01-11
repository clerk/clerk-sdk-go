package clerk

import "errors"

type SMSService service

type SMSMessage struct {
	Message       string  `json:"message"`
	ToPhoneNumber *string `json:"to_phone_number"`
	PhoneNumberID *string `json:"phone_number_id"`
}

type SMSMessageResponse struct {
	Object          string `json:"object"`
	ID              string `json:"id"`
	FromPhoneNumber string `json:"from_phone_number"`
	Status          string `json:"status"`
	SMSMessage
}

func (s *SMSService) Create(message SMSMessage) (*SMSMessageResponse, error) {
	if message.ToPhoneNumber == nil && message.PhoneNumberID == nil {
		return nil, errors.New("one of ToPhoneNumber or PhoneNumberID must be supplied")
	}

	req, _ := s.client.NewRequest("POST", "sms_messages", &message)

	var smsResponse SMSMessageResponse
	_, err := s.client.Do(req, &smsResponse)
	if err != nil {
		return nil, err
	}
	return &smsResponse, nil
}
