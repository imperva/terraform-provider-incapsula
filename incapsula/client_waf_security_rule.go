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
const endpointWAFRuleConfigure = "sites/configure/security"

// WAF Rule Enumerations
const backdoorRuleId = "api.threats.backdoor"
const crossSiteScriptingRuleId = "api.threats.cross_site_scripting"
const illegalResourceAccessRuleId = "api.threats.illegal_resource_access"
const remoteFileInclusionRuleId = "api.threats.remote_file_inclusion"
const sqlInjectionRuleId = "api.threats.sql_injection"
const ddosRuleId = "api.threats.ddos"
const botAccessControlRuleId = "api.threats.bot_access_control"

// ConfigureWAFSecurityRule adds an WAF rule
func (c *Client) ConfigureWAFSecurityRule(siteID int, ruleID, security_rule_action, activation_mode, ddos_traffic_threshold, block_bad_bots, challenge_suspected_bots string) (*SiteStatusResponse, error) {

	// Base URL values
	values := url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
		"rule_id": {ruleID},
	}

	// Additional URL values for specific rule ids
	if ruleID == backdoorRuleId || ruleID == crossSiteScriptingRuleId || ruleID == illegalResourceAccessRuleId || ruleID == remoteFileInclusionRuleId || ruleID == sqlInjectionRuleId {
		values.Add("security_rule_action", security_rule_action)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with security_rule_action (%s) for site id (%d)\n", ruleID, security_rule_action, siteID)
	} else if ruleID == ddosRuleId {
		values.Add("activation_mode", activation_mode)
		values.Add("ddos_traffic_threshold", ddos_traffic_threshold)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with activation_mode (%s) and ddos_traffic_threshold (%s) for site id (%d)\n", ruleID, activation_mode, ddos_traffic_threshold, siteID)
	} else if ruleID == botAccessControlRuleId {
		values.Add("block_bad_bots", block_bad_bots)
		values.Add("challenge_suspected_bots", challenge_suspected_bots)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with block_bad_bots (%s) and challenge_suspected_bots (%s) for site id (%d)\n", ruleID, block_bad_bots, challenge_suspected_bots, siteID)
	} else {
		return nil, fmt.Errorf("Error - invalid WAF security rule rule_id (%s)", ruleID)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointWAFRuleConfigure), values)
	if err != nil {
		return nil, fmt.Errorf("Error configuring WAF security rule rule_id (%s) for site_id (%d)", ruleID, siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula configure WAF security rule JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing configure WAF rule JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)
	}

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d): %s", ruleID, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}
