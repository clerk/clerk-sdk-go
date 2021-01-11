package clerk

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClientBaseUrl(t *testing.T) {
	c, err := NewClient("token")
	if err != nil {
		t.Errorf("NewClient failed")
	}

	if got, want := c.baseURL.String(), clerkBaseUrl; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClientCreatesDifferenceClients(t *testing.T) {
	token := "token"
	c, _ := NewClient(token)
	c2, _ := NewClient(token)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	client, _ := NewClient("token")

	inputUrl, outputUrl := "test", clerkBaseUrl+"test"
	method := "GET"
	req, err := client.NewRequest(method, inputUrl)
	if err != nil {
		t.Errorf("NewRequest(%q, %s) method is generated error %v", inputUrl, method, err)
	}

	if got, want := req.Method, method; got != want {
		t.Errorf("NewRequest(%q, %s) method is %v, want %v", inputUrl, method, got, want)
	}

	if got, want := req.URL.String(), outputUrl; got != want {
		t.Errorf("NewRequest(%q, %s) URL is %v, want %v", inputUrl, method, got, want)
	}
}

func TestNewRequest_invalidUrl(t *testing.T) {
	client, _ := NewClient("token")
	_, err := client.NewRequest("GET", ":")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewRequest_invalidMethod(t *testing.T) {
	client, _ := NewClient("token")
	invalidMethod := "ΠΟΣΤ"
	_, err := client.NewRequest(invalidMethod, "/test")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestDo_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "test")
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_sendsTokenInRequest(t *testing.T) {
	token := "token"
	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		testHeader(t, req, "Authorization", "Bearer "+token)
		w.WriteHeader(204)
	})

	req, _ := client.NewRequest("GET", "test")
	_, err := client.Do(req, nil)
	if err != nil {
		t.Errorf("Was not expecting any errors")
	}
}

func TestDo_invalidServer(t *testing.T) {
	client, _ := NewClientWithBaseUrl("token", "http://dummy_url:1337")

	req, _ := client.NewRequest("GET", "test")

	// No server setup, should result in an error
	_, err := client.Do(req, nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "test")
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestDo_unexpectedHttpError(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	req, _ := client.NewRequest("GET", "test")
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 500 error, got no error.")
	}
	if resp.StatusCode != 500 {
		t.Errorf("Expected HTTP 500 error, got %d status code.", resp.StatusCode)
	}
}

func TestDo_failToReadBody(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		// Lying about the body, telling client the length is 1 but not sending anything back
		w.Header().Set("Content-Length", "1")
	})

	req, _ := client.NewRequest("GET", "test")
	body := new(foo)
	_, err := client.Do(req, body)
	if err == nil {
		t.Fatal("Expected EOF error, got no error.")
	}
}

func TestDo_failToUnmarshalBody(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		fmt.Fprint(w, `{invalid}`)
	})

	req, _ := client.NewRequest("GET", "test")
	body := new(foo)
	_, err := client.Do(req, body)
	if err == nil {
		t.Fatal("Expected JSON encoding error, got no error.")
	}
}
