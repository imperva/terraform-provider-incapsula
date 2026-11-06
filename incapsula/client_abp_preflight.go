package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// AbpPreflightWarning is one reason why publishing a preflight will not have
// the effect the user expects.
type AbpPreflightWarning struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type AbpPreflightLink struct {
	Href string `json:"href"`
}

type AbpPreflight struct {
	Id         string                      `json:"id"`
	AccountId  string                      `json:"account_id"`
	CreatedAt  string                      `json:"created_at"`
	ModifiedAt string                      `json:"modified_at"`
	Warnings   []AbpPreflightWarning       `json:"warnings"`
	Links      map[string]AbpPreflightLink `json:"_links"`
}

// abpPreflightPublishLink is the transition the API offers only while the
// preflight can still be published.
const abpPreflightPublishLink = "create_publish"

// CanPublish reports whether the service still accepts a publish of this
// preflight. A preflight that is older than the configuration last published
// to the account can never be published again, so this is false forever once
// it turns false.
func (p *AbpPreflight) CanPublish() bool {
	_, ok := p.Links[abpPreflightPublishLink]
	return ok
}

func (c *Client) AbpAccountPreflightUrl(accountId string) string {
	return fmt.Sprintf("%s/botmanagement/v1/account/%s/preflight", c.config.BaseURLAPI, accountId)
}

func (c *Client) AbpPreflightUrl(preflightId string) string {
	return fmt.Sprintf("%s/botmanagement/v1/preflight/%s", c.config.BaseURLAPI, preflightId)
}

func (c *Client) AbpPreflightPublishUrl(preflightId string) string {
	return fmt.Sprintf("%s/botmanagement/v1/preflight/%s/publish", c.config.BaseURLAPI, preflightId)
}

const getPreflightAction = "reading preflight"

// GetAbpPreflight reads a preflight. It returns (nil, nil) when the preflight
// no longer exists: the service discards preflights as they age, so a caller
// holding an old id has to treat its disappearance as an ordinary outcome.
func (c *Client) GetAbpPreflight(preflightId string) (*AbpPreflight, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Reading ABP preflight %s\n", preflightId)

	reqURL := c.AbpPreflightUrl(preflightId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadAbpPreflight)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightAction, &err, nil))
		return nil, diags
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightAction, &err, responseBody))
		return nil, diags
	}

	log.Printf("[DEBUG] abp %s JSON response: %s\n", getPreflightAction, string(responseBody))

	if resp.StatusCode == http.StatusNotFound {
		return nil, diags
	}

	if resp.StatusCode != http.StatusOK {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightAction, nil, responseBody))
		return nil, diags
	}

	var preflight AbpPreflight
	if err = json.Unmarshal(responseBody, &preflight); err != nil {
		diags = append(diags, httpSourcedErrorDiagnostic(getPreflightAction, &err, responseBody))
		return nil, diags
	}

	return &preflight, diags
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
