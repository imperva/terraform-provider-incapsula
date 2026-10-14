package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ListV3Sites retrieves sites for an account, optionally filtering by name
func (c *Client) ListV3Sites(accountId string, siteName string) (*SiteV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] listing v3 sites for account id: %s (filter by name: %s)", accountId, siteName)

	listUrl := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointSiteV3)
	queryParams := map[string]string{}
	if accountId != "" {
		queryParams["caid"] = accountId
	}
	if siteName != "" {
		queryParams["names"] = siteName
	}
	log.Printf("[DEBUG] ListV3Sites - listUrl: %s, queryParams: %+v, apiID: %s", listUrl, queryParams, c.config.APIID)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, listUrl, nil, queryParams, ReadSite)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on list v3 sites",
			Detail:   fmt.Sprintf("Failed to list v3 sites account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	log.Printf("[DEBUG] ListV3Sites - HTTP Status: %d", resp.StatusCode)
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] ListV3Sites - Response Body: %s", string(responseBody))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on list v3 sites",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva list v3 sites JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on list v3 sites",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteV3Response SiteV3Response
	err = json.Unmarshal(responseBody, &siteV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse list v3 sites response",
			Detail:   fmt.Sprintf("Failed to parse v3 sites JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva list v3 sites ended successfully for account id: %s", accountId)

	return &siteV3Response, nil
}

// GetV3SiteByName retrieves a site by name within an account
func (c *Client) GetV3SiteByName(siteName string, accountId string) (*SiteV3Request, diag.Diagnostics) {
	var diags diag.Diagnostics

	if siteName == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Site name is required",
			Detail:   "Site name cannot be empty",
		})
		return nil, diags
	}

	log.Printf("[INFO] getting v3 site by name: %s (account_id: %s)", siteName, accountId)

	listResponse, diags := c.ListV3Sites(accountId, siteName)
	if len(diags) > 0 {
		return nil, diags
	}

	log.Printf("[DEBUG] GetV3SiteByName - listResponse.Data length: %d, response: %+v", len(listResponse.Data), listResponse)

	if len(listResponse.Data) != 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Site lookup error",
			Detail:   fmt.Sprintf("Expected exactly 1 site with name '%s' in account %s, but found %d", siteName, accountId, len(listResponse.Data)),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] GetV3SiteByName - Found site: %+v", listResponse.Data[0])
	return &listResponse.Data[0], nil
}
