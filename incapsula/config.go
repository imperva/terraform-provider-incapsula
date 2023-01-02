package incapsula

import (
	"errors"
	"log"
	"strings"
)

// Config represents the configuration required for the Incapsula Client
type Config struct {
	// API Identifier
	APIID string

	// API Key
	APIKey string

	// Base URL (no trailing slash)
	// This endpoint is unlikely to change in the near future
	BaseURL string

	// Base URL Revision 2 (no trailing slash)
	// Updates to APIv1 are underway and newer resources are supported
	// Rev2 includes the move to Swagger, appropriate method verbs (not everything is a post)
	// The other endpoints will eventually move over but we'll need the following for now
	BaseURLRev2 string

	// Base URL API
	// API V2
	// Same as revision 2 but with a different subdomain
	BaseURLAPI string
}

var missingAPIIDMessage = "API Identifier (api_id) must be provided"
var missingAPIKeyMessage = "API Key (api_key) must be provided"
var missingBaseURLMessage = "Base URL must be provided"
var missingBaseURLRev2Message = "Base URL Revision 2 must be provided"
var missingBaseURLAPIMessage = "Base URL API must be provided"

// Client configures and returns a fully initialized Incapsula Client
func (c *Config) Client() (interface{}, error) {
	log.Println("[INFO] Checking API credentials for client instantiation")

	// Check API Identifier
	if strings.TrimSpace(c.APIID) == "" {
		return nil, errors.New(missingAPIIDMessage)
	}

	// Check API Key
	if strings.TrimSpace(c.APIKey) == "" {
		return nil, errors.New(missingAPIKeyMessage)
	}

	// Check Base URL
	if strings.TrimSpace(c.BaseURL) == "" {
		return nil, errors.New(missingBaseURLMessage)
	}

	// Check Base URL Revision 2
	if strings.TrimSpace(c.BaseURLRev2) == "" {
		return nil, errors.New(missingBaseURLRev2Message)
	}

	// Check Base URL API
	if strings.TrimSpace(c.BaseURLAPI) == "" {
		return nil, errors.New(missingBaseURLAPIMessage)
	}

	// Create client
	client := NewClient(c)

	// Verify client credentials
	accountStatusResponse, err := client.Verify()
	client.accountStatus = accountStatusResponse
	if err != nil {
		return nil, err
	}

	return client, nil
}
