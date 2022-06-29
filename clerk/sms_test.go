package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSMSService_Create_happyPath(t *testing.T) {
	token := "token"
	var message SMSMessage
	_ = json.Unmarshal([]byte(dummySMSMessageResponseJson), &message)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/sms_messages", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummySMSMessageResponseJson)
	})

	got, _ := client.SMS().Create(message)

	var want SMSMessageResponse
	_ = json.Unmarshal([]byte(dummySMSMessageResponseJson), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, message)
	}
}

func TestSMSService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	var message SMSMessage
	_ = json.Unmarshal([]byte(dummySMSMessageResponseJson), &message)

	_, err := client.SMS().Create(message)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummySMSMessageResponseJson = `{
    "from_phone_number": "12345678",
    "id": "sms_1mvjlpmFtRaoee3pm7lS8c3NuAX",
    "message": "This is a test message",
    "object": "sms_message",
    "phone_number_id": "idn_1mebQ9KkZWrhb9rL6iEiXQGF8Yj",
    "status": "queued",
    "to_phone_number": "87654321",
	"data": { "baz": "xyz" },
	"delivered_by_clerk": true
}`
