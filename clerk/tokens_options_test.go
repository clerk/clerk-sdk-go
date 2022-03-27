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

func arrayToMap(t *testing.T, input []string) map[string]struct{} {
	t.Helper()

	output := make(map[string]struct{})
	for _, s := range input {
		output[s] = struct{}{}
	}

	return output
}
