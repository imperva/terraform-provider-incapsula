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
	BaseURL string
}

var missingAPIIDMessage = "API Identifier (api_id) must be provided"
var missingAPIKeyMessage = "API Key (api_key) must be provided"
var missingBaseURLMessage = "Base URL must be provided"

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

	// Create client
	client := NewClient(c)

	// Verify client credentials
	_, err := client.Verify()
	if err != nil {
		return nil, err
	}

	return client, nil
}
