//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyChecks(t *testing.T) {
	client := createClient()
	// Create a domain first
	domainName := gofakeit.DomainName()
	domain, err := client.Domains().Create(clerk.CreateDomainParams{
		Name:        domainName,
		IsSatellite: true,
	})
	require.NoError(t, err)

	// Now trigger a proxy check. Most likely a proxy is not configured.
	_, err = client.ProxyChecks().Create(clerk.CreateProxyCheckParams{
		DomainID: domain.ID,
		ProxyURL: "https://" + domainName + "/__clerk",
	})
	assert.Error(t, err)
}
