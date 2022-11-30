package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const endpointSiemConnection = "siem-config-service/v3/connections/"

type ConnectionInfo struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Path      string `json:"path"`
}

type SiemConnectionInfo struct {
	AssetID        string         `json:"assetId"`
	ConnectionName string         `json:"connectionName"`
	StorageType    string         `json:"storageType"`
	ConnectionInfo ConnectionInfo `json:"connectionInfo"`
}

type SiemConnectionWithIdAndVersionInfo struct {
	ID             string         `json:"id"`
	Version        string         `json:"version"`
	AssetID        string         `json:"assetId"`
	ConnectionName string         `json:"connectionName"`
	StorageType    string         `json:"storageType"`
	ConnectionInfo ConnectionInfo `json:"connectionInfo"`
}

type SiemConnection struct {
	Data []SiemConnectionInfo `json:"data"`
}

type SiemConnectionWithIdAndVersion struct {
	Data []SiemConnectionWithIdAndVersionInfo `json:"data"`
}

func (c *Client) CreateSiemConnection(connection *SiemConnection) (*SiemConnectionWithIdAndVersion, *int, error) {
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemConnection: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURLAPI, endpointSiemConnection)
	return siemConnectionRequestWithResponse(c, CreateSiemConnection, http.MethodPost, reqURL, connectionJSON, 201)
}

func (c *Client) ReadSiemConnection(ID string) (*SiemConnectionWithIdAndVersion, *int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	return siemConnectionRequestWithResponse(c, ReadSiemConnection, http.MethodGet, reqURL, nil, 200)
}

func (c *Client) UpdateSiemConnection(siemConnectionWithIdAndVersion *SiemConnectionWithIdAndVersion) (*SiemConnectionWithIdAndVersion, *int, error) {
	siemConnectionWithIDJSON, err := json.Marshal(siemConnectionWithIdAndVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemConnectionWithID: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, siemConnectionWithIdAndVersion.Data[0].ID)
	return siemConnectionRequestWithResponse(c, UpdateSiemConnection, http.MethodPut, reqURL, siemConnectionWithIDJSON, 200)
}

func (c *Client) DeleteSiemConnection(ID string) (*int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	_, _, responseStatusCode, err := siemConnectionRequest(c, DeleteSiemConnection, http.MethodDelete, reqURL, nil, 200)
	return responseStatusCode, err
}

func dClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func siemConnectionRequest(c *Client, operation string, method string, reqURL string, data []byte, expectedSuccessStatusCode int) (*string, *[]byte, *int, error) {

	log.Printf("[INFO] Executing operation %s on SIEM connection", operation)

	resp, err := c.DoJsonRequestWithHeaders(method, reqURL, data, operation)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error from Incapsula service when executing %s operation on SIEM connection: %s", operation, err)
	}

	defer dClose(resp.Body)
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

func siemConnectionRequestWithResponse(c *Client, operation string, method string, reqURL string, data []byte, expectedSuccessStatusCode int) (*SiemConnectionWithIdAndVersion, *int, error) {

	body, responseBody, responseStatusCode, err := siemConnectionRequest(c, operation, method, reqURL, data, expectedSuccessStatusCode)
	if responseBody == nil {
		return nil, responseStatusCode, err
	}

	var response SiemConnectionWithIdAndVersion
	err = json.Unmarshal(*responseBody, &response)
	if err != nil {
		return nil, responseStatusCode, fmt.Errorf("error obtained %s\n when constructing response for %s operation on SIEM connection from: %p",
			err, operation, body)
	}

	return &response, responseStatusCode, nil
}
