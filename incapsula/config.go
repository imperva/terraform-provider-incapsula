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
}

// Client configures and returns a fully initialized Incapsula Client
func (c *Config) Client() (interface{}, error) {
	log.Println("[INFO] Checking API credentials for client instantiation")

	// Check API Identifier
	if strings.TrimSpace(c.APIID) == "" {
		return nil, errors.New("API Identifier (api_id) must be provided")
	}

	// Check API Key
	if strings.TrimSpace(c.APIKey) == "" {
		return nil, errors.New("API Key (api_key) must be provided")
	}

	// Create client
	client := Client{Config: c}

	// Verify client credentials
	err := client.Verify()
	if err != nil {
		return nil, err
	}

	return client, nil
}
