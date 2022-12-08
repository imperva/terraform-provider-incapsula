package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointDomainManagement = "/site-domain-manager/v2/sites/"

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
	log.Printf("[INFO] list domains for given website")
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

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get domain management JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get domain details %s\n: %s\n%s", resp.StatusCode, siteId, err, string(responseBody))
	}

	// Dump JSON
	var siteDomainDetailsResponse SiteDomainDetailsDto
	err = json.Unmarshal([]byte(responseBody), &siteDomainDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain details JSON response for site ID %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if len(siteDomainDetailsResponse.Data) > 0 {
		return &siteDomainDetailsResponse, nil
	} else {
		return nil, fmt.Errorf("domains for siteId %s not found", siteId)
	}
}

func (c *Client) BulkUpdateDomainsToSite(siteID string, siteDomainDetails []SiteDomainDetails) (*SiteDomainDetailsDto, error) {
	siteExtraDetails, err := GetSiteExtraDetails(c, siteID)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}

	var domainsAmountLimit = siteExtraDetails.MaxAllowedDomains
	var autoDiscoveredAmount = siteExtraDetails.NumberOfAutoDiscoveredDomains

	if (len(siteDomainDetails) + autoDiscoveredAmount) >= domainsAmountLimit {
		log.Printf("[DEBUG] the site has currently %d auto-discovered domains", autoDiscoveredAmount)
		log.Printf("[DEBUG] you are trying to add %d domains", len(siteDomainDetails))
		return nil, fmt.Errorf("amount of domains is above the limit")
	}

	domainNames := make([]DomainNameDto, len(siteDomainDetails))
	var i = 0
	for _, siteDomainDetailsItem := range siteDomainDetails {
		domainNames[i] = DomainNameDto{Name: siteDomainDetailsItem.Domain}
		i++
	}
	addBulkDomainsDtoA := BulkAddDomainsDto{Data: domainNames} //a full domain slice
	addBulkDomainsDtoB := BulkAddDomainsDto{}                  //a partial domain slices
	var splitThreshold = 500
	if len(addBulkDomainsDtoA.Data) > splitThreshold {
		//due to BE 1 min connection limitation need to split into 2 requests
		//todo -  run in a single request, after optimizing the BE
		addBulkDomainsDtoB = BulkAddDomainsDto{Data: addBulkDomainsDtoA.Data[0:splitThreshold]}
		handleAddBulkRequest(c, addBulkDomainsDtoB, siteID)
	}
	return handleAddBulkRequest(c, addBulkDomainsDtoA, siteID)
}

func handleAddBulkRequest(c *Client, bulkAddDomainsDto BulkAddDomainsDto, siteId string) (*SiteDomainDetailsDto, error) {
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteId, "/domains")

	body, err := json.Marshal(bulkAddDomainsDto)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal bulkAddDomainsDto: %s ", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, body, UpdateDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when bulk adding domains %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula add domain management JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula bulk add domains to site %s\n: %s\n%s", resp.StatusCode, siteId, err, string(responseBody))
	}

	// Dump JSON
	var siteDomainDto SiteDomainDetailsDto
	err = json.Unmarshal([]byte(responseBody), &siteDomainDto)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing bulk add domains to site JSON response for site ID %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}
	return &siteDomainDto, nil
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
