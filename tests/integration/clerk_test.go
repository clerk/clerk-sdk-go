//go:build integration
// +build integration

package integration

import (
	"net/http"
	"os"
	"time"

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

	httpClient := &http.Client{Timeout: time.Second * 20}
	client, err := clerk.NewClient(apiKey, clerk.WithBaseURL(apiUrl), clerk.WithHTTPClient(httpClient))
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
