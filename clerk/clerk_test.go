package clerk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient_baseUrl(t *testing.T) {
	c, err := NewClient("token")
	if err != nil {
		t.Errorf("NewClient failed")
	}

	if got, want := c.(*client).baseURL.String(), ProdUrl; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClient_baseUrlWithoutSlash(t *testing.T) {
	input, want := "http://host/v1", "http://host/v1/"
	c, _ := NewClientWithBaseUrl("token", input)

	if got := c.(*client).baseURL.String(); got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClient_createsDifferentClients(t *testing.T) {
	token := "token"
	c, _ := NewClient(token)
	c2, _ := NewClient(token)
	if c.(*client).client == c2.(*client).client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	client, _ := NewClient("token")

	inputUrl, outputUrl := "test", ProdUrl+"test"
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

func TestNewRequest_noBody(t *testing.T) {
	client, _ := NewClient("token")
	req, _ := client.NewRequest("GET", ".")
	if req.Body != nil {
		t.Fatalf("Expected nil Body but request contains a non-nil Body")
	}
}

func TestNewRequest_nilBody(t *testing.T) {
	client, _ := NewClient("token")
	req, _ := client.NewRequest("GET", ".", nil)
	if req.Body != nil {
		t.Fatalf("Expected nil Body but request contains a non-nil Body")
	}
}

func TestNewRequest_withBody(t *testing.T) {
	client, _ := NewClient("token")

	type Foo struct {
		Key string `json:"key"`
	}

	inBody, outBody := Foo{Key: "value"}, `{"key":"value"}`+"\n"
	req, _ := client.NewRequest("GET", ".", inBody)

	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}
}

func TestNewRequest_invalidBody(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.NewRequest("GET", ".", make(chan int))
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
		w.WriteHeader(http.StatusNoContent)
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

func TestDo_handlesClerkErrors(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expected := &ErrorResponse{
		Errors: []Error{{
			Message:     "Error message",
			LongMessage: "Error long message",
			Code:        "error_message",
		}},
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(expected)
		w.Write(data)
	})

	req, _ := client.NewRequest("GET", "test")
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}

	errorResponse, isClerkErr := err.(*ErrorResponse)
	if !isClerkErr {
		t.Fatal("Expected Clerk error response.")
	}
	if errorResponse.Response != nil {
		t.Fatal("Expected error response to contain the HTTP response")
	}
	if !reflect.DeepEqual(errorResponse.Errors, expected.Errors) {
		t.Fatalf("Actual = %v, want %v", errorResponse.Errors, expected.Errors)
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
