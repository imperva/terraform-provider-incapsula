package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const endpointSiteV3 = "/sites-mgmt/v3/sites"

// SiteAddResponse contains the relevant site information when adding an Incapsula managed site
type SiteV3Request struct {
	Id           int    `json:"id,omitempty"`
	AccountId    int    `json:"accountId,omitempty"`
	CreationTime int64  `json:"creationTime,omitempty"`
	Cname        string `json:"cname,omitempty"`
	Name         string `json:"name,omitempty"`
	SiteType     string `json:"type,omitempty"`
}

type SiteV3Response struct {
	Data   []SiteV3Request `json:"data"`
	Errors []APIErrors     `json:"errors"`
}

// AddV3Site adds a v3 site to be managed by Incapsula
func (c *Client) AddV3Site(siteV3Request *SiteV3Request, accountId string) (*SiteV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] adding v3 site to: %v ", siteV3Request)

	updateUrl := getSiteV3Url("", "", c.config.BaseURLAPI)
	siteV3RequestJson, err := json.Marshal(siteV3Request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse add site v3 dto",
			Detail:   fmt.Sprintf("Failed to parse add site v3 dto account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, updateUrl, siteV3RequestJson, nil, AddV3Site)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on add v3 site",
			Detail:   fmt.Sprintf("Failed to add v3 site account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on add v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva add v3 site JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on add v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteV3Response SiteV3Response
	err = json.Unmarshal(responseBody, &siteV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse add site v3 response",
			Detail:   fmt.Sprintf("Failed to parse add v3 site JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva add v3 site ended successfully for account id: %s", accountId)

	return &siteV3Response, nil
}

// UpdateV3Site update a v3 site currently managed by Incapsula
func (c *Client) UpdateV3Site(siteV3Request *SiteV3Request, accountId string) (*SiteV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] updating v3 site %d to: %v ", siteV3Request.Id, siteV3Request)

	updateUrl := getSiteV3Url("", "/"+strconv.Itoa(siteV3Request.Id), c.config.BaseURLAPI)
	siteV3RequestJson, err := json.Marshal(siteV3Request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse update site v3 dto",
			Detail:   fmt.Sprintf("Failed to parse update site v3 dto account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPatch, updateUrl, siteV3RequestJson, nil, UpdateV3Site)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on update v3 site",
			Detail:   fmt.Sprintf("Failed to update v3 site account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on update v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva update v3 site JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on update v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteV3Response SiteV3Response
	err = json.Unmarshal(responseBody, &siteV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse update site v3 response",
			Detail:   fmt.Sprintf("Failed to parse update v3 site JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva update v3 site ended successfully for account id: %s", accountId)

	return &siteV3Response, nil
}

// DeleteV3Site deletes a site currently managed by Incapsula
func (c *Client) DeleteV3Site(siteV3Request *SiteV3Request, accountId string) (*SiteV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] deleting v3 site %d to: %v ", siteV3Request.Id, siteV3Request)

	updateUrl := getSiteV3Url("", "/"+strconv.Itoa(siteV3Request.Id), c.config.BaseURLAPI)
	siteV3RequestJson, err := json.Marshal(siteV3Request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse delete site v3 dto",
			Detail:   fmt.Sprintf("Failed to parse delete site v3 dto account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, updateUrl, siteV3RequestJson, nil, UpdateV3Site)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on delete v3 site",
			Detail:   fmt.Sprintf("Failed to delete v3 site account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on delete v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva delete v3 site JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on delete v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteV3Response SiteV3Response
	err = json.Unmarshal(responseBody, &siteV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse delete site v3 response",
			Detail:   fmt.Sprintf("Failed to parse delete v3 site JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva delete v3 site ended successfully for account id: %s", accountId)

	return &siteV3Response, nil
}

// GetV3Site deletes a site currently managed by Incapsula
func (c *Client) GetV3Site(siteV3Request *SiteV3Request, accountId string) (*SiteV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] getting v3 site %d to: %v ", siteV3Request.Id, siteV3Request)

	updateUrl := getSiteV3Url("", "/"+strconv.Itoa(siteV3Request.Id), c.config.BaseURLAPI)
	siteV3RequestJson, err := json.Marshal(siteV3Request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse get site v3 dto",
			Detail:   fmt.Sprintf("Failed to parse get site v3 dto account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, updateUrl, siteV3RequestJson, nil, UpdateV3Site)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on get v3 site",
			Detail:   fmt.Sprintf("Failed to get v3 site account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on get v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva get v3 site JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on get v3 site",
			Detail:   fmt.Sprintf("Failed to read response for account id %s, got response status %d, %s", accountId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteV3Response SiteV3Response
	err = json.Unmarshal(responseBody, &siteV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse get site v3 response",
			Detail:   fmt.Sprintf("Failed to parse get v3 site JSON response for account %s, %s", accountId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva get v3 site ended successfully for account id: %s", accountId)

	return &siteV3Response, nil
}

func getSiteV3Url(accountId string, siteId string, baseUrl string) string {
	url := fmt.Sprintf("%s%s%s", baseUrl, endpointSiteV3, siteId)
	if accountId != "" {
		url = fmt.Sprintf("%s%s%s?caid=%s", baseUrl, endpointSiteV3, siteId, accountId)
	}
	return url
}
