package emailaddress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestEmailAddressClientCreate(t *testing.T) {
	t.Parallel()
	email := "foo@bar.com"
	userID := "user_123"
	id := "idn_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s","user_id":"%s","verified":false}`, email, userID)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s"}`, id, email)),
			Method: http.MethodPost,
			Path:   "/v1/email_addresses",
		},
	}
	client := NewClient(config)
	emailAddress, err := client.Create(context.Background(), &CreateParams{
		UserID:       clerk.String(userID),
		EmailAddress: clerk.String(email),
		Verified:     clerk.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, id, emailAddress.ID)
	require.Equal(t, email, emailAddress.EmailAddress)
}

func TestEmailAddressClientCreate_Error(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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

func TestEmailAddressClientUpdate(t *testing.T) {
	t.Parallel()
	id := "idn_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"verified":true}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","verification":{"status":"verified"}}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/email_addresses/" + id,
		},
	}
	client := NewClient(config)
	emailAddress, err := client.Update(context.Background(), "idn_123", &UpdateParams{
		Verified: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, id, emailAddress.ID)
	require.Equal(t, "verified", emailAddress.Verification.Status)
}

func TestEmailAddressClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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
	_, err := client.Update(context.Background(), "idn_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestEmailAddressClientGet(t *testing.T) {
	t.Parallel()
	id := "idn_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","verification":{"status":"verified"}}`, id)),
			Method: http.MethodGet,
			Path:   "/v1/email_addresses/" + id,
		},
	}
	client := NewClient(config)
	emailAddress, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, emailAddress.ID)
	require.Equal(t, "verified", emailAddress.Verification.Status)
}

func TestEmailAddressClientDelete(t *testing.T) {
	t.Parallel()
	id := "idn_456"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/email_addresses/" + id,
		},
	}
	client := NewClient(config)
	emailAddress, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, emailAddress.ID)
	require.True(t, emailAddress.Deleted)
}
