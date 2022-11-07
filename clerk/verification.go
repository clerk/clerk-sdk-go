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

type Verification struct {
	Status           string `json:"status"`
	Strategy         string `json:"strategy"`
	Attempts         *int   `json:"attempts"`
	ExpireAt         *int64 `json:"expire_at"`
	VerifiedAtClient string `json:"verified_at_client,omitempty"`

	// needed for Web3
	Nonce *string `json:"nonce,omitempty"`

	// needed for OAuth
	ExternalVerificationRedirectURL *string `json:"external_verification_redirect_url,omitempty"`
	Error                           []byte  `json:"error,omitempty"`
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

	for _, session := range clientResponse.Sessions {
		if session.ID == *clientResponse.LastActiveSessionID {
			return session, nil
		}
	}

	return nil, errors.New("active session not included in client's sessions")
}

func doVerify(client Client, url, token string, response interface{}) error {
	tokenPayload := verifyRequest{Token: token}
	req, _ := client.NewRequest("POST", url, &tokenPayload)

	_, err := client.Do(req, response)
	return err
}
