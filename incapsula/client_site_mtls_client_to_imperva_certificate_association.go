package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID int) (*ClientCaCertificateWithSites, error) {
	log.Printf("[INFO] Getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d", siteID)
	//reqURL := fmt.Sprintf("%s/sites/%d/client-certificates", c.config.BaseURLAPI, siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/accounts/%d/client-certificates/%d", c.config.BaseURLAPI, accountID, certificateID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, "")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting Site to mutual TLS Client to Imperva Certificate association for Site ID %d, certificate ID %d, account ID: %s", siteID, certificateID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get Site to mutual TLS Client to Imperva Certificate association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching Incapsula Site to mutual TLS Client to Imperva Certificate association for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Dump JSON
	var clientCaCertificate ClientCaCertificateWithSites
	err = json.Unmarshal([]byte(responseBody), &clientCaCertificate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Incapsula Site to mutual TLS Client to Imperva Certificate association JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	//var mtlsClientCertificate ClientCaCertificate
	if len(clientCaCertificate.AssignedSites) > 0 {
		log.Printf("%v", clientCaCertificate.AssignedSites)

		//filter out
		var found bool
		for _, entry := range clientCaCertificate.AssignedSites {
			log.Printf("%d", entry)
			if entry == siteID {
				found = true
				//mtlsClientCertificate = entry
				break
			}
		}
		if !found {
			return nil, nil
		}
		return &clientCaCertificate, nil
	} else {
		return nil, nil
	}
}

func (c *Client) CreateSiteMtlsClientToImpervaCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Updating Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d/client-certificates/%d", c.config.BaseURLAPI, siteID, certificateID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, "")
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Incapsula Site to mutual TLS Client to Imperva Certificate Association for certificate ID %d, Site ID %d\n%s", certificateID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula update Site to mutual TLS Client to Imperva Certificate Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on creating Site to mutual TLS Client to Imperva Certificate Association for Site ID %d, Certificate ID %d:\n%s", resp.StatusCode, siteID, certificateID, string(responseBody))
	}
	return nil
}

func (c *Client) DeleteSiteMtlsClientToImpervaCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Unassigning Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d/client-certificates/%d", c.config.BaseURLAPI, siteID, certificateID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, "")
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Incapsula Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d\n%s", certificateID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula delete Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d JSON response: %s\n", certificateID, siteID, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on deleting site to mutual TLS Client to Imperva Certificate Association for certificate ID %d for Site ID %d\n%s", resp.StatusCode, certificateID, siteID, string(responseBody))
	}
	//add logic for 404
	return nil
}
