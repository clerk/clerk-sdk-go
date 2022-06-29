//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func TestEmails(t *testing.T) {
	client := createClient()

	users, _ := client.Users().ListAll(clerk.ListAllUsersParams{})
	if users == nil || len(users) == 0 {
		return
	}

	user := users[0]

	if user.PrimaryEmailAddressID == nil {
		return
	}

	email := clerk.Email{
		FromEmailName:  "integration-test",
		Subject:        "Testing Go SDK",
		Body:           "Testing email functionality for Go SDK",
		EmailAddressID: *user.PrimaryEmailAddressID,
	}

	emailResponse, err := client.Emails().Create(email)
	if err != nil {
		t.Fatalf("Emails.Create returned error: %v", err)
	}

	assert.Equal(t, "email", emailResponse.Object)
	assert.Equal(t, "queued", emailResponse.Status)
	assert.Equal(t, email.FromEmailName, emailResponse.FromEmailName)
	assert.Equal(t, email.EmailAddressID, emailResponse.EmailAddressID)
	assert.Equal(t, email.Subject, emailResponse.Subject)
	assert.Equal(t, email.Body, emailResponse.Body)
	assert.True(t, emailResponse.DeliveredByClerk)
	// assert.Nil(t, emailResponse.Data)
}
