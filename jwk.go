package clerk

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/go-jose/go-jose/v3"
)

type JSONWebKeySet struct {
	APIResource
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	APIResource
	raw       jose.JSONWebKey
	Key       any    `json:"key"`
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
	Use       string `json:"use"`
}

func (k *JSONWebKey) UnmarshalJSON(data []byte) error {
	err := k.raw.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	k.Key = k.raw.Key
	k.KeyID = k.raw.KeyID
	k.Algorithm = k.raw.Algorithm
	k.Use = k.raw.Use
	return nil
}

// JSONWebKeyFromPEM returns a JWK from an RSA key.
func JSONWebKeyFromPEM(key string) (*JSONWebKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM-encoded block")
	}

	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid key type, expected a public key")
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &JSONWebKey{
		Key:       rsaPublicKey,
		Algorithm: "RS256",
	}, nil
}
