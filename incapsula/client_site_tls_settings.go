package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointSiteTlsSettings = "/configuration/client-certificates"

type SiteTlsSettings struct {
	Mandatory                  bool     `json:"mandatory"`
	Ports                      []int    `json:"ports"`
	IsPortsException           bool     `json:"isPortsException"`
	Hosts                      []string `json:"hosts"`
	IsHostsException           bool     `json:"isHostsException"`
	Fingerprints               []string `json:"fingerprints"`
	ForwardToOrigin            bool     `json:"forwardToOrigin"`
	HeaderName                 string   `json:"headerName,omitempty"`
	HeaderValue                string   `json:"headerValue,omitempty"`
	IsDisableSessionResumption bool     `json:"isDisableSessionResumption"`
}

func (c *Client) GetSiteTlsSettings(siteID int) (*SiteTlsSettings, error) {
	log.Printf("[INFO] Getting Site TLS Settings for Site ID %d", siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d%s", c.config.BaseURLAPI, siteID, endpointSiteTlsSettings)

	//todo KATRIN add operation
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, "")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting Site TLS Settings for Site ID %d: %s", siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get Site TLS Settings JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on fetching Incapsula Site TLS Settings for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Dump JSON
	var siteTlsSettings SiteTlsSettings
	err = json.Unmarshal([]byte(responseBody), &siteTlsSettings)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Incapsula Site to mutual TLS Imperva to Origin Certificate association JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &siteTlsSettings, nil
}

func (c *Client) UpdateSiteTlsSetings(siteID int, siteTlsSettingsPayload SiteTlsSettings) (*SiteTlsSettings, error) {
	siteTlsSettingsJSON, err := json.Marshal(siteTlsSettingsPayload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal site TLS settings: %s", err)
	}
	log.Printf("ssl settings JSON:\n%s", siteTlsSettingsJSON)
	log.Printf("[INFO] Updating Site TLS Settings for Site ID %d", siteID)
	reqURL := fmt.Sprintf("%s/certificate-manager/v2/sites/%d/configuration/client-certificates", c.config.BaseURLAPI, siteID)

	//todo KATRIN add operation
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, siteTlsSettingsJSON, "")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error updating Site TLS Settings for Site ID %d: %s", siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Update Site TLS Settings for Site ID %d  JSON response: %s\n", siteID, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service on update Incapsula Site TLS Settings for Site ID %d\n: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Dump JSON
	var siteTlsSettings SiteTlsSettings
	err = json.Unmarshal([]byte(responseBody), &siteTlsSettings)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Incapsula Site TLS Settings JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &siteTlsSettings, nil
}
