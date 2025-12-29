package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Endpoints (unexported consts)
const endpointRole = "user-management/v1/roles"
const endpointRoleAdd = endpointRole
const endpointRoleGet = endpointRole
const endpointRoleUpdate = endpointRole
const endpointRoleDelete = endpointRole
const endpointAbilitiesGet = "user-management/v1/abilities/accounts"
const endpointAccountRolesGet = endpointRole

type RoleAbility struct {
	AbilityKey              string `json:"abilityKey"`
	AbilityDisplayName      string `json:"abilityDisplayName"`
	IsRelevantForSubAccount bool   `json:"isRelevantForSubAccount"`
}

type UserAssignment struct {
	UserEmail string `json:"userEmail"`
	AccountId int    `json:"accountId"`
}

type RoleDetailsBasicDTO struct {
	RoleName        string   `json:"roleName"`
	RoleDescription string   `json:"roleDescription"`
	RoleAbilities   []string `json:"roleAbilities"`
}

type RoleDetailsCreateDTO struct {
	AccountId int `json:"accountId"`
	RoleDetailsBasicDTO
}

// RoleDetailsDTO - Same DTO for: GET response, POST request, and POST response
type RoleDetailsDTO struct {
	RoleId          int              `json:"roleId"`
	RoleName        string           `json:"roleName"`
	RoleDescription string           `json:"roleDescription"`
	AccountId       int              `json:"accountId"`
	AccountName     string           `json:"accountName"`
	RoleAbilities   []RoleAbility    `json:"roleAbilities"`
	UserAssignment  []UserAssignment `json:"userAssignment"`
	UpdateDate      string           `json:"updateDate"`
	IsEditable      bool             `json:"isEditable"`
	ErrorCode       int              `json:"errorCode"`
	Description     string           `json:"description"`
}

// AddAccountRole Adds an Account Role to be managed by Incapsula
func (c *Client) AddAccountRole(requestDTO RoleDetailsCreateDTO) (*RoleDetailsDTO, error) {
	log.Printf("[INFO] Adding Incapsula account role %s (account ID %d)\n", requestDTO.RoleName, requestDTO.AccountId)

	roleJSON, err := json.Marshal(requestDTO)
	log.Printf("[INFO]  roleJSON: %v\n", string(roleJSON))

	reqURL := fmt.Sprintf("%s/%s", c.config.BaseURLAPI, endpointRoleAdd)
	log.Printf("[INFO]  reqURL: %v\n", reqURL)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, roleJSON, CreateAccountRole)
	if err != nil {
		return nil, fmt.Errorf("Error adding account role %s: %s", requestDTO.RoleName, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add account role JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var roleResponse RoleDetailsDTO
	err = json.Unmarshal(responseBody, &roleResponse)

	if err != nil {
		return nil, fmt.Errorf("Error parsing add account role JSON response: %s", err)
	}

	// Look at the response status code from Incapsula
	if roleResponse.ErrorCode != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding account role: %s", string(responseBody))
	}

	return &roleResponse, nil
}

// GetAccountRole - Retrieve the Account Role for a given role ID
func (c *Client) GetAccountRole(roleId int) (*RoleDetailsDTO, error) {
	log.Printf("[INFO] Getting Account Role (Id: %d)\n", roleId)

	// Get request to Incapsula
	reqURL := fmt.Sprintf("%s/%s/%d", c.config.BaseURLAPI, endpointRoleGet, roleId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAccountRole)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Account Role request for role with id %d: %s", roleId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Account Role JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO RoleDetailsDTO
	err = json.Unmarshal(responseBody, &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Role JSON response for role with Id: %d %s\nresponse: %s", roleId, err, string(responseBody))
	}

	return &responseDTO, nil
}

// UpdateAccountRole - Update the Account Role for a given role ID
func (c *Client) UpdateAccountRole(roleId int, accountId int, requestDTO RoleDetailsBasicDTO) (*RoleDetailsDTO, error) {
	log.Printf("[INFO] Updating Incapsula account role (Id: %d, Account Id: %d)\n", roleId, accountId)

	log.Printf("[INFO]  requestDTO: %+v\n", requestDTO)

	roleJSON, err := json.Marshal(requestDTO)
	log.Printf("[INFO]  roleJSON: %v\n", string(roleJSON))

	reqURL := fmt.Sprintf("%s/%s/%d", c.config.BaseURLAPI, endpointRoleUpdate, roleId)
	log.Printf("[INFO]  reqURL: %v\n", reqURL)

	params := GetRequestParamsWithCaid(accountId)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, roleJSON, params, UpdateAccountRole)
	if err != nil {
		return nil, fmt.Errorf("Error updating account role with Id %d: %s", roleId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update account role JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var roleResponse RoleDetailsDTO
	err = json.Unmarshal(responseBody, &roleResponse)

	if err != nil {
		return nil, fmt.Errorf("Error parsing update account role JSON response: %s", err)
	}

	// Look at the response status code from Incapsula
	if roleResponse.ErrorCode != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating account role: %s", string(responseBody))
	}

	return &roleResponse, nil
}

// DeleteAccountRole - Delete the Account Role for a given role ID
func (c *Client) DeleteAccountRole(roleId int) error {
	log.Printf("[INFO] Delete Account Role (Id: %d)\n", roleId)

	// Get request to Incapsula
	reqURL := fmt.Sprintf("%s/%s/%d", c.config.BaseURLAPI, endpointRoleDelete, roleId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteAccountRole)
	if err != nil {
		return fmt.Errorf("Error executing delete Account Role request for role with id %d: %s", roleId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Account Role JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO RoleDetailsDTO
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return fmt.Errorf("Error parsing Account Role JSON response for role with Id: %d %s\nresponse: %s", roleId, err, string(responseBody))
	}

	return nil
}

// GetAccountAbilities - Retrieve the Account Abilities for a given account ID
func (c *Client) GetAccountAbilities(accountId int) (*[]RoleAbility, error) {
	log.Printf("[INFO] Getting Account Abilities for account Id: %d\n", accountId)

	// Get request to Incapsula
	reqURL := fmt.Sprintf("%s/%s/%d", c.config.BaseURLAPI, endpointAbilitiesGet, accountId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAccountAbilities)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Account Abilities request for account with id %d: %s", accountId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Account Abilities JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var roleAbility []RoleAbility
	err = json.Unmarshal(responseBody, &roleAbility)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Abilities JSON response for account with id: %d %s\nresponse: %s", accountId, err, string(responseBody))
	}

	return &roleAbility, nil
}

// GetAccountRoles - Retrieve all the Roles for a given Account ID
func (c *Client) GetAccountRoles(accountId int) (*[]RoleDetailsDTO, error) {
	log.Printf("[INFO] Getting Account Roles (Account Id: %d)\n", accountId)

	// Get request to Incapsula
	reqURL := fmt.Sprintf("%s/%s?accountId=%d", c.config.BaseURLAPI, endpointAccountRolesGet, accountId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAccountRoles)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Account Roles request for account with id %d: %s", accountId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Account Roles JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO []RoleDetailsDTO
	err = json.Unmarshal(responseBody, &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Roles JSON response for account with Id: %d %s\nresponse: %s", accountId, err, string(responseBody))
	}

	return &responseDTO, nil
}
