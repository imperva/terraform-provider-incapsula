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
	blacklistedCountriesExceptionRuleId: []string{"client_app_types", "ips", "url_patterns", "urls"},
	blacklistedIPsExceptionRuleId:       []string{"client_apps", "countries", "ips", "url_patterns", "urls"},
	blacklistedURLsExceptionRuleId:      []string{"client_apps", "countries", "ips", "url_patterns", "urls"},
	// WAF RuleIDs
	backdoorExceptionRuleId:              []string{"client_apps", "countries", "ips", "url_patterns", "urls", "user_agents", "parameters"},
	botAccessControlExceptionRuleId:      []string{"client_app_types", "ips", "url_patterns", "urls", "user_agents"},
	crossSiteScriptingExceptionRuleId:    []string{"client_apps", "countries", "url_patterns", "urls", "parameters"},
	ddosExceptionRuleId:                  []string{"client_apps", "countries", "ips", "url_patterns", "urls"},
	illegalResourceAccessExceptionRuleId: []string{"client_apps", "countries", "ips", "url_patterns", "urls", "parameters"},
	remoteFileInclusionExceptionRuleId:   []string{"client_apps", "countries", "ips", "url_patterns", "urls", "user_agents", "parameters"},
	sqlInjectionExceptionRuleId:          []string{"client_apps", "countries", "ips", "url_patterns", "urls"},
}

// SecurityRuleExceptionCreateResponse provides exception_id of rule exception
type SecurityRuleExceptionCreateResponse struct {
	Res         string `json:"res"`
	ExceptionID string `json:"exception_id"`
	Status      string `json:"status"`
}

// AddSecurityRuleException adds an exception to a security rule
func (c *Client) AddSecurityRuleException(siteID int, ruleID, client_app_types, client_apps, countries, continents, ips, url_patterns, urls, user_agents, parameters string) (*SecurityRuleExceptionCreateResponse, error) {

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
			switch param {
			case "client_app_types":
				if client_app_types != "" {
					values.Add("client_app_types", client_app_types)
				}
			case "client_apps":
				if client_apps != "" {
					values.Add("client_apps", client_apps)
				}
			case "countries":
				if countries != "" {
					values.Add("countries", countries)
				}
			case "continents":
				if continents != "" {
					values.Add("continents", continents)
				}
			case "ips":
				if ips != "" {
					values.Add("ips", ips)
				}
			case "parameters":
				if parameters != "" {
					values.Add("parameters", parameters)
				}
			case "url_patterns":
				if url_patterns != "" {
					values.Add("url_patterns", url_patterns)
				}
			case "urls":
				if urls != "" {
					values.Add("urls", urls)
				}
			case "user_agents":
				if user_agents != "" {
					values.Add("user_agents", user_agents)
				}
			}
		}
	} else {
		return nil, fmt.Errorf("Error - invalid security rule exception rule_id (%s)", ruleID)
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

// EditSecurityRuleException adds an exception to a security rule
func (c *Client) EditSecurityRuleException(siteID int, ruleID, client_app_types, client_apps, countries, continents, ips, url_patterns, urls, user_agents, parameters, whitelistID string) (*SiteStatusResponse, error) {
	//whitelistID, _ := strconv.Atoi(whitelist_id)
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
			switch param {
			case "client_app_types":
				if client_app_types != "" {
					values.Add("client_app_types", client_app_types)
				}
			case "client_apps":
				if client_apps != "" {
					values.Add("client_apps", client_apps)
				}
			case "countries":
				if countries != "" {
					values.Add("countries", countries)
				}
			case "continents":
				if continents != "" {
					values.Add("continents", continents)
				}
			case "ips":
				if ips != "" {
					values.Add("ips", ips)
				}
			case "parameters":
				if parameters != "" {
					values.Add("parameters", parameters)
				}
			case "url_patterns":
				if url_patterns != "" {
					values.Add("url_patterns", url_patterns)
				}
			case "urls":
				if urls != "" {
					values.Add("urls", urls)
				}
			case "user_agents":
				if user_agents != "" {
					values.Add("user_agents", user_agents)
				}
			}
		}
	} else {
		return nil, fmt.Errorf("Error - invalid security rule exception rule_id (%s)", ruleID)
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

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting data centers list (site_id: %s): %s", siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}

func (c *Client) DeleteSecurityRuleException(siteID int, ruleID, whitelistID string) error {
	type ExceptionDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
		Status     string `json:"status"`
	}

	//whitelistID, _ := strconv.Atoi(whitelist_id)
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
