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

func TestActorTokens(t *testing.T) {
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

	// create actor token
	createParams := clerk.CreateActorTokenParams{
		UserID: user.ID,
		Actor:  []byte(`{"sub":"my_actor_id"}`),
	}
	actorTokenResponse, err := client.ActorTokens().Create(createParams)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, actorTokenResponse.ID)
	assert.Equal(t, "pending", actorTokenResponse.Status)
	assert.JSONEq(t, string(createParams.Actor), string(actorTokenResponse.Actor))
	assert.Equal(t, createParams.UserID, actorTokenResponse.UserID)

	// revoke the previously created token
	actorTokenResponse, err = client.ActorTokens().Revoke(actorTokenResponse.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "revoked", actorTokenResponse.Status)
	assert.Empty(t, actorTokenResponse.Token)
}
