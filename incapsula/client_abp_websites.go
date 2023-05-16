package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const resourceName = "ABP Websites"

func (c *Client) AbpBaseUrl() string {
	// Useful for local testing of the ABP API
	url, ok := os.LookupEnv("INCAPSULA_ABP_BASE_URL_OVERRIDE")
	if ok {
		return url
	} else {
		return c.config.BaseURLAPI
	}
}

var authHeader = map[string]string{
	"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2ODMyMDY4ODksImlzcyI6Im1hcmt1cy53ZXN0ZXJsaW5kIiwic3ViIjoiYjYzNGNkMTMtNjQzZS01YzM0LThiMGEtNzU5MDgxZGM2ODM4IiwiZW1haWwiOiJtYXJrdXMud2VzdGVybGluZEBpbXBlcnZhLmNvbSIsImF1ZCI6ImJvbl9hY2NvdW50X2NvbmZpZyIsImV4cCI6MTY4NTc5ODg4OSwicm9sZSI6InN1cGVyYWRtaW4iLCJhY2NvdW50X2lkIjoiZjUxZTEwYTMtM2RjZC00NDM5LWI4ZTgtZGNiMTQ4N2Q0ZDMwIn0.EZmMBQCg7jTqCVvRUrCVN3NU-wrstRvFeOWRV0z_JXg",
}

func (c *Client) CreateAbpWebsites(accountId string, account AbpTerraformAccount) (*AbpTerraformAccount, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Creating Abp websites Account ID %s\n", accountId)

	accountJson, err := json.Marshal(account)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure generating %s create request", resourceName),
			Detail:   fmt.Sprintf("Failed to JSON marshal AbpWebsites: %s", err.Error()),
		})
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] %s payload: %s\n", resourceName, string(accountJson))

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/v1/account/%s/terraform", c.AbpBaseUrl(), accountId)
	resp, err := c.DoJsonRequestWithCustomHeaders(http.MethodPost, reqURL, accountJson, authHeader, CreateAbpWebsites)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure Creating %s", resourceName),
			Detail:   fmt.Sprintf("Error from Incapsula service when creating %s for Account ID %s: %s", resourceName, accountId, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Create %s JSON response: %s\n", resourceName, string(responseBody))

	// Check the response code
	if resp.StatusCode != 201 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure Creating %s", resourceName),
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when creating %s for Account ID %s: %s", resp.StatusCode, resourceName, accountId, string(responseBody)),
		})
	}

	// Parse the JSON
	var newAbpWebsites AbpTerraformAccount
	err = json.Unmarshal([]byte(responseBody), &newAbpWebsites)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure parsing %s create response", resourceName),
			Detail:   fmt.Sprintf("Error parsing %s JSON response for Account ID %s: %s\nresponse: %s", resourceName, accountId, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &newAbpWebsites, diags
}

func (c *Client) UpdateAbpWebsites(accountId string, account AbpTerraformAccount) (*AbpTerraformAccount, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Updating Abp websites Account ID %s\n", accountId)

	accountJson, err := json.Marshal(account)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure generating %s create request", resourceName),
			Detail:   fmt.Sprintf("Failed to JSON marshal AbpWebsites: %s", err.Error()),
		})
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] %s payload: %s\n", resourceName, string(accountJson))

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/v1/account/%s/terraform", c.AbpBaseUrl(), accountId)
	resp, err := c.DoJsonRequestWithCustomHeaders(http.MethodPut, reqURL, accountJson, authHeader, UpdateAbpWebsites)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure Updating %s", resourceName),
			Detail:   fmt.Sprintf("Error from Incapsula service when creating %s for Account ID %s: %s", resourceName, accountId, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update %s JSON response: %s\n", resourceName, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure Updating %s", resourceName),
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when creating %s for Account ID %s: %s", resp.StatusCode, resourceName, accountId, string(responseBody)),
		})
	}

	// Parse the JSON
	var newAbpWebsites AbpTerraformAccount
	err = json.Unmarshal([]byte(responseBody), &newAbpWebsites)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure parsing %s create response", resourceName),
			Detail:   fmt.Sprintf("Error parsing %s JSON response for Account ID %s: %s\nresponse: %s", resourceName, accountId, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &newAbpWebsites, diags
}

func (c *Client) ReadAbpWebsites(accountId string) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsites(accountId, false, http.MethodGet, ReadAbpWebsites, "Reading", 200)
}

func (c *Client) DeleteAbpWebsites(accountId string, autoPublish bool) (*AbpTerraformAccount, diag.Diagnostics) {
	return c.RequestAbpWebsites(accountId, autoPublish, http.MethodDelete, DeleteAbpWebsites, "Deleting", 200)
}

func (c *Client) RequestAbpWebsites(accountId string, autoPublish bool, method string, operation string, action string, successStatus int) (*AbpTerraformAccount, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] %s Abp websites Account ID %s\n", action, accountId)

	// Post form to Incapsula
	var reqURL string
	if autoPublish {
		reqURL = fmt.Sprintf("%s/v1/account/%s/terraform?autoPublish=1", c.AbpBaseUrl(), accountId)
	} else {
		reqURL = fmt.Sprintf("%s/v1/account/%s/terraform", c.AbpBaseUrl(), accountId)
	}
	resp, err := c.DoJsonRequestWithCustomHeaders(method, reqURL, nil, authHeader, operation)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure %s %s", action, resourceName),
			Detail:   fmt.Sprintf("Error from Incapsula service when %s %s for Account ID %s: %s", strings.ToLower(action), resourceName, accountId, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula %s %s JSON response: %s\n", method, resourceName, string(responseBody))

	// Check the response code
	if resp.StatusCode != successStatus {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure Creating %s", resourceName),
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when %s %s for Account ID %s: %s", resp.StatusCode, strings.ToLower(action), resourceName, accountId, string(responseBody)),
		})
	}

	// Parse the JSON
	var newAbpWebsites AbpTerraformAccount
	err = json.Unmarshal([]byte(responseBody), &newAbpWebsites)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failure parsing %s %s response", resourceName, method),
			Detail:   fmt.Sprintf("Error parsing %s JSON response for Account ID %s: %s\nresponse: %s", resourceName, accountId, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &newAbpWebsites, diags
}
