package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MaskingSettings is a struct that encompasses all the properties of Site Masking Settings
type MaskingSettings struct {
	HashingEnabled bool   `json:"hashing_enabled,omitempty"`
	HashSalt       string `json:"hash_salt,omitempty"`
}

// GetMaskingSettings gets the site masking settings
func (c *Client) GetMaskingSettings(siteID string) (*MaskingSettings, error) {
	log.Printf("[INFO] Getting Incapsula Masking Settings for Site ID %s\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/sites/%s/settings/masking?api_id=%s&api_key=%s", c.config.APIV2BaseURL, siteID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading masking settings for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula ReadMaskingSettings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading masking settings for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var maskingSettings MaskingSettings
	err = json.Unmarshal([]byte(responseBody), &maskingSettings)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap masking settings JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &maskingSettings, nil
}

// UpdateMaskingSettings updates the site masking settings
func (c *Client) UpdateMaskingSettings(siteID string, maskingSettings *MaskingSettings) error {
	log.Printf("[INFO] Updating Incapsula masking settings for Site ID %s\n", siteID)

	maskingSettingsJSON, err := json.Marshal(maskingSettings)
	if err != nil {
		return fmt.Errorf("Failed to JSON marshal MaskingSettings: %s", err)
	}

	// Put request to Incapsula
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/sites/%s/settings/masking?api_id=%s&api_key=%s", c.config.APIV2BaseURL, siteID, c.config.APIID, c.config.APIKey),
		bytes.NewReader(maskingSettingsJSON))
	if err != nil {
		return fmt.Errorf("Error preparing HTTP PUT for updating Incap masking settings for Site ID %s: %s", siteID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when updating masking settings for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula UpdateMaskingSettings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when updating masking settings for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	return nil
}
