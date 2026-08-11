package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestDomainCreate(t *testing.T) {
	name := "clerk.com"
	id := "dmn_123"
	dmarcHost := "_dmarc." + name
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:  t,
				In: json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"name":"%s",
					"cname_targets":[{"host":"clerk.%s","value":"frontend-api.clerk.services","required":false}],
					"dns_targets":[
						{"host":"clerk.%s","value":"frontend-api.clerk.services","required":false,"record_type":"CNAME"},
						{"host":"%s","value":"v=DMARC1; p=none;","required":true,"record_type":"TXT","automation_disposition":"safe_to_create"}
					]
				}`, id, name, name, name, dmarcHost)),
				Path:   "/v1/domains",
				Method: http.MethodPost,
			},
		},
	}))

	dmn, err := Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.Equal(t, name, dmn.Name)
	require.Equal(t, []clerk.CNAMETarget{{
		Host:  "clerk." + name,
		Value: "frontend-api.clerk.services",
	}}, dmn.CNAMETargets)
	require.Equal(t, []clerk.DNSTarget{
		{
			Host:       "clerk." + name,
			Value:      "frontend-api.clerk.services",
			Required:   false,
			RecordType: "CNAME",
		},
		{
			Host:                  dmarcHost,
			Value:                 "v=DMARC1; p=none;",
			Required:              true,
			RecordType:            "TXT",
			AutomationDisposition: "safe_to_create",
		},
	}, dmn.DNSTargets)
}

func TestDomainCreate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
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
		},
	}))

	_, err := Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestDomainUpdate(t *testing.T) {
	id := "dmn_456"
	name := "clerk.dev"
	dmarcHost := "_dmarc." + name
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:  t,
				In: json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"name":"%s",
					"dns_targets":[{"host":"%s","value":"v=DMARC1; p=none;","required":true,"record_type":"TXT"}]
				}`, id, name, dmarcHost)),
				Path:   fmt.Sprintf("/v1/domains/%s", id),
				Method: http.MethodPatch,
			},
		},
	}))

	dmn, err := Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.Equal(t, name, dmn.Name)
	require.Equal(t, []clerk.DNSTarget{{
		Host:       dmarcHost,
		Value:      "v=DMARC1; p=none;",
		Required:   true,
		RecordType: "TXT",
	}}, dmn.DNSTargets)
}

func TestDomainUpdate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
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
		},
	}))

	_, err := Update(context.Background(), "dmn_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestDomainDelete(t *testing.T) {
	id := "dmn_789"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
				Path:   fmt.Sprintf("/v1/domains/%s", id),
				Method: http.MethodDelete,
			},
		},
	}))

	dmn, err := Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.True(t, dmn.Deleted)
}

func TestDomainList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [{"id":"dmn_123","name":"clerk.com"}],
	"total_count": 1
}`),
				Path:   "/v1/domains",
				Method: http.MethodGet,
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Domains))
	require.Equal(t, "dmn_123", list.Domains[0].ID)
	require.Equal(t, "clerk.com", list.Domains[0].Name)
}
