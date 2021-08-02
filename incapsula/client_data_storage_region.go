package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

const endpointDataStorageRegionGet = "sites/data-privacy/show"
const endpointDataStorageRegionUpdate = "sites/data-privacy/region-change"

// DataStorageRegionResponse contains the relevant information when getting/setting a data storage region
type DataStorageRegionResponse struct {
	Region     string `json:"region"`
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	DebugInfo  struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// GetDataStorageRegion gets the data storage region for the site
func (c *Client) GetDataStorageRegion(siteID string) (*DataStorageRegionResponse, error) {
	log.Printf("[INFO] Getting Incapsula data storage region for site: %s\n", siteID)

	// Post form to Incapsula
	values := url.Values{"site_id": {siteID}}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataStorageRegionGet)
	resp, err := c.PostFormWithHeaders(reqURL, values)
	if err != nil {
		return nil, fmt.Errorf("Error getting data storage region for site id: %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula data storage region JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataStorageRegionResponse DataStorageRegionResponse
	err = json.Unmarshal([]byte(responseBody), &dataStorageRegionResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing site data storage region JSON response for site id: %s: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if dataStorageRegionResponse.Res != 0 {
		return &dataStorageRegionResponse, fmt.Errorf("Error from Incapsula service when getting site data storage region for site id: %s: %s", siteID, string(responseBody))
	}

	return &dataStorageRegionResponse, nil
}

// UpdateDataStorageRegion will update the data storage region on the site
func (c *Client) UpdateDataStorageRegion(siteID, region string) (*DataStorageRegionResponse, error) {
	log.Printf("[INFO] Updating Incapsula site data storage region (%s) for siteID: %s\n", region, siteID)

	// Post form to Incapsula
	values := url.Values{
		"site_id":             {siteID},
		"data_storage_region": {region},
	}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataStorageRegionUpdate)
	resp, err := c.PostFormWithHeaders(reqURL, values)
	if err != nil {
		return nil, fmt.Errorf("Error updating data storage region with value (%s) on site_id: %s: %s", region, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update site data storage region JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataStorageRegionResponse DataStorageRegionResponse
	err = json.Unmarshal([]byte(responseBody), &dataStorageRegionResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update site data storage region JSON response for siteID %s: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if dataStorageRegionResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating site data storage region for siteID %s: %s", siteID, string(responseBody))
	}

	return &dataStorageRegionResponse, nil
}
