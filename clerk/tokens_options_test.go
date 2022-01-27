package clerk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAuthorizedPartyNone(t *testing.T) {
	opts := &verifyTokenOptions{}
	WithAuthorizedParty()(opts)

	assert.Len(t, opts.authorizedParties, 0)
}

func TestWithAuthorizedPartySingle(t *testing.T) {
	opts := &verifyTokenOptions{}
	WithAuthorizedParty("test-party")(opts)

	assert.Len(t, opts.authorizedParties, 1)
	assert.Equal(t, arrayToMap(t, []string{"test-party"}), opts.authorizedParties)
}

func TestWithAuthorizedPartyMultiple(t *testing.T) {
	authorizedParties := []string{"test-party", "another_party", "yet-another-party"}

	opts := &verifyTokenOptions{}
	WithAuthorizedParty(authorizedParties...)(opts)

	assert.Len(t, opts.authorizedParties, len(authorizedParties))
	assert.Equal(t, arrayToMap(t, authorizedParties), opts.authorizedParties)
}

func arrayToMap(t *testing.T, input []string) map[string]struct{} {
	t.Helper()

	output := make(map[string]struct{})
	for _, s := range input {
		output[s] = struct{}{}
	}

	return output
}
