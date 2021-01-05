package clerk

import (
	"net/url"
	"testing"
)

func TestNewClientBaseUrl(t *testing.T) {
	c, err := NewClient("token")
	if err != nil {
		t.Errorf("NewClient failed")
	}

	if got, want := c.BaseURL.String(), clerkBaseUrl; got != want {
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

	inputUrl, outputUrl := "/test", clerkBaseUrl + "/test"
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
