package clerk

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	CookieSession       = "__session"
	QueryParamSessionId = "_clerk_session_id"
	OriginHeader        = "Origin"
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

	if !isXHR(req) && req.Method == "GET" && sessionId == "" {
		return s.useClientActiveSession(sessionToken)
	}

	if sessionId == "" {
		return nil, errors.New(fmt.Sprintf("no session id is specified via the %s query parameter", QueryParamSessionId))
	}

	return s.client.Sessions().Verify(sessionId, sessionToken)
}

func isXHR(req *http.Request) bool {
	headers := req.Header
	origin := headers.Get(OriginHeader)
	return len(origin) > 0
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
	if err != nil {
		return err
	}
	return nil
}
