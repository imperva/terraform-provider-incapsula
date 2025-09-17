package incapsula

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const endpointAPIClient = "/v3/api-client/"

// APIClientUpdateRequest represents the PATCH request body for updating an API client
// Only the fields that are set will be sent in the PATCH request

type APIClientUpdateRequest struct {
	ExpirationPeriod *string `json:"expiration_period,omitempty"`
	Enabled          *bool   `json:"enabled,omitempty"`
	GracePeriod      *int    `json:"grace_period,omitempty"`
	Regenerate       *bool   `json:"regenerate,omitempty"`
}

type APIClientResponse struct {
	APIClientID    string `json:"api_client_id"`
	APIKey         string `json:"api_key"`
	Enabled        bool   `json:"enabled"`
	ExpirationDate string `json:"expiration_date"`
	LastUsedAt     string `json:"last_used_at"`
	GracePeriod    int    `json:"grace_period"`
}

// PatchAPIClient updates or regenerates an API client using the unified PATCH endpoint
func (c *Client) PatchAPIClient(ctx context.Context, clientID string, req *APIClientUpdateRequest) (*APIClientResponse, error) {
	url := fmt.Sprintf("%s%s%s", c.config.BaseURL, endpointAPIClient, clientID)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create PATCH request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	c.setAuthHeaders(request)
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("PATCH request failed: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PATCH response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("PATCH /v3/api-client/%s failed: %s, body: %s", clientID, resp.Status, string(responseBody))
	}
	var apiResp APIClientResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PATCH response: %w", err)
	}
	return &apiResp, nil
}

// GetAPIClient fetches API client metadata (GET /v3/api-client/{client_id})
func (c *Client) GetAPIClient(ctx context.Context, clientID string) (*APIClientResponse, error) {
	url := fmt.Sprintf("%s%s%s", c.config.BaseURL, endpointAPIClient, clientID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	c.setAuthHeaders(request)
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GET response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GET /v3/api-client/%s failed: %s, body: %s", clientID, resp.Status, string(responseBody))
	}
	var apiResp APIClientResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GET response: %w", err)
	}
	return &apiResp, nil
}

// DeleteAPIClient deletes an API client (DELETE /v3/api-client/{client_id})
func (c *Client) DeleteAPIClient(ctx context.Context, clientID string) error {
	url := fmt.Sprintf("%s%s%s", c.config.BaseURL, endpointAPIClient, clientID)
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	c.setAuthHeaders(request)
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read DELETE response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("DELETE /v3/api-client/%s failed: %s, body: %s", clientID, resp.Status, string(responseBody))
	}
	return nil
}

// setAuthHeaders sets authentication headers for the request
func (c *Client) setAuthHeaders(req *http.Request) {
	// Example: req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	// Implement according to your authentication scheme
}
