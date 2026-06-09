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

		Description: "Looks up an existing ABP Site (a.k.a. Website Group) by its ID. Use this " +
			"to reference a Site that is not managed by this Terraform configuration, " +
			"for example when building an `incapsula_abp_account_site_priority` list.",

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description:  "ID of the Site to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"account_id": {
				Description: "ABP account UUID this Site belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Human-readable name of the Site.",
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
	siteId := data.Get("site_id").(string)

	site, err := client.ReadAbpSite(siteId)
	if err != nil {
		return diag.FromErr(err)
	}
	if site == nil {
		return diag.Errorf("no ABP Site with id %q found", siteId)
	}

	data.SetId(site.Id)
	if err := serializeAbpSite(data, site); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
