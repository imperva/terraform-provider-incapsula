package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

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

// Endpoints (unexported consts)
const endpointAccount = "account"
const endpointSiteAdd = "sites/add"
const endpointSiteDelete = "sites/delete"
const endpointSiteStatus = "sites/status"

// DNSRecord contains the relevant DNS information for the Incapsula managed site
type DNSRecord struct {
	DNSRecordName string   `json:"dns_record_name"`
	SetTypeTo     string   `json:"set_type_to"`
	SetDataTo     []string `json:"set_data_to"`
}

// SiteAddResponse contains the relevant site information when adding an Incapsula managed site
type SiteAddResponse struct {
	SiteID int `json:"site_id"`
	Res    int `json:"res"`
}

// SiteStatusResponse contains the relevant site information when getting an Incapsula managed site's status
type SiteStatusResponse struct {
	SiteCreationDate int         `json:"site_creation_date"`
	DNS              []DNSRecord `json:"dns"`
	Res              int         `json:"res"`
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

// AddSite adds a site to be managed by Incapsula
func (c *Client) AddSite(domain, accountID, refID, sendSiteSetupEmails, siteIP, forceSSL, logLevel, logsAccountID string) (*SiteAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula site for domain: %s\n", domain)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteAdd), url.Values{
		"api_id":                 {c.config.APIID},
		"api_key":                {c.config.APIKey},
		"domain":                 {domain},
		"account_id":             {accountID},
		"ref_id":                 {refID},
		"send_site_setup_emails": {sendSiteSetupEmails},
		"site_ip":                {siteIP},
		"force_ssl":              {forceSSL},
		"log_level":              {logLevel},
		"logs_account_id":        {logsAccountID},
	})
	if err != nil {
		return nil, fmt.Errorf("Error adding site for domain %s: %s", domain, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteAddResponse SiteAddResponse
	err = json.Unmarshal([]byte(responseBody), &siteAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add site JSON response for domain %s: %s", domain, err)
	}

	// Look at the response status code from Incapsula
	if siteAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding site for domain %s: %s", domain, string(responseBody))
	}

	return &siteAddResponse, nil
}

// SiteStatus gets the Incapsula managed site's status
func (c *Client) SiteStatus(domain string, siteID int) (*SiteStatusResponse, error) {
	log.Printf("[INFO] Getting Incapsula site status for domain: %s (site id: %d)\n", domain, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteStatus), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting site status for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula site status JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing site status JSON response for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting site status for domain %q (site id: %d): %s", domain, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}

// DeleteSite deletes a site currently managed by Incapsula
func (c *Client) DeleteSite(domain string, siteID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type SiteDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula site for domain: %s (site id: %d)\n", domain, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting site for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteDeleteResponse SiteDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &siteDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete site JSON response for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Look at the response status code from Incapsula
	if siteDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting site for domain %s (site id: %d): %s", domain, siteID, string(responseBody))
	}

	return nil
}
