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

// Incap Action Enumerations
const actionAlert = "RULE_ACTION_ALERT"
const actionBlockIP = "RULE_ACTION_BLOCK_IP"
const actionBlockRequest = "RULE_ACTION_BLOCK_R" // todo: verify if copy paste error
const actionBlockSession = "RULE_ACTION_BLOCK_S"
const actionDeleteCookie = "RULE_ACTION_DELETE_COOKIE"
const actionDeleteHeader = "RULE_ACTION_DELETE_HEADER"
const actionFwdToDataCenter = "RULE_ACTION_FORWARD_TO_DC"
const actionRedirect = "RULE_ACTION_REDIRECT"
const actionCaptcha = "RULE_ACTION_CAPTCHA"
const actionRetry = "RULE_ACTION_RETRY"
const actionIntrusiveHtml = "RULE_ACTION_INTRUSIVE_HTML"
const actionRewriteCookie = "RULE_ACTION_REWRITE_COOKIE"
const actionRewriteHeader = "RULE_ACTION_REWRITE_HEADER"
const actionRewriteUrl = "RULE_ACTION_REWRITE_URL"

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
func (c *Client) AddIncapRule(siteID, ruleID, dcID int, enabled, priority, name, action, filter, allowCaching, responseCode, from, to, addMissing, rewriteName string) (*IncapRuleAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula incap rule for siteID: %d\n", siteID)

	// Base URL values
	values := url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"enabled": {enabled},
		"name":    {name},
		"action":  {action},
		"filter":  {filter},
	}

	// Additional URL values for specific action types
	switch action {
	case actionAlert:
		fallthrough
	case actionBlockIP:
		fallthrough
	case actionBlockRequest:
		fallthrough
	case actionBlockSession:
		fallthrough
	case actionCaptcha:
		fallthrough
	case actionRetry:
		fallthrough
	case actionIntrusiveHtml:
		values.Add("site_id", strconv.Itoa(siteID))
		values.Add("priority", priority)
	case actionDeleteCookie:
		fallthrough
	case actionDeleteHeader:
		fallthrough
	case actionRewriteUrl:
		values.Add("rule_id", strconv.Itoa(ruleID))
	case actionFwdToDataCenter:
		values.Add("site_id", strconv.Itoa(siteID))
		values.Add("priority", priority)
		values.Add("dc_id", strconv.Itoa(dcID))
		values.Add("allow_caching", allowCaching)
	case actionRedirect:
		values.Add("site_id", strconv.Itoa(siteID))
		values.Add("priority", priority)
		values.Add("response_code", responseCode)
		values.Add("from", from)
		values.Add("to", to)
	case actionRewriteCookie:
		fallthrough
	case actionRewriteHeader:
		values.Add("site_id", strconv.Itoa(siteID))
		values.Add("priority", priority)
		values.Add("add_missing", addMissing)
		values.Add("from", from)
		values.Add("to", to)
		values.Add("allow_caching", allowCaching)
		values.Add("rewrite_name", rewriteName)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleAdd), values)
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
func (c *Client) DeleteIncapRule(ruleID string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type IncapRuleDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula incap rule (rule_id: %s)\n", ruleID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointIncapRuleDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"rule_id": {ruleID},
	})
	if err != nil {
		return fmt.Errorf("Error deleting incap rule (rule_id: %s): %s", ruleID, err)
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
		return fmt.Errorf("Error parsing delete incap rule JSON response (rule_id: %s): %s", ruleID, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting incap rule (rule_id: %s): %s", ruleID, string(responseBody))
	}

	return nil
}
