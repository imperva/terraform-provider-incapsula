package incapsula

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const contentTypeApplicationUrlEncoded = "application/x-www-form-urlencoded"
const contentTypeApplicationJson = "application/json"

// Retry defaults for transient API failures (502, 503, 504, 429, HTML error pages).
// These are vars (not consts) so tests can override them for fast execution.
var maxRetries = 4
var retryWaitMinSeconds = 1
var retryWaitMaxSeconds = 30

// Client represents an internal client that brokers calls to the Incapsula API
type Client struct {
	config          *Config
	httpClient      *http.Client
	providerVersion string
	accountStatus   *AccountStatusResponse
}

// NewClient creates a new client with the provided configuration
func NewClient(config *Config) *Client {
	client := &http.Client{}

	return &Client{config: config, httpClient: client, providerVersion: "3.39.3"}
}

func (c *Client) CreateFormDataBody(bodyMap map[string]interface{}) ([]byte, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range bodyMap {
		switch value.(type) {
		case string:
			fw, err := writer.CreateFormField(key)
			if err != nil {
				log.Printf("failed to create %s formdata field", key)
			}
			_, err = io.Copy(fw, strings.NewReader(fmt.Sprintf("%v", value)))
			break
		case []byte:
			fw, err := writer.CreateFormFile(key, filepath.Base(key+".pfx")) //todo KATRIN try to remove .pfx
			if err != nil {
				log.Printf("failed to create %s formdata field", key)
			}
			fw.Write(value.([]byte))
			break
		default:
			//throw error
		}
	}
	writer.Close()

	return body.Bytes(), writer.FormDataContentType()
}

// Verify the API credentials using the lightweight verify endpoint
func (c *Client) Verify() (*AccountStatusResponse, error) {
	log.Println("[INFO] Checking API credentials against Incapsula API")

	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountVerify)
	data := url.Values{}

	resp, err := c.PostFormWithHeaders(reqURL, data, VerifyAccount)
	if err != nil {
		return nil, fmt.Errorf("Error checking account: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Successful test of API credentials.")

	// Parse the JSON using the lightweight verify response
	var accountVerifyResponse AccountVerifyResponse
	err = json.Unmarshal([]byte(responseBody), &accountVerifyResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing account JSON response: %s", err)
	}

	var resString string

	if resNumber, ok := accountVerifyResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = accountVerifyResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString != "0" {
		return nil, fmt.Errorf("Error from Incapsula service when checking account: %s", string(responseBody))
	}

	// Convert the lightweight verify response to AccountStatusResponse for backward compatibility
	accountStatusResponse := &AccountStatusResponse{
		AccountType: accountVerifyResponse.AccountType,
		AccountID:   accountVerifyResponse.AccountID,
		ParentID:    accountVerifyResponse.ParentID,
		AccountName: accountVerifyResponse.AccountName,
		PlanID:      accountVerifyResponse.PlanID,
		PlanName:    accountVerifyResponse.PlanName,
		Res:         accountVerifyResponse.Res,
		ResMessage:  accountVerifyResponse.ResMessage,
		DebugInfo:   accountVerifyResponse.DebugInfo,
	}

	return accountStatusResponse, nil
}

func (c *Client) PostFormWithHeaders(url string, data url.Values, operation string) (*http.Response, error) {
	encoded := []byte(data.Encode())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentTypeApplicationUrlEncoded, operation, nil)
	return c.executeRequest(req)
}

func (c *Client) GetWithHeaders(url string, queryParams url.Values, operation string) (*http.Response, error) {
	reqURL := url
	if len(queryParams) > 0 {
		reqURL = fmt.Sprintf("%s?%s", url, queryParams.Encode())
	}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentTypeApplicationJson, operation, nil)
	return c.executeRequest(req)
}

func (c *Client) DoJsonRequestWithCustomHeaders(method string, url string, data []byte, headers map[string]string, operation string) (*http.Response, error) {
	req, err := PrepareJsonRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentTypeApplicationJson, operation, headers)

	return c.executeRequest(req)
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

	return c.executeRequest(req)
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

func (c *Client) DoFormDataRequestWithHeaders(method string, url string, data []byte, contentType string, operation string) (*http.Response, error) {
	req, err := PrepareJsonRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Error preparing request: %s", err)
	}

	SetHeaders(c, req, contentType, operation, nil)
	return c.executeRequest(req)
}

func PrepareJsonRequest(method string, url string, data []byte) (*http.Request, error) {
	if data == nil {
		return http.NewRequest(method, url, nil)
	}

	return http.NewRequest(method, url, bytes.NewReader(data))
}

func SetHeaders(c *Client, req *http.Request, contentType string, operation string, customHeaders map[string]string) {
	if req.Method != http.MethodGet {
		req.Header.Set("Content-Type", contentType)
	}
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

func (c *Client) executeRequest(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			if req.GetBody != nil {
				req.Body, _ = req.GetBody()
			}
			delay := time.Duration(retryWaitMinSeconds) * time.Second * (1 << (attempt - 1))
			maxDelay := time.Duration(retryWaitMaxSeconds) * time.Second
			if delay > maxDelay {
				delay = maxDelay
			}
			var jitter time.Duration
			if delay > 0 {
				jitter = time.Duration(rand.Int63n(int64(delay) / 4))
			}
			if err != nil {
				log.Printf("[WARN] Transient network error, retry %d/%d for %s %s (backoff %s): %v",
					attempt, maxRetries, req.Method, req.URL.Path, delay+jitter, err)
			} else {
				log.Printf("[WARN] Transient error (status %d), retry %d/%d for %s %s (backoff %s)",
					resp.StatusCode, attempt, maxRetries, req.Method, req.URL.Path, delay+jitter)
			}
			time.Sleep(delay + jitter)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if isTransientNetError(err) && attempt < maxRetries {
				continue
			}
			return nil, err
		}

		if !c.isRetryableResponse(req, resp) {
			return resp, nil
		}

		if attempt == maxRetries {
			log.Printf("[WARN] Retries exhausted (status %d) for %s %s, returning last response",
				resp.StatusCode, req.Method, req.URL.Path)
			return resp, nil
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	return nil, fmt.Errorf("request to %s %s failed after %d retries: last status %d", req.Method, req.URL.Path, maxRetries, resp.StatusCode)
}

func (c *Client) isRetryableResponse(req *http.Request, resp *http.Response) bool {
	if resp.StatusCode == 429 {
		return true
	}

	isRead := req.Method == http.MethodGet ||
		strings.HasPrefix(strings.ToLower(req.Header.Get("x-tf-operation")), "read")

	if resp.StatusCode >= 500 {
		if isRead {
			return true
		}
		return c.responseBodyIsHTML(resp)
	}

	if resp.StatusCode == 200 && c.responseBodyIsHTML(resp) {
		return true
	}

	return false
}

func isTransientNetError(err error) bool {
	if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

func (c *Client) responseBodyIsHTML(resp *http.Response) bool {
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "text/html") {
		return true
	}
	if strings.Contains(ct, "application/json") {
		return false
	}

	peek := make([]byte, 1)
	n, err := resp.Body.Read(peek)
	if err != nil || n == 0 {
		return false
	}
	resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(peek[:n]), resp.Body))
	return peek[0] == '<'
}
