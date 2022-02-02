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
const endpointSubAccountList = "accounts/listSubAccounts"
const endpointSubAccountDelete = "subaccounts/delete"

type SubAccount struct {
	SubAccountID   int    `json:"sub_account_id"`
	SubAccountName string `json:"sub_account_name"`
	RefID          string `json:"ref_id"`
	LogLevel       string `json:"log_level"`
	SupportLevel   string `json:"support_level"`
	ParentID       int    `json:"parent_id"`
	LogsAccountID  int    `json:"logs_account_id"`
}

// SubAccountAddResponse contains the relevant information when adding an Incapsula SubAccount
type SubAccountAddResponse struct {
	SubAccount SubAccount `json:"sub_account"`
	Res        int        `json:"res"`
}

// SubAccountListResponse contains list of Incapsula SubAccount
type SubAccountListResponse struct {
	SubAccounts []SubAccount `json:"resultList"`
	Res         int          `json:"res"`
}

// SubAccountPayload contains the payload for Incapsula SubAccount creation
type SubAccountPayload struct {
	subAccountName string
	refID          string
	logLevel       string
	logsAccountID  int
	parentID       int
}

// AddSubAccount adds a SubAccount to be managed by Incapsula
func (c *Client) AddSubAccount(subAccountPayload *SubAccountPayload) (*SubAccountAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula subaccount: %s\n", subAccountPayload.subAccountName)

	values := url.Values{
		"sub_account_name": {subAccountPayload.subAccountName},
	}

	if subAccountPayload.refID != "" {
		values["ref_id"] = make([]string, 1)
		values["ref_id"][0] = fmt.Sprint(subAccountPayload.refID)
	}

	if subAccountPayload.parentID != 0 {
		values["parent_id"] = make([]string, 1)
		values["parent_id"][0] = fmt.Sprint(subAccountPayload.parentID)
	}

	if subAccountPayload.logsAccountID != 0 {
		values["logs_account_id"] = make([]string, 1)
		values["logs_account_id"][0] = fmt.Sprint(subAccountPayload.logsAccountID)
	}

	if subAccountPayload.logLevel != "" {
		values["log_level"] = make([]string, 1)
		values["log_level"][0] = fmt.Sprint(subAccountPayload.logLevel)
	}

	log.Printf("[DEBUG] parent_id %d\n", subAccountPayload.parentID)
	log.Printf("[DEBUG] logsAccountID %d\n", subAccountPayload.logsAccountID)
	log.Printf("[DEBUG] logLevel %s\n", subAccountPayload.logLevel)
	log.Printf("[DEBUG] values %s\n", values)

	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountAdd), values)
	if err != nil {
		return nil, fmt.Errorf("Error adding subaccount %s: %s", subAccountPayload.subAccountName, err)
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
		return nil, fmt.Errorf("Error parsing add subaccount JSON response for subaccount %s: %s", subAccountPayload.subAccountName, err)
	}

	// Look at the response status code from Incapsula
	if subAccountAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding subaccount %s: %s", subAccountPayload.subAccountName, string(responseBody))
	}

	return &subAccountAddResponse, nil
}

// ListSubAccounts gets the Incapsula list of SubAccounts
func (c *Client) ListSubAccounts(accountID int) (*SubAccountListResponse, error) {

	log.Printf("[INFO] Getting Incapsula subaccounts for: %d)\n", accountID)

	values := map[string][]string{}

	if accountID != 0 {
		values["account_id"] = make([]string, 1)
		values["account_id"][0] = fmt.Sprint(accountID)
	}

	// Post form to Incapsula
	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountList), values)
	if err != nil {
		return nil, fmt.Errorf("Error getting subaccounts for account %d: %s", accountID, err)
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
		return nil, fmt.Errorf("Error parsing subaccounts list JSON response for accountid: %d %s\nresponse: %s", accountID, err, string(responseBody))
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

		resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountList), values)
		if err != nil {
			return nil, fmt.Errorf("Error getting subaccounts for account %d: %s, page num: %d", accountID, err, count)
		}

		// Read the body
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)

		log.Printf("[DEBUG] Incapsula subaccounts JSON page %d response: %s\n", count, string(responseBody))

		// Clean Var and parse json response
		tempSubAccountListResponse = SubAccountListResponse{}
		err = json.Unmarshal([]byte(responseBody), &tempSubAccountListResponse)
		if err != nil {
			return nil, fmt.Errorf("Error parsing subaccounts list JSON response for accountid: %d %s\nresponse: %s", accountID, err, string(responseBody))
		}

		// Look at the response status code from Incapsula
		if tempSubAccountListResponse.Res != 0 {
			return &tempSubAccountListResponse, fmt.Errorf("Error from Incapsula service when getting sub accounts list %d): %s", accountID, string(responseBody))
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

	// Look at the response status code from Incapsula
	if subAccountListResponse.Res != 0 {
		return &subAccountListResponse, fmt.Errorf("Error from Incapsula service when getting sub accounts list %d): %s", accountID, string(responseBody))
	}

	return &subAccountListResponse, nil
}

// DeleteSubAccount deletes a SubAcccount currently managed by Incapsula
func (c *Client) DeleteSubAccount(subAccountID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type SubAccountDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula subaccount id: %d\n", subAccountID)

	// Post form to Incapsula
	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSubAccountDelete), url.Values{
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
