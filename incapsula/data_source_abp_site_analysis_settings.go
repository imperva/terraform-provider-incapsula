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
` + "`analysis_settings`" + ` field.`,

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
				Description:  "Override the site default for this selector.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_requests_per_session": {
				Description:  "Override the site default for this selector.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_session_length": {
				Description: "Override the site default for this selector. Moi duration format (e.g. \"2d1h\").",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"use_site_rate_limiting_parameters": {
				Description: "When true, fall back to the site-level defaults for rate-limiting parameters not set on this selector.",
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
