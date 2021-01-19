package clerk

import "fmt"

type SessionsService service

type Session struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	ClientID     string `json:"client_id"`
	UserID       string `json:"user_id"`
	Status       string `json:"status"`
	LastActiveAt int64  `json:"last_active_at"`
	ExpireAt     int64  `json:"expire_at"`
	AbandonAt    int64  `json:"abandon_at"`
}

func (s *SessionsService) ListAll() ([]Session, error) {
	sessionsUrl := "sessions"
	req, _ := s.client.NewRequest("GET", sessionsUrl)

	var sessions []Session
	_, err := s.client.Do(req, &sessions)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *SessionsService) Read(sessionId string) (*Session, error) {
	sessionUrl := fmt.Sprintf("%s/%s", SessionsUrl, sessionId)
	req, _ := s.client.NewRequest("GET", sessionUrl)

	var session Session
	_, err := s.client.Do(req, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionsService) Revoke(sessionId string) (*Session, error) {
	sessionUrl := fmt.Sprintf("%s/%s/revoke", SessionsUrl, sessionId)
	req, _ := s.client.NewRequest("POST", sessionUrl)

	var session Session
	_, err := s.client.Do(req, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionsService) Verify(sessionId string, token string) (*Session, error) {
	verifyUrl := fmt.Sprintf("%s/%s/verify", SessionsUrl, sessionId)
	var sessionResponse Session

	err := doVerify(s.client, verifyUrl, token, &sessionResponse)
	if err != nil {
		return nil, err
	}
	return &sessionResponse, nil
}
