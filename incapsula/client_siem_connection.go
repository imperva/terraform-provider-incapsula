package incapsula

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const endpointSiemConnection = "siem-config-service/v3/connections"

type S3ConnectionInfo struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
	Path      string `json:"path"`
}

type SplunkConnectionInfo struct {
	Host                    string `json:"host"`
	Port                    int    `json:"port"`
	Token                   string `json:"token,omitempty"`
	DisableCertVerification bool   `json:"disableCertVerification"`
}

type SftpConnectionInfo struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

const StorageTypeCustomerHttp = "CUSTOMER_HTTP"

const HttpAuthTypeHeaders = "HEADERS"
const HttpAuthTypeBasicAuth = "BASIC_AUTH"

// AuthHeaderInfo is one entry of a HEADERS authentication.
type AuthHeaderInfo struct {
	HeaderName   string `json:"headerName"`
	SecretFormat string `json:"secretFormat"`
	Secret       string `json:"secret,omitempty"`
}

// HttpAuthenticationInfo is the polymorphic authentication object, discriminated by Type.
type HttpAuthenticationInfo struct {
	Type        string           `json:"type"`
	AuthHeaders []AuthHeaderInfo `json:"authHeaders,omitempty"`
	Username    string           `json:"username,omitempty"`
	Password    string           `json:"password,omitempty"`
}

type HttpConnectionInfo struct {
	Url                     string                 `json:"url"`
	Authentication          HttpAuthenticationInfo `json:"authentication"`
	Headers                 map[string]string      `json:"headers,omitempty"`
	DisableCertVerification bool                   `json:"disableCertVerification"`
}

type ConnectionInfo interface {
	getConnectionInfo() any
}

func (s SplunkConnectionInfo) getConnectionInfo() any {
	return SplunkConnectionInfo{
		Host:                    s.Host,
		Port:                    s.Port,
		Token:                   s.Token,
		DisableCertVerification: s.DisableCertVerification,
	}
}

func (s SftpConnectionInfo) getConnectionInfo() any {
	return SftpConnectionInfo{
		Host:     s.Host,
		Username: s.Username,
		Password: s.Password,
		Path:     s.Path,
	}
}

func (c HttpConnectionInfo) getConnectionInfo() any {
	return c
}

func (s S3ConnectionInfo) getConnectionInfo() any {
	return S3ConnectionInfo{
		AccessKey: s.AccessKey,
		SecretKey: s.SecretKey,
		Path:      s.Path,
	}
}

type SiemConnectionData struct {
	ID             string         `json:"id,omitempty"`
	AssetID        string         `json:"assetId,omitempty"`
	ConnectionName string         `json:"connectionName"`
	StorageType    string         `json:"storageType"`
	ConnectionInfo ConnectionInfo `json:"connectionInfo"`
}

func (s *SiemConnectionData) UnmarshalJSON(input []byte) error {
	body := string(input)
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(body), &jsonMap)
	if err != nil {
		return err
	}
	s.ID = jsonMap["id"].(string)
	s.ConnectionName = jsonMap["connectionName"].(string)
	s.AssetID = jsonMap["assetId"].(string)
	s.StorageType = jsonMap["storageType"].(string)
	if s.StorageType == "CUSTOMER_S3" {
		s.ConnectionInfo = S3ConnectionInfo{
			AccessKey: jsonMap["connectionInfo"].(map[string]interface{})["accessKey"].(string),
			SecretKey: jsonMap["connectionInfo"].(map[string]interface{})["secretKey"].(string),
			Path:      jsonMap["connectionInfo"].(map[string]interface{})["path"].(string),
		}
	} else if s.StorageType == "CUSTOMER_S3_ARN" {
		s.ConnectionInfo = S3ConnectionInfo{
			Path: jsonMap["connectionInfo"].(map[string]interface{})["path"].(string),
		}
	} else if s.StorageType == "CUSTOMER_SPLUNK" {
		s.ConnectionInfo = SplunkConnectionInfo{
			Host:                    jsonMap["connectionInfo"].(map[string]interface{})["host"].(string),
			Port:                    int(jsonMap["connectionInfo"].(map[string]interface{})["port"].(float64)),
			Token:                   jsonMap["connectionInfo"].(map[string]interface{})["token"].(string),
			DisableCertVerification: jsonMap["connectionInfo"].(map[string]interface{})["disableCertVerification"].(bool),
		}
	} else if s.StorageType == "CUSTOMER_SFTP" {
		s.ConnectionInfo = SftpConnectionInfo{
			Host:     jsonMap["connectionInfo"].(map[string]interface{})["host"].(string),
			Username: jsonMap["connectionInfo"].(map[string]interface{})["username"].(string),
			Password: jsonMap["connectionInfo"].(map[string]interface{})["password"].(string),
			Path:     jsonMap["connectionInfo"].(map[string]interface{})["path"].(string),
		}
	} else if s.StorageType == StorageTypeCustomerHttp {
		// Re-marshal into the typed struct: optional members (headers, authHeaders, username/password)
		// would panic under the map-assertion idiom used above.
		rawConnectionInfo, marshalErr := json.Marshal(jsonMap["connectionInfo"])
		if marshalErr != nil {
			return fmt.Errorf("marshal CUSTOMER_HTTP connectionInfo: %w", marshalErr)
		}
		var httpConnectionInfo HttpConnectionInfo
		if unmarshalErr := json.Unmarshal(rawConnectionInfo, &httpConnectionInfo); unmarshalErr != nil {
			return fmt.Errorf("unmarshal CUSTOMER_HTTP connectionInfo: %w", unmarshalErr)
		}
		s.ConnectionInfo = httpConnectionInfo
	} else {
		err = errors.New("unsupported ConnectionInfo type")
	}
	return error(err)
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

func (c *Client) ReadSiemConnection(ID string, accountId string) (*SiemConnection, *int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	return siemConnectionRequestWithResponse(c, ReadSiemConnection, http.MethodGet, reqURL, nil, accountId, 200)
}

func (c *Client) UpdateSiemConnection(siemConnection *SiemConnection) (*SiemConnection, *int, error) {
	siemConnectionJSON, err := json.Marshal(siemConnection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to produce JSON from SiemConnectionWithID: %s", err)
	}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, siemConnection.Data[0].ID)
	return siemConnectionRequestWithResponse(c, UpdateSiemConnection, http.MethodPut, reqURL, siemConnectionJSON, siemConnection.Data[0].AssetID, 200)
}

func (c *Client) DeleteSiemConnection(ID string, accountId string) (*int, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.BaseURLAPI, endpointSiemConnection, ID)
	_, _, statusCode, err := siemConnectionRequest(c, DeleteSiemConnection, http.MethodDelete, reqURL, nil, accountId, 200)
	return statusCode, err
}

func dSiemConnectionResponseClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

var siemConnectionSecretKeys = map[string]bool{"secret": true, "password": true}

const redactedSecretValue = "***"

// redactSiemConnectionSecrets masks secret values so request and response bodies can be logged safely.
func redactSiemConnectionSecrets(data []byte) string {
	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return string(data)
	}
	redacted, err := json.Marshal(redactSecretValues(parsed))
	if err != nil {
		return string(data)
	}
	return string(redacted)
}

func redactSecretValues(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, nested := range typed {
			if siemConnectionSecretKeys[key] {
				typed[key] = redactedSecretValue
			} else {
				typed[key] = redactSecretValues(nested)
			}
		}
	case []interface{}:
		for index, nested := range typed {
			typed[index] = redactSecretValues(nested)
		}
	}
	return value
}

func siemConnectionRequest(c *Client, operation string, method string, reqURL string, data []byte, accountIdStr string, expectedSuccessStatusCode int) (*string, *[]byte, *int, error) {
	log.Printf("[INFO] Executing operation %s on SIEM connection with data: %s", operation, redactSiemConnectionSecrets(data))

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
	log.Printf("[DEBUG] Incapsula returned response: %s\nfor %s operation on SIEM connection", redactSiemConnectionSecrets(responseBody), operation)

	if resp.StatusCode != expectedSuccessStatusCode {
		return nil, nil, &resp.StatusCode, formatSiemApiError(operation, "SIEM connection", resp.StatusCode, body)
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
