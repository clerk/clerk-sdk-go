package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceService_CreateOrgRole(t *testing.T) {
	expectedResponse := dummyOrgRoleJson

	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/organization_roles", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, expectedResponse)
	})

	createParams := CreateInstanceOrganizationRoleParams{
		Name:        "custom role",
		Key:         "org:custom_role",
		Description: "my org custom role",
		Permissions: []string{},
	}

	got, err := client.Instances().CreateOrganizationRole(createParams)
	assert.NoError(t, err)

	var want Role
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationRolesService_Read(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyOrgRoleJson

	mux.HandleFunc(fmt.Sprintf("/organization_roles/%s", dummyOrgRoleID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Instances().ReadOrganizationRole(dummyOrgRoleID)
	if err != nil {
		t.Fatal(err)
	}

	var want Role
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationRolesService_Update(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()
	var payload UpdateInstanceOrganizationRoleParams
	_ = json.Unmarshal([]byte(dummyUpdateOrgRoleJson), &payload)

	expectedResponse := dummyOrgRoleJson
	mux.HandleFunc(fmt.Sprintf("/organization_roles/%s", dummyOrgRoleID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Instances().UpdateOrganizationRole(dummyOrgRoleID, payload)
	if err != nil {
		t.Fatal(err)
	}

	var want Role
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationRolesService_Update_invalidServer(t *testing.T) {
	client, _ := NewClient("token")
	var payload UpdateInstanceOrganizationRoleParams
	_ = json.Unmarshal([]byte(dummyUpdateOrgRoleJson), &payload)

	_, err := client.Instances().UpdateOrganizationRole("someOrgRoleId", payload)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestOrganizationsService_List_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrgRoleJson)

	mux.HandleFunc("/organization_roles", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		expectedQuery := url.Values{
			"limit":    {"5"},
			"offset":   {"6"},
			"query":    {"my-query"},
			"order_by": {"created_at"},
		}
		assert.Equal(t, expectedQuery, req.URL.Query())

		fmt.Fprint(w, expectedResponse)
	})

	var want *RolesResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Instances().ListOrganizationRole(ListInstanceOrganizationRoleParams{
		Limit:   intToPtr(5),
		Offset:  intToPtr(6),
		Query:   stringToPtr("my-query"),
		OrderBy: stringToPtr("created_at"),
	})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organization roles to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_List_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrgRoleJson)

	mux.HandleFunc("/organization_roles", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"limit":  {"5"},
			"offset": {"6"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want *RolesResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Instances().ListOrganizationRole(ListInstanceOrganizationRoleParams{
		Limit:  intToPtr(5),
		Offset: intToPtr(6),
	})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organization roles to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_List_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	orgRoles, err := client.Instances().ListOrganizationRole(ListInstanceOrganizationRoleParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if orgRoles != nil {
		t.Errorf("Was not expecting any organization roles to be returned, instead got %v", orgRoles)
	}
}

func TestOrganizationRolesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/organization_roles/%s", dummyOrgRoleID),
		func(w http.ResponseWriter, req *http.Request) {
			testHttpMethod(t, req, http.MethodDelete)
			testHeader(t, req, "Authorization", "Bearer token")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s"}`, dummyOrgRoleID))
		},
	)

	_, err := client.Instances().DeleteOrganizationRole(dummyOrgRoleID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrganizationRolesService_AssignPermission(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyOrgRoleJson
	mux.HandleFunc(fmt.Sprintf("/organization_roles/%s/permissions/%s", dummyOrgRoleID, dummyOrgPermissionID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Instances().AssignOrganizationRolePermission(dummyOrgRoleID, dummyOrgPermissionID)
	if err != nil {
		t.Fatal(err)
	}

	var want Role
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationRolesService_RemovePermission(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyOrgRoleJson
	mux.HandleFunc(fmt.Sprintf("/organization_roles/%s/permissions/%s", dummyOrgRoleID, dummyOrgPermissionID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})
	got, err := client.Instances().RemoveOrganizationRolePermission(dummyOrgRoleID, dummyOrgPermissionID)
	if err != nil {
		t.Fatal(err)
	}

	var want Role
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

const dummyOrgRoleID = "role_1mebQggrD3xO5JfuHk7clQ94ysA"
const dummyOrgPermissionID = "perm_1mebQggrD3xO5JfuHk7clQ94ysA"

const dummyOrgRoleJson = `{
	"object": "role",
	"id": "role_1mebQggrD3xO5JfuHk7clQ94ysA",
	"name": "custom role",
	"key": "org:custom_role",
	"description": "my org custom role",
	"permissions": [
		{
			"object": "permission",
			"id": "perm_2YDDmdVaAyCjmQq6RhrvNsyuUbC",
			"name": "Read members",
			"key": "org:sys_memberships:read",
			"description": "Permission to read the members of an organization.",
			"type": "system",
			"created_at": 1700051416222,
			"updated_at": 1700051416222
		},
		{
			"object": "permission",
			"id": "perm_2YDDmcGvIxBiLyPB8mUhLbWhWqG",
			"name": "Manage members",
			"key": "org:sys_memberships:manage",
			"description": "Permission to manage the members of an organization.",
			"type": "system",
			"created_at": 1700051416224,
			"updated_at": 1700051416224
		},
		{
			"object": "permission",
			"id": "perm_2YDDmbh2v71OSrvPK5WtUrh8FyX",
			"name": "Delete organization",
			"key": "org:sys_profile:delete",
			"description": "Permission to delete an organization.",
			"type": "system",
			"created_at": 1700051416221,
			"updated_at": 1700051416221
		},
		{
			"object": "permission",
			"id": "perm_2YDDmeXfQl3rdJr0FgkIclA1v5I",
			"name": "View reports",
			"key": "org:report:view",
			"description": "Permission to have the ability to view reports.",
			"type": "user",
			"created_at": 1700051416227,
			"updated_at": 1700051416227
		}
	],
	"is_creator_eligible": true,
	"created_at": 1610783813,
	"updated_at": 1610783813
}`

const dummyUpdateOrgRoleJson = `{
	"object": "role",
	"id": "role_1mebQggrD3xO5JfuHk7clQ94ysA",
	"name": "custom org 2",
	"key": "org:custom_role_2",
	"description": "my org custom role",
	"permissions": [
		{
			"object": "permission",
			"id": "perm_2YDDmdVaAyCjmQq6RhrvNsyuUbC",
			"name": "Read members",
			"key": "org:sys_memberships:read",
			"description": "Permission to read the members of an organization.",
			"type": "system",
			"created_at": 1700051416222,
			"updated_at": 1700051416222
		},
		{
			"object": "permission",
			"id": "perm_2YDDmeXfQl3rdJr0FgkIclA1v5I",
			"name": "View reports",
			"key": "org:report:view",
			"description": "Permission to have the ability to view reports.",
			"type": "user",
			"created_at": 1700051416227,
			"updated_at": 1700051416227
		}
	],
	"is_creator_eligible": false,
	"created_at": 1610783813,
	"updated_at": 1610783813
}`
