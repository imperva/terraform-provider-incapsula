package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

const endpointAccountDataStorageRegionGet = "accounts/data-privacy/show"
const endpointAccountDataStorageRegionUpdate = "accounts/data-privacy/set-region-default"

// AccountDataStorageRegionResponse contains the relevant information when getting/setting a default data storage region
type AccountDataStorageRegionResponse struct {
	Region     string `json:"region"`
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	DebugInfo  struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// GetAccountDataStorageRegion gets the default data storage region for sites in the account
func (c *Client) GetAccountDataStorageRegion(accountID string) (*AccountDataStorageRegionResponse, error) {
	log.Printf("[INFO] Getting default Incapsula data storage region for account: %s\n", accountID)

	// Post form to Incapsula
	values := url.Values{"account_id": {accountID}}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountDataStorageRegionGet)
	resp, err := c.PostFormWithHeaders(reqURL, values, ReadAccountDataStorageRegion)
	if err != nil {
		return nil, fmt.Errorf("Error getting default data storage region for account id: %s: %s", accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula default data storage region JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountDataStorageRegionResponse AccountDataStorageRegionResponse
	err = json.Unmarshal([]byte(responseBody), &accountDataStorageRegionResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing default data storage region JSON response for account id: %s: %s", accountID, err)
	}

	// Look at the response status code from Incapsula
	if accountDataStorageRegionResponse.Res != 0 {
		return &accountDataStorageRegionResponse, fmt.Errorf("Error from Incapsula service when getting default data storage region for account id: %s: %s", accountID, string(responseBody))
	}

	return &accountDataStorageRegionResponse, nil
}

// UpdateAccountDataStorageRegion will update the default data storage region on the account
func (c *Client) UpdateAccountDataStorageRegion(accountID, region string) (*AccountDataStorageRegionResponse, error) {
	log.Printf("[INFO] Updating Incapsula default data storage region (%s) for accountID: %s\n", region, accountID)

	// Post form to Incapsula
	values := url.Values{
		"account_id":          {accountID},
		"data_storage_region": {region},
	}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountDataStorageRegionUpdate)
	resp, err := c.PostFormWithHeaders(reqURL, values, UpdateAccountDataStorageRegion)
	if err != nil {
		return nil, fmt.Errorf("Error updating data storage region with value (%s) on account_id: %s: %s", region, accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update account default data storage region JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountDataStorageRegionResponse AccountDataStorageRegionResponse
	err = json.Unmarshal([]byte(responseBody), &accountDataStorageRegionResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update default data storage region JSON response for accountID %s: %s", accountID, err)
	}

	// Look at the response status code from Incapsula
	if accountDataStorageRegionResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating default data storage region for accountID %s: %s", accountID, string(responseBody))
	}

	return &accountDataStorageRegionResponse, nil
}
