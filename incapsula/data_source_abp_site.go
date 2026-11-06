package incapsula

import (
	"context"
	"maps"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpSite() *schema.Resource {
	siteSchema := map[string]*schema.Schema{
		"account_id": {
			Description:  "ABP account UUID to search within.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
		"site_id": {
			Description:  "ID of the Site to look up. One of (but not both) `site_id` or `name` is required.",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsUUID,
			AtLeastOneOf: []string{"site_id", "name"},
		},
		"name": {
			Description:  "Name of the Site to look up. Matched exactly and case-sensitively. One of (but not both) `site_id` or `name` is required.",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			AtLeastOneOf: []string{"site_id", "name"},
		},
	}
	maps.Copy(siteSchema, abpSiteReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpSiteRead,

		Description: "Looks up an existing ABP Site (a.k.a. Website Group) within an account " +
			"by its `site_id`, or (XOR) its `name`. Exactly one of `site_id` or `name` " +
			"must be set. Use this to " +
			"reference a Site that is not managed by this Terraform configuration, for " +
			"example when building an `incapsula_abp_account_site_priority` list. The " +
			"lookup fails if zero or more than one Site matches. Use " +
			"`incapsula_abp_sites` to retrieve every Site of an account instead.",

		Schema: siteSchema,
	}
}

// abpSiteReadOnlyAttributes returns the read-only Site attributes shared by the
// `incapsula_abp_site` and `incapsula_abp_sites` data sources, keeping the two
// in lockstep.
func abpSiteReadOnlyAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"industry": {
			Description: "Industry the Site operates in. Informational metadata, unused by the runtime.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"default_max_requests_per_minute": {
			Description: "Default maximum number of requests without a token per minute.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"default_max_requests_per_session": {
			Description: "Default maximum number of requests without a token per session.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"default_max_session_length": {
			Description: "Default maximum length of a session without a token. Duration format, e.g. \"2d1h\".",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"selector": {
			Description: "Ordered list of user-defined Selectors for this Site. Each Selector " +
				"matches requests one way: exactly one of the `path_prefix`, `path_regex`, " +
				"`postback` and `api_security_group` blocks is present.",
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "Server-assigned Selector ID.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"policy_id": {
						Description: "Policy applied when this Selector matches.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"purpose": {
						Description: "What this Selector matches, e.g. the user flow or endpoint kind. Informational metadata, unused by the runtime. Empty when unset.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"path_prefix": {
						Description: "Present when the Selector matches requests whose path begins with a prefix.",
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Description: "The path prefix that is matched.",
									Type:        schema.TypeString,
									Computed:    true,
								},
								"case_insensitive": {
									Description: "Whether the prefix is matched without regard to case.",
									Type:        schema.TypeBool,
									Computed:    true,
								},
							},
						},
					},
					"path_regex": {
						Description: "Present when the Selector matches requests whose path matches a regular expression.",
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Description: "The regular expression the path is matched against.",
									Type:        schema.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"postback": {
						Description: "Present when the Selector matches a Postback request type.",
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Description: "The Postback request type that is matched. One of: " + strings.Join(abpSelectorPostbacks, ", ") + ".",
									Type:        schema.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"api_security_group": {
						Description: "Present when the Selector matches requests belonging to an API Security group.",
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Description: "Name of the API Security group that is matched.",
									Type:        schema.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"analysis_settings": {
						Description: "JSON-encoded analysis settings for this selector.",
						Type:        schema.TypeString,
						Computed:    true,
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
			Description: "Catch-all selector matching all request paths after all other selectors.",
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
						Description: "The path prefix the default Selector matches, always `/`.",
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Description: "Always `/`.",
									Type:        schema.TypeString,
									Computed:    true,
								},
								"case_insensitive": {
									Description: "Whether the prefix is matched without regard to case.",
									Type:        schema.TypeBool,
									Computed:    true,
								},
							},
						},
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
	}
}

func dataSourceAbpSiteRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	siteId := data.Get("site_id").(string)
	name := data.Get("name").(string)
	if siteId != "" && name != "" {
		return diag.Errorf("exactly one of site_id or name should be set, both were set")
	}

	sites, err := client.ListAbpSites(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	var match *AbpSite
	for i := range sites {
		if siteId != "" && sites[i].Id != siteId {
			continue
		}
		if name != "" && sites[i].Name != name {
			continue
		}
		if match != nil {
			return diag.Errorf("multiple ABP Sites match the given criteria in account %s; refine site_id/name", accountId)
		}
		match = &sites[i]
	}
	if match == nil {
		return diag.Errorf("no ABP Site matching the given criteria found in account %s", accountId)
	}

	data.SetId(match.Id)
	if err := data.Set("site_id", match.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := serializeAbpSite(data, match); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
