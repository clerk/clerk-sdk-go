package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func returnActiveSession(w http.ResponseWriter, req *http.Request) {
	sessionClaims, ok := clerk.SessionFromContext(req.Context())
	if ok {
		jsonResp, _ := json.Marshal(sessionClaims)
		fmt.Fprintf(w, string(jsonResp))
	} else {
		// handle non-authenticated user
	}

}

func main() {
	fmt.Print("Clerk secret key: ")
	var apiKey string
	fmt.Scanf("%s", &apiKey)

	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	injectActiveSession := clerk.WithSessionV2(client)
	mux.Handle("/session", injectActiveSession(http.HandlerFunc(returnActiveSession)))

	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
