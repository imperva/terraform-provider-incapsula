package incapsula

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncapsulaGcpPscTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncapsulaGcpPscTargetRead,
		Description: "Retrieves the PSC target service address for a specified Google Cloud region.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Description: "The Google Cloud region for which you want to retrieve the PSC target service address.",
				Required:    true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account ID (caid) used to scope the request. When omitted, the account of the API credentials is used.",
				Optional:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The PSC target service address for the specified Google Cloud region.",
				Computed:    true,
			},
		},
	}
}

func dataSourceIncapsulaGcpPscTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	region := d.Get("region").(string)
	accountID := d.Get("account_id").(string)

	resp, err := client.GetPscTarget(region, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("address", resp.Address); err != nil {
		return diag.FromErr(err)
	}
	if accountID != "" {
		d.SetId(accountID + "/" + region)
	} else {
		d.SetId(region)
	}

	return nil
}
