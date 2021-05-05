package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhooksService_CreateDiahook_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyDiahookResponseJson

	mux.HandleFunc("/webhooks/diahook", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want DiahookResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Webhooks().CreateDiahook()
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("response = %v, want %v", got, want)
	}
}

func TestWebhooksService_CreateDiahook_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	diahookResponse, err := client.Webhooks().CreateDiahook()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
	if diahookResponse != nil {
		t.Errorf("was not expecting any users to be returned, instead got %v", diahookResponse)
	}
}

func TestWebhooksService_DeleteDiahook_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/webhooks/diahook", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer token")
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Webhooks().DeleteDiahook()
	if err != nil {
		t.Errorf("was not expecting error, found %v instead", err)
	}
}

func TestWebhooksService_DeleteDiahook_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	err := client.Webhooks().DeleteDiahook()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
}

func TestWebhooksService_RefreshDiahookURL_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyDiahookResponseJson

	mux.HandleFunc("/webhooks/diahook_url", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want DiahookResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Webhooks().RefreshDiahookURL()
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("response = %v, want %v", got, want)
	}
}

func TestWebhooksService_RefreshDiahookURL_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	diahookResponse, err := client.Webhooks().RefreshDiahookURL()
	if err == nil {
		t.Errorf("expected error to be returned")
	}
	if diahookResponse != nil {
		t.Errorf("was not expecting any users to be returned, instead got %v", diahookResponse)
	}
}

const dummyDiahookResponseJson = `{
	"diahook_url": "http://example.diahook.com"
}`
