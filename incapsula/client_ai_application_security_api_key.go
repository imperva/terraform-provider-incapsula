package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// API key endpoints (relative to c.config.BaseURLAPI). The api-key, policy and application
// controllers all share the /v3/applications base path (the application controller dropped its
// former /api/ segment in AIFW-1432).
const (
	// aiApplicationSecurityApiKeyAppEndpoint is the application-scoped path for create and delete.
	aiApplicationSecurityApiKeyAppEndpoint = "/ai-application-security/v3/applications/%s/api-keys"
	// aiApplicationSecurityApiKeyAccountEndpoint is the account-level list path used by Read (and import).
	aiApplicationSecurityApiKeyAccountEndpoint = "/ai-application-security/v3/api-keys"
)

// AiApplicationSecurityApiKeyRequest is the write payload for Create. Wrapped by aiApplicationSecurityWrapData
// into the {"data": {...}} envelope. Only name is settable (CreateApiKeyRequestDto).
type AiApplicationSecurityApiKeyRequest struct {
	Name string `json:"name"`
}

// AiApplicationSecurityApiKeyInfo matches ApiKeyDto. The plaintext key is never present here — it is
// returned once, separately, in the create response's fullKey field.
type AiApplicationSecurityApiKeyInfo struct {
	Id            int64  `json:"id"`
	HashedApiKey  string `json:"hashedApiKey,omitempty"`
	MaskedApiKey  string `json:"maskedApiKey,omitempty"`
	AccountId     int64  `json:"accountId,omitempty"`
	Name          string `json:"name"`
	CreatedAt     int64  `json:"createdAt,omitempty"`
	Active        bool   `json:"active"`
	LastUsedAt    int64  `json:"lastUsedAt,omitempty"`
	ApplicationId string `json:"applicationId,omitempty"` // UUID
}

// AiApplicationSecurityApiKeyCreateResponse matches CreateApiKeyResponseDto: the created key's
// metadata plus the plaintext fullKey (returned exactly once, unrecoverable afterward).
type AiApplicationSecurityApiKeyCreateResponse struct {
	ApiKey  AiApplicationSecurityApiKeyInfo `json:"apiKey"`
	FullKey string                          `json:"fullKey"`
}

// AiApplicationSecurityApiKeyListResponse matches ApiKeyListResponseDto (account-level list).
type AiApplicationSecurityApiKeyListResponse struct {
	ApiKeys    []AiApplicationSecurityApiKeyInfo `json:"apiKeys"`
	TotalCount int64                             `json:"totalCount"`
}

type aiApplicationSecurityApiKeyCreateResponseBody struct {
	Data AiApplicationSecurityApiKeyCreateResponse `json:"data"`
}

type aiApplicationSecurityApiKeyListResponseBody struct {
	Data AiApplicationSecurityApiKeyListResponse `json:"data"`
}

// CreateAiApplicationSecurityApiKey creates an API key for an application.
// POST /ai-application-security/v3/applications/{applicationID}/api-keys?caid={accountID}
// A 400 (e.g. the max-5-keys-per-account limit) is surfaced as an actionable error.
func (c *Client) CreateAiApplicationSecurityApiKey(accountID int, applicationID string, req *AiApplicationSecurityApiKeyRequest) (*AiApplicationSecurityApiKeyCreateResponse, error) {
	reqURL := c.config.BaseURLAPI + fmt.Sprintf(aiApplicationSecurityApiKeyAppEndpoint, applicationID)
	params := GetRequestParamsWithCaid(accountID)

	body, err := aiApplicationSecurityWrapData(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Create AI Application Security API key URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, body, params, CreateAiApplicationSecurityApiKey)
	if err != nil {
		return nil, fmt.Errorf("Error creating AI Application Security API key for application %s: %s", applicationID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading create AI Application Security API key response: %s", err)
	}

	if resp.StatusCode == 400 {
		return nil, fmt.Errorf("Error creating AI Application Security API key for application %s: the backend rejected the request (HTTP 400). An account may hold at most 5 API keys; delete an existing key before creating another. Backend response: %s", applicationID, string(responseBody))
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when creating API key for application %s: %s", resp.StatusCode, applicationID, string(responseBody))
	}

	var parsed aiApplicationSecurityApiKeyCreateResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing create AI Application Security API key JSON response: %s", err)
	}

	return &parsed.Data, nil
}

// GetAiApplicationSecurityApiKey reads an API key by its numeric id from the account-level list.
// GET /ai-application-security/v3/api-keys?caid={accountID}
// Returns (nil, nil) if no key with the given id exists (or the list is empty / 404).
func (c *Client) GetAiApplicationSecurityApiKey(accountID int, apiKeyID int64) (*AiApplicationSecurityApiKeyInfo, error) {
	reqURL := c.config.BaseURLAPI + aiApplicationSecurityApiKeyAccountEndpoint
	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Read AI Application Security API key URL: %s, params: %s, id: %d", reqURL, params, apiKeyID)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadAiApplicationSecurityApiKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security API key with id %d: %s", apiKeyID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security API key response body for id %d: %s", apiKeyID, err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when reading API key %d: %s", resp.StatusCode, apiKeyID, string(responseBody))
	}

	var parsed aiApplicationSecurityApiKeyListResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing read AI Application Security API key JSON response for id %d: %s", apiKeyID, err)
	}

	for i := range parsed.Data.ApiKeys {
		if parsed.Data.ApiKeys[i].Id == apiKeyID {
			return &parsed.Data.ApiKeys[i], nil
		}
	}

	return nil, nil
}

// DeleteAiApplicationSecurityApiKey deletes an API key.
// DELETE /ai-application-security/v3/applications/{applicationID}/api-keys/{apiKeyID}?caid={accountID}
func (c *Client) DeleteAiApplicationSecurityApiKey(accountID int, applicationID string, apiKeyID int64) error {
	reqURL := fmt.Sprintf("%s/%d", c.config.BaseURLAPI+fmt.Sprintf(aiApplicationSecurityApiKeyAppEndpoint, applicationID), apiKeyID)
	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Delete AI Application Security API key URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteAiApplicationSecurityApiKey)
	if err != nil {
		return fmt.Errorf("Error deleting AI Application Security API key with id %d: %s", apiKeyID, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error status code %d from AI Application Security service when deleting API key %d: %s", resp.StatusCode, apiKeyID, string(responseBody))
	}

	return nil
}
