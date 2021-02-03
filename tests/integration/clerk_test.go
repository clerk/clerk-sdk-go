// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"os"
)

type key string

const (
	APIKey       = "CLERK_API_KEY"
	SessionToken = "CLERK_SESSION_TOKEN"
	SessionID    = "CLERK_SESSION_ID"
)

func createClient() clerk.Client {
	apiKey := getEnv(APIKey)
	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic("Unable to create Clerk client")
	}
	return client
}

func getEnv(k string) string {
	envValue := os.Getenv(k)
	if envValue == "" {
		panic("Missing env variable " + k)
	}
	return envValue
}
