package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Endpoints (unexported consts)
const endpointAccount = "account"

// Client represents an internal client that brokers calls to the Incapsula API
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new client with the provided configuration
func NewClient(config *Config) *Client {
	client := &http.Client{}
	return &Client{config: config, httpClient: client}
}

// Verify checks the API credentials
func (c *Client) Verify() error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type AccountResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Println("[INFO] Checking API credentials against Incapsula API")

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccount), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
	})
	if err != nil {
		return fmt.Errorf("Error checking account: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula acount JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountResponse AccountResponse
	err = json.Unmarshal([]byte(responseBody), &accountResponse)
	if err != nil {
		return fmt.Errorf("Error parsing account JSON response: %s", err)
	}

	// Look at the response status code from Incapsula
	if accountResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when checking account: %s", string(responseBody))
	}

	return nil
}
