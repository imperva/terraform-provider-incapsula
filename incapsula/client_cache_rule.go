package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// CacheRule is a struct that encompasses all the properties of a CacheRule
type CacheRule struct {
	Name                 string `json:"name"`
	Action               string `json:"action"`
	Filter               string `json:"filter"`
	Enabled              bool   `json:"enabled"`
	TTL                  int    `json:"ttl"`
	IgnoredParams        string `json:"ignored_params,omitempty"`
	Text                 string `json:"text,omitempty"`
	DifferentiateByValue string `json:"differentiate_by_value,omitempty"`
}

// CacheRuleWithID contains the CacheRule as well as the rule identifier
type CacheRuleWithID struct {
	CacheRule
	RuleID int `json:"rule_id"`
}

// AddCacheRule adds an incap rule to be managed by Incapsula
func (c *Client) AddCacheRule(siteID string, rule *CacheRule) (*CacheRuleWithID, error) {
	log.Printf("[INFO] Adding Incapsula Cache Rule for Site ID %s\n", siteID)

	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal CacheRule: %s", err)
	}

	// Dump Request JSON
	log.Printf("[DEBUG] Incapsula Add Cache Rule JSON request body: %s\n", string(ruleJSON))

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache/rules", c.config.BaseURLRev2, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, ruleJSON, CreateCacheRule)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding Cache Rule for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add Cache Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding Cache Rule for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var cacheRuleWithID CacheRuleWithID
	err = json.Unmarshal([]byte(responseBody), &cacheRuleWithID)
	if err != nil || !strings.Contains(string(responseBody), "\"rule_id\":") {
		return nil, fmt.Errorf("Error parsing Cache Rule JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &cacheRuleWithID, nil
}

// ReadCacheRule gets the specific Incap Rule
func (c *Client) ReadCacheRule(siteID string, ruleID int) (*CacheRuleWithID, int, error) {
	log.Printf("[INFO] Getting Incapsula Cache Rule %d for Site ID %s\n", ruleID, siteID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadCacheRule)
	if err != nil {
		return nil, 0, fmt.Errorf("Error from Incapsula service when reading Cache Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Cache Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("Error status code %d from Incapsula service when reading Cache Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	// Parse the JSON
	var cacheRuleWithID CacheRuleWithID
	err = json.Unmarshal([]byte(responseBody), &cacheRuleWithID)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error parsing Cache Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	return &cacheRuleWithID, resp.StatusCode, nil
}

// UpdateCacheRule updates the Incapsula Incap Rule
func (c *Client) UpdateCacheRule(siteID string, ruleID int, rule *CacheRule) error {
	log.Printf("[INFO] Updating Incapsula Cache Rule %d for Site ID %s\n", ruleID, siteID)

	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return fmt.Errorf("Failed to JSON marshal CacheRule: %s", err)
	}

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, ruleJSON, UpdateCacheRule)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when updating Cache Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Cache Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when updating Cache Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	// Parse the JSON
	var cacheRuleWithID CacheRuleWithID
	err = json.Unmarshal([]byte(responseBody), &cacheRuleWithID)
	if err != nil {
		return fmt.Errorf("Error parsing Cache Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	return nil
}

// DeleteCacheRule deletes a site currently managed by Incapsula
func (c *Client) DeleteCacheRule(siteID string, ruleID int) error {
	type DeleteCacheRuleResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
		DebugInfo  struct {
			RuleID string `json:"rule_id"`
			IDInfo string `json:"id-info"`
		} `json:"debug_info"`
	}

	log.Printf("[INFO] Deleting Incapsula Cache Rule %d for Site ID %s\n", ruleID, siteID)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/settings/cache/rules/%d", c.config.BaseURLRev2, siteID, ruleID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteCacheRule)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting Cache Rule %d for Site ID %s: %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Cache Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	// Unfortunately, this API endpoint is not RESTful and we return 200's back for failures (instead of 40X - joy)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting Cache Rule %d for Site ID %s: %s", resp.StatusCode, ruleID, siteID, string(responseBody))
	}

	// Parse the JSON
	var deleteCacheRuleResponse DeleteCacheRuleResponse
	err = json.Unmarshal([]byte(responseBody), &deleteCacheRuleResponse)
	if err != nil {
		return fmt.Errorf("Error parsing Delete Cache Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	if deleteCacheRuleResponse.Res != 0 {
		return fmt.Errorf("Error deleting Cache Rule %d JSON response for Site ID %s: %s\nresponse: %s", ruleID, siteID, err, string(responseBody))
	}

	return nil
}
