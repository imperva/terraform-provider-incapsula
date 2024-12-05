package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const endpointSiteCertV3BasePath = "/certificates-ui/v3/sites/"
const endpointSiteCertV3Suffix = "/certificates/managed"
const endpointSSLValidationSuffix = "/certificates/managed/validate"

// SiteAddResponse contains the relevant site information when adding an Incapsula managed site

type SanDTO struct {
	SanId                int    `json:"sanId,omitempty"`
	SanValue             string `json:"sanValue,omitempty"`
	ValidationMethod     string `json:"validationMethod,omitempty"`
	ExpirationDate       int    `json:"expirationDate,omitempty"`
	Status               string `json:"status,omitempty"`
	StatusDate           int    `json:"statusDate,omitempty"`
	NumSitesCovered      int    `json:"numSitesCovered,omitempty"`
	VerificationCode     string `json:"verificationCode,omitempty"`
	CnameValidationValue string `json:"cnameValidationValue,omitempty"`
	AutoValidation       bool   `json:"autoValidation,omitempty"`
	ApproverFqdn         string `json:"approverFqdn,omitempty"`
	ValidationEmail      string `json:"validationEmail,omitempty"`
	DomainIds            []int  `json:"domainIds,omitempty"`
}
type CertificateDTO struct {
	Id                 int      `json:"id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Status             string   `json:"status,omitempty"`
	Type               string   `json:"type,omitempty"`
	ExpirationDate     int64    `json:"expirationDate,omitempty"`
	InRenewal          bool     `json:"inRenewal,omitempty"`
	RenewalCertOrderId string   `json:"renewalCertOrderId,omitempty"`
	OriginCertOrderId  string   `json:"originCertOrderId,omitempty"`
	Sans               []SanDTO `json:"sans,omitempty"`
}
type SiteCertificateDTO struct {
	SiteId                  int              `json:"siteId,omitempty"`
	DefaultValidationMethod string           `json:"defaultValidationMethod,omitempty"`
	CertificatesDetails     []CertificateDTO `json:"certificateDetails,omitempty"`
}
type SiteCertificateV3Response struct {
	Data   []SiteCertificateDTO `json:"data"`
	Errors []APIErrors          `json:"errors"`
}

// RequestSiteCertificate request site certificate
func (c *Client) RequestSiteCertificate(siteId int, validationMethod string, accountId *int) (*SiteCertificateV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] request site certificate to: %d ", siteId)
	siteCertificateDTO := SiteCertificateDTO{}
	siteCertificateDTO.DefaultValidationMethod = validationMethod
	siteCertificateDTOJSON, err := json.Marshal(siteCertificateDTO)
	url := fmt.Sprintf("%s%s%d%s", c.config.BaseURLAPI, endpointSiteCertV3BasePath, siteId, endpointSiteCertV3Suffix)
	if accountId != nil {
		url = fmt.Sprintf("%s?caid=%d", url, *accountId)
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, url, []byte(siteCertificateDTOJSON), nil, RequestSiteCert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on request site certificate",
			Detail:   fmt.Sprintf("Failed to request site certificate for site %d,%s", siteId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva request site certificate JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, got response status %d, %s", siteId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteCertificateV3Response SiteCertificateV3Response
	err = json.Unmarshal(responseBody, &siteCertificateV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse request site certificate response",
			Detail:   fmt.Sprintf("Failed to parse request site certificate JSON response for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva request site certificate ended successfully for site id: %d", siteId)

	return &siteCertificateV3Response, nil
}

// DeleteRequestSiteCertificate deletes a site certificate request
func (c *Client) DeleteRequestSiteCertificate(siteId int, accountId *int) (*SiteCertificateV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] deleting request site certificate %d", siteId)

	url := fmt.Sprintf("%s%s%d%s", c.config.BaseURLAPI, endpointSiteCertV3BasePath, siteId, endpointSiteCertV3Suffix)
	if accountId != nil {
		url = fmt.Sprintf("%s?caid=%d", url, *accountId)
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, url, []byte("{}"), nil, RequestSiteCert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on delete request site certificate",
			Detail:   fmt.Sprintf("Failed to delete request site certificate for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on delete request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva delete request site certificate JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on delete request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, got response status %d, %s", siteId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteCertificateV3Response SiteCertificateV3Response
	err = json.Unmarshal(responseBody, &siteCertificateV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse delete request site certificate response",
			Detail:   fmt.Sprintf("Failed to parse delete request site certificate JSON response for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva delete request site certificate ended successfully for site id: %d", siteId)

	return &siteCertificateV3Response, nil
}

// GetSiteCertificateRequestStatus get site cert request
func (c *Client) GetSiteCertificateRequestStatus(siteId int, accountId *int) (*SiteCertificateV3Response, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] get request site certificate status %d", siteId)

	url := fmt.Sprintf("%s%s%d%s", c.config.BaseURLAPI, endpointSiteCertV3BasePath, siteId, endpointSiteCertV3Suffix)
	if accountId != nil {
		url = fmt.Sprintf("%s?caid=%d", url, *accountId)
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, url, nil, nil, RequestSiteCert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on get request site certificate",
			Detail:   fmt.Sprintf("Failed to get request site certificate for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on get request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva get request site certificate JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on get request site certificate",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, got response status %d, %s", siteId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var siteCertificateV3Response SiteCertificateV3Response
	err = json.Unmarshal(responseBody, &siteCertificateV3Response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse get request site certificate response",
			Detail:   fmt.Sprintf("Failed to parse get request site certificate JSON response for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva get request site certificate ended successfully for site id: %d", siteId)

	return &siteCertificateV3Response, nil
}

// ValidateDomains SSL validation of domains
func (c *Client) ValidateDomains(siteId int, domainIds []int) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[INFO] request ssl validation to: %v of site %d", domainIds, siteId)

	url := fmt.Sprintf("%s%s%d%s", c.config.BaseURLAPI, endpointSiteCertV3BasePath, siteId, endpointSSLValidationSuffix)
	domainJson, err := json.Marshal(domainIds)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse domain Ids for ssl validation",
			Detail:   fmt.Sprintf("Failed to parse domain Ids %v for ssl validation %d, %s", domainIds, siteId, err.Error()),
		})
		return diags
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, url, domainJson, nil, RequestSiteCert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on request for ssl validation",
			Detail:   fmt.Sprintf("Failed to request site SSL validation for domains %v and site %d, %s", siteId, domainIds, err.Error()),
		})
		return diags
	}
	log.Printf("[DEBUG] Imperva ssl validation response: %d\n", resp.StatusCode)
	if resp.StatusCode != 201 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to request ssl validation",
			Detail:   fmt.Sprintf("Failed to request ssl validation for site id %d and domains %v, got response status %d", siteId, domainIds, resp.StatusCode),
		})
		return diags
	}
	defer resp.Body.Close()
	return diags
}
