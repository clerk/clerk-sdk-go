//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func TestSMS(t *testing.T) {
	client := createClient()

	users, _ := client.Users().ListAll(clerk.ListAllUsersParams{})
	if users == nil || len(users) == 0 {
		return
	}

	user := users[0]

	if user.PrimaryPhoneNumberID == nil {
		return
	}

	message := clerk.SMSMessage{
		Message:       "Go SDK test message",
		PhoneNumberID: *user.PrimaryPhoneNumberID,
	}

	_, err := client.SMS().Create(message)
	if err != nil {
		t.Fatalf("SMS.Create returned error: %v", err)
	}
}
