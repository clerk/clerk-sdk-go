package clerk

import "fmt"

type ClientsService service

type ClientResponse struct {
	Object              string  `json:"object"`
	ID                  string  `json:"id"`
	LastActiveSessionID *string `json:"last_active_session_id"`
	SignInAttemptID     *string `json:"sign_in_attempt_id"`
	SignUpAttemptID     *string `json:"sign_up_attempt_id"`
	Ended               bool    `json:"ended"`
}

func (s *ClientsService) ListAll() ([]ClientResponse, error) {
	clientsUrl := "clients"
	req, _ := s.client.NewRequest("GET", clientsUrl)

	var clients []ClientResponse
	_, err := s.client.Do(req, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *ClientsService) Read(clientId string) (*ClientResponse, error) {
	clientUrl := fmt.Sprintf("clients/%v", clientId)
	req, _ := s.client.NewRequest("GET", clientUrl)

	var clientResponse ClientResponse
	_, err := s.client.Do(req, &clientResponse)
	if err != nil {
		return nil, err
	}
	return &clientResponse, nil
}

func (s *ClientsService) Verify(token string) (*ClientResponse, error) {
	verifyUrl := "clients/verify"
	var clientResponse ClientResponse

	err := doVerify(s.client, verifyUrl, token, &clientResponse)
	if err != nil {
		return nil, err
	}
	return &clientResponse, nil
}
