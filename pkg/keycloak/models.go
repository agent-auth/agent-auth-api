package keycloak

import "errors"

var ErrNotFound = errors.New("resource not found")

type ResourceService struct {
	client *Client
}

type ResourceRepresentation struct {
	ID          string                `json:"_id,omitempty"`
	Attributes  map[string][]string   `json:"attributes,omitempty"`
	DisplayName string                `json:"displayName,omitempty"`
	IconURI     string                `json:"icon_uri,omitempty"`
	Name        string                `json:"name,omitempty"`
	Owner       map[string]string     `json:"owner,omitempty"`
	Scopes      []ScopeRepresentation `json:"scopes,omitempty"`
	Type        string                `json:"type,omitempty"`
	URIs        []string              `json:"uris,omitempty"`
}

type ScopeRepresentation struct {
	DisplayName string                   `json:"displayName,omitempty"`
	IconURI     string                   `json:"iconUri,omitempty"`
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Policies    []PolicyRepresentation   `json:"policies,omitempty"`
	Resources   []ResourceRepresentation `json:"resources,omitempty"`
}

type PolicyRepresentation struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Type             string `json:"type,omitempty"`
	Logic            string `json:"logic,omitempty"`
	DecisionStrategy string `json:"decisionStrategy,omitempty"`
}

// Query parameters struct for filtering resources
type ResourceQuery struct {
	Deep        *bool  `query:"deep"`
	ExactName   *bool  `query:"exactName"`
	First       *int32 `query:"first"`
	MatchingURI *bool  `query:"matchingUri"`
	Max         *int32 `query:"max"`
	Name        string `query:"name,omitempty"`
	Owner       string `query:"owner,omitempty"`
	Scope       string `query:"scope,omitempty"`
	Type        string `query:"type,omitempty"`
	URI         string `query:"uri,omitempty"`
}
