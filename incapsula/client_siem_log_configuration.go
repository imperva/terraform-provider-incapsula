package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const endpointSiemLogConfiguration = "siem-config-service/v3/log-configurations/"

type SiemLogConfigurationInfo struct {
	AssetID           string   `json:"assetId"`
	ConfigurationName string   `json:"configurationName"`
	Provider          string   `json:"provider"`
	Datasets          []string `json:"datasets"`
	Enabled           bool     `json:"enabled"`
	ConnectionId      string   `json:"connectionId"`
}

type SiemLogConfigurationWithIdAndVersionInfo struct {
	ID                string   `json:"id"`
	Version           string   `json:"version"`
	AssetID           string   `json:"assetId"`
	ConfigurationName string   `json:"configurationName"`
	Provider          string   `json:"provider"`
	Datasets          []string `json:"datasets"`
	Enabled           bool     `json:"enabled"`
	ConnectionId      string   `json:"connectionId"`
}

type SiemLogConfiguration struct {
	Data []SiemLogConfigurationInfo `json:"data"`
}

type SiemLogConfigurationWithIdAndVersion struct {
	Data []SiemLogConfigurationWithIdAndVersionInfo `json:"data"`
}

func (c *Client) CreateSiemLogConfiguration(siemLogConfiguration *SiemLogConfiguration) (*SiemLogConfigurationWithIdAndVersion, *int, error) {
	logConfigurationJSON, err := json.Marshal(siemLogConfiguration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemLogConfiguration: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration)
	accountId := 0
	if len(siemLogConfiguration.Data[0].AssetID) > 0 {
		accountId, _ = strconv.Atoi(siemLogConfiguration.Data[0].AssetID)
	}
	return siemLogConfigurationRequestWithResponse(c, CreateSiemLogConfiguration, http.MethodPost, reqURL, logConfigurationJSON, accountId, 201)
}

func (c *Client) ReadSiemLogConfiguration(ID string) (*SiemLogConfigurationWithIdAndVersion, *int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, ID)
	return siemLogConfigurationRequestWithResponse(c, ReadSiemLogConfiguration, http.MethodGet, reqURL, nil, 0, 200)
}

func (c *Client) UpdateSiemLogConfiguration(logConfigurationWithIdAndVersion *SiemLogConfigurationWithIdAndVersion) (*SiemLogConfigurationWithIdAndVersion, *int, error) {
	siemLogConfigurationWithIDJSON, err := json.Marshal(logConfigurationWithIdAndVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemLogConfigurationWithID: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, logConfigurationWithIdAndVersion.Data[0].ID)
	accountId := 0
	if len(logConfigurationWithIdAndVersion.Data[0].AssetID) > 0 {
		accountId, _ = strconv.Atoi(logConfigurationWithIdAndVersion.Data[0].AssetID)
	}
	return siemLogConfigurationRequestWithResponse(c, UpdateSiemLogConfiguration, http.MethodPut, reqURL, siemLogConfigurationWithIDJSON, accountId, 200)
}

func (c *Client) DeleteSiemLogConfiguration(ID string) (*int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, ID)
	_, _, responseStatusCode, err := siemLogConfigurationRequest(c, DeleteSiemLogConfiguration, http.MethodDelete, reqURL, nil, 0, 200)
	return responseStatusCode, err
}

func dSiemLogConfigurationResponseClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func siemLogConfigurationRequest(c *Client, operation string, method string, reqURL string, data []byte, accountId int, expectedSuccessStatusCode int) (*string, *[]byte, *int, error) {

	log.Printf("[INFO] Executing operation %s on SIEM log configuration", operation)

	params := GetRequestParamsWithCaid(accountId)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(method, reqURL, data, params, operation)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error from Incapsula service when executing %s operation on SIEM log configuration: %s", operation, err)
	}

	defer dSiemLogConfigurationResponseClose(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	body := string(responseBody)

	if err != nil {
		return nil, nil, &resp.StatusCode, fmt.Errorf("error occurred: %s\n when reading response from body: %s", err, body)
	}
	log.Printf("[DEBUG] Incapsula returned response: %s\nfor %s operation on SIEM log configuration", body, operation)

	if resp.StatusCode != expectedSuccessStatusCode {
		return nil, nil, &resp.StatusCode, fmt.Errorf("received failure response for operation: %s on SIEM log configuration\nstatus code: %d\nbody: %s",
			operation, resp.StatusCode, body)
	}

	return &body, &responseBody, &resp.StatusCode, nil
}

func siemLogConfigurationRequestWithResponse(c *Client, operation string, method string, reqURL string, data []byte, accountId int, expectedSuccessStatusCode int) (*SiemLogConfigurationWithIdAndVersion, *int, error) {

	body, responseBody, responseStatusCode, err := siemLogConfigurationRequest(c, operation, method, reqURL, data, accountId, expectedSuccessStatusCode)
	if responseBody == nil {
		return nil, responseStatusCode, err
	}

	var response SiemLogConfigurationWithIdAndVersion
	err = json.Unmarshal(*responseBody, &response)
	if err != nil {
		return nil, responseStatusCode, fmt.Errorf("error obtained %s\n when constructing response for %s operation on SIEM log configuration from: %p",
			err, operation, body)
	}

	return &response, responseStatusCode, nil
}
