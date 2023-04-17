package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const endpointATOSiteBase = "/ato/v2/sites"
const endpointAtoAllowlist = "/allowlist"

type AtoAllowlistItem struct {
	Ip   string `json:"ip"`
	Mask string `json:"mask"`
	Desc string `json:"desc"`
}

type ATOAllowlistDTO struct {
	AccountId int                `json:"accountId"`
	SiteId    int                `json:"siteId"`
	Allowlist []AtoAllowlistItem `json:"allowlist"`
}

func formEmptyAllowlistDTO(accountId, siteId int) *ATOAllowlistDTO {
	return &ATOAllowlistDTO{
		AccountId: accountId,
		SiteId:    siteId,
		Allowlist: make([]AtoAllowlistItem, 0),
	}
}

func (atoAllowlistDTO *ATOAllowlistDTO) toMap() (map[string]interface{}, error) {

	// Initialize the data map that terraform uses
	var atoAllowlistMap = make(map[string]interface{})

	// Set site id
	atoAllowlistMap["site_id"] = atoAllowlistDTO.SiteId
	atoAllowlistMap["accountId"] = atoAllowlistDTO.AccountId

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistDTO.Allowlist != nil {

		atoAllowlistMap["allowlist"] = make([]map[string]interface{}, len(atoAllowlistDTO.Allowlist))

		for i, allowlistItem := range atoAllowlistDTO.Allowlist {

			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = map[string]interface{}{
				"ip":   allowlistItem.Ip,
				"mask": allowlistItem.Mask,
				"desc": allowlistItem.Desc,
			}
			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = atoAllowlistMap["allowlist"].([]map[string]interface{})[i]
		}

	} else {
		atoAllowlistMap["allowlist"] = make([]interface{}, 0)
	}

	return atoAllowlistMap, nil
}

func formAtoAllowlistDTOFromMap(atoAllowlistMap map[string]interface{}) (*ATOAllowlistDTO, error) {

	atoAllowlistDTO := ATOAllowlistDTO{}

	// Validate site_id
	switch atoAllowlistMap["site_id"].(type) {
	case int:
		break
	default:
		return nil, fmt.Errorf("site_id should be of type int")
	}

	// validate account_id
	switch atoAllowlistMap["account_id"].(type) {
	case int:
		break
	default:
		return nil, fmt.Errorf("account_id should be of type int")
	}

	// Assign site ID
	atoAllowlistDTO.SiteId = atoAllowlistMap["site_id"].(int)

	// Assign account ID
	atoAllowlistDTO.AccountId = atoAllowlistMap["account_id"].(int)

	// Assign the allowlist
	if atoAllowlistMap["allowlist"] == nil {
		atoAllowlistDTO.Allowlist = make([]AtoAllowlistItem, 0)
	}

	// Verify that the allowlist is an array
	if _, ok := atoAllowlistMap["allowlist"].([]interface{}); !ok {
		return nil, fmt.Errorf("allowlist should have type array")
	}

	allowlistItemsInMap := atoAllowlistMap["allowlist"].([]interface{})
	atoAllowlistDTO.Allowlist = make([]AtoAllowlistItem, len(allowlistItemsInMap))

	// Convert each allowlist entry in the map to the allowlist item for the DTO
	for i, allowlistItemMap := range allowlistItemsInMap {
		allowListItemMap := allowlistItemMap.(map[string]interface{})

		// Initialize allowlist item
		allowlistItem := AtoAllowlistItem{}

		// Check that IP is not empty
		if allowListItemMap["ip"] == nil {
			return nil, fmt.Errorf("IP cannot be empty in allowlist items")
		}

		allowlistItem.Ip = allowListItemMap["ip"].(string)

		// Extract description
		if allowListItemMap["desc"] != nil {
			allowlistItem.Desc = allowListItemMap["desc"].(string)
		}

		// Extract subnet from map
		if allowListItemMap["mask"] != nil {
			allowlistItem.Mask = allowListItemMap["mask"].(string)
		}

		atoAllowlistDTO.Allowlist[i] = allowlistItem
	}

	return &atoAllowlistDTO, nil
}

func (c *Client) GetAtoSiteAllowlistWithRetries(accountId, siteId int) (*ATOAllowlistDTO, error) {
	// Since the newly created site can take upto 30 seconds to be fully configured, we per.si a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		atoAllowlistDTO, err := c.GetAtoSiteAllowlist(accountId, siteId)
		if err == nil {
			return atoAllowlistDTO, nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return nil, lastError
}

func (c *Client) GetAtoSiteAllowlist(accountId, siteId int) (*ATOAllowlistDTO, error) {
	log.Printf("[INFO] Getting IP allowlist for (Site Id: %d)\n", siteId)

	// Get request to ATO
	var reqURL string
	if accountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointAtoAllowlist)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, siteId, endpointAtoAllowlist, accountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadATOSiteAllowlistOperation)
	if err != nil {
		return nil, fmt.Errorf("[Error] Error executing get ATO allowlist request for site with id %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] ATO allowlist JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var atoAllowlistItems []AtoAllowlistItem
	var atoAllowlistDTO ATOAllowlistDTO
	err = json.Unmarshal(responseBody, &atoAllowlistItems)
	atoAllowlistDTO.SiteId = siteId
	atoAllowlistDTO.AccountId = accountId
	atoAllowlistDTO.Allowlist = atoAllowlistItems
	if err != nil {
		return nil, fmt.Errorf("[Error] Q Error parsing ATO allowlist response for site with ID: %d %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &atoAllowlistDTO, nil
}

func (c *Client) UpdateATOSiteAllowlistWithRetries(atoSiteAllowlistDTO *ATOAllowlistDTO) error {
	// Since the newly created site can take upto 30 seconds to be fully configured, we perform a simple backoff
	var backoffSchedule = []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
	var lastError error

	for _, backoff := range backoffSchedule {
		err := c.UpdateATOSiteAllowlist(atoSiteAllowlistDTO)
		if err == nil {
			return nil
		}
		lastError = err
		time.Sleep(backoff)
	}
	return lastError
}

func (c *Client) UpdateATOSiteAllowlist(atoSiteAllowlistDTO *ATOAllowlistDTO) error {

	log.Printf("[INFO] Updating ATO IP allowlist for (Site Id: %d)\n", atoSiteAllowlistDTO.SiteId)

	// Form the request body
	atoAllowlistJSON, err := json.Marshal(atoSiteAllowlistDTO.Allowlist)

	// verify site ID and account ID are not the default value for int type
	if atoSiteAllowlistDTO.SiteId == 0 {
		return fmt.Errorf("site_id is not specified in updating ATO allowlist")
	}
	var reqURL string
	if atoSiteAllowlistDTO.AccountId == 0 {
		reqURL = fmt.Sprintf("%s%s/%d%s", c.config.BaseURLAPI, endpointATOSiteBase, atoSiteAllowlistDTO.SiteId, endpointAtoAllowlist)
	} else {
		reqURL = fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointATOSiteBase, atoSiteAllowlistDTO.SiteId, endpointAtoAllowlist, atoSiteAllowlistDTO.AccountId)
	}

	// Update request to ATO
	response, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, atoAllowlistJSON, UpdateATOSiteAllowlistOperation)

	// Read the body
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)

	log.Printf("Updated ATO allowlist with response : %s", responseBody)

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing update ATO allowlist request for site with id %d: %s", atoSiteAllowlistDTO.SiteId, err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("[Error] Error executing update ATO allowlist request for site with status %d: %s", atoSiteAllowlistDTO.SiteId, response.Status)
	}

	return nil

}

func (c *Client) DeleteATOSiteAllowlist(accountId, siteId int) error {
	log.Printf("[INFO] Deleting IP allowlist for (Site Id: %d)\n", siteId)

	err := c.UpdateATOSiteAllowlist(formEmptyAllowlistDTO(accountId, siteId))

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing delete ATO allowlist request for site with id %d: %s", siteId, err)
	}

	return nil
}
