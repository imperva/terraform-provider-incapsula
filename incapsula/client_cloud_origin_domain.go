package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CloudOriginDomainConfig struct {
	Port            int    `json:"port,omitempty"`
	OriginTlsPolicy string `json:"originTlsPolicy,omitempty"`
}

type CloudOriginDomainData struct {
	ID                  int                     `json:"id"`
	SiteID              int                     `json:"siteId"`
	OriginDomain        string                  `json:"originDomain"`
	Region              string                  `json:"region"`
	ImpervaOriginDomain string                  `json:"impervaOriginDomain"`
	OriginConfig        CloudOriginDomainConfig `json:"originConfig"`
	CreatedAt           string                  `json:"createdAt"`
	UpdatedAt           string                  `json:"updatedAt"`
}

type CloudOriginDomainResponse struct {
	Data   []CloudOriginDomainData `json:"data"`
	Errors []APIErrors             `json:"errors"`
}

type CloudOriginDomainCreateRequest struct {
	OriginDomain string                   `json:"originDomain"`
	Region       string                   `json:"region"`
	DomainConfig *CloudOriginDomainConfig `json:"domainConfig,omitempty"`
}

type CloudOriginDomainUpdateRequest struct {
	Region       string                   `json:"region,omitempty"`
	DomainConfig *CloudOriginDomainConfig `json:"domainConfig,omitempty"`
}

func getCloudOriginUrl(baseURL string, siteID int, path string, accountID string) string {
	url := fmt.Sprintf("%s/sites/%d/cloud-origins%s", baseURL, siteID, path)
	if accountID != "" {
		url = fmt.Sprintf("%s?caid=%s", url, accountID)
	}
	return url
}

func (c *Client) CreateCloudOriginDomain(siteID int, accountID string, domain, region string, port int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Creating Incapsula cloud origin domain: %s for site: %d\n", domain, siteID)

	payload := CloudOriginDomainCreateRequest{
		OriginDomain: domain,
		Region:       region,
		DomainConfig: &CloudOriginDomainConfig{
			Port: port,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal cloud origin domain: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost,
		getCloudOriginUrl(c.config.BaseURLRev3, siteID, "", accountID),
		payloadJSON,
		CreateCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while creating cloud origin domain %s for site %d: %s", domain, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula create cloud origin domain JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when creating cloud origin domain %s for site %d: %s", resp.StatusCode, domain, siteID, string(responseBody))
	}

	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error from Incapsula service when creating cloud origin domain %s for site %d: %s", domain, siteID, response.Errors[0].Detail)
	}

	return &response, nil
}

func (c *Client) GetCloudOriginDomain(siteID, originID int, accountID string) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Getting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		getCloudOriginUrl(c.config.BaseURLRev3, siteID, fmt.Sprintf("/%d", originID), accountID),
		nil,
		ReadCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while reading cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula get cloud origin domain JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error from Incapsula service when reading cloud origin domain %d for site %d: %s", originID, siteID, response.Errors[0].Detail)
	}

	return &response, nil
}

func (c *Client) UpdateCloudOriginDomain(siteID, originID int, accountID string, region string, port int) (*CloudOriginDomainResponse, error) {
	log.Printf("[INFO] Updating Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	payload := CloudOriginDomainUpdateRequest{
		Region: region,
		DomainConfig: &CloudOriginDomainConfig{
			Port: port,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal cloud origin domain update: %s", err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut,
		getCloudOriginUrl(c.config.BaseURLRev3, siteID, fmt.Sprintf("/%d", originID), accountID),
		payloadJSON,
		UpdateCloudOriginDomain)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while updating cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula update cloud origin domain JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	var response CloudOriginDomainResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cloud origin domain JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating cloud origin domain %d for site %d: %s", originID, siteID, response.Errors[0].Detail)
	}

	return &response, nil
}

func (c *Client) DeleteCloudOriginDomain(siteID, originID int, accountID string) error {
	log.Printf("[INFO] Deleting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete,
		getCloudOriginUrl(c.config.BaseURLRev3, siteID, fmt.Sprintf("/%d", originID), accountID),
		nil,
		DeleteCloudOriginDomain)

	if err != nil {
		return fmt.Errorf("Error from Incapsula service while deleting cloud origin domain %d for site %d: %s", originID, siteID, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula delete cloud origin domain JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting cloud origin domain %d for site %d: %s", resp.StatusCode, originID, siteID, string(responseBody))
	}

	return nil
}
