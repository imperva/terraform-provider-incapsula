package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpSiteRead,

		Description: "Looks up an existing ABP Site (a.k.a. Website Group) within an account " +
			"by its `site_id`, its `name`, or both. At least one of `site_id` or `name` " +
			"must be set; when both are given, a Site must match both. Use this to " +
			"reference a Site that is not managed by this Terraform configuration, for " +
			"example when building an `incapsula_abp_account_site_priority` list. The " +
			"lookup fails if zero or more than one Site matches.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to search within.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"site_id": {
				Description:  "ID of the Site to look up. At least one of `site_id` or `name` is required.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				AtLeastOneOf: []string{"site_id", "name"},
			},
			"name": {
				Description:  "Name of the Site to look up. Matched exactly and case-sensitively. At least one of `site_id` or `name` is required.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"site_id", "name"},
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
				Description: "Default maximum length of a session without a token, in moi duration format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"selector": {
				Description: "Ordered list of user-defined Selectors for this Site.",
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
							Description: "Policy applied when this Selector matches.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path_prefix": {
							Description: "Match requests whose path begins with this prefix.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path_regex": {
							Description: "Match requests whose path matches this regular expression.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"postback": {
							Description: "Match a specific Postback request type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"analysis_settings": {
							Description: "JSON-encoded analysis settings for this selector.",
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

func dataSourceAbpSiteRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	siteId := data.Get("site_id").(string)
	name := data.Get("name").(string)

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
