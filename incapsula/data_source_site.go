package incapsula

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiteRead,
		Description: "Provides data about user Account",

		Schema: map[string]*schema.Schema{
			// Computed Attributes
			"account_id": {
				Type:        schema.TypeString,
				Description: "Account ID",
				Optional:    true,
			},
			"site_url": {
				Type:        schema.TypeString,
				Description: "site url to find ID of",
				Required:    true,
			},
			// "site_id": {
			// 	Type:        schema.TypeString,
			// 	Description: "site ID",
			// 	Computed:    true,
			// },
		},
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var accountStatusResponse *SiteListResponse
	var err error
	if accountID, ok := d.GetOk("account_id"); ok {
		accountStatusResponse, err = client.ListAllSitesInAccount(accountID.(string))
		if err != nil {
			return diag.Errorf("Error listing all sites: %v", err)
		}
	} else {
		accountStatusResponse, err = client.ListAllSites()
		if err != nil {
			return diag.Errorf("Error listing all sites: %v", err)
		}
	}

	foundSite := false
	for i := range accountStatusResponse.Sites {
		site := accountStatusResponse.Sites[i]
		if site.Domain == d.Get("site_url") {
			d.SetId(strconv.Itoa(site.SiteID))
			d.Set("account_id", strconv.Itoa(site.AccountID))
			foundSite = true
		}
	}

	if !foundSite {
		return diag.Errorf("Error finding site: %v, this domain name does not exist in account: %v", d.Get("site_url"), d.Get("account_id"))
	}

	return nil
}
