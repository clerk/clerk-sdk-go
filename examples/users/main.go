package main

import (
	"fmt"
	"github.com/clerkinc/clerk_server_sdk_go/v1/clerk"
)

func main() {
	fmt.Print("Clerk API Key: ")
	var apiKey string
	fmt.Scanf("%s", &apiKey)

	users := retrieveUsers(apiKey)
	for i, user := range users {
		fmt.Printf("%v. %v\n", i+1, user.ID)
	}
}

func retrieveUsers(apiKey string) []clerk.User {
	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	users, err := client.Users.ListAll()
	if err != nil {
		panic(err)
	}

	return users
}
