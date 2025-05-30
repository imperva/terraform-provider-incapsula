package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// IncapRule is a struct that encompasses all the properties of an IncapRule
type IncapRule struct {
	Name                  string                `json:"name"`
	Action                string                `json:"action"`
	Filter                string                `json:"filter,omitempty"`
	ResponseCode          int                   `json:"response_code,omitempty"`
	AddMissing            bool                  `json:"add_missing,omitempty"`
	RewriteExisting       *bool                 `json:"rewrite_existing,omitempty"`
	From                  string                `json:"from,omitempty"`
	To                    string                `json:"to,omitempty"`
	RewriteName           string                `json:"rewrite_name,omitempty"`
	DCID                  int                   `json:"dc_id,omitempty"`
	PortForwardingContext string                `json:"port_forwarding_context,omitempty"`
	PortForwardingValue   string                `json:"port_forwarding_value,omitempty"`
	RateContext           string                `json:"rate_context,omitempty"`
	RateInterval          int                   `json:"rate_interval,omitempty"`
	ErrorType             string                `json:"error_type,omitempty"`
	ErrorResponseFormat   string                `json:"error_response_format,omitempty"`
	ErrorResponseData     string                `json:"error_response_data,omitempty"`
	MultipleDeletions     bool                  `json:"multiple_deletions,omitempty"`
	OverrideWafRule       string                `json:"overrideWafRule,omitempty"`
	OverrideWafAction     string                `json:"overrideWafAction,omitempty"`
	Enabled               bool                  `json:"enabled"`
	SendNotifications     *bool                 `json:"sendNotifications,omitempty"`
	BlockDurationDetails  *BlockDurationDetails `json:"blockDurationDetails,omitempty"`
}

type BlockDurationDetails struct {
	BlockDurationType string `json:"blockDurationPeriodType,omitempty"`
	BlockDuration     int    `json:"blockFixedDurationValue,omitempty"`
	BlockDurationMin  int    `json:"blockRandomizedDurationMinValue,omitempty"`
	BlockDurationMax  int    `json:"blockRandomizedDurationMaxValue,omitempty"`
}

// IncapRuleWithID contains the IncapRule as well as the rule identifier
type IncapRuleWithID struct {
	IncapRule
	RuleID int `json:"rule_id"`
}

// AddIncapRule adds an incap rule to be managed by Incapsula
func (c *Client) AddIncapRule(siteID string, rule *IncapRule) (*IncapRuleWithID, error) {
	log.Printf("[INFO] Adding Incapsula Incap Rule for Site ID %s\n", siteID)

	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	log.Printf("[DEBUG] Create rule DTO request: %v\n", string(ruleJSON[:]))

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/rules", c.config.BaseURLRev2, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, ruleJSON, CreateIncapRule)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding Incap Rule for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding Incap Rule for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var incapRuleWithID IncapRuleWithID
	err = json.Unmarshal([]byte(responseBody), &incapRuleWithID)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Rule JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &incapRuleWithID, nil
}

// ReadIncapRule gets the specific Incap Rule
func (c *Client) ReadIncapRule(siteID string, ruleID int) (*IncapRuleWithID, int, error) {
	log.Printf("[INFO] Getting Incapsula Incap Rule %d for Site ID %s\n", ruleID, siteID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadIncapRule)
	if err != nil {
		return nil, 0, fmt.Errorf("Error from Incapsula service when reading Incap Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("Error status code %d from Incapsula service when reading Incap Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	// Parse the JSON
	var incapRuleWithID IncapRuleWithID
	err = json.Unmarshal([]byte(responseBody), &incapRuleWithID)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error parsing Incap Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	return &incapRuleWithID, resp.StatusCode, nil
}

// UpdateIncapRule updates the Incapsula Incap Rule
func (c *Client) UpdateIncapRule(siteID string, ruleID int, rule *IncapRule) (*IncapRuleWithID, error) {
	log.Printf("[INFO] Updating Incapsula Incap Rule %d for Site ID %s\n", ruleID, siteID)

	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	log.Printf("[DEBUG] Update rule DTO request: %v\n", string(ruleJSON[:]))

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, ruleJSON, UpdateIncapRule)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating Incap Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating Incap Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	// Parse the JSON
	var incapRuleWithID IncapRuleWithID
	err = json.Unmarshal([]byte(responseBody), &incapRuleWithID)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	return &incapRuleWithID, nil
}

// DeleteIncapRule deletes a site currently managed by Incapsula
func (c *Client) DeleteIncapRule(siteID string, ruleID int) error {
	log.Printf("[INFO] Deleting Incapsula Incap Rule %d for Site ID %s\n", ruleID, siteID)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteIncapRule)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting Incap Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting Incap Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	return nil
}
