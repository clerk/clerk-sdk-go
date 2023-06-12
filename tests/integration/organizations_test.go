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

	// Delete all pre existing org before running the tests.
	deleteAllOrgs, err := client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		IncludeMembersCount: true,
	})
	for _, organization := range deleteAllOrgs.Data {
		_, err := client.Organizations().Delete(organization.ID)
		assert.NoError(t, err)
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

	slug := "my-org-slug"
	updatedOrganization, err = client.Organizations().Update(newOrganization.ID, clerk.UpdateOrganizationParams{
		Slug: &slug,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, slug, *updatedOrganization.Slug)

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

	newOrganization2, err := client.Organizations().Create(clerk.CreateOrganizationParams{
		Name:      "my-org-2",
		CreatedBy: users[0].ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, newOrganization2.ID)
	assert.Equal(t, "my-org-2", newOrganization2.Name)

	// Should short  Organizations by name DESC
	orderBy := "-name"
	organizations, err = client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		OrderBy: &orderBy,
	})
	assert.NoError(t, err)
	assert.Equal(t, organizations.Data[0].ID, newOrganization2.ID)

	orderBy = "+name"
	// Should short  Organizations by name ASC
	organizations, err = client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		OrderBy: &orderBy,
	})
	assert.NoError(t, err)
	assert.Equal(t, organizations.Data[0].ID, newOrganization.ID)

	// Should return 1 organization - exact match on ID
	organizations, err = client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		Query: newOrganization.ID,
	})

	if err != nil {
		t.Fatalf("Organizations.ListAll returned error: %v", err)
	}
	if organizations == nil {
		t.Fatalf("Organizations.ListAll returned nil")
	}

	assert.Equal(t, len(organizations.Data), 1)
	assert.Equal(t, organizations.Data[0].ID, newOrganization.ID)
	assert.Greater(t, organizations.TotalCount, int64(0))

	// Should return organizations with name containing "my-org"
	organizations, err = client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		Query: "my-org",
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
		assert.Contains(t, organization.Name, "my-org")
	}

	// Should return 2 results when using the userID as search param with ('+') prefix.
	filteredOrganizationByUserID, err := client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
		UserIDs: []string{"+" + users[0].ID},
	})
	assert.NoError(t, err)

	assert.Equal(t, len(filteredOrganizationByUserID.Data), 2)
	assert.Equal(t, filteredOrganizationByUserID.TotalCount, int64(2))

	// Should return non empty list
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

	// Should return 1 result when providing matching query
	organizationMembershipsWithQuery, err := client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		Query:          &organizationMemberships.Data[0].PublicUserData.UserID,
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 1)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(1))

	// Should return 0 results when providing non matching query
	query := "somequery"
	organizationMembershipsWithQuery, err = client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		Query:          &query,
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 0)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(0))

	// Should return 1 results when using the email as search param and email exist
	organizationMembershipsWithQuery, err = client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		EmailAddresses: []string{organizationMemberships.Data[0].PublicUserData.Identifier},
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 1)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(1))

	// Should return 0 results when using the email as search param and email does not exist
	organizationMembershipsWithQuery, err = client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		EmailAddresses: []string{"justanemptyemai@clerk.com"},
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 0)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(0))

	// Should change the role of a user to admin
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

	// Should return 2 results when using Role admin as a query param
	organizationMembershipsWithQuery, err = client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		Roles:          []string{"admin"},
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 2)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(2))

	// Should change the role of a user to basic_member
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

	// Should return 1 results when using Role admin as a query param
	organizationMembershipsWithQuery, err = client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		Roles:          []string{"admin"},
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if organizationMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(organizationMembershipsWithQuery.Data), 1)
	assert.Equal(t, organizationMembershipsWithQuery.TotalCount, int64(1))

	// Should delete member
	deletedOrganizationMembership, err := client.Organizations().DeleteMembership(newOrganization.ID,
		organizationMemberships.Data[0].PublicUserData.UserID,
	)
	if err != nil {
		t.Fatalf("Organizations.DeleteMembership returned error: %v", err)
	}
	if deletedOrganizationMembership == nil {
		t.Fatalf("Organizations.DeleteMembership returned nil")
	}
	assert.Equal(t, deletedOrganizationMembership.Role, "basic_member")

	// Should return 0 result
	afterDeletionMembershipsWithQuery, err := client.Organizations().ListMemberships(clerk.ListOrganizationMembershipsParams{
		OrganizationID: newOrganization.ID,
		Query:          &organizationMemberships.Data[0].PublicUserData.UserID,
	})
	if err != nil {
		t.Fatalf("Organizations.ListMemberships with query returned error: %v", err)
	}
	if afterDeletionMembershipsWithQuery == nil {
		t.Fatalf("Organizations.ListMemberships with query returned nil")
	}
	assert.Equal(t, len(afterDeletionMembershipsWithQuery.Data), 0)
	assert.Equal(t, afterDeletionMembershipsWithQuery.TotalCount, int64(0))

	deleteResponse, err := client.Organizations().Delete(newOrganization.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newOrganization.ID, deleteResponse.ID)
}
