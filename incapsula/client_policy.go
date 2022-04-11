package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// PolicySubmitted is struct that encompasses all the properties of a policy object to submit
type PolicySubmitted struct {
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Enabled             bool            `json:"enabled"`
	AccountID           int             `json:"accountId,omitempty"`
	PolicyType          string          `json:"policyType"`
	PolicySettings      []PolicySetting `json:"policySettings"`
	DefaultPolicyConfig []struct {
		AccountID int    `json:"accountId"`
		AssetType string `json:"assetType"`
		PolicyID  int    `json:"policyId"`
	} `json:"defaultPolicyConfig,omitempty"`
}

// PolicyExtended is a struct that encompasses all the properties of an extended policy setting
type PolicyExtended struct {
	Value struct {
		ID                  int             `json:"id"`
		Name                string          `json:"name"`
		Description         string          `json:"description"`
		Enabled             bool            `json:"enabled"`
		AccountID           int             `json:"accountId,omitempty"`
		PolicyType          string          `json:"policyType"`
		PolicySettings      []PolicySetting `json:"policySettings"`
		DefaultPolicyConfig []struct {
			AccountID int    `json:"accountId"`
			AssetType string `json:"assetType"`
			PolicyID  int    `json:"policyId"`
		} `json:"defaultPolicyConfig"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

// PolicySetting is a struct that encompasses all the properties of a policy setting
type PolicySetting struct {
	SettingsAction    string `json:"settingsAction"`
	PolicySettingType string `json:"policySettingType"`
	Data              struct {
		Geo *struct {
			Countries  []string `json:"countries,omitempty"`
			Continents []string `json:"continents,omitempty"`
		} `json:"geo,omitempty"`
		Ips  []string `json:"ips,omitempty"`
		Urls []struct {
			Pattern string `json:"pattern,omitempty"`
			URL     string `json:"url,omitempty"`
		} `json:"urls,omitempty"`
		HeaderValue string `json:"headerValue,omitempty"`
	} `json:"data"`
	PolicyDataExceptions []struct {
		Data []struct {
			ValidateExceptionData bool     `json:"validateExceptionData,omitempty"`
			ExceptionType         string   `json:"exceptionType,omitempty"`
			Values                []string `json:"values,omitempty"`
		} `json:"data,omitempty"`
		Comment string `json:"comment,omitempty"`
	} `json:"policyDataExceptions,omitempty"`
}

// AddPolicy adds a policy to be managed by Incapsula
func (c *Client) AddPolicy(policySubmitted *PolicySubmitted) (*PolicyExtended, error) {
	log.Printf("[INFO] Adding Incapsula Policy\n")

	policyJSON, err := json.Marshal(policySubmitted)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	// Post form to Incapsula
	log.Printf("[DEBUG] Incapsula Add Incap Policy JSON request: %s\n", string(policyJSON))
	reqURL := fmt.Sprintf("%s/policies/v2/policies", c.config.BaseURLAPI)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, policyJSON, CreatePolicy)
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
	var policyExtended PolicyExtended
	err = json.Unmarshal([]byte(responseBody), &policyExtended)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &policyExtended, nil
}

// GetPolicy gets the policy
func (c *Client) GetPolicy(policyID string) (*PolicyExtended, error) {
	log.Printf("[INFO] Getting Incapsula Policy: %s\n", policyID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/policies/v2/policies/%s?extended=true", c.config.BaseURLAPI, policyID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadPolicy)
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
	var policyExtended PolicyExtended
	err = json.Unmarshal([]byte(responseBody), &policyExtended)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response for Policy ID %s: %s\nresponse: %s", policyID, err, string(responseBody))
	}

	return &policyExtended, nil
}

// UpdatePolicy updates the Incapsula Policy
func (c *Client) UpdatePolicy(policyID int, policySubmitted *PolicySubmitted) (*PolicyExtended, error) {
	log.Printf("[INFO] Updating Incapsula Policy with ID %d\n", policyID)

	policyJSON, err := json.Marshal(policySubmitted)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Policy: %s", err)
	}

	// Post form to Incapsula
	log.Printf("[DEBUG] Incapsula Update Incap Policy JSON request: %s\n", string(policyJSON))
	reqURL := fmt.Sprintf("%s/policies/v2/policies/%d", c.config.BaseURLAPI, policyID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, policyJSON, UpdatePolicy)
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
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating Policy with ID %d: %s", resp.StatusCode, policyID, string(responseBody))
	}

	// Parse the JSON
	var updatedPolicyExtended PolicyExtended
	err = json.Unmarshal([]byte(responseBody), &updatedPolicyExtended)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Policy JSON response for Policy ID %d: %s\nresponse: %s", policyID, err, string(responseBody))
	}

	return &updatedPolicyExtended, nil
}

// DeletePolicy deletes a policy currently managed by Incapsula
func (c *Client) DeletePolicy(policyID string) error {
	log.Printf("[INFO] Deleting Incapsula Policy for ID %s\n", policyID)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s/policies/v2/policies/%s", c.config.BaseURLAPI, policyID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeletePolicy)
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
