package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

// Endpoints (unexported consts)
const endpointDataCenterAdd = "sites/dataCenters/add"
const endpointDataCenterList = "sites/dataCenters/list"
const endpointDataCenterEdit = "sites/dataCenters/edit"
const endpointDataCenterDelete = "sites/dataCenters/delete"

// todo: get data center responses
// DataCenterAddResponse contains todo
type DataCenterAddResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// DataCenterListResponse contains todo
type DataCenterListResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// DataCenterEditResponse contains todo
type DataCenterEditResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// DataCenterDeleteResponse contains todo
type DataCenterDeleteResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

// AddDataCenter adds an incap rule to be managed by Incapsula
func (c *Client) AddDataCenter(siteID int, name string, serverAddress string, isStandby string, isContent string) (*DataCenterAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula data center for siteID: %d\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterAdd), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"site_id":        {strconv.Itoa(siteID)},
		"name":           {name},
		"server_address": {serverAddress},
		"is_standby":     {isStandby},
		"is_content":     {isContent},
	})
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center for siteID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterAddResponse DataCenterAddResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add data center JSON response for siteID %d: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center for siteID %d: %s", siteID, string(responseBody))
	}

	return &dataCenterAddResponse, nil
}

// DataCenterList gets the Incapsula list of incap rules
func (c *Client) ListDataCenters(siteID int) (*DataCenterListResponse, error) {
	log.Printf("[INFO] Getting Incapsula data centers (site_id: %d)\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterList), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting data centers (site_id: %d): %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula data centers JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var incapRuleListResponse DataCenterListResponse
	err = json.Unmarshal([]byte(responseBody), &incapRuleListResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing data centers list JSON response (site_id: %d): %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if incapRuleListResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting data centers list (site_id: %d): %s", siteID, string(responseBody))
	}

	return &incapRuleListResponse, nil
}

// EditDataCenter edits the Incapsula incap rule
func (c *Client) EditDataCenter(dcID int, name, isStandby, isContent string) (*DataCenterEditResponse, error) {
	log.Printf("[INFO] Editing Incapsula data center for dcID: %d\n", name, dcID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterEdit), url.Values{
		"api_id":     {c.config.APIID},
		"api_key":    {c.config.APIKey},
		"dc_id":      {strconv.Itoa(dcID)},
		"name":       {name},
		"is_standby": {isStandby},
		"is_content": {isContent},
	})
	if err != nil {
		return nil, fmt.Errorf("Error editing data center  for dcID: %d: %s", dcID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula edit data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterEditResponse DataCenterEditResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterEditResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing edit dta center JSON response for dcID %d: %s", dcID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterEditResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when editing data center for dcID %d: %s", dcID, string(responseBody))
	}

	return &dataCenterEditResponse, nil
}

// DeleteDataCenter deletes a site currently managed by Incapsula
func (c *Client) DeleteDataCenter(dcID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type DataCenterDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula data center id: %d)\n", dcID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"dc_id":   {strconv.Itoa(dcID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting data center (dc_id: %d): %s", dcID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterDeleteResponse DataCenterDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete data center JSON response (dc_id: %d): %s", dcID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting data center (dc_id: %d): %s", dcID, string(responseBody))
	}

	return nil
}
