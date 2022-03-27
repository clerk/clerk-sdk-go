package clerk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithAuthorizedPartyNone(t *testing.T) {
	opts := &verifyTokenOptions{}
	err := WithAuthorizedParty()(opts)

	if assert.NoError(t, err) {
		assert.Len(t, opts.authorizedParties, 0)
	}
}

func TestWithAuthorizedPartySingle(t *testing.T) {
	opts := &verifyTokenOptions{}
	err := WithAuthorizedParty("test-party")(opts)

	if assert.NoError(t, err) {
		assert.Len(t, opts.authorizedParties, 1)
		assert.Equal(t, arrayToMap(t, []string{"test-party"}), opts.authorizedParties)
	}
}

func TestWithAuthorizedPartyMultiple(t *testing.T) {
	authorizedParties := []string{"test-party", "another_party", "yet-another-party"}

	opts := &verifyTokenOptions{}
	err := WithAuthorizedParty(authorizedParties...)(opts)

	if assert.NoError(t, err) {
		assert.Len(t, opts.authorizedParties, len(authorizedParties))
		assert.Equal(t, arrayToMap(t, authorizedParties), opts.authorizedParties)
	}
}

func TestWithLeeway(t *testing.T) {
	leeway := 5 * time.Second

	opts := &verifyTokenOptions{}
	err := WithLeeway(leeway)(opts)

	if assert.NoError(t, err) {
		assert.Equal(t, opts.leeway, leeway)
	}
}

func TestWithJWTVerificationKey(t *testing.T) {
	key := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm7Zs5PFGrsrmvys1hHkSDOYoghz9+z9o+E6WgMqR+R/Af0/QRqQo/YwCmzB+01+5Us1NdSa32YuQYiMxV4T+g3eebSiBqPNiCyjl2wttCm5LAV5iHyVqwnBNcrXlA5mRFQz8lmyfpoksNDEVzJPwwHzPjKSIKsGgsrPnw6XsyOPJY/8UocscEcHptTmahHrbfNZLN0FrMneHw9tnn2AiUctuU9bw80KwPd+WFdZ6UZF/kPxVFsANJpz1aMpz7Lxi3Sz1ztUCdHvNJitRUO1Qewby4xi9DfIEECMq78LLmwGaTiKxutC6KwHLJEcbUblOJHpYVEXdBex9xGJ/2DHrBQIDAQAB"

	opts := &verifyTokenOptions{}
	err := WithJWTVerificationKey(key)(opts)

	if assert.NoError(t, err) {
		assert.Equal(t, "RS256", opts.jwk.Algorithm)
	}
}

func TestWithJWTVerificationKey_PEM(t *testing.T) {
	key := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm7Zs5PFGrsrmvys1hHkS
DOYoghz9+z9o+E6WgMqR+R/Af0/QRqQo/YwCmzB+01+5Us1NdSa32YuQYiMxV4T+
g3eebSiBqPNiCyjl2wttCm5LAV5iHyVqwnBNcrXlA5mRFQz8lmyfpoksNDEVzJPw
wHzPjKSIKsGgsrPnw6XsyOPJY/8UocscEcHptTmahHrbfNZLN0FrMneHw9tnn2Ai
UctuU9bw80KwPd+WFdZ6UZF/kPxVFsANJpz1aMpz7Lxi3Sz1ztUCdHvNJitRUO1Q
ewby4xi9DfIEECMq78LLmwGaTiKxutC6KwHLJEcbUblOJHpYVEXdBex9xGJ/2DHr
BQIDAQAB
-----END PUBLIC KEY-----`

	opts := &verifyTokenOptions{}
	err := WithJWTVerificationKey(key)(opts)

	if assert.NoError(t, err) {
		assert.Equal(t, "RS256", opts.jwk.Algorithm)
	}
}

func arrayToMap(t *testing.T, input []string) map[string]struct{} {
	t.Helper()

	output := make(map[string]struct{})
	for _, s := range input {
		output[s] = struct{}{}
	}

	return output
}
