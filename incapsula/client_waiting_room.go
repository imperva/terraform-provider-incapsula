package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ThresholdSettings struct {
	EntranceRateEnabled         bool `json:"isEntranceRateEnabled"`
	EntranceRateThreshold       int  `json:"entranceRateThreshold,omitempty"`
	ConcurrentSessionsEnabled   bool `json:"isConcurrentSessionsEnabled"`
	ConcurrentSessionsThreshold int  `json:"concurrentSessionsThreshold,omitempty"`
	InactivityTimeout           int  `json:"inactivityTimeout,omitempty"`
}

type WaitingRoomDTO struct {
	Id                      int64             `json:"id,omitempty"`
	AccountId               int64             `json:"accountId,omitempty"`
	Name                    string            `json:"name"`
	Description             string            `json:"description,omitempty"`
	Enabled                 bool              `json:"enabled"`
	Filter                  string            `json:"filter,omitempty"`
	HtmlTemplateBase64      string            `json:"htmlTemplateBase64,omitempty"`
	CreatedAt               int64             `json:"createdAt,omitempty"`
	LastModifiedAt          int64             `json:"lastModifiedAt,omitempty"`
	LastModifiedBy          string            `json:"lastModifiedBy,omitempty"`
	Mode                    string            `json:"mode,omitempty"`
	BotsActionInQueuingMode string            `json:"botsActionInQueuingMode,omitempty"`
	QueueInactivityTimeout  int               `json:"queueInactivityTimeout,omitempty"`
	HidePositionInLine      bool              `json:"hidePositionInLine"`
	ThresholdSettings       ThresholdSettings `json:"thresholdSettings"`
}

type WaitingRoomDTOResponse struct {
	Data   []WaitingRoomDTO `json:"data"`
	Errors []APIErrors      `json:"errors"`
}

func (c *Client) CreateWaitingRoom(accountId string, siteID string, waitingRoom *WaitingRoomDTO) (*WaitingRoomDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Creating Waiting Room for Site ID %s\n", siteID)

	waitingRoomJSON, err := json.Marshal(waitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure generating Waiting Room create request",
			Detail:   fmt.Sprintf("Failed to JSON marshal WaitingRoom: %s", err.Error()),
		})
		return nil, diags
	}

	// Dump JSON
	log.Printf("[DEBUG] Waiting Room payload: %s\n", string(waitingRoomJSON))

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/waiting-room-settings/v3/sites/%s/waiting-rooms", c.config.BaseURLAPI, siteID)
	if accountId != "" {
		reqURL += "?caid=" + accountId
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, waitingRoomJSON, CreateWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Creating Waiting Room",
			Detail:   fmt.Sprintf("Error from Incapsula service when creating Waiting Room for Site ID %s: %s", siteID, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Create Waiting Room JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 201 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Creating Waiting Room",
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when creating Waiting Room for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody)),
		})
	}

	// Parse the JSON
	var newWaitingRoom WaitingRoomDTOResponse
	err = json.Unmarshal([]byte(responseBody), &newWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure parsing Waiting Room create response",
			Detail:   fmt.Sprintf("Error parsing Waiting Room JSON response for Site ID %s: %s\nresponse: %s", siteID, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &newWaitingRoom, diags
}

func (c *Client) ReadWaitingRoom(accountId string, siteID string, waitingRoomID int64) (*WaitingRoomDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Getting Incapsula Waiting Room %d for Site ID %s\n", waitingRoomID, siteID)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", c.config.BaseURLAPI, siteID, waitingRoomID)
	if accountId != "" {
		reqURL += "?caid=" + accountId
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure sending Waiting Room read request",
			Detail:   fmt.Sprintf("Error from Incapsula service when reading Waiting Room %d for Site ID %s: %s", waitingRoomID, siteID, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Read Waiting Room JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Reading Waiting Room",
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when reading Waiting Room %d for Site ID %s: %s", resp.StatusCode, waitingRoomID, siteID, string(responseBody)),
		})
	}

	// Parse the JSON
	var waitingRoom WaitingRoomDTOResponse
	err = json.Unmarshal([]byte(responseBody), &waitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure parsing Waiting Room read response",
			Detail:   fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s: %s\nresponse: %s", waitingRoomID, siteID, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &waitingRoom, diags
}

func (c *Client) UpdateWaitingRoom(accountId string, siteID string, waitingRoomID int64, waitingRoom *WaitingRoomDTO) (*WaitingRoomDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Updating Incapsula Waiting Room %d for Site ID %s\n", waitingRoomID, siteID)

	waitingRoomJSON, err := json.Marshal(waitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure generating Waiting Room update request",
			Detail:   fmt.Sprintf("Failed to JSON marshal WaitingRoom: %s", err.Error()),
		})
		return nil, diags
	}

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", c.config.BaseURLAPI, siteID, waitingRoomID)
	if accountId != "" {
		reqURL += "?caid=" + accountId
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, waitingRoomJSON, UpdateWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Updating Waiting Room",
			Detail:   fmt.Sprintf("Error from Incapsula service when updating Waiting Room %d for Site ID %s: %s", waitingRoomID, siteID, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Waiting Room JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Updating Waiting Room",
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when updating Waiting Room %d for Site ID %s: %s", resp.StatusCode, waitingRoomID, siteID, string(responseBody)),
		})
	}

	// Parse the JSON
	var updatedWaitingRoom WaitingRoomDTOResponse
	err = json.Unmarshal([]byte(responseBody), &updatedWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure parsing Waiting Room create response",
			Detail:   fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s: %s\nresponse: %s", waitingRoomID, siteID, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &updatedWaitingRoom, diags
}

func (c *Client) DeleteWaitingRoom(accountId string, siteID string, waitingRoomID int64) (*WaitingRoomDTOResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Deleting Incapsula Waiting Room %d for Site ID %s\n", waitingRoomID, siteID)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", c.config.BaseURLAPI, siteID, waitingRoomID)
	if accountId != "" {
		reqURL += "?caid=" + accountId
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeleteWaitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure sending Waiting Room delete request",
			Detail:   fmt.Sprintf("Error from Incapsula service when deleting Waiting Room %d for Site ID %s: %s", waitingRoomID, siteID, err.Error()),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Waiting Room JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure Deleting Waiting Room",
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when deleting Waiting Room %d for Site ID %s: %s", resp.StatusCode, waitingRoomID, siteID, string(responseBody)),
		})
	}

	// Parse the JSON
	var waitingRoom WaitingRoomDTOResponse
	err = json.Unmarshal([]byte(responseBody), &waitingRoom)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failure parsing Waiting Room delete response",
			Detail:   fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s: %s\nresponse: %s", waitingRoomID, siteID, err.Error(), string(responseBody)),
		})
		return nil, diags
	}

	return &waitingRoom, diags
}
