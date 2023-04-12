package incapsula

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/openpgp/errors"
	"io"
	"log"
	"net/http"
)

const endpointSiteBase = "/v2/sites"
const endpointAtoAllowlist = "/allowlist"

type AtoAllowlistItem struct {
	Ip      string `json:"ip"`
	Mask    string `json:"mask"`
	Desc    string `json:"desc"`
	Updated int64  `json:"updated"`
}

type AtoAllowlistDTO struct {
	accountId int                `json:"accountId"`
	siteId    int                `json:"siteId"`
	allowlist []AtoAllowlistItem `json:"allowlist"`
}

func (atoAllowlistDTO *AtoAllowlistDTO) toMap() (map[string]interface{}, error) {

	// Initialize the data map that terraform uses
	var atoAllowlistMap = make(map[string]interface{})

	// Set site id
	atoAllowlistMap["site_id"] = atoAllowlistDTO.siteId
	atoAllowlistMap["accountId"] = atoAllowlistDTO.accountId

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistDTO.allowlist != nil {

		atoAllowlistMap["allowlist"] = make([]map[string]interface{}, len(atoAllowlistDTO.allowlist))

		for i, allowlistItem := range atoAllowlistDTO.allowlist {

			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = map[string]interface{}{
				"Ip":      allowlistItem.Ip,
				"Mask":    allowlistItem.Mask,
				"Desc":    allowlistItem.Desc,
				"Updated": allowlistItem.Updated,
			}
			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = atoAllowlistMap["allowlist"].([]map[string]interface{})[i]
		}

	} else {
		atoAllowlistMap["allowlist"] = make([]interface{}, 0)
	}

	return atoAllowlistMap, nil
}

func formAtoAllowlistDTOFromMap(atoAllowlistMap map[string]interface{}) (*AtoAllowlistDTO, error) {

	atoAllowlistDTO := AtoAllowlistDTO{}

	// Assign site ID
	if _, err := atoAllowlistMap["site_id"].(int); err {
		return nil, errors.InvalidArgumentError("site_id should be of type int")
	}
	atoAllowlistDTO.siteId = atoAllowlistMap["site_id"].(int)

	// Assign account ID
	if _, err := atoAllowlistMap["account_id"].(int); err {
		return nil, errors.InvalidArgumentError("account_id should be of type int")
	}
	atoAllowlistDTO.accountId = atoAllowlistMap["account_id"].(int)

	// Assign the allowlist
	if atoAllowlistMap["allowlist"] == nil {
		atoAllowlistDTO.allowlist = make([]AtoAllowlistItem, 0)
	}

	if _, err := atoAllowlistMap["allowlist"].([]interface{}); err {
		return nil, errors.InvalidArgumentError("allowlist should have type array")
	}

	allowlistItems := atoAllowlistMap["allowlist"].([]map[string]interface{})
	atoAllowlistDTO.allowlist = make([]AtoAllowlistItem, len(allowlistItems))

	for i, allowlistItemMap := range allowlistItems {
		allowlistItem := AtoAllowlistItem{
			Ip:      allowlistItemMap["Ip"].(string),
			Mask:    allowlistItemMap["Mask"].(string),
			Desc:    allowlistItemMap["Desc"].(string),
			Updated: allowlistItemMap["Updated"].(int64),
		}
		atoAllowlistDTO.allowlist[i] = allowlistItem
	}

	return &atoAllowlistDTO, nil
}

func (c *Client) GetAtoSiteAllowlist(accountId, siteId int) (*AtoAllowlistDTO, error) {
	log.Printf("[INFO] Getting IP allowlist for (Site Id: %d)\n", siteId)

	// Get request to ATO
	reqURL := fmt.Sprintf("%s%s/%d%s?caid=%d", c.config.BaseURLAPI, endpointSiteBase, siteId, endpointAtoAllowlist, accountId)
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
	var atoAllowlistDTO AtoAllowlistDTO
	err = json.Unmarshal(responseBody, &atoAllowlistItems)
	atoAllowlistDTO.siteId = siteId
	atoAllowlistDTO.accountId = accountId
	atoAllowlistDTO.allowlist = atoAllowlistItems
	if err != nil {
		return nil, fmt.Errorf("[Error] Error parsing ATO allowlist response for site with ID: %d %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &atoAllowlistDTO, nil
}

func (c *Client) UpdateATOSiteAllowlist(atoSiteAllowlistDTO *AtoAllowlistDTO) error {

	log.Printf("[INFO] Updating IP allowlist for (Site Id: %d)\n", atoSiteAllowlistDTO.siteId)

	// Form the request body
	atoAllowlistJSON, err := json.Marshal(atoSiteAllowlistDTO)

	// verify site ID and account ID are not the default value for int type
	if atoSiteAllowlistDTO.siteId == 0 {
		return errors.InvalidArgumentError("site_id is not specified in updating ATO allowlist")
	}
	if atoSiteAllowlistDTO.accountId == 0 {
		return errors.InvalidArgumentError("account_id is not specified in updating ATO allowlist")
	}

	// Update request to ATO
	reqURL := fmt.Sprintf("%s%s%d%s?caid=%d", c.config.BaseURLAPI, endpointSiteBase, atoSiteAllowlistDTO.siteId, endpointAtoAllowlist, atoSiteAllowlistDTO.accountId)
	_, err = c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, atoAllowlistJSON, UpdateATOSiteAllowlistOperation)

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing update ATO allowlist request for site with id %d: %s", atoSiteAllowlistDTO.siteId, err)
	}

	return nil

}

func (c *Client) DeleteATOSiteAllowlist(accountId, siteId int) error {
	log.Printf("[INFO] Deleting IP allowlist for (Site Id: %d)\n", siteId)

	// Delete request to ATO
	reqURL := fmt.Sprintf("%s%s%d%s?caid=%d", c.config.BaseURLAPI, endpointSiteBase, siteId, endpointAtoAllowlist, accountId)
	_, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteATOSiteAllowlistOperation)

	// Handle request error
	if err != nil {
		return fmt.Errorf("[Error] Error executing update ATO allowlist request for site with id %d: %s", siteId, err)
	}

	return nil
}
