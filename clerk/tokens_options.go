package clerk

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
)

// VerifyTokenOption describes a functional parameter for the VerifyToken method
type VerifyTokenOption func(*verifyTokenOptions) error

// WithAuthorizedParty allows to set the authorized parties to check against the azp claim of the session token
func WithAuthorizedParty(parties ...string) VerifyTokenOption {
	return func(o *verifyTokenOptions) error {
		authorizedParties := make(map[string]struct{})
		for _, party := range parties {
			authorizedParties[party] = struct{}{}
		}

		o.authorizedParties = authorizedParties
		return nil
	}
}

// WithLeeway allows to set a custom leeway that gives some extra time to the token to accomodate for clock skew, etc.
func WithLeeway(leeway time.Duration) VerifyTokenOption {
	return func(o *verifyTokenOptions) error {
		o.leeway = leeway
		return nil
	}
}

// WithJWTVerificationKey allows to set the JWK to use for verifying tokens without the need to download or cache any JWKs at runtime
func WithJWTVerificationKey(key string) VerifyTokenOption {
	return func(o *verifyTokenOptions) error {
		// From the Clerk docs: "Note that the JWT Verification key is not in
		// PEM format, the header and footer are missing, in order to be shorter
		// and single-line for easier setup."
		if !strings.HasPrefix(key, "-----BEGIN") {
			key = "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
		}

		jwk, err := pemToJWK(key)
		if err != nil {
			return err
		}

		o.jwk = jwk
		return nil
	}
}

// WithCustomClaims allows to pass a type (e.g. struct), which will be populated with the token claims based on json tags.
// For this option to work you must pass a pointer.
func WithCustomClaims(customClaims interface{}) VerifyTokenOption {
	return func(o *verifyTokenOptions) error {
		o.customClaims = customClaims
		return nil
	}
}

func pemToJWK(key string) (*jose.JSONWebKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM-encoded block")
	}

	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid key type, expected a public key")
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return &jose.JSONWebKey{Key: rsaPublicKey, Algorithm: "RS256"}, nil
}
