package clerk

import (
	"net/http"

	"gopkg.in/square/go-jose.v2"
)

type JWKSService service

type JWKS jose.JSONWebKeySet

func (s *JWKSService) ListAll() (*JWKS, error) {
	req, _ := s.client.NewRequest(http.MethodGet, "jwks", nil)

	jwks := JWKS{}
	_, err := s.client.Do(req, &jwks)
	if err != nil {
		return nil, err
	}

	return &jwks, nil
}
