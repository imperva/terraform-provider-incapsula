package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"path/filepath"
	//"strings"
	//"bytes"
	//"mime/multipart"
	//"io"
)

const endpointMTLSCertificate = "/certificates-ui/v3/mtls-origin/certificates"

type MTLSCertificateGetById struct {
	Hash      string `json:"hash"`
	Id        int    `json:"id"`
	AccountId int    `json:"accountId"`
}

type MTLSCertificate struct {
	Id        int    `json:"certificateId"`
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	AccountId int    `json:"accountId"`
}

type MTLSCertificateResponse struct {
	Data []MTLSCertificate `json:"data"`
}

func (c *Client) AddMTLSCertificate(certificate, privateKey []byte, passphrase, certificateName, inputHash, accountID string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Adding mutual TLS Imperva to Origin Certificate")
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointMTLSCertificate)
	if accountID != "" {
		reqURL = fmt.Sprintf("%s%s?caid=%s", c.config.BaseURLAPI, endpointMTLSCertificate, accountID)
	}
	return c.editMTLSCertificate(http.MethodPost, reqURL, certificate, privateKey, passphrase, certificateName, inputHash, "Create", CreateMtlsImpervaToOriginCertifiate)
}

func (c *Client) UpdateMTLSCertificate(certificateID string, certificate, privateKey []byte, passphrase, certificateName, inputHash, accountID string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Updating mutual TLS Imperva to Origin Certificate with ID %s", certificateID)
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)
	if accountID != "" {
		reqURL = fmt.Sprintf("%s%s/%s?caid=%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, accountID)
	}
	return c.editMTLSCertificate(http.MethodPut, reqURL, certificate, privateKey, passphrase, certificateName, inputHash, "Update", UpdateMtlsImpervaToOriginCertifiate)
}

func (c *Client) GetMTLSCertificate(certificateID, accountID string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Reading mutual TLS Imperva to Origin Certificate with ID %s", certificateID)
	//todo refactor !! move to separate method
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)
	if accountID != "" {
		reqURL = fmt.Sprintf("%s%s/%s?caid=%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, accountID)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadMtlsImpervaToOriginCertifiate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading mTLS Imperva to Origin Certificate ID %s: %s", certificateID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get mutual TLS Imperva to Origin Certificate JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching mutual TLS Imperva to Origin certificate ID %s\n: %s\n%s", resp.StatusCode, certificateID, err, string(responseBody))
	}

	// Dump JSON
	var mtlsCertificate MTLSCertificateResponse
	err = json.Unmarshal([]byte(responseBody), &mtlsCertificate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing mutual TLS Imperva to Origin Certificate JSON response for certificate ID %s: %s\nresponse: %s", certificateID, err, string(responseBody))
	}
	if len(mtlsCertificate.Data) > 0 {
		return &mtlsCertificate.Data[0], nil
	} else {
		return nil, fmt.Errorf("No cerificate with ID %s found", certificateID)
	}

}
func (c *Client) editMTLSCertificate(hhtpMethod, reqURL string, certificate, privateKey []byte, passphrase, certificateName, inputHash, action, operation string) (*MTLSCertificate, error) {
	bodyMap := map[string]interface{}{}
	bodyMap["certificateFile"] = []byte(certificate)

	if privateKey != nil && len(privateKey) > 0 {
		bodyMap["privateKeyFile"] = []byte(privateKey)
	}
	if passphrase != "" && passphrase != ignoreSensitivaeVariableString {
		bodyMap["passphrase"] = passphrase
	}
	//certificateName
	if certificateName != "" {
		bodyMap["certificateName"] = certificateName
	}
	//don't update hash if ignore all sensitive fields from account-export
	if inputHash != "" {
		bodyMap["hash"] = inputHash
	}

	bodyNew, contentTypeNew := c.CreateFormDataBody(bodyMap)
	resp, err := c.DoFormDataRequestWithHeaders(hhtpMethod, reqURL, bodyNew, contentTypeNew, CreateMtlsClientToImpervaCertifiate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error while %s mTLS Imperva to Origin Certificate: %s", action, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula %s mutual TLS Imperva to Origin Certificate JSON response: %s\n", action, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on %s mutual TLS Imperva to Origin certificate: %s", resp.StatusCode, action, string(responseBody))
	}

	// Dump JSON
	var mtlsCertificate MTLSCertificateResponse
	err = json.Unmarshal([]byte(responseBody), &mtlsCertificate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing mutual TLS Imperva to Origin Certificate JSON response: %s\nresponse: %s", err, string(responseBody))
	}
	if len(mtlsCertificate.Data) > 0 {
		return &mtlsCertificate.Data[0], nil
	} else {
		return nil, fmt.Errorf("No cerificate found")
	}

}
func (c *Client) DeleteMTLSCertificate(certificateID, accountID string) error {
	log.Printf("[INFO] Deleting mTLS certificate with ID %s", certificateID)

	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteMtlsImpervaToOriginCertifiate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error from Incapsula service when deleting mTLS Imperva to Origin Certificate ID %s: %s", certificateID, err)
	}

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on deleting mutual TLS Imperva to Origin certificate ID %s\n: %s", resp.StatusCode, certificateID, err)
	}

	// Read the body
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting mTLS Imperva to Origin Certificate: %s", err)
	}
	return nil
}
