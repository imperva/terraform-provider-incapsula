package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// PerformanceSettings is a struct that encompasses all the properties for performance settings
type PerformanceSettings struct {
	Mode struct {
		Level string `json:"level"`
		HTTPS string `json:"https"`
		Time  int    `json:"time,omitempty"`
	} `json:"mode"`
	Key struct {
		UniteNakedFullCache bool `json:"unite_naked_full_cache"`
		ComplyVary          bool `json:"comply_vary"`
	} `json:"key"`
	Response struct {
		StaleContent struct {
			Mode string `json:"mode"`
			Time int    `json:"time,omitempty"`
		} `json:"stale_content"`
		CacheShield         bool `json:"cache_shield"`
		CacheResponseHeader struct {
			Mode    string        `json:"mode"`
			Headers []interface{} `json:"headers"`
		} `json:"cache_response_header"`
		TagResponseHeader    string `json:"tag_response_header"`
		CacheEmptyResponses  bool   `json:"cache_empty_responses"`
		Cache300X            bool   `json:"cache_300x"`
		CacheHTTP10Responses bool   `json:"cache_http_10_responses"`
		Cache404             struct {
			Enabled bool `json:"enabled"`
			Time    int  `json:"time,omitempty"`
		} `json:"cache_404"`
	} `json:"response"`
	TTL struct {
		UseShortestCaching bool `json:"use_shortest_caching"`
		PreferLastModified bool `json:"prefer_last_modified"`
	} `json:"ttl"`
	ClientSide struct {
		EnableClientSideCaching bool `json:"enable_client_side_caching"`
		ComplyNoCache           bool `json:"comply_no_cache"`
		SendAgeHeader           bool `json:"send_age_header"`
	} `json:"client_side"`
}

// GetPerformanceSettings gets the site performance settings
func (c *Client) GetPerformanceSettings(siteID string) (*PerformanceSettings, int, error) {
	log.Printf("[INFO] Getting Incapsula Performance Settings for Site ID %s\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/sites/%s/settings/cache?api_id=%s&api_key=%s", c.config.APIV2BaseURL, siteID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, 0, fmt.Errorf("Error from Incapsula service when reading Incap Performance Settings for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Incap Performance Settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("Error status code %d from Incapsula service when reading Incap Performance Settings for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var performanceSettings PerformanceSettings
	err = json.Unmarshal([]byte(responseBody), &performanceSettings)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error parsing Incap Performance Settings JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &performanceSettings, resp.StatusCode, nil
}

// UpdatePerformanceSettings updates the site performance settings
func (c *Client) UpdatePerformanceSettings(siteID string, performanceSettings *PerformanceSettings) (*PerformanceSettings, error) {
	log.Printf("[INFO] Updating Incapsula Performance Settings for Site ID %s\n", siteID)

	performanceSettingsJSON, err := json.Marshal(performanceSettings)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal PerformanceSettings: %s", err)
	}

	// Post request to Incapsula
	log.Printf("[DEBUG] Incapsula Update Incap Performance Settings JSON request: %s\n", string(performanceSettingsJSON))
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/sites/%s/settings/cache?api_id=%s&api_key=%s", c.config.APIV2BaseURL, siteID, c.config.APIID, c.config.APIKey),
		bytes.NewReader(performanceSettingsJSON))
	if err != nil {
		return nil, fmt.Errorf("Error preparing HTTP POST for updating Incap Performance Settings for Site ID %s: %s", siteID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating Incap Performance Settings for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Incap Performance Settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating Incap Performance Settings for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var updatedPerformanceSettings PerformanceSettings
	err = json.Unmarshal([]byte(responseBody), &updatedPerformanceSettings)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Performance Settings JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &updatedPerformanceSettings, nil
}
