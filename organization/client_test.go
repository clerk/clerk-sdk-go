package organization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestOrganizationClientCreate(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPost,
			Path:   "/v1/organizations",
		},
	}
	client := NewClient(config)
	organization, err := client.Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
}

func TestOrganizationClientCreate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"create-error-code"
	}],
	"clerk_trace_id":"create-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationClientGet(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	membersCount := int64(1)
	hasMemberWithElevatedPermissions := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(
				`{"id":"%s","name":"%s","members_count":%d,"has_member_with_elevated_permissions":%t}`,
				id,
				name,
				membersCount,
				hasMemberWithElevatedPermissions,
			)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + id,
			Query: &url.Values{
				"include_members_count":                        []string{"true"},
				"include_has_member_with_elevated_permissions": []string{"true"},
			},
		},
	}
	params := &GetParams{
		IncludeMembersCount:                     clerk.Bool(true),
		IncludeHasMemberWithElevatedPermissions: clerk.Bool(true),
	}
	client := NewClient(config)
	organization, err := client.Get(context.Background(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
	require.Equal(t, membersCount, *organization.MembersCount)
	require.NotNil(t, organization.HasMemberWithElevatedPermissions)
	require.Equal(t, hasMemberWithElevatedPermissions, *organization.HasMemberWithElevatedPermissions)

}

func TestOrganizationClientUpdate(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + id,
		},
	}
	client := NewClient(config)
	organization, err := client.Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
}

func TestOrganizationClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"update-error-code"
	}],
	"clerk_trace_id":"update-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Update(context.Background(), "org_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationClientDelete(t *testing.T) {
	t.Parallel()
	id := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + id,
		},
	}
	client := NewClient(config)
	organization, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.True(t, organization.Deleted)
}

func TestOrganizationClientList(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"id":"org_123","name":"Acme Inc","members_count":1,"has_member_with_elevated_permissions":true}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/organizations",
			Query: &url.Values{
				"limit":                 []string{"1"},
				"offset":                []string{"2"},
				"order_by":              []string{"-created_at"},
				"query":                 []string{"Acme"},
				"user_id":               []string{"user_123", "user_456"},
				"include_members_count": []string{"true"},
				"include_has_member_with_elevated_permissions": []string{"true"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		OrderBy:                                 clerk.String("-created_at"),
		Query:                                   clerk.String("Acme"),
		UserIDs:                                 []string{"user_123", "user_456"},
		IncludeMembersCount:                     clerk.Bool(true),
		IncludeHasMemberWithElevatedPermissions: clerk.Bool(true),
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Organizations))
	require.Equal(t, "org_123", list.Organizations[0].ID)
	require.Equal(t, "Acme Inc", list.Organizations[0].Name)
	require.NotNil(t, list.Organizations[0].MembersCount)
	require.Equal(t, int64(1), *list.Organizations[0].MembersCount)
	require.NotNil(t, *list.Organizations[0].HasMemberWithElevatedPermissions)
	require.Equal(t, true, *list.Organizations[0].HasMemberWithElevatedPermissions)
}

type testFile struct {
	bytes.Reader
}

func (_ *testFile) Close() error {
	return nil
}

func TestOrganizationClientUpdateLogo(t *testing.T) {
	t.Parallel()
	id := "org_123"
	userID := "user_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodPut,
			Path:   "/v1/organizations/" + id + "/logo",
		},
	}
	client := NewClient(config)
	organization, err := client.UpdateLogo(context.Background(), id, &UpdateLogoParams{
		UploaderUserID: &userID,
		File:           &testFile{Reader: *bytes.NewReader([]byte{})},
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
}

func TestOrganizationClientDeleteLogo(t *testing.T) {
	t.Parallel()
	id := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + id + "/logo",
		},
	}
	client := NewClient(config)
	organization, err := client.DeleteLogo(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
}

func TestOrganizationClientUpdateMetadata(t *testing.T) {
	t.Parallel()
	id := "org_123"
	metadata := `{"foo":"bar"}`
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"private_metadata":%s}`, metadata)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","private_metadata":%s}`, id, metadata)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + id + "/metadata",
		},
	}
	client := NewClient(config)
	metadataParam := json.RawMessage(metadata)
	organization, err := client.UpdateMetadata(context.Background(), id, &UpdateMetadataParams{
		PrivateMetadata: &metadataParam,
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.JSONEq(t, metadata, string(organization.PrivateMetadata))
}
