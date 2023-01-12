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

func TestOrganizationsService_Read(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyOrganizationJson
	orgID := "randomIDorSlug"

	mux.HandleFunc(fmt.Sprintf("/organizations/%s", orgID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Organizations().Read(orgID)
	if err != nil {
		t.Fatal(err)
	}

	var want Organization
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationsService_Update(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()
	var payload UpdateOrganizationParams
	_ = json.Unmarshal([]byte(dummyUpdateOrganizationJson), &payload)

	expectedResponse := dummyOrganizationJson
	orgID := "randomIDorSlug"

	mux.HandleFunc(fmt.Sprintf("/organizations/%s", orgID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Organizations().Update(orgID, payload)
	if err != nil {
		t.Fatal(err)
	}

	var want Organization
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Response = %v, want %v", got, &want)
	}
}

func TestOrganizationsService_invalidServer(t *testing.T) {
	client, _ := NewClient("token")
	var payload UpdateOrganizationParams
	_ = json.Unmarshal([]byte(dummyUpdateOrganizationJson), &payload)

	_, err := client.Organizations().Update("someOrgId", payload)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestOrganizationsService_ListAll_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrganizationJson)

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want *OrganizationsResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Organizations().ListAll(ListAllOrganizationsParams{})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organizations to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_ListAll_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrganizationJson)

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"limit":                 {"5"},
			"offset":                {"6"},
			"include_members_count": {"true"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want *OrganizationsResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	limit := 5
	offset := 6
	got, _ := client.Organizations().ListAll(ListAllOrganizationsParams{
		Limit:               &limit,
		Offset:              &offset,
		IncludeMembersCount: true,
	})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organizations to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	organizations, err := client.Organizations().ListAll(ListAllOrganizationsParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if organizations != nil {
		t.Errorf("Was not expecting any organizations to be returned, instead got %v", organizations)
	}
}

const dummyOrganizationJson = `{
        "object": "organization",
        "id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
        "name": "test-org",
        "slug": "org_slug",
		"members_count": 42,
        "created_at": 1610783813,
        "updated_at": 1610783813,
		"public_metadata": {
			"address": {
				"street": "Pennsylvania Avenue",
				"number": "1600"
			}
		},
		"private_metadata": {
			"app_id": 5
		}
    }`

const dummyUpdateOrganizationJson = `{
        "object": "organization",
        "id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
        "name": "test-org",
        "slug": "org_slug",
		"members_count": 42,
        "created_at": 1610783813,
        "updated_at": 1610783813,
		"public_metadata": {},
		"private_metadata": {
			"app_id": 8,
		}
    }`
