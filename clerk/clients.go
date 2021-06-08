package clerk

import "fmt"

type ClientsService service

type ClientResponse struct {
	Object              string     `json:"object"`
	ID                  string     `json:"id"`
	LastActiveSessionID *string    `json:"last_active_session_id"`
	SessionIDs          []string   `json:"session_ids"`
	Sessions            []*Session `json:"sessions"`
	SignInID            *string    `json:"sign_in_id"`
	SignUpID            *string    `json:"sign_up_id"`
	Ended               bool       `json:"ended"`
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
	clientUrl := fmt.Sprintf("%s/%s", ClientsUrl, clientId)
	req, _ := s.client.NewRequest("GET", clientUrl)

	var clientResponse ClientResponse
	_, err := s.client.Do(req, &clientResponse)
	if err != nil {
		return nil, err
	}
	return &clientResponse, nil
}

func (s *ClientsService) Verify(token string) (*ClientResponse, error) {
	var clientResponse ClientResponse

	err := doVerify(s.client, ClientsVerifyUrl, token, &clientResponse)
	if err != nil {
		return nil, err
	}
	return &clientResponse, nil
}
