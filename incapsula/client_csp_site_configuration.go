package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// CspSiteConfig is the struct describing a csp site config response
type CspSiteConfig struct {
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	Discovery string `json:"discovery"`
	Settings  struct {
		Emails []CspSiteConfigEmail `json:"emails"`
	} `json:"settings"`
	TrackingIDs []struct {
		TrackingId   string `json:"trackingId"`
		DiscoveredMS int    `json:"discoveredMs"`
	} `json:"tracking-ids"`
}

type CspSiteConfigEmail struct {
	Email string `json:"email"`
}

const (
	CspSiteApiPath = "/csp-api/v1/sites"
)

// GetCspSite gets the csp site config
func (c *Client) GetCspSite(siteID int) (*CspSiteConfig, error) {
	log.Printf("[INFO] Getting CSP site configuration for site ID: %d\n", siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s%s/%d", c.config.BaseURLAPI, CspSiteApiPath, siteID),
		nil)
	if err != nil {
		return nil, fmt.Errorf("Error from CSP API for when reading site ID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API Read Site Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from CSP API when reading site config for ID %d: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var cspSiteConfig CspSiteConfig
	err = json.Unmarshal([]byte(responseBody), &cspSiteConfig)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &cspSiteConfig, nil
}

// UpdateCspSite gets the csp site config
func (c *Client) UpdateCspSite(siteID int, config *CspSiteConfig) (*CspSiteConfig, error) {
	log.Printf("[INFO] Updating CSP site configuration for site ID: %d\n%v", siteID, config)
	configJSON, err := json.Marshal(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal csp api site config: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut,
		fmt.Sprintf("%s/%s/%d", c.config.BaseURLAPI, CspSiteApiPath, siteID),
		configJSON)

	if err != nil {
		return nil, fmt.Errorf("Error from CSP API while updating site configuration for site ID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API Update Site Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from CSP API when reading site config for ID %d: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var cspSiteConfig CspSiteConfig
	err = json.Unmarshal([]byte(responseBody), &cspSiteConfig)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &cspSiteConfig, nil
}
