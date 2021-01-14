// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"testing"
)

func TestSMS(t *testing.T) {
	client := createClient()

	dummyPhoneNumberId := "12345678"
	message := clerk.SMSMessage{
		Message:       "Go SDK test message",
		PhoneNumberID: dummyPhoneNumberId,
	}

	_, err := client.SMS().Create(message)
	if err != nil {
		t.Fatalf("SMS.Create returned error: %v", err)
	}
}
