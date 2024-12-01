package authz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type JWKSCache struct {
	mu      sync.RWMutex
	keys    map[string]interface{}
	jwksURL string
}

func NewJWKSCache(jwksURL string) *JWKSCache {
	return &JWKSCache{
		keys:    make(map[string]interface{}),
		jwksURL: jwksURL,
	}
}

func (j *JWKSCache) FetchKey(kid string) (interface{}, error) {
	j.mu.RLock()
	if key, found := j.keys[kid]; found {
		j.mu.RUnlock()
		return key, nil
	}
	j.mu.RUnlock()

	// Fetch and parse JWKS
	resp, err := http.Get(j.jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// Update cache
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, key := range jwks.Keys {
		if key["kid"] == kid {
			j.keys[kid] = key
			return key, nil
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}
