package incapsula

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// abpSiteIndustries is SiteIndustryV1.
var abpSiteIndustries = []string{
	"financial_services", "retail", "travel", "media", "telecom", "crypto",
	"healthcare", "government", "ticketing", "real_estate", "insurance",
	"automotive", "nonprofit", "utilities", "other",
}

// abpSelectorPurposes is SelectorPurposeV1.
var abpSelectorPurposes = []string{
	"login", "signup", "password_reset", "mfa", "search", "browse", "cart",
	"checkout", "payment", "form_submission", "booking", "backend_api",
	"unknown", "other",
}

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
				Description: "ABP account UUID this Site belongs to. Changing it replaces the " +
					"Site — Terraform deletes the existing Site, with its selectors and Domains, " +
					"and creates a new one — because the API cannot move a Site to another account.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Human-readable name of the Site. 1..100 characters, without leading or trailing whitespace.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(1, 100), validateAbpName),
			},
			"industry": {
				Description:  "Industry the Site operates in. Informational metadata, unused by the runtime. One of: " + strings.Join(abpSiteIndustries, ", ") + ".",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(abpSiteIndustries, false),
			},
			"default_max_requests_per_minute": {
				Description:  "Default maximum number of requests without a token per minute. Used by every selector whose analysis settings keep `use_site_rate_limiting_parameters = true` (the default).",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"default_max_requests_per_session": {
				Description:  "Default maximum number of requests without a token per session. Used by every selector whose analysis settings keep `use_site_rate_limiting_parameters = true` (the default).",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"default_max_session_length": {
				Description: "Default maximum length of a session without a token. Used by every selector whose analysis settings keep `use_site_rate_limiting_parameters = true` (the default). Duration format, e.g. \"2d1h\". Use the largest fitting units and omit zero components: write \"2h\", not \"120m\".",
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
							Description:  "Policy applied when this Selector matches. Omit to apply no Policy (e.g. for static assets).",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"purpose": {
							Description:  "What this Selector matches, e.g. the user flow or endpoint kind. Informational metadata, unused by the runtime. One of: " + strings.Join(abpSelectorPurposes, ", ") + ". `unknown` means \"explicitly not known\", which is distinct from omitting the attribute.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(abpSelectorPurposes, false),
						},
						"kind": {
							Description: "Match criteria for this Selector. Exactly one of `path_prefix`, `path_regex`, `postback`, `api_security_group` must be set.",
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_prefix": {
										Description: "Match requests whose path begins with this prefix. Mutually exclusive with `path_regex`, `postback` and `api_security_group`.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"case_insensitive": {
										Description: "Match `path_prefix` without regard to case. Applies to `path_prefix` only.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"path_regex": {
										Description: "Match requests whose path matches this regular expression. Mutually exclusive with `path_prefix`, `postback` and `api_security_group`.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"postback": {
										Description:  "Match a specific Postback request type. One of: web_interrogation, ios_interrogation, web_automation, android_interrogation. Mutually exclusive with `path_prefix`, `path_regex` and `api_security_group`.",
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"web_interrogation", "ios_interrogation", "web_automation", "android_interrogation"}, false),
									},
									"api_security_group": {
										Description: "Match requests belonging to this API Security group. Mutually exclusive with `path_prefix`, `path_regex` and `postback`.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
						"analysis_settings": {
							Description:  "JSON-encoded analysis settings for this selector, typically produced by an `incapsula_abp_site_analysis_settings` data source. Use that data source rather than hand-written JSON: the value is compared as a string, so key order and whitespace must match its output exactly.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"derived_id": {
							Description: "Server-derived ID, recomputed from the Selector's fields whenever any of them change. ABP's request counters are reported against it.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "RFC3339 timestamp at which the Selector was created.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"modified_at": {
							Description: "RFC3339 timestamp at which the Selector was last modified.",
							Type:        schema.TypeString,
							Computed:    true,
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
						"case_insensitive": {
							Description: "Whether `path_prefix` is matched without regard to case.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"purpose": {
							Description: "What this Selector matches. Unset for the generated default.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"analysis_settings": {
							Description: "JSON-encoded analysis settings of the default selector.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"derived_id": {
							Description: "Server-derived ID the Selector's request counters are reported against.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "RFC3339 timestamp at which the Selector was created.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"modified_at": {
							Description: "RFC3339 timestamp at which the Selector was last modified.",
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
	if purpose, ok := raw["purpose"].(string); ok && purpose != "" {
		sel.Purpose = &purpose
	}
	return sel, nil
}

func extractAbpSelectorCriteria(raw map[string]any) (AbpSelectorCriteria, error) {
	prefix, _ := raw["path_prefix"].(string)
	regex, _ := raw["path_regex"].(string)
	postback, _ := raw["postback"].(string)
	group, _ := raw["api_security_group"].(string)
	caseInsensitive, _ := raw["case_insensitive"].(bool)

	set := 0
	c := AbpSelectorCriteria{}
	if prefix != "" {
		c.PathPrefix = &prefix
		c.CaseInsensitive = caseInsensitive
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
	if group != "" {
		c.ApiSecurityGroup = &group
		set++
	}
	if set != 1 {
		return c, fmt.Errorf("selector requires exactly one of path_prefix, path_regex, postback, api_security_group (got %d)", set)
	}
	if caseInsensitive && prefix == "" {
		return c, fmt.Errorf("case_insensitive applies to path_prefix only")
	}
	return c, nil
}

func extractAbpSite(data *schema.ResourceData) (AbpSite, error) {
	site := AbpSite{
		Name:      data.Get("name").(string),
		Selectors: []AbpSelector{},
	}
	if v, ok := data.GetOk("industry"); ok {
		s := v.(string)
		site.Industry = &s
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

func serializeAbpSite(data *schema.ResourceData, site *AbpSite) error {
	flat, err := flattenAbpSite(site)
	if err != nil {
		return err
	}

	delete(flat, "id")

	for _, key := range slices.Sorted(maps.Keys(flat)) {
		if err := data.Set(key, flat[key]); err != nil {
			return fmt.Errorf("setting %s: %w", key, err)
		}
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
