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

func TestInstanceService_List_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrganizationPermissionsJson)

	mux.HandleFunc("/organization_permissions", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"limit":  {"3"},
			"offset": {"2"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	want := &PermissionsResponse{}
	_ = json.Unmarshal([]byte(expectedResponse), want)

	got, _ := client.Instances().ListOrganizationPermissions(ListInstanceOrganizationPermissionsParams{
		Limit:  intToPtr(3),
		Offset: intToPtr(2),
	})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organization permissions to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestInstanceService_List_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	orgPermissions, err := client.Instances().ListOrganizationPermissions(ListInstanceOrganizationPermissionsParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if orgPermissions != nil {
		t.Errorf("Was not expecting any organization permissions to be returned, instead got %v", orgPermissions)
	}
}

const dummyOrganizationPermissionID = "perm_1mebQggrD3xO5JfuHk7clQ94ysA"

const dummyOrganizationPermissionsJson = `{
	"object": "permission",
	"id": "perm_1mebQggrD3xO5JfuHk7clQ94ysA",
	"name": "Manage organization",
	"key": "org:sys_profile:manage",
	"description": "Permission to manage an organization.",
	"type": "system",
	"created_at": 1610783813,
	"updated_at": 1610783813
}`
