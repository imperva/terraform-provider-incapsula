package incapsula

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomainsValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSLValidationAdd,
		ReadContext:   resourceSSLValidationRead,
		UpdateContext: resourceSSLValidationAdd,
		DeleteContext: resourceSSLValidationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				siteId := d.Id()

				d.Set("site_id", siteId)
				log.Printf("[DEBUG] site v3 resource: Import  Site Config JSON for Site ID %s", siteId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"domain_ids": {
				Description: "domain ids.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSSLValidationAdd(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteId := d.Get("site_id").(int)
	log.Printf("[INFO] requesting site cert to site ID: %d to %v", siteId, d)
	domains := d.Get("domain_ids").(*schema.Set)
	var index = 0
	var domainIds = make([]int, len(domains.List()))
	for _, domain := range domains.List() {
		dom, _ := strconv.Atoi(domain.(string))
		domainIds[index] = dom
		index++
	}
	for i := 0; i < 50; i++ {
		siteCertificateV3Response, _ := client.GetSiteCertificateRequestStatus(siteId, nil)
		b := siteCertificateV3Response != nil && siteCertificateV3Response.Data != nil && len(siteCertificateV3Response.Data) > 0
		b = b && siteCertificateV3Response.Data[0].CertificatesDetails != nil && len(siteCertificateV3Response.Data[0].CertificatesDetails) > 0
		b = b && siteCertificateV3Response.Data[0].CertificatesDetails[0].Sans != nil && len(siteCertificateV3Response.Data[0].CertificatesDetails[0].Sans) == len(domainIds)
		for _, value := range siteCertificateV3Response.Data[0].CertificatesDetails[0].Sans {
			b = b && value.Status == "VALIDATED"
		}
		if b {
			diags = client.ValidateDomains(siteId, domainIds)
			d.SetId(strconv.Itoa(siteId))
			return diags
		}
		client.ValidateDomains(siteId, domainIds)
		time.Sleep(10 * time.Second)
	}
	panic("ssl validation was not completed as expected")
}

func resourceSSLValidationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	siteId, _ := d.Get("site_id").(int)
	d.SetId(strconv.Itoa(siteId))
	return nil
}

func resourceSSLValidationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
