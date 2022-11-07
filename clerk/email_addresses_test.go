package clerk

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestEmailAddressesService_Create_HappyPath(t *testing.T) {
	token := "token"
	expectedResponse := unverifiedEmailAddressJSON

	verified := false
	primary := false

	payload := CreateEmailAddressParams{
		UserID:       "user_abcdefg",
		EmailAddress: "banana@cherry.com",
		Verified:     &verified,
		Primary:      &primary,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := "/email_addresses"

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want EmailAddress
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.EmailAddresses().Create(payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestEmailAddressesService_Read_HappyPath(t *testing.T) {
	token := "token"
	emailAddressID := "idn_banana"
	expectedResponse := unverifiedEmailAddressJSON

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/email_addresses/%s", emailAddressID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want EmailAddress
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.EmailAddresses().Read(emailAddressID)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestEmailAddressesService_Update_HappyPath(t *testing.T) {
	token := "token"
	emailAddressID := "idn_banana"
	expectedResponse := verifiedEmailAddressJSON

	verified := true
	primary := true

	payload := UpdateEmailAddressParams{
		Verified: &verified,
		Primary:  &primary,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/email_addresses/%s", emailAddressID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want EmailAddress
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.EmailAddresses().Update(emailAddressID, payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestEmailAddressesService_Delete_HappyPath(t *testing.T) {
	token := "token"
	emailAddressID := "idn_banana"

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/email_addresses/%s", emailAddressID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "id": "%v", "object": "email_address" }`, emailAddressID)
		fmt.Fprint(w, response)
	})

	want := DeleteResponse{ID: emailAddressID, Object: "email_address", Deleted: true}

	got, _ := client.EmailAddresses().Delete(emailAddressID)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

const unverifiedEmailAddressJSON = `{
	"id": "idn_banana",
	"object": "email_address",
	"email_address": "banana@cherry.com",
	"reserved": true,
	"linked_to": []
}`

const verifiedEmailAddressJSON = `{
	"id": "idn_banana",
	"object": "email_address",
	"email_address": "banana@cherry.com",
	"reserved": true,
	"verification": {
		"status": "verified",
		"strategy": "admin",
		"attempts": null,
		"expire_at": null
	},
	"linked_to": []
}`
