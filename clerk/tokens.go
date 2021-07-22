package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type TokensService service

//type Session struct {
//	Object       string `json:"object"`
//	ID           string `json:"id"`
//	ClientID     string `json:"client_id"`
//	UserID       string `json:"user_id"`
//	Status       string `json:"status"`
//	LastActiveAt int64  `json:"last_active_at"`
//	ExpireAt     int64  `json:"expire_at"`
//	AbandonAt    int64  `json:"abandon_at"`
//}

// token is the short-lived JWT
func (s *TokensService) Verify(token string) error {
	// TODO: refresh periodically in the background
	jwks, err := s.fetchJWKs()
	if err != nil {
		return err
	}

	token, err := jwt.Parse(token, getKey)
	if err != nil {
		panic(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range claims {
		fmt.Printf("%s\t%v\n", key, value)
	}
	// verify

	//sessionsUrl := "sessions"
	//req, _ := s.client.NewRequest("GET", sessionsUrl)

	//var sessions []Session
	//_, err := s.client.Do(req, &sessions)
	//if err != nil {
	//	return nil, err
	//}
	//return sessions, nil
	return nil
}

func (s *TokensService) fetchJWKs() ([]*jwks, error) {
	// TODO: this is the dashboard's instance jwks
	req, err := http.NewRequest("GET", "https://clerk.prod.lclclerk.com/v1/.well-known/jwks.json", nil)
	if err != nil {
		panic(err)
		return err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		panic(err)
		return err
	}
	defer resp.Body.Close()

	jwks := make([]*jwk, 0)

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		panic(err)
		return err
	}

	if len(jwks) == 0 {
		panic("invalid length")
	}

	return jwks, nil
}
