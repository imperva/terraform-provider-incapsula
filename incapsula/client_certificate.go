package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

// Endpoints (unexported consts)
const endpointCertificateAdd = "sites/customCertificate/upload"
const endpointCertificateList = "sites/status"
const endpointCertificateEdit = "sites/customCertificate/upload"
const endpointCertificateDelete = "sites/customCertificate/remove"

// CertificateAddResponse contains confirmation of successful upload of certificate
type CertificateAddResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// CertificateListResponse contains site object with details of custom certificate
type CertificateListResponse struct {
	Res int `json:"res"`
	SSL SSL `json:"ssl"`
}

// CertificateEditResponse contains confirmation of successful upload of certificate
type CertificateEditResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	SSL        SSL    `json:"ssl"`
}

type SSL struct {
	CustomCertificate CustomCertificate `json:"custom_certificate"`
}

type CustomCertificate struct {
	InputHash string `json:"inputHash"`
}

// AddCertificate adds a custom SSL certificate to a site in Incapsula
func (c *Client) AddCertificate(siteID, certificate, privateKey, passphrase, authType, inputHash string) (*CertificateAddResponse, error) {

	log.Printf("[INFO] Adding custom certificate for site_id: %s", siteID)

	values := url.Values{
		"site_id":     {siteID},
		"certificate": {certificate},
		"input_hash":  {inputHash},
	}

	if privateKey != "" {
		values.Set("private_key", privateKey)
	}
	if passphrase != "" {
		values.Set("passphrase", passphrase)
	}
	if authType != "" {
		values.Set("auth_type", authType)
	}
	log.Printf("AddCertificate certificate\n%v", values)
	log.Printf("certificate\n%v", certificate)
	// Post to Incapsula
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateAdd)
	resp, err := c.PostFormWithHeaders(reqURL, values, CreateCustomCertificate)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding custom certificate for site_id %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add custom certificate JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var certificateAddResponse CertificateAddResponse
	err = json.Unmarshal([]byte(responseBody), &certificateAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add custom certificate JSON response for site_id %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if certificateAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding custom certificate for site_id %s: %s", siteID, string(responseBody))
	}

	return &certificateAddResponse, nil
}

// ListCertificates gets the list of custom certificates for a site
func (c *Client) ListCertificates(siteID, operation string) (*CertificateListResponse, error) {
	log.Printf("[INFO] Getting Incapsula site custom certificates (site_id: %s)\n", siteID)

	// Post form to Incapsula
	values := url.Values{"site_id": {siteID}}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateList)
	resp, err := c.PostFormWithHeaders(reqURL, values, operation)
	if err != nil {
		return nil, fmt.Errorf("Error getting custom certificates for site_id %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula list certificate (site status) JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var certificateListResponse CertificateListResponse
	err = json.Unmarshal([]byte(responseBody), &certificateListResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing certificates list JSON response for site_id: %s %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if certificateListResponse.Res != 0 {
		return &certificateListResponse, fmt.Errorf("Error from Incapsula service when getting custom certificates list for site_id %s: %s", siteID, string(responseBody))
	}

	return &certificateListResponse, nil
}

// EditCertificate updates the custom certifiacte on an Incapsula site
func (c *Client) EditCertificate(siteID, certificate, privateKey, passphrase, authType, inputHash string) (*CertificateEditResponse, error) {

	log.Printf("[INFO] Editing custom certificate for Incapsula site_id: %s\n", siteID)

	values := url.Values{
		"site_id":     {siteID},
		"certificate": {certificate},
		"input_hash":  {inputHash},
	}

	if privateKey != "" {
		values.Set("private_key", privateKey)

	}
	if passphrase != "" {
		values.Set("passphrase", passphrase)
	}
	if authType != "" {
		values.Set("auth_type", authType)
	}

	// Post to Incapsula
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateEdit)
	resp, err := c.PostFormWithHeaders(reqURL, values, UpdateCustomCertificate)
	if err != nil {
		return nil, fmt.Errorf("Error editing custom certificate for site_id: %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula edit custom certificate JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var certificateEditResponse CertificateEditResponse
	err = json.Unmarshal([]byte(responseBody), &certificateEditResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing edit custom certificarte JSON response for site_id: %s: %s)", siteID, err)
	}

	// Look at the response status code from Incapsula
	if certificateEditResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when editing custom certificarte for site_id %s: %s", siteID, string(responseBody))
	}

	return &certificateEditResponse, nil
}

// DeleteCertificate deletes a custom certificate for a specific site in Incapsula
func (c *Client) DeleteCertificate(siteID, authType string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type CertificateDeleteResponse struct {
		Res        interface{} `json:"res"`
		ResMessage string      `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula custom certificate for site_id: %s\n", siteID)

	// Post form to Incapsula
	values := url.Values{"site_id": {siteID}}

	if authType != "" {
		values.Set("auth_type", authType)
	}

	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateDelete)
	resp, err := c.PostFormWithHeaders(reqURL, values, DeleteCustomCertificate)
	if err != nil {
		return fmt.Errorf("Error deleting custom certificate for site_id: %s %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete custom certificate JSON response for site_id %s: %s\n", siteID, string(responseBody))

	// Parse the JSON
	var certificateDeleteResponse CertificateDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &certificateDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error deleting custom certificate for site_id: %s %s", siteID, err)
	}

	// Res can sometimes oscillate between a string and number
	// We need to add safeguards for this inside the provider
	var resString string

	if resNumber, ok := certificateDeleteResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = certificateDeleteResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString == "0" {
		return nil
	}

	return fmt.Errorf("Error from Incapsula service when deleting custom certificate for site_id %s %s", siteID, string(responseBody))
}
