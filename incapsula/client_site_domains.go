package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
)

const endpointSiteDomainsV3 = "/site-domain-manager/v3/sites/"

const siteDomainsV3PageSize = 100

type SiteDomainV3 struct {
	Id           int      `json:"id"`
	SiteId       int      `json:"siteId"`
	Domain       string   `json:"domain"`
	Status       string   `json:"status,omitempty"`
	Cname        string   `json:"cname"`
	CreationDate int64    `json:"creationDate"`
	ARecords     []string `json:"aRecords"`
}

type SiteDomainsV3Meta struct {
	Page          int `json:"page"`
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

type SiteDomainsV3Response struct {
	Data   []SiteDomainV3    `json:"data"`
	Meta   SiteDomainsV3Meta `json:"meta"`
	Errors []APIErrors       `json:"errors"`
}

func (c *Client) GetSiteDomains(siteId string, accountId string) (*SiteDomainsV3Response, error) {
	if siteId == "" {
		return nil, fmt.Errorf("[ERROR] site ID was not provided")
	}

	log.Printf("[INFO] getting domains for site: %s\n", siteId)

	result := &SiteDomainsV3Response{Data: []SiteDomainV3{}}
	for pageNumber := 0; ; pageNumber++ {
		page, err := c.getSiteDomainsPage(siteId, accountId, pageNumber)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, page.Data...)
		result.Meta = page.Meta
		if len(page.Data) == 0 || pageNumber+1 >= page.Meta.TotalPages {
			break
		}
	}

	return result, nil
}

func (c *Client) getSiteDomainsPage(siteId string, accountId string, pageNumber int) (*SiteDomainsV3Response, error) {
	reqURL := fmt.Sprintf("%s%s%s/domains?excludeAutoDiscovered=false&pageSize=%d&pageNumber=%d", c.config.BaseURLAPI, endpointSiteDomainsV3, siteId, siteDomainsV3PageSize, pageNumber)
	if accountId != "" {
		reqURL = fmt.Sprintf("%s&caid=%s", reqURL, neturl.QueryEscape(accountId))
	}

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, nil, ReadSiteDomains)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domains for site %s: %s", siteId, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}
	log.Printf("[DEBUG] Incapsula get site domains JSON response: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service when reading domains for site %s: %s", resp.StatusCode, siteId, string(responseBody))
	}

	var response SiteDomainsV3Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get site domains response for site %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading domains for site %s: %s", siteId, response.Errors[0].Detail)
	}

	return &response, nil
}
