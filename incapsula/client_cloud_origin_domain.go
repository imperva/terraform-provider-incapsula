package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// CloudOriginDomainResponse contains the response for a single cloud origin domain
type CloudOriginDomainResponse struct {
	Value struct {
		OriginID             int    `json:"originId"`
		Domain               string `json:"domain"`
		Region               string `json:"region"`
		Port                 int    `json:"port"`
		ImpervaOriginDomain  string `json:"impervaOriginDomain"`
		Status               string `json:"status"`
		CreatedAt            string `json:"createdAt"`
		UpdatedAt            string `json:"updatedAt"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

// CloudOriginDomainCreateRequest contains the request payload for creating a cloud origin domain
type CloudOriginDomainCreateRequest struct {
	Domain string `json:"domain"`
	Region string `json:"region"`
	Port   int    `json:"port"`
}

// CloudOriginDomainUpdateRequest contains the request payload for updating a cloud origin domain
type CloudOriginDomainUpdateRequest struct {
	Region string `json:"region"`
	Port   int    `json:"port"`
}

// CreateCloudOriginDomain creates a new cloud origin domain
func (c *Client) CreateCloudOriginDomain(siteID, accountID int, domain, region string, port int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Creating Incapsula cloud origin domain: %s for site: %d\n", domain, siteID)

	payload := CloudOriginDomainCreateRequest{
		Domain: domain,
		Region: region,
		Port:   port,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal cloud origin domain: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s/sites/%d/cloud-origins", c.config.BaseURLRev3, siteID),
		payloadJSON,
		CreateCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while creating cloud origin domain %s for site %d: %s", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula create cloud origin domain JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when creating cloud origin domain %s for site %d: %s", resp.StatusCode, domain, siteID, string(responseBody))
	}

	// Parse the JSON
	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if response.IsError {
		return nil, fmt.Errorf("Error from Incapsula service when creating cloud origin domain %s for site %d: %s", domain, siteID, string(responseBody))
	}

	return &response, nil
}

// GetCloudOriginDomain retrieves a cloud origin domain by ID
func (c *Client) GetCloudOriginDomain(siteID, accountID, originID int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Getting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s/sites/%d/cloud-origins/%d", c.config.BaseURLRev3, siteID, originID),
		nil,
		ReadCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while reading cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula get cloud origin domain JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	// Parse the JSON
	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if response.IsError {
		return nil, fmt.Errorf("Error from Incapsula service when reading cloud origin domain %d for site %d: %s", originID, siteID, string(responseBody))
	}

	return &response, nil
}

// UpdateCloudOriginDomain updates an existing cloud origin domain
func (c *Client) UpdateCloudOriginDomain(siteID, accountID, originID int, region string, port int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Updating Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	payload := CloudOriginDomainUpdateRequest{
		Region: region,
		Port:   port,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal cloud origin domain update: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut,
		fmt.Sprintf("%s/sites/%d/cloud-origins/%d", c.config.BaseURLRev3, siteID, originID),
		payloadJSON,
		UpdateCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while updating cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update cloud origin domain JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	// Parse the JSON
	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if response.IsError {
		return nil, fmt.Errorf("Error from Incapsula service when updating cloud origin domain %d for site %d: %s", originID, siteID, string(responseBody))
	}

	return &response, nil
}

// DeleteCloudOriginDomain deletes a cloud origin domain
func (c *Client) DeleteCloudOriginDomain(siteID, accountID, originID int) error {
	log.Printf("[INFO] Deleting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete,
		fmt.Sprintf("%s/sites/%d/cloud-origins/%d", c.config.BaseURLRev3, siteID, originID),
		nil,
		DeleteCloudOriginDomain)

	if err != nil {
		return fmt.Errorf("Error from Incapsula service while deleting cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete cloud origin domain JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	return nil
}

// MoveCloudOriginDomain moves a cloud origin domain to another site
func (c *Client) MoveCloudOriginDomain(siteID, originID, targetSiteID int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Moving Incapsula cloud origin domain: %d from site: %d to site: %d\n", originID, siteID, targetSiteID)

	payload := map[string]int{"targetSiteId": targetSiteID}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal cloud origin domain move request: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s/sites/%d/cloud-origins/%d/move", c.config.BaseURLRev3, siteID, originID),
		payloadJSON,
		UpdateCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while moving cloud origin domain %d: %s", originID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula move cloud origin domain JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when moving cloud origin domain %d: %s", resp.StatusCode, originID, string(responseBody))
	}

	// Parse the JSON
	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if response.IsError {
		return nil, fmt.Errorf("Error from Incapsula service when moving cloud origin domain %d: %s", originID, string(responseBody))
	}

	return &response, nil
}
