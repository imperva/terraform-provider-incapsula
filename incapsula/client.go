package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const contentTypeApplicationUrlEncoded = "application/x-www-form-urlencoded"
const contentTypeApplicationJson = "application/json"

// Client represents an internal client that brokers calls to the Incapsula API
type Client struct {
	config          *Config
	httpClient      *http.Client
	providerVersion string
}

// NewClient creates a new client with the provided configuration
func NewClient(config *Config) *Client {
	client := &http.Client{}

	return &Client{config: config, httpClient: client, providerVersion: "3.8.4"}
}

// Verify checks the API credentials
func (c *Client) Verify() (*AccountStatusResponse, error) {
	log.Println("[INFO] Checking API credentials against Incapsula API")

	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountStatus)
	data := url.Values{}

	resp, err := c.PostFormWithHeaders(reqURL, data, VerifyAccount)
	if err != nil {
		return nil, fmt.Errorf("Error checking account: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula account JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountStatusResponse AccountStatusResponse
	err = json.Unmarshal([]byte(responseBody), &accountStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing account JSON response: %s", err)
	}

	var resString string

	if resNumber, ok := accountStatusResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = accountStatusResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString != "0" {
		return &accountStatusResponse, fmt.Errorf("Error from Incapsula service when checking account: %s", string(responseBody))
	}

	return &accountStatusResponse, nil
}

func (c *Client) PostFormWithHeaders(url string, data url.Values, operation string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentTypeApplicationUrlEncoded, operation, nil)
	return c.httpClient.Do(req)
}

func (c *Client) DoJsonRequestWithCustomHeaders(method string, url string, data []byte, headers map[string]string, operation string) (*http.Response, error) {
	req, err := PrepareJsonRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentTypeApplicationJson, operation, headers)

	return c.httpClient.Do(req)
}

func (c *Client) DoJsonRequestWithHeaders(method string, url string, data []byte, operation string) (*http.Response, error) {
	return c.DoJsonRequestWithCustomHeaders(method, url, data, nil, operation)
}

func (c *Client) DoJsonAndQueryParamsRequestWithHeaders(method string, url string, data []byte, params map[string]string, operation string) (*http.Response, error) {
	req, err := PrepareJsonRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}
	q := req.URL.Query()
	for name, value := range params {
		q.Add(name, value)
	}
	req.URL.RawQuery = q.Encode()

	SetHeaders(c, req, contentTypeApplicationJson, operation, nil)

	return c.httpClient.Do(req)
}

// GetRequestParamsWithCaid Use this function if you want to add caid to your request as a query param.
// you need to send caid if you want to preform action on resources belong to child account (example: reseller -> account)
func GetRequestParamsWithCaid(accountId int) map[string]string {
	var params = map[string]string{}
	if accountId != 0 {
		params["caid"] = strconv.Itoa(accountId)
	}

	return params
}

func (c *Client) DoJsonRequestWithHeadersForm(method string, url string, data []byte, contentType string, operation string) (*http.Response, error) {
	req, err := PrepareJsonRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentType, operation, nil)
	return c.httpClient.Do(req)
}

func PrepareJsonRequest(method string, url string, data []byte) (*http.Request, error) {
	if data == nil {
		return http.NewRequest(method, url, nil)
	}

	return http.NewRequest(method, url, bytes.NewReader(data))
}

func SetHeaders(c *Client, req *http.Request, contentType string, operation string, customHeaders map[string]string) {
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-id", c.config.APIID)
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("x-tf-provider-ver", c.providerVersion)
	req.Header.Set("x-tf-operation", operation)

	if customHeaders != nil {
		for name, value := range customHeaders {
			req.Header.Set(name, value)
		}
	}

}
