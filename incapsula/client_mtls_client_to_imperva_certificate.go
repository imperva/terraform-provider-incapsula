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

type ClientCaCertificate struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ClientCaCertificateWithSites struct {
	//ClientCaCertificate ClientCaCertificate
	Id            int    `json:"id"`
	Name          string `json:"name"`
	AssignedSites []int  `json:"assignedSites"`
}

func (c *Client) GetClientCaCertificate(accountID, certificateID string) (*ClientCaCertificateWithSites, bool, error) {
	log.Printf("[INFO] Reading mutual TLS Client To Imperva Certificate for ID %s, Account ID %s", certificateID, accountID)

	//todo refactor !! move to separate method baseURLAPI
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/accounts/%s/client-certificates/%s", c.config.BaseURLAPI, accountID, certificateID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadMtlsClientToImpervaCertifiate)
	if err != nil {
		return nil, true, fmt.Errorf("[ERROR] Error from Incapsula service when reading mTLS Client CA to Imperva Certificate ID %s: %s", certificateID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get mutual TLS Client To Imperva Certificate ID %s JSON response: %s\n", accountID, string(responseBody))

	// Check if certificate exists
	if resp.StatusCode == 406 && strings.HasPrefix(string(responseBody), "{") {
		return nil, false, nil
	}
	if resp.StatusCode != 200 {
		return nil, true, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching TLS Client to Imperva certificate ID %s\n: %s\n%s", resp.StatusCode, certificateID, err, string(responseBody))
	}

	// Dump JSON
	var clientCaCertificateWithSites ClientCaCertificateWithSites
	err = json.Unmarshal([]byte(responseBody), &clientCaCertificateWithSites)
	if err != nil {
		return nil, true, fmt.Errorf("[ERROR] Error parsing mutual GET TLS Client To Imperva Certificate for Account ID %s JSON response: %s\nresponse: %s", accountID, err, string(responseBody))
	}

	return &clientCaCertificateWithSites, true, nil
}

func (c *Client) AddClientCaCertificate(certificate []byte, accountID, certificateName string) (*ClientCaCertificate, error) {
	log.Printf("[INFO] Creating mutual TLS Client To Imperva Certificate for Account ID %s", accountID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/accounts/%s/client-certificates", c.config.BaseURLAPI, accountID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//certificate file
	fw, err := writer.CreateFormFile("ca_file", filepath.Base("certificateFile.pfx"))
	if err != nil {
		log.Printf("failed to create %s formdata field", "certificateFile")
	}
	fw.Write([]byte(certificate))

	//certificate name
	if certificateName != "" {
		fw, err := writer.CreateFormField("name")
		if err != nil {
			log.Printf("failed to create %s formdata field", "passphrase")
		}
		_, err = io.Copy(fw, strings.NewReader(certificateName))
		if err != nil {
			log.Printf("failed to write %s formdata field", "cartificate_name")
		}
	}

	writer.Close()

	contentType := writer.FormDataContentType()
	//rename method
	resp, err := c.DoJsonRequestWithHeadersForm(http.MethodPost, reqURL, body.Bytes(), contentType, CreateMtlsClientToImpervaCertifiate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula while creating mutual TLS Client To Imperva Certificate for Account ID %s", accountID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula create mutual TLS Client To Imperva Certificate for Account ID %s JSON response: %s\n", accountID, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on create mutual TLS Client To Imperva certificate for account ID %s : %s", resp.StatusCode, accountID, string(responseBody))
	}

	// Dump JSON
	var clientCaCertificateList []ClientCaCertificate
	err = json.Unmarshal([]byte(responseBody), &clientCaCertificateList)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing ADD mutual TLS Client To Imperva Certificate for Account ID %s JSON response: %s\nresponse: %s", accountID, err, string(responseBody))
	}

	if len(clientCaCertificateList) < 1 {
		return nil, fmt.Errorf("[ERROR] Failed to create mutual TLS Client To Imperva Certificate for Account ID %s")
	}
	return &clientCaCertificateList[0], nil
}

func (c *Client) DeleteClientCaCertificate(accountID, certificateID string) error {
	log.Printf("[INFO] Deleting mutual TLS Client To Imperva Certificate ID %s, Account ID %s", certificateID, accountID)

	reqURL := fmt.Sprintf("%s/certificate-manager/v2/accounts/%s/client-certificates/%s", c.config.BaseURLAPI, accountID, certificateID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteMtlsClientToImpervaCertifiate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error from Incapsula service when deletingmutual TLS Client To Imperva Certificate ID %s: %s", certificateID, err)
	}

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on deleting mutual TLS Client To Imperva Certificate ID %s\n: %v", resp.StatusCode, certificateID, err)
	}

	// Read the body
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting mutual TLS Client To Imperva Certificate ID %s: %s", certificateID, err)
	}
	return nil
}
