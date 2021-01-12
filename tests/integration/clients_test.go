// +build integration

package integration

import (
	"testing"
)

func TestClients(t *testing.T) {
	client := createClient()

	clients, err := client.Clients().ListAll()
	if err != nil {
		t.Fatalf("Clients.ListAll returned error: %v", err)
	}
	if clients == nil {
		t.Fatalf("Clients.ListAll returned nil")
	}

	for _, response := range clients {
		if response.LastActiveSessionID == nil {
			continue
		}
		clientId := response.ID
		clientResponse, err := client.Clients().Read(clientId)
		if err != nil {
			t.Fatalf("Clients.Read returned error: %v", err)
		}
		if clientResponse == nil {
			t.Fatalf("Clients.Read returned nil")
		}
	}
}
