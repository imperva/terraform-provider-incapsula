package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const endPointNotificationCenterPolicy = "notification-settings/v3/policies/"

type AssetDto struct {
	AssetType string `json:"assetType"`
	AssetId   int    `json:"assetId"`
}

type SubAccountDTO struct {
	SubAccountId int `json:"subAccountId"`
}

type SubAccountPolicyInfo struct {
	ApplyToNewSubAccounts string          `json:"applyToNewSubAccounts"`
	SubAccountList        []SubAccountDTO `json:"subAccountList"`
}

type RecipientDto struct {
	RecipientType string `json:"recipientType"`
	Id            int    `json:"id,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
}

type NotificationChannelEmailDto struct {
	RecipientToList []RecipientDto `json:"recipientToList"`
	ChannelType     string         `json:"channelType"`
}

type NotificationPolicyFullDto struct {
	PolicyId                int                           `json:"policyId,omitempty"`
	AccountId               int                           `json:"accountId"`
	PolicyName              string                        `json:"policyName"`
	Status                  string                        `json:"status"`
	SubCategory             string                        `json:"subCategory"`
	NotificationChannelList []NotificationChannelEmailDto `json:"notificationChannelList"`
	AssetList               []AssetDto                    `json:"assetList"`
	ApplyToNewAssets        string                        `json:"applyToNewAssets"`
	PolicyType              string                        `json:"policyType"`
	SubAccountPolicyInfo    SubAccountPolicyInfo          `json:"subAccountPolicyInfo"`
}

type NotificationPolicy struct {
	Data NotificationPolicyFullDto `json:"data"`
}

func (c *Client) AddNotificationCenterPolicy(notificationPolicyFullDto *NotificationPolicyFullDto) (*NotificationPolicy, error) {
	notificationPolicy := NotificationPolicy{
		Data: *notificationPolicyFullDto,
	}

	params := getAccountIdRequestParams(notificationPolicyFullDto.AccountId)
	reqURL := getRequestUrl(c)
	policyJSON, err := json.Marshal(notificationPolicy)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal NotificationCenterPolicy: %s ", err)
	}

	log.Printf("[DEBUG] Add NotificationCenterPolicy with params %s and JSON request: %s\n", params, string(policyJSON))
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, policyJSON, params)
	log.Printf("[DEBUG] client_notification_center_policy Post rest response:\n%+v", resp)
	if err != nil {
		return nil, fmt.Errorf("Error from NotificationCenter service when adding policy: %s ", err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Add NotificationCenterPolicy JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from NotificationCenter service when adding policy: %s ", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var policy NotificationPolicy
	err = json.Unmarshal(responseBody, &policy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing NotificationCenterPolicy JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &policy, nil

}

func (c *Client) UpdateNotificationCenterPolicy(notificationPolicyFullDto *NotificationPolicyFullDto) (*NotificationPolicy, error) {
	notificationPolicy := NotificationPolicy{
		Data: *notificationPolicyFullDto,
	}

	reqURL := getRequestUrlWithId(c, notificationPolicy.Data.PolicyId)
	policyJSON, err := json.Marshal(notificationPolicy)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal NotificationCenterPolicy: %s ", err)
	}
	params := getAccountIdRequestParams(notificationPolicyFullDto.AccountId)

	log.Printf("[DEBUG] Update NotificationCenterPolicy JSON request: %s\n", string(policyJSON))
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPut, reqURL, policyJSON, params)
	log.Printf("[DEBUG] client_notification_center_policy Put rest response:\n%+v", resp)
	if err != nil {
		return nil, fmt.Errorf("Error from NotificationCenter service when updateing policy: %s ", err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Update NotificationCenterPolicy JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from NotificationCenter service when updateing policy: %s ", resp.StatusCode, string(responseBody))
	}

	// Parse the JSON
	var policy NotificationPolicy
	err = json.Unmarshal(responseBody, &policy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing NotificationCenterPolicy JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	return &policy, nil
}

func (c *Client) DeleteNotificationCenterPolicy(policyId int, accountId int) error {
	log.Printf("[INFO] Deleting NotificationCenterPolicy with ID %d ", policyId)
	requestUrl := getRequestUrlWithId(c, policyId)
	params := getAccountIdRequestParams(accountId)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, requestUrl, nil, params)
	log.Printf("[DEBUG] client_notification_center_policy Delete rest response:\n%+v", resp)
	if err != nil {
		return fmt.Errorf("Error from NotificationCenterPolicy service when deleting Policy with Id %d: %s ", policyId, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] NotificationCenter Delete policy JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from NotificationCenter service when deleting policy with Id %d: %s ", resp.StatusCode, policyId, string(responseBody))
	}

	return nil
}

func getRequestUrl(c *Client) string {
	requestUrl := fmt.Sprintf("%s/%s", c.config.BaseURLAPI, endPointNotificationCenterPolicy)

	return requestUrl
}

func getRequestUrlWithId(c *Client, policyId int) string {
	requestUrl := fmt.Sprintf("%s/%d", getRequestUrl(c), policyId)

	return requestUrl
}

func (c *Client) GetNotificationCenterPolicy(policyId int, accountId int) (*NotificationPolicy, error) {
	log.Printf("[INFO] Getting  NotificationCenterPolicy with policyId: %d and accountId: %d", policyId, accountId)
	requestUrl := getRequestUrlWithId(c, policyId)

	params := getAccountIdRequestParams(accountId)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, requestUrl, nil, params)
	log.Printf("[DEBUG] client_notification_center_policy Get rest response:\n%+v", resp)
	if err != nil {
		return nil, fmt.Errorf("Error from NotificationCenter service when reading policy with Id %d: %s ", policyId, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] NotificationCenter Read policy JSON response: %s\n", string(responseBody))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from NotificationCenter service when reading policy for ID %d: %s ", resp.StatusCode, policyId, string(responseBody))
	}

	var notificationCenterPolicy NotificationPolicy
	err = json.Unmarshal(responseBody, &notificationCenterPolicy)
	if err != nil {
		return nil, fmt.Errorf("Error parsing NotificationCenterPolicy JSON response with policy ID %d: %s\nresponse: %s", policyId, err, string(responseBody))
	}

	return &notificationCenterPolicy, nil
}

func getAccountIdRequestParams(accountId int) map[string]string {
	var params = map[string]string{}
	params["caid"] = strconv.Itoa(accountId)

	return params
}
