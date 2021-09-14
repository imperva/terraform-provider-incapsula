package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const siteConfigUrl = "/api-security/config/site/"

type ApiSecuritySiteConfigGetResponse struct {
	Value struct {
		SiteId                                    int              `json:"siteId"`
		AccountId                                 int              `json:"accountId"`
		SiteName                                  string           `json:"siteName"`
		ApiOnlySite                               bool             `json:"apiOnlySite"`
		NonApiRequestViolationAction              string           `json:"nonApiRequestViolationAction"`
		LastModified                              int64            `json:"lastModified"`
		ViolationActions                          ViolationActions `json:"violationActions"`
		IsAutomaticDiscoveryApiIntegrationEnabled bool             `json:"isAutomaticDiscoveryApiIntegrationEnabled"`
		DiscoveryEnabled                          bool             `json:"discoveryEnabled"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

type ApiSecuritySiteConfigPostResponse struct {
	Value struct {
		SiteId int `json:"siteId"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

type ApiSecuritySiteConfigPostPayload struct {
	ApiOnlySite                               bool             `json:"apiOnlySite"`
	SiteName                                  string           `json:"siteName"`
	NonApiRequestViolationAction              string           `json:"nonApiRequestViolationAction"`
	ViolationActions                          ViolationActions `json:"violationActions"`
	IsAutomaticDiscoveryApiIntegrationEnabled bool             `json:"isAutomaticDiscoveryApiIntegrationEnabled"`
}

// ReadApiSecuritySiteConfig gets the Api-Security Site Config
func (c *Client) ReadApiSecuritySiteConfig(siteId int) (*ApiSecuritySiteConfigGetResponse, error) {
	log.Printf("[INFO] Getting Incapsula Api-Security Site Config: %d\n", siteId)

	// Post form to Incapsula
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s%s%d", c.config.BaseURLAPI, siteConfigUrl, siteId),
		nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]Error from Incapsula service while reading Api-Security Site Config for site ID %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Api-Security Site Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Api-Security Site Config for site ID %d: %s", resp.StatusCode, siteId, string(responseBody))
	}

	// Parse the JSON
	var siteConfigGetResponse ApiSecuritySiteConfigGetResponse
	err = json.Unmarshal(responseBody, &siteConfigGetResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing GET Api-Security Site Config JSON response for site ID %d: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &siteConfigGetResponse, nil
}

// UpdateApiSecuritySiteConfig updates an Api-Security Site Config
func (c *Client) UpdateApiSecuritySiteConfig(siteId int, siteConfigPayload *ApiSecuritySiteConfigPostPayload) (*ApiSecuritySiteConfigPostResponse, error) {
	siteConfigJSON, err := json.Marshal(siteConfigPayload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal api security site config: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s"+siteConfigUrl+"%d", c.config.BaseURLAPI, siteId),
		siteConfigJSON)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while updating API security site configuration for site ID %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update api-security site configuration JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating api-security site configuration: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var response ApiSecuritySiteConfigPostResponse
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing API security JSON response: %s\nresponse: %s", err, string(responseBody))
	}
	return &response, nil
}
