package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// aiApplicationSecurityGuardrailExceptionEndpoint is the policy-scoped collection path of the management
// service's GuardianExceptionController. Format it with the application and policy UUIDs; append "/{id}" for items.
const aiApplicationSecurityGuardrailExceptionEndpoint = "/ai-application-security/v3/applications/%s/policies/%s/guardian-exceptions"

// AiApplicationSecurityGuardrailExceptionRequest is the create (POST) and update (PUT) body
// (CreateGuardianExceptionRequestDto / UpdateGuardianExceptionRequestDto). Every field is sent in full: the update
// replaces them, and an empty note clears it.
type AiApplicationSecurityGuardrailExceptionRequest struct {
	GuardrailId      int64  `json:"guardianId"`
	ExceptionPattern string `json:"exceptionPattern"`
	ExceptionReason  string `json:"exceptionReason"`
	Note             string `json:"note"`
	Active           bool   `json:"active"`
}

// AiApplicationSecurityGuardrailException matches GuardianExceptionResponseDto (the fields the resource uses).
// The service trims leading and trailing whitespace of exceptionPattern and note when it saves them.
type AiApplicationSecurityGuardrailException struct {
	Id               string  `json:"id"`
	GuardrailId      int64   `json:"guardianId"`
	ExceptionPattern string  `json:"exceptionPattern"`
	ExceptionReason  string  `json:"exceptionReason"`
	Note             *string `json:"note"`
	Active           bool    `json:"active"`
	Author           string  `json:"author"`
	UpdatedBy        *string `json:"updatedBy"`
}

type aiApplicationSecurityGuardrailExceptionResponseBody struct {
	Data AiApplicationSecurityGuardrailException `json:"data"`
}

func (c *Client) aiApplicationSecurityGuardrailExceptionURL(applicationID string, policyID string, guardrailExceptionID string) string {
	url := c.config.BaseURLAPI + fmt.Sprintf(aiApplicationSecurityGuardrailExceptionEndpoint, applicationID, policyID)
	if guardrailExceptionID != "" {
		url += "/" + guardrailExceptionID
	}
	return url
}

// doAiApplicationSecurityGuardrailException sends one request and parses the {"data": {...}} response.
// Returns (nil, nil) on 404 when notFoundIsNil is set.
func (c *Client) doAiApplicationSecurityGuardrailException(method string, reqURL string, accountID int, payload interface{}, operation string,
	okStatuses []int, notFoundIsNil bool) (*AiApplicationSecurityGuardrailException, error) {
	var body []byte
	if payload != nil {
		var err error
		if body, err = aiApplicationSecurityWrapData(payload); err != nil {
			return nil, err
		}
	}
	params := GetRequestParamsWithCaid(accountID)
	log.Printf("[DEBUG] %s AI Application Security guardrail exception URL: %s, params: %s", method, reqURL, params)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(method, reqURL, body, params, operation)
	if err != nil {
		return nil, fmt.Errorf("Error calling %s %s: %s", method, reqURL, err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading the response of %s %s: %s", method, reqURL, err)
	}

	if resp.StatusCode == http.StatusNotFound && notFoundIsNil {
		return nil, nil
	}
	ok := false
	for _, status := range okStatuses {
		ok = ok || resp.StatusCode == status
	}
	if !ok {
		return nil, fmt.Errorf("Error status code %d from AI Application Security service on %s %s: %s", resp.StatusCode, method, reqURL, string(responseBody))
	}
	if method == http.MethodDelete {
		return nil, nil
	}

	var parsed aiApplicationSecurityGuardrailExceptionResponseBody
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("Error parsing the response of %s %s: %s", method, reqURL, err)
	}
	return &parsed.Data, nil
}

// CreateAiApplicationSecurityGuardrailException creates a guardrail exception on a guardrail of the policy.
// POST .../v3/applications/{applicationID}/policies/{policyID}/guardian-exceptions?caid={accountID}
func (c *Client) CreateAiApplicationSecurityGuardrailException(accountID int, applicationID string, policyID string,
	req *AiApplicationSecurityGuardrailExceptionRequest) (*AiApplicationSecurityGuardrailException, error) {
	return c.doAiApplicationSecurityGuardrailException(http.MethodPost, c.aiApplicationSecurityGuardrailExceptionURL(applicationID, policyID, ""),
		accountID, req, CreateAiApplicationSecurityGuardrailException, []int{http.StatusOK, http.StatusCreated}, false)
}

// GetAiApplicationSecurityGuardrailException reads one guardrail exception of the policy.
// GET .../guardian-exceptions/{id}. Returns (nil, nil) on 404: it does not exist, or not under this policy.
func (c *Client) GetAiApplicationSecurityGuardrailException(accountID int, applicationID string, policyID string,
	guardrailExceptionID string) (*AiApplicationSecurityGuardrailException, error) {
	return c.doAiApplicationSecurityGuardrailException(http.MethodGet, c.aiApplicationSecurityGuardrailExceptionURL(applicationID, policyID, guardrailExceptionID),
		accountID, nil, ReadAiApplicationSecurityGuardrailException, []int{http.StatusOK}, true)
}

// UpdateAiApplicationSecurityGuardrailException edits a guardrail exception in place - the id stays the same,
// also when it moves to another guardrail of the policy. PUT .../guardian-exceptions/{id}
func (c *Client) UpdateAiApplicationSecurityGuardrailException(accountID int, applicationID string, policyID string, guardrailExceptionID string,
	req *AiApplicationSecurityGuardrailExceptionRequest) (*AiApplicationSecurityGuardrailException, error) {
	return c.doAiApplicationSecurityGuardrailException(http.MethodPut, c.aiApplicationSecurityGuardrailExceptionURL(applicationID, policyID, guardrailExceptionID),
		accountID, req, UpdateAiApplicationSecurityGuardrailException, []int{http.StatusOK}, false)
}

// DeleteAiApplicationSecurityGuardrailException deletes a guardrail exception. DELETE .../guardian-exceptions/{id}
// A 404 counts as success: it is already gone (deleted in the UI, or with its policy).
func (c *Client) DeleteAiApplicationSecurityGuardrailException(accountID int, applicationID string, policyID string, guardrailExceptionID string) error {
	_, err := c.doAiApplicationSecurityGuardrailException(http.MethodDelete, c.aiApplicationSecurityGuardrailExceptionURL(applicationID, policyID, guardrailExceptionID),
		accountID, nil, DeleteAiApplicationSecurityGuardrailException, []int{http.StatusOK, http.StatusNoContent}, true)
	return err
}
