package enterprisedomainverification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestEnterpriseDomainVerificationClientPrepareDNSTxt(t *testing.T) {
	t.Parallel()
	id := "entd_ver_123"
	domain := "example.com"
	strategy := "dns_txt"
	record := "clerk-verification=abc123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"strategy":"%s","domain":"%s"}`, strategy, domain)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_domain_verification","domain":"%s","enterprise_connection_id":null,"enterprise_domain_id":null,"verified":false,"verification":{"status":"unverified","strategy":"%s","attempts":0,"expire_at":1715520000000,"dns_txt_record":"%s"}}`, id, domain, strategy, record)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_domain_verifications",
		},
	}
	client := NewClient(config)
	res, err := client.Prepare(context.Background(), &PrepareParams{
		Strategy: clerk.String(strategy),
		Domain:   clerk.String(domain),
	})
	require.NoError(t, err)
	require.Equal(t, id, res.ID)
	require.Equal(t, domain, res.Domain)
	require.False(t, res.Verified)
	require.Nil(t, res.EnterpriseConnectionID)
	require.Nil(t, res.EnterpriseDomainID)
	require.NotNil(t, res.Verification)
	require.Equal(t, "unverified", res.Verification.Status)
	require.Equal(t, strategy, res.Verification.Strategy)
	require.NotNil(t, res.Verification.DNSTxtRecord)
	require.Equal(t, record, *res.Verification.DNSTxtRecord)
}

func TestEnterpriseDomainVerificationClientPrepareEmailCode(t *testing.T) {
	t.Parallel()
	id := "entd_ver_456"
	domain := "example.com"
	strategy := "email_code"
	storedStrategy := "enterprise_email_code"
	email := "admin@example.com"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"strategy":"%s","domain":"%s","affiliation_email_address":"%s"}`, strategy, domain, email)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_domain_verification","domain":"%s","enterprise_connection_id":null,"enterprise_domain_id":null,"verified":false,"verification":{"status":"unverified","strategy":"%s","attempts":0,"expire_at":1715520000000}}`, id, domain, storedStrategy)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_domain_verifications",
		},
	}
	client := NewClient(config)
	res, err := client.Prepare(context.Background(), &PrepareParams{
		Strategy:                clerk.String(strategy),
		Domain:                  clerk.String(domain),
		AffiliationEmailAddress: clerk.String(email),
	})
	require.NoError(t, err)
	require.Equal(t, id, res.ID)
	require.Equal(t, domain, res.Domain)
	require.False(t, res.Verified)
	require.NotNil(t, res.Verification)
	require.Equal(t, storedStrategy, res.Verification.Strategy)
	require.Nil(t, res.Verification.DNSTxtRecord)
}

func TestEnterpriseDomainVerificationClientAttemptDNSTxt(t *testing.T) {
	t.Parallel()
	id := "entd_ver_789"
	domain := "example.com"
	strategy := "dns_txt"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"strategy":"%s"}`, strategy)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_domain_verification","domain":"%s","enterprise_connection_id":null,"enterprise_domain_id":null,"verified":true,"verification":{"status":"verified","strategy":"%s","attempts":1,"expire_at":1715520000000}}`, id, domain, strategy)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_domain_verifications/" + id + "/attempt_verification",
		},
	}
	client := NewClient(config)
	res, err := client.Attempt(context.Background(), &AttemptParams{
		VerificationID: id,
		Strategy:       clerk.String(strategy),
	})
	require.NoError(t, err)
	require.Equal(t, id, res.ID)
	require.True(t, res.Verified)
	require.Equal(t, "verified", res.Verification.Status)
}

func TestEnterpriseDomainVerificationClientAttemptEmailCode(t *testing.T) {
	t.Parallel()
	id := "entd_ver_abc"
	domain := "example.com"
	strategy := "email_code"
	storedStrategy := "enterprise_email_code"
	code := "424242"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"strategy":"%s","code":"%s"}`, strategy, code)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_domain_verification","domain":"%s","enterprise_connection_id":null,"enterprise_domain_id":null,"verified":true,"verification":{"status":"verified","strategy":"%s","attempts":1,"expire_at":1715520000000}}`, id, domain, storedStrategy)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_domain_verifications/" + id + "/attempt_verification",
		},
	}
	client := NewClient(config)
	res, err := client.Attempt(context.Background(), &AttemptParams{
		VerificationID: id,
		Strategy:       clerk.String(strategy),
		Code:           clerk.String(code),
	})
	require.NoError(t, err)
	require.Equal(t, id, res.ID)
	require.True(t, res.Verified)
}
