// +build integration

package integration

import (
	"github.com/clerkinc/clerk_server_sdk_go/clerk"
	"testing"
)

func TestUsers(t *testing.T) {
	client := createClient()

	users, err := client.Users().ListAll()
	if err != nil {
		t.Fatalf("Users.ListAll returned error: %v", err)
	}
	if users == nil {
		t.Fatalf("Users.ListAll returned nil")
	}

	for _, user := range users {
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
