//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

type organizationMetadata struct {
	AppID int `json:"app_id"`
}

func TestOrganizations(t *testing.T) {
	client := createClient()

	limit := 2
	users, err := client.Users().ListAll(clerk.ListAllUsersParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("Users.ListAll returned error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("Users.ListAll returned %d results, expected 2", len(users))
	}

	newOrganization, err := client.Organizations().Create(clerk.CreateOrganizationParams{
		Name:      "my-org",
		CreatedBy: users[0].ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, newOrganization.ID)
	assert.Equal(t, "my-org", newOrganization.Name)

	membershipLimit := 20
	updatedOrganization, err := client.Organizations().Update(newOrganization.ID, clerk.UpdateOrganizationParams{
		MaxAllowedMemberships: &membershipLimit,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, membershipLimit, updatedOrganization.MaxAllowedMemberships)

	privateMetadata, err := json.Marshal(
		organizationMetadata{
			AppID: 6,
		},
	)
	publicMetadata, err := json.Marshal(
		organizationMetadata{
			AppID: 2,
		},
	)
	updatedOrganization, err = client.Organizations().Update(newOrganization.ID, clerk.UpdateOrganizationParams{
		PrivateMetadata: privateMetadata,
		PublicMetadata:  publicMetadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, updatedOrganization.PrivateMetadata, json.RawMessage(privateMetadata))
	assert.Equal(t, updatedOrganization.PublicMetadata, json.RawMessage(publicMetadata))

	organizations, err := client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		IncludeMembersCount: true,
	})
	if err != nil {
		t.Fatalf("Organizations.ListAll returned error: %v", err)
	}
	if organizations == nil {
		t.Fatalf("Organizations.ListAll returned nil")
	}

	assert.Greater(t, len(organizations.Data), 0)
	assert.Greater(t, organizations.TotalCount, int64(0))
	for _, organization := range organizations.Data {
		assert.Greater(t, *organization.MembersCount, 0)
	}

	organizationMemberships, err := client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships returned error: %v", err)
	}
	if organizationMemberships == nil {
		t.Fatalf("Organizations.ListMemberships returned nil")
	}
	assert.Greater(t, len(organizationMemberships.Data), 0)
	assert.Greater(t, organizationMemberships.TotalCount, int64(0))
	for _, organizationMembership := range organizationMemberships.Data {
		assert.NotEmpty(t, organizationMembership.ID)
	}

	createdOrganizationMembership, err := client.Organizations().CreateMembership(newOrganization.ID, clerk.CreateOrganizationMembershipParams{
		UserID: users[1].ID,
		Role:   "admin",
	})
	if err != nil {
		t.Fatalf("Organizations.CreateMembership returned error: %v", err)
	}
	if createdOrganizationMembership == nil {
		t.Fatalf("Organizations.CreateMembership returned nil")
	}
	assert.Equal(t, createdOrganizationMembership.Role, "admin")

	updatedOrganizationMembership, err := client.Organizations().UpdateMembership(newOrganization.ID, clerk.UpdateOrganizationMembershipParams{
		UserID: organizationMemberships.Data[0].PublicUserData.UserID,
		Role:   "basic_member",
	})
	if err != nil {
		t.Fatalf("Organizations.UpdateMembership returned error: %v", err)
	}
	if updatedOrganizationMembership == nil {
		t.Fatalf("Organizations.UpdateMembership returned nil")
	}
	assert.Equal(t, updatedOrganizationMembership.Role, "basic_member")

	deleteResponse, err := client.Organizations().Delete(newOrganization.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newOrganization.ID, deleteResponse.ID)
}
