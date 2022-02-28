//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type addressDetails struct {
	Street string `json:"street"`
	Number string `json:"number"`
}

type userAddress struct {
	Address addressDetails `json:"address"`
}

type userAppID struct {
	AppID int `json:"app_id"`
}

func TestUsers(t *testing.T) {
	client := createClient()

	users, err := client.Users().ListAll(clerk.ListAllUsersParams{})
	if err != nil {
		t.Fatalf("Users.ListAll returned error: %v", err)
	}
	if users == nil {
		t.Fatalf("Users.ListAll returned nil")
	}

	for i, user := range users {
		userId := user.ID
		user, err := client.Users().Read(userId)
		if err != nil {
			t.Fatalf("Users.Read returned error: %v", err)
		}
		if user == nil {
			t.Fatalf("Users.Read returned nil")
		}

		updateRequest := clerk.UpdateUser{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			PublicMetadata: userAddress{Address: addressDetails{
				Street: "Fifth Avenue",
				Number: "890",
			}},
			PrivateMetadata: userAppID{AppID: i},
		}
		updatedUser, err := client.Users().Update(userId, &updateRequest)
		if err != nil {
			t.Fatalf("Users.Update returned error: %v", err)
		}
		if updatedUser == nil {
			t.Errorf("Users.Update returned nil")
		}
	}
}
