// +build integration

package integration

import (
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
	}
}
