package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type BotStruct struct {
	ID          *int   `json:"id"`
	DisplayName string `json:"displayName"`
}

type BotsStruct struct {
	CanceledGoodBots []BotStruct `json:"canceledGoodBots"`
	BadBots          []BotStruct `json:"badBots"`
}

// BotsConfigurationDTO - Same DTO for: GET response, POST request, and POST response
type BotsConfigurationDTO struct {
	Errors []ApiError   `json:"errors"`
	Data   []BotsStruct `json:"data"`
}

// UpdateBotAccessControlConfiguration - Update the Bot Access Control configuration for a given website
func (c *Client) UpdateBotAccessControlConfiguration(siteID string, requestDTO BotsConfigurationDTO) (*BotsConfigurationDTO, error) {
	log.Printf("[INFO] Updating the Bot Access Control configuration for siteID: %s\n", siteID)

	log.Printf("[INFO]  requestDTO: %+v\n", requestDTO)

	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	botsJSON, err := json.Marshal(requestDTO)
	log.Printf("[INFO]  botsJSON: %v\n", string(botsJSON))
	reqURL := fmt.Sprintf("%s/sites/%s/settings/botConfiguration", baseURLv3, siteID)
	log.Printf("[INFO]  reqURL: %v\n", reqURL)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, botsJSON, CreateBotConfiguration)
	if err != nil {
		return nil, fmt.Errorf("Error executing update Bot Access Control configuration request for siteID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Bot Access Control configuration JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO BotsConfigurationDTO
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Update Bot Access Control configuration JSON response for siteID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &responseDTO, nil
}

// GetBotAccessControlConfiguration - Retrieve the Bot Access Control configuration for a given website
func (c *Client) GetBotAccessControlConfiguration(siteID string) (*BotsConfigurationDTO, error) {
	log.Printf("[INFO] Getting Bot Access Control configuration (site_id: %s)\n", siteID)

	// Get request to Incapsula
	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	reqURL := fmt.Sprintf("%s/sites/%s/settings/botConfiguration", baseURLv3, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadBotConfiguration)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Bot Access Control configuration request for siteID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Bot Access Control JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO BotsConfigurationDTO
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Bot Access Control list JSON response for siteID: %s %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &responseDTO, nil
}

// GetClientApplicationsMetadata - Retrieve the Bot Access Control Metadata
func (c *Client) GetClientApplicationsMetadata() (*ClientApps, error) {
	log.Printf("[INFO] Getting Client Applications Metadata\n")

	// Get request to Incapsula
	baseURLIntegration := strings.Replace(c.config.BaseURL, "/prov/", "/integration/", 1)
	reqURL := fmt.Sprintf("%s/clapps", baseURLIntegration)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, ReadClientApplications)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Bot Access Control Metadata request %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Client Applications Metadata JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO ClientApps
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Client Applications Metadata list JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &responseDTO, nil
}
