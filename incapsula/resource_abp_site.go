package incapsula

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAbpSite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpSiteCreate,
		ReadContext:   resourceAbpSiteRead,
		UpdateContext: resourceAbpSiteUpdate,
		DeleteContext: resourceAbpSiteDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpSiteImport,
		},

		Description: `Provides an ABP Site (a.k.a. Website Group) resource. A Site groups one or
more Domains together and maps incoming requests to Policies via an ordered
list of Selectors.

NOTE: The API automatically appends a default catch-all selector
(` + "`path_prefix = \"/\"`" + `) at the lowest priority on create. You manage only
the user-defined selectors in the ` + "`selector`" + ` list; the default is exposed
read-only as ` + "`default_selector`" + ` and is preserved across updates.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Site belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Human-readable name of the Site. 1..100 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"default_max_requests_per_minute": {
				Description:  "Default maximum number of requests without a token per minute. Applied to selectors that opt in via `use_site_rate_limiting_parameters`.",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"default_max_requests_per_session": {
				Description:  "Default maximum number of requests without a token per session.",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"default_max_session_length": {
				Description: "Default maximum length of a session without a token, in moi duration format (e.g. \"2d1h\").",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"selector": {
				Description: "Ordered list of Selectors evaluated top-down for this Site. The first matching Selector decides which Policy is applied.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Server-assigned Selector ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"policy_id": {
							Description:  "Policy applied when this Selector matches. Omit to apply no policy (e.g. for static assets).",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"kind": {
							Description: "Match criteria for this Selector. Exactly one of `path_prefix`, `path_regex`, `postback` must be set.",
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_prefix": {
										Description: "Match requests whose path begins with this prefix. Mutually exclusive with `path_regex` and `postback`.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"path_regex": {
										Description: "Match requests whose path matches this regular expression. Mutually exclusive with `path_prefix` and `postback`.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"postback": {
										Description:  "Match a specific Postback request type. One of: web_interrogation, ios_interrogation, web_automation, android_interrogation. Mutually exclusive with `path_prefix` and `path_regex`.",
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"web_interrogation", "ios_interrogation", "web_automation", "android_interrogation"}, false),
									},
								},
							},
						},
						"analysis_settings": {
							Description:  "JSON-encoded analysis settings for this selector, typically produced by an `incapsula_abp_site_analysis_settings` data source.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
						},
					},
				},
			},
			"default_selector": {
				Description: "Catch-all selector matching all request paths after all other selectors. Added automatically by the backend",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Server-assigned Selector ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"policy_id": {
							Description: "Default Policy applied when no user-defined selector matches.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path_prefix": {
							Description: "Always `/`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"analysis_settings": {
							Description: "JSON-encoded analysis settings of the default selector.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Site was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Site was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// extractAbpSelector builds a selector from its config map. User-managed
// selectors carry their criteria in a nested `kind` block (nestedKind=true);
// the server-managed default selector keeps them flat (nestedKind=false).
func extractAbpSelector(raw map[string]any, nestedKind bool) (AbpSelector, error) {
	criteriaRaw := raw
	if nestedKind {
		kind, ok := raw["kind"].([]any)
		if !ok || len(kind) != 1 {
			return AbpSelector{}, fmt.Errorf("selector requires exactly one kind block")
		}
		criteriaRaw = kind[0].(map[string]any)
	}
	criteria, err := extractAbpSelectorCriteria(criteriaRaw)
	if err != nil {
		return AbpSelector{}, err
	}

	encoded, _ := raw["analysis_settings"].(string)
	var settings AbpAnalysisSettings
	if err := json.Unmarshal([]byte(encoded), &settings); err != nil {
		return AbpSelector{}, fmt.Errorf("decoding analysis_settings: %w", err)
	}

	sel := AbpSelector{
		Criteria:         criteria,
		AnalysisSettings: settings,
	}
	// Carry the server-assigned id when present (Update path) so backend
	// preserves selector identity across the PUT. Empty on Create.
	if id, ok := raw["id"].(string); ok && id != "" {
		sel.Id = id
	}
	if pid, ok := raw["policy_id"].(string); ok && pid != "" {
		sel.PolicyId = &pid
	}
	return sel, nil
}

func extractAbpSelectorCriteria(raw map[string]any) (AbpSelectorCriteria, error) {
	prefix, _ := raw["path_prefix"].(string)
	regex, _ := raw["path_regex"].(string)
	postback, _ := raw["postback"].(string)

	set := 0
	c := AbpSelectorCriteria{}
	if prefix != "" {
		c.PathPrefix = &prefix
		set++
	}
	if regex != "" {
		c.PathRegex = &regex
		set++
	}
	if postback != "" {
		c.Postback = &postback
		set++
	}
	if set != 1 {
		return c, fmt.Errorf("selector requires exactly one of path_prefix, path_regex, postback (got %d)", set)
	}
	return c, nil
}

func extractAbpSite(data *schema.ResourceData) (AbpSite, error) {
	site := AbpSite{
		Name:      data.Get("name").(string),
		Selectors: []AbpSelector{},
	}
	if v, ok := data.GetOk("default_max_requests_per_minute"); ok {
		n := v.(int)
		site.DefaultMaxRequestsPerMinute = &n
	}
	if v, ok := data.GetOk("default_max_requests_per_session"); ok {
		n := v.(int)
		site.DefaultMaxRequestsPerSession = &n
	}
	if v, ok := data.GetOk("default_max_session_length"); ok {
		s := v.(string)
		site.DefaultMaxSessionLength = &s
	}

	rawSelectors := data.Get("selector").([]any)
	site.Selectors = make([]AbpSelector, 0, len(rawSelectors))
	for i, item := range rawSelectors {
		sel, err := extractAbpSelector(item.(map[string]any), true)
		if err != nil {
			return AbpSite{}, fmt.Errorf("selector[%d]: %w", i, err)
		}
		site.Selectors = append(site.Selectors, sel)
	}

	// Carry the auto-managed default forward into the Update payload so
	// backend preserves its id and policy. On Create this list is empty.
	if rawDefault, ok := data.Get("default_selector").([]any); ok && len(rawDefault) == 1 {
		def, err := extractAbpSelector(rawDefault[0].(map[string]any), false)
		if err != nil {
			return AbpSite{}, fmt.Errorf("default_selector: %w", err)
		}
		site.DefaultSelector = &def
	}
	return site, nil
}

func flattenAbpSelectors(selectors []AbpSelector, nestedKind bool) ([]any, error) {
	out := make([]any, len(selectors))
	for i, s := range selectors {
		encoded, err := json.Marshal(s.AnalysisSettings)
		if err != nil {
			return nil, fmt.Errorf("encoding selector[%d].analysis_settings: %w", i, err)
		}
		m := map[string]any{
			"id":                s.Id,
			"analysis_settings": string(encoded),
		}
		if s.PolicyId != nil {
			m["policy_id"] = *s.PolicyId
		}

		criteria := map[string]any{}
		if s.Criteria.PathPrefix != nil {
			criteria["path_prefix"] = *s.Criteria.PathPrefix
		}
		if s.Criteria.PathRegex != nil {
			criteria["path_regex"] = *s.Criteria.PathRegex
		}
		if s.Criteria.Postback != nil {
			criteria["postback"] = *s.Criteria.Postback
		}
		if nestedKind {
			m["kind"] = []any{criteria}
		} else {
			for k, v := range criteria {
				m[k] = v
			}
		}
		out[i] = m
	}
	return out, nil
}

func serializeAbpSite(data *schema.ResourceData, site *AbpSite) error {
	if err := data.Set("account_id", site.AccountId); err != nil {
		return err
	}
	if err := data.Set("name", site.Name); err != nil {
		return err
	}
	if site.DefaultMaxRequestsPerMinute != nil {
		if err := data.Set("default_max_requests_per_minute", *site.DefaultMaxRequestsPerMinute); err != nil {
			return err
		}
	}
	if site.DefaultMaxRequestsPerSession != nil {
		if err := data.Set("default_max_requests_per_session", *site.DefaultMaxRequestsPerSession); err != nil {
			return err
		}
	}
	if site.DefaultMaxSessionLength != nil {
		if err := data.Set("default_max_session_length", *site.DefaultMaxSessionLength); err != nil {
			return err
		}
	}
	flat, err := flattenAbpSelectors(site.Selectors, true)
	if err != nil {
		return err
	}
	if err := data.Set("selector", flat); err != nil {
		return err
	}

	var flatDefault []any
	if site.DefaultSelector != nil {
		flatDefault, err = flattenAbpSelectors([]AbpSelector{*site.DefaultSelector}, false)
		if err != nil {
			return err
		}
	}
	if err := data.Set("default_selector", flatDefault); err != nil {
		return err
	}

	if err := data.Set("created_at", site.CreatedAt); err != nil {
		return err
	}
	if err := data.Set("modified_at", site.ModifiedAt); err != nil {
		return err
	}
	return nil
}

func resourceAbpSiteCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	site, err := extractAbpSite(data)
	if err != nil {
		return diag.FromErr(err)
	}

	created, err := client.CreateAbpSite(accountId, site)
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Site create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpSite(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Site %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpSiteRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	site, err := client.ReadAbpSite(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if site == nil {
		log.Printf("[INFO] ABP Site %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpSite(data, site); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpSiteUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	site, err := extractAbpSite(data)
	if err != nil {
		return diag.FromErr(err)
	}

	updated, err := client.UpdateAbpSite(id, site)
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpSite(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpSiteDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpSite(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpSiteImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<site_id>'")
	}

	client := m.(*Client)
	site, err := client.ReadAbpSite(id)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, fmt.Errorf("ABP Site %s not found", id)
	}

	data.SetId(id)
	if err := data.Set("account_id", site.AccountId); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
