package incapsula

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceManagedCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceManagedCertificateAdd,
		ReadContext:   resourceManagedCertificateRead,
		UpdateContext: resourceManagedCertificateAdd,
		DeleteContext: resourceManagedCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				siteId := d.Id()

				d.Set("site_id", siteId)
				log.Printf("[DEBUG] cloudwaf site resource: Import  Site Config JSON for Site ID %s", siteId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default_validation_method": {
				Description: "The default validation method.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "CNAME",
			},
			"account_id": {
				Description: "(Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceManagedCertificateAdd(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteId, _ = d.Get("site_id").(string)
	var accountId, _ = d.Get("account_id").(int)
	validationMethod, _ := d.Get("default_validation_method").(string)
	id, _ := strconv.Atoi(siteId)
	log.Printf("[INFO] requesting site cert to site ID: %d to %v", id, d)
	siteCertificateV3Response, diags := client.RequestSiteCertificate(id, validationMethod, accountId)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed request site cert to site ID: %d, %v\n", id, diags)
		return diags
	} else if siteCertificateV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to request site cert to site ID: %d, %v\n", id, siteCertificateV3Response.Errors[0].Detail)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to request site cert to site",
			Detail:   fmt.Sprintf("Failed to request site cert to site ID%d, %s", id, siteCertificateV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("site_id", strconv.Itoa(siteCertificateV3Response.Data[0].SiteId))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site id after delete v3 site: %s\n", err)
		return diag.FromErr(err)
	}
	resourceManagedCertificateRead(ctx, d, m)
	return diags
}

func resourceManagedCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteId, _ = d.Get("site_id").(string)
	id, _ := strconv.Atoi(siteId)
	log.Printf("[INFO] get site cert status of site ID: %d to %v", id, d)
	siteCertificateV3Response, diags := client.GetSiteCertificateRequestStatus(id)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed get site cert status of site ID: %d, %v\n", id, diags)
		return diags
	} else if siteCertificateV3Response.Errors != nil {
		log.Printf("[ERROR] Failed get site cert status of site ID: %d, %v\n", id, siteCertificateV3Response.Errors[0].Detail)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to request site cert to site",
			Detail:   fmt.Sprintf("Failed to get site cert status of site ID%d, %s", id, siteCertificateV3Response.Errors[0].Detail),
		}}
	}
	siteId := siteCertificateV3Response.Data[0].SiteId
	err := d.Set("site_id", strconv.Itoa(siteId))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site id after request site cert to site ID: %d, %s\n", id, err)
		return diag.FromErr(err)
	}
	err = d.Set("default_validation_method", siteCertificateV3Response.Data[0].DefaultValidationMethod)
	if err != nil {
		log.Printf("[ERROR] Could not read Default vlidation method after request site cert to site ID: %d, %s\n", id, err)
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(siteCertificateV3Response.Data[0].SiteId))
	return diags
}

func resourceManagedCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteId, _ = d.Get("site_id").(string)
	id, _ := strconv.Atoi(siteId)
	log.Printf("[INFO] deleting site cert request of site ID: %d to %v", id, d)
	siteCertificateV3Response, diags := client.DeleteRequestSiteCertificate(id)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed delete site cert request of site ID: %d, %v\n", id, diags)
		return diags
	} else if siteCertificateV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to delete site cert request of site ID: %d, %v\n", id, siteCertificateV3Response.Errors[0].Detail)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete site cert request of",
			Detail:   fmt.Sprintf("Failed to delete site cert request of site ID%d, %s", id, siteCertificateV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("site_id", strconv.Itoa(siteCertificateV3Response.Data[0].SiteId))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site id after delete v3 site: %s\n", err)
		return diag.FromErr(err)
	}
	resourceManagedCertificateRead(ctx, d, m)
	return diags
}
