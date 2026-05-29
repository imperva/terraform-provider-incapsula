package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpDomainResourceName = "ABP Domain"

type AbpDomain struct {
	Id                              string                 `json:"id,omitempty"`
	AccountId                       string                 `json:"account_id,omitempty"`
	SiteId                          string                 `json:"site_id"`
	Criteria                        AbpDomainCriteria      `json:"criteria,omitempty"`
	Cookiescope                     string                 `json:"cookiescope"`
	ChallengeIpLookupMode           AbpIpLookupMode        `json:"challenge_ip_lookup_mode"`
	AnalysisIpLookupMode            AbpIpLookupMode        `json:"analysis_ip_lookup_mode"`
	CaptchaSettings                 AbpCaptchaSettings     `json:"captcha_settings"`
	LogRegion                       string                 `json:"log_region"`
	NoJsInjectionPaths              []AbpNoJsInjectionPath `json:"no_js_injection_paths"`
	ObfuscatePath                   *string                `json:"obfuscate_path,omitempty"`
	MobileApiObfuscatePath          *string                `json:"mobile_api_obfuscate_path,omitempty"`
	CookieMode                      string                 `json:"cookie_mode"`
	UnmaskedHeaders                 []string               `json:"unmasked_headers"`
	ProxyFlags                      []string               `json:"proxy_flags"`
	FilterOutStaticAssets           *bool                  `json:"filter_out_static_assets,omitempty"`
	EnableMitigation                bool                   `json:"enable_mitigation"`
	EnableMobileSdkToken            *bool                  `json:"enable_mobile_sdk_token,omitempty"`
	InterstitialInprogressIframeSrc *string                `json:"interstitial_inprogress_iframe_src,omitempty"`
	DivertHost                      *string                `json:"divert_host,omitempty"`
	EncryptionKeyId                 *string                `json:"encryption_key_id,omitempty"`
	CreatedAt                       string                 `json:"created_at,omitempty"`
	ModifiedAt                      string                 `json:"modified_at,omitempty"`
}

// AbpDomainCriteria is a oneOf:
//
//	{exact: <fqdn>} | {prefix: ...} | {suffix: ...} | {cloudwaf_id: <int>} |
//	{cloudwaf_website_id: <int>, fqdn: <string>}  (read-only, returned by the server)
//
// On create, only `cloudwaf_id` (preferred) or `cloudwaf_website_id` may be sent
// for CloudWAF; on read, the server may include both plus `fqdn`. PUT/Update
// omits criteria entirely (it's not in UpdateDomainV1), so this struct is only
// used on Create and Read.
type AbpDomainCriteria struct {
	Exact             string
	Prefix            string
	Suffix            string
	CloudwafId        *int64
	CloudwafWebsiteId *int64
	Fqdn              string
}

func (c AbpDomainCriteria) MarshalJSON() ([]byte, error) {
	switch {
	case c.Exact != "":
		return json.Marshal(map[string]any{"exact": c.Exact})
	case c.Prefix != "":
		return json.Marshal(map[string]any{"prefix": c.Prefix})
	case c.Suffix != "":
		return json.Marshal(map[string]any{"suffix": c.Suffix})
	case c.CloudwafId != nil:
		return json.Marshal(map[string]any{"cloudwaf_id": *c.CloudwafId})
	case c.CloudwafWebsiteId != nil:
		return json.Marshal(map[string]any{"cloudwaf_website_id": *c.CloudwafWebsiteId})
	default:
		return nil, fmt.Errorf("domain criteria: exactly one of exact, prefix, suffix, cloudwaf_id must be set")
	}
}

func (c *AbpDomainCriteria) UnmarshalJSON(data []byte) error {
	var obj struct {
		Exact             string `json:"exact"`
		Prefix            string `json:"prefix"`
		Suffix            string `json:"suffix"`
		CloudwafId        *int64 `json:"cloudwaf_id"`
		CloudwafWebsiteId *int64 `json:"cloudwaf_website_id"`
		Fqdn              string `json:"fqdn"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("domain criteria: %w", err)
	}
	c.Exact = obj.Exact
	c.Prefix = obj.Prefix
	c.Suffix = obj.Suffix
	c.CloudwafId = obj.CloudwafId
	c.CloudwafWebsiteId = obj.CloudwafWebsiteId
	c.Fqdn = obj.Fqdn
	return nil
}

// AbpIpLookupMode is a oneOf:
//
//	"none" | {lookup: {header_name, reverse_index}}
type AbpIpLookupMode struct {
	HeaderName   string
	ReverseIndex int
	HasLookup    bool
}

func (m AbpIpLookupMode) MarshalJSON() ([]byte, error) {
	if !m.HasLookup {
		return json.Marshal("none")
	}
	return json.Marshal(map[string]any{
		"lookup": map[string]any{
			"header_name":   m.HeaderName,
			"reverse_index": m.ReverseIndex,
		},
	})
}

func (m *AbpIpLookupMode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.HasLookup = false
		m.HeaderName = ""
		m.ReverseIndex = 0
		return nil
	}
	var obj struct {
		Lookup *struct {
			HeaderName   string `json:"header_name"`
			ReverseIndex int    `json:"reverse_index"`
		} `json:"lookup"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("ip_lookup_mode: expected \"none\" or {lookup}: %w", err)
	}
	if obj.Lookup == nil {
		return fmt.Errorf("ip_lookup_mode: lookup object missing")
	}
	m.HasLookup = true
	m.HeaderName = obj.Lookup.HeaderName
	m.ReverseIndex = obj.Lookup.ReverseIndex
	return nil
}

// AbpCaptchaSettings is a oneOf:
//
//	"none" | {geetest:{...}} | {managed_geetest:{difficulty}} | {managed_hcaptcha:{difficulty}}
//
// Mode selects which alternative is used. Other fields apply only to the
// matching mode.
type AbpCaptchaSettings struct {
	Mode              string
	GeetestCaptchaId  string
	GeetestPrivateKey string
	ManagedDifficulty string
}

const (
	AbpCaptchaModeNone            = "none"
	AbpCaptchaModeGeetest         = "geetest"
	AbpCaptchaModeManagedGeetest  = "managed_geetest"
	AbpCaptchaModeManagedHcaptcha = "managed_hcaptcha"
)

func (c AbpCaptchaSettings) MarshalJSON() ([]byte, error) {
	switch c.Mode {
	case "", AbpCaptchaModeNone:
		return json.Marshal("none")
	case AbpCaptchaModeGeetest:
		return json.Marshal(map[string]any{
			"geetest": map[string]any{
				"geetest_captcha_id":  c.GeetestCaptchaId,
				"geetest_private_key": c.GeetestPrivateKey,
			},
		})
	case AbpCaptchaModeManagedGeetest:
		return json.Marshal(map[string]any{
			"managed_geetest": map[string]any{"difficulty": c.ManagedDifficulty},
		})
	case AbpCaptchaModeManagedHcaptcha:
		return json.Marshal(map[string]any{
			"managed_hcaptcha": map[string]any{"difficulty": c.ManagedDifficulty},
		})
	default:
		return nil, fmt.Errorf("captcha_settings: unknown mode %q", c.Mode)
	}
}

func (c *AbpCaptchaSettings) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Mode = AbpCaptchaModeNone
		return nil
	}
	var obj struct {
		Geetest *struct {
			GeetestCaptchaId  string `json:"geetest_captcha_id"`
			GeetestPrivateKey string `json:"geetest_private_key"`
		} `json:"geetest"`
		ManagedGeetest *struct {
			Difficulty string `json:"difficulty"`
		} `json:"managed_geetest"`
		ManagedHcaptcha *struct {
			Difficulty string `json:"difficulty"`
		} `json:"managed_hcaptcha"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("captcha_settings: %w", err)
	}
	switch {
	case obj.Geetest != nil:
		c.Mode = AbpCaptchaModeGeetest
		c.GeetestCaptchaId = obj.Geetest.GeetestCaptchaId
		c.GeetestPrivateKey = obj.Geetest.GeetestPrivateKey
	case obj.ManagedGeetest != nil:
		c.Mode = AbpCaptchaModeManagedGeetest
		c.ManagedDifficulty = obj.ManagedGeetest.Difficulty
	case obj.ManagedHcaptcha != nil:
		c.Mode = AbpCaptchaModeManagedHcaptcha
		c.ManagedDifficulty = obj.ManagedHcaptcha.Difficulty
	default:
		return fmt.Errorf("captcha_settings: object did not match any known variant")
	}
	return nil
}

// AbpNoJsInjectionPath is a oneOf:
//
//	{path_prefix: ...} | {incap_rule: ...}
type AbpNoJsInjectionPath struct {
	PathPrefix string
	IncapRule  string
}

func (p AbpNoJsInjectionPath) MarshalJSON() ([]byte, error) {
	switch {
	case p.PathPrefix != "" && p.IncapRule == "":
		return json.Marshal(map[string]any{"path_prefix": p.PathPrefix})
	case p.IncapRule != "" && p.PathPrefix == "":
		return json.Marshal(map[string]any{"incap_rule": p.IncapRule})
	default:
		return nil, fmt.Errorf("no_js_injection_path: exactly one of path_prefix, incap_rule must be set")
	}
}

func (p *AbpNoJsInjectionPath) UnmarshalJSON(data []byte) error {
	var obj struct {
		PathPrefix string `json:"path_prefix"`
		IncapRule  string `json:"incap_rule"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("no_js_injection_path: %w", err)
	}
	p.PathPrefix = obj.PathPrefix
	p.IncapRule = obj.IncapRule
	return nil
}

func (c *Client) abpDomainAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/domain", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpDomainUrl(domainId string) string {
	return fmt.Sprintf("%s/v1/domain/%s", c.config.BaseURLAPI, domainId)
}

func (c *Client) CreateAbpDomain(accountId string, domain AbpDomain) (*AbpDomain, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpDomainResourceName, accountId)

	body, err := json.Marshal(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpDomainResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpDomainAccountUrl(accountId), body, CreateAbpDomain)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpDomainResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpDomainResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpDomainResourceName, accountId, string(responseBody))
	}

	var created AbpDomain
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpDomainResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpDomain(domainId string) (*AbpDomain, error) {
	log.Printf("[INFO] Reading %s with id %s", abpDomainResourceName, domainId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpDomainUrl(domainId), nil, ReadAbpDomain)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpDomainResourceName, domainId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpDomainResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpDomainResourceName, domainId, string(responseBody))
	}

	var domain AbpDomain
	if err := json.Unmarshal(responseBody, &domain); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpDomainResourceName, err, string(responseBody))
	}
	return &domain, nil
}

func (c *Client) UpdateAbpDomain(domainId string, domain AbpDomain) (*AbpDomain, error) {
	log.Printf("[INFO] Updating %s with id %s", abpDomainResourceName, domainId)

	// criteria is not part of UpdateDomainV1; clear it so it's omitted.
	updateBody := domain
	updateBody.Criteria = AbpDomainCriteria{}

	type updatePayload struct {
		SiteId                          string                 `json:"site_id"`
		ChallengeIpLookupMode           AbpIpLookupMode        `json:"challenge_ip_lookup_mode"`
		AnalysisIpLookupMode            AbpIpLookupMode        `json:"analysis_ip_lookup_mode"`
		Cookiescope                     string                 `json:"cookiescope"`
		CaptchaSettings                 AbpCaptchaSettings     `json:"captcha_settings"`
		LogRegion                       string                 `json:"log_region"`
		NoJsInjectionPaths              []AbpNoJsInjectionPath `json:"no_js_injection_paths"`
		ObfuscatePath                   *string                `json:"obfuscate_path,omitempty"`
		CookieMode                      string                 `json:"cookie_mode"`
		UnmaskedHeaders                 []string               `json:"unmasked_headers"`
		ProxyFlags                      []string               `json:"proxy_flags"`
		FilterOutStaticAssets           *bool                  `json:"filter_out_static_assets,omitempty"`
		EnableMitigation                bool                   `json:"enable_mitigation"`
		EnableMobileSdkToken            *bool                  `json:"enable_mobile_sdk_token,omitempty"`
		InterstitialInprogressIframeSrc *string                `json:"interstitial_inprogress_iframe_src,omitempty"`
		DivertHost                      *string                `json:"divert_host,omitempty"`
	}
	body, err := json.Marshal(updatePayload{
		SiteId:                          domain.SiteId,
		ChallengeIpLookupMode:           domain.ChallengeIpLookupMode,
		AnalysisIpLookupMode:            domain.AnalysisIpLookupMode,
		Cookiescope:                     domain.Cookiescope,
		CaptchaSettings:                 domain.CaptchaSettings,
		LogRegion:                       domain.LogRegion,
		NoJsInjectionPaths:              domain.NoJsInjectionPaths,
		ObfuscatePath:                   domain.ObfuscatePath,
		CookieMode:                      domain.CookieMode,
		UnmaskedHeaders:                 domain.UnmaskedHeaders,
		ProxyFlags:                      domain.ProxyFlags,
		FilterOutStaticAssets:           domain.FilterOutStaticAssets,
		EnableMitigation:                domain.EnableMitigation,
		EnableMobileSdkToken:            domain.EnableMobileSdkToken,
		InterstitialInprogressIframeSrc: domain.InterstitialInprogressIframeSrc,
		DivertHost:                      domain.DivertHost,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpDomainResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpDomainUrl(domainId), body, UpdateAbpDomain)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpDomainResourceName, domainId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpDomainResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpDomainResourceName, domainId, string(responseBody))
	}

	var updated AbpDomain
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpDomainResourceName, err, string(responseBody))
	}
	return &updated, nil
}

func (c *Client) DeleteAbpDomain(domainId string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpDomainResourceName, domainId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpDomainUrl(domainId), nil, DeleteAbpDomain)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpDomainResourceName, domainId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpDomainResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpDomainResourceName, domainId, string(responseBody))
	}
	return nil
}
