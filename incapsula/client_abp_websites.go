package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const resourceName = "ABP Websites"

func (c *Client) AbpTerraformUrl(accountId int) string {
	return fmt.Sprintf("%s/botmanagement/v1/account/%d/terraform", c.config.BaseURLAPI, accountId)
}

func (c *Client) CreateAbpWebsites(accountId int, account AbpTerraformAccount) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsitesWithBody(accountId, account, http.MethodPost, CreateAbpWebsites, "Creating", http.StatusCreated)
}

func (c *Client) UpdateAbpWebsites(accountId int, account AbpTerraformAccount) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsitesWithBody(accountId, account, http.MethodPut, UpdateAbpWebsites, "Updating", http.StatusOK)
}

func (c *Client) RequestAbpWebsitesWithBody(accountId int, account AbpTerraformAccount, method string, operation string, action string, successStatus int) (*AbpTerraformAccount, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] %s Abp websites Account ID %d\n", action, accountId)

	accountJson, err := json.Marshal(account)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure generating %s update request", resourceName),
			Detail:   fmt.Sprintf("Failed to JSON marshal AbpWebsites: %s", err.Error()),
		})
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] %s payload: %s\n", resourceName, string(accountJson))

	// Post form to Incapsula
	reqURL := c.AbpTerraformUrl(accountId)
	resp, err := c.DoJsonRequestWithHeaders(method, reqURL, accountJson, UpdateAbpWebsites)
	if err != nil {
		diags = append(diags, httpErrorDiagnostic(err, resourceName, accountId, method, action))
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		diags = append(diags, httpBodyErrorDiagnostic(err, resourceName, accountId, method, action, responseBody))
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update %s JSON response: %s\n", resourceName, string(responseBody))

	// Check the response code
	if resp.StatusCode != successStatus {
		diags = append(diags, httpStatusErrorDiagnostic(err, resourceName, accountId, method, action, resp, responseBody))
		return nil, diags
	}

	// Parse the JSON
	var newAbpWebsites AbpTerraformAccount
	err = json.Unmarshal([]byte(responseBody), &newAbpWebsites)
	if err != nil {
		diags = append(diags, jsonErrorDiagnostic(err, resourceName, accountId, method, responseBody))
		return nil, diags
	}

	return &newAbpWebsites, diags
}

func (c *Client) ReadAbpWebsites(accountId int) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsites(accountId, false, http.MethodGet, ReadAbpWebsites, "Reading", http.StatusOK)
}

func (c *Client) DeleteAbpWebsites(accountId int, autoPublish bool) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsites(accountId, autoPublish, http.MethodDelete, DeleteAbpWebsites, "Deleting", http.StatusOK)
}

func (c *Client) RequestAbpWebsites(accountId int, autoPublish bool, method string, operation string, action string, successStatus int) (*AbpTerraformAccount, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] %s Abp websites Account ID %d\n", action, accountId)

	// Post form to Incapsula
	var reqURL string
	if autoPublish {
		reqURL = fmt.Sprintf("%s?autoPublish=1", c.AbpTerraformUrl(accountId))
	} else {
		reqURL = c.AbpTerraformUrl(accountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(method, reqURL, nil, operation)
	if err != nil {
		diags = append(diags, httpErrorDiagnostic(err, resourceName, accountId, method, action))
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		diags = append(diags, httpBodyErrorDiagnostic(err, resourceName, accountId, method, action, responseBody))
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] Incapsula %s %s JSON response: %s\n", method, resourceName, string(responseBody))

	// Check the response code
	if resp.StatusCode != successStatus {
		diags = append(diags, httpStatusErrorDiagnostic(err, resourceName, accountId, method, action, resp, responseBody))
		return nil, diags
	}

	// Parse the JSON
	var newAbpWebsites AbpTerraformAccount
	err = json.Unmarshal([]byte(responseBody), &newAbpWebsites)
	if err != nil {
		diags = append(diags, jsonErrorDiagnostic(err, resourceName, accountId, method, responseBody))
		return nil, diags
	}

	return &newAbpWebsites, diags
}

func httpErrorDiagnostic(err error, resourceName string, accountId int, method string, action string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Failure %s %s", action, resourceName),
		Detail:   fmt.Sprintf("Error from Incapsula service when %s %s for Account ID %d: %s", strings.ToLower(action), resourceName, accountId, err.Error()),
	}
}

func httpBodyErrorDiagnostic(err error, resourceName string, accountId int, method string, action string, responseBody []byte) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Failure %s %s HTTP body", strings.ToLower(action), resourceName),
		Detail:   fmt.Sprintf("Error %s %s HTTP body for Account ID %d: %s\nresponse: %s", strings.ToLower(action), resourceName, accountId, err.Error(), string(responseBody)),
	}
}

func httpStatusErrorDiagnostic(err error, resourceName string, accountId int, method string, action string, resp *http.Response, responseBody []byte) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Failure %s %s", action, resourceName),
		Detail:   fmt.Sprintf("Error status code %d from Incapsula service when %s %s for Account ID %d: %s", resp.StatusCode, strings.ToLower(action), resourceName, accountId, string(responseBody)),
	}
}

func jsonErrorDiagnostic(err error, resourceName string, accountId int, method string, responseBody []byte) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Failure parsing %s %s response", resourceName, method),
		Detail:   fmt.Sprintf("Error parsing %s JSON response for Account ID %d: %s\nresponse: %s", resourceName, accountId, err.Error(), string(responseBody)),
	}
}
