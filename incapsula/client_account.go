package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

// Endpoints (unexported consts)
const endpointAccountAdd = "accounts/add"
const endpointAccountStatus = "account"
const endpointAccountUpdate = "accounts/configure"
const endpointAccountDelete = "accounts/delete"

// AccountAddResponse contains the relevant account information when adding an Incapsula Account
type AccountAddResponse struct {
	Account struct {
		ParentID    int    `json:"parent_id"`
		Email       string `json:"email"`
		PlanID      string `json:"plan_id"`
		AccountID   int    `json:"account_id"`
		UserName    string `json:"user_name"`
		AccountName string `json:"account_name"`
		Logins      []struct {
			LoginID       float64 `json:"login_id"`
			Email         string  `json:"email"`
			EmailVerified bool    `json:"email_verified"`
		} `json:"logins"`
	} `json:"account"`
	Res int `json:"res"`
}

// AccountUpdateResponse contains the relevant account information when updating an Incapsula Account
type AccountUpdateResponse struct {
	AccountID int `json:"account_id"`
	Res       int `json:"res"`
}

// AccountResponse contains account id
type AccountStatusResponse struct {
	Account struct {
		ParentID     int    `json:"parent_id"`
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
	ParentID    int    `json:"parent_id"`
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
	SupportLevel                 string      `json:"support_level"`
	SupportAllTLSVersions        bool        `json:"supprt_all_tls_versions"`
	WildcardSANForNewSites       string      `json:"wildcard_san_for_new_sites"`
	NakedDomainSANForNewWWWSites bool        `json:"naked_domain_san_for_new_www_sites"`
	Res                          interface{} `json:"res"`
	ResMessage                   string      `json:"res_message"`
	DebugInfo                    struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// AddAccount adds an account to be managed by Incapsula
func (c *Client) AddAccount(email, refID, userName, planID, accountName, logLevel string, logsAccountID int, parentID int) (*AccountAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula account for email: %s (account ID %d)\n", email, parentID)

	values := url.Values{
		"email":        {email},
		"user_name":    {userName},
		"plan_id":      {planID},
		"ref_id":       {refID},
		"account_name": {accountName},
		"log_level":    {logLevel},
	}
	if parentID != 0 {
		values["parent_id"] = make([]string, 1)
		values["parent_id"][0] = fmt.Sprint(parentID)
	}

	if logsAccountID != 0 {
		values["logs_account_id"] = make([]string, 1)
		values["logs_account_id"][0] = fmt.Sprint(logsAccountID)
	}

	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountAdd)
	resp, err := c.PostFormWithHeaders(reqURL, values, CreateAccount)
	if err != nil {
		return nil, fmt.Errorf("Error adding account for email %s: %s", email, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add account JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountAddResponse AccountAddResponse
	err = json.Unmarshal([]byte(responseBody), &accountAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add account JSON response for email %s: %s", email, err)
	}

	// Look at the response status code from Incapsula
	if accountAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding account for email %s: %s", email, string(responseBody))
	}

	return &accountAddResponse, nil
}

// AccountStatus gets the Incapsula managed account's status
func (c *Client) AccountStatus(accountID int, operation string) (*AccountStatusResponse, error) {
	log.Printf("[INFO] Getting Incapsula account status for account id: %d\n", accountID)

	// Post form to Incapsula
	values := url.Values{"account_id": {strconv.Itoa(accountID)}}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountStatus)
	resp, err := c.PostFormWithHeaders(reqURL, values, operation)
	if err != nil {
		return nil, fmt.Errorf("Error getting account status for account id %d: %s", accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula account status JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountStatusResponse AccountStatusResponse
	err = json.Unmarshal([]byte(responseBody), &accountStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing account status JSON response for account id %d: %s", accountID, err)
	}

	var resString string

	if resNumber, ok := accountStatusResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = accountStatusResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString != "0" {
		return &accountStatusResponse, fmt.Errorf("Error from Incapsula service when getting account status for account id %d: %s", accountID, string(responseBody))
	}

	return &accountStatusResponse, nil
}

// UpdateAccount will update the specific param/value on the account resource
func (c *Client) UpdateAccount(accountID, param, value string) (*AccountUpdateResponse, error) {
	log.Printf("[INFO] Updating Incapsula account for accountID: %s\n", accountID)

	values := url.Values{
		"account_id": {accountID},
		"param":      {param},
		"value":      {value},
	}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountUpdate)
	resp, err := c.PostFormWithHeaders(reqURL, values, UpdateAccount)
	if err != nil {
		return nil, fmt.Errorf("Error updating param (%s) with value (%s) on account_id: %s: %s", param, value, accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update account JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountUpdateResponse AccountUpdateResponse
	err = json.Unmarshal([]byte(responseBody), &accountUpdateResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update account JSON response for accountID %s: %s", accountID, err)
	}

	// Look at the response status code from Incapsula
	if accountUpdateResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating account for accountID %s: %s", accountID, string(responseBody))
	}

	return &accountUpdateResponse, nil
}

// DeleteAccount deletes a account currently managed by Incapsula
func (c *Client) DeleteAccount(accountID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type AccountDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula account id: %d\n", accountID)

	// Post form to Incapsula
	values := url.Values{"account_id": {strconv.Itoa(accountID)}}
	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURL, endpointAccountDelete)
	resp, err := c.PostFormWithHeaders(reqURL, values, DeleteAccount)
	if err != nil {
		return fmt.Errorf("Error deleting account id: %d: %s", accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete account JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var accountDeleteResponse AccountDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &accountDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete account JSON response for account id: %d: %s", accountID, err)
	}

	// Look at the response status code from Incapsula
	if accountDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting account id: %d: %s", accountID, string(responseBody))
	}

	return nil
}
