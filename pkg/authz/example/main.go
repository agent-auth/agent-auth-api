package main

import (
	"fmt"
	"log"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
)

func main() {
	// Initialize the token provider with the JWKS URL
	provider := authz.NewTokenProvider("http://localhost:8081/realms/saas-ui-api-users/protocol/openid-connect/certs")

	// Example JWT token - replace with your actual token
	testToken := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ1Q2R0bE52YnFxbE5TWDg3VUNBQ2M4NlBFcDltVjRZNm44Tzh6QUstcmswIn0.eyJleHAiOjE3MzMwOTU1MDQsImlhdCI6MTczMzA5NTIwNCwianRpIjoiZTYxNTVmZTEtMGNhMC00MjMyLWJkZDctYjQxZWZlMzAyMWE1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgxL3JlYWxtcy9zYWFzLXVpLWFwaS11c2VycyIsInN1YiI6ImE5NjYzZmUwLTIxYjQtNDhlNS1iODU5LTlkNzZkMTliNGIwNSIsInR5cCI6IkJlYXJlciIsImF6cCI6InNhYXMtdWktYXBpLXVzZXJzIiwic2lkIjoiNTA1MmI5YWItZDVmOS00MTdiLTllNGQtYmZkNTNjN2JjODEwIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJzY29wZSI6ImVtYWlsIHByb2ZpbGUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJHVUZSQU4gQkFJRyIsInByZWZlcnJlZF91c2VybmFtZSI6Imd1ZnJhbm1pcnphMUBnbWFpbC5jb20iLCJnaXZlbl9uYW1lIjoiR1VGUkFOIiwiZmFtaWx5X25hbWUiOiJCQUlHIiwiZW1haWwiOiJndWZyYW5taXJ6YTFAZ21haWwuY29tIn0.dfRMvRwaSksU6HTEE2827_OUTnnmpILwN6_0da3PbspH7-_yKw9Ay2zMmnLTB70lvqitu1nzrFiAFpbgzqaYWLxuYv1YFSeXA_JLHFRWCjULDzUtaizkpBSdBvNe9B4UPUlmZOH9hi9Bs6QSpYarN9owr_NQwU6JUTpG93jbcEp7NG9mrIot-fpe-U2393AfHtHIJs6LDWhyFRkrzUsm85xN57QKm29TsRwHbgXoYKa4X_qPFkbNyVz8FiXhZOxDyhNvunxxXUMQ_BzStsLYuDHIJLxyqthJgoBB87YvCvozZB8NpB6BnvRlhj6FagOjeB2BnYdzeQbog6oNVqgTWw"

	// Validate the token
	claims, err := provider.ValidateToken(testToken, "", "")
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return
	}

	fmt.Println(claims)
	// Print the validated claims
	// fmt.Printf("Token validated successfully!\n")
	// fmt.Printf("Subject: %s\n", claims.Subject)
	// fmt.Printf("Issuer: %s\n", claims.Issuer)
	// fmt.Printf("Expires at: %v\n", claims.ExpiresAt)
}
