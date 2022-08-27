package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (c *Client) GetSiteMtlsClientToImpervaCertificateAssociation(siteID, certificateID int) (*ClientCaCertificateWithSites, bool, error) {
	log.Printf("[INFO] Getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d", siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d/client-certificates", c.config.BaseURLAPI, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadMtlsClientToImpervaCertifiateSiteAssociation)
	if err != nil {
		return nil, true, fmt.Errorf("[ERROR] Error getting Site to mutual TLS Client to Imperva Certificate association for Site ID %d, certificate ID %d", siteID, certificateID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get Site to mutual TLS Client to Imperva Certificate association JSON response: %s\n", string(responseBody))

	//check if association exists
	if resp.StatusCode == 401 && strings.HasPrefix(string(responseBody), "{") {
		return nil, false, nil
	}
	if resp.StatusCode != 200 {
		return nil, true, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching Incapsula Site to mutual TLS Client to Imperva Certificate association for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Dump JSON
	var clientCaCertificateList []ClientCaCertificateWithSites
	err = json.Unmarshal([]byte(responseBody), &clientCaCertificateList)
	if err != nil {
		return nil, true, fmt.Errorf("[ERROR] Error parsing Incapsula Site to mutual TLS Client to Imperva Certificate association JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	if len(clientCaCertificateList) > 0 {
		//filter out the certificate with relevant iD
		var found bool
		index := 0
		for i, entry := range clientCaCertificateList {
			if entry.Id == certificateID {
				found = true
				index = i
				break
			}
		}
		if !found {
			return nil, false, nil
		}
		return &clientCaCertificateList[index], true, nil
	} else {
		return nil, false, nil
	}
}

func (c *Client) CreateSiteMtlsClientToImpervaCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Updating Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d/client-certificates/%d", c.config.BaseURLAPI, siteID, certificateID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, CreateMtlsClientToImpervaCertifiateSiteAssociation)
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

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteMtlsClientToImpervaCertifiateSiteAssociation)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Incapsula Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d\n%s", certificateID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula delete Site to mutual TLS Client to Imperva Certificate Association certificate ID %d for Site ID %d JSON response: %s\n", certificateID, siteID, string(responseBody))

	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on deleting site to mutual TLS Client to Imperva Certificate Association for certificate ID %d for Site ID %d\n%s", resp.StatusCode, certificateID, siteID, string(responseBody))
	}
	return nil
}
