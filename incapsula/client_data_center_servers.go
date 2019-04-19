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
const endpointDataCenterServersAdd = "sites/dataCenters/servers/add"
const endpointDataCenterServersEdit = "sites/dataCenters/servers/edit"
const endpointDataCenterServersDelete = "sites/dataCenters/servers/delete"

// DataCenterServersAddResponse contains id of server
type DataCenterServersAddResponse struct {
	ServerID string `json:"server_id"`
	Res      string `json:"res"`
}

// DataCenterServersEditResponse contains data center id
type DataCenterServersEditResponse struct {
	Res          string `json:"res"`
	DataCenterID string `json:"datacenter_id"`
}

// AddDataCenterServers adds an incap data center server to be managed by Incapsula
func (c *Client) AddDataCenterServers(dcID, serverAddress, isStandby string) (*DataCenterServersAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula data center server for dcID: %s\n", dcID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServersAdd), url.Values{
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
	var dataCenterServerAddResponse DataCenterServersAddResponse
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

// EditDataCenterServers edits the Incapsula data center servers
func (c *Client) EditDataCenterServers(serverID, serverAddress, isStandby, isContent string) (*DataCenterServersEditResponse, error) {
	log.Printf("[INFO] Editing Incapsula data center servers for serverID: %s\n", serverID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServersEdit), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"server_id":      {serverID},
		"server_address": {serverAddress},
		"is_standby":     {isStandby},
		"is_content":     {isContent},
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
	var dataCenterServersEditResponse DataCenterServersEditResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterServersEditResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing edit data center server JSON response for serverID %s: %s", serverID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterServersEditResponse.Res != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when editing data center server for serverID %s: %s", serverID, string(responseBody))
	}

	return &dataCenterServersEditResponse, nil
}

// DeleteDataCenterServers deletes a data center servers currently managed by Incapsula
func (c *Client) DeleteDataCenterServers(serverID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type DataCenterServersDeleteResponse struct {
		Res      string `json:"res"`
		ServerID string `json:"server_id"`
	}

	log.Printf("[INFO] Deleting Incapsula data center server serverID: %d)\n", serverID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointDataCenterServersDelete), url.Values{
		"api_id":    {c.config.APIID},
		"api_key":   {c.config.APIKey},
		"server_id": {strconv.Itoa(serverID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting data center server (server_id: %d): %s", serverID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete data center JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var dataCenterServerDeleteResponse DataCenterServersDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &dataCenterServerDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete data center server JSON response (server_id: %d): %s", serverID, err)
	}

	// Look at the response status code from Incapsula
	if dataCenterServerDeleteResponse.Res != "0" {
		return fmt.Errorf("Error from Incapsula service when deleting data center server (server_id: %d): %s", serverID, string(responseBody))
	}

	return nil
}
