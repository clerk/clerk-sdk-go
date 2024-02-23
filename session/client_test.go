package session

import (
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

func TestSessionClientGet(t *testing.T) {
	t.Parallel()
	id := "sess_123"
	status := "active"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"%s"}`, id, status)),
			Method: http.MethodGet,
			Path:   "/v1/sessions/" + id,
		},
	}
	client := NewClient(config)
	session, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, session.ID)
	require.Equal(t, status, session.Status)
}

func TestSessionClientList(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"id":"sess_123","status":"active"}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/sessions",
			Query: &url.Values{
				"paginated": []string{"true"},
				"limit":     []string{"1"},
				"offset":    []string{"2"},
				"status":    []string{"active"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		Status: clerk.String("active"),
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Sessions))
	require.Equal(t, "sess_123", list.Sessions[0].ID)
	require.Equal(t, "active", list.Sessions[0].Status)
}

func TestSessionClientRevoke(t *testing.T) {
	t.Parallel()
	id := "sess_123"
	status := "revoked"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"%s"}`, id, status)),
			Method: http.MethodPost,
			Path:   "/v1/sessions/" + id + "/revoke",
		},
	}
	client := NewClient(config)
	session, err := client.Revoke(context.Background(), &RevokeParams{
		ID: id,
	})
	require.NoError(t, err)
	require.Equal(t, id, session.ID)
	require.Equal(t, status, session.Status)
}

func TestSessionClientVerify(t *testing.T) {
	t.Parallel()
	id := "sess_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodPost,
			Path:   "/v1/sessions/" + id + "/verify",
		},
	}
	client := NewClient(config)
	session, err := client.Verify(context.Background(), &VerifyParams{
		ID:    id,
		Token: clerk.String("the-token"),
	})
	require.NoError(t, err)
	require.Equal(t, id, session.ID)
}
