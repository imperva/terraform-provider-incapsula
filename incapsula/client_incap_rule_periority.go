package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-providers/terraform-provider-incapsula/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var diags diag.Diagnostics

func (c *Client) ReadIncapRulePriorities(siteID string, ruleType string) (*utils.Response, int, diag.Diagnostics) {
	log.Printf("[INFO] Getting Incapsula Incap ruleType Rule %s for Site ID %s\n", ruleType, siteID)

	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, ruleType)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadIncapRule)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error from Incapsula service when reading Incap Rule Categories :" + ruleType,
			Detail:   fmt.Sprintf("error from Incapsula service when reading Incap Rule Catagories %s for Site ID %s: %s", ruleType, siteID, err),
		})
		return nil, 0, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error status code from Incapsula service when reading Incap Rule Catagorie",
			Detail:   fmt.Sprintf("error status code %d from Incapsula service when reading Incap Rule Catagorie %s for Site ID %s: %s", resp.StatusCode, ruleType, siteID, string(responseBody)),
		})
		return nil, resp.StatusCode, diags
	}
	var rulesPriorities utils.Response
	err = json.Unmarshal(responseBody, &rulesPriorities)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Incap Rule Catagorie %s JSON response", ruleType),
			Detail:   fmt.Sprintf("Error parsing Incap Rule Catagorie %s JSON response for Site ID %s: %s\nresponse: %s", ruleType, siteID, err, string(responseBody)),
		})
		return nil, 0, diags
	}
	log.Printf("[INFO] Getting Incapsula Incap ruleType Rule %s for Site ID %s\n - finished", ruleType, siteID)

	return &rulesPriorities, resp.StatusCode, nil
}

func (c *Client) UpdateIncapRulePriorities(siteID string, ruleType string, rule []utils.RuleDetails) (*utils.Response, diag.Diagnostics) {
	log.Printf("[INFO] Updating Incapsula Incap Rule ruleType %s for Site ID %s\n", ruleType, siteID)
	ruleJSON, err := json.Marshal(rule)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to JSON marshal IncapRule priority"),
			Detail:   fmt.Sprintf("failed to JSON marshal IncapRule: %s", err),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Update rule DTO request: %v\n", string(ruleJSON[:]))

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, ruleType)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, ruleJSON, UpdateIncapRule)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error from Incapsula service when updating Incap Rule ruleType : %s", ruleType),
			Detail:   fmt.Sprintf("Error from Incapsula service when updating Incap Rule ruleType %s for Site ID %s: %s", ruleType, siteID, err),
		})
		return nil, diags
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("error status code %d from Incapsula service", resp.StatusCode),
			Detail:   fmt.Sprintf("error status code %d from Incapsula service when updating Incap Rule ruleType %s for Site ID %s: %s", resp.StatusCode, ruleType, siteID, string(responseBody)),
		})
		return nil, diags
	}
	var rulesPriorities utils.Response
	err = json.Unmarshal(responseBody, &rulesPriorities)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Incap Rule ruleType %s JSON response", ruleType),
			Detail:   fmt.Sprintf("Error parsing Incap Rule ruleType %s JSON response for Site ID %s: %s\nresponse: %s", ruleType, siteID, err, string(responseBody)),
		})
		return nil, diags
	}
	return &rulesPriorities, nil
}
