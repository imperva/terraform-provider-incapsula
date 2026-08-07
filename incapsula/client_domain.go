package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const endpointDomain = "/site-domain-manager/v2/sites/"

// addDomainMaxAttempts and addDomainRetryDelay govern the scoped retry
// behavior in AddDomainToSite for transient 401 responses caused by an
// authorization-scope propagation race shortly after site/cert-settings
// creation (CWMS-7457). They are vars (not consts) so tests can override
// them for fast execution.
var addDomainMaxAttempts = 3
var addDomainRetryDelay = 2 * time.Second

type AddSiteDetails struct {
	Domain     string `json:"domain"`
	StrictMode bool   `json:"strictMode"`
}

func (c *Client) GetDomain(siteId string, domainId string) (*SiteDomainDetails, error) {

	if siteId == "" {
		return nil, fmt.Errorf("[ERROR] site ID was not provided")
	}

	if domainId == "" {
		return nil, fmt.Errorf("[ERROR] domain ID was not provided")
	}

	reqURL := fmt.Sprintf("%s%s%s%s%s", c.config.BaseURLAPI, endpointDomain, siteId, "/domains/", domainId)

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, nil, ReadDomain)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domain details. domain id %s, site id %s: %s", siteId, domainId, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get domain response: %s\n", string(responseBody))

	var siteDomainDetailsResponse SiteDomainDetails
	err = json.Unmarshal(responseBody, &siteDomainDetailsResponse)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get domain response for site ID %s: Domain id: %s %s\nresponse: %s", siteId, domainId, err, string(responseBody))
	}

	return &siteDomainDetailsResponse, nil
}

func (c *Client) AddDomainToSite(siteID string, domainName string) (*SiteDomainDetails, error) {

	addDomainDto := AddSiteDetails{domainName, true}

	var resp *SiteDomainDetails
	var statusCode int
	var err error

	for attempt := 1; attempt <= addDomainMaxAttempts; attempt++ {
		resp, statusCode, err = handleAddDomainRequest(c, addDomainDto, siteID)

		if statusCode != http.StatusUnauthorized {
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		// Transient 401: the backend (site-domain-manager) occasionally
		// rejects the add-domain call shortly after site/cert-settings
		// creation while authorization scope is still propagating. Retry a
		// bounded number of times with a short delay before giving up.
		if attempt < addDomainMaxAttempts {
			log.Printf("[WARN] Incapsula add domain received 401 (unauthorized) for site %s, attempt %d/%d - retrying after transient auth propagation delay\n", siteID, attempt, addDomainMaxAttempts)
			time.Sleep(addDomainRetryDelay)
		}
	}

	return nil, fmt.Errorf("add domain request failed after %d attempts: 401 unauthorized (transient auth propagation delay - retry apply if this persists)", addDomainMaxAttempts)
}

func (c *Client) DeleteDomain(siteID string, domainId string) error {

	err := handleDeleteDomainRequest(c, siteID, domainId)

	if err != nil {
		return err
	}

	return nil
}

func handleAddDomainRequest(c *Client, addDomainsDto AddSiteDetails, siteId string) (*SiteDomainDetails, int, error) {
	reqURL := fmt.Sprintf("%s%s%s%s", c.config.BaseURLAPI, endpointDomain, siteId, "/domains")
	body, err := json.Marshal(addDomainsDto)

	if err != nil {
		return nil, 0, fmt.Errorf("Failed to parse addDomainsDto: %s ", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, body, CreateDomain)
	if err != nil {
		return nil, 0, fmt.Errorf("[ERROR] Error from Incapsula service when creating domains for site %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %s", err)
	}
	log.Printf("[DEBUG] Incapsula add domain response: %s\n", string(responseBody))

	var siteDomainDetails SiteDomainDetails
	err = json.Unmarshal(responseBody, &siteDomainDetails)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("[ERROR] Error parsing create domain response for siteId %s: %s\n response: %s", siteId, err, string(responseBody))
	}

	if siteDomainDetails.Errors != nil && len(siteDomainDetails.Errors) > 0 {
		log.Printf("[ERROR] Incapsula create domain failed for site: %s \n", siteId)
		return nil, resp.StatusCode, fmt.Errorf("add domain request failed (status %d): %s", resp.StatusCode, siteDomainDetails.Errors[0].Detail)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Incapsula create domain failed for site: %s \n", siteId)
		return nil, resp.StatusCode, fmt.Errorf("create request failed: %d", resp.StatusCode)
	}

	return &siteDomainDetails, resp.StatusCode, nil
}

func handleDeleteDomainRequest(c *Client, siteId string, domainId string) error {
	reqURL := fmt.Sprintf("%s%s%s%s%s", c.config.BaseURLAPI, endpointDomain, siteId, "/domains/", domainId)

	var params = map[string]string{}
	params["deleteLastDomain"] = "true"

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteDomain)

	if err != nil {
		return fmt.Errorf("[ERROR] Error from Incapsula service when deleting domains for site %s: %s", siteId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Incapsula delete domain failed for site: %s domain: %s \n", siteId, domainId)
		return fmt.Errorf("delete domain request failed: %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	log.Printf("[DEBUG] Incapsula delete domain response: %s\n", string(responseBody))

	return nil
}
