//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestAllowlists(t *testing.T) {
	client := createClient()

	allowlistIdentifiers, err := client.Allowlists().ListAllIdentifiers()
	assert.Nil(t, err)

	previousCount := allowlistIdentifiers.TotalCount

	identifier := fmt.Sprintf("email_%d@example.com", time.Now().Unix())
	allowlistIdentifier, err := client.Allowlists().CreateIdentifier(clerk.CreateAllowlistIdentifierParams{
		Identifier: identifier,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, allowlistIdentifier.ID)
	assert.Equal(t, identifier, allowlistIdentifier.Identifier)
	assert.Equal(t, "allowlist_identifier", allowlistIdentifier.Object)

	allowlistIdentifiers, err = client.Allowlists().ListAllIdentifiers()
	assert.Nil(t, err)
	assert.Equal(t, previousCount+1, allowlistIdentifiers.TotalCount)

	deletedResponse, err := client.Allowlists().DeleteIdentifier(allowlistIdentifier.ID)
	assert.Nil(t, err)
	assert.Equal(t, allowlistIdentifier.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)
}
