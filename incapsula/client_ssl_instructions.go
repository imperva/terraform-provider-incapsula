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

const endpointSSLInstructions = "/certificates-ui/v3/instructions/all"

// SiteAddResponse contains the relevant site information when adding an Incapsula managed site

type InstructionsDTO struct {
	Domain                         string          `json:"domain,omitempty"`
	ValidationMethod               string          `json:"validationMethod,omitempty"`
	RecordType                     string          `json:"recordType,omitempty"`
	VerificationCode               string          `json:"verificationCode,omitempty"`
	VerificationCodeExpirationDate int64           `json:"verificationCodeExpirationDate,omitempty"`
	LastNotificationDate           int64           `json:"lastNotificationDate,omitempty"`
	RelatedSansDetails             []SanDetailsDTO `json:"relatedSansDetails,omitempty"`
}
type SanDetailsDTO struct {
	SanId     int    `json:"sanId,omitempty"`
	SanValue  string `json:"sanValue,omitempty"`
	DomainIds []int  `json:"domainIds,omitempty"`
}
type SSLInstructionsResponse struct {
	Data   []InstructionsDTO `json:"data"`
	Errors []APIErrors       `json:"errors"`
}

// GetSiteSSLInstructions request site ssl instructions
func (c *Client) GetSiteSSLInstructions(siteId int) (*SSLInstructionsResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] request site SSL instructions to: %d ", siteId)

	url := fmt.Sprintf("%s%s?extSiteId=%s", c.config.BaseURLAPI, endpointSSLInstructions, strconv.Itoa(siteId))
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, url, []byte("{}"), nil, RequestSiteCert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error response from Imperva service on request site SSL instructions",
			Detail:   fmt.Sprintf("Failed to request site SSL instructions site %d,%s", siteId, err.Error()),
		})
		return nil, diags
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on request site SSL instructions",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}
	log.Printf("[DEBUG] Imperva request site SSL instructions JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read response from Imperva service on request site SSL instructions",
			Detail:   fmt.Sprintf("Failed to read response for site id %d, got response status %d, %s", siteId, resp.StatusCode, string(responseBody)),
		})
		return nil, diags
	}
	var sSLInstructionsResponse SSLInstructionsResponse
	err = json.Unmarshal(responseBody, &sSLInstructionsResponse)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse request site SSL instructions response",
			Detail:   fmt.Sprintf("Failed to parse request site SSL instructions JSON response for site %d, %s", siteId, err.Error()),
		})
		return nil, diags
	}

	log.Printf("[DEBUG] Imperva request site SSL instructions ended successfully for site id: %d", siteId)

	return &sSLInstructionsResponse, nil
}
