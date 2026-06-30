package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type AbpPreflightStatus struct {
	PreflightId string `json:"preflight_id"`
	CanPublish  bool   `json:"can_publish"`
}

type AbpPreflight struct {
	Id        string `json:"id"`
	AccountId string `json:"account_id"`
}

func (c *Client) AbpAccountPreflightUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/preflight", c.config.BaseURLAPI, accountId)
}

func (c *Client) AbpPreflightStatusUrl(preflightId string) string {
	return fmt.Sprintf("%s/v1/preflight/%s/status", c.config.BaseURLAPI, preflightId)
}

func (c *Client) AbpPreflightPublishUrl(preflightId string) string {
	return fmt.Sprintf("%s/v1/preflight/%s/publish", c.config.BaseURLAPI, preflightId)
}

const getPreflightStatusAction = "reading preflight status"

func (c *Client) GetAbpPreflightStatus(preflightId string) (*AbpPreflightStatus, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Reading ABP preflight status for Preflight ID %s\n", preflightId)

	reqURL := c.AbpPreflightStatusUrl(preflightId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAbpPreflightStatus)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightStatusAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightStatusAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] abp %s JSON response: %s\n", getPreflightStatusAction, string(responseBody))

	if resp.StatusCode != http.StatusOK {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightStatusAction, nil, responseBody))
		return nil, diags
	}

	var status AbpPreflightStatus
	if err = json.Unmarshal(responseBody, &status); err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightStatusAction, &err, responseBody))
		return nil, diags
	}

	return &status, diags
}

const createAbpPreflightAction = "creating preflight"

func (c *Client) CreateAbpPreflight(accountId string) (*AbpPreflight, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Creating ABP preflight for Account ID %s\n", accountId)

	reqURL := c.AbpAccountPreflightUrl(accountId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, CreateAbpPreflight)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPreflightAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPreflightAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] Incapsula POST %s preflight JSON response: %s\n", createAbpPreflightAction, string(responseBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPreflightAction, nil, responseBody))
		return nil, diags
	}

	var preflight AbpPreflight
	if err = json.Unmarshal(responseBody, &preflight); err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(createAbpPreflightAction, &err, responseBody))
		return nil, diags
	}

	return &preflight, diags
}

const publishPreflightAction = "publishing preflight"

func (c *Client) PublishAbpPreflight(preflightId string) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[INFO] Publishing ABP preflight ID %s\n", preflightId)

	reqURL := c.AbpPreflightPublishUrl(preflightId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, PublishAbpPreflight)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(publishPreflightAction, &err, nil))
		return diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(publishPreflightAction, &err, responseBody))
		return diags
	}

	log.Printf("[DEBUG] Incapsula POST %s JSON response: %s\n", publishPreflightAction, string(responseBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		diags = append(diags, httpSourcedErrorDiagnostic(publishPreflightAction, nil, responseBody))
		return diags
	}

	return diags
}
