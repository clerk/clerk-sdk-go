package main

import (
	"encoding/json"
	"fmt"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"log"
	"net/http"
)

func returnActiveSession(w http.ResponseWriter, req *http.Request) {
	session := req.Context().Value(clerk.ActiveSession)
	jsonResp, _ := json.Marshal(session)

	fmt.Fprintf(w, string(jsonResp))
}

func main() {
	fmt.Print("Clerk API Key: ")
	var apiKey string
	fmt.Scanf("%s", &apiKey)

	client, err := clerk.NewClient(apiKey)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	injectActiveSession := clerk.WithSession(client)
	mux.Handle("/session", injectActiveSession(http.HandlerFunc(returnActiveSession)))

	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
