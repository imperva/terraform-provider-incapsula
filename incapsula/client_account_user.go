package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Endpoints (unexported consts)

const endpointUserOperationNew = "identity-management/v3/idm-users"

// UserApisResponse contains the relevant user information when adding, getting or updating a user
type UserApisResponse struct {
	Data []struct {
		UserID      string   `json:"id"`
		AccountID   int      `json:"accountId"`
		FirstName   string   `json:"firstName"`
		LastName    string   `json:"lastName"`
		Email       string   `json:"email"`
		ApprovedIps []string `json:"approvedIps"`
		Roles       []struct {
			RoleID   int    `json:"id"`
			RoleName string `json:"name"`
		} `json:"roles"`
	} `json:"data"`
}

type UserApisUpdateResponse struct {
	Data []struct {
		UserID      string   `json:"id"`
		AccountID   int      `json:"accountId"`
		FirstName   string   `json:"firstName"`
		LastName    string   `json:"lastName"`
		Email       string   `json:"email"`
		ApprovedIps []string `json:"approvedIps"`
		Roles       []struct {
			RoleID   int    `json:"id"`
			RoleName string `json:"name"`
		} `json:"roles"`
	} `json:"data"`
}

type UserAddReq struct {
	UserEmail   string   `json:"email"`
	RoleIds     []int    `json:"roleIds"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	ApprovedIps []string `json:"approvedIps"`
}

type UserUpdateReq struct {
	RoleIds     *[]int    `json:"roleIds,omitempty"`
	ApprovedIps *[]string `json:"approvedIps,omitempty"`
}

// AddAccountUser adds a user to Incapsula Account
func (c *Client) AddAccountUser(accountID int, email, firstName, lastName string, roleIds []interface{}, approvedIps []interface{}) (*UserApisResponse, error) {
	log.Printf("[INFO] Adding Incapsula account user for email: %s (account ID %d)\n", email, accountID)

	listRoles := make([]int, len(roleIds))
	for i, v := range roleIds {
		listRoles[i] = v.(int)
	}

	listApprovedIps := make([]string, len(approvedIps))
	for i, v := range approvedIps {
		listApprovedIps[i] = v.(string)
	}

	userAddReq := UserAddReq{UserEmail: email, RoleIds: listRoles, FirstName: firstName, LastName: lastName, ApprovedIps: listApprovedIps}

	userJSON, err := json.Marshal(userAddReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	endpointUserAdd := endpointUserOperationNew
	operation := CreateAccountUser
	accountStatusResponse, err := c.AccountStatus(accountID, ReadAccount)
	if accountStatusResponse != nil && accountStatusResponse.AccountType == "Sub Account" {
		endpointUserAdd = endpointUserOperationNew + "/" + email
		operation = CreateSubAccountUser
	}

	reqURL := fmt.Sprintf("%s/%s?caid=%d", c.config.BaseURLAPI, endpointUserAdd, accountID)
	log.Printf("[INFO] Values: %s\n", userJSON)
	log.Printf("[INFO] Req: %s\n", reqURL)
	log.Printf("[INFO] json: %s\n", userJSON)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, userJSON, operation)

	if err != nil {
		return nil, fmt.Errorf("Error adding user email %s: %s", email, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add user JSON response: %s\n", string(responseBody))

	// Look at the response status code from Incapsula
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when adding User %s: %s", resp.StatusCode, email, string(responseBody))
	}

	// Parse the JSON
	var userAddResponse UserApisResponse
	err = json.Unmarshal(responseBody, &userAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add user JSON response for email %s: %s", email, err)
	}

	log.Printf("[INFO] ResponseStruct : %+v\n", userAddResponse)
	return &userAddResponse, nil
}

// GetAccountUser gets the Incapsula user status
func (c *Client) GetAccountUser(accountID int, email string) (*UserApisResponse, error) {
	log.Printf("[INFO] Getting Incapsula user status for email id: %s\n", email)

	// Get to Incapsula
	reqURL := fmt.Sprintf("%s/%s/%s?caid=%d", c.config.BaseURLAPI, endpointUserOperationNew, email, accountID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAccountUser)

	if err != nil {
		return nil, fmt.Errorf("Error getting user %s: %s", email, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula user status JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when getting User %s: %s", resp.StatusCode, email, string(responseBody))
	}

	// Parse the JSON
	var userStatusResponse UserApisResponse
	err = json.Unmarshal(responseBody, &userStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing user status JSON response for user id %s: %s", email, err)
	}

	log.Printf("[INFO] ResponseStruct : %+v\n", userStatusResponse)
	return &userStatusResponse, nil
}

// UpdateAccountUser User Roles
// Pass nil for roleIds or approvedIps to leave them unchanged (for PATCH semantics)
func (c *Client) UpdateAccountUser(accountID int, email string, roleIds []interface{}, approvedIps []interface{}) (*UserApisUpdateResponse, error) {
	log.Printf("[INFO] Update Incapsula User for email: %s (account ID %d)\n", email, accountID)

	userUpdateReq := UserUpdateReq{}

	// Only include roleIds if provided (not nil)
	if roleIds != nil {
		listRoles := make([]int, len(roleIds))
		for i, v := range roleIds {
			listRoles[i] = v.(int)
		}
		userUpdateReq.RoleIds = &listRoles
	}

	// Only include approvedIps if provided (not nil)
	if approvedIps != nil {
		listApprovedIps := make([]string, len(approvedIps))
		for i, v := range approvedIps {
			listApprovedIps[i] = v.(string)
		}
		userUpdateReq.ApprovedIps = &listApprovedIps
	}

	userJSON, err := json.Marshal(userUpdateReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	reqURL := fmt.Sprintf("%s/%s/%s?caid=%d", c.config.BaseURLAPI, endpointUserOperationNew, email, accountID)

	log.Printf("[INFO] Req: %s\n", reqURL)
	log.Printf("[INFO] json: %s\n", userJSON)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPatch, reqURL, userJSON, UpdateAccountUser)

	if err != nil {
		return nil, fmt.Errorf("Error updating user email %s: %s", email, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update user JSON response: %s\n", string(responseBody))

	// Look at the response status code from Incapsula
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating User %s: %s", resp.StatusCode, email, string(responseBody))
	}

	// Parse the JSON
	var userUpdateResponse UserApisUpdateResponse
	err = json.Unmarshal(responseBody, &userUpdateResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update user JSON response for email %s: %s", email, err)
	}

	log.Printf("[INFO] ResponseStruct : %+v\n", userUpdateResponse)
	return &userUpdateResponse, nil
}

// DeleteAccountUser deletes a user from Incapsula
func (c *Client) DeleteAccountUser(accountID int, email string) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type UserDeleteResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	log.Printf("[INFO] Deleting Incapsula user: %s, account Id: %d\n", email, accountID)

	// Delete form to Incapsula

	reqURL := fmt.Sprintf("%s/%s/%s?caid=%d", c.config.BaseURLAPI, endpointUserOperationNew, email, accountID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteAccountUser)

	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting USER: %s %s", email, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete user JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting User %s: %s", resp.StatusCode, email, string(responseBody))
	}

	// Parse the JSON
	var userDeleteResponse UserDeleteResponse
	err = json.Unmarshal(responseBody, &userDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete user JSON response for user %s : %s", email, err)
	}

	return nil
}
