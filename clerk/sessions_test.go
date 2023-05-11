package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSessionsService_ListAll_happyPath_noParams(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := "[" + dummySessionJson + "]"

	mux.HandleFunc("/sessions", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want []Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Sessions().ListAll(ListAllSessionsParams{})
	if len(got) != len(want) {
		t.Errorf("Was expecting %d sessions to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestSessionsService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	sessions, err := client.Sessions().ListAll(ListAllSessionsParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if sessions != nil {
		t.Errorf("Was not expecting any sessions to be returned, instead got %v", sessions)
	}
}

func TestSessionsService_Read_happyPath(t *testing.T) {
	token := "token"
	sessionId := "someSessionId"
	expectedResponse := dummySessionJson

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sessions/"+sessionId, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Sessions().Read(sessionId)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestSessionsService_Read_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	session, err := client.Sessions().Read("someSessionId")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if session != nil {
		t.Errorf("Was not expecting any session to be returned, instead got %v", session)
	}
}

func TestSessionsService_Revoke_happyPath(t *testing.T) {
	token := "token"
	sessionId := "someSessionId"
	expectedResponse := dummySessionJson

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sessions/"+sessionId+"/revoke", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Sessions().Revoke(sessionId)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestSessionsService_Revoke_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	session, err := client.Sessions().Revoke("someSessionId")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if session != nil {
		t.Errorf("Was not expecting any session to be returned, instead got %v", session)
	}
}

func TestSessionsService_Verify_happyPath(t *testing.T) {
	token := "token"
	sessionId := "someSessionId"
	sessionToken := "someSessionToken"
	expectedResponse := dummySessionJson

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sessions/"+sessionId+"/verify", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Sessions().Verify(sessionId, sessionToken)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestSessionsService_Verify_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	session, err := client.Sessions().Verify("someSessionId", "someSessionToken")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if session != nil {
		t.Errorf("Was not expecting any session to be returned, instead got %v", session)
	}
}

const dummySessionJson = `{
        "abandon_at": 1612448988,
        "client_id": "client_1mebPYz8NFNA17fi7NemNXIwp1p",
        "expire_at": 1610461788,
        "id": "sess_1mebQdHlQI14cjxln4e2eXNzwzi",
        "last_active_at": 1609857251,
        "object": "session",
        "status": "ended",
        "user_id": "user_1mebQggrD3xO5JfuHk7clQ94ysA"
    }`
