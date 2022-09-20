package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) GetSiteMtlsCertificateAssociation(siteID int) (*MTLSCertificate, error) {
	log.Printf("[INFO] Getting mTLS certificate for Site ID %d", siteID)
	reqURL := fmt.Sprintf("%s%s?siteId=%d", c.config.BaseURLAPI, endpointMTLSCertificate, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, "")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting Incapsula Site to mTLS certificate association for Site ID %d: %s", siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get mutual TLS Imperva to Origin Certificate JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching Incapsula Site to mutual TLS Imperva to Origin certificate association for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Dump JSON
	var mtlsCertificate MTLSCertificateResponse
	err = json.Unmarshal([]byte(responseBody), &mtlsCertificate)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Incapsula Site to mutual TLS Imperva to Origin Certificate association JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}
	if len(mtlsCertificate.Data) > 0 {
		return &mtlsCertificate.Data[0], nil
	} else {
		return nil, nil
	}
}

func (c *Client) CreateSiteMtlsCertificateAssociation(certificateID, siteID int) error {
	log.Printf("[INFO] Updating mTLS certificate ID %d for Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s%s/%d/associated-sites/%d", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, nil, "")
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
	log.Printf("[INFO] Unassigning mTLS certificate ID %d for Site ID %d", certificateID, siteID)
	reqURL := fmt.Sprintf("%s%s/%d/associated-sites/%d", c.config.BaseURLAPI, endpointMTLSCertificate, certificateID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, "")
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
	return nil
}