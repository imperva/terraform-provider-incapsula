package incapsula

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncapsulaGcpRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncapsulaGcpRegionsRead,
		Description: "Retrieves the supported Google Cloud regions for Imperva for Google Cloud. Use a region from this list when retrieving a PSC target service address.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account ID (caid) used to scope the request. When omitted, the account of the API credentials is used.",
				Optional:    true,
			},
			"regions": {
				Type:        schema.TypeList,
				Description: "The supported Google Cloud region identifiers, sorted alphabetically.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIncapsulaGcpRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountID := d.Get("account_id").(string)

	resp, err := client.GetGcpRegions(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	sorted := append([]string(nil), resp.Data...)
	sort.Strings(sorted)
	regions := make([]interface{}, 0, len(sorted))
	for _, r := range sorted {
		regions = append(regions, r)
	}

	if err := d.Set("regions", regions); err != nil {
		return diag.FromErr(err)
	}

	if accountID != "" {
		d.SetId("gcp-regions-" + accountID)
	} else {
		d.SetId("gcp-regions")
	}

	return nil
}
