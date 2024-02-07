package client

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

func TestClientClientGet(t *testing.T) {
	t.Parallel()
	id := "client_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodGet,
			Path:   "/v1/clients/" + id,
		},
	}
	c := NewClient(config)
	client, err := c.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, client.ID)
}

func TestClientClientVerify(t *testing.T) {
	t.Parallel()
	id := "client_123"
	token := "the-token"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"token":"%s"}`, token)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodPost,
			Path:   "/v1/clients/verify",
		},
	}
	c := NewClient(config)
	client, err := c.Verify(context.Background(), &VerifyParams{
		Token: clerk.String("the-token"),
	})
	require.NoError(t, err)
	require.Equal(t, id, client.ID)
}

func TestClientClientList(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
	"data": [{"id":"client_123","last_active_session_id":"sess_123"}],
	"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/clients",
			Query: &url.Values{
				"limit":  []string{"1"},
				"offset": []string{"2"},
			},
		},
	}
	c := NewClient(config)
	params := &ListParams{}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := c.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Clients))
	require.Equal(t, "client_123", list.Clients[0].ID)
	require.Equal(t, "sess_123", *list.Clients[0].LastActiveSessionID)
}
