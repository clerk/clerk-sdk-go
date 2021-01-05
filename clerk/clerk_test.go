package clerk

import "testing"

func TestNewClientBaseUrl(t *testing.T) {
	c, err := NewClient("token")
	if err != nil {
		t.Errorf("NewClient failed")
	}

	if got, want := c.BaseURL.String(), clerkBaseUrl; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestNewClientCreatesDifferenceClients(t *testing.T) {
	token := "token"
	c, _ := NewClient(token)
	c2, _ := NewClient(token)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}
