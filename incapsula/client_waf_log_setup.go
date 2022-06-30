package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
)

// Endpoints
const endpointCreateDefault = "accounts/setDefaultSiemStorage"
const endpointCreateS3 = "accounts/setAmazonSiemStorage"
const endpointCreateSFTP = "accounts/setSftpSiemStorage"
const endpointTestCreateS3 = "accounts/testS3Connection"
const endpointTestCreateSFTP = "accounts/testSftpConnection"
const endpointWAFLogsActivate = "waf-log-setup/activate"
const endpointWAFLogsChangeStatus = "waf-log-setup/change/status"

const saveOnSuccess = true

type WAFLogSetup struct {
	*WAFLogSetupPayload
}

type WAFLogSetupResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
}

type WAFLogSetupResponseActivate struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	ConfigID   int    `json:"logs_collector_config_id"`
}

// WAFLogSetupPayload contains the S3/SFTP payload for Incapsula WAF Log Setup creation
type WAFLogSetupPayload struct {
	AccountID         int    `json:"account_id,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	BucketName        string `json:"s3_bucket_name,omitempty"`
	AccessKey         string `json:"s3_access_key,omitempty"`
	SecretKey         string `json:"s3_secret_key,omitempty"`
	Host              string `json:"sftp_host,omitempty"`
	UserName          string `json:"sftp_user_name,omitempty"`
	Password          string `json:"sftp_password,omitempty"`
	DestinationFolder string `json:"sftp_destination_folder,omitempty"`
}

func (c *Client) activateAndUpdateStatus(wafLogSetupPayload *WAFLogSetupPayload) error {

	//******************************************
	//Activate WAF Logs
	//******************************************

	valuesActivate := url.Values{
		"account_id": {strconv.Itoa(wafLogSetupPayload.AccountID)},
	}

	respActivate, errActivate := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointWAFLogsActivate), valuesActivate, ActivateWAFLogSetup)
	if errActivate != nil {
		return fmt.Errorf("Error activating WAF Log Setup for account  %d: %s", wafLogSetupPayload.AccountID, errActivate)
	}

	// Read the body
	defer respActivate.Body.Close()
	responseBodyActivate, err := ioutil.ReadAll(respActivate.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula activate WAF Log Setup JSON response: %s\n", string(responseBodyActivate))

	// Parse the JSON
	var wafLogSetupResponseActivate WAFLogSetupResponseActivate
	err = json.Unmarshal([]byte(responseBodyActivate), &wafLogSetupResponseActivate)
	if err != nil {
		return fmt.Errorf("Error parsing activate WAF Log Setup JSON response for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}

	//******************************************
	//Update Status WAF Logs (Active/Suspended)
	//******************************************

	var logsConfigNewStatus string
	if wafLogSetupPayload.Enabled {
		logsConfigNewStatus = "ACTIVE"
	} else {
		logsConfigNewStatus = "SUSPENDED"
	}

	log.Printf("[DEBUG] accountId %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] configId %d\n", wafLogSetupResponseActivate.ConfigID)
	log.Printf("[DEBUG] logsConfigNewStatus %s\n", logsConfigNewStatus)

	valuesStatus := url.Values{
		"account_id":             {strconv.Itoa(wafLogSetupPayload.AccountID)},
		"config_id":              {strconv.Itoa(wafLogSetupResponseActivate.ConfigID)},
		"logs_config_new_status": {logsConfigNewStatus},
	}

	respStatus, errStatus := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointWAFLogsChangeStatus), valuesStatus, UpdateStatusWAFLogSetup)
	if errActivate != nil {
		return fmt.Errorf("Error changing WAF Log Setup status for account  %d: %s", wafLogSetupPayload.AccountID, errStatus)
	}

	// Read the body
	defer respStatus.Body.Close()
	responseBodyStatus, err := ioutil.ReadAll(respStatus.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula change WAF Log Setup status JSON response: %s\n", string(responseBodyStatus))

	// Parse the JSON
	var wafLogSetupResponseStatus WAFLogSetupResponse
	err = json.Unmarshal([]byte(responseBodyStatus), &wafLogSetupResponseStatus)
	if err != nil {
		return fmt.Errorf("Error parsing change WAF Log Setup status JSON response for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}
	return nil
}

// AddWAFLogSetupS3 adds a WAFLogSetup to be managed by Incapsula
func (c *Client) AddWAFLogSetupS3(wafLogSetupPayload *WAFLogSetupPayload) (*WAFLogSetupResponse, error) {
	log.Printf("[INFO] Adding Incapsula WAF S3 Log Setup for account: %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] accountId %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] enabled %s\n", strconv.FormatBool(wafLogSetupPayload.Enabled))
	log.Printf("[DEBUG] accessKey %s\n", wafLogSetupPayload.AccessKey)
	log.Printf("[DEBUG] bucketName %s\n", wafLogSetupPayload.BucketName)
	//log.Printf("[DEBUG] secretKey %s\n", wafLogSetupPayload.SecretKey)

	err := c.activateAndUpdateStatus(wafLogSetupPayload)
	if err != nil {
		return nil, err
	}

	//******************************************
	//Setup Connection WAF Logs
	//************************************

	values := url.Values{
		"account_id":  {strconv.Itoa(wafLogSetupPayload.AccountID)},
		"bucket_name": {wafLogSetupPayload.BucketName},
		"access_key":  {wafLogSetupPayload.AccessKey},
		"secret_key":  {wafLogSetupPayload.SecretKey},
	}

	var endpoint string
	if v := os.Getenv("TESTING_PROFILE"); v != "" {
		endpoint = endpointCreateS3
	} else {
		values["save_on_success"] = make([]string, 1)
		values["save_on_success"][0] = fmt.Sprint(saveOnSuccess)
		endpoint = endpointTestCreateS3
	}
	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint), values, CreateWAFLogSetup)
	if err != nil {
		return nil, fmt.Errorf("Error creating S3 WAF Log Setup for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula create S3 WAF Log Setup JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var wafLogSetupResponse WAFLogSetupResponse
	err = json.Unmarshal([]byte(responseBody), &wafLogSetupResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing create S3 WAF Log Setup JSON response for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}

	// Look at the response status code from Incapsula
	if wafLogSetupResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when creating S3 WAF Log Setup for account %d: %s", wafLogSetupPayload.AccountID, string(responseBody))
	}

	return &wafLogSetupResponse, nil
}

// AddWAFLogSetupSFTP adds a WAFLogSetup to be managed by Incapsula
func (c *Client) AddWAFLogSetupSFTP(wafLogSetupPayload *WAFLogSetupPayload) (*WAFLogSetupResponse, error) {
	log.Printf("[INFO] Adding Incapsula WAF S3 Log Setup for account: %d\n", wafLogSetupPayload.AccountID)

	err := c.activateAndUpdateStatus(wafLogSetupPayload)
	if err != nil {
		return nil, err
	}

	values := url.Values{
		"account_id":         {strconv.Itoa(wafLogSetupPayload.AccountID)},
		"host":               {wafLogSetupPayload.Host},
		"user_name":          {wafLogSetupPayload.UserName},
		"password":           {wafLogSetupPayload.Password},
		"destination_folder": {wafLogSetupPayload.DestinationFolder},
	}

	log.Printf("[DEBUG] accountId %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] host %s\n", wafLogSetupPayload.Host)
	log.Printf("[DEBUG] destinationFolder %s\n", wafLogSetupPayload.DestinationFolder)
	log.Printf("[DEBUG] userName %s\n", wafLogSetupPayload.UserName)
	//log.Printf("[DEBUG] password %s\n", wafLogSetupPayload.Password)

	var endpoint string
	if v := os.Getenv("TESTING_PROFILE"); v != "" {
		endpoint = endpointCreateSFTP
	} else {
		values["save_on_success"] = make([]string, 1)
		values["save_on_success"][0] = fmt.Sprint(saveOnSuccess)
		endpoint = endpointTestCreateSFTP
	}

	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint), values, CreateWAFLogSetup)
	if err != nil {
		return nil, fmt.Errorf("Error creating SFTP WAF Log Setup for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula create SFTP WAF Log Setup JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var wafLogSetupResponse WAFLogSetupResponse
	err = json.Unmarshal([]byte(responseBody), &wafLogSetupResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing create SFTP WAF Log Setup JSON response for account  %d: %s", wafLogSetupPayload.AccountID, err)
	}

	// Look at the response status code from Incapsula
	if wafLogSetupResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when creating SFTP WAF Log Setup for account %d: %s", wafLogSetupPayload.AccountID, string(responseBody))
	}

	return &wafLogSetupResponse, nil
}

// AddWAFLogSetupDefault turns WAF Log Setup to Default with enablement option ACTIVE/SUSPENDED
func (c *Client) AddWAFLogSetupDefault(wafLogSetupPayload *WAFLogSetupPayload) (*WAFLogSetupResponse, error) {
	log.Printf("[INFO] Adding Incapsula WAF Default Log Setup for account: %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] accountId %d\n", wafLogSetupPayload.AccountID)
	log.Printf("[DEBUG] enabled %s\n", strconv.FormatBool(wafLogSetupPayload.Enabled))

	err := c.activateAndUpdateStatus(wafLogSetupPayload)
	if err != nil {
		return nil, err
	}

	//******************************************
	//Setup Connection WAF Logs
	//************************************

	return c.RestoreWAFLogSetupDefault(wafLogSetupPayload.AccountID)
}

// RestoreWAFLogSetupDefault restores WAF Log Setup to Default for a given account ID
func (c *Client) RestoreWAFLogSetupDefault(accountID int) (*WAFLogSetupResponse, error) {
	log.Printf("[INFO] Restoring Incapsula WAF Log Setup to default for account: %d\n", accountID)

	values := url.Values{
		"account_id": {strconv.Itoa(accountID)},
	}

	resp, err := c.PostFormWithHeaders(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointCreateDefault), values, DeleteWAFLogSetup)
	if err != nil {
		return nil, fmt.Errorf("Error restoring WAF Log Setup to default for account  %d: %s", accountID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula create default WAF Log Setup JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var wafLogSetupResponse WAFLogSetupResponse
	err = json.Unmarshal([]byte(responseBody), &wafLogSetupResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing create default WAF Log Setup JSON response for account  %d: %s", accountID, err)
	}

	// Look at the response status code from Incapsula
	if wafLogSetupResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when creating default WAF Log Setup for account %d: %s", accountID, string(responseBody))
	}

	return &wafLogSetupResponse, nil
}
