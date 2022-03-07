package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CspDomainNote struct {
	Note string `json:"note"`
}

type CspDomainInfo struct {
	BaseDomain            string             `json:"baseDomain"`
	CompanyName           string             `json:"companyName"`
	DomainCategory        string             `json:"domainCategory"`
	Countries             []string           `json:"countries"`
	SSLCertificateInfo    string             `json:"sslCertificateInfo"`
	RegistrationTime      string             `json:"registrationTime"`
	Registrar             string             `json:"registrar"`
	OrgOwner              string             `json:"orgOwner"`
	DynamicDNSBased       bool               `json:"dynamicDnsBased"`
	DomainQuality         map[string]float64 `json:"domainQuality"`
	AdditionalInsights    []string           `json:"additionalInsights"`
	DomainCategorySemrush string             `json:"domainCategorySemrush"`
}

type CspDomainReport struct {
	DocumentUri string `json:"documentUri"`
	SourceFile  string `json:"sourceFile"`
	BlockedUri  string `json:"blockedUri"`
	LineNumber  int    `json:"lineNumber"`
	SourceType  string `json:"sourceType"`
}

// CspDomainData is the struct describing a csp site config response
type CspDomainData struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Status struct {
		Blocked  bool `json:"blocked"`
		Reviewed bool `json:"reviewed"`
	} `json:"status"`
	DomainRisk    string            `json:"domainRisk"`
	Notes         []CspDomainNote   `json:"notes"`
	TimeBucket    int64             `json:"timeBucket"`
	Significance  int               `json:"significance"`
	ResourceTypes []string          `json:"resourceTypes"`
	BrowserStats  map[string]int    `json:"browserStats"`
	CountryStats  map[string]int    `json:"countryStats"`
	IPSamples     []string          `json:"ipsSample"`
	Sources       int               `json:"sources"`
	DiscoveredAt  int64             `json:"discoveredAt"`
	LastSeenMs    int64             `json:"LastSeenMs"`
	DomainInfo    CspDomainInfo     `json:"domainInfo"`
	DomainReports []CspDomainReport `json:"domainReports"`
	PartOfProfile bool              `json:"partOfProfile"`
	Frequent      bool              `json:"frequent"`
}

type CspPreApprovedDomain struct {
	Domain      string `json:"domain"`
	Subdomains  bool   `json:"subdomains"`
	ReferenceID string `json:"referenceId"`
}

//type CspPreApprovedDomainsList []CspPreApprovedDomain
type CspPreApprovedDomainsMap map[string]CspPreApprovedDomain

func (c *Client) getCspDomainData(siteID int, domainRef string) (*CspDomainData, error) {
	log.Printf("[INFO] Getting CSP domain data for domain %s from site ID: %d\n", domainRef, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s/%s/%d/domains/%s", c.config.BaseURLAPI, CspSiteApiPath, siteID, domainRef),
		nil)
	if err != nil {
		return nil, fmt.Errorf("Error from CSP API for when getting domain %s from site ID %d: %s\n", domainRef, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API Get Domain Data JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from CSP API when getting domain data for domain %s from site %d: %s\n",
			resp.StatusCode, domainRef, siteID, string(responseBody))
	}

	// Parse the JSON
	var domainData CspDomainData
	err = json.Unmarshal([]byte(responseBody), &domainData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for domains %s data from site ID %d: %s\nresponse: %s\n",
			domainRef, siteID, err, string(responseBody))
	}

	return &domainData, nil
}

func (c *Client) getCspPreApprovedDomains(siteID int) (CspPreApprovedDomainsMap, error) {
	log.Printf("[INFO] Getting CSP pre-approved domains for site ID: %d\n", siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s%s/%d/preapprovedlist", c.config.BaseURLAPI, CspSiteApiPath, siteID),
		nil)
	if err != nil {
		return nil, fmt.Errorf("Error from CSP API for when getting pre-approved domains list for site ID %d: %s\n", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API Get Pre-Approved Domain Data JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from CSP API when getting pre-approved domains list for site %d: %s\n",
			resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var preApprovedList []CspPreApprovedDomain
	err = json.Unmarshal([]byte(responseBody), &preApprovedList)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for pre-approved domains list for site ID %d: %s\nresponse: %s\n",
			siteID, err, string(responseBody))
	}

	// var domMap = CspPreApprovedDomainsMap{}
	domMap := make(CspPreApprovedDomainsMap, len(preApprovedList))
	for i := range preApprovedList {
		dom := preApprovedList[i]
		domMap[dom.ReferenceID] = dom
	}

	return domMap, nil
}

func (c *Client) updateCspPreApprovedDomain(siteID int, dom *CspPreApprovedDomain) (*CspPreApprovedDomain, error) {
	log.Printf("[INFO] Updating CSP pre-approved domain for site ID: %d , domain: %v", siteID, dom)

	domJSON, err := json.Marshal(dom)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal CSP pre-approved domain %v: %s\n", dom, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s%s/%d/preapprovedlist", c.config.BaseURLAPI, CspSiteApiPath, siteID),
		domJSON)

	if err != nil {
		return nil, fmt.Errorf("Error from CSP API while updating pre-approved domain %v for site ID %d: %s\n", dom, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API Post Pre-Approved Domain Data JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Error status code %d from CSP API when updating pre-approved domain for site %d: %s\n",
			resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var updatedDom CspPreApprovedDomain
	err = json.Unmarshal([]byte(responseBody), &updatedDom)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for pre-approved domain %v for site ID %d: %s\nresponse: %s\n",
			dom, siteID, err, string(responseBody))
	}

	return &updatedDom, nil
}

func (c *Client) deleteCspPreApprovedDomains(siteID int, domainRef string) error {
	log.Printf("[INFO] Deleting CSP pre-approved domain %s for site ID: %d\n", domainRef, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete,
		fmt.Sprintf("%s%s/%d/preapprovedlist/%s", c.config.BaseURLAPI, CspSiteApiPath, siteID, domainRef),
		nil)
	if err != nil {
		return fmt.Errorf("Error from CSP API for when deleting pre-approved domain %s from site ID %d: %s\n", domainRef, siteID, err)
	}

	// Read the body
	defer resp.Body.Close() // Do I still need to?

	// Check the response code - no content for DELETE
	if resp.StatusCode != 204 {
		return fmt.Errorf("Error status code %d from CSP API when deleting pre-approved domain %s for site ID %d\n",
			resp.StatusCode, domainRef, siteID)
	}
	log.Printf("[DEBUG] CSP API Delete Pre-Approved Domain %s was successful\n", domainRef)

	return nil
}
