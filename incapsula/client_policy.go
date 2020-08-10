package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Policy is a struct that encompasses all the properties for a policy
type Policy struct {
	Value struct {
		ID                int    `json:"id"`
		PolicyType        string `json:"policyType"`
		Name              string `json:"name"`
		Enabled           bool   `json:"enabled"`
		Description       string `json:"description"`
		AccountID         int    `json:"accountId"`
		LastModified      string `json:"lastModified"`
		LastModifiedBy    int    `json:"lastModifiedBy"`
		LastUserModified  string `json:"lastUserModified"`
		NumberOfAssets    int    `json:"numberOfAssets"`
		IsMarkedAsDefault bool   `json:"isMarkedAsDefault"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

// PolicyLite is a struct that encompasses all the properties for a policy to be used in the creation process
type PolicyLite struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Enabled        bool   `json:"enabled"`
	AccountID      int    `json:"accountId"`
	PolicyType     string `json:"policyType"`
	PolicySettings []int  `json:"policySettings"`
}

// AddPolicy adds a policy to be managed by Incapsula
func (c *Client) AddPolicy(policyLite *PolicyLite) (*Policy, error) {
	log.Printf("[INFO] Adding Incapsula Policy\n")

	policyJSON, err := json.Marshal(policyLite)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/policies/v2/policies?api_id=%s&api_key=%s", c.config.BaseURLAPI, c.config.APIID, c.config.APIKey),
		"application/json",
		bytes.NewReader(policyJSON))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding Policy: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add Policy JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding Policy: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var policy Policy
	err = json.Unmarshal([]byte(responseBody), &policy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &policy, nil
}

// GetPolicy gets the policy
func (c *Client) GetPolicy(policyID string) (*Policy, error) {
	log.Printf("[INFO] Getting Incapsula Policy: %s\n", policyID)

	// Post form to Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/policies/v2/policies/%s?extended=true&api_id=%s&api_key=%s", c.config.BaseURLAPI, policyID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading Policy for ID %s: %s", policyID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Policy JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Policy for ID %s: %s", resp.StatusCode, policyID, string(responseBody))
	}

	// Parse the JSON
	var policy Policy
	err = json.Unmarshal([]byte(responseBody), &policy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response for Policy ID %s: %s\nresponse: %s", policyID, err, string(responseBody))
	}

	return &policy, nil
}

// UpdatePolicy updates the Incapsula Policy
func (c *Client) UpdatePolicy(policy *PolicyLite) (*Policy, error) {
	log.Printf("[INFO] Updating Incapsula Policy with ID %d\n", policy.ID)

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Policy: %s", err)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/policies/v2/policies/%d?api_id=%s&api_key=%s", c.config.BaseURLAPI, policy.ID, c.config.APIID, c.config.APIKey),
		"application/json",
		bytes.NewReader(policyJSON))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating Policy: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Policy JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating Policy with ID %d: %s", resp.StatusCode, policy.ID, string(responseBody))
	}

	// Parse the JSON
	var updatedPolicy Policy
	err = json.Unmarshal([]byte(responseBody), &policy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response for Policy ID %d: %s\nresponse: %s", policy.ID, err, string(responseBody))
	}

	return &updatedPolicy, nil
}

// DeletePolicy deletes a policy currently managed by Incapsula
func (c *Client) DeletePolicy(policyID string) error {
	log.Printf("[INFO] Deleting Incapsula Policy for ID %s\n", policyID)

	// Delete request to Incapsula
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/policies/v2/policies/%s?api_id=%s&api_key=%s", c.config.BaseURLAPI, policyID, c.config.APIID, c.config.APIKey),
		nil)
	if err != nil {
		return fmt.Errorf("Error preparing HTTP DELETE for deleting Policy with ID %s: %s", policyID, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting Policy with ID %s: %s", policyID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Policy JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting Policy with ID %s: %s", resp.StatusCode, policyID, string(responseBody))
	}

	return nil
}
