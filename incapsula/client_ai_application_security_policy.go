package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// aiApplicationSecurityPolicyEndpoint is the application-scoped policy collection path (PolicyController,
// @RequestMapping("/v3/applications/{applicationId}/policies")). Format it with the
// application UUID to get the collection URL; append "/{policyId}" for item operations.
const aiApplicationSecurityPolicyEndpoint = "/ai-application-security/v3/applications/%s/policies"

// aiApplicationSecurityPolicyImportApplicationPlaceholder is used as the {applicationId} path segment
// when the real application id is unknown (policyId-only import). The backend
// getPolicyById(accountId, policyId) ignores the path segment and looks the policy up by id,
// so any syntactically valid UUID works; the response then carries the real applicationId.
const aiApplicationSecurityPolicyImportApplicationPlaceholder = "00000000-0000-0000-0000-000000000000"

// AiApplicationSecurityGuardrail matches GuardrailRequestDto (write) and GuardrailDto (read) — they share
// the guardrailType/guardrailMode/guardrailPhase/config/active JSON field names. Config carries
// the "type" discriminator required by the backend's polymorphic BaseGuardrailConfigDto
// (@JsonTypeInfo(property="type")); the resource layer injects it on write and strips it on read.
type AiApplicationSecurityGuardrail struct {
	GuardrailType  string          `json:"guardrailType"`
	GuardrailMode  string          `json:"guardrailMode"`
	GuardrailPhase string          `json:"guardrailPhase"`
	Active         bool            `json:"active"`
	Config         json.RawMessage `json:"config,omitempty"`
}

// AiApplicationSecurityPolicyRequest is the write payload for Create (POST) and Update (PATCH). It matches
// CreatePolicyRequestDto and PatchPolicyRequestDto: because Terraform always sends the full
// desired state, a single PATCH covers field-only, guardrail-only, and combined changes.
type AiApplicationSecurityPolicyRequest struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Active      bool                             `json:"active"`
	Request     []AiApplicationSecurityGuardrail `json:"request"`
	Response    []AiApplicationSecurityGuardrail `json:"response"`
}

// AiApplicationSecurityPolicyResponse matches PolicyDto (the subset the resource consumes).
type AiApplicationSecurityPolicyResponse struct {
	Id            string                           `json:"id"`
	AccountId     int64                            `json:"accountId"`
	ApplicationId string                           `json:"applicationId"`
	Name          string                           `json:"name"`
	Description   string                           `json:"description"`
	Active        bool                             `json:"active"`
	Request       []AiApplicationSecurityGuardrail `json:"request"`
	Response      []AiApplicationSecurityGuardrail `json:"response"`
}

type aiApplicationSecurityPolicyResponseBody struct {
	Data AiApplicationSecurityPolicyResponse `json:"data"`
}

// aiApplicationSecurityPolicyCollectionURL builds the application-scoped policy collection URL.
func (c *Client) aiApplicationSecurityPolicyCollectionURL(applicationID string) string {
	return c.config.BaseURLAPI + fmt.Sprintf(aiApplicationSecurityPolicyEndpoint, applicationID)
}

// CreateAiApplicationSecurityPolicy creates the (single) policy for an application.
// POST /ai-application-security/v3/applications/{applicationID}/policies?caid={accountID}
func (c *Client) CreateAiApplicationSecurityPolicy(accountID int, applicationID string, req *AiApplicationSecurityPolicyRequest) (*AiApplicationSecurityPolicyResponse, error) {
	reqURL := c.aiApplicationSecurityPolicyCollectionURL(applicationID)
	params := GetRequestParamsWithCaid(accountID)

	body, err := aiApplicationSecurityWrapData(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Create AI Application Security policy URL: %s, params: %s, body: %s", reqURL, params, string(body))

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, body, params, CreateAiApplicationSecurityPolicy)
	if err != nil {
		return nil, fmt.Errorf("Error creating AI Application Security policy: %s", err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading create AI Application Security policy response: %s", err)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when creating policy for application %s: %s", resp.StatusCode, applicationID, string(responseBody))
	}

	var parsed aiApplicationSecurityPolicyResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing create AI Application Security policy JSON response: %s", err)
	}

	return &parsed.Data, nil
}

// GetAiApplicationSecurityPolicy reads a policy by id.
// GET /ai-application-security/v3/applications/{applicationID}/policies/{policyID}?caid={accountID}
// Returns (nil, nil) if the policy does not exist.
func (c *Client) GetAiApplicationSecurityPolicy(accountID int, applicationID string, policyID string) (*AiApplicationSecurityPolicyResponse, error) {
	reqURL := fmt.Sprintf("%s/%s", c.aiApplicationSecurityPolicyCollectionURL(applicationID), policyID)
	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Read AI Application Security policy URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadAiApplicationSecurityPolicy)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security policy with id %s: %s", policyID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading AI Application Security policy response body for id %s: %s", policyID, err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when reading policy %s: %s", resp.StatusCode, policyID, string(responseBody))
	}

	var parsed aiApplicationSecurityPolicyResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing read AI Application Security policy JSON response for id %s: %s", policyID, err)
	}

	return &parsed.Data, nil
}

// UpdateAiApplicationSecurityPolicy updates a policy with the full desired state via PATCH.
// PATCH /ai-application-security/v3/applications/{applicationID}/policies/{policyID}?caid={accountID}
func (c *Client) UpdateAiApplicationSecurityPolicy(accountID int, applicationID string, policyID string, req *AiApplicationSecurityPolicyRequest) (*AiApplicationSecurityPolicyResponse, error) {
	reqURL := fmt.Sprintf("%s/%s", c.aiApplicationSecurityPolicyCollectionURL(applicationID), policyID)
	params := GetRequestParamsWithCaid(accountID)

	body, err := aiApplicationSecurityWrapData(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Update AI Application Security policy URL: %s, params: %s, body: %s", reqURL, params, string(body))

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPatch, reqURL, body, params, UpdateAiApplicationSecurityPolicy)
	if err != nil {
		return nil, fmt.Errorf("Error updating AI Application Security policy with id %s: %s", policyID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading update AI Application Security policy response body for id %s: %s", policyID, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service when updating policy %s: %s", resp.StatusCode, policyID, string(responseBody))
	}

	var parsed aiApplicationSecurityPolicyResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing update AI Application Security policy JSON response for id %s: %s", policyID, err)
	}

	return &parsed.Data, nil
}

// DeleteAiApplicationSecurityPolicy deletes a policy by id.
// DELETE /ai-application-security/v3/applications/{applicationID}/policies/{policyID}?caid={accountID}
func (c *Client) DeleteAiApplicationSecurityPolicy(accountID int, applicationID string, policyID string) error {
	reqURL := fmt.Sprintf("%s/%s", c.aiApplicationSecurityPolicyCollectionURL(applicationID), policyID)
	params := GetRequestParamsWithCaid(accountID)

	log.Printf("[DEBUG] Delete AI Application Security policy URL: %s, params: %s", reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteAiApplicationSecurityPolicy)
	if err != nil {
		return fmt.Errorf("Error deleting AI Application Security policy with id %s: %s", policyID, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error status code %d from AI Application Security service when deleting policy %s: %s", resp.StatusCode, policyID, string(responseBody))
	}

	return nil
}
