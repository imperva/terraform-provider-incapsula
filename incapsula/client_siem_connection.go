package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const endpointSiemConnection = "siem-config-service/v3/connections"

type ConnectionInfo struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
	Path      string `json:"path"`
}

type SiemConnectionData struct {
	ID             string         `json:"id,omitempty"`
	Version        string         `json:"version,omitempty"`
	AssetID        string         `json:"assetId,omitempty"`
	ConnectionName string         `json:"connectionName"`
	StorageType    string         `json:"storageType"`
	ConnectionInfo ConnectionInfo `json:"connectionInfo"`
}

type SiemConnection struct {
	Data []SiemConnectionData `json:"data"`
}

func (c *Client) CreateSiemConnection(connection *SiemConnection) (*SiemConnection, *int, error) {
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemConnection: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/", c.config.BaseURLAPI, endpointSiemConnection)
	return siemConnectionRequestWithResponse(c, CreateSiemConnection, http.MethodPost, reqURL, connectionJSON, connection.Data[0].AssetID, 201)
}

func (c *Client) ReadSiemConnection(ID string) (*SiemConnection, *int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	return siemConnectionRequestWithResponse(c, ReadSiemConnection, http.MethodGet, reqURL, nil, "", 200)
}

func (c *Client) UpdateSiemConnection(siemConnection *SiemConnection) (*SiemConnection, *int, error) {
	siemConnectionJSON, err := json.Marshal(siemConnection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemConnectionWithID: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, siemConnection.Data[0].ID)
	return siemConnectionRequestWithResponse(c, UpdateSiemConnection, http.MethodPut, reqURL, siemConnectionJSON, siemConnection.Data[0].AssetID, 200)
}

func (c *Client) DeleteSiemConnection(ID string) (*int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	_, _, statusCode, err := siemConnectionRequest(c, DeleteSiemConnection, http.MethodDelete, reqURL, nil, "", 200)
	return statusCode, err
}

func dSiemConnectionResponseClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func siemConnectionRequest(c *Client, operation string, method string, reqURL string, data []byte, accountIdStr string, expectedSuccessStatusCode int) (*string, *[]byte, *int, error) {
	log.Printf("[INFO] Executing operation %s on SIEM connection with data: %s", operation, data)

	var params = map[string]string{}
	accountId, err := strconv.Atoi(accountIdStr)
	if err == nil && accountId > 0 {
		params["caid"] = accountIdStr
	}

	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(method, reqURL, data, params, operation)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error from Incapsula service when executing %s operation on SIEM connection: %s", operation, err)
	}

	defer dSiemConnectionResponseClose(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	body := string(responseBody)

	if err != nil {
		return nil, nil, &resp.StatusCode, fmt.Errorf("error occurred: %s\n when reading response from body: %s", err, body)
	}
	log.Printf("[DEBUG] Incapsula returned response: %s\nfor %s operation on SIEM connection", body, operation)

	if resp.StatusCode != expectedSuccessStatusCode {
		return nil, nil, &resp.StatusCode, fmt.Errorf("received failure response for operation: %s on SIEM connection\nstatus code: %d\nbody: %s",
			operation, resp.StatusCode, body)
	}

	return &body, &responseBody, &resp.StatusCode, nil
}

func siemConnectionRequestWithResponse(c *Client, operation string, method string, reqURL string, data []byte, accountIdStr string, expectedSuccessStatusCode int) (*SiemConnection, *int, error) {
	body, responseBody, responseStatusCode, err := siemConnectionRequest(c, operation, method, reqURL, data, accountIdStr, expectedSuccessStatusCode)
	if responseBody == nil {
		return nil, responseStatusCode, err
	}

	var response SiemConnection
	err = json.Unmarshal(*responseBody, &response)
	if err != nil {
		return nil, responseStatusCode, fmt.Errorf("error obtained %s\n when constructing response for %s operation on SIEM connection from: %p",
			err, operation, body)
	}

	return &response, responseStatusCode, nil
}
