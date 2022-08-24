package incapsula

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) GetSiteMtlsCertificateAssociation(certificateID, siteID int) (bool, error) {
	log.Printf("[INFO] Getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d", siteID)
	reqURL := fmt.Sprintf("%s%s/%d/associated-sites/%d", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, "ReadSiteMtlsImpervaToOriginCertifiateAssociation")
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d: %s", siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get  Site to mutual TLS Imperva to Origin Certificate association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode == 404 {
		return false, err
	} else if resp.StatusCode == 200 {
		return true, err
	} else {
		return false, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching Incapsula Site to mutual TLS Imperva to Origin certificate association for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}
}

func (c *Client) CreateSiteMtlsCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Updating Site to mutual TLS Imperva to Origin Certificate association for certificate ID %d, Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s%s/%d/associated-sites/%d", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, nil, CreateSiteMtlsImpervaToOriginCertifiateAssociation)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Incapsula Site to Imperva to Origin mutual TLS Certificate Association for certificate ID %d, Site ID %d\n%s", certificateID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula update Site to mutual TLS Imperva to Origin Certificate association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on creating Incapsula Site to mutual TLS Imperva to Origin certificate Association for Site ID %d, Certificate ID %d:\n%s", resp.StatusCode, siteID, certificateID, string(responseBody))
	}
	return nil
}

func (c *Client) DeleteSiteMtlsCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Unassigning Site to mutual TLS Imperva to Origin Certificate association for certificate ID %d, Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s%s/%d/associated-sites/%d", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteSiteMtlsImpervaToOriginCertifiateAssociation)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Incapsula Site to Imperva to Origin mutual TLS Certificate Association for certificate ID %d for Site ID %d\n%s", certificateID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula delete Site to Imperva to Origin mutual TLS Certificate Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching site to mutual TLS Imperva to Origin certificate Association for certificate ID %d for Site ID %d\n%s", resp.StatusCode, certificateID, siteID, string(responseBody))
	}
	//add logic for 404
	return nil
}
