package proxycheck

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

func TestProxyCheckClientCreate(t *testing.T) {
	t.Parallel()
	id := "proxchk_123"
	proxyURL := "https://clerk.com/__proxy"
	domainID := "dmn_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"domain_id":"%s","proxy_url":"%s"}`, domainID, proxyURL)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","domain_id":"%s","proxy_url":"%s","successful":true}`, id, domainID, proxyURL)),
			Method: http.MethodPost,
			Path:   "/v1/proxy_checks",
		},
	}
	client := NewClient(config)
	proxyCheck, err := client.Create(context.Background(), &CreateParams{
		ProxyURL: clerk.String(proxyURL),
		DomainID: clerk.String(domainID),
	})
	require.NoError(t, err)
	require.Equal(t, id, proxyCheck.ID)
	require.Equal(t, proxyURL, proxyCheck.ProxyURL)
	require.Equal(t, domainID, proxyCheck.DomainID)
	require.True(t, proxyCheck.Successful)
}
