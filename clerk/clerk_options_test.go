package clerk

import (
	"net/http"
	"testing"
	"time"
)

func TestWithHTTPClient(t *testing.T) {
	expectedHTTPClient := &http.Client{Timeout: time.Second * 10}

	got, err := NewClient("token", WithHTTPClient(expectedHTTPClient))
	if err != nil {
		t.Fatal(err)
	}

	if got.(*client).client != expectedHTTPClient {
		t.Fatalf("Expected the http client to have been overriden")
	}
}

func TestWithHTTPClientNil(t *testing.T) {
	_, err := NewClient("token", WithHTTPClient(nil))
	if err == nil {
		t.Fatalf("Expected an error with a nil http client provided")
	}
}

func TestWithBaseURL(t *testing.T) {
	expectedBaseURL := "https://api.example.com/"

	got, err := NewClient("token", WithBaseURL(expectedBaseURL))
	if err != nil {
		t.Fatal(err)
	}

	if got.(*client).baseURL.String() != expectedBaseURL {
		t.Fatalf("Expected the base URL to have been overriden")
	}
}

func TestWithBaseURLEmpty(t *testing.T) {
	_, err := NewClient("token", WithBaseURL(""))
	if err == nil {
		t.Fatalf("Expected an error with an empty base URL provided")
	}
}

func TestWithBaseURLInvalid(t *testing.T) {
	invalidBaseURL := "https:// api.example.com"

	_, err := NewClient("token", WithBaseURL(invalidBaseURL))
	if err == nil {
		t.Fatalf("Expected an error with an invalid base URL provided")
	}
}
