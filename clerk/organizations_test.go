package clerk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
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

func TestOrganizationsService_ListAll_happyPathWithQuery(t *testing.T) {
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
			"query": {"test"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want *OrganizationsResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)
	got, _ := client.Organizations().ListAll(ListAllOrganizationsParams{
		Query: "test",
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
func TestOrganizationsService_ListAllInvitations_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
			"data": %s,
			"total_count": 2
	}`, dummyOrganizationInvitationListJson)

	organizationID := "org_1mebQggrD3xO5JfuHk7clQ94ysA"

	mux.HandleFunc(fmt.Sprintf("/organizations/%s/invitations", organizationID), func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Organizations().ListInvitations(ListAllOrganizationInvitationParams{
		OrganizationID: organizationID,
	})
	if err != nil {
		t.Fatalf("Organizations.ListInvitations returned error: %v", err)
	}

	var want OrganizationInvitationResponse
	err = json.Unmarshal([]byte(expectedResponse), &want)
	if err != nil {
		t.Fatalf("Error unmarshaling expected response: %v", err)
	}

	if len(got.Data) != len(want.Data) {
		t.Errorf("Organizations.ListInvitations returned %d invitations, want %d", len(got.Data), len(want.Data))
	}

	if got.TotalCount != want.TotalCount {
		t.Errorf("Organizations.ListInvitations returned total_count %d, want %d", got.TotalCount, want.TotalCount)
	}
}

func TestOrganizationsService_DeleteInvitation(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	organizationID := "org_123"
	invitationID := "inv_456"
	requestingUserID := "user_789"

	expectedResponse := dummyOrganizationRevokedInvitationJson

	mux.HandleFunc(fmt.Sprintf("/organizations/%s/invitations/%s", organizationID, invitationID), func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, "DELETE")
		testHeader(t, r, "Authorization", "Bearer token")

		if got := r.URL.Query().Get("requesting_user_id"); got != requestingUserID {
			t.Errorf("Request URL query 'requesting_user_id' = %v, want %v", got, requestingUserID)
		}

		fmt.Fprint(w, expectedResponse)
	})

	params := DeleteOrganizationInvitationParams{
		OrganizationID:   organizationID,
		InvitationID:     invitationID,
		RequestingUserID: requestingUserID,
	}

	invitation, err := client.Organizations().DeleteInvitation(params)
	if err != nil {
		t.Errorf("Organizations.DeleteInvitation returned error: %v", err)
	}

	want := &OrganizationInvitation{}
	err = json.Unmarshal([]byte(expectedResponse), want)
	if err != nil {
		t.Fatalf("Error unmarshaling expected response: %v", err)
	}

	if !reflect.DeepEqual(invitation, want) {
		t.Errorf("Organizations.DeleteInvitation returned %+v, want %+v", invitation, want)
	}
}
func TestOrganizationsService_UpdateLogo(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	organizationID := "org_123"
	uploaderUserID := "user_123"
	expectedResponse := fmt.Sprintf(`{"id":"%s"}`, organizationID)
	filename := "200x200-grayscale.jpg"
	file, err := os.Open(path.Join("..", "testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	mux.HandleFunc(
		fmt.Sprintf("/organizations/%s/logo", organizationID),
		func(w http.ResponseWriter, req *http.Request) {
			testHttpMethod(t, req, http.MethodPut)
			testHeader(t, req, "Authorization", "Bearer token")
			// Assert that the request is sent as multipart/form-data
			if !strings.Contains(req.Header["Content-Type"][0], "multipart/form-data") {
				t.Errorf("expected content-type to be multipart/form-data, got %s", req.Header["Content-Type"])
			}
			defer req.Body.Close()

			// Check that the file is sent correctly
			fileParam, header, err := req.FormFile("file")
			if err != nil {
				t.Fatal(err)
			}
			if header.Filename != filename {
				t.Errorf("expected %s, got %s", filename, header.Filename)
			}
			defer fileParam.Close()

			got := make([]byte, header.Size)
			gotSize, err := fileParam.Read(got)
			if err != nil {
				t.Fatal(err)
			}
			fileInfo, err := file.Stat()
			if err != nil {
				t.Fatal(err)
			}
			want := make([]byte, fileInfo.Size())
			_, err = file.Seek(0, 0)
			if err != nil {
				t.Fatal(err)
			}
			wantSize, err := file.Read(want)
			if err != nil {
				t.Fatal(err)
			}
			if gotSize != wantSize {
				t.Errorf("read different size of files")
			}
			if !bytes.Equal(got, want) {
				t.Errorf("file was not sent correctly")
			}

			// Check the uploader user ID
			if got, ok := req.MultipartForm.Value["uploader_user_id"]; !ok || got[0] != uploaderUserID {
				t.Errorf("expected %s, got %s", uploaderUserID, got)
			}

			fmt.Fprint(w, expectedResponse)
		},
	)

	// Trigger a request to update the logo with the file
	org, err := client.Organizations().UpdateLogo(organizationID, UpdateOrganizationLogoParams{
		File:           file,
		Filename:       &filename,
		UploaderUserID: "user_123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if org.ID != organizationID {
		t.Errorf("expected %s, got %s", organizationID, org.ID)
	}
}

func TestOrganizationsService_DeleteLogo(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	organizationID := "org_123"
	mux.HandleFunc(
		fmt.Sprintf("/organizations/%s/logo", organizationID),
		func(w http.ResponseWriter, req *http.Request) {
			testHttpMethod(t, req, http.MethodDelete)
			testHeader(t, req, "Authorization", "Bearer token")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s"}`, organizationID))
		},
	)

	// Trigger a request to delete the logo
	_, err := client.Organizations().DeleteLogo(organizationID)
	if err != nil {
		t.Fatal(err)
	}
}

const dummyOrganizationJson = `{
	"object": "organization",
	"id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
	"name": "test-org",
	"slug": "org_slug",
	"members_count": 42,
	"created_by": "user_1mebQggrD3xO5JfuHk7clQ94ysA",
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
	"created_by": "user_1mebQggrD3xO5JfuHk7clQ94ysA",
	"created_at": 1610783813,
	"updated_at": 1610783813,
	"public_metadata": {},
	"private_metadata": {
		"app_id": 8,
	}
}`

const dummyOrganizationInvitationListJson = `[
	{
	"object": "organization_invitation",
	"id": "orginv_1mebQggrD3xO5JfuHk7clQ94ysA",
	"email_address": "test1@test.com",
	"role": "admin",
	"organization_id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
	"status": "pending",
	"public_metadata": { },
	"private_metadata": { },
	"created_at": 1721806882553,
	"updated_at": 1721806882553
	},
	{
	"object": "organization_invitation",
	"id": "orginv_1mebQggrD3xO5JfuHk7clQ94ysA",
	"email_address": "test2@test.com",
	"role": "admin",
	"organization_id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
	"status": "pending",
	"public_metadata": { },
	"private_metadata": { },
	"created_at": 1721806882553,
	"updated_at": 1721806882553
	}
]`

const dummyOrganizationRevokedInvitationJson = `{
			"object": "organization_invitation",
			"id": "orginv_1mebQggrD3xO5JfuHk7clQ94ysA",
			"email_address": "test@example.com",
			"role": "admin",
			"organization_id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
			"status": "revoked",
			"public_metadata": {},
			"private_metadata": {},
			"created_at": 1621234567890,
			"updated_at": 1621234567890
}`
