package incapsula

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpSiteAnalysisSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpSiteAnalysisSettingsRead,

		Description: `Builds an analysis-settings document for use in an ` + "`incapsula_abp_site`" + `
selector. Pass the resulting ` + "`json`" + ` attribute to the selector's
` + "`analysis_settings`" + ` field.

NOTE: ` + "`use_site_rate_limiting_parameters`" + ` defaults to ` + "`true`" + `,
in which case the Site-level defaults are used and
` + "`max_requests_per_minute`" + `, ` + "`max_requests_per_session`" + ` and
` + "`max_session_length`" + ` have no effect at all. Set it to
` + "`false`" + ` whenever you set any of the three.`,

		Schema: map[string]*schema.Schema{
			"rate_limiting": {
				Description:  "Rate limiting scope. One of: none, per_site, custom_scope. When `custom_scope`, set `rate_limiting_custom_scope`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{AbpRateLimitingModeNone, AbpRateLimitingModePerSite, AbpRateLimitingModeCustomScope}, false),
			},
			"rate_limiting_custom_scope": {
				Description: "Custom scope name; required iff `rate_limiting = custom_scope`. Selectors sharing the same scope share rate-limit counters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"max_requests_per_minute": {
				Description:  "Maximum number of requests without a token per minute for selectors using these settings. Only takes effect when `use_site_rate_limiting_parameters` is set to `false`.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_requests_per_session": {
				Description:  "Maximum number of requests without a token per session for selectors using these settings. Only takes effect when `use_site_rate_limiting_parameters` is set to `false`.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_session_length": {
				Description: "Maximum length of a session without a token for selectors using these settings. Only takes effect when `use_site_rate_limiting_parameters` is set to `false`. Duration format, e.g. \"2d1h\". Use the largest fitting units and omit zero components: write \"2h\", not \"120m\".",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"use_site_rate_limiting_parameters": {
				Description: "Whether the Site-level defaults are used instead of the three values above. Defaults to `true`, which means the Site's `default_max_requests_per_minute`, `default_max_requests_per_session` and `default_max_session_length` apply and `max_requests_per_minute`, `max_requests_per_session` and `max_session_length` are **ignored entirely** — it is not a per-value fallback. Set it to `false` whenever you want those three to take effect.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"json": {
				Description: "Canonical JSON encoding of the analysis settings. Pass this to a selector's `analysis_settings` field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAbpSiteAnalysisSettingsRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	mode := data.Get("rate_limiting").(string)
	scope := data.Get("rate_limiting_custom_scope").(string)

	rl := AbpRateLimiting{Mode: mode}
	switch mode {
	case AbpRateLimitingModeCustomScope:
		if scope == "" {
			return diag.Errorf("rate_limiting_custom_scope must be set when rate_limiting = %q", AbpRateLimitingModeCustomScope)
		}
		rl.CustomScope = scope
	case AbpRateLimitingModeNone, AbpRateLimitingModePerSite:
		if scope != "" {
			return diag.Errorf("rate_limiting_custom_scope must not be set when rate_limiting = %q", mode)
		}
	}

	settings := AbpAnalysisSettings{RateLimiting: rl}
	if v, ok := data.GetOk("max_requests_per_minute"); ok {
		n := v.(int)
		settings.MaxRequestsPerMinute = &n
	}
	if v, ok := data.GetOk("max_requests_per_session"); ok {
		n := v.(int)
		settings.MaxRequestsPerSession = &n
	}
	if v, ok := data.GetOk("max_session_length"); ok {
		s := v.(string)
		settings.MaxSessionLength = &s
	}
	useDefaults := data.Get("use_site_rate_limiting_parameters").(bool)
	settings.UseSiteRateLimitingParameters = &useDefaults

	encoded, err := json.Marshal(settings)
	if err != nil {
		return diag.FromErr(fmt.Errorf("encoding analysis_settings: %w", err))
	}
	if err := data.Set("json", string(encoded)); err != nil {
		return diag.FromErr(err)
	}

	sum := sha256.Sum256(encoded)
	data.SetId(hex.EncodeToString(sum[:]))
	return nil
}
