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

	sessionId := "sess_1n8u0DZwnk3N4X9iDZ4eSxPfwQ0"
	request := buildRequest(&sessionId)

	session, err := client.Verification().Verify(request)
	if err != nil {
		t.Errorf("Was not expecting error, found %v instead", err)
	}

	if session == nil {
		t.Errorf("Was expecting session to be returned")
	}
}

const dummySessionToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXYiOiJkdmJfMW44dHNvcWh6QkdBSTAwNWRVSXQ0YUNPbXhwIiwiaWQiOiJjbGllbnRfMW44dHpCZzZPMFV0bDIzQUxzSnNPOHd6UGltIiwicm90YXRpbmdfdG9rZW4iOiI5ZjJ2c2pocGhtb3J3enZ0ZXRhNGt5YjVjODViMDZtYzg0emI4ZHM1In0.03WXffBCv8MmRRqx_GzfKcpo8tPQ12gx8eZ_v7w2thZJmLPL-5N_-dFSDICpP2_0gOErnl9NKhFDjgC6TPpery0sLzNqSpxc12gb0VWJlYP78vI5daF8aqhUzefTY3buwFsSLJrdm1nxQViWBGc7HEJGsLvD5sYZxF1OYzi7eOy2k-k6ASFAg3CGJnjBmZy-jhUzgjntbflekfaSfBeYgWpWsuVXwieZ6ZdEeaypTQ20VZzAvhFX-TZocyf3p2aIS9yZRf5OL2n5p8Vy8yLoxhrkP9TWf7k6et7O3gKLaq4ym7cZL_iFOOiMOAIZcXWTfe37F6VA4qaNsUaVvd74Sg"

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
		Value: dummySessionToken,
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
