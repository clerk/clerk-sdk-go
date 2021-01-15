package main

import (
	"fmt"
	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func main() {
	fmt.Print("Clerk API Key: ")
	var apiKey string
	fmt.Scanf("%s", &apiKey)

	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	retrieveUsers(client)
	retrieveSessions(client)
}

func retrieveUsers(client clerk.Client) {
	users, err := client.Users().ListAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Users:")
	for i, user := range users {
		fmt.Printf("%v. %v %v\n", i+1, *user.FirstName, *user.LastName)
	}
}

func retrieveSessions(client clerk.Client) {
	sessions, err := client.Sessions().ListAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nSessions:")
	for i, session := range sessions {
		fmt.Printf("%v. %v (%v)\n", i+1, session.ID, session.Status)
	}
}
