package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"gopkg.in/square/go-jose.v2"
)

func TestWithSession_addSessionToContext(t *testing.T) {
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

func TestWithSession_returnsErrorIfVerificationFails(t *testing.T) {
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

func TestWithNetworkless_addSessionToContext_Header(t *testing.T) {
	c, mux, serverUrl, teardown := setup("apiToken")
	defer teardown()

	expectedClaims := dummyClaims

	token, pubKey := testGenerateTokenJWT(t, expectedClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveClaims`
		activeClaims := r.Context().Value(ActiveClaims)
		_ = json.NewEncoder(w).Encode(activeClaims)
	})

	mux.Handle("/claims", WithNetworkless(c)(dummyHandler))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/claims", serverUrl), nil)
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var got Claims
	_ = json.NewDecoder(resp.Body).Decode(&got)

	if !reflect.DeepEqual(got, expectedClaims) {
		t.Errorf("Response = %v, want %v", got, expectedClaims)
	}
}

func TestWithNetworkless_addSessionToContext_Cookie(t *testing.T) {
	c, mux, serverUrl, teardown := setup("apiToken")
	defer teardown()

	expectedClaims := dummyClaims

	token, pubKey := testGenerateTokenJWT(t, expectedClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveClaims`
		activeClaims := r.Context().Value(ActiveClaims)
		_ = json.NewEncoder(w).Encode(activeClaims)
	})

	mux.Handle("/claims", WithNetworkless(c)(dummyHandler))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/claims", serverUrl), nil)
	req.AddCookie(&http.Cookie{
		Name:     "__session",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var got Claims
	_ = json.NewDecoder(resp.Body).Decode(&got)

	if !reflect.DeepEqual(got, expectedClaims) {
		t.Errorf("Response = %v, want %v", got, expectedClaims)
	}
}

func TestWithNetworkless_returnsErrorIfVerificationFails(t *testing.T) {
	apiToken := "apiToken"

	c, mux, serverUrl, teardown := setup(apiToken)
	defer teardown()

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveSession`
		t.Errorf("This should never be called!")
	})

	mux.Handle("/claims", WithNetworkless(c)(dummyHandler))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/claims", serverUrl), nil)
	req.Header.Set("Authorization", "invalid-token")

	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Was expecting 401 error code, found %v instead", resp.StatusCode)
	}
}

func TestWithNetworkless_returnsErrorIfTokenMissing(t *testing.T) {
	apiToken := "apiToken"

	c, mux, serverUrl, teardown := setup(apiToken)
	defer teardown()

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this handler should be called after the middleware has added the `ActiveSession`
		t.Errorf("This should never be called!")
	})

	mux.Handle("/claims", WithNetworkless(c)(dummyHandler))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/claims", serverUrl), nil)

	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Was expecting 401 error code, found %v instead", resp.StatusCode)
	}
}
