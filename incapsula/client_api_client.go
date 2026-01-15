package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointAPIClient = "/authorization/v3/api-clients"

type APIClientUpdateRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
	Regenerate     bool   `json:"regenerate,omitempty"`
}

type APIClientResponse struct {
	APIClientID    int    `json:"id"`
	UserEmail      string `json:"userEmail"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	APIKey         string `json:"apiKey"`
	Enabled        bool   `json:"enabled"`
	ExpirationDate string `json:"expirationDate"`
	LastUsedAt     string `json:"lastActionTime"`
}

type APIClientResponseTemp struct {
	Meta MetaData            `json:"meta"`
	Data []APIClientResponse `json:"data"`
}

type MetaData struct {
	Total          int `json:"total"`
	PageNumber     int `json:"pageNumber"`
	PageSize       int `json:"pageSize"`
	MaxApiKeyLimit int `json:"maxApiKeyLimit"`
}

// PatchAPIClient updates or regenerates an API client using the unified PATCH endpoint
func (c *Client) PatchAPIClient(accountID int, clientID string, req *APIClientUpdateRequest) (*APIClientResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointAPIClient, clientID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Patch API client URL: %s, Request:%+v, Params: %s, Body: %s\n", url, req, params, body)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPatch, url, body, params, UpdateApiClient)
	if err != nil {
		return nil, fmt.Errorf("Error updating api_client with Id %s: %s", clientID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

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
	log.Printf("[DEBUG]  Create API Response : %+v", apiClientResponse)
	return &apiClientResponse, nil
}

func (c *Client) GetAPIClient(accountID int, clientID string) (*APIClientResponse, error) {

	log.Printf("[DEBUG] Reading incapsula api_client with account_id:%d, client_id:%s", accountID, clientID)

	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointAPIClient)
	params := GetRequestParamsWithCaid(accountID)
	params["id"] = clientID

	log.Printf("[DEBUG] GET URL: %s, params: %s", reqURL, params)

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
	var apiClientResponseTemp APIClientResponseTemp
	err = json.Unmarshal(responseBody, &apiClientResponseTemp)
	if err != nil {
		return nil, fmt.Errorf("Error parsing api_client status JSON response for api_client id %s: %s", clientID, err)
	}
	log.Printf("[INFO] GET Response temp Struct : %+v", apiClientResponseTemp)

	if len(apiClientResponseTemp.Data) == 0 {
		return nil, nil
	}

	return &apiClientResponseTemp.Data[0], nil

}

func (c *Client) CreateAPIClient(accountID int, userEmail string, req *APIClientUpdateRequest) (*APIClientResponse, error) {
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointAPIClient)

	params := GetRequestParamsWithCaid(accountID)
	if userEmail != "" {
		params["userEmail"] = userEmail
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[DEBUG] Create API Client URL: %s, params: %s, body:%s", reqURL, params, string(body))

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, body, params, CreateApiClient)

	if err != nil {
		return nil, fmt.Errorf("Error creating api_client: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when creating api_client: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var apiClientResponse APIClientResponse
	err = json.Unmarshal(responseBody, &apiClientResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing api_client status JSON response for api_client: %s", err)
	}

	log.Printf("[DEBUG]  Create API Response : %+v", apiClientResponse)
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

func (resp APIClientResponse) String() string {
	return fmt.Sprintf("ApiClientResponse{APIClientID:%d, APIKey:<SECRET>, Name:%s, Description:%s, Enabled:%v, ExpirationDate:%s, LastUsedAt:%s}", resp.APIClientID, resp.Name, resp.Description, resp.Enabled, resp.ExpirationDate, resp.LastUsedAt)
}

func Bool(v bool) *bool { return &v }
