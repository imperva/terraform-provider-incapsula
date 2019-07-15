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

// AccountResponse contains account id
type AccountResponse struct {
	Account struct {
		Email        string `json:"email"`
		PlanID       string `json:"plan_id"`
		PlanName     string `json:"plan_name"`
		TrialEndDate string `json:"trial_end_date"`
		AccountID    int    `json:"account_id"`
		RefID        string `json:"ref_id"`
		UserName     string `json:"user_name"`
		AccountName  string `json:"account_name"`
		Logins       []struct {
			LoginID       float64 `json:"login_id"`
			Email         string  `json:"email"`
			EmailVerified bool    `json:"email_verified"`
		} `json:"logins"`
		SupportLevel                 string `json:"support_level"`
		SupportAllTLSVersions        bool   `json:"supprt_all_tls_versions"`
		WildcardSANForNewSites       string `json:"wildcard_san_for_new_sites"`
		NakedDomainSANForNewWWWSites bool   `json:"naked_domain_san_for_new_www_sites"`
	} `json:"account"`
	Email       string `json:"email"`
	PlanID      string `json:"plan_id"`
	PlanName    string `json:"plan_name"`
	AccountID   int    `json:"account_id"`
	UserName    string `json:"user_name"`
	AccountName string `json:"account_name"`
	RefID       string `json:"ref_id"`
	Logins      []struct {
		LoginID       float64 `json:"login_id"`
		Email         string  `json:"email"`
		EmailVerified bool    `json:"email_verified"`
	} `json:"logins"`
	SupportLevel                 string `json:"support_level"`
	SupportAllTLSVersions        bool   `json:"supprt_all_tls_versions"`
	WildcardSANForNewSites       string `json:"wildcard_san_for_new_sites"`
	NakedDomainSANForNewWWWSites bool   `json:"naked_domain_san_for_new_www_sites"`
	Res                          int    `json:"res"`
	ResMessage                   string `json:"res_message"`
	DebugInfo                    struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

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
func (c *Client) Verify() (*AccountResponse, error) {
	log.Println("[INFO] Checking API credentials against Incapsula API")

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccount), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
	})
	if err != nil {
		return nil, fmt.Errorf("Error checking account: %s", err)
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
		return nil, fmt.Errorf("Error parsing account JSON response: %s", err)
	}

	// Look at the response status code from Incapsula
	if accountResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when checking account: %s", string(responseBody))
	}

	return &accountResponse, nil
}
