package clerk

import "strings"

type issuer struct {
	iss         string
	isSatellite bool
	proxyURL    string
}

func newIssuer(iss string) *issuer {
	return &issuer{
		iss: iss,
	}
}

func (iss *issuer) WithSatelliteDomain(isSatellite bool) *issuer {
	iss.isSatellite = isSatellite
	return iss
}

func (iss *issuer) WithProxyURL(proxyURL string) *issuer {
	iss.proxyURL = proxyURL
	return iss
}

func (iss *issuer) IsValid() bool {
	if iss.isSatellite {
		return true
	}

	if iss.proxyURL != "" {
		return iss.iss == iss.proxyURL
	}

	return strings.HasPrefix(iss.iss, "https://clerk.") ||
		strings.Contains(iss.iss, ".clerk.accounts")
}
