//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestImpersonationTokens(t *testing.T) {
	client := createClient()

	password := "my-extremely_Str0ng_P445wOrd"
	firstName := "John"
	lastName := "Doe"
	user, err := client.Users().Create(clerk.CreateUserParams{
		EmailAddresses: []string{
			fmt.Sprintf("email_%d@example.com", time.Now().Unix()),
		},
		Password:  &password,
		FirstName: &firstName,
		LastName:  &lastName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create impersonation token
	createParams := clerk.CreateImpersonationTokenParams{
		SubjectID: user.ID,
		ActorID:   "my_actor_id",
	}
	impersonationTokenResponse, err := client.ImpersonationTokens().Create(createParams)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, impersonationTokenResponse.ID)
	assert.Equal(t, "pending", impersonationTokenResponse.Status)
	assert.Equal(t, createParams.ActorID, impersonationTokenResponse.ActorID)
	assert.Equal(t, createParams.SubjectID, impersonationTokenResponse.SubjectID)

	// revoke the previously created token
	impersonationTokenResponse, err = client.ImpersonationTokens().Revoke(impersonationTokenResponse.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "revoked", impersonationTokenResponse.Status)
	assert.Empty(t, impersonationTokenResponse.Token)
}
