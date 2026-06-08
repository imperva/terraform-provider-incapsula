package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpSiteResourceName = "ABP Site"

type AbpSite struct {
	Id                           string        `json:"id,omitempty"`
	AccountId                    string        `json:"account_id,omitempty"`
	Name                         string        `json:"name"`
	Selectors                    []AbpSelector `json:"selectors"`
	DefaultMaxRequestsPerMinute  *int          `json:"default_max_requests_per_minute,omitempty"`
	DefaultMaxRequestsPerSession *int          `json:"default_max_requests_per_session,omitempty"`
	DefaultMaxSessionLength      *string       `json:"default_max_session_length,omitempty"`
	CreatedAt                    string        `json:"created_at,omitempty"`
	ModifiedAt                   string        `json:"modified_at,omitempty"`

	// The default selector is split off the wire payload before unmarshalling
	// completes and is re-attached on update preserving its id and policy.
	DefaultSelector *AbpSelector `json:"-"`
}

type AbpSelector struct {
	Id               string              `json:"id,omitempty"`
	PolicyId         *string             `json:"policy_id"`
	Criteria         AbpSelectorCriteria `json:"criteria"`
	AnalysisSettings AbpAnalysisSettings `json:"analysis_settings"`
}

// AbpSelectorCriteria is a oneOf — exactly one of the three fields must be set.
type AbpSelectorCriteria struct {
	Postback   *string `json:"postback,omitempty"`
	PathPrefix *string `json:"path_prefix,omitempty"`
	PathRegex  *string `json:"path_regex,omitempty"`
}

type AbpAnalysisSettings struct {
	RateLimiting                  AbpRateLimiting `json:"rate_limiting"`
	MaxRequestsPerMinute          *int            `json:"max_requests_per_minute,omitempty"`
	MaxRequestsPerSession         *int            `json:"max_requests_per_session,omitempty"`
	MaxSessionLength              *string         `json:"max_session_length,omitempty"`
	UseSiteRateLimitingParameters *bool           `json:"use_site_rate_limiting_parameters,omitempty"`
}

// AbpRateLimiting is a oneOf:
//
//	"none" | "per_site" | { "custom_scope": "<scope>" }
//
// Mode is one of "none", "per_site", "custom_scope". CustomScope is only
// meaningful when Mode == "custom_scope".
type AbpRateLimiting struct {
	Mode        string
	CustomScope string
}

const (
	AbpRateLimitingModeNone        = "none"
	AbpRateLimitingModePerSite     = "per_site"
	AbpRateLimitingModeCustomScope = "custom_scope"
)

func (r AbpRateLimiting) MarshalJSON() ([]byte, error) {
	switch r.Mode {
	case AbpRateLimitingModeCustomScope:
		return json.Marshal(map[string]any{"custom_scope": r.CustomScope})
	case AbpRateLimitingModeNone, AbpRateLimitingModePerSite:
		return json.Marshal(r.Mode)
	default:
		return nil, fmt.Errorf("invalid rate_limiting mode %q", r.Mode)
	}
}

func (r *AbpRateLimiting) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		r.Mode = s
		r.CustomScope = ""
		return nil
	}
	var obj struct {
		CustomScope string `json:"custom_scope"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("rate_limiting: expected string or {custom_scope}: %w", err)
	}
	r.Mode = AbpRateLimitingModeCustomScope
	r.CustomScope = obj.CustomScope
	return nil
}

func (c *Client) abpSiteAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/site", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpSiteUrl(siteId string) string {
	return fmt.Sprintf("%s/v1/site/%s", c.config.BaseURLAPI, siteId)
}

func isDefaultSelectorShape(s AbpSelector) bool {
	return s.Criteria.PathPrefix != nil && *s.Criteria.PathPrefix == "/"
}

// splitDefaultSelector splits the default selector (if present) from others.
// Whatever the user-managed selectors are, the backend will always append one
// default at the bottom, so we trust the shape on the trailing entry and pop it
// into `site.DefaultSelector`.
//
// If the trailing selector doesn't match the default shape, that's a site the
// user has explicitly cleared the default off. In that case we leave `Selectors`
// untouched and `DefaultSelector` nil; subsequent updates won't re-add a default.
func splitDefaultSelector(site *AbpSite) {
	n := len(site.Selectors)
	if n == 0 {
		return
	}
	last := site.Selectors[n-1]
	if isDefaultSelectorShape(last) {
		site.DefaultSelector = &last
		site.Selectors = site.Selectors[:n-1]
	}
}

func (c *Client) CreateAbpSite(accountId string, site AbpSite) (*AbpSite, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpSiteResourceName, accountId)

	body, err := json.Marshal(site)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpSiteResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpSiteAccountUrl(accountId), body, CreateAbpSite)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpSiteResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpSiteResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpSiteResourceName, accountId, string(responseBody))
	}

	var created AbpSite
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpSiteResourceName, err, string(responseBody))
	}

	splitDefaultSelector(&created)

	return &created, nil
}

func (c *Client) ReadAbpSite(siteId string) (*AbpSite, error) {
	log.Printf("[INFO] Reading %s with id %s", abpSiteResourceName, siteId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpSiteUrl(siteId), nil, ReadAbpSite)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpSiteResourceName, siteId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpSiteResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpSiteResourceName, siteId, string(responseBody))
	}

	var site AbpSite
	if err := json.Unmarshal(responseBody, &site); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpSiteResourceName, err, string(responseBody))
	}

	splitDefaultSelector(&site)

	return &site, nil
}

func (c *Client) UpdateAbpSite(siteId string, site AbpSite) (*AbpSite, error) {
	log.Printf("[INFO] Updating %s with id %s", abpSiteResourceName, siteId)

	payload := site
	if site.DefaultSelector != nil {
		payload.Selectors = make([]AbpSelector, 0, len(site.Selectors)+1)
		payload.Selectors = append(payload.Selectors, site.Selectors...)
		payload.Selectors = append(payload.Selectors, *site.DefaultSelector)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpSiteResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpSiteUrl(siteId), body, UpdateAbpSite)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpSiteResourceName, siteId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpSiteResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpSiteResourceName, siteId, string(responseBody))
	}

	var updated AbpSite
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpSiteResourceName, err, string(responseBody))
	}

	splitDefaultSelector(&updated)

	return &updated, nil
}

func (c *Client) DeleteAbpSite(siteId string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpSiteResourceName, siteId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpSiteUrl(siteId), nil, DeleteAbpSite)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpSiteResourceName, siteId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpSiteResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpSiteResourceName, siteId, string(responseBody))
	}
	return nil
}
