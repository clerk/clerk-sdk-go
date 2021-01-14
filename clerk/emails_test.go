package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEmailService_Create_happyPath(t *testing.T) {
	token := "token"
	var email Email
	_ = json.Unmarshal([]byte(dummyEmailJson), &email)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyEmailJson)
	})

	got, _ := client.Emails().Create(email)

	var want EmailResponse
	_ = json.Unmarshal([]byte(dummyEmailJson), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, email)
	}
}

func TestEmailService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	var email Email
	_ = json.Unmarshal([]byte(dummyEmailJson), &email)

	_, err := client.Emails().Create(email)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyEmailJson = `{
    "body": "This is the body of a test email",
    "email_address_id": "idn_1mebQ9KkZWrhb9rL6iEiXQGF8Yj",
    "from_email_name": "info",
    "id": "ema_1mvFol71HiKCcypBd6xxg0IpMBN",
    "object": "email",
    "status": "queued",
    "subject": "This is a test email"
}`
