package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

// Endpoints (unexported consts)
const endpointDataCenterServerAdd = "sites/dataCenters/servers/add"
const endpointDataCenterServerEdit = "sites/dataCenters/servers/edit"
const endpointDataCenterServerDelete = "sites/dataCenters/servers/delete"

// DataCenterServerAddResponse contains id of server
type DataCenterServerAddResponse struct {
	ServerID string `json:"server_id"`
	Res      string `json:"res"`
}

// DataCenterServerEditResponse contains data center id
type DataCenterServerEditResponse struct {
	Res          string `json:"res"`
	DataCenterID string `json:"datacenter_id"`
}

// AddDataCenterServer adds an incap data center server to be managed by Incapsula
func (c *Client) AddDataCenterServer(dcID, serverAddress, isStandby string) (*DataCenterServerAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula data center server for dcID: %s\n", dcID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServerAdd), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"dc_id":          {dcID},
		"server_address": {serverAddress},
		"is_standby":     {isStandby},
	})
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center server for dcID %s: %s", dcID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterServerAddResponse DataCenterServerAddResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterServerAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add data center server JSON response for dcID %s: %s\nresponse: %s", dcID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if dataCenterServerAddResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when adding data center server for dcID %s: %s", dcID, string(responseBody))
	}

	return &dataCenterServerAddResponse, nil
}

// EditDataCenterServer edits the Incapsula data center server
func (c *Client) EditDataCenterServer(serverID, serverAddress, isStandby, isEnabled string) (*DataCenterServerEditResponse, error) {
	log.Printf("[INFO] Editing Incapsula data center server for serverID: %s\n", serverID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServerEdit), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"server_id":      {serverID},
		"server_address": {serverAddress},
		"is_standby":     {isStandby},
		"is_enabled":     {isEnabled},
	})
	if err != nil {
		return nil, fmt.Errorf("Error editing data center server for serverID: %s: %s", serverID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula edit data center server JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterServerEditResponse DataCenterServerEditResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterServerEditResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing edit data center server JSON response for serverID %s: %s", serverID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterServerEditResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when editing data center server for serverID %s: %s", serverID, string(responseBody))
	}

	return &dataCenterServerEditResponse, nil
}

// DeleteDataCenterServer deletes a data center server currently managed by Incapsula
func (c *Client) DeleteDataCenterServer(serverID string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type DataCenterServerDeleteResponse struct {
		Res      string `json:"res"`
		ServerID string `json:"server_id"`
	}

	log.Printf("[INFO] Deleting Incapsula data center server serverID: %s)\n", serverID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServerDelete), url.Values{
		"api_id":    {c.config.APIID},
		"api_key":   {c.config.APIKey},
		"server_id": {serverID},
	})
	if err != nil {
		return fmt.Errorf("Error deleting data center server (server_id: %s): %s", serverID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterServerDeleteResponse DataCenterServerDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterServerDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete data center server JSON response (server_id: %s): %s", serverID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterServerDeleteResponse.Res != "0" {
		return fmt.Errorf("Error from Incapsula service when deleting data center server (server_id: %s): %s", serverID, string(responseBody))
	}

	return nil
}
