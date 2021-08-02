package clerk

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type jwksCache struct {
	sync.RWMutex
	jwks      *JWKS
	expiresAt time.Time
}

func (c *jwksCache) isInvalid() bool {
	c.RLock()
	defer c.RUnlock()

	return c.jwks == nil || len(c.jwks.Keys) == 0 || time.Now().After(c.expiresAt)
}

func (c *jwksCache) set(jwks *JWKS) {
	c.Lock()
	defer c.Unlock()

	c.jwks = jwks
	c.expiresAt = time.Now().Add(time.Hour)
}

func (c *jwksCache) get(kid string) (*jose.JSONWebKey, error) {
	c.RLock()
	defer c.RUnlock()

	for _, key := range c.jwks.Keys {
		if key.KeyID == kid {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("no jwk key found for kid %s", kid)
}
