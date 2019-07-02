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
const endpointACLRuleConfigure = "sites/configure/acl"

// ACL Rule Enumerations
const blacklistedCountries = "api.acl.blacklisted_countries"
const blacklistedURLs = "api.acl.blacklisted_urls"
const blacklistedIPs = "api.acl.blacklisted_ips"
const whitelistedIPs = "api.acl.whitelisted_ips"

// ConfigureACLSecurityRule adds an ACL rule
func (c *Client) ConfigureACLSecurityRule(siteID int, ruleID, continents, countries, ips, urls, urlPatterns string) (*SiteStatusResponse, error) {
	log.Printf("[INFO] Configuring Incapsula ACL rule id: %s for site id: %d\n", ruleID, siteID)

	// Base URL values
	values := url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
		"rule_id": {ruleID},
	}

	// Additional URL values for specific rule ids
	if ruleID == blacklistedCountries {
		if countries != "" {
			values.Add("countries", countries)
		}
		if continents != "" {
			values.Add("continents", continents)
		}
	} else if ruleID == blacklistedURLs {
		values.Add("urls", urls)
		values.Add("url_patterns", urlPatterns)
	} else if ruleID == blacklistedIPs || ruleID == whitelistedIPs {
		values.Add("ips", ips)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointACLRuleConfigure), values)
	if err != nil {
		return nil, fmt.Errorf("Error adding ACL for rule id %s and site id %d", ruleID, siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add ACL rule JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add ACL rule JSON response for rule id %s and site id %d", ruleID, siteID)
	}

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding ACL rule for rule id %s and site id %d: %s", ruleID, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}
