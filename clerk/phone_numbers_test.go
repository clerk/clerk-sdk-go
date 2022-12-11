package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNumbersService_Create_HappyPath(t *testing.T) {
	token := "token"
	expectedResponse := unverifiedPhoneNumberJSON

	verified := false
	primary := false

	payload := CreatePhoneNumberParams{
		UserID:      "user_abcdefg",
		PhoneNumber: "+15555555555",
		Verified:    &verified,
		Primary:     &primary,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := "/phone_numbers"

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want PhoneNumber
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.PhoneNumbers().Create(payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestPhoneNumbersService_Read_HappyPath(t *testing.T) {
	token := "token"
	phoneNumberID := "idn_banana"
	expectedResponse := unverifiedPhoneNumberJSON

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/phone_numbers/%s", phoneNumberID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want PhoneNumber
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.PhoneNumbers().Read(phoneNumberID)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestPhoneNumbersService_Update_HappyPath(t *testing.T) {
	token := "token"
	phoneNumberID := "idn_banana"
	expectedResponse := verifiedPhoneNumberJSON

	verified := true
	primary := true

	payload := UpdatePhoneNumberParams{
		Verified: &verified,
		Primary:  &primary,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/phone_numbers/%s", phoneNumberID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want PhoneNumber
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.PhoneNumbers().Update(phoneNumberID, payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestPhoneNumbersService_Delete_HappyPath(t *testing.T) {
	token := "token"
	phoneNumberID := "idn_avocado"

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/phone_numbers/%s", phoneNumberID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "id": "%v", "object": "phone_number" }`, phoneNumberID)
		fmt.Fprint(w, response)
	})

	want := DeleteResponse{ID: phoneNumberID, Object: "phone_number", Deleted: true}

	got, _ := client.PhoneNumbers().Delete(phoneNumberID)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

const unverifiedPhoneNumberJSON = `{
	"id": "idn_avocado",
	"object": "phone_number",
	"phone_number": "+15555555555",
	"reserved_for_second_factor": false,
	"default_second_factor": false,
	"reserved": false,
	"linked_to": [],
	"backup_codes": null
}`

const verifiedPhoneNumberJSON = `{
	"id": "idn_avocado",
	"object": "phone_number",
	"phone_number": "+15555555555",
	"reserved_for_second_factor": false,
	"default_second_factor": false,
	"reserved": false,
	"verification": {
		"status": "verified",
		"strategy": "admin",
		"attempts": null,
		"expire_at": null
	},
	"linked_to": [],
	"backup_codes": null
}`
