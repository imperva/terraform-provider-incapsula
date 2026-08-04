package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (c *Client) CreateAbpPolicy(accountId string, policy AbpPolicy) (*AbpPolicy, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpPolicyResourceName, accountId)

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpPolicyResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.AbpPolicyCreateUrl(accountId), body, CreateAbpPolicy)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpPolicyResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpPolicyResourceName, accountId, string(responseBody))
	}

	var created AbpPolicy
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpPolicyResourceName, err, string(responseBody))
	}
	return &created, nil
}

// UpdateAbpPolicy updates a policy by id. Returns (nil, nil) if the policy
// no longer exists upstream so callers can clear local state.
func (c *Client) UpdateAbpPolicy(policyId string, policy AbpPolicy) (*AbpPolicy, error) {
	log.Printf("[INFO] Updating %s with id %s", abpPolicyResourceName, policyId)

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpPolicyResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.AbpPolicyUrl(policyId), body, UpdateAbpPolicy)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpPolicyResourceName, policyId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpPolicyResourceName, policyId, string(responseBody))
	}

	var updated AbpPolicy
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpPolicyResourceName, err, string(responseBody))
	}
	return &updated, nil
}

// ReadAbpPolicy fetches a policy by id. If the policy does not exist, returns (nil, nil)
// so callers can distinguish 404 from transport/parse errors and clear local state.
func (c *Client) ReadAbpPolicy(policyId string) (*AbpPolicy, error) {
	log.Printf("[INFO] Reading %s with id %s", abpPolicyResourceName, policyId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.AbpPolicyUrl(policyId), nil, ReadAbpPolicy)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpPolicyResourceName, policyId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpPolicyResourceName, policyId, string(responseBody))
	}

	var policy AbpPolicy
	if err := json.Unmarshal(responseBody, &policy); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpPolicyResourceName, err, string(responseBody))
	}
	return &policy, nil
}

func (c *Client) AbpAccountGlobalPolicyUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/global_policy", c.config.BaseURLAPI, accountId)
}

// ReadAbpAccountGlobalPolicy fetches the account global policy. If the account does
// not exist, returns (nil, nil) so callers can distinguish 404 from other errors.
func (c *Client) ReadAbpAccountGlobalPolicy(accountId string) (*AbpPolicy, error) {
	log.Printf("[INFO] Reading %s global policy for account %s", abpPolicyResourceName, accountId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.AbpAccountGlobalPolicyUrl(accountId), nil, ReadAbpAccountGlobalPolicy)
	if err != nil {
		return nil, fmt.Errorf("error reading %s global policy for account %s: %w", abpPolicyResourceName, accountId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s global policy: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s global policy for account %s: %s", resp.StatusCode, abpPolicyResourceName, accountId, string(responseBody))
	}

	var policy AbpPolicy
	if err := json.Unmarshal(responseBody, &policy); err != nil {
		return nil, fmt.Errorf("error parsing %s global policy read response: %w; body: %s", abpPolicyResourceName, err, string(responseBody))
	}
	return &policy, nil
}

func (c *Client) DeleteAbpPolicy(policyId string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpPolicyResourceName, policyId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.AbpPolicyUrl(policyId), nil, DeleteAbpPolicy)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpPolicyResourceName, policyId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpPolicyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpPolicyResourceName, policyId, string(responseBody))
	}
	return nil
}
