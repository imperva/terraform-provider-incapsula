package incapsula

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpSites() *schema.Resource {
	siteAttributes := map[string]*schema.Schema{
		"id": {
			Description: "ID of the Site.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"account_id": {
			Description: "ABP account UUID the Site belongs to.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name of the Site.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
	maps.Copy(siteAttributes, abpSiteReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpSitesRead,

		Description: "Lists every ABP Site (a.k.a. Website Group) belonging to an account. " +
			"Use this to reference Sites that are not managed by this Terraform " +
			"configuration when all of them are of interest, for example when building " +
			"an `incapsula_abp_account_site_priority` list. Use `incapsula_abp_site` to " +
			"look up a single Site by `site_id` or `name` instead.\n\n" +
			"Note: The time of reading this datasource is not defined by the terraform module. That means that" +
			"the read is subject to race-conditions. If a site is created or deleted in the same plan as this " +
			"datasource is read, it may yield unexpected results, like an inexhaustive listing, which would cause " +
			"e.g. `site_priority` to fail (as that requires an exhaustive list). " +
			"If deterministic operation is required, specify resource dependencies explicitly to force terraform " +
			"plan order, e.g. make `incapsula_abp_account_site_priority` depend on sites that you are creating.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to list the Sites of.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"sites": {
				Description: "All Sites belonging to the account, ordered from most to least significant priority.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: siteAttributes,
				},
			},
		},
	}
}

func dataSourceAbpSitesRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	sites, err := client.ListAbpSites(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	flat := make([]any, 0, len(sites))
	for i := range sites {
		site, err := flattenAbpSite(&sites[i])
		if err != nil {
			return diag.FromErr(fmt.Errorf("ABP Site %s: %w", sites[i].Id, err))
		}
		flat = append(flat, site)
	}

	data.SetId(accountId)
	if err := data.Set("sites", flat); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
