package authz

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: validate audience and issuer

// TokenProvider represents a JWT token provider (Auth0, Keycloak, etc.)
type TokenProvider struct {
	jwksCache *JWKSCache
}

// NewTokenProvider creates a new token provider instance
func NewTokenProvider(jwksURL string) *TokenProvider {
	return &TokenProvider{
		jwksCache: NewJWKSCache(jwksURL),
	}
}

// FetchPublicKey returns a function that fetches and parses public keys
func (tp *TokenProvider) FetchPublicKey() func(string) (interface{}, error) {
	return func(kid string) (interface{}, error) {
		key, err := tp.jwksCache.FetchKey(kid)
		if err != nil {
			return nil, err
		}

		keyData, ok := key.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid key format")
		}

		return parseRSAPublicKey(keyData)
	}
}

// parseRSAPublicKey converts a JWKS key into an RSA public key
func parseRSAPublicKey(jwk map[string]interface{}) (*rsa.PublicKey, error) {
	// The "x5c" (X.509 Certificate Chain) contains the public key in base64-encoded DER format
	x5c, ok := jwk["x5c"].([]interface{})
	if !ok || len(x5c) == 0 {
		return nil, fmt.Errorf("x5c field is missing or invalid")
	}

	// Decode the first certificate in the chain
	certPEM := "-----BEGIN CERTIFICATE-----\n" + x5c[0].(string) + "\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the certificate to extract the public key
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain an RSA public key")
	}

	return rsaKey, nil
}

func (p *TokenProvider) ValidateToken(tokenString string, expectedAudience, expectedIssuer string) (map[string]interface{}, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithLeeway(30*time.Second), // Add some leeway for clock skew
	)

	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Get the key ID from the token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		// Get the public key from JWKS
		key, err := p.jwksCache.FetchKey(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch public key: %w", err)
		}

		keyData, ok := key.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid key format")
		}

		return parseRSAPublicKey(keyData)
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate issuer explicitly
	iss, _ := claims["iss"].(string)
	if expectedIssuer != "" && iss != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	// Validate token type in header (optional but recommended)
	if tokenType, ok := token.Header["typ"].(string); ok && tokenType != "JWT" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Validate audience if provided
	if expectedAudience != "" {
		aud, ok := claims["aud"]
		if !ok {
			return nil, fmt.Errorf("audience claim missing")
		}

		// Handle both string and array audience formats
		switch v := aud.(type) {
		case string:
			if v != expectedAudience {
				return nil, fmt.Errorf("invalid audience")
			}
		case []interface{}:
			valid := false
			for _, a := range v {
				if a.(string) == expectedAudience {
					valid = true
					break
				}
			}
			if !valid {
				return nil, fmt.Errorf("invalid audience")
			}
		default:
			return nil, fmt.Errorf("invalid audience format")
		}
	}

	// Validate additional security claims if present
	if nonce, ok := claims["nonce"].(string); ok && nonce == "" {
		return nil, fmt.Errorf("invalid nonce")
	}

	return claims, nil
}
