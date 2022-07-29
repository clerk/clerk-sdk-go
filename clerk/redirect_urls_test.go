package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectURLService_Create_happyPath(t *testing.T) {
	token := "token"
	var redirectURLResponse RedirectURLResponse
	_ = json.Unmarshal([]byte(dummyRedirectURLJson), &redirectURLResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/redirect_urls", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyRedirectURLJson)
	})

	got, err := client.RedirectURLs().Create(CreateRedirectURLParams{
		URL: redirectURLResponse.URL,
	})

	assert.Nil(t, err)
	assert.Equal(t, *got, redirectURLResponse)
}

func TestRedirectURLService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.RedirectURLs().Create(CreateRedirectURLParams{
		URL: "example.com",
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestRedirectURLService_ListAll_happyPath(t *testing.T) {
	token := "token"
	var redirectURLResponse RedirectURLResponse
	_ = json.Unmarshal([]byte(dummyRedirectURLJson), &redirectURLResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/redirect_urls", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, "["+dummyRedirectURLJson+"]")
	})

	got, err := client.RedirectURLs().ListAll()

	assert.Nil(t, err)
	assert.Equal(t, got, []*RedirectURLResponse{&redirectURLResponse})
}

func TestRedirectURLService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.RedirectURLs().ListAll()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestRedirectURLService_Delete_happyPath(t *testing.T) {
	token := "token"
	client, mux, _, teardown := setup(token)
	defer teardown()

	id := "some_id"
	mux.HandleFunc("/redirect_urls/"+id, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer "+token)

		response := fmt.Sprintf(`{ "deleted": true, "id": "%v", "object": "user" }`, id)
		fmt.Fprint(w, response)
	})

	got, err := client.RedirectURLs().Delete(id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestRedirectURLService_Delete_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.RedirectURLs().Delete("random_id")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyRedirectURLJson = `{
    "object": "redirect_url",
	"id": "ru_1mvFol71HiKCcypBd6xxg0IpMBN",
	"url": "example.com"
}`
