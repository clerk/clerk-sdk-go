// +build integration

package integration

import (
	"github.com/clerkinc/clerk_server_sdk_go/clerk"
	"testing"
)

func TestSMS(t *testing.T) {
	client := createClient()

	dummyPhoneNumber := "12345678"
	message := clerk.SMSMessage{
		Message:       "Go SDK test message",
		ToPhoneNumber: &dummyPhoneNumber,
	}

	_, err := client.SMS().Create(message)
	if err != nil {
		t.Fatalf("SMS.Create returned error: %v", err)
	}
}
