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

type userAppAndContactID struct {
	AppID     int `json:"app_id"`
	ContactID int `json:"contact_id"`
}

type userEvent struct {
	ViewedProfile bool `json:"viewed_profile"`
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

	userCount, err := client.Users().Count(clerk.ListAllUsersParams{})
	if err != nil {
		t.Fatalf("Users.Count returned error: %v", err)
	}
	if userCount.TotalCount != len(users) {
		t.Fatalf("Users.Count returned %d, expected %d", userCount.TotalCount, len(users))
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
			PrivateMetadata: userAppAndContactID{AppID: i},
			UnsafeMetadata: userEvent{
				ViewedProfile: true,
			},
		}
		updatedUser, err := client.Users().Update(userId, &updateRequest)
		if err != nil {
			t.Fatalf("Users.Update returned error: %v", err)
		}
		if updatedUser == nil {
			t.Errorf("Users.Update returned nil")
		}

		updatedUser, err = client.Users().UpdateMetadata(userId, &clerk.UpdateUserMetadata{
			PrivateMetadata: userAppAndContactID{
				ContactID: i,
			},
		})
		if err != nil {
			t.Fatalf("Users.UpdateMetadata returned error: %v", err)
		}
		if updatedUser == nil {
			t.Errorf("Users.UpdateMetadata returned nil")
		}
	}
}
