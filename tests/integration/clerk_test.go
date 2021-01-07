// +build integration

package integration

import (
	"github.com/clerkinc/clerk_server_sdk_go/v1/clerk"
	"os"
)

func createClient() clerk.Client {
	apiKey := os.Getenv("CLERK_API_KEY")
	if apiKey == "" {
		panic("Missing env variable CLERK_API_KEY")
	}

	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic("Unable to create Clerk client")
	}
	return client
}
