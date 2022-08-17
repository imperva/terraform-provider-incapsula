package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io/ioutil"
	"log"
	"net/http"
)

const accountSSLSettingsUrl = "/certificates-ui/v3/account/ssl-settings"

type Delegation struct {
	ValueForCNAMEValidation          string   `json:"valueForCNAMEValidation,omitempty"`
	AllowedDomainsForCNAMEValidation []string `json:"allowedDomainsForCNAMEValidation"`
	AllowCNAMEValidation             *bool    `json:"allowCNAMEValidation,omitempty"`
}

type ImpervaCertificate struct {
	AddNakedDomainSanForWWWSites *bool       `json:"addNakedDomainSanForWWWSites,omitempty"`
	UseWildCardSanInsteadOfFQDN  *bool       `json:"useWildCardSanInsteadOfFQDN,omitempty"`
	Delegation                   *Delegation `json:"delegation,omitempty"`
}

// AccountSSLSettingsDTO contains the SSL settings of an account
type AccountSSLSettingsDTO struct {
	ImpervaCertificate *ImpervaCertificate `json:"impervaCertificate,omitempty"`
}

type AccountSSLSettingsDTOResponse struct {
	Data   []AccountSSLSettingsDTO `json:"data"`
	Errors []APIErrors             `json:"errors"`
}

// UpdateAccountSSLSettings update account SSL settings
func (c *Client) UpdateAccountSSLSettings(accountSSLSettingsDTO *AccountSSLSettingsDTO, accountId string) (*AccountSSLSettingsDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] updating account SSL settings to: %v ", accountSSLSettingsDTO)

	updateUrl := getUrl(accountId, c.config.BaseURLAPI)
	accountSSLSettingsDTOJSON, err := json.Marshal(accountSSLSettingsDTO)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse update account SSL settings properties",
			Detail:   fmt.Sprintf("Failed to parse update account SSL settings properties for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, updateUrl, accountSSLSettingsDTOJSON, nil, UpdateAccountSSLSettings)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on update account SSL settings",
			Detail:   fmt.Sprintf("Failed to update account SSL settings for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on update account SSL settings",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva update account SSL settings JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on update account SSL settings",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var accountSSLSettingsDTOResponse AccountSSLSettingsDTOResponse
	err = json.Unmarshal(responseBody, &accountSSLSettingsDTOResponse)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse update account SSL settings response",
			Detail:   fmt.Sprintf("Failed to parse update account SSL settings JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva update account SSL settings ended successfully for account id: %s", accountId)

	return &accountSSLSettingsDTOResponse, nil
}

// GetAccountSSLSettings gets the Incapsula managed account's status
func (c *Client) GetAccountSSLSettings(accountId string) (*AccountSSLSettingsDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Getting account SSL settings of: %s ", accountId)

	getUrl := getUrl(accountId, c.config.BaseURLAPI)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, getUrl, nil, nil, GetAccountSSLSettings)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on getting account SSL settings",
			Detail:   fmt.Sprintf("Failed to get account SSL settings for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on getting account SSL settings",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva get account SSL settings for account %s response: %s\n", accountId, string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on getting account SSL settings",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var accountSSLSettingsDTOResponse AccountSSLSettingsDTOResponse
	err = json.Unmarshal(responseBody, &accountSSLSettingsDTOResponse)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse get account SSL settings response",
			Detail:   fmt.Sprintf("Failed to parse get account SSL settings JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] get account SSL settings ended successfully for account id: %s", accountId)
	return &accountSSLSettingsDTOResponse, nil
}

// DeleteAccountSSLSettings gets the Incapsula managed account's status
func (c *Client) DeleteAccountSSLSettings(accountId string) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[INFO] Reseting account SSL settings of: %s ", accountId)

	getUrl := getUrl(accountId, c.config.BaseURLAPI)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, getUrl, nil, nil, DeleteAccountSSLSettings)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete account SSL settings",
			Detail:   fmt.Sprintf("error from Imperva service when deleting Account SSL certificate for account_id  %s: %s", accountId, err.Error()),
		})
		return diags
	}

	// Read the body
	defer resp.Body.Close()

	log.Printf("[DEBUG] delete account SSL settings ended successfully for account id: %s", accountId)
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on update account SSL settings",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d", accountId, resp.StatusCode),
		})
		return diags
	}
	return nil
}

func getUrl(accountId string, baseUrl string) string {
	url := fmt.Sprintf("%s%s", baseUrl, accountSSLSettingsUrl)
	if accountId != "" {
		url = fmt.Sprintf("%s%s?caid=%s", baseUrl, accountSSLSettingsUrl, accountId)
	}
	return url
}
