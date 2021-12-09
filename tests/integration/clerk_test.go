// +build integration

package integration

import (
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type key string

const (
	APIUrl = "CLERK_API_URL"
	APIKey = "CLERK_API_KEY"
)

func createClient() clerk.Client {
	apiUrl := getEnv(APIUrl)
	apiKey := getEnv(APIKey)

	client, err := clerk.NewClient(apiKey, clerk.WithBaseURL(apiUrl))
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
