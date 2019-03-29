package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

// Endpoints (unexported consts)
const endpointIncapRuleAdd = "sites/incapRules/add"
const endpointIncapRuleList = "sites/incapRules/list"
const endpointIncapRuleEdit = "sites/incapRules/edit"
const endpointIncapRuleDelete = "sites/incapRules/delete"

// todo: get incap rule responses
// IncapRuleAddResponse contains todo
type IncapRuleAddResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// IncapRuleListResponse contains todo
type IncapRuleListResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// IncapRuleEditResponse contains todo
type IncapRuleEditResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// IncapRuleDeleteResponse contains todo
type IncapRuleDeleteResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// AddIncapRule adds an incap rule to be managed by Incapsula
func (c *Client) AddIncapRule(siteID int, enabled, priority, name, action, filter string, ruleID int) (*IncapRuleAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula incap rule for siteID: %d\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleAdd), url.Values{
		"api_id":   {c.config.APIID},
		"api_key":  {c.config.APIKey},
		"site_id":  {strconv.Itoa(siteID)},
		"enabled":  {enabled},
		"priority": {priority},
		"name":     {name},
		"action":   {action},
		"filter":   {filter},
		"rule_id":  {strconv.Itoa(ruleID)},
	})
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding incap rule for siteID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add incap rule JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var incapRuleAddResponse IncapRuleAddResponse
	err = json.Unmarshal([]byte(responseBody), &incapRuleAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add incap rule JSON response for siteID %d: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding incap rule for siteID %d: %s", siteID, string(responseBody))
	}

	return &incapRuleAddResponse, nil
}

// IncapRuleList gets the Incapsula list of incap rules
func (c *Client) ListIncapRules(includeAdRules, includeIncapRules string) (*IncapRuleListResponse, error) {
	log.Printf("[INFO] Getting Incapsula incaprules (include_ad_rules: %s, include_incap_rules: %s)\n", includeAdRules, includeIncapRules)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleList), url.Values{
		"api_id":              {c.config.APIID},
		"api_key":             {c.config.APIKey},
		"include_ad_rules":    {includeAdRules},
		"include_incap_rules": {includeIncapRules},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting incap rules (include_ad_rules: %s, include_incap_rules: %s): %s", includeAdRules, includeIncapRules, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula incap rules JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var incapRuleListResponse IncapRuleListResponse
	err = json.Unmarshal([]byte(responseBody), &incapRuleListResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing incap rule list JSON response (include_ad_rules: %s, include_incap_rules: %s): %s", includeAdRules, includeIncapRules, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleListResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting incap rule list (include_ad_rules: %s, include_incap_rules: %s): %s", includeAdRules, includeIncapRules, string(responseBody))
	}

	return &incapRuleListResponse, nil
}

// EditIncapRule edits the Incapsula incap rule
func (c *Client) EditIncapRule(siteID int, enabled, priority, name, action, filter string, ruleID int) (*IncapRuleEditResponse, error) {
	log.Printf("[INFO] Editing Incapsula incap rule name: %s for siteID: %d\n", name, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleEdit), url.Values{
		"api_id":   {c.config.APIID},
		"api_key":  {c.config.APIKey},
		"site_id":  {strconv.Itoa(siteID)},
		"enabled":  {enabled},
		"priority": {priority},
		"name":     {name},
		"action":   {action},
		"filter":   {filter},
		"rule_id":  {strconv.Itoa(ruleID)},
	})
	if err != nil {
		return nil, fmt.Errorf("Error editing incap rule name: %s for siteID: %d: %s", name, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula edit incap rule JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var incapRuleEditResponse IncapRuleEditResponse
	err = json.Unmarshal([]byte(responseBody), &incapRuleEditResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing edit incap rule JSON response for siteID %d: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleEditResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when editing incap rule for siteID %d, ruleID: %d: %s", siteID, ruleID, string(responseBody))
	}

	return &incapRuleEditResponse, nil
}

// DeleteIncapRule deletes a site currently managed by Incapsula
func (c *Client) DeleteIncapRule(ruleID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type IncapRuleDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula incap rule (rule_id: %d)\n", ruleID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"rule_id": {strconv.Itoa(ruleID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting incap rule (rule_id: %d): %s", ruleID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete incap rule JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var incapRuleDeleteResponse IncapRuleDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &incapRuleDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete incap rule JSON response (rule_id: %d): %s", ruleID, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting incap rule (rule_id: %d): %s", ruleID, string(responseBody))
	}

	return nil
}
