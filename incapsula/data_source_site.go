package incapsula

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiteRead,
		Description: "Provides a site ID based on the site name. This allows referencing the same site across multiple Terraform states and accounts.",

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the site to look up.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "Numeric identifier of the account containing the site. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"site_id": {
				Description: "The external ID of the site.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"cname": {
				Description: "The CNAME of the site.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the site.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	siteName := d.Get("name").(string)
	accountID := d.Get("account_id").(int)

	var accountIDStr string
	if accountID > 0 {
		accountIDStr = strconv.Itoa(accountID)
	}

	log.Printf("[DEBUG] dataSourceSiteRead - siteName: %s, accountID: %d, accountIDStr: %s", siteName, accountID, accountIDStr)

	site, diags := client.GetV3SiteByName(siteName, accountIDStr)
	if len(diags) > 0 {
		return diags
	}

	d.SetId(strconv.Itoa(site.Id))
	d.Set("site_id", site.Id)
	d.Set("cname", site.Cname)
	d.Set("type", site.SiteType)

	return nil
}
