package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhooksService_CreateSvix_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummySvixResponseJson

	mux.HandleFunc("/webhooks/svix", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want SvixResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Webhooks().CreateSvix()
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("response = %v, want %v", got, want)
	}
}

func TestWebhooksService_CreateSvix_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	svixResponse, err := client.Webhooks().CreateSvix()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
	if svixResponse != nil {
		t.Errorf("was not expecting any users to be returned, instead got %v", svixResponse)
	}
}

func TestWebhooksService_DeleteSvix_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/webhooks/svix", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer token")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Webhooks().DeleteSvix()
	if err != nil {
		t.Errorf("was not expecting error, found %v instead", err)
	}
}

func TestWebhooksService_DeleteSvix_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	err := client.Webhooks().DeleteSvix()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
}

func TestWebhooksService_RefreshSvixURL_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummySvixResponseJson

	mux.HandleFunc("/webhooks/svix_url", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want SvixResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Webhooks().RefreshSvixURL()
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("response = %v, want %v", got, want)
	}
}

func TestWebhooksService_RefreshSvixURL_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	svixResponse, err := client.Webhooks().RefreshSvixURL()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
	if svixResponse != nil {
		t.Errorf("was not expecting any users to be returned, instead got %v", svixResponse)
	}
}

const dummySvixResponseJson = `{
	"svix_url": "http://example.svix.com"
}`
