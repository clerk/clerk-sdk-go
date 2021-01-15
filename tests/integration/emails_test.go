// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"reflect"
	"testing"
)

func TestEmails(t *testing.T) {
	client := createClient()

	users, _ := client.Users().ListAll()
	if users == nil || len(users) == 0 {
		return
	}

	user := users[0]

	email := clerk.Email{
		FromEmailName:  "integration-test",
		Subject:        "Testing Go SDK",
		Body:           "Testing email functionality for Go SDK",
		EmailAddressID: user.PrimaryEmailAddressID,
	}

	got, err := client.Emails().Create(email)
	if err != nil {
		t.Fatalf("Emails.Create returned error: %v", err)
	}

	want := clerk.EmailResponse{
		ID:     got.ID,
		Object: "email",
		Status: "queued",
		Email: clerk.Email{
			FromEmailName:  email.FromEmailName,
			Subject:        email.Subject,
			Body:           email.Body,
			ToEmailAddress: got.ToEmailAddress,
			EmailAddressID: email.EmailAddressID,
		},
	}

	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("Emails.Create(%v) got: %v, wanted %v", email, got, want)
	}
}
