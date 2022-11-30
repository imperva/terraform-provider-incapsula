package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type AccountPolicyAssociationV3RequestResponse struct {
	Data []AccountPolicyAssociationV3 `json:"data"`
}
type AccountPolicyAssociationV3 struct {
	AccountID                               int   `json:"accountId"`
	AvailablePolicyIds                      []int `json:"availablePolicyIds"`
	DefaultNonMandatoryNonDistinctPolicyIds []int `json:"defaultNonMandatoryNonDistinctPolicyIds"`
	DefaultWafPolicyId                      int   `json:"defaultWafPolicyId,omitempty"`
}

// GetAccountPolicyAssociation get the account policy association for the specified_account
func (c *Client) GetAccountPolicyAssociation(accountId string) (*AccountPolicyAssociationV3, error) {
	log.Printf("[INFO] Getting Policy Association for account: %s\n", accountId)
	//
	// Get the association
	reqURL := fmt.Sprintf("%s/policies/v3/accounts/associated-policies?caid=%s", c.config.BaseURLAPI, accountId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadPolicyAccountAssociatiation)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading Policies Assocication for Account ID %s: %s", accountId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Policy Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service when reading Policy Association for Account ID %s: %s", resp.StatusCode, accountId, string(responseBody))
	}

	// Parse the JSON
	var accountPolicyAssociationV3RequestResponse AccountPolicyAssociationV3RequestResponse
	err = json.Unmarshal([]byte(responseBody), &accountPolicyAssociationV3RequestResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Policies Association JSON response for Account ID %s: %s\nresponse: %s", accountId, err, string(responseBody))
	}
	if accountPolicyAssociationV3RequestResponse.Data == nil || len(accountPolicyAssociationV3RequestResponse.Data) == 0 {
		return nil, fmt.Errorf("[ERROR] got epmty response for Account ID %s\nresponse: %s", accountId, string(responseBody))
	}
	return &accountPolicyAssociationV3RequestResponse.Data[0], nil
}

// PatchAccountPolicyAssociation get the account policy association for the specified_account
func (c *Client) PatchAccountPolicyAssociation(accountId string, availablePolicyIds []int, defaultNonMandatoryPolicyIds []int, wafPolicyIdStr string) (*AccountPolicyAssociationV3, error) {
	log.Printf("[INFO] Setting Policy Association for account: %s, WAF Rules Policy: %s, Default non mandatory non distinct: %v\n", accountId, wafPolicyIdStr, defaultNonMandatoryPolicyIds)

	//Build the policy association request data
	var accountPolicyAssociationV3 AccountPolicyAssociationV3
	accountIdInt, err := strconv.Atoi(accountId)
	if err != nil {
		log.Printf("[ERROR] Could not convert Account ID. Error: is not numeric: %s", accountId)
		return nil, err
	}
	accountPolicyAssociationV3.AccountID = accountIdInt
	if wafPolicyIdStr != "" {
		wafPolicyId, err := strconv.Atoi(wafPolicyIdStr)
		if err != nil {
			log.Printf("[ERROR] Could not convert WAF Rule Policy ID. Error: is not numeric: %s", wafPolicyIdStr)
			return nil, err
		}
		accountPolicyAssociationV3.DefaultWafPolicyId = wafPolicyId
	}
	if defaultNonMandatoryPolicyIds != nil {
		accountPolicyAssociationV3.DefaultNonMandatoryNonDistinctPolicyIds = defaultNonMandatoryPolicyIds
	}

	if availablePolicyIds != nil {
		accountPolicyAssociationV3.AvailablePolicyIds = availablePolicyIds
	}
	var accountPolicyAssociationV3RequestResponse AccountPolicyAssociationV3RequestResponse
	accountPolicyAssociationV3RequestResponse.Data = make([]AccountPolicyAssociationV3, 1)
	accountPolicyAssociationV3RequestResponse.Data[0] = accountPolicyAssociationV3

	// Patch the association
	reqURL := fmt.Sprintf("%s/policies/v3/accounts/associated-policies?caid=%s", c.config.BaseURLAPI, accountId)
	byteJSON, err := json.Marshal(accountPolicyAssociationV3RequestResponse)
	if err != nil {
		log.Printf("[ERROR] Failed to create body for request %+v", accountPolicyAssociationV3RequestResponse)
		return nil, err
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPatch, reqURL, byteJSON, UpdatePolicyAccountAssociatiation)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when setting Policies Assocication for Account ID %s with body %+v: %s",
			accountId, accountPolicyAssociationV3RequestResponse, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Policy Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Error status code %d from Incapsula service when setting Policy Association for Account ID %s with body %+v: %s",
			resp.StatusCode, accountId, accountPolicyAssociationV3RequestResponse, string(responseBody))
	}

	// Parse the JSON
	err = json.Unmarshal([]byte(responseBody), &accountPolicyAssociationV3RequestResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Policies Association JSON response for Account ID %s: %s\nresponse: %s", accountId, err, string(responseBody))
	}
	if accountPolicyAssociationV3RequestResponse.Data == nil || len(accountPolicyAssociationV3RequestResponse.Data) == 0 {
		return nil, fmt.Errorf("[ERROR] got epmty response for Account ID %s\nresponse: %s", accountId, string(responseBody))
	}
	return &accountPolicyAssociationV3RequestResponse.Data[0], nil
}
