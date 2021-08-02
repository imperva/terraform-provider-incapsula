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
const backdoorRuleID = "api.threats.backdoor"
const crossSiteScriptingRuleID = "api.threats.cross_site_scripting"
const illegalResourceAccessRuleID = "api.threats.illegal_resource_access"
const remoteFileInclusionRuleID = "api.threats.remote_file_inclusion"
const sqlInjectionRuleID = "api.threats.sql_injection"
const ddosRuleID = "api.threats.ddos"
const botAccessControlRuleID = "api.threats.bot_access_control"
const customRuleDefaultActionID = "api.threats.customRule"

// ConfigureWAFSecurityRule adds an WAF rule
func (c *Client) ConfigureWAFSecurityRule(siteID int, ruleID, securityRuleAction, activationMode, ddosTrafficThreshold, blockBadBots, challengeSuspectedBots string) (*SiteStatusResponse, error) {
	// Base URL values
	values := url.Values{
		"site_id": {strconv.Itoa(siteID)},
		"rule_id": {ruleID},
	}

	// Additional URL values for specific rule ids
	if ruleID == backdoorRuleID || ruleID == crossSiteScriptingRuleID || ruleID == illegalResourceAccessRuleID || ruleID == remoteFileInclusionRuleID || ruleID == sqlInjectionRuleID {
		values.Add("security_rule_action", securityRuleAction)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with security rule action (%s) for site id (%d)\n", ruleID, securityRuleAction, siteID)
	} else if ruleID == ddosRuleID {
		values.Add("activation_mode", activationMode)
		values.Add("ddos_traffic_threshold", ddosTrafficThreshold)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with activation mode (%s) and DDoS traffic threshold (%s) for site id (%d)\n", ruleID, activationMode, ddosTrafficThreshold, siteID)
	} else if ruleID == botAccessControlRuleID {
		values.Add("block_bad_bots", blockBadBots)
		values.Add("challenge_suspected_bots", challengeSuspectedBots)
		log.Printf("[INFO] Configuring Incapsula WAF rule id (%s) with block_bad_bots (%s) and challenge suspected bots (%s) for site id (%d)\n", ruleID, blockBadBots, challengeSuspectedBots, siteID)
	} else {
		return nil, fmt.Errorf("Error - invalid WAF security rule rule_id (%s)", ruleID)
	}

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointWAFRuleConfigure)
	resp, err := c.PostFormWithHeaders(reqURL, values)
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

	// Res can sometimes oscillate between a string and number
	// We need to add safeguards for this inside the provider
	var resString string

	if resNumber, ok := siteStatusResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = siteStatusResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d): %s", ruleID, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}
