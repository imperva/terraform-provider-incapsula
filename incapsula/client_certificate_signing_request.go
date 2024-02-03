package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

// Endpoints
const endpointCertificateSigningRequestCreate = "sites/customCertificate/csr"

// CertificateSigningRequestCreateResponse contains confirmation of successful csr
type CertificateSigningRequestCreateResponse struct {
	Res        int    `json:"res,string"`
	ResMessage string `json:"res_message"`
	Status     string `json:"status"`
	CsrContent string `json:"csr_content"`
}

// CreateCertificateSigningRequest creates a Certificate Signing Request (CSR)
func (c *Client) CreateCertificateSigningRequest(siteID, domain, email, country, state, city, organization, organizationUnit string) (*CertificateSigningRequestCreateResponse, error) {

	log.Printf("[INFO] Creating certificate signing request for site_id: %s", siteID)

	values := url.Values{
		"site_id":           {siteID},
		"domain":            {domain},
		"email":             {email},
		"country":           {country},
		"state":             {state},
		"city":              {city},
		"organization":      {organization},
		"organization_unit": {organizationUnit},
	}

	if domain != "" {
		values.Set("domain", domain)
	}
	if email != "" {
		values.Set("email", email)
	}
	if country != "" {
		values.Set("country", country)
	}
	if state != "" {
		values.Set("state", state)
	}
	if city != "" {
		values.Set("city", city)
	}
	if organization != "" {
		values.Set("organization", organization)
	}
	if organizationUnit != "" {
		values.Set("organization_unit", organizationUnit)
	}
	log.Printf("CertificateSigningRequest\n%v", values)
	// Post to Incapsula
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateSigningRequestCreate)
	resp, err := c.PostFormWithHeaders(reqURL, values, CreateCertificateSigningRequest)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when creating certificate signing request for site_id %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula create certificate signing request JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var certificateSigningRequestCreateResponse CertificateSigningRequestCreateResponse
	err = json.Unmarshal([]byte(responseBody), &certificateSigningRequestCreateResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing create certificate signing request JSON response for site_id %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if certificateSigningRequestCreateResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when creating certificate signing request for site_id %s: %s", siteID, string(responseBody))
	}

	return &certificateSigningRequestCreateResponse, nil
}
