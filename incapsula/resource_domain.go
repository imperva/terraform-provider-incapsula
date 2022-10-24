package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDomainRead,
		Create: resourceDomainCreate,
		Delete: resourceDomainDelete,
		Update: resourceDomainUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	domain := d.Get("domain").(string)

	domainManagementData, err := client.AddDomainToSite(siteID, domain)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(domainManagementData.Id))
	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteDomainDetails, err := client.GetDomainDetails(d.Get("site_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.Set("site_id", siteDomainDetails.SiteId)

	return nil
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}
