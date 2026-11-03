package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncapsulaGcpDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncapsulaGcpDomainsRead,
		Description: "Retrieves the domains associated with a Google Cloud website in Imperva. For a load balancer website, returns the domains protected through that website. For a domain website, returns the domain managed separately by that website.",

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:        schema.TypeString,
				Description: "The numeric identifier of the site whose domains should be listed.",
				Required:    true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account ID (caid) used to scope the request. When omitted, the account of the API credentials is used.",
				Optional:    true,
			},
			"domains": {
				Type:        schema.TypeList,
				Description: "The domains attached to the site.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"a_records": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceIncapsulaGcpDomainsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	siteId := d.Get("site_id").(string)
	accountID := d.Get("account_id").(string)

	resp, err := client.GetSiteDomains(siteId, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	domains := make([]interface{}, 0, len(resp.Data))
	for _, dom := range resp.Data {
		domains = append(domains, map[string]interface{}{
			"id":            dom.Id,
			"domain":        dom.Domain,
			"status":        dom.Status,
			"cname":         dom.Cname,
			"creation_date": dom.CreationDate,
			"a_records":     dom.ARecords,
		})
	}

	if err := d.Set("domains", domains); err != nil {
		return diag.FromErr(err)
	}
	if accountID != "" {
		d.SetId(accountID + "/" + siteId)
	} else {
		d.SetId(siteId)
	}

	return nil
}
