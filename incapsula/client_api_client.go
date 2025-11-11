package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointAPIClient = "/authorization/v3/api-clients"

// APIClientUpdateRequest represents the PATCH request body for updating an API client
// Only the fields that are set will be sent in the PATCH request

type APIClientUpdateRequest struct {
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description,omitempty"`
	ExpirationPeriod *string `json:"expirationDate,omitempty"`
	Enabled          *bool   `json:"enabled,omitempty"`
	GracePeriod      *int    `json:"gracePeriodInSeconds,omitempty"`
	Regenerate       *bool   `json:"regenerate,omitempty"`
}

type APIClientResponse struct {
	APIClientID    string  `json:"id"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	APIKey         string  `json:"apiKey"`
	Enabled        bool    `json:"enabled"`
	ExpirationDate string  `json:"expirationDate"`
	LastUsedAt     string  `json:"lastActionTime"`
	GracePeriod    int     `json:"gracePeriodInSeconds"`
}

// PatchAPIClient updates or regenerates an API client using the unified PATCH endpoint
func (c *Client) PatchAPIClient(accountID int, clientID string, req *APIClientUpdateRequest) (*APIClientResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointAPIClient, clientID)
	log.Printf("[DEBUG] Patch API client URL: %s\n", url)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[DEBUG] Patch API client Body: %s\n", body)

	params := GetRequestParamsWithCaid(accountID)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPatch, url, body, params, UpdateApiClient)
	if err != nil {
		return nil, fmt.Errorf("Error updating api_client with Id %s: %s", clientID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update api_client JSON response: %s", string(responseBody))

	// Parse the JSON
	var apiClientResponse APIClientResponse
	err = json.Unmarshal([]byte(responseBody), &apiClientResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update api_client JSON response for id %s: %s", clientID, err)
	}

	// Look at the response status code from Incapsula
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating api_client %s: %s", resp.StatusCode, clientID, string(responseBody))
	}

	return &apiClientResponse, nil
}

func (c *Client) GetAPIClient(accountID int, userEmail string, clientID string) (*APIClientResponse, error) {
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointAPIClient)
	log.Printf("[DEBUG] **** URL: %s", string(reqURL))
	params := GetRequestParamsWithCaid(accountID)
	params["id"] = clientID
	if userEmail != "" {
		params["userEmail"] = userEmail
	}

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadApiClient)

	if err != nil {
		return nil, fmt.Errorf("Error getting api_client with id %s: %s", clientID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula api_client status JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when getting api_client %s: %s", resp.StatusCode, clientID, string(responseBody))
	}

	// Parse the JSON
	var apiClientResponse APIClientResponse
	err = json.Unmarshal(responseBody, &apiClientResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing api_client status JSON response for api_client id %s: %s", clientID, err)
	}

	log.Printf("[INFO] ResponseStruct : %+v\n", apiClientResponse)
	return &apiClientResponse, nil

}

func (c *Client) CreateAPIClient(accountID int, userEmail string) (*APIClientResponse, error) {
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointAPIClient)
	log.Printf("[DEBUG] **** URL: %s\n", string(reqURL))

	params := GetRequestParamsWithCaid(accountID)
	if userEmail != "" {
		params["userEmail"] = userEmail
	}

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, nil, params, CreateApiClient)

	if err != nil {
		return nil, fmt.Errorf("Error creating api_client", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula api_client status JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when creating api_client: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var apiClientResponse APIClientResponse
	err = json.Unmarshal(responseBody, &apiClientResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing api_client status JSON response for api_client: %s", err)
	}

	log.Printf("[INFO] ResponseStruct : %+v\n", apiClientResponse)
	return &apiClientResponse, nil

}

func (c *Client) DeleteAPIClient(accountID int, clientID string) error {

	log.Printf("[INFO] Deleting api client with ID: %s", clientID)

	requestUrl := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointAPIClient, clientID)
	log.Printf("[DEBUG] Deleting api client URL: %s\n", string(requestUrl))

	params := GetRequestParamsWithCaid(accountID)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, requestUrl, nil, params, DeleteApiClient)

	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting api-client: %s %s", clientID, err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting api-client %s:%s", resp.StatusCode, resp.Body, clientID)
	}

	return nil
}
