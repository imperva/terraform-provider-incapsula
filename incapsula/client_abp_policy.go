package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func (c *Client) AbpPolicyUrl(policyId string) string {
	return fmt.Sprintf("%s/v1/policy/%s", c.config.BaseURLAPI, policyId)
}

func (c *Client) AbpPolicyCreateUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/policy", c.config.BaseURLAPI, accountId)
}

const abpPolicyResourceName = "ABP Policy"

func (c *Client) ListAbpPolicies(accountId string) ([]AbpPolicy, error) {
	log.Printf("[INFO] Listing %s in ABP account %s", abpPolicyResourceName, accountId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.AbpPolicyCreateUrl(accountId), nil, ListAbpPolicies)
	if err != nil {
		return nil, fmt.Errorf("error listing %ss in ABP account %s: %w", abpPolicyResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when listing %ss: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when listing %ss in ABP account %s: %s", resp.StatusCode, abpPolicyResourceName, accountId, string(responseBody))
	}

	var listResp struct {
		Items []AbpPolicy `json:"items"`
	}
	if err := json.Unmarshal(responseBody, &listResp); err != nil {
		return nil, fmt.Errorf("error parsing %s list response: %w; body: %s", abpPolicyResourceName, err, string(responseBody))
	}
	return listResp.Items, nil
}

const createAbpPolicyAction = "creating abp policy"

func (c *Client) CreateAbpPolicy(accountId string, policy AbpPolicy) (*AbpPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Creating Abp Policy for Account ID %s\n", accountId)

	policyJson, err := json.Marshal(policy)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure generating abp_policy create request"),
			Detail:   fmt.Sprintf("Failed to JSON marshal abp_policy: %s", err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] abp_policy payload: %s\n", string(policyJson))

	reqURL := c.AbpPolicyCreateUrl(accountId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, policyJson, CreateAbpPolicy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPolicyAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] Incapsula Create abp_policy JSON response: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusCreated {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPolicyAction, nil, responseBody))
		return nil, diags
	}

	var newPolicy AbpPolicy
	err = json.Unmarshal(responseBody, &newPolicy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	return &newPolicy, diags
}

const updateAbpPolicyAction = "updating abp_policy"

func (c *Client) UpdateAbpPolicy(policyId string, policy AbpPolicy) (*AbpPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Updating Abp Policy ID %s\n", policyId)

	policyJson, err := json.Marshal(policy)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure generating abp_policy update request"),
			Detail:   fmt.Sprintf("Failed to JSON marshal AbpPolicy: %s", err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] abp_policy_update payload: %s\n", string(policyJson))

	reqURL := c.AbpPolicyUrl(policyId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, policyJson, UpdateAbpPolicy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(updateAbpPolicyAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(updateAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] Incapsula Update abp_policy JSON response: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		diags = append(diags, httpSourcedErrorDiagnostic(updateAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	var updated AbpPolicy
	err = json.Unmarshal(responseBody, &updated)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(updateAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	return &updated, diags
}

const readAbpPolicyAction = "reading abp_policy"

// ReadAbpPolicy fetches a policy by id. If the policy does not exist, returns (nil, nil)
// so callers can distinguish 404 from transport/parse errors and clear local state.
func (c *Client) ReadAbpPolicy(policyId string) (*AbpPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Reading Abp Policy ID %s\n", policyId)

	reqURL := c.AbpPolicyUrl(policyId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAbpPolicy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(readAbpPolicyAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(readAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] Incapsula GET abp_policy JSON response: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		diags = append(diags, httpSourcedErrorDiagnostic(readAbpPolicyAction, nil, responseBody))
		return nil, diags
	}

	var policy AbpPolicy
	err = json.Unmarshal(responseBody, &policy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(readAbpPolicyAction, &err, responseBody))
		return nil, diags
	}

	return &policy, diags
}

const deleteAbpPolicyAction = "deleting abp_policy"

func (c *Client) DeleteAbpPolicy(policyId string) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[INFO] Deleting Abp Policy ID %s\n", policyId)

	reqURL := c.AbpPolicyUrl(policyId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteAbpPolicy)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(deleteAbpPolicyAction, &err, nil))
		return diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(deleteAbpPolicyAction, &err, responseBody))
		return diags
	}

	log.Printf("[DEBUG] Incapsula DELETE abp_policy response: %s\n", string(responseBody))

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[INFO] ABP policy ID %s already gone upstream", policyId)
		return diags
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		diags = append(diags, httpSourcedErrorDiagnostic(deleteAbpPolicyAction, nil, responseBody))
		return diags
	}

	return diags
}
