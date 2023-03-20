package clerk

import (
	"net/http"

	"github.com/go-jose/go-jose/v3"
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
