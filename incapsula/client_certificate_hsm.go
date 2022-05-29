package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointHsmCertificateAdd = "hsmCertificate"

type HSMDetailsDTO struct {
	KeyId    string `json:"key_id"`
	ApiKey   string `json:"api_key"`
	HostName string `json:"host_name"`
}

type HSMDataDTO struct {
	Certificate    string          `json:"certificate"`
	HsmDetailsList []HSMDetailsDTO `json:"hsmDetails"`
}

type HsmCustomCertificate struct {
	Data HSMDataDTO `json:"data"`
}

//-----

type HsmCertificatePutResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

func (c *Client) AddHsmCertificate(siteId, inputHash string, hSMDataDTO *HSMDataDTO) (*HsmCertificatePutResponse, error) {
	log.Printf("[INFO] Adding HSM certificate for site_id: %s with inputHash: %s ", siteId, inputHash)

	// Put to MY (This API using put, not post)
	reqURL := getHsmUrl(siteId, c)
	hsmCustomCertificate := HsmCustomCertificate{
		Data: *hSMDataDTO,
	}

	log.Printf("[INFO] Adding HSM certificate for site_id: %d with inputHash: %s with json: +%v", siteId, inputHash, hSMDataDTO)
	hSMDataDTOJSON, err := json.Marshal(hsmCustomCertificate)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal HSMDataDTO: %s ", err)
	}

	var params = map[string]string{}
	params["input_hash"] = inputHash
	log.Printf("[DEBUG] Add HSM certificate with params %s and JSON request: %s\n", params, string(hSMDataDTOJSON))
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPut, reqURL, hSMDataDTOJSON, params, CreateHSMCustomCertificate)
	if err != nil {
		return nil, fmt.Errorf("error from Imperva service when adding HSM certificate for site_id %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Imperva add HSM certificate JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var hsmCertificateAddResponse HsmCertificatePutResponse
	err = json.Unmarshal(responseBody, &hsmCertificateAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add HSM certificate JSON response for siteId %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if hsmCertificateAddResponse.Res != 0 {
		return nil, fmt.Errorf("error adding HSM certificate- res not 0. siteId: %s resposne:%s", siteId, string(responseBody))
	}

	log.Printf("[DEBUG] Imperva add HSM certificate clent pat ended successfully for site id: %s", siteId)

	return &hsmCertificateAddResponse, nil
}

func getHsmUrl(siteId string, c *Client) string {
	return fmt.Sprintf("%s/sites/%s/%s", c.config.BaseURLRev2, siteId, endpointHsmCertificateAdd)
}

//TODO: complete
// EditCertificate updates the custom certifiacte on an Incapsula site
//func (c *Client) EditCertificate(siteID, certificate, privateKey, passphrase, inputHash string) (*CertificateEditResponse, error) {
//
//	log.Printf("[INFO] Editing custom certificate for Incapsula site_id: %s\n", siteID)
//
//	values := url.Values{
//		"site_id":     {siteID},
//		"certificate": {certificate},
//		"input_hash":  {inputHash},
//	}
//
//	if privateKey != "" {
//		values.Set("private_key", privateKey)
//
//	}
//	if passphrase != "" {
//		values.Set("passphrase", passphrase)
//	}
//
//	// Post to Incapsula
//	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCertificateEdit)
//	resp, err := c.PostFormWithHeaders(reqURL, values, UpdateCustomCertificate)
//	if err != nil {
//		return nil, fmt.Errorf("Error editing custom certificate for site_id: %s: %s", siteID, err)
//	}
//
//	// Read the body
//	defer resp.Body.Close()
//	responseBody, err := ioutil.ReadAll(resp.Body)
//
//	// Dump JSON
//	log.Printf("[DEBUG] Incapsula edit custom certificate JSON response: %s\n", string(responseBody))
//
//	// Parse the JSON
//	var certificateEditResponse CertificateEditResponse
//	err = json.Unmarshal([]byte(responseBody), &certificateEditResponse)
//	if err != nil {
//		return nil, fmt.Errorf("Error parsing edit custom certificarte JSON response for site_id: %s: %s)", siteID, err)
//	}
//
//	// Look at the response status code from Incapsula
//	if certificateEditResponse.Res != 0 {
//		return nil, fmt.Errorf("Error from Incapsula service when editing custom certificarte for site_id %s: %s", siteID, string(responseBody))
//	}
//
//	return &certificateEditResponse, nil
//}

// DeleteHsmCustomCertificate deletes a hsm certificate for a specific site in Imperva
func (c *Client) DeleteHsmCertificate(siteId string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type CertificateDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Imperva HSM certificate for siteId: %s\n", siteId)

	// Post form to Incapsula
	reqURL := getHsmUrl(siteId, c)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, nil, DeleteHsmCustomCertificate)
	if err != nil {
		return fmt.Errorf("error deleting HSM certificate while sending request. siteId: %s %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Imperva service when deleting hsm certificate for site id %s: %s ", resp.StatusCode, siteId, string(responseBody))
	}

	if err != nil {
		return fmt.Errorf("Error reading response when deleting hsm certificate for site id %s: %s ", siteId, err)
	}

	log.Printf("[DEBUG] Imperva delete HSM certificate JSON response for siteId %s: %s\n", siteId, string(responseBody))

	// Parse the JSON
	var hsmCertificateDeleteResponse CertificateDeleteResponse
	err = json.Unmarshal(responseBody, &hsmCertificateDeleteResponse)
	if err != nil {
		return fmt.Errorf("error deleting HSM certificate, json parse error. siteId: %s %s", siteId, err)
	}

	if hsmCertificateDeleteResponse.Res != 0 {
		log.Printf("[DEBUG] response: %+v", hsmCertificateDeleteResponse)
		return fmt.Errorf("error deleting HSM certificate- res not 0. siteId: %s resposne:%s", siteId, string(responseBody))
	}

	return nil
}
