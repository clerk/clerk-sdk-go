package main

import (
	"fmt"
	"github.com/clerkinc/clerk_server_sdk_go/v1/clerk"
)

func main() {
	fmt.Print("Clerk API Key: ")
	var apiKey string
	fmt.Scanf("%s", &apiKey)

	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	users, err := client.Users.ListAll()
	if err != nil {
		panic(err)
	}

	for i, user := range users {
		userDetails, err := client.Users.Read(user.ID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v. %v %v\n", i+1, *userDetails.FirstName, *userDetails.LastName)
	}
}
