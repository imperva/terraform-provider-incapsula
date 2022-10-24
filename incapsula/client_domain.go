package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointDomainManagement = "/site-domain-manager/v2/sites/"

type AddSiteDomainDetails struct {
	Domain     string `json:"domain"`
	StrictMode bool   `json:"strictMode"`
}

type SiteDomainDetails struct {
	Id             int    `json:"id"`
	SiteId         int    `json:"siteId"`
	Domain         string `json:"domain"`
	AutoDiscovered bool   `json:"autoDiscovered"`
	MainDomain     bool   `json:"mainDomain"`
	Managed        bool   `json:"managed"`
	SubDomains     []struct {
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

type SiteDomainDetailsResponse struct {
	Data []SiteDomainDetails `json:"data"`
}

func (c *Client) GetWebsiteDomains(siteID string) (*SiteDomainDetailsResponse, error) {
	log.Printf("[INFO] list domains for given website")
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteID, "/domains")
	if siteID != "" {
		//todo - print error
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domnain management details %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get domain management JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get domain details %s\n: %s\n%s", resp.StatusCode, siteID, err, string(responseBody))
	}

	// Dump JSON
	var siteDomainDetailsResponse SiteDomainDetailsResponse
	err = json.Unmarshal([]byte(responseBody), &siteDomainDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain details JSON response for site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	if len(siteDomainDetailsResponse.Data) > 0 {
		return &siteDomainDetailsResponse, nil
	} else {
		return nil, fmt.Errorf("domains for siteId %s not found", siteID)
	}
}

func (c *Client) GetDomainDetails(siteID string, domainID string) (*SiteDomainDetails, error) {
	log.Printf("[INFO] get domain details")
	reqURL := fmt.Sprintf("%s%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteID, "/domains/", domainID)
	if siteID != "" {
		//todo - print error
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domnain management details %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get domain management JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula get domain details %s\n: %s\n%s", resp.StatusCode, siteID, err, string(responseBody))
	}

	// Dump JSON
	var siteDomainDetails SiteDomainDetails
	err = json.Unmarshal([]byte(responseBody), &SiteDomainDetails{})
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain details JSON response for site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	} else {
		return &siteDomainDetails, nil
	}
	//if len(siteDomainDetailsResponse.Data) > 0 {
	//	return &siteDomainDetailsResponse.Data[0], nil
	//} else {
	//	return nil, fmt.Errorf("Results for site ID %s not found 1", siteID)
	//}
}

func (c *Client) AddDomainToSite(siteID string, domain string) (*SiteDomainDetails, error) {
	log.Printf("[INFO] Adding domain management")
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomainManagement, siteID, "/domains")
	bodyMap := map[string]interface{}{}
	bodyMap["domain"] = domain
	bodyMap["strictMode"] = true //strictMode must be true for TF.

	//bodyNew, contentTypeNew := c.CreateFormDataBody(bodyMap)

	AddSiteDomainDetails, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal AddSiteDomainDetails: %s ", err)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, AddSiteDomainDetails, CreateDomain)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when creating domnain management details %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula add domain management JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula add domain to site %s\n: %s\n%s", resp.StatusCode, siteID, err, string(responseBody))
	}

	// Dump JSON
	var siteDomainDetails SiteDomainDetails
	err = json.Unmarshal([]byte(responseBody), &siteDomainDetails)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing add domain to site JSON response for site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}
	return &siteDomainDetails, nil
}
