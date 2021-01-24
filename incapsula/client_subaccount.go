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
const endpointSubAccountAdd = "subaccounts/add"
const endpointSubAccountStatus = "accounts/listSubAccounts"
const endpointSubAccountDelete = "subaccounts/delete"

// SubAccountAddResponse contains the relevant information when adding an Incapsula SubAccount
type SubAccountAddResponse struct {
	SubAccount struct {
		SubAccountID   int    `json:"sub_account_id"`
		SubAccountName string `json:"sub_account_name"`
		SpeicalSSL     bool   `json:"is_for_special_ssl_configuration"`
		SupportLevel   string `json:"support_level"`
	} `json:"sub_account"`
	Res int `json:"res"`
}

// DataCenterListResponse contains list of data centers and servers
type SubAccountListResponse struct {
	Res         interface{} `json:"res"`
	SubAccounts []struct {
		SubAccountID   int    `json:"sub_account_id"`
		SubAccountName string `json:"sub_account_name"`
		SpeicalSSL     bool   `json:"is_for_special_ssl_configuration"`
		SupportLevel   string `json:"support_level"`
		Logins         []struct {
			LoginID       float64 `json:"login_id"`
			Email         string  `json:"email"`
			EmailVerified bool    `json:"email_verified"`
		} `json:"logins"`
	} `json:"resultList"`
}

// AddSubAccount adds an subaccount to be managed by Incapsula
func (c *Client) AddSubAccount(subAccountName, refID, logLevel string, logsAccountID int, parentID int) (*SubAccountAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula subaccount: %s\n", subAccountName)

	values := url.Values{
		"api_id":           {c.config.APIID},
		"api_key":          {c.config.APIKey},
		"ref_id":           {refID},
		"sub_account_name": {subAccountName},
	}
	if parentID != 0 {
		values["parent_id"] = make([]string, 1)
		values["parent_id"][0] = fmt.Sprint(parentID)
	}

	if logsAccountID != 0 {
		values["logs_account_id"] = make([]string, 1)
		values["logs_account_id"][0] = fmt.Sprint(logsAccountID)
	}

	if logLevel != "" {
		values["log_level"] = make([]string, 1)
		values["log_level"][0] = fmt.Sprint(logLevel)
	}

	log.Printf("[DEBUG] parent_id %d\n", parentID)
	log.Printf("[DEBUG] logsAccountID %d\n", logsAccountID)
	log.Printf("[DEBUG] logLevel %s\n", logLevel)
	log.Printf("[DEBUG] values %s\n", values)

	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountAdd), values)
	if err != nil {
		return nil, fmt.Errorf("Error adding subaccount %s: %s", subAccountName, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add subaccount JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var subAccountAddResponse SubAccountAddResponse
	err = json.Unmarshal([]byte(responseBody), &subAccountAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add subaccount JSON response for subaccount %s: %s", subAccountName, err)
	}

	// Look at the response status code from Incapsula
	if subAccountAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding subaccount %s: %s", subAccountName, string(responseBody))
	}

	return &subAccountAddResponse, nil
}

// ListSubAccounts gets the Incapsula list of sub accounts
func (c *Client) ListSubAccounts(AccountID int) (*SubAccountListResponse, error) {

	log.Printf("[INFO] Getting Incapsula subaccounts for: %d)\n", AccountID)

	values := url.Values{
		"api_id":    {c.config.APIID},
		"api_key":   {c.config.APIKey},
		"page_size": {strconv.Itoa(50)},
	}

	if AccountID != 0 {
		values["account_id"] = make([]string, 1)
		values["account_id"][0] = fmt.Sprint(AccountID)
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountStatus), values)
	if err != nil {
		return nil, fmt.Errorf("Error getting subaccounts for account %d: %s", AccountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula subaccounts JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var subAccountListResponse SubAccountListResponse
	err = json.Unmarshal([]byte(responseBody), &subAccountListResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing subaccounts list JSON response for accountid: %d %s\nresponse: %s", AccountID, err, string(responseBody))
	}

	log.Printf("[INFO] Array before loop : %v)\n", subAccountListResponse.SubAccounts)

	// Pagination (default page size 50)
	tempSubAccountListResponse := subAccountListResponse
	var count int = 1
	for len(tempSubAccountListResponse.SubAccounts) == 50 {
		log.Printf("[INFO] Pagination loop, page : %d)\n", count)
		values["page_num"] = make([]string, count)
		values["page_num"][0] = fmt.Sprint(count)
		log.Printf("[INFO] values : %s)\n", values)

		resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountStatus), values)
		if err != nil {
			return nil, fmt.Errorf("Error getting subaccounts for account %d: %s, page num: %d", AccountID, err, count)
		}

		// Read the body
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)

		log.Printf("[DEBUG] Incapsula subaccounts JSON page %d response: %s\n", count, string(responseBody))

		// Clean Var and parse json response
		tempSubAccountListResponse = SubAccountListResponse{}
		err = json.Unmarshal([]byte(responseBody), &tempSubAccountListResponse)
		if err != nil {
			return nil, fmt.Errorf("Error parsing subaccounts list JSON response for accountid: %d %s\nresponse: %s", AccountID, err, string(responseBody))
		}

		// add sub-accounts to reponse
		for _, subAccount := range tempSubAccountListResponse.SubAccounts {
			log.Printf("[INFO] Array to add : %v", subAccount)
			subAccountListResponse.SubAccounts = append(subAccountListResponse.SubAccounts, subAccount)
		}
		log.Printf("[INFO] Length of SubAccounts Array : %d)\n", len(subAccountListResponse.SubAccounts))
		log.Printf("[INFO] Array : %v", subAccountListResponse.SubAccounts)
		count += 1
	}

	// Res can sometimes oscillate between a string and number
	// We need to add safeguards for this inside the provider
	var resString string

	if resNumber, ok := subAccountListResponse.Res.(float64); ok {
		resString = fmt.Sprintf("%d", int(resNumber))
	} else {
		resString = subAccountListResponse.Res.(string)
	}

	// Look at the response status code from Incapsula
	if resString != "0" {
		return &subAccountListResponse, fmt.Errorf("Error from Incapsula service when getting sub accounts list %d): %s", AccountID, string(responseBody))
	}

	return &subAccountListResponse, nil
}

// DeleteAccount deletes a account currently managed by Incapsula
func (c *Client) DeleteSubAccount(subAccountID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type SubAccountDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula subaccount id: %d\n", subAccountID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountDelete), url.Values{
		"api_id":         {c.config.APIID},
		"api_key":        {c.config.APIKey},
		"sub_account_id": {strconv.Itoa(subAccountID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting subaccount id: %d: %s", subAccountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete subaccount JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var subaccountDeleteResponse SubAccountDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &subaccountDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete account JSON response for subaccount id: %d: %s", subAccountID, err)
	}

	// Look at the response status code from Incapsula
	if subaccountDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting subaccount id: %d: %s", subAccountID, string(responseBody))
	}

	return nil
}
