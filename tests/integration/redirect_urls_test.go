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

func TestRedirectURLs(t *testing.T) {
	client := createClient()

	redirectURLs, err := client.RedirectURLs().ListAll()
	assert.Nil(t, err)

	previousRedirectURLsCount := len(redirectURLs)

	url := fmt.Sprintf("http://www.%d.com", time.Now().Unix())
	redirectURL, err := client.RedirectURLs().Create(clerk.CreateRedirectURLParams{
		URL: url,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, redirectURL.ID)
	assert.Equal(t, url, redirectURL.URL)
	assert.Equal(t, "redirect_url", redirectURL.Object)

	redirectURLs, err = client.RedirectURLs().ListAll()
	assert.Nil(t, err)
	assert.Equal(t, previousRedirectURLsCount+1, len(redirectURLs))

	deletedResponse, err := client.RedirectURLs().Delete(redirectURL.ID)
	assert.Nil(t, err)
	assert.Equal(t, redirectURL.ID, deletedResponse.ID)
	assert.True(t, deletedResponse.Deleted)
}
