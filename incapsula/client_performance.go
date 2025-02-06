package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const FORCE_RISKY_OP_HEADER_NAME = "force-risky-operation"

// PerformanceSettings is a struct that encompasses all the properties for performance settings
type PerformanceSettings struct {
	Mode struct {
		Level string `json:"level"`
		HTTPS string `json:"https,omitempty"`
		Time  int    `json:"time,omitempty"`
	} `json:"mode"`
	Key struct {
		UniteNakedFullCache bool `json:"unite_naked_full_cache"`
		ComplyVary          bool `json:"comply_vary"`
	} `json:"key,omitempty"`
	Response struct {
		StaleContent struct {
			Mode string `json:"mode,omitempty"`
			Time int    `json:"time,omitempty"`
		} `json:"stale_content,omitempty"`
		CacheShield         bool `json:"cache_shield"`
		CacheResponseHeader struct {
			Mode    string        `json:"mode,omitempty"`
			Headers []interface{} `json:"headers,omitempty"`
		} `json:"cache_response_header,omitempty"`
		TagResponseHeader    string `json:"tag_response_header,omitempty"`
		CacheEmptyResponses  bool   `json:"cache_empty_responses"`
		Cache300X            bool   `json:"cache_300x"`
		CacheHTTP10Responses bool   `json:"cache_http_10_responses"`
		Cache404             struct {
			Enabled bool `json:"enabled"`
			Time    int  `json:"time,omitempty"`
		} `json:"cache_404,omitempty"`
	} `json:"response,omitempty"`
	TTL struct {
		UseShortestCaching bool `json:"use_shortest_caching"`
		PreferLastModified bool `json:"prefer_last_modified"`
	} `json:"ttl,omitempty"`
	ClientSide struct {
		EnableClientSideCaching bool `json:"enable_client_side_caching"`
		ComplyNoCache           bool `json:"comply_no_cache"`
		SendAgeHeader           bool `json:"send_age_header"`
	} `json:"client_side,omitempty"`
}

// GetPerformanceSettings gets the site performance settings
func (c *Client) GetPerformanceSettings(siteID string) (*PerformanceSettings, error) {
	log.Printf("[INFO] Getting Incapsula Performance Settings for Site ID %s\n", siteID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache", c.config.BaseURLRev2, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadSitePerformance)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading Incap Performance Settings for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Incap Performance Settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Incap Performance Settings for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var performanceSettings PerformanceSettings
	err = json.Unmarshal([]byte(responseBody), &performanceSettings)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Performance Settings JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &performanceSettings, nil
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

	var headers = map[string]string{}
	if performanceSettings.Mode.HTTPS == "include_all_resources" && performanceSettings.Mode.Level == "all_resources" {
		log.Printf("[DEBUG] Incapsula Update Incap Performance Settings  - adding header: %s\n", FORCE_RISKY_OP_HEADER_NAME)
		headers[FORCE_RISKY_OP_HEADER_NAME] = strconv.FormatBool(true)
	}

	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache", c.config.BaseURLRev2, siteID)
	resp, err := c.DoJsonRequestWithCustomHeaders(http.MethodPut, reqURL, performanceSettingsJSON, headers, UpdateSitePerformance)
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
