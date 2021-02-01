package clerk

import (
	"errors"
	"net/http"
)

const (
	CookieSession       = "__session"
	QueryParamSessionId = "_clerk_session_id"
)

type VerificationService service

type verifyRequest struct {
	Token string `json:"token"`
}

func (s *VerificationService) Verify(req *http.Request) (*Session, error) {
	if req == nil {
		return nil, errors.New("cannot verify empty request")
	}
	cookie, err := req.Cookie(CookieSession)
	if err != nil {
		return nil, errors.New("couldn't find cookie " + CookieSession)
	}

	sessionToken := cookie.Value
	sessionId := req.URL.Query().Get(QueryParamSessionId)

	if sessionId == "" {
		return s.useClientActiveSession(sessionToken)
	}

	return s.client.Sessions().Verify(sessionId, sessionToken)
}

func (s *VerificationService) useClientActiveSession(token string) (*Session, error) {
	clientResponse, err := s.client.Clients().Verify(token)
	if err != nil {
		return nil, err
	}

	if clientResponse.LastActiveSessionID == nil {
		return nil, errors.New("no active sessions for given client")
	}

	return s.client.Sessions().Read(*clientResponse.LastActiveSessionID)
}

func doVerify(client Client, url string, token string, response interface{}) error {
	tokenPayload := verifyRequest{Token: token}
	req, _ := client.NewRequest("POST", url, &tokenPayload)

	_, err := client.Do(req, response)
	return err
}
