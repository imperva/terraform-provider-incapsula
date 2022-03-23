package incapsula

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CSPDomainNote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
	Date   int64  `json:"date"`
}

type CSPDomainStatus struct {
	Blocked  *bool `json:"blocked"`
	Reviewed *bool `json:"reviewed"`
}

type CSPDomainInfo struct {
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

type CSPDomainReport struct {
	DocumentUri string `json:"documentUri"`
	SourceFile  string `json:"sourceFile"`
	BlockedUri  string `json:"blockedUri"`
	LineNumber  int    `json:"lineNumber"`
	SourceType  string `json:"sourceType"`
}

// CSPDomainData is the struct describing a csp site config response
type CSPDomainData struct {
	ID            string            `json:"id"`
	Domain        string            `json:"domain"`
	Status        CSPDomainStatus   `json:"status"`
	DomainRisk    string            `json:"domainRisk"`
	Notes         []CSPDomainNote   `json:"notes"`
	TimeBucket    int64             `json:"timeBucket"`
	Significance  int               `json:"significance"`
	ResourceTypes []string          `json:"resourceTypes"`
	BrowserStats  map[string]int    `json:"browserStats"`
	CountryStats  map[string]int    `json:"countryStats"`
	IPSamples     []string          `json:"ipsSample"`
	Sources       int               `json:"sources"`
	DiscoveredAt  int64             `json:"discoveredAt"`
	LastSeenMs    int64             `json:"LastSeenMs"`
	DomainInfo    CSPDomainInfo     `json:"domainInfo"`
	DomainReports []CSPDomainReport `json:"domainReports"`
	PartOfProfile bool              `json:"partOfProfile"`
	Frequent      bool              `json:"frequent"`
}

type CSPPreApprovedDomain struct {
	Domain      string `json:"domain"`
	Subdomains  bool   `json:"subdomains"`
	ReferenceID string `json:"referenceId"`
}

type CSPPreApprovedDomainsMap map[string]CSPPreApprovedDomain

func (c *Client) getCSPDomainAPI(siteID int, domain string, APIPath string, ret interface{}) error {
	log.Printf("[INFO] Getting CSP domain %s for domain %s from site ID: %d\n", APIPath, domain, siteID)

	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		strings.Trim(fmt.Sprintf("%s%s/%d/domains/%s/%s", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef, APIPath),
			"/"),
		nil)
	if err != nil {
		return fmt.Errorf("Error from CSP API for when getting domain %s for domain %s from site ID %d: %s\n", APIPath, domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API get domain %s data JSON response: %s\n", APIPath, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from CSP API when getting domain %s for domain %s from site %d: %s\n",
			resp.StatusCode, APIPath, domain, siteID, string(responseBody))
	}

	// Parse the JSON
	err = json.Unmarshal([]byte(responseBody), ret)
	if err != nil {
		return fmt.Errorf("Error parsing JSON response for domain %s for domain %s from site ID %d: %s\nresponse: %s\n",
			APIPath, domain, siteID, err, string(responseBody))
	}

	return nil
}

func (c *Client) getCSPDomainData(siteID int, domain string) (*CSPDomainData, error) {
	ret := &CSPDomainData{}
	if err := c.getCSPDomainAPI(siteID, domain, "", ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) getCSPDomainStatus(siteID int, domain string) (*CSPDomainStatus, error) {
	ret := &CSPDomainStatus{}
	if err := c.getCSPDomainAPI(siteID, domain, "status", ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) updateCSPDomainStatus(siteID int, domain string, status *CSPDomainStatus) (*CSPDomainStatus, error) {
	log.Printf("[INFO] Updating CSP domain status for domain %s from site ID: %d to: %v\n", domain, siteID, status)

	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal CSP domain status %v: %s\n", status, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut,
		fmt.Sprintf("%s%s/%d/domains/%s/status", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef),
		statusJSON)
	if err != nil {
		return nil, fmt.Errorf("Error from CSP API for when updating domain status for domain %s from site ID %d: %s\n", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API update domain status data JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from CSP API when updating domain status for domain %s from site %d: %s\n",
			resp.StatusCode, domain, siteID, string(responseBody))
	}

	// Parse the JSON
	st := &CSPDomainStatus{}
	err = json.Unmarshal([]byte(responseBody), st)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for domain status for domain %s from site ID %d: %s\nresponse: %s\n",
			domain, siteID, err, string(responseBody))
	}

	return st, nil
}

func (c *Client) getCSPDomainNotes(siteID int, domain string) ([]CSPDomainNote, error) {
	var ret []CSPDomainNote
	if err := c.getCSPDomainAPI(siteID, domain, "notes", &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) addCSPDomainNote(siteID int, domain string, note string) error {
	log.Printf("[INFO] Getting CSP domain notes for domain %s from site ID: %d\n", domain, siteID)

	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s%s/%d/domains/%s/notes", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef),
		[]byte(note))
	if err != nil {
		return fmt.Errorf("Error from CSP API for when getting domain notes for domain %s from site ID %d: %s\n", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] CSP API get domain notes data JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 201 {
		return fmt.Errorf("Error status code %d from CSP API when getting domain notes for domain %s from site %d: %s\n",
			resp.StatusCode, domain, siteID, string(responseBody))
	}

	// Parse the JSON
	var notes []CSPDomainNote
	err = json.Unmarshal([]byte(responseBody), &notes)
	if err != nil {
		return fmt.Errorf("Error parsing JSON response for domain notes for domain %s from site ID %d: %s\nresponse: %s\n",
			domain, siteID, err, string(responseBody))
	}

	return nil
}

func (c *Client) deleteCSPDomainNotes(siteID int, domain string) error {
	log.Printf("[INFO] Deleting CSP domain notes for domain %s from site ID: %d\n", domain, siteID)

	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete,
		fmt.Sprintf("%s%s/%d/domains/%s/notes", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef),
		nil)
	if err != nil {
		return fmt.Errorf("Error from CSP API for when deleting domain notes for domain %s from site ID %d: %s\n", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()

	// Check the response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("Error status code %d from CSP API when getting domain notes for domain %s from site %d\n",
			resp.StatusCode, domain, siteID)
	}

	return nil
}

func (c *Client) getCSPPreApprovedDomain(siteID int, domain string) (*CSPPreApprovedDomain, error) {
	log.Printf("[INFO] Getting CSP pre-approved domain %s from site ID: %d\n", domain, siteID)

	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s%s/%d/preapprovedlist/%s", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef),
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
		return nil, fmt.Errorf("Error status code %d from CSP API when getting pre-approved domain %s for site %d: %s\n",
			resp.StatusCode, domain, siteID, string(responseBody))
	}

	// Parse the JSON
	var preApprovedDomain CSPPreApprovedDomain
	err = json.Unmarshal([]byte(responseBody), &preApprovedDomain)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for pre-approved domain %s for site ID %d: %s\nresponse: %s\n",
			domain, siteID, err, string(responseBody))
	}

	return &preApprovedDomain, nil
}

func (c *Client) getCSPPreApprovedDomains(siteID int) (CSPPreApprovedDomainsMap, error) {
	log.Printf("[INFO] Getting CSP pre-approved domains for site ID: %d\n", siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		fmt.Sprintf("%s%s/%d/preapprovedlist", c.config.BaseURLAPI, CSPSiteApiPath, siteID),
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
	var preApprovedList []CSPPreApprovedDomain
	err = json.Unmarshal([]byte(responseBody), &preApprovedList)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for pre-approved domains list for site ID %d: %s\nresponse: %s\n",
			siteID, err, string(responseBody))
	}

	domainsMap := make(CSPPreApprovedDomainsMap, len(preApprovedList))
	for i := range preApprovedList {
		dom := preApprovedList[i]
		domainsMap[dom.ReferenceID] = dom
	}

	return domainsMap, nil
}

func (c *Client) updateCSPPreApprovedDomain(siteID int, dom *CSPPreApprovedDomain) (*CSPPreApprovedDomain, error) {
	log.Printf("[INFO] Updating CSP pre-approved domain for site ID: %d , domain: %v", siteID, dom)

	domJSON, err := json.Marshal(dom)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal CSP pre-approved domain %v: %s\n", dom, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		fmt.Sprintf("%s%s/%d/preapprovedlist", c.config.BaseURLAPI, CSPSiteApiPath, siteID),
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
	var updatedDom CSPPreApprovedDomain
	err = json.Unmarshal([]byte(responseBody), &updatedDom)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response for pre-approved domain %v for site ID %d: %s\nresponse: %s\n",
			dom, siteID, err, string(responseBody))
	}

	return &updatedDom, nil
}

func (c *Client) deleteCSPPreApprovedDomains(siteID int, domainRef string) error {
	log.Printf("[INFO] Deleting CSP pre-approved domain %s for site ID: %d\n", domainRef, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete,
		fmt.Sprintf("%s%s/%d/preapprovedlist/%s", c.config.BaseURLAPI, CSPSiteApiPath, siteID, domainRef),
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
