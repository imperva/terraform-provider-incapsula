package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const endpointATOMitigation = "/mitigation"

type AtoMitigationItem struct {
	EndpointId   int    `json:"endpointId"`
	LowAction    string `json:"lowAction"`
	MediumAction string `json:"mediumAction"`
	HighAction   string `json:"highAction"`
}

type ATOSiteMitigationConfigurationDTO struct {
	AccountId               int                 `json:"accountId"`
	SiteId                  int                 `json:"siteId"`
	MitigationConfiguration []AtoMitigationItem `json:"mitigationConfiguration"`
}

func formNoMitigationConfigurationDTO(accountId, siteId int) *ATOSiteMitigationConfigurationDTO {
	return &ATOSiteMitigationConfigurationDTO{
		AccountId:               accountId,
		SiteId:                  siteId,
		MitigationConfiguration: make([]AtoMitigationItem, 0),
	}
}

func (atoMitigationItem *AtoMitigationItem) toMap() (map[string]interface{}, error) {

	// Initialize the data map that terraform uses
	var atoMitigationItemMap = make(map[string]interface{})

	// Set properties
	atoMitigationItemMap["endpointId"] = atoMitigationItem.EndpointId
	atoMitigationItemMap["lowAction"] = atoMitigationItem.LowAction
	atoMitigationItemMap["mediumAction"] = atoMitigationItem.MediumAction
	atoMitigationItemMap["highAction"] = atoMitigationItem.HighAction

	return atoMitigationItemMap, nil
}

func formAtoMitigationItemFromMap(atoMitigationItemMap map[string]interface{}) (*AtoMitigationItem, error) {

	atoMitigationItem := AtoMitigationItem{}

	// Set endpointId
	if atoMitigationItemMap["endpointId"] != nil {
		endpointId, ok := atoMitigationItemMap["endPointId"].(int)
		if ok {
			atoMitigationItem.EndpointId = endpointId
		} else {
			return nil, fmt.Errorf("endpointId is not a number. Received : %d", endpointId)
		}

	} else {
		return nil, fmt.Errorf("endpointId cannot be empty")
	}

	atoMitigationItem.LowAction = atoMitigationItemMap["lowAction"].(string)
	atoMitigationItem.MediumAction = atoMitigationItemMap["mediumAction"].(string)
	atoMitigationItem.HighAction = atoMitigationItemMap["highAction"].(string)

	return &atoMitigationItem, nil
}

func (c *Client) GetAtoSiteMitigationConfigurationWithRetries(accountId, siteId int) (*ATOSiteMitigationConfigurationDTO, error) {
	// Since the newly created site can take upto 30 seconds to be fully configured, we per.si a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		aatoSiteMitigationConfigurationDTO, err := c.GetAtoSiteMitigationConfiguration(accountId, siteId)
		if err == nil {
			return aatoSiteMitigationConfigurationDTO, nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return nil, lastError
}

func (c *Client) GetAtoSiteMitigationConfiguration(accountId, siteId int) (*ATOSiteMitigationConfigurationDTO, error) {
	log.Printf("[INFO] Getting ATO mitigation configuration for (Site Id: %d)\n", siteId)

	// Get request to ATO
	var reqURL string
	if accountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointATOMitigation)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointATOMitigation, accountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadATOSiteMitigationConfigurationOperation)
	if err != nil {
		return nil, fmt.Errorf("[Error] Error executing get ATO mitigation configuration request for site with id %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] ATO mitigation configuration JSON response: %s\n", string(responseBody))

	// Check for internal server error
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("[Error] Error response from server for fetching ATO mitigation configuration for site : %d , Error : %s", siteId, responseBody)
	}

	// Parse the JSON
	var atoMitigationItems []AtoMitigationItem
	var atoSiteMitigationConfigurationDTO ATOSiteMitigationConfigurationDTO
	err = json.Unmarshal(responseBody, &atoMitigationItems)

	if err != nil {
		return nil, fmt.Errorf("Error in parsing JSON response for ATO mitigation configuration : %s", responseBody)
	}

	atoSiteMitigationConfigurationDTO.SiteId = siteId
	atoSiteMitigationConfigurationDTO.AccountId = accountId
	atoSiteMitigationConfigurationDTO.MitigationConfiguration = atoMitigationItems

	return &atoSiteMitigationConfigurationDTO, nil
}

func (c *Client) UpdateATOSiteMitigationConfigurationWithRetries(atoSiteMitigationConfigurationDTO *ATOSiteMitigationConfigurationDTO) error {
	// Since the newly created site can take upto 30 seconds to be fully configured, we perform a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		err := c.UpdateATOSiteMitigationConfigurationWithRetries(atoSiteMitigationConfigurationDTO)
		if err == nil {
			return nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return lastError
}

func (c *Client) UpdateATOSiteMitigationConfiguration(aTOSiteMitigationConfigurationDTO *ATOSiteMitigationConfigurationDTO) error {

	log.Printf("[INFO] Updating ATO mitigation configuration for (Site Id: %d)\n", aTOSiteMitigationConfigurationDTO.SiteId)

	// Form the request body
	mitigationConfigurationJSON, err := json.Marshal(aTOSiteMitigationConfigurationDTO.MitigationConfiguration)

	// verify site ID and account ID are not the default value for int type
	if aTOSiteMitigationConfigurationDTO.SiteId == 0 {
		return fmt.Errorf("site_id is not specified in updating ATO Mitigation configuration")
	}
	var reqURL string
	if aTOSiteMitigationConfigurationDTO.AccountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, aTOSiteMitigationConfigurationDTO.SiteId, endpointATOMitigation)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, aTOSiteMitigationConfigurationDTO.SiteId, endpointAtoAllowlist, aTOSiteMitigationConfigurationDTO.AccountId)
	}

	// Update request to ATO
	response, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, mitigationConfigurationJSON, UpdateATOSiteAllowlistOperation)

	// Read the body
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)

	log.Printf("Updated ATO mitigation configuration with response : %s", responseBody)

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing update ATO mitigation configuratgion request for site with id %d: %s", aTOSiteMitigationConfigurationDTO.SiteId, err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("[Error] Error executing update ATO mitigation configuration request for site with status %d: %s", aTOSiteMitigationConfigurationDTO.SiteId, response.Status)
	}

	return nil

}

func (c *Client) DisableATOSiteMitigationConfiguration(accountId, siteId int) error {
	log.Printf("[INFO] Disabling ATO site mitigation configuration for (Site Id: %d)\n", siteId)

	err := c.UpdateATOSiteMitigationConfiguration(formNoMitigationConfigurationDTO(accountId, siteId))

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing disable ATO mitigation configuration request for site with id %d: %s", siteId, err)
	}

	return nil
}
