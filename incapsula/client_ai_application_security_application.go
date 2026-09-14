package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// The application controller path dropped its /api/ segment (AIFW-1432): the base path moved
// from /v3/api/applications to /v3/applications, matching the policy and api-key controllers.
const aiApplicationSecurityApplicationEndpoint = "/ai-application-security/v3/applications"

// AiApplicationSecurityImpervaApiBody wraps every AI Application Security write payload / response body
// (backend ImpervaApiBody<T> envelope, AIFW-1321). Shared by the policy and api-key resources.
type AiApplicationSecurityImpervaApiBody struct {
	Data json.RawMessage `json:"data"`
}

// AiApplicationSecurityApplicationRequest is the write payload for Create (POST) and Update (PATCH).
type AiApplicationSecurityApplicationRequest struct {
	Name            string                                  `json:"name,omitempty"`
	ApplicationType string                                  `json:"applicationType,omitempty"`
	Region          string                                  `json:"region,omitempty"`
	Configuration   *AiApplicationSecurityApplicationConfig `json:"configuration,omitempty"`
}

// AiApplicationSecurityApplicationConfig matches ApplicationConfigurationDto.
type AiApplicationSecurityApplicationConfig struct {
	SiteId                   int64                                   `json:"siteId,omitempty"`
	Path                     string                                  `json:"path,omitempty"`
	ContentType              string                                  `json:"contentType,omitempty"`
	PromptLocation           string                                  `json:"promptLocation,omitempty"`
	BlockedResponseStructure string                                  `json:"blockedResponseStructure,omitempty"`
	IsStreaming              bool                                    `json:"isStreaming"`
	Request                  *AiApplicationSecurityStreamingRequest  `json:"request,omitempty"`
	Response                 *AiApplicationSecurityStreamingResponse `json:"response,omitempty"`
}

// AiApplicationSecurityStreamingRequest matches StreamingRequestConfigDto.
type AiApplicationSecurityStreamingRequest struct {
	MessagePath string `json:"messagePath,omitempty"`
	ContentPath string `json:"contentPath,omitempty"`
	RolePath    string `json:"rolePath,omitempty"`
}

// AiApplicationSecurityStreamingResponse matches StreamingResponseConfigDto.
type AiApplicationSecurityStreamingResponse struct {
	RolePath          string `json:"rolePath,omitempty"`
	ContentPath       string `json:"contentPath,omitempty"`
	EndOfStreamMarker string `json:"endOfStreamMarker,omitempty"`
	FinishReasonPath  string `json:"finishReasonPath,omitempty"`
	FinishReasonValue string `json:"finishReasonValue,omitempty"`
}

// AiApplicationSecurityApplicationDetails matches ApplicationDetailsDto (flattened response shape).
type AiApplicationSecurityApplicationDetails struct {
	ApplicationId     string                                  `json:"applicationId"`
	Name              string                                  `json:"name"`
	AccountId         int64                                   `json:"accountId"`
	Region            string                                  `json:"region"`
	Status            string                                  `json:"status"`
	StatusDescription string                                  `json:"statusDescription"`
	ApplicationType   string                                  `json:"applicationType"`
	Configuration     *AiApplicationSecurityApplicationConfig `json:"configuration,omitempty"`
}

type aiApplicationSecurityApplicationListResponse struct {
	Data []AiApplicationSecurityApplicationDetails `json:"data"`
}

type aiApplicationSecurityApplicationResponse struct {
	Data AiApplicationSecurityApplicationDetails `json:"data"`
}

func aiApplicationSecurityWrapData(payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}
	return json.Marshal(AiApplicationSecurityImpervaApiBody{Data: data})
}

// CreateAiApplicationSecurityApplication creates a new AI Application Security application.
// POST /ai-application-security/v3/applications?caid={accountID}
func (c *Client) CreateAiApplicationSecurityApplication(accountID int, req *AiApplicationSecurityApplicationRequest) (*AiApplicationSecurityApplicationDetails, error) {
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, aiApplicationSecurityApplicationEndpoint)
	params := GetRequestParamsWithCaid(accountID)

	body, err := aiApplicationSecurityWrapData(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Create AI Application Security application URL: %s, params: %s, body: %s", reqURL, params, string(body))

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, body, params, CreateAiApplicationSecurityApplication)
	if err != nil {
		return nil, fmt.Errorf("Error creating AI Application Security application: %s", err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading create AI Application Security application response: %s", err)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when creating application: %s", resp.StatusCode, string(responseBody))
	}

	var parsed aiApplicationSecurityApplicationResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing create AI Application Security application JSON response: %s", err)
	}

	return &parsed.Data, nil
}

// GetAiApplicationSecurityApplication reads an AI Application Security application by ID.
// GET /ai-application-security/v3/applications?caid={accountID}&applicationId={applicationID}
// Returns (nil, nil) if the application does not exist.
func (c *Client) GetAiApplicationSecurityApplication(accountID int, applicationID string) (*AiApplicationSecurityApplicationDetails, error) {
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, aiApplicationSecurityApplicationEndpoint)
	params := GetRequestParamsWithCaid(accountID)
	params["applicationId"] = applicationID

	log.Printf("[DEBUG] Read AI Application Security application URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadAiApplicationSecurityApplication)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security application with id %s: %s", applicationID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security application response body for id %s: %s", applicationID, err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when reading application %s: %s", resp.StatusCode, applicationID, string(responseBody))
	}

	var parsed aiApplicationSecurityApplicationListResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing read AI Application Security application JSON response for id %s: %s", applicationID, err)
	}

	if len(parsed.Data) == 0 {
		return nil, nil
	}

	return &parsed.Data[0], nil
}

// UpdateAiApplicationSecurityApplication partially updates an AI Application Security application.
// PATCH /ai-application-security/v3/applications/{applicationID}?caid={accountID}
func (c *Client) UpdateAiApplicationSecurityApplication(accountID int, applicationID string, req *AiApplicationSecurityApplicationRequest) (*AiApplicationSecurityApplicationDetails, error) {
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, aiApplicationSecurityApplicationEndpoint, applicationID)
	params := GetRequestParamsWithCaid(accountID)

	body, err := aiApplicationSecurityWrapData(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Update AI Application Security application URL: %s, params: %s, body: %s", reqURL, params, string(body))

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPatch, reqURL, body, params, UpdateAiApplicationSecurityApplication)
	if err != nil {
		return nil, fmt.Errorf("Error updating AI Application Security application with id %s: %s", applicationID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading update AI Application Security application response body for id %s: %s", applicationID, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when updating application %s: %s", resp.StatusCode, applicationID, string(responseBody))
	}

	var parsed aiApplicationSecurityApplicationResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing update AI Application Security application JSON response for id %s: %s", applicationID, err)
	}

	return &parsed.Data, nil
}

// DeleteAiApplicationSecurityApplication deletes an AI Application Security application.
// DELETE /ai-application-security/v3/applications/{applicationID}?caid={accountID}
func (c *Client) DeleteAiApplicationSecurityApplication(accountID int, applicationID string) error {
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, aiApplicationSecurityApplicationEndpoint, applicationID)
	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Delete AI Application Security application URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteAiApplicationSecurityApplication)
	if err != nil {
		return fmt.Errorf("Error deleting AI Application Security application with id %s: %s", applicationID, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error status code %d from AI Application Security service when deleting application %s: %s", resp.StatusCode, applicationID, string(responseBody))
	}

	return nil
}
