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

func TestBlocklists(t *testing.T) {
	client := createClient()

	blocklistIdentifiers, err := client.Blocklists().ListAllIdentifiers()
	assert.Nil(t, err)

	previousCount := blocklistIdentifiers.TotalCount

	identifier := fmt.Sprintf("email_%d@example.com", time.Now().Unix())
	blocklistIdentifier, err := client.Blocklists().CreateIdentifier(clerk.CreateBlocklistIdentifierParams{
		Identifier: identifier,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, blocklistIdentifier.ID)
	assert.Equal(t, identifier, blocklistIdentifier.Identifier)
	assert.Equal(t, "blocklist_identifier", blocklistIdentifier.Object)

	blocklistIdentifiers, err = client.Blocklists().ListAllIdentifiers()
	assert.Nil(t, err)
	assert.Equal(t, previousCount+1, blocklistIdentifiers.TotalCount)

	deletedResponse, err := client.Blocklists().DeleteIdentifier(blocklistIdentifier.ID)
	assert.Nil(t, err)
	assert.Equal(t, blocklistIdentifier.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)
}
