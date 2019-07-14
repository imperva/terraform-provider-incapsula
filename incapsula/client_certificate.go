package incapsula

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
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
}

// CertificateEditResponse contains confirmation of successful upload of certificate
type CertificateEditResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// AddCertificate adds a custom SSL certificate to a site in Incapsula
func (c *Client) AddCertificate(site_id, certificate, private_key, passphrase string) (*CertificateAddResponse, error) {
	b64_certificate := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(certificate)))
	b64_private_key := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(private_key)))

	//log.Printf("[INFO] Adding custom certificate for site_id: %s\n base_64 certificate: %s\n base_64 private_key: %s", site_id, b64_certificate, b64_private_key)
	log.Printf("[INFO] Adding custom certificate for site_id: %s", site_id)

	// Post to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateAdd), url.Values{
		"api_id":      {c.config.APIID},
		"api_key":     {c.config.APIKey},
		"site_id":     {site_id},
		"certificate": {b64_certificate},
		"private_key": {b64_private_key},
		"passphrase":  {passphrase},
	})
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding custom certificate for site_id %s: %s", site_id, err)
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
		return nil, fmt.Errorf("Error parsing add custom certificate JSON response for site_id %s: %s\nresponse: %s", site_id, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if certificateAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding custom certificate for site_id %s: %s", site_id, string(responseBody))
	}

	return &certificateAddResponse, nil
}

// CertificateList gets the list of custom certificates for a site
func (c *Client) ListCertificates(site_id string) (*CertificateListResponse, error) {
	log.Printf("[INFO] Getting Incapsula site custom certificates (site_id: %s)\n", site_id)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateList), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {site_id},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting custom certificates for site_id %s: %s", site_id, err)
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
		return nil, fmt.Errorf("Error parsing certificates list JSON response for site_id: %s %s\nresponse: %s", site_id, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if certificateListResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting custom certificates list for site_id %s: %s", site_id, string(responseBody))
	}

	return &certificateListResponse, nil
}

// EditCertificate updates the custom certifiacte on an Incapsula site
func (c *Client) EditCertificate(site_id, certificate, private_key, passphrase string) (*CertificateEditResponse, error) {
	b64_certificate := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(certificate)))
	b64_private_key := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(private_key)))

	log.Printf("[INFO] Editing custom certificate for Incapsula site_id: %s\n", site_id)

	values := url.Values{
		"api_id":      {c.config.APIID},
		"api_key":     {c.config.APIKey},
		"site_id":     {site_id},
		"certificate": {b64_certificate},
		"private_key": {b64_private_key},
	}

	if passphrase != "" {
		values.Add("passphrase", passphrase)
	}

	// Post to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateEdit), values)
	if err != nil {
		return nil, fmt.Errorf("Error editing custom certificate for site_id: %s: %s", site_id, err)
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
		return nil, fmt.Errorf("Error parsing edit custom certificarte JSON response for site_id: %s: %s)", site_id, err)
	}

	// Look at the response status code from Incapsula
	if certificateEditResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when editing custom certificarte for site_id %s: %s", site_id, string(responseBody))
	}

	return &certificateEditResponse, nil
}

// DeleteCertificate deletes a custom certificate for a specific site in Incapsula
func (c *Client) DeleteCertificate(site_id string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type CertificateDeleteResponse struct {
		Res        interface{} `json:"res"`
		ResMessage string      `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula custom certificate for site_id: %s\n", site_id)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {site_id},
	})
	if err != nil {
		return fmt.Errorf("Error deleting custom certificate for site_id: %s %s", site_id, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete custom certificate JSON response for site_id %s: %s\n", site_id, string(responseBody))

	// Parse the JSON
	var certificateDeleteResponse CertificateDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &certificateDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error deleting custom certificate for site_id: %s %s", site_id, err)
	}

	// Res can sometimes oscillate between a string and number
	// We need to add safeguards for this inside the provider
	var resString string

	if resNumber, ok := certificateDeleteResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = certificateDeleteResponse.Res.(string)
	}

	// Look at the response status code from Incapsula data center
	if resString == "0" {
		return nil
	}

	return fmt.Errorf("Error from Incapsula service when deleting custom certificate for site_id %s %s\n", site_id, string(responseBody))
}
