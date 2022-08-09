package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const endpointMTLSCertificate = "/certificates-ui/v3/mtls-origin/certificates"

type MTLSCertificateGetById struct {
	Hash string `json:"hash"`
	Id   int    `json:"id"`
}

type MTLSCertificate struct {
	Id   int    `json:"certificateId"`
	Hash string `json:"hash"`
	Name string `json:"name"`
}

type MTLSCertificateResponse struct {
	Data []MTLSCertificate `json:"data"`
}

func (c *Client) AddMTLSCertificate(certificate, privateKey []byte, passphrase, certificateName, inputHash string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Adding mutual TLS Imperva to Origin Certificate")
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpointMTLSCertificate)
	return c.editMTLSCertificate(http.MethodPost, reqURL, certificate, privateKey, passphrase, certificateName, inputHash, "Create", CreateMtlsCertifiate)
}

func (c *Client) UpdateMTLSCertificate(certificateID string, certificate, privateKey []byte, passphrase, certificateName, inputHash string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Updating mutual TLS Imperva to Origin Certificate with ID %s", certificateID)
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)
	return c.editMTLSCertificate(http.MethodPut, reqURL, certificate, privateKey, passphrase, certificateName, inputHash, "Update", UpdateMtlsCertifiate)
}

func (c *Client) GetMTLSCertificate(certificateID string) (*MTLSCertificate, error) {
	log.Printf("[INFO] Reading mutual TLS Imperva to Origin Certificate with ID %s", certificateID)
	//todo refactor !! move to separate method
	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)

	//todo add operation!!!!!!!!!!!
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadMtlsCertifiate)
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
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("certificateFile", filepath.Base("certificateFile.pfx"))
	if err != nil {
		log.Printf("failed to create %s formdata field", "certificateFile")
	}
	fw.Write([]byte(certificate))

	if len(privateKey) > 0 {
		fw, err := writer.CreateFormFile("privateKeyFile", filepath.Base("privateKeyFile"))
		if err != nil {
			log.Printf("failed to create %s formdata field", "privateKeyFile")
		}
		fw.Write([]byte(privateKey))
	}

	//passphrase
	if passphrase != "" {
		fw, err := writer.CreateFormField("passphrase")
		if err != nil {
			log.Printf("failed to create %s formdata field", "passphrase")
		}
		_, err = io.Copy(fw, strings.NewReader(passphrase))
		if err != nil {
			log.Printf("failed to write %s formdata field", "passphrase")
		}
	}

	//certificateName
	if certificateName != "" {
		fw, err := writer.CreateFormField("certificateName")
		if err != nil {
			log.Printf("failed to create %s formdata field", "certificateName")
		}
		_, err = io.Copy(fw, strings.NewReader(certificateName))
		if err != nil {
			log.Printf("failed to write %s formdata field", "certificateName")
		}
	}

	//certificateName
	if inputHash != "" {
		fw, err := writer.CreateFormField("hash")
		if err != nil {
			log.Printf("failed to create %s formdata field", "hash")
		}
		_, err = io.Copy(fw, strings.NewReader(inputHash))
		if err != nil {
			log.Printf("failed to write %s formdata field", "hash")
		}
	}

	writer.Close()

	contentType := writer.FormDataContentType()
	resp, err := c.DoJsonRequestWithHeadersForm(hhtpMethod, reqURL, body.Bytes(), contentType, operation)
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
		// todo add ID if we have KATRIN!!!!!!
		return nil, fmt.Errorf("No cerificate found")
	}

}
func (c *Client) DeleteMTLSCertificate(certificateID string) error {
	log.Printf("[INFO] Deleting mTLS certificate with ID %s", certificateID)

	reqURL := fmt.Sprintf("%s%s/%s", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID)
	//todo add operation!!!!!!!!!!!

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteMtlsCertifiate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error from Incapsula service when deleting mTLS Imperva to Origin Certificate ID %s: %s", certificateID, err)
	}

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on deleting mutual TLS Imperva to Origin certificate ID %s\n: %s", resp.StatusCode, certificateID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula delete mutual TLS Imperva to Origin Certificate JSON response: %s\n", string(responseBody))

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting mTLS Imperva to Origin Certificate: %s", err)
	}
	return nil
}
