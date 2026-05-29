package incapsula

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAbpDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpDomainCreate,
		ReadContext:   resourceAbpDomainRead,
		UpdateContext: resourceAbpDomainUpdate,
		DeleteContext: resourceAbpDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpDomainImport,
		},

		Description: `Provides an ABP Domain (a.k.a. Website) resource. A Domain maps one or more
domain names to a Site, and configures how requests are matched, captured, and
analyzed.

NOTE: ` + "`criteria`" + ` cannot be changed in-place; modifying it forces resource
replacement. ` + "`obfuscate_path`" + ` is server-generated when omitted on create
and is preserved on update.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Domain belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"site_id": {
				Description:  "Site (a.k.a. Website Group) UUID this Domain is attached to.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"criteria": {
				Description: "Specifies which domain names this Domain matches. Exactly one of `exact`, `prefix`, `suffix`, `cloudwaf_id` must be set.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exact": {
							Description: "Match exactly this fully-qualified domain name.",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
						"prefix": {
							Description: "Match any domain name with this prefix.",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
						"suffix": {
							Description: "Match any domain name with this suffix.",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
						"cloudwaf_id": {
							Description: "CloudWAF website ID; use this when onboarding CloudWAF domains.",
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"cookiescope": {
				Description: "The Domain attribute of the Set-Cookie header set by the ABP JavaScript. Use `$suffix` as the TLD with prefix criteria.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"log_region": {
				Description:  "Region in which ABP logs are stored. One of: apac, australia, eu, usa, india.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "usa",
				ValidateFunc: validation.StringInSlice([]string{"apac", "australia", "eu", "usa", "india"}, false),
			},
			"cookie_mode": {
				Description:  "SameSite policy of the ABP cookies. One of: lax, legacy, none_secure, lax_and_none_secure, lax_and_legacy, legacy_and_none_secure.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "lax",
				ValidateFunc: validation.StringInSlice([]string{"lax", "legacy", "none_secure", "lax_and_none_secure", "lax_and_legacy", "legacy_and_none_secure"}, false),
			},
			"enable_mitigation": {
				Description: "If false, all active conditions are treated as monitor-only. If true (default), conditions behave per their state.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"enable_mobile_sdk_token": {
				Description: "Enable to allow mobile SDK tokens. Only enable if you use the ABP mobile SDK with this domain.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_out_static_assets": {
				Description: "CWAF only. Prevents common static asset paths from being analyzed by ABP.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"obfuscate_path": {
				Description: "Recommended path under which the ABP JavaScript is loaded. Server-generated when omitted on create.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"mobile_api_obfuscate_path": {
				Description: "Server-managed path used for the mobile API challenge endpoint.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"interstitial_inprogress_iframe_src": {
				Description: "URL of the iframe used to display the PoW challenge in progress.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"divert_host": {
				Description: "Domain or IP to which traffic is diverted when the `divert` directive triggers.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"unmasked_headers": {
				Description: "Header names whose values should be visible to ABP (CloudWAF normally masks them). Compared case-insensitively.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"proxy_flags": {
				Description: "CloudWAF configuration flags. Allowed values: enable_referrer_fix, dont_minify_post_resubmit_error_page, inject_js_into_body, set_shared_cookie_on_apex.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"enable_referrer_fix", "dont_minify_post_resubmit_error_page", "inject_js_into_body", "set_shared_cookie_on_apex"}, false),
				},
			},
			"encryption_key_id": {
				Description:  "Existing encryption key to copy when creating this Domain. If unset, a copy of the account default key is used.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"analysis_ip_lookup_mode":  ipLookupModeSchema("Controls how the analysis host determines the end-user's IP. Omit the block for `none`."),
			"challenge_ip_lookup_mode": ipLookupModeSchema("Controls how the challenge host determines the end-user's IP. Omit the block for `none`."),
			"captcha_settings": {
				Description: "CAPTCHA configuration. Omit the block for `none`. Set exactly one of `geetest`, `managed_geetest`, `managed_hcaptcha`.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geetest": {
							Description: "Use Geetest with your own credentials.",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"geetest_captcha_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"geetest_private_key": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
						"managed_geetest": {
							Description: "Use Imperva-managed Geetest. CloudWAF only.",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"difficulty": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"easy", "normal", "hard"}, false),
									},
								},
							},
						},
						"managed_hcaptcha": {
							Description: "Use Imperva-managed hCaptcha.",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"difficulty": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"auto"}, false),
									},
								},
							},
						},
					},
				},
			},
			"no_js_injection_path": {
				Description: "Rules describing where the ABP JavaScript injection should not occur. Each block must set exactly one of `path_prefix`, `incap_rule`.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_prefix": {
							Description: "Disable injection for paths starting with this prefix (must begin with `/`).",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"incap_rule": {
							Description: "Disable injection for paths matching this Incapsula rule.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Domain was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Domain was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func ipLookupModeSchema(description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"header_name": {
					Description: "Name of the HTTP header carrying the originating IP.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"reverse_index": {
					Description:  "Index from the right (0-based) when the header has multiple comma-separated IPs.",
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
	}
}

func extractAbpDomainCriteria(data *schema.ResourceData) (AbpDomainCriteria, error) {
	raw := data.Get("criteria").([]any)
	if len(raw) == 0 {
		return AbpDomainCriteria{}, fmt.Errorf("criteria block is required")
	}
	m := raw[0].(map[string]any)

	exact, _ := m["exact"].(string)
	prefix, _ := m["prefix"].(string)
	suffix, _ := m["suffix"].(string)
	cloudwafId, _ := m["cloudwaf_id"].(int)

	c := AbpDomainCriteria{}
	set := 0
	if exact != "" {
		c.Exact = exact
		set++
	}
	if prefix != "" {
		c.Prefix = prefix
		set++
	}
	if suffix != "" {
		c.Suffix = suffix
		set++
	}
	if cloudwafId != 0 {
		v := int64(cloudwafId)
		c.CloudwafId = &v
		set++
	}
	if set != 1 {
		return c, fmt.Errorf("criteria requires exactly one of exact, prefix, suffix, cloudwaf_id (got %d)", set)
	}
	return c, nil
}

func extractAbpIpLookupMode(raw []any) AbpIpLookupMode {
	if len(raw) == 0 {
		return AbpIpLookupMode{}
	}
	m := raw[0].(map[string]any)
	header, _ := m["header_name"].(string)
	idx, _ := m["reverse_index"].(int)
	return AbpIpLookupMode{
		HasLookup:    true,
		HeaderName:   header,
		ReverseIndex: idx,
	}
}

func extractAbpCaptchaSettings(raw []any) (AbpCaptchaSettings, error) {
	if len(raw) == 0 {
		return AbpCaptchaSettings{Mode: AbpCaptchaModeNone}, nil
	}
	m := raw[0].(map[string]any)
	geetest, _ := m["geetest"].([]any)
	managedGeetest, _ := m["managed_geetest"].([]any)
	managedHcaptcha, _ := m["managed_hcaptcha"].([]any)

	set := 0
	c := AbpCaptchaSettings{}
	if len(geetest) > 0 {
		set++
		g := geetest[0].(map[string]any)
		c.Mode = AbpCaptchaModeGeetest
		c.GeetestCaptchaId, _ = g["geetest_captcha_id"].(string)
		c.GeetestPrivateKey, _ = g["geetest_private_key"].(string)
	}
	if len(managedGeetest) > 0 {
		set++
		mg := managedGeetest[0].(map[string]any)
		c.Mode = AbpCaptchaModeManagedGeetest
		c.ManagedDifficulty, _ = mg["difficulty"].(string)
	}
	if len(managedHcaptcha) > 0 {
		set++
		mh := managedHcaptcha[0].(map[string]any)
		c.Mode = AbpCaptchaModeManagedHcaptcha
		c.ManagedDifficulty, _ = mh["difficulty"].(string)
	}
	if set != 1 {
		return c, fmt.Errorf("captcha_settings requires exactly one of geetest, managed_geetest, managed_hcaptcha (got %d)", set)
	}
	return c, nil
}

func extractAbpNoJsInjectionPaths(raw []any) ([]AbpNoJsInjectionPath, error) {
	out := make([]AbpNoJsInjectionPath, 0, len(raw))
	for i, item := range raw {
		m := item.(map[string]any)
		prefix, _ := m["path_prefix"].(string)
		rule, _ := m["incap_rule"].(string)
		set := 0
		p := AbpNoJsInjectionPath{}
		if prefix != "" {
			p.PathPrefix = prefix
			set++
		}
		if rule != "" {
			p.IncapRule = rule
			set++
		}
		if set != 1 {
			return nil, fmt.Errorf("no_js_injection_path[%d]: exactly one of path_prefix, incap_rule must be set (got %d)", i, set)
		}
		out = append(out, p)
	}
	return out, nil
}

func setToStringSlice(s any) []string {
	if s == nil {
		return []string{}
	}
	set, ok := s.(*schema.Set)
	if !ok {
		return []string{}
	}
	raw := set.List()
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if str, ok := v.(string); ok {
			out = append(out, str)
		}
	}
	return out
}

func extractAbpDomain(data *schema.ResourceData) (AbpDomain, error) {
	criteria, err := extractAbpDomainCriteria(data)
	if err != nil {
		return AbpDomain{}, err
	}
	captcha, err := extractAbpCaptchaSettings(data.Get("captcha_settings").([]any))
	if err != nil {
		return AbpDomain{}, err
	}
	noJs, err := extractAbpNoJsInjectionPaths(data.Get("no_js_injection_path").([]any))
	if err != nil {
		return AbpDomain{}, err
	}

	domain := AbpDomain{
		SiteId:                data.Get("site_id").(string),
		Cookiescope:           data.Get("cookiescope").(string),
		Criteria:              criteria,
		CaptchaSettings:       captcha,
		NoJsInjectionPaths:    noJs,
		LogRegion:             data.Get("log_region").(string),
		CookieMode:            data.Get("cookie_mode").(string),
		EnableMitigation:      data.Get("enable_mitigation").(bool),
		AnalysisIpLookupMode:  extractAbpIpLookupMode(data.Get("analysis_ip_lookup_mode").([]any)),
		ChallengeIpLookupMode: extractAbpIpLookupMode(data.Get("challenge_ip_lookup_mode").([]any)),
		UnmaskedHeaders:       setToStringSlice(data.Get("unmasked_headers")),
		ProxyFlags:            setToStringSlice(data.Get("proxy_flags")),
	}

	if v, ok := data.GetOk("obfuscate_path"); ok {
		s := v.(string)
		domain.ObfuscatePath = &s
	}
	if v, ok := data.GetOk("interstitial_inprogress_iframe_src"); ok {
		s := v.(string)
		domain.InterstitialInprogressIframeSrc = &s
	}
	if v, ok := data.GetOk("divert_host"); ok {
		s := v.(string)
		domain.DivertHost = &s
	}
	if v, ok := data.GetOk("encryption_key_id"); ok {
		s := v.(string)
		domain.EncryptionKeyId = &s
	}
	if v, ok := data.GetOkExists("filter_out_static_assets"); ok {
		b := v.(bool)
		domain.FilterOutStaticAssets = &b
	}
	if v, ok := data.GetOkExists("enable_mobile_sdk_token"); ok {
		b := v.(bool)
		domain.EnableMobileSdkToken = &b
	}

	return domain, nil
}

func flattenAbpDomainCriteria(c AbpDomainCriteria) []any {
	m := map[string]any{}
	if c.Exact != "" {
		m["exact"] = c.Exact
	}
	if c.Prefix != "" {
		m["prefix"] = c.Prefix
	}
	if c.Suffix != "" {
		m["suffix"] = c.Suffix
	}
	if c.CloudwafId != nil {
		m["cloudwaf_id"] = int(*c.CloudwafId)
	} else if c.CloudwafWebsiteId != nil {
		m["cloudwaf_id"] = int(*c.CloudwafWebsiteId)
	}
	return []any{m}
}

func flattenAbpIpLookupMode(m AbpIpLookupMode) []any {
	if !m.HasLookup {
		return []any{}
	}
	return []any{map[string]any{
		"header_name":   m.HeaderName,
		"reverse_index": m.ReverseIndex,
	}}
}

func flattenAbpCaptchaSettings(c AbpCaptchaSettings) []any {
	switch c.Mode {
	case "", AbpCaptchaModeNone:
		return []any{}
	case AbpCaptchaModeGeetest:
		return []any{map[string]any{
			"geetest": []any{map[string]any{
				"geetest_captcha_id":  c.GeetestCaptchaId,
				"geetest_private_key": c.GeetestPrivateKey,
			}},
		}}
	case AbpCaptchaModeManagedGeetest:
		return []any{map[string]any{
			"managed_geetest": []any{map[string]any{"difficulty": c.ManagedDifficulty}},
		}}
	case AbpCaptchaModeManagedHcaptcha:
		return []any{map[string]any{
			"managed_hcaptcha": []any{map[string]any{"difficulty": c.ManagedDifficulty}},
		}}
	default:
		return []any{}
	}
}

func flattenAbpNoJsInjectionPaths(paths []AbpNoJsInjectionPath) []any {
	out := make([]any, len(paths))
	for i, p := range paths {
		m := map[string]any{}
		if p.PathPrefix != "" {
			m["path_prefix"] = p.PathPrefix
		}
		if p.IncapRule != "" {
			m["incap_rule"] = p.IncapRule
		}
		out[i] = m
	}
	return out
}

func serializeAbpDomain(data *schema.ResourceData, domain *AbpDomain) error {
	if err := data.Set("account_id", domain.AccountId); err != nil {
		return err
	}
	if err := data.Set("site_id", domain.SiteId); err != nil {
		return err
	}
	if err := data.Set("criteria", flattenAbpDomainCriteria(domain.Criteria)); err != nil {
		return err
	}
	if err := data.Set("cookiescope", domain.Cookiescope); err != nil {
		return err
	}
	if err := data.Set("log_region", domain.LogRegion); err != nil {
		return err
	}
	if err := data.Set("cookie_mode", domain.CookieMode); err != nil {
		return err
	}
	if err := data.Set("enable_mitigation", domain.EnableMitigation); err != nil {
		return err
	}
	if err := data.Set("unmasked_headers", domain.UnmaskedHeaders); err != nil {
		return err
	}
	if err := data.Set("proxy_flags", domain.ProxyFlags); err != nil {
		return err
	}
	if err := data.Set("analysis_ip_lookup_mode", flattenAbpIpLookupMode(domain.AnalysisIpLookupMode)); err != nil {
		return err
	}
	if err := data.Set("challenge_ip_lookup_mode", flattenAbpIpLookupMode(domain.ChallengeIpLookupMode)); err != nil {
		return err
	}
	if err := data.Set("captcha_settings", flattenAbpCaptchaSettings(domain.CaptchaSettings)); err != nil {
		return err
	}
	if err := data.Set("no_js_injection_path", flattenAbpNoJsInjectionPaths(domain.NoJsInjectionPaths)); err != nil {
		return err
	}
	if domain.ObfuscatePath != nil {
		if err := data.Set("obfuscate_path", *domain.ObfuscatePath); err != nil {
			return err
		}
	}
	if domain.MobileApiObfuscatePath != nil {
		if err := data.Set("mobile_api_obfuscate_path", *domain.MobileApiObfuscatePath); err != nil {
			return err
		}
	}
	if domain.InterstitialInprogressIframeSrc != nil {
		if err := data.Set("interstitial_inprogress_iframe_src", *domain.InterstitialInprogressIframeSrc); err != nil {
			return err
		}
	}
	if domain.DivertHost != nil {
		if err := data.Set("divert_host", *domain.DivertHost); err != nil {
			return err
		}
	}
	if domain.FilterOutStaticAssets != nil {
		if err := data.Set("filter_out_static_assets", *domain.FilterOutStaticAssets); err != nil {
			return err
		}
	}
	if domain.EnableMobileSdkToken != nil {
		if err := data.Set("enable_mobile_sdk_token", *domain.EnableMobileSdkToken); err != nil {
			return err
		}
	}
	if err := data.Set("created_at", domain.CreatedAt); err != nil {
		return err
	}
	if err := data.Set("modified_at", domain.ModifiedAt); err != nil {
		return err
	}
	return nil
}

func resourceAbpDomainCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	domain, err := extractAbpDomain(data)
	if err != nil {
		return diag.FromErr(err)
	}

	created, err := client.CreateAbpDomain(accountId, domain)
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Domain create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpDomain(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Domain %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpDomainRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	domain, err := client.ReadAbpDomain(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if domain == nil {
		log.Printf("[INFO] ABP Domain %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpDomain(data, domain); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpDomainUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	domain, err := extractAbpDomain(data)
	if err != nil {
		return diag.FromErr(err)
	}

	updated, err := client.UpdateAbpDomain(id, domain)
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpDomain(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpDomainDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpDomain(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpDomainImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<domain_id>'")
	}

	client := m.(*Client)
	domain, err := client.ReadAbpDomain(id)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, fmt.Errorf("ABP Domain %s not found", id)
	}

	data.SetId(id)
	if err := data.Set("account_id", domain.AccountId); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
