package authz

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Add this at the package level
type contextKey string

const ClaimsContextKey = contextKey("claims")

// AuthMiddleware creates a new authentication middleware using the provided token provider
func AuthMiddleware(provider *TokenProvider, audience, issuer string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader || tokenString == "" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Validate the token
			claims, err := provider.ValidateToken(tokenString, audience, issuer)
			if err != nil {
				fmt.Println("Error validating token:", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Store claims in context for later use
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles creates middleware that checks if the user has any of the specified roles
func RequireRoles(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(map[string]interface{})
			if !ok {
				http.Error(w, "Unauthorized: no claims found", http.StatusUnauthorized)
				return
			}

			// Get roles from claims (adjust the claim key based on your JWT structure)
			userRoles, ok := claims["roles"].([]interface{})
			if !ok {
				http.Error(w, "Unauthorized: no roles found", http.StatusForbidden)
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, userRole := range userRoles {
				roleStr, ok := userRole.(string)
				if !ok {
					continue
				}
				for _, requiredRole := range roles {
					if Role(roleStr) == requiredRole {
						hasRole = true
						break
					}
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: insufficient roles", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireScopes creates middleware that checks if the user has all the specified scopes
func RequireScopes(scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(map[string]interface{})
			if !ok {
				http.Error(w, "Unauthorized: no claims found", http.StatusUnauthorized)
				return
			}

			// Get scopes from claims (adjust based on your JWT structure)
			// Some providers use space-separated string, others use array
			var userScopes []string
			switch scopeClaim := claims["scope"].(type) {
			case string:
				userScopes = strings.Split(scopeClaim, " ")
			case []interface{}:
				for _, s := range scopeClaim {
					if str, ok := s.(string); ok {
						userScopes = append(userScopes, str)
					}
				}
			default:
				http.Error(w, "Unauthorized: invalid scope format", http.StatusForbidden)
				return
			}

			// Check if user has all required scopes
			for _, requiredScope := range scopes {
				hasScope := false
				for _, userScope := range userScopes {
					if userScope == requiredScope {
						hasScope = true
						break
					}
				}
				if !hasScope {
					http.Error(w, "Forbidden: missing required scope", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
