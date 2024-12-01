package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewResourceService creates a new instance of ResourceService
func NewResourceService(client *Client) *ResourceService {
	return &ResourceService{
		client: client,
	}
}

// GetAuthzResources gets all resources for the client's resource server
func (s *ResourceService) GetAuthzResources(realm, clientUUID string, params *ResourceQuery) ([]ResourceRepresentation, error) {
	url := s.client.buildURL("admin", "realms", realm, "clients", clientUUID, "authz", "resource-server", "resource")

	// Add query parameters if provided
	if params != nil {
		query := url.Query()
		if params.Deep != nil {
			query.Set("deep", fmt.Sprintf("%t", *params.Deep))
		}
		if params.ExactName != nil {
			query.Set("exactName", fmt.Sprintf("%t", *params.ExactName))
		}
		if params.First != nil {
			query.Set("first", fmt.Sprintf("%d", *params.First))
		}
		if params.Max != nil {
			query.Set("max", fmt.Sprintf("%d", *params.Max))
		}
		if params.Name != "" {
			query.Set("name", params.Name)
		}
		if params.Owner != "" {
			query.Set("owner", params.Owner)
		}
		if params.Type != "" {
			query.Set("type", params.Type)
		}
		if params.URI != "" {
			query.Set("uri", params.URI)
		}
		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var resources []ResourceRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return resources, nil
}

// CreateAuthzResource creates a new resource for the client's resource server
func (s *ResourceService) CreateAuthzResource(realm, clientUUID string, resource *ResourceRepresentation) (*ResourceRepresentation, error) {
	url := s.client.buildURL("admin", "realms", realm, "clients", clientUUID, "authz", "resource-server", "resource")

	body, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdResource ResourceRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&createdResource); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdResource, nil
}

// GetAuthzResource gets a specific resource by ID from the client's resource server
func (s *ResourceService) GetAuthzResource(realm, clientUUID, resourceID string) (*ResourceRepresentation, error) {
	url := s.client.buildURL("admin", "realms", realm, "clients", clientUUID, "authz", "resource-server", "resource", resourceID)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var resource ResourceRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&resource); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resource, nil
}

// UpdateAuthzResource updates a specific resource in the client's resource server
func (s *ResourceService) UpdateAuthzResource(realm, clientUUID, resourceID string, resource *ResourceRepresentation) error {
	url := s.client.buildURL("admin", "realms", realm, "clients", clientUUID, "authz", "resource-server", "resource", resourceID)

	body, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteAuthzResource deletes a specific resource from the client's resource server
func (s *ResourceService) DeleteAuthzResource(realm, clientUUID, resourceID string) error {
	url := s.client.buildURL("admin", "realms", realm, "clients", clientUUID, "authz", "resource-server", "resource", resourceID)

	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.Token)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
