package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const endpointDomainManagement = "/site-domain-manager/v2/sites/"

type AsyncResponseDataDto struct {
	Handler string `json:"handler"`
	Status  string `json:"status"`
}

type AsyncResponseDetailsDto struct {
	Data []AsyncResponseDataDto `json:"data"`
}

type SiteDomainsExtraDetailsDto struct {
	NumberOfAutoDiscoveredDomains int `json:"numberOfAutoDiscoveredDomains"`
	MaxAllowedDomains             int `json:"maxAllowedDomains"`
}

type SiteDomainsExtraDetailsResponse struct {
	Data []SiteDomainsExtraDetailsDto `json:"data"`
}

type SiteDomainDetails struct {
	Id                     int    `json:"id"`
	SiteId                 int    `json:"siteId"`
	Domain                 string `json:"domain"`
	AutoDiscovered         bool   `json:"autoDiscovered"`
	MainDomain             bool   `json:"mainDomain"`
	Managed                bool   `json:"managed"`
	CnameRedirectionRecord string `json:"cnameRedirectionRecord"`
	SubDomains             []struct {
		Id                 int    `json:"id"`
		SubDomain          string `json:"subDomain"`
		LastDiscoveredTime int64  `json:"lastDiscoveredTime"`
		CreationTime       int64  `json:"creationTime"`
	} `json:"subDomains"`
	CantDetachSubDomains bool   `json:"cantDetachSubDomains"`
	ValidationMethod     string `json:"validationMethod"`
	ValidationCode       string `json:"validationCode"`
	Status               string `json:"status"`
	CreationDate         int64  `json:"creationDate"`
}

type SiteDomainDetailsDto struct {
	Errors []ApiError          `json:"errors"`
	Data   []SiteDomainDetails `json:"data"`
}

type BulkAddDomainsDto struct {
	Data []DomainNameDto `json:"data"`
}

type DomainNameDto struct {
	Name string `json:"name"`
}

func (c *Client) GetWebsiteDomains(siteId string) (*SiteDomainDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteId, "/domains")
	if siteId != "" {
		fmt.Errorf("[ERROR] site ID was not provided")
	}
	var params = map[string]string{}
	params["pageSize"] = "-1"
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domain configuration details %s: %s", siteId, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get domain management response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get domain details %s\n: %s\n%s", resp.StatusCode, siteId, err, string(responseBody))
	}

	var siteDomainDetailsResponse SiteDomainDetailsDto
	err = json.Unmarshal([]byte(responseBody), &siteDomainDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain details response for site ID %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if len(siteDomainDetailsResponse.Data) > 0 {
		return &siteDomainDetailsResponse, nil
	} else {
		return nil, fmt.Errorf("domains for siteId %s not found", siteId)
	}
}

func (c *Client) BulkUpdateDomainsToSite(siteID string, siteDomainDetails []SiteDomainDetails) error {
	err := verifyDomainsAmountBelowMaxAllowed(c, siteID, siteDomainDetails)
	if err != nil {
		return err
	}

	domainNames := make([]DomainNameDto, len(siteDomainDetails))
	var i = 0
	for _, siteDomainDetailsItem := range siteDomainDetails {
		domainNames[i] = DomainNameDto{Name: siteDomainDetailsItem.Domain}
		i++
	}

	addBulkDomainsDto := BulkAddDomainsDto{Data: domainNames}
	resp, err := handleAddBulkRequest(c, addBulkDomainsDto, siteID)
	if err != nil {
		return err
	}

	isStatusCompleted := false
asyncStatusLoop:
	for i := 1; i <= 15; i++ {
		asyncResponseDataDto, err := checkForAsyncRequestStatus(c, siteID, resp.Data[0].Handler)
		if err != nil {
			return err
		}

		status := asyncResponseDataDto.Data[0].Status
		log.Printf("iteration %d: update domains for site status: %s", i, status)
		switch status {
		case "IN_PROGRESS":
			time.Sleep(10 * time.Second)
		case "FAILED":
			return fmt.Errorf("async update domains for site returned FAILED status")
		case "COMPLETED_SUCCESSFULLY":
			isStatusCompleted = true
			break asyncStatusLoop
		}
	}

	if isStatusCompleted {
		return nil
	}
	return fmt.Errorf("async request status timeout")
}

func checkForAsyncRequestStatus(c *Client, siteId string, requestUuid string) (*AsyncResponseDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteId, "/domains/status/", requestUuid)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, UpdateDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when update domains for siteId %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get status for async request %s\n: %s\n%s", resp.StatusCode, requestUuid, err, string(responseBody))
	}

	var asyncResponseDetailsDto AsyncResponseDetailsDto
	err = json.Unmarshal([]byte(responseBody), &asyncResponseDetailsDto)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing async request status %s: %s\nresponse: %s", requestUuid, err, string(responseBody))
	}
	return &asyncResponseDetailsDto, nil
}

func handleAddBulkRequest(c *Client, bulkAddDomainsDto BulkAddDomainsDto, siteId string) (*AsyncResponseDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteId, "/domains")
	body, err := json.Marshal(bulkAddDomainsDto)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse bulkAddDomainsDto: %s ", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, body, UpdateDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when updating domains for site %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula update domain management response: %s\n", string(responseBody))

	if resp.StatusCode != 202 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula update domains for siteId %s\n: %s\n%s", resp.StatusCode, siteId, err, string(responseBody))
	}

	var asyncResponseDetailsDto AsyncResponseDetailsDto
	err = json.Unmarshal([]byte(responseBody), &asyncResponseDetailsDto)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing async request for update domains for siteId %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}
	return &asyncResponseDetailsDto, nil
}

func verifyDomainsAmountBelowMaxAllowed(c *Client, siteId string, siteDomainDetails []SiteDomainDetails) error {
	siteExtraDetails, err := GetSiteExtraDetails(c, siteId)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	var domainsAmountLimit = siteExtraDetails.MaxAllowedDomains
	var autoDiscoveredAmount = siteExtraDetails.NumberOfAutoDiscoveredDomains

	if (len(siteDomainDetails) + autoDiscoveredAmount) >= domainsAmountLimit {
		log.Printf("[INFO] the site has currently %d auto-discovered domains", autoDiscoveredAmount)
		log.Printf("[ERROR] you are trying to add %d domains which is above the limit", len(siteDomainDetails))
		return fmt.Errorf("amount of domains is above the limit")
	}
	return nil
}

func GetSiteExtraDetails(c *Client, siteID string) (*SiteDomainsExtraDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteID, "/domains/extraDetails")
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadDomainExtraDetails)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when geting site domains extra details %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula get site domains extra details JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get site domains extra details for siteId %s\n: %s\n%s", resp.StatusCode, siteID, err, string(responseBody))
	}

	// Dump JSON
	var siteExtraDetailsResponse SiteDomainsExtraDetailsResponse
	err = json.Unmarshal([]byte(responseBody), &siteExtraDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing add domains to site JSON response for site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &siteExtraDetailsResponse.Data[0], nil
}
