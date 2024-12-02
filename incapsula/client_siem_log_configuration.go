package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const endpointSiemLogConfiguration = "siem-config-service/v3/log-configurations"

type SiemLogConfigurationData struct {
	ID                string        `json:"id,omitempty"`
	AssetID           string        `json:"assetId,omitempty"`
	ConfigurationName string        `json:"configurationName"`
	Provider          string        `json:"provider"`
	Datasets          []interface{} `json:"datasets"`
	Enabled           bool          `json:"enabled"`
	ConnectionId      string        `json:"connectionId"`
	CompressLogs      bool          `json:"compressLogs,omitempty"`
	Format            string        `json:"format,omitempty"`
	LogsLevel         string        `json:"logsLevel,omitempty"`
	PublicKey         string        `json:"publicKey,omitempty"`
	PublicKeyFileNAme string        `json:"publicKeyFileName,omitempty"`
}

type SiemLogConfiguration struct {
	Data []SiemLogConfigurationData `json:"data"`
}

func (c *Client) CreateSiemLogConfiguration(siemLogConfiguration *SiemLogConfiguration) (*SiemLogConfiguration, *int, error) {
	logConfigurationJSON, err := json.Marshal(siemLogConfiguration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemLogConfiguration: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/", c.config.BaseURLAPI, endpointSiemLogConfiguration)
	return siemLogConfigurationRequestWithResponse(c, CreateSiemLogConfiguration, http.MethodPost, reqURL, logConfigurationJSON, siemLogConfiguration.Data[0].AssetID, 201)
}

func (c *Client) ReadSiemLogConfiguration(ID string, accountId string) (*SiemLogConfiguration, *int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, ID)
	return siemLogConfigurationRequestWithResponse(c, ReadSiemLogConfiguration, http.MethodGet, reqURL, nil, accountId, 200)
}

func (c *Client) UpdateSiemLogConfiguration(siemLogConfiguration *SiemLogConfiguration) (*SiemLogConfiguration, *int, error) {
	siemLogConfigurationJSON, err := json.Marshal(siemLogConfiguration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemLogConfigurationWithID: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, siemLogConfiguration.Data[0].ID)
	return siemLogConfigurationRequestWithResponse(c, UpdateSiemLogConfiguration, http.MethodPut, reqURL, siemLogConfigurationJSON, siemLogConfiguration.Data[0].AssetID, 200)
}

func (c *Client) DeleteSiemLogConfiguration(ID string, accountId string) (*int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemLogConfiguration, ID)
	_, _, responseStatusCode, err := siemLogConfigurationRequest(c, DeleteSiemLogConfiguration, http.MethodDelete, reqURL, nil, accountId, 200)
	return responseStatusCode, err
}

func dSiemLogConfigurationResponseClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func siemLogConfigurationRequest(c *Client, operation string, method string, reqURL string, data []byte, accountIdStr string, expectedSuccessStatusCode int) (*string, *[]byte, *int, error) {
	log.Printf("[INFO] Executing operation %s on SIEM log configuration with data: %s", operation, data)

	var params = map[string]string{}
	accountId, err := strconv.Atoi(accountIdStr)
	if err == nil && accountId > 0 {
		params["caid"] = accountIdStr
	}

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

func siemLogConfigurationRequestWithResponse(c *Client, operation string, method string, reqURL string, data []byte, accountIdStr string, expectedSuccessStatusCode int) (*SiemLogConfiguration, *int, error) {
	body, responseBody, responseStatusCode, err := siemLogConfigurationRequest(c, operation, method, reqURL, data, accountIdStr, expectedSuccessStatusCode)
	if responseBody == nil {
		return nil, responseStatusCode, err
	}

	var response SiemLogConfiguration
	err = json.Unmarshal(*responseBody, &response)
	if err != nil {
		return nil, responseStatusCode, fmt.Errorf("error obtained %s\n when constructing response for %s operation on SIEM log configuration from: %p",
			err, operation, body)
	}

	return &response, responseStatusCode, nil
}
