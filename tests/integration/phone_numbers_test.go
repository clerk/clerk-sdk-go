//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestPhoneNumbers(t *testing.T) {
	client := createClient()

	users, _ := client.Users().ListAll(clerk.ListAllUsersParams{})
	if users == nil || len(users) == 0 {
		return
	}

	user := users[0]

	// create

	fakePhoneNumber := gofakeit.Phone()

	verified := false
	primary := false

	createPhoneNumberParams := clerk.CreatePhoneNumberParams{
		UserID:      user.ID,
		PhoneNumber: fakePhoneNumber,
		Verified:    &verified,
		Primary:     &primary,
	}

	phoneNumber, err := client.PhoneNumbers().Create(createPhoneNumberParams)
	assert.Nil(t, err)

	assert.Equal(t, "phone_number", phoneNumber.Object)
	assert.Equal(t, fakePhoneNumber, phoneNumber.PhoneNumber)

	// read

	phoneNumber, err = client.PhoneNumbers().Read(phoneNumber.ID)
	assert.Nil(t, err)

	assert.Equal(t, "phone_number", phoneNumber.Object)
	assert.Equal(t, fakePhoneNumber, phoneNumber.PhoneNumber)

	// update

	verified = true
	primary = true

	updatePhoneNumberParams := clerk.UpdatePhoneNumberParams{
		Verified: &verified,
		Primary:  &primary,
	}

	phoneNumber, err = client.PhoneNumbers().Update(phoneNumber.ID, updatePhoneNumberParams)
	assert.Nil(t, err)

	assert.Equal(t, "phone_number", phoneNumber.Object)
	assert.Equal(t, fakePhoneNumber, phoneNumber.PhoneNumber)

	// delete

	deletedResponse, err := client.PhoneNumbers().Delete(phoneNumber.ID)
	assert.Nil(t, err)

	assert.Equal(t, phoneNumber.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)
}
