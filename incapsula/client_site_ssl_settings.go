package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type HSTSConfiguration struct {
	IsEnabled          bool `json:"isEnabled"`
	MaxAge             int  `json:"maxAge"`
	SubDomainsIncluded bool `json:"subDomainsIncluded"`
	PreLoaded          bool `json:"preLoaded"`
}

type Data struct {
	HstsConfiguration HSTSConfiguration `json:"hstsConfiguration"`
}

type SSLSettingsDTO struct {
	Data []Data `json:"data"`
}

func (c *Client) UpdateSiteSSLSettings(siteID int, mySSLSettings SSLSettingsDTO) (*SSLSettingsDTO, error) {
	log.Printf("[INFO] Updating Incapsula Site SSL settings for Site ID %d\n", siteID)

	requestJSON, err := json.Marshal(mySSLSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal HSTSConfiguration: %s", err)
	}

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%d/settings/TLSConfiguration", c.config.BaseURLRev3, siteID)
	log.Printf("[INFO] HSTS request json looks like this %s\n", requestJSON)
	log.Printf("[INFO] HSTS request URL looks like this %s\n", reqURL)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, requestJSON, UpdateSiteSSLSettings)
	if err != nil {
		return nil, fmt.Errorf("error from Incapsula service when updating Site SSL settings %s for Site ID %d: %s", requestJSON, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Site SSL settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating Site SSL settings %s for Site ID %d: %s", resp.StatusCode, requestJSON, siteID, string(responseBody))
	}

	// Parse the JSON
	var sslSettingsDTO SSLSettingsDTO
	err = json.Unmarshal([]byte(responseBody), &sslSettingsDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Site settings JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &sslSettingsDTO, nil
}

func (c *Client) ReadSiteSSLSettings(siteID int) (*SSLSettingsDTO, int, error) {
	log.Printf("[INFO] Getting Incapsula Incap SSL settings for Site ID %d\n", siteID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%d/settings/TLSConfiguration", c.config.BaseURLRev3, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadSiteSSLSettings)
	if err != nil {
		return nil, 0, fmt.Errorf("error from Incapsula service when reading SSL Settings for Site ID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Site SSL settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("error status code %d from Incapsula service when reading SSL settings for Site ID %d: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var sslSettingsDTO SSLSettingsDTO
	err = json.Unmarshal([]byte(responseBody), &sslSettingsDTO)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error parsing Site SSL settings JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &sslSettingsDTO, resp.StatusCode, nil
}
