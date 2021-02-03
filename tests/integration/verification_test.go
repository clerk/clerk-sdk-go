// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
	"net/url"
	"testing"
)

func TestVerification_clientActiveSession(t *testing.T) {
	client := createClient()

	request := buildRequest(nil)

	session, err := client.Verification().Verify(request)
	if err != nil {
		t.Errorf("Was not expecting error, found %v instead", err)
	}

	if session == nil {
		t.Errorf("Was expecting session to be returned")
	}
}

func TestVerification_verifySessionId(t *testing.T) {
	client := createClient()

	sessionId := getEnv(SessionID)
	request := buildRequest(&sessionId)

	session, err := client.Verification().Verify(request)
	if err != nil {
		t.Errorf("Was not expecting error, found %v instead", err)
	}

	if session == nil {
		t.Errorf("Was expecting session to be returned")
	}
}

func TestVerification_returnsClerkErrorForInvalidSessionID(t *testing.T) {
	client := createClientWithKey("invalid_key")

	request := buildRequest(nil)

	session, err := client.Verification().Verify(request)
	if err == nil {
		t.Fatal("Was expecting error")
	}

	if session != nil {
		t.Fatalf("Was not expecting session to be returned, found %v instead", session)
	}

	if _, isClerkErr := err.(*clerk.ErrorResponse); !isClerkErr {
		t.Fatalf("Was expecting a Clerk error response, got %v instead", err)
	}
}

func buildRequest(sessionId *string) *http.Request {
	var request http.Request
	request.Method = "GET"

	url := url.URL{
		Scheme: "http",
		Host:   "host.com",
		Path:   "path",
	}
	request.URL = &url

	// add session token as cookie
	sessionCookie := http.Cookie{
		Name:  clerk.CookieSession,
		Value: getEnv(SessionToken),
	}
	request.Header = make(map[string][]string)
	request.AddCookie(&sessionCookie)

	if sessionId != nil {
		// add session id as query parameter
		query := url.Query()
		query.Add(clerk.QueryParamSessionId, *sessionId)
		url.RawQuery = query.Encode()
	}

	return &request
}
