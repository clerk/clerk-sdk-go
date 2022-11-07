//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestEmailAddresses(t *testing.T) {
	client := createClient()

	users, _ := client.Users().ListAll(clerk.ListAllUsersParams{})
	if users == nil || len(users) == 0 {
		return
	}

	user := users[0]

	// create

	fakeEmailAddress := gofakeit.Email()

	verified := false
	primary := false

	createEmailAddressParams := clerk.CreateEmailAddressParams{
		UserID:       user.ID,
		EmailAddress: fakeEmailAddress,
		Verified:     &verified,
		Primary:      &primary,
	}

	emailAddress, err := client.EmailAddresses().Create(createEmailAddressParams)
	assert.Nil(t, err)

	assert.Equal(t, "email_address", emailAddress.Object)
	assert.Equal(t, fakeEmailAddress, emailAddress.EmailAddress)

	// read

	emailAddress, err = client.EmailAddresses().Read(emailAddress.ID)
	assert.Nil(t, err)

	assert.Equal(t, "email_address", emailAddress.Object)
	assert.Equal(t, fakeEmailAddress, emailAddress.EmailAddress)

	// update

	verified = true
	primary = true

	updateEmailAddressParams := clerk.UpdateEmailAddressParams{
		Verified: &verified,
		Primary:  &primary,
	}

	emailAddress, err = client.EmailAddresses().Update(emailAddress.ID, updateEmailAddressParams)
	assert.Nil(t, err)

	assert.Equal(t, "email_address", emailAddress.Object)
	assert.Equal(t, fakeEmailAddress, emailAddress.EmailAddress)

	// delete

	deletedResponse, err := client.EmailAddresses().Delete(emailAddress.ID)
	assert.Nil(t, err)

	assert.Equal(t, emailAddress.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)
}
