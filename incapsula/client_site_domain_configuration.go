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
	Data   []AsyncResponseDataDto `json:"data"`
	Errors []ApiErrorResponse     `json:"errors"`
}

type SiteDomainsExtraDetailsDto struct {
	NumberOfAutoDiscoveredDomains int `json:"numberOfAutoDiscoveredDomains"`
	MaxAllowedDomains             int `json:"maxAllowedDomains"`
}

type SiteDomainsExtraDetailsResponse struct {
	Data   []SiteDomainsExtraDetailsDto `json:"data"`
	Errors []ApiErrorResponse           `json:"errors"`
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
	CantDetachSubDomains bool               `json:"cantDetachSubDomains"`
	ValidationMethod     string             `json:"validationMethod"`
	ValidationCode       string             `json:"validationCode"`
	Status               string             `json:"status"`
	CreationDate         int64              `json:"creationDate"`
	Errors               []ApiErrorResponse `json:"errors"`
}

type SiteDomainDetailsDto struct {
	Errors []ApiErrorResponse  `json:"errors"`
	Data   []SiteDomainDetails `json:"data"`
}

type BulkAddDomainsDto struct {
	Data   []DomainNameDto    `json:"data"`
	Errors []ApiErrorResponse `json:"errors"`
}

type DomainNameDto struct {
	Name string `json:"name"`
}

type ApiErrorResponse struct {
	ID     string         `json:"id"`
	Status int            `json:"status"`
	Title  string         `json:"title"`
	Detail string         `json:"detail"`
	Source ApiErrorSource `json:"source"`
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

	var siteDomainDetailsResponse SiteDomainDetailsDto
	err = json.Unmarshal([]byte(responseBody), &siteDomainDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain details response for site ID %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &siteDomainDetailsResponse, nil
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

	for i := 0; i < 15; i++ {
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
			return fmt.Errorf("async update domains for siteId %s returned FAILED status", siteID)
		case "COMPLETED_SUCCESSFULLY":
			return nil
		}
	}

	return fmt.Errorf("async request status for siteId %s reached timeout", siteID)
}

func checkForAsyncRequestStatus(c *Client, siteId string, requestUuid string) (*AsyncResponseDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteId, "/domains/status/", requestUuid)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, UpdateDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when update domains for siteId %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	var asyncResponseDetailsDto AsyncResponseDetailsDto
	err = json.Unmarshal([]byte(responseBody), &asyncResponseDetailsDto)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing async request status %s: %s\nresponse: %s", requestUuid, err, string(responseBody))
	}

	if asyncResponseDetailsDto.Errors != nil && len(asyncResponseDetailsDto.Errors) > 0 {
		return nil, fmt.Errorf("got error when trying to get async update domains requst for siteId %s: %s", siteId, asyncResponseDetailsDto.Errors[0].Detail)
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

	var asyncResponseDetailsDto AsyncResponseDetailsDto
	err = json.Unmarshal([]byte(responseBody), &asyncResponseDetailsDto)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing async request for update domains for siteId %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if asyncResponseDetailsDto.Errors != nil && len(asyncResponseDetailsDto.Errors) > 0 {
		return nil, fmt.Errorf("update domains async request failed: %s", asyncResponseDetailsDto.Errors[0].Detail)
	}
	return &asyncResponseDetailsDto, nil
}

func verifyDomainsAmountBelowMaxAllowed(c *Client, siteId string, siteDomainDetails []SiteDomainDetails) error {
	siteExtraDetails, err := GetSiteExtraDetails(c, siteId)
	if err != nil {
		return err
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

	var siteExtraDetailsResponse SiteDomainsExtraDetailsResponse
	err = json.Unmarshal([]byte(responseBody), &siteExtraDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing add domains to site JSON response for site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}
	if siteExtraDetailsResponse.Errors != nil && len(siteExtraDetailsResponse.Errors) > 0 {
		return nil, fmt.Errorf("got error when trying to get site extra details: %s", siteExtraDetailsResponse.Errors[0].Detail)
	}
	if len(siteExtraDetailsResponse.Data) > 0 {
		return &siteExtraDetailsResponse.Data[0], nil
	} else {
		return nil, fmt.Errorf("fail to get extra data for siteId %s", siteID)
	}

}
