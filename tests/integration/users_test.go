//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
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
	if userCount.TotalCount == 0 {
		t.Fatalf("Users.Count returned 0, expected %d", len(users))
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

		privateMetadata := userAppAndContactID{ContactID: i}
		privateMetadataJSON, _ := json.Marshal(privateMetadata)

		updatedUser, err = client.Users().UpdateMetadata(userId, &clerk.UpdateUserMetadata{
			PrivateMetadata: privateMetadataJSON,
		})
		if err != nil {
			t.Fatalf("Users.UpdateMetadata returned error: %v", err)
		}
		if updatedUser == nil {
			t.Errorf("Users.UpdateMetadata returned nil")
		}

		updatedUser, err = client.Users().Ban(userId)
		if err != nil {
			t.Fatalf("Users.Ban returned error: %v", err)
		}
		assert.True(t, updatedUser.Banned)

		updatedUser, err = client.Users().Unban(userId)
		if err != nil {
			t.Fatalf("Users.Unban returned error: %v", err)
		}
		assert.False(t, updatedUser.Banned)
	}

	// Should return all memberships of a user
	newOrganization, err := client.Organizations().Create(clerk.CreateOrganizationParams{
		Name:      "my-org",
		CreatedBy: users[0].ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	organizationMemberships, err := client.Users().ListMemberships(clerk.ListMembershipsParams{UserID: users[0].ID})
	assert.Equal(t, len(organizationMemberships.Data), 2)
	assert.Equal(t, organizationMemberships.TotalCount, int64(2))
	assert.Equal(t, newOrganization.ID, organizationMemberships.Data[0].Organization.ID)

	// delete previous created organization to not create conflict with future tests
	deleteResponse, err := client.Organizations().Delete(newOrganization.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newOrganization.ID, deleteResponse.ID)
}
