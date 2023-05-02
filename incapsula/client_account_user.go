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
const endpointAccountUserAdd = "identity-management/v3/users"
const endpointSubAccountUserAdd = "identity-management/v3/users/sub-account"
const endpointUserStatus = "identity-management/v3/users"
const endpointUserDelete = "identity-management/v3/users"
const endpointUserUpdate = "identity-management/v3/roles/assignments"

// UserApisResponse contains the relevant user information when adding, getting or updating a user
type UserApisResponse struct {
	Data struct {
		UserID    string `json:"id"`
		AccountID int    `json:"accountId"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Roles     []struct {
			RoleID   int    `json:"id"`
			RoleName string `json:"name"`
		} `json:"roles"`
	} `json:"data"`
}

type UserApisUpdateResponse struct {
	Data []struct {
		UserID    string `json:"id"`
		AccountID int    `json:"accountId"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Roles     []struct {
			RoleID   int    `json:"id"`
			RoleName string `json:"name"`
		} `json:"roles"`
	} `json:"data"`
}

type UserAddReq struct {
	UserEmail string `json:"email"`
	RoleIds   []int  `json:"roleIds"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserUpdateReq struct {
	UserEmail string `json:"email"`
	AccountId int    `json:"accountId"`
	RoleIds   []int  `json:"roleIds"`
}

// AddAccountUser adds a user to Incapsula Account
func (c *Client) AddAccountUser(accountID int, email, firstName, lastName string, roleIds []interface{}) (*UserApisResponse, error) {
	log.Printf("[INFO] Adding Incapsula account user for email: %s (account ID %d)\n", email, accountID)

	listRoles := make([]int, len(roleIds))
	for i, v := range roleIds {
		listRoles[i] = v.(int)
	}

	userAddReq := UserAddReq{UserEmail: email, RoleIds: listRoles, FirstName: firstName, LastName: lastName}

	userJSON, err := json.Marshal(userAddReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	endpointUserAdd := endpointAccountUserAdd
	operation := CreateAccountUser
	accountStatusResponse, err := c.AccountStatus(accountID, ReadAccount)
	if accountStatusResponse != nil && accountStatusResponse.AccountType == "Sub Account" {
		endpointUserAdd = endpointSubAccountUserAdd
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
	reqURL := fmt.Sprintf("%s/%s?caid=%d&email=%s", c.config.BaseURLAPI, endpointUserStatus, accountID, url.QueryEscape(email))
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
func (c *Client) UpdateAccountUser(accountID int, email string, roleIds []interface{}) (*UserApisUpdateResponse, error) {
	log.Printf("[INFO] Update Incapsula User for email: %s (account ID %d)\n", email, accountID)

	listRoles := make([]int, len(roleIds))
	for i, v := range roleIds {
		listRoles[i] = v.(int)
	}

	userUpdateReq := []UserUpdateReq{{AccountId: accountID, UserEmail: email, RoleIds: listRoles}}

	userJSON, err := json.Marshal(userUpdateReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal IncapRule: %s", err)
	}

	reqURL := fmt.Sprintf("%s/%s?caid=%d", c.config.BaseURLAPI, endpointUserUpdate, accountID)

	log.Printf("[INFO] Req: %s\n", reqURL)
	log.Printf("[INFO] json: %s\n", userJSON)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, userJSON, UpdateAccountUser)

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

	reqURL := fmt.Sprintf("%s/%s?caid=%d&email=%s", c.config.BaseURLAPI, endpointUserDelete, accountID, url.QueryEscape(email))
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
