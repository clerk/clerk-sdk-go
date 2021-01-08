package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSessionsService_ListAll_happyPath(t *testing.T) {
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

	got, _ := client.Sessions().ListAll()
	if len(got) != len(want) {
		t.Errorf("Was expecting %d sessions to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestSessionsService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	sessions, err := client.Sessions().ListAll()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if sessions != nil {
		t.Errorf("Was not expecting any sessions to be returned, instead got %v", sessions)
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
