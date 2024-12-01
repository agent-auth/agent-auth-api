package keycloak

import (
	"net/http"
	"net/url"
	"path"
)

// Client represents a Keycloak API client
type Client struct {
	// BaseURL is the base URL of the Keycloak server
	BaseURL *url.URL

	// HTTPClient is the HTTP client used for making requests
	HTTPClient *http.Client

	// Token is the bearer token used for authentication
	Token string
}

// NewClient creates a new Keycloak client instance
func NewClient(baseURL string, token string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL:    parsedURL,
		HTTPClient: &http.Client{},
		Token:      token,
	}, nil
}

// buildURL constructs the full URL for an API request
func (c *Client) buildURL(parts ...string) *url.URL {
	result := *c.BaseURL

	// Combine all parts into a single path
	elements := append([]string{result.Path}, parts...)
	result.Path = path.Join(elements...)

	return &result
}
