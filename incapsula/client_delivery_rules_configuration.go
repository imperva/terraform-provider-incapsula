package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io/ioutil"
	"log"
	"net/http"
)

type DeliveryRulesListResponseDTO struct {
	RuleDetails []DeliveryRuleDto `json:"data"`
	Errors      []APIErrors       `json:"errors"`
}

type DeliveryRuleDto struct {
	Name                    string `json:"rule_name"`
	Action                  string `json:"action"`
	Filter                  string `json:"filter,omitempty"`
	AddMissing              bool   `json:"add_if_missing,omitempty"`
	From                    string `json:"from,omitempty"`
	To                      string `json:"to,omitempty"`
	ResponseCode            int    `json:"response_code,omitempty"`
	RewriteExisting         bool   `json:"rewrite_existing,omitempty"`
	RewriteName             string `json:"rewrite_name,omitempty"`
	CookieName              string `json:"cookie_name"`
	HeaderName              string `json:"header_name"`
	DCID                    int    `json:"dc_id,omitempty"`
	PortForwardingContext   string `json:"port_forwarding_context,omitempty"`
	PortForwardingValue     string `json:"port_forwarding_value,omitempty"`
	ErrorType               string `json:"error_type,omitempty"`
	ErrorResponseFormat     string `json:"error_response_format,omitempty"`
	ErrorResponseData       string `json:"error_response_data,omitempty"`
	MultipleHeaderDeletions bool   `json:"multiple_headers_deletion"`
	Enabled                 bool   `json:"enabled"`
}

var diags diag.Diagnostics

func (c *Client) ReadIncapRulePriorities(siteID string, category string) (*DeliveryRulesListResponseDTO, int, diag.Diagnostics) {
	log.Printf("[INFO] Getting Delivery rules Type Rule %s for Site ID %s\n", category, siteID)

	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, category)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadIncapRule)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error from Incapsula service when reading Incap Rule Categories :" + category,
			Detail:   fmt.Sprintf("error from Incapsula service when reading Incap Rule Catagories %s for Site ID %s: %s", category, siteID, err),
		})
		return nil, 0, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error status code from Incapsula service when reading Incap Rule Catagorie",
			Detail:   fmt.Sprintf("error status code %d from Incapsula service when reading Incap Rule Catagorie %s for Site ID %s: %s", resp.StatusCode, category, siteID, string(responseBody)),
		})
		return nil, resp.StatusCode, diags
	}
	var rulesPriorities DeliveryRulesListResponseDTO
	err = json.Unmarshal(responseBody, &rulesPriorities)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Incap Rule Catagorie %s JSON response", category),
			Detail:   fmt.Sprintf("Error parsing Incap Rule Catagorie %s JSON response for Site ID %s: %s\nresponse: %s", category, siteID, err, string(responseBody)),
		})
		return nil, 0, diags
	}
	log.Printf("[INFO] Getting Delivery rules Type Rule %s for Site ID %s\n - finished", category, siteID)

	return &rulesPriorities, resp.StatusCode, nil
}

func (c *Client) UpdateIncapRulePriorities(siteID string, category string, rule []DeliveryRuleDto) (*DeliveryRulesListResponseDTO, diag.Diagnostics) {
	log.Printf("[INFO] Updating Delivery rules Type %s for Site ID %s\n", category, siteID)
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
	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, category)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, ruleJSON, UpdateIncapRule)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error from Incapsula service when updating Incap Rule category : %s", category),
			Detail:   fmt.Sprintf("Error from Incapsula service when updating Incap Rule category %s for Site ID %s: %s", category, siteID, err),
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
			Detail:   fmt.Sprintf("error status code %d from Incapsula service when updating Incap Rule category %s for Site ID %s: %s", resp.StatusCode, category, siteID, string(responseBody)),
		})
		return nil, diags
	}
	var rulesPriorities DeliveryRulesListResponseDTO
	err = json.Unmarshal(responseBody, &rulesPriorities)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Incap Rule category %s JSON response", category),
			Detail:   fmt.Sprintf("Error parsing Incap Rule category %s JSON response for Site ID %s: %s\nresponse: %s", category, siteID, err, string(responseBody)),
		})
		return nil, diags
	}
	return &rulesPriorities, nil
}
