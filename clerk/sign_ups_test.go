package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSignUpService_Read_HappyPath(t *testing.T) {
	token := "token"
	signUpID := "someSignUpID"
	expectedResponse := dummySignUpJSON

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sign_ups/"+signUpID, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want SignUp
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.SignUps().Read(signUpID)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestSignUpService_Update_HappyPath(t *testing.T) {
	token := "token"
	signUpID := "someSignUpID"

	customAction := true
	externalID := "eternia"
	payload := UpdateSignUp{
		CustomAction: &customAction,
		ExternalID:   &externalID,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sign_ups/"+signUpID, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummySignUpJSON)
	})

	got, _ := client.SignUps().Update(signUpID, &payload)

	var want SignUp
	_ = json.Unmarshal([]byte(dummySignUpJSON), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, payload)
	}
}

const dummySignUpJSON = `{
	"object": "sign_up_attempt",
	"id": "sua_ekekekekekekek",
	"status": "missing_requirements",
	"required_fields": [
		"email_address",
		"password"
	],
	"optional_fields": [],
	"missing_fields": [
		"email_address",
		"password"
	],
	"unverified_fields": [],
	"verifications": {
		"email_address": null,
		"phone_number": null,
		"web3_wallet": null,
		"external_account": null
	},
	"username": null,
	"email_address": null,
	"phone_number": null,
	"web3_wallet": null,
	"password_enabled": false,
	"first_name": null,
	"last_name": null,
	"unsafe_metadata": {},
	"public_metadata": {},
	"custom_action": false,
	"external_id": null,
	"created_session_id": null,
	"created_user_id": null,
	"abandon_at": 449971200000,
	"identification_requirements": [
		["email_address"],
		[]
	],
	"missing_requirements": [
		"password",
		"email_address"
	],
	"email_address_verification": null,
	"phone_number_verification": null,
	"external_account_strategy": null,
	"external_account_verification": null
}`
