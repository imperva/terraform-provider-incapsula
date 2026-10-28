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
			"cloud_type": {
				Description: "The cloud provider type of the site, AWS or GCP. Empty for sites that are not PUBLIC_CLOUD.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_load_balancer_site": {
				Description: "Indicates whether the site is an Imperva for Google Cloud load balancer website. False for domain websites, legacy Imperva for Google Cloud default sites, and sites that are not part of Imperva for Google Cloud.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"authority_header": {
				Description: "The authority header assigned by Imperva for an Imperva for Google Cloud load balancer website. Available for load balancer websites and legacy Imperva for Google Cloud default sites. Empty for domain websites and sites that are not part of Imperva for Google Cloud.",
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
	if accountID != 0 {
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

	if err := d.Set("cloud_type", site.CloudType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_load_balancer_site", site.IsLoadBalancerSite != nil && *site.IsLoadBalancerSite); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("authority_header", site.AuthorityHeader); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
