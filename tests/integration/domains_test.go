//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestDomains(t *testing.T) {
	client := createClient()

	domains, err := client.Domains().ListAll()
	assert.Nil(t, err)

	domainCount := domains.TotalCount

	// create (satellite)

	name := gofakeit.DomainName()
	createDomainParams := clerk.CreateDomainParams{
		Name:        name,
		IsSatellite: true,
	}

	domain, err := client.Domains().Create(createDomainParams)
	assert.Nil(t, err)

	assert.Equal(t, "domain", domain.Object)
	assert.Equal(t, name, domain.Name)

	// list

	domains, err = client.Domains().ListAll()
	assert.Nil(t, err)
	assert.Equal(t, domainCount+1, domains.TotalCount)

	// update

	name = gofakeit.DomainName()
	updateDomainParams := clerk.UpdateDomainParams{
		Name: &name,
	}

	domain, err = client.Domains().Update(domain.ID, updateDomainParams)
	assert.Nil(t, err)

	assert.Equal(t, "domain", domain.Object)
	assert.Equal(t, name, domain.Name)

	// delete

	deletedResponse, err := client.Domains().Delete(domain.ID)
	assert.Nil(t, err)

	assert.Equal(t, domain.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)

	// list

	domains, err = client.Domains().ListAll()
	assert.Nil(t, err)
	assert.Equal(t, domainCount, domains.TotalCount)
}
