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

type ATOEndpointMitigationConfigurationDTO struct {
	AccountId    int    `json:"accountId"`
	SiteId       int    `json:"siteId"`
	EndpointId   string `json:"endpointId"`
	LowAction    string `json:"lowAction"`
	MediumAction string `json:"mediumAction"`
	HighAction   string `json:"highAction"`
}

func formNoMitigationConfigurationDTO(accountId, siteId int, endpointId string) *ATOEndpointMitigationConfigurationDTO {
	return &ATOEndpointMitigationConfigurationDTO{
		AccountId:    accountId,
		SiteId:       siteId,
		EndpointId:   endpointId,
		LowAction:    "NONE",
		MediumAction: "NONE",
		HighAction:   "NONE",
	}
}

// GetAtoEndpointMitigationConfigurationWithRetries Fetch the mitigation configuration for an endpoint
func (c *Client) GetAtoEndpointMitigationConfigurationWithRetries(accountId, siteId int, endpointId string) (*ATOEndpointMitigationConfigurationDTO, int, error) {
	// Since the newly created site can take upto 30 seconds to be fully configured, we per.si a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		//60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		atoEndpointMitigationConfigurationDTO, status, err := c.GetAtoEndpointMitigationConfiguration(accountId, siteId, endpointId)
		if err == nil {
			return atoEndpointMitigationConfigurationDTO, status, nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return nil, 0, lastError
}

func (c *Client) GetAtoEndpointMitigationConfiguration(accountId, siteId int, endpointId string) (*ATOEndpointMitigationConfigurationDTO, int, error) {
	log.Printf("[INFO] Getting ATO mitigation configuration for (Site Id: %d)\n", siteId)

	// Get request to ATO
	var reqURL string
	if accountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointATOMitigation)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointATOMitigation, accountId)
	}
	// Adding specific endpoint ID from the API spec at https://docs.imperva.com/bundle/account-takeover/page/account-takeover/ato-api-definition.htm
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, map[string]string{"endpointIds": endpointId}, ReadATOSiteMitigationConfigurationOperation)
	if err != nil {
		return nil, 0, fmt.Errorf("[Error] Error executing get ATO mitigation configuration request for site with id %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] ATO mitigation configuration JSON response: %s\n", string(responseBody))

	// Check for internal server error
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[Error] Error response from server for fetching ATO mitigation configuration for site : %d , endpointId : %s , Error : %s", siteId, endpointId, responseBody)
	}

	// Parse the JSON
	var atoMitigationItems []ATOEndpointMitigationConfigurationDTO
	err = json.Unmarshal(responseBody, &atoMitigationItems)

	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error in parsing JSON response for ATO mitigation configuration : %s", responseBody)
	}

	// Get the desired mitigation configuration for the endpoint specified
	var atoEndpointMitigationConfigurationDTO ATOEndpointMitigationConfigurationDTO

	for _, atoMitigationItem := range atoMitigationItems {
		if atoMitigationItem.EndpointId == endpointId {
			atoEndpointMitigationConfigurationDTO.AccountId = accountId
			atoEndpointMitigationConfigurationDTO.SiteId = siteId
			atoEndpointMitigationConfigurationDTO.EndpointId = atoMitigationItem.EndpointId
			atoEndpointMitigationConfigurationDTO.LowAction = atoMitigationItem.LowAction
			atoEndpointMitigationConfigurationDTO.MediumAction = atoMitigationItem.MediumAction
			atoEndpointMitigationConfigurationDTO.HighAction = atoMitigationItem.HighAction
			break
		}
	}

	// Endpoint was missing in the mitigation configuration
	if atoEndpointMitigationConfigurationDTO.EndpointId == "" {
		return formNoMitigationConfigurationDTO(accountId, siteId, endpointId), resp.StatusCode, nil
	}

	return &atoEndpointMitigationConfigurationDTO, resp.StatusCode, nil
}

func (c *Client) UpdateATOEndpointMitigationConfigurationWithRetries(atoSiteMitigationConfigurationDTO *ATOEndpointMitigationConfigurationDTO) error {
	// Since the newly created site can take upto 30 seconds to be fully configured, we perform a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		//60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		err := c.UpdateATOSiteMitigationConfiguration(atoSiteMitigationConfigurationDTO)
		if err == nil {
			return nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return lastError
}

func (c *Client) UpdateATOSiteMitigationConfiguration(atoSiteMitigationConfigurationDTO *ATOEndpointMitigationConfigurationDTO) error {

	log.Printf("[INFO] Updating ATO mitigation configuration for (Site Id: %d)\n", atoSiteMitigationConfigurationDTO.SiteId)

	// Form the request body
	var mitigationConfigBody = []ATOEndpointMitigationConfigurationDTO{*atoSiteMitigationConfigurationDTO}
	mitigationConfigurationJSON, err := json.Marshal(mitigationConfigBody)

	// verify site ID and account ID are not the default value for int type
	if atoSiteMitigationConfigurationDTO.SiteId == 0 {
		return fmt.Errorf("site_id is not specified in updating ATO Mitigation configuration")
	}
	var reqURL string
	if atoSiteMitigationConfigurationDTO.AccountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, atoSiteMitigationConfigurationDTO.SiteId, endpointATOMitigation)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, atoSiteMitigationConfigurationDTO.SiteId, endpointATOMitigation, atoSiteMitigationConfigurationDTO.AccountId)
	}

	// Update request to ATO
	response, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, mitigationConfigurationJSON, UpdateATOSiteMitigationConfigurationOperation)

	// Read the body
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)

	log.Printf("Updated ATO mitigation configuration with response : %s", responseBody)

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing update ATO mitigation configuratgion request for site with id %d: %s", atoSiteMitigationConfigurationDTO.SiteId, err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("[Error] Error executing update ATO mitigation configuration request for site with status %d: %s", atoSiteMitigationConfigurationDTO.SiteId, response.Status)
	}

	return nil

}

func (c *Client) DisableATOEndpointMitigationConfiguration(accountId, siteId int, endpointId string) error {
	log.Printf("[INFO] Disabling ATO site mitigation configuration for (Site Id: %d)\n", siteId)

	// We are using empty mitigation config array instead of assigning 'NONE' to all risk levels
	// This has the advantage of resetting config entirely instead of possible enum conversion issues in the future
	err := c.UpdateATOSiteMitigationConfiguration(formNoMitigationConfigurationDTO(accountId, siteId, endpointId))

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing disable ATO mitigation configuration request for site with id %d, endpoint with id %s: %s", siteId, endpointId, err)
	}

	return nil
}
