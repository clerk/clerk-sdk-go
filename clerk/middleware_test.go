package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMiddleware_addSessionToContext(t *testing.T) {
	apiToken := "apiToken"
	sessionId := "someSessionId"
	sessionToken := "someSessionToken"

	client, mux, serverUrl, teardown := setup(apiToken)
	defer teardown()

	expectedResponse := dummySessionJson

	mux.HandleFunc("/sessions/"+sessionId+"/verify", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, expectedResponse)
	})

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveSession`
		activeSession := r.Context().Value(ActiveSession)
		resp, _ := json.Marshal(activeSession)
		fmt.Fprint(w, string(resp))
	})

	mux.Handle("/session", WithSession(client)(dummyHandler))

	request := setupRequest(&sessionId, &sessionToken)
	request.URL.Host = serverUrl.Host
	request.URL.Path = "/v1/session"

	var got Session
	_, _ = client.Do(request, &got)

	var want Session
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestMiddleware_returnsErrorIfVerificationFails(t *testing.T) {
	apiToken := "apiToken"
	sessionId := "someSessionId"
	sessionToken := "someSessionToken"

	client, mux, serverUrl, teardown := setup(apiToken)
	defer teardown()

	mux.HandleFunc("/sessions/"+sessionId+"/verify", func(w http.ResponseWriter, req *http.Request) {
		// return error
		w.WriteHeader(400)
	})

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveSession`
		t.Errorf("This should never be called!")
	})

	mux.Handle("/session", WithSession(client)(dummyHandler))

	request := setupRequest(&sessionId, &sessionToken)
	request.URL.Host = serverUrl.Host
	request.URL.Path = "/v1/session"

	resp, _ := client.Do(request, nil)

	if resp.StatusCode != 400 {
		t.Errorf("Was expecting 400 error code, found %v instead", resp.StatusCode)
	}
}
