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
const endpointExceptionConfigure = "sites/configure/whitelists"
const endpointExceptionList = "sites/status"

// Exception param mapping by ruleID
// NOTE: no exceptions for whitelistedIPsExceptionRuleId
var securityRuleExceptionParamMapping = map[string][]string{
	// ACL RuleIDs
	blacklistedCountriesExceptionRuleID: {"client_app_types", "ips", "url_patterns", "urls"},
	blacklistedIPsExceptionRuleID:       {"client_apps", "countries", "ips", "url_patterns", "urls"},
	blacklistedURLsExceptionRuleID:      {"client_apps", "countries", "ips", "url_patterns", "urls"},
	// WAF RuleIDs
	backdoorExceptionRuleID:              {"client_apps", "countries", "ips", "url_patterns", "urls", "user_agents", "parameters"},
	botAccessControlExceptionRuleID:      {"client_app_types", "ips", "url_patterns", "urls", "user_agents"},
	crossSiteScriptingExceptionRuleID:    {"client_apps", "countries", "url_patterns", "urls", "parameters"},
	ddosExceptionRuleID:                  {"client_apps", "countries", "ips", "url_patterns", "urls"},
	illegalResourceAccessExceptionRuleID: {"client_apps", "countries", "ips", "url_patterns", "urls", "parameters"},
	remoteFileInclusionExceptionRuleID:   {"client_apps", "countries", "ips", "url_patterns", "urls", "user_agents", "parameters"},
	sqlInjectionExceptionRuleID:          {"client_apps", "countries", "ips", "url_patterns", "urls"},
}

// SecurityRuleExceptionCreateResponse provides exception_id of rule exception
type SecurityRuleExceptionCreateResponse struct {
	Res         string `json:"res"`
	ExceptionID string `json:"exception_id"`
	Status      string `json:"status"`
}

// AddSecurityRuleException adds a security rule exception
func (c *Client) AddSecurityRuleException(siteID int, ruleID, clientAppTypes, clientApps, countries, continents, ips, urlPatterns, urls, userAgents, parameters string) (*SecurityRuleExceptionCreateResponse, error) {
	// Base URL values
	values := url.Values{
		"api_id":            {c.config.APIID},
		"api_key":           {c.config.APIKey},
		"site_id":           {strconv.Itoa(siteID)},
		"rule_id":           {ruleID},
		"exception_id_only": {"true"},
	}

	log.Printf("[INFO] Adding new security rule exception for rule_id (%s) for site id (%d)\n", ruleID, siteID)

	// Check to see if ruleID is correct, then iterate rule specific parameters
	if ruleParams, ok := securityRuleExceptionParamMapping[ruleID]; ok {
		for i := 0; i < len(ruleParams); i++ {
			// Add param values for specific ruleID based on securityRuleExceptionParamMapping
			param := ruleParams[i]
			if param == "client_app_types" && clientAppTypes != "" {
				values.Add("client_app_types", clientAppTypes)
			} else if param == "client_apps" && clientApps != "" {
				values.Add("client_apps", clientApps)
			} else if param == "countries" && countries != "" {
				values.Add("countries", countries)
			} else if param == "continents" && continents != "" {
				values.Add("continents", continents)
			} else if param == "ips" && ips != "" {
				values.Add("ips", ips)
			} else if param == "parameters" && parameters != "" {
				values.Add("parameters", parameters)
			} else if param == "url_patterns" && urlPatterns != "" {
				values.Add("url_patterns", urlPatterns)
			} else if param == "urls" && urls != "" {
				values.Add("urls", urls)
			} else if param == "user_agents" && userAgents != "" {
				values.Add("user_agents", userAgents)
			}
		}
	} else {
		return nil, fmt.Errorf("Error configuring security rule exception: invalid rule_id (%s)", ruleID)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointExceptionConfigure), values)
	if err != nil {
		return nil, fmt.Errorf("Error configuring security rule exception rule_id (%s) for site_id (%d)", ruleID, siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula configure SecurityRuleExceptionCreateResponse JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var securityRuleExceptionCreateResponse SecurityRuleExceptionCreateResponse
	err = json.Unmarshal([]byte(responseBody), &securityRuleExceptionCreateResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing SecurityRuleExceptionCreateResponse JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)
	}

	// Look at the response status code from Incapsula
	if securityRuleExceptionCreateResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when adding security rule exception for rule_id (%s) and site_id (%d): %s", ruleID, siteID, string(responseBody))
	}

	return &securityRuleExceptionCreateResponse, nil
}

// EditSecurityRuleException edits a security rule exception
func (c *Client) EditSecurityRuleException(siteID int, ruleID, clientAppTypes, clientApps, countries, continents, ips, urlPatterns, urls, userAgents, parameters, whitelistID string) (*SiteStatusResponse, error) {
	// Base URL values
	values := url.Values{
		"api_id":       {c.config.APIID},
		"api_key":      {c.config.APIKey},
		"site_id":      {strconv.Itoa(siteID)},
		"rule_id":      {ruleID},
		"whitelist_id": {whitelistID},
	}

	log.Printf("[INFO] Updating existing security rule exception for rule_id (%s) whitelist_id (%s) for site_id (%d)\n", ruleID, whitelistID, siteID)

	// Check to see if ruleID is correct, then iterate rule specific parameters
	if ruleParams, ok := securityRuleExceptionParamMapping[ruleID]; ok {
		for i := 0; i < len(ruleParams); i++ {
			// Add param values for specific ruleID based on securityRuleExceptionParamMapping
			param := ruleParams[i]
			if param == "client_app_types" && clientAppTypes != "" {
				values.Add("client_app_types", clientAppTypes)
			} else if param == "client_apps" && clientApps != "" {
				values.Add("client_apps", clientApps)
			} else if param == "countries" && countries != "" {
				values.Add("countries", countries)
			} else if param == "continents" && continents != "" {
				values.Add("continents", continents)
			} else if param == "ips" && ips != "" {
				values.Add("ips", ips)
			} else if param == "parameters" && parameters != "" {
				values.Add("parameters", parameters)
			} else if param == "url_patterns" && urlPatterns != "" {
				values.Add("url_patterns", urlPatterns)
			} else if param == "urls" && urls != "" {
				values.Add("urls", urls)
			} else if param == "user_agents" && userAgents != "" {
				values.Add("user_agents", userAgents)
			}
		}
	} else {
		return nil, fmt.Errorf("Error configuring security rule exception: invalid rule_id (%s)", ruleID)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointExceptionConfigure), values)
	if err != nil {
		return nil, fmt.Errorf("Error configuring security rule exception rule_id (%s) for site_id (%d)", ruleID, siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula configure security rule exception JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing configure security rule exception JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)
	}

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding security rule exception for rule_id (%s) and site_id (%d): %s", ruleID, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}

// ListSecurityRuleExceptions gets the site status including the list of exceptions for security rules
func (c *Client) ListSecurityRuleExceptions(siteID, ruleID string) (*SiteStatusResponse, error) {
	log.Printf("[INFO] Getting Incapsula security rule exeptions for rule_id (%s) on site_id (%s)\n", ruleID, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointExceptionList), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {siteID},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting security rule exceptions for rule_id (%s) on siteID (%s): %s", ruleID, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula ListSecurityRuleExceptions JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing ListSecurityRuleExceptions JSON response for siteID: %s %s\nresponse: %s", siteID, err, string(responseBody))
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
		return &siteStatusResponse, fmt.Errorf("Error from Incapsula service when getting security rule exceptions (site_id: %s): %s", siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}

// DeleteSecurityRuleException deletes a security rule exception
func (c *Client) DeleteSecurityRuleException(siteID int, ruleID, whitelistID string) error {
	type ExceptionDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
		Status     string `json:"status"`
	}

	// Base URL values
	values := url.Values{
		"api_id":           {c.config.APIID},
		"api_key":          {c.config.APIKey},
		"site_id":          {strconv.Itoa(siteID)},
		"rule_id":          {ruleID},
		"whitelist_id":     {whitelistID},
		"delete_whitelist": {"true"},
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointExceptionConfigure), values)
	if err != nil {
		return fmt.Errorf("Error deleting security rule exception whitelist_id (%s) for rule_id (%s) for site_id (%d)", whitelistID, ruleID, siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Deleting Incapsula security rule exception JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var exceptionDeleteResponse ExceptionDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &exceptionDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete security rule exception JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)
	}

	// Look at the response status code from Incapsula
	if exceptionDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting security rule exception for rule_id (%s) and site_id (%d): %s", ruleID, siteID, string(responseBody))
	}

	return nil
}
