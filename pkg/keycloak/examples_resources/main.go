package main

import (
	"log"

	"github.com/agent-auth/agent-auth-api/pkg/keycloak"
)

func main() {
	// Initialize the client
	client, err := keycloak.NewClient("http://localhost:8081", "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJpczFVZzAwSF95ZUFMczdMbXY4R2IzaUJCQXQ3emxJcjhSby1ES0g5aUJFIn0.eyJleHAiOjE3MzMwNTcxMjEsImlhdCI6MTczMzA1NTMyMSwianRpIjoiYTAyOGM0ZjMtMmVhOS00YWMzLWE4NjQtMDdlMzk5MGFmOTE5IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgxL3JlYWxtcy9tYXN0ZXIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhZG1pbi1jbGkiLCJzaWQiOiJlMjI5YzAzZC0wYWE0LTQyMWQtYTA0Yi0wYzVjYjA0YzMxNGQiLCJzY29wZSI6InByb2ZpbGUgZW1haWwifQ.HgS6n273j-yHfTELdqVew2bFzz7XinMbUgLlWcG-nbUyBAzHR3OvQYD9WJ0NU0YOeZh3RicQDx2ixUusv4zVFqfLCwiprVxLdrHX_nsJUrgnEZuwR6EfCV0VIJJ2541SX7aZ7B6daVnsl2213BDD0GQWdkAcCSs7dakpNu35-zay_3o9qmXovFFPknEVa5zk5uxBMJVfxvtR6AezI9DB7RIyry6pvMMfCWxraZ2C1ABsiswA1AqhwanLft7qQmfPvbnvpfFg5dcf-YLQ5U9aeKmqJAJfkSoAfJnihURePdesOWY8hYMArJso7bs1aEoTi3PTUhmWgn7p13v2CILPyg")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create resource service
	resourceService := keycloak.NewResourceService(client)

	// Define test parameters
	realmName := "master"
	clientID := "0d37d66f-dffd-4c87-9301-28b49abc9c7a"

	// 1. Get all resources
	resources, err := resourceService.GetAuthzResources(realmName, clientID, &keycloak.ResourceQuery{})
	if err != nil {
		log.Fatalf("Failed to get resources: %v", err)
	}
	log.Printf("All resources: %v\n", resources)

	// 2. Create a new resource
	newResource := &keycloak.ResourceRepresentation{
		Name:        "test-resource-xx",
		DisplayName: "Test Resource-xx",
		Type:        "urn:test-resource:resources:test-xx",
		URIs:        []string{"/api/test"},
		Scopes:      []keycloak.ScopeRepresentation{},
	}

	createdResource, err := resourceService.CreateAuthzResource(realmName, clientID, newResource)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}
	log.Printf("Created resource: %v\n", createdResource)

	// 3. Get specific resource by ID
	resource, err := resourceService.GetAuthzResource(realmName, clientID, createdResource.ID)
	if err != nil {
		log.Fatalf("Failed to get resource: %v", err)
	}
	log.Printf("Retrieved resource: %v\n", resource)

	// 4. Update resource
	resource.DisplayName = "Updated Test Resource"
	err = resourceService.UpdateAuthzResource(realmName, clientID, resource.ID, resource)
	if err != nil {
		log.Fatalf("Failed to update resource: %v", err)
	}
	log.Printf("Updated resource: %v\n", resource)

	// 5. Delete resource
	err = resourceService.DeleteAuthzResource(realmName, clientID, resource.ID)
	if err != nil {
		log.Fatalf("Failed to delete resource: %v", err)
	}
	log.Println("Resource deleted successfully")

	// 6. Search resources with query parameters
	query := &keycloak.ResourceQuery{}
	searchResults, err := resourceService.GetAuthzResources(realmName, clientID, query)
	if err != nil {
		log.Fatalf("Failed to search resources: %v", err)
	}
	log.Printf("Search results: %v\n", searchResults)
}
