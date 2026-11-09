package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"slices"
)

const abpSiteResourceName = "ABP Site"

type AbpSite struct {
	Id                           string        `json:"id,omitempty"`
	AccountId                    string        `json:"account_id,omitempty"`
	Name                         string        `json:"name"`
	Industry                     *string       `json:"industry,omitempty"`
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
	Purpose          *string             `json:"purpose,omitempty"`
	AnalysisSettings AbpAnalysisSettings `json:"analysis_settings"`
	DerivedId        string              `json:"derived_id,omitempty"`
	CreatedAt        string              `json:"created_at,omitempty"`
	ModifiedAt       string              `json:"modified_at,omitempty"`
}

// AbpSelectorCriteria is a oneOf — exactly one of the match fields must be
// set. `case_insensitive` applies to the `path_prefix` variant only; the
// backend omits it when false.
type AbpSelectorCriteria struct {
	Postback         *string `json:"postback,omitempty"`
	PathPrefix       *string `json:"path_prefix,omitempty"`
	PathRegex        *string `json:"path_regex,omitempty"`
	ApiSecurityGroup *string `json:"api_security_group,omitempty"`
	CaseInsensitive  bool    `json:"case_insensitive,omitempty"`
}

func (AbpSelectorCriteria) abpWireOneOf() abpWireOneOf {
	return abpWireOneOf{
		Selectors: []string{"postback", "path_prefix", "path_regex", "api_security_group"},
		Extra:     []string{"case_insensitive"},
	}
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

func (AbpRateLimiting) abpWireOneOf() abpWireOneOf {
	return abpWireOneOf{Selectors: []string{"custom_scope"}}
}

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
	return fmt.Sprintf("%s/botmanagement/v1/account/%s/site", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpSiteUrl(siteId string) string {
	return fmt.Sprintf("%s/botmanagement/v1/site/%s", c.config.BaseURLAPI, siteId)
}

// abpSiteSelectorIdentity is what the provider already knows about a Site's
// Selectors, taken from what it last recorded about them.
type abpSiteSelectorIdentity struct {
	// DefaultId is the catch-all Selector's id, empty when no catch-all is
	// known.
	DefaultId string
	// OwnedIds are the ids of the Selectors the configuration manages.
	OwnedIds []string
}

func isDefaultSelectorShape(s AbpSelector) bool {
	return s.Criteria.PathPrefix != nil && *s.Criteria.PathPrefix == "/"
}

// splitDefaultSelector moves the Site's catch-all Selector out of `Selectors`
// into `DefaultSelector`.
//
// The API does not mark it, so it has to be recognised, and its `/` path
// prefix is no help: a configuration may legitimately write a Selector that
// matches every path. It is therefore recognised by its id, which `known`
// carries. A catch-all the Site no longer holds is simply gone, and a trailing
// Selector the configuration owns is left where it belongs.
//
// Only when nothing is known about the Site — an import, or a data source read
// — is the catch-all guessed from the shape of the trailing Selector, which is
// the rule the API itself applies when deciding which Selector serves the
// paths no other one claims. The guess cannot tell a configuration's own
// catch-all from the API's, which is why it is used nowhere an answer is
// available.
func splitDefaultSelector(site *AbpSite, known abpSiteSelectorIdentity) {
	n := len(site.Selectors)
	if n == 0 {
		return
	}

	if known.DefaultId != "" {
		for i, selector := range site.Selectors {
			if selector.Id != known.DefaultId {
				continue
			}
			found := site.Selectors[i]
			site.DefaultSelector = &found
			site.Selectors = slices.Delete(slices.Clone(site.Selectors), i, i+1)
			return
		}
		return
	}

	last := site.Selectors[n-1]
	if slices.Contains(known.OwnedIds, last.Id) {
		return
	}
	if isDefaultSelectorShape(last) {
		site.DefaultSelector = &last
		site.Selectors = site.Selectors[:n-1]
	}
}

// splitDefaultSelectorAfter moves whatever a write came back with beyond the
// `sent` Selectors it submitted into `DefaultSelector`.
//
// A create has the API's catch-all appended to it — always, even when the
// Selectors sent already end with one of their own — so the extra trailing
// Selector of the response is it. An update carries the catch-all itself, so
// there is nothing beyond what it sent; an update that carried none leaves the
// Site without one, which is what a configuration whose own trailing Selector
// matches every path gets: that Selector serves the paths no other one claims,
// and the configuration owns it.
func splitDefaultSelectorAfter(site *AbpSite, sent int) {
	n := len(site.Selectors)
	if n <= sent {
		return
	}
	last := site.Selectors[n-1]
	site.DefaultSelector = &last
	site.Selectors = site.Selectors[:n-1]
}

func flattenAbpSelectors(selectors []AbpSelector) ([]any, error) {
	out := make([]any, len(selectors))
	for i, s := range selectors {
		encoded, err := json.Marshal(s.AnalysisSettings)
		if err != nil {
			return nil, fmt.Errorf("encoding selector[%d].analysis_settings: %w", i, err)
		}
		m := map[string]any{
			"id":                s.Id,
			"analysis_settings": string(encoded),
			"derived_id":        s.DerivedId,
			"created_at":        s.CreatedAt,
			"modified_at":       s.ModifiedAt,
			"purpose":           "",
		}
		if s.PolicyId != nil {
			m["policy_id"] = *s.PolicyId
		}
		if s.Purpose != nil {
			m["purpose"] = *s.Purpose
		}

		criteria, err := flattenAbpSelectorCriteria(s.Criteria)
		if err != nil {
			return nil, fmt.Errorf("selector[%d]: %w", i, err)
		}
		for name, block := range criteria {
			m[name] = block
		}
		out[i] = m
	}
	return out, nil
}

// flattenAbpSelectorCriteria turns a criteria into the single block holding
// it. An unrecognised criteria is an error rather than no block at all, so
// that a criteria this provider version does not know about is never silently
// rewritten into a different one.
func flattenAbpSelectorCriteria(c AbpSelectorCriteria) (map[string]any, error) {
	switch {
	case c.PathPrefix != nil:
		return map[string]any{"path_prefix": []any{map[string]any{
			"value":            *c.PathPrefix,
			"case_insensitive": c.CaseInsensitive,
		}}}, nil
	case c.PathRegex != nil:
		return map[string]any{"path_regex": []any{map[string]any{"value": *c.PathRegex}}}, nil
	case c.Postback != nil:
		return map[string]any{"postback": []any{map[string]any{"value": *c.Postback}}}, nil
	case c.ApiSecurityGroup != nil:
		return map[string]any{"api_security_group": []any{map[string]any{"value": *c.ApiSecurityGroup}}}, nil
	default:
		return nil, fmt.Errorf("selector uses a match criteria this provider version does not support; upgrade the provider to manage this Site")
	}
}

func flattenAbpSite(site *AbpSite) (map[string]any, error) {
	selectors, err := flattenAbpSelectors(site.Selectors)
	if err != nil {
		return nil, err
	}

	var defaultSelector []any
	if site.DefaultSelector != nil {
		defaultSelector, err = flattenAbpSelectors([]AbpSelector{*site.DefaultSelector})
		if err != nil {
			return nil, err
		}
	}

	m := map[string]any{
		"id":               site.Id,
		"account_id":       site.AccountId,
		"name":             site.Name,
		"industry":         "",
		"selector":         selectors,
		"default_selector": defaultSelector,
		"created_at":       site.CreatedAt,
		"modified_at":      site.ModifiedAt,
	}
	if site.Industry != nil {
		m["industry"] = *site.Industry
	}
	if site.DefaultMaxRequestsPerMinute != nil {
		m["default_max_requests_per_minute"] = *site.DefaultMaxRequestsPerMinute
	}
	if site.DefaultMaxRequestsPerSession != nil {
		m["default_max_requests_per_session"] = *site.DefaultMaxRequestsPerSession
	}
	if site.DefaultMaxSessionLength != nil {
		m["default_max_session_length"] = *site.DefaultMaxSessionLength
	}
	return m, nil
}

func (c *Client) ListAbpSites(accountId string) ([]AbpSite, error) {
	log.Printf("[INFO] Listing %ss in ABP account %s", abpSiteResourceName, accountId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpSiteAccountUrl(accountId), nil, ListAbpSites)
	if err != nil {
		return nil, fmt.Errorf("error listing %ss in ABP account %s: %w", abpSiteResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when listing %ss: %w", abpSiteResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when listing %ss in ABP account %s: %s", resp.StatusCode, abpSiteResourceName, accountId, string(responseBody))
	}

	var listResp struct {
		Items []AbpSite `json:"items"`
	}
	if err := json.Unmarshal(responseBody, &listResp); err != nil {
		return nil, fmt.Errorf("error parsing %s list response: %w; body: %s", abpSiteResourceName, err, string(responseBody))
	}

	for i := range listResp.Items {
		// A listing knows nothing about any Site it returns, so every
		// catch-all is the guess the API's own rule makes.
		splitDefaultSelector(&listResp.Items[i], abpSiteSelectorIdentity{})
	}
	return listResp.Items, nil
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

	splitDefaultSelectorAfter(&created, len(site.Selectors))

	return &created, nil
}

func (c *Client) ReadAbpSite(siteId string, known abpSiteSelectorIdentity) (*AbpSite, error) {
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

	splitDefaultSelector(&site, known)

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

	body, err = c.abpPreserveUnmodelledFields(c.abpSiteUrl(siteId), body, reflect.TypeOf(AbpSite{}), abpSiteResourceName, ReadAbpSite)
	if err != nil {
		return nil, err
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

	splitDefaultSelectorAfter(&updated, len(site.Selectors))

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
