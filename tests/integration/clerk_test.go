// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"os"
)

type key string

const (
	APIUrl       = "CLERK_API_URL"
	APIKey       = "CLERK_API_KEY"
	SessionToken = "CLERK_SESSION_TOKEN"
	SessionID    = "CLERK_SESSION_ID"
)

func createClient() clerk.Client {
	apiUrl := getEnv(APIUrl)
	apiKey := getEnv(APIKey)
	return createClientWithKey(apiUrl, apiKey)
}

func createClientWithKey(apiUrl string, apiKey string) clerk.Client {
	client, err := clerk.NewClientWithBaseUrl(apiKey, apiUrl)
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
