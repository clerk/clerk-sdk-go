package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestVerificationService_Verify_useSessionId(t *testing.T) {
	apiToken := "apiToken"
	sessionId := "someSessionId"
	sessionToken := "someSessionToken"
	request := setupRequest(&sessionId, &sessionToken)

	client, mux, _, teardown := setup(apiToken)
	defer teardown()

	expectedResponse := dummySessionJson

	mux.HandleFunc("/sessions/"+sessionId+"/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+apiToken)
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Verification().Verify(request)
	if err != nil {
		t.Errorf("Was not expecting error to be returned, got %v instead", err)
	}

	var want Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestVerificationService_Verify_handleServerErrorWhenUsingSessionId(t *testing.T) {
	sessionId := "someSessionId"
	sessionToken := "someSessionToken"
	request := setupRequest(&sessionId, &sessionToken)

	client, mux, _, teardown := setup("apiToken")
	defer teardown()

	mux.HandleFunc("/sessions/"+sessionId+"/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		w.WriteHeader(400)
	})

	_, err := client.Verification().Verify(request)
	if err == nil {
		t.Errorf("Was expecting error to be returned")
	}
}

func TestVerificationService_Verify_useClientActiveSession(t *testing.T) {
	apiToken := "apiToken"
	sessionToken := "someSessionToken"
	request := setupRequest(nil, &sessionToken)

	client, mux, _, teardown := setup(apiToken)
	defer teardown()

	clientResponseJson := dummyClientResponseJson
	var clientResponse ClientResponse
	_ = json.Unmarshal([]byte(clientResponseJson), &clientResponse)

	sessionJson := dummySessionJson

	mux.HandleFunc("/clients/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+apiToken)
		fmt.Fprint(w, clientResponseJson)
	})

	got, err := client.Verification().Verify(request)
	if err != nil {
		t.Errorf("Was not expecting error to be returned, got %v instead", err)
	}

	var want Session
	_ = json.Unmarshal([]byte(sessionJson), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestVerificationService_Verify_handleServerErrorWhenUsingClientActiveSession(t *testing.T) {
	sessionToken := "someSessionToken"
	request := setupRequest(nil, &sessionToken)

	client, mux, _, teardown := setup("apiToken")
	defer teardown()

	mux.HandleFunc("/clients/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		w.WriteHeader(400)
	})

	_, err := client.Verification().Verify(request)
	if err == nil {
		t.Errorf("Was expecting error to be returned")
	}
}

func TestVerificationService_Verify_noActiveSessionWhenUsingClientActiveSession(t *testing.T) {
	apiToken := "apiToken"
	sessionToken := "someSessionToken"
	request := setupRequest(nil, &sessionToken)

	client, mux, _, teardown := setup(apiToken)
	defer teardown()

	var clientResponse ClientResponse
	_ = json.Unmarshal([]byte(dummyClientResponseJson), &clientResponse)
	clientResponse.LastActiveSessionID = nil

	mux.HandleFunc("/clients/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		jsonResp, _ := json.Marshal(clientResponse)
		fmt.Fprint(w, string(jsonResp))
	})

	_, err := client.Verification().Verify(request)
	if err == nil {
		t.Errorf("Was expecting error to be returned")
	}
}

func TestVerificationService_Verify_activeSessionNotIncludedInSessions(t *testing.T) {
	apiToken := "apiToken"
	sessionToken := "someSessionToken"
	request := setupRequest(nil, &sessionToken)

	client, mux, _, teardown := setup(apiToken)
	defer teardown()

	var clientResponse ClientResponse
	_ = json.Unmarshal([]byte(dummyClientResponseJson), &clientResponse)
	clientResponse.Sessions = make([]*Session, 0)

	mux.HandleFunc("/clients/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		jsonResp, _ := json.Marshal(clientResponse)
		fmt.Fprint(w, string(jsonResp))
	})

	_, err := client.Verification().Verify(request)
	if err == nil {
		t.Errorf("Was expecting error to be returned")
	}
}

func TestVerificationService_Verify_notFailForEmptyRequest(t *testing.T) {
	client, _, _, teardown := setup("apiToken")
	defer teardown()

	session, err := client.Verification().Verify(nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if session != nil {
		t.Errorf("Expected no session to returned, got %v instead", session)
	}
}

func TestVerificationService_Verify_noSessionCookie(t *testing.T) {
	client, _, _, teardown := setup("apiToken")
	defer teardown()

	request := setupRequest(nil, nil)

	session, err := client.Verification().Verify(request)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if session != nil {
		t.Errorf("Expected no session to returned, got %v instead", session)
	}
}

func setupRequest(sessionId, sessionToken *string) *http.Request {
	var request http.Request
	request.Method = "GET"

	url := url.URL{
		Scheme: "http",
		Host:   "host.com",
		Path:   "/path",
	}
	request.URL = &url

	if sessionToken != nil {
		// add session token as cookie
		sessionCookie := http.Cookie{
			Name:  CookieSession,
			Value: *sessionToken,
		}
		request.Header = make(map[string][]string)
		request.AddCookie(&sessionCookie)
	}

	if sessionId != nil {
		// add session id as query parameter
		addQueryParam(&request, QueryParamSessionId, *sessionId)
	}

	return &request
}

func addQueryParam(req *http.Request, key, value string) {
	url := req.URL
	query := url.Query()
	query.Add(key, value)
	url.RawQuery = query.Encode()
}
