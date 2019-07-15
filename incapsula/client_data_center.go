package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

// Endpoints (unexported consts)
const endpointDataCenterAdd = "sites/dataCenters/add"
const endpointDataCenterList = "sites/dataCenters/list"
const endpointDataCenterEdit = "sites/dataCenters/edit"
const endpointDataCenterDelete = "sites/dataCenters/delete"

// DataCenterAddResponse contains id of data center
type DataCenterAddResponse struct {
	Res          string `json:"res"`
	DataCenterID string `json:"datacenter_id"`
}

// DataCenterListResponse contains list of data centers and servers
type DataCenterListResponse struct {
	Res string `json:"res"`
	DCs []struct {
		ID      string `json:"id"`
		Enabled string `json:"enabled"`
		Servers []struct {
			ID        string `json:"id"`
			Enabled   string `json:"enabled"`
			Address   string `json:"address"`
			IsStandBy string `json:"isStandby"`
		} `json:"servers"`
		Name        string `json:"name"`
		ContentOnly string `json:"contentOnly"`
		IsActive    string `json:"isActive"`
	} `json:"DCs"`
}

// DataCenterEditResponse contains edit response message
type DataCenterEditResponse struct {
	Res        string `json:"res"`
	ResMessage string `json:"res_message"`
}

// AddDataCenter adds an incap rule to be managed by Incapsula
func (c *Client) AddDataCenter(siteID, name, serverAddress, isStandby, isContent string) (*DataCenterAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula data center for siteID: %s\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterAdd), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"site_id":        {siteID},
		"name":           {name},
		"server_address": {serverAddress},
		"is_standby":     {isStandby},
		"is_content":     {isContent},
	})
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center for siteID %s: %s", siteID, err)
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
		return nil, fmt.Errorf("Error parsing add data center JSON response for siteID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if dataCenterAddResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center for siteID %s: %s", siteID, string(responseBody))
	}

	return &dataCenterAddResponse, nil
}

// ListDataCenters gets the Incapsula list of data centers
func (c *Client) ListDataCenters(siteID string) (*DataCenterListResponse, error) {
	log.Printf("[INFO] Getting Incapsula data centers (site_id: %s)\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterList), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {siteID},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting data centers for siteID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula data centers JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterListResponse DataCenterListResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterListResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing data centers list JSON response for siteID: %s %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if dataCenterListResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when getting data centers list (site_id: %s): %s", siteID, string(responseBody))
	}

	return &dataCenterListResponse, nil
}

// EditDataCenter edits the Incapsula incap rule
func (c *Client) EditDataCenter(dcID, name, isStandby, isContent, isEnabled string) (*DataCenterEditResponse, error) {
	log.Printf("[INFO] Editing Incapsula data center for dcID: %s\n", dcID)

	values := url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"dc_id":   {dcID},
	}

	if name != "" {
		values.Add("name", name)
	}

	if isStandby != "" {
		values.Add("is_standby", isStandby)
	}

	if isContent != "" {
		values.Add("is_content", isContent)
	}

	if isEnabled != "" {
		values.Add("is_enabled", isEnabled)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterEdit), values)
	if err != nil {
		return nil, fmt.Errorf("Error editing data center  for dcID: %s: %s", dcID, err)
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
		return nil, fmt.Errorf("Error parsing edit dta center JSON response for dcID %s: %s", dcID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterEditResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when editing data center for dcID %s: %s", dcID, string(responseBody))
	}

	return &dataCenterEditResponse, nil
}

// DeleteDataCenter deletes a site currently managed by Incapsula
func (c *Client) DeleteDataCenter(dcID string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type DataCenterDeleteResponse struct {
		Res        interface{} `json:"res"`
		ResMessage string      `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula data center id: %s)\n", dcID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"dc_id":   {dcID},
	})
	if err != nil {
		return fmt.Errorf("Error deleting data center (dc_id: %s): %s", dcID, err)
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
		return fmt.Errorf("Error parsing delete data center JSON response (dc_id: %s): %s", dcID, err)
	}

	// Res can sometimes oscillate between a string and number
	// We need to add safeguards for this inside the provider
	var resString string

	if resNumber, ok := dataCenterDeleteResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = dataCenterDeleteResponse.Res.(string)
	}

	// Look at the response status code from Incapsula data center
	if resString == "0" || resString == "9413" {
		return nil
	}

	return fmt.Errorf("Error from Incapsula service when deleting data center (dc_id: %s): %s", dcID, string(responseBody))
}
