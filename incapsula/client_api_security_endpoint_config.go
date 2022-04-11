package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	//"io"
	"io/ioutil"
	"log"
	//"strings"
)

const endpointConfigUrl = "/api-security/endpoint/"

type ApiSecurityEndpointConfigGetResponse struct {
	Value   EndpointResponse `json:"value"`
	IsError bool             `json:"is_error"`
}

type ApiSecurityEndpointConfigGetAllResponse struct {
	Value   []EndpointResponse `json:"value"`
	IsError bool               `json:"is_error"`
}

type ApiSecurityEndpointConfigPostResponse struct {
	Value struct {
		EndpointId int `json:"endpoint_id"`
	} `json:"value"`
	IsError bool `json:"is_error"`
}

type ApiSecurityEndpointConfigPostPayload struct {
	ViolationActions             UserViolationActions
	SpecificationViolationAction string
}

//PostApiSecurityEndpointConfig updates an Api-Security Endpoint Config
func (c *Client) PostApiSecurityEndpointConfig(apiId, endpointId int, endpointConfigPayload *ApiSecurityEndpointConfigPostPayload) (*ApiSecurityEndpointConfigPostResponse, error) {
	log.Printf("[INFO] Updating Incapsula API security Enpoint Configuration\n")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	violationActionsStr, err := json.Marshal(endpointConfigPayload.ViolationActions)
	if err != nil {
		fmt.Println(err)
	}
	if endpointConfigPayload.ViolationActions != (UserViolationActions{}) {
		fw, err := writer.CreateFormField("violationActions")
		if err != nil {
		}
		_, err = io.Copy(fw, strings.NewReader(string(violationActionsStr)))
		if err != nil {
		}
	}
	if endpointConfigPayload.SpecificationViolationAction != "" {
		fw, err := writer.CreateFormField("specificationViolationAction")
		if err != nil {
		}
		_, err = io.Copy(fw, strings.NewReader(endpointConfigPayload.SpecificationViolationAction))
		if err != nil {
		}
	}

	writer.Close()
	url := fmt.Sprintf("%s%s%d"+"/"+"%d", c.config.BaseURLAPI, endpointConfigUrl, apiId, endpointId)
	contentType := writer.FormDataContentType()
	resp, err := c.DoJsonRequestWithHeadersForm(http.MethodPost, url, body.Bytes(), contentType, UpdateApiSecEndpointConfig)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service while updating Api Security Endpoint Configuration for API Config Id %d, API Config Id %d : %s", apiId, endpointId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Update Api-Security Endpoint Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service while updating Api Security Endpoint configuration for API Config Id %d, Endpoint Config Id: %d. Error: %s", resp.StatusCode, apiId, endpointId, string(responseBody))
	}

	// Parse the JSON
	var response ApiSecurityEndpointConfigPostResponse
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing api-security JSON response for create/update Api Security Endpoint Configuration for API Config Id %d, Endpoint Config Id %d : %s\nresponse: %s", apiId, endpointId, err, string(responseBody))
	}

	return &response, nil
}

// GetApiSecurityEndpointConfig gets the Api-Security Endpoint Config
func (c *Client) GetApiSecurityEndpointConfig(apiId int, endpointId string) (*ApiSecurityEndpointConfigGetResponse, error) {
	log.Printf("[INFO] Getting Incapsula Api-Security Endpoint Config on API: %d and Endpoint: %s\n", apiId, endpointId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, fmt.Sprintf("%s%s%d/%s", c.config.BaseURLAPI, endpointConfigUrl, apiId, endpointId), nil, ReadApiSecEndpointConfig)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service while reading Api-Security Endpoint Config for API ID %d and Endpoint ID %s: %s", apiId, endpointId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Read Api-Security Endpoint Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service when reading Api-Security Endpoint Config for API ID %d and Endpoint ID %s: %s", resp.StatusCode, apiId, endpointId, string(responseBody))
	}

	// Parse the JSON
	var apiSecurityEndpointConfigGetResponse ApiSecurityEndpointConfigGetResponse
	err = json.Unmarshal([]byte(responseBody), &apiSecurityEndpointConfigGetResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing GET Api-Security Endpoint Config JSON response for API ID %d and endpoint ID %s: %s\nresponse: %s", apiId, endpointId, err, string(responseBody))
	}

	return &apiSecurityEndpointConfigGetResponse, nil
}

// GetApiSecurityAllEndpointsConfig gets all the Api-Security Endpoints for API Config ID
func (c *Client) GetApiSecurityAllEndpointsConfig(apiId int) (*ApiSecurityEndpointConfigGetAllResponse, error) {
	log.Printf("[INFO] Getting Incapsula Api-Security all Endpoints Config on API: %d\n", apiId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, fmt.Sprintf("%s%s%d", c.config.BaseURLAPI, endpointConfigUrl, apiId), nil, ReadApiSecEndpointConfig)
	if err != nil {
		return nil, fmt.Errorf("error from Incapsula service when reading Api-Security all Endpoints Config for API ID %d: %s", apiId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula Read All Api-Security Endpoint Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error status code %d from Incapsula service when reading Api-Security all Endpoints Config for API ID %d: %s", resp.StatusCode, apiId, string(responseBody))
	}

	// Parse the JSON
	var apiSecurityEndpointConfigGetAllResponse ApiSecurityEndpointConfigGetAllResponse
	err = json.Unmarshal([]byte(responseBody), &apiSecurityEndpointConfigGetAllResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing GET Api-Security all Endpoints Config JSON response for API ID %d: %s\nresponse: %s", apiId, err, string(responseBody))
	}

	return &apiSecurityEndpointConfigGetAllResponse, nil
}
