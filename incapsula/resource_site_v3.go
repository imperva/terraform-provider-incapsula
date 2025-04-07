package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSiteV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiteV3Add,
		ReadContext:   resourceSiteV3Read,
		UpdateContext: resourceSiteV3Update,
		DeleteContext: resourceSiteV3Delete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

				idSlice := strings.Split(d.Id(), "/")
				log.Printf("[DEBUG] Starting to import site v3. Parameters: %s\n", d.Id())

				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected account_id/site_id", d.Id())
				}

				err := d.Set("account_id", idSlice[0])

				if err != nil {
					return nil, err
				}

				_, err = strconv.Atoi(idSlice[1])
				if err != nil {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected account_id/site_id", d.Id())
				}

				d.SetId(idSlice[1])

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "The site name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The website type. Indicates which kind of website is created, e.g. CLOUD_WAF for a website onboarded to Imperva Cloud WAF.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "CLOUD_WAF",
			},
			"creation_time": {
				Description: "Creation time of the site.",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"cname": {
				Description: "The CNAME provided by Imperva that is used for pointing your website traffic to the Imperva network.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ref_id": {
				Description: "(Optional) Sets the Reference ID. A free-text field that enables you to add a unique identifier to correlate a website in our service with an object on the customer side.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "(Optional) Whether the site is active or bypassed by the Imperva network.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceSiteV3Add(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	log.Printf("[DEBUG] ************************* test2")
	accountID, _ := d.Get("account_id").(string)
	log.Printf("[INFO] adding v3 site to Account ID: %s to %v", accountID, d)
	siteV3Request := SiteV3Request{}
	siteV3Request.SiteType = d.Get("type").(string)
	siteV3Request.Name = d.Get("name").(string)
	siteV3Request.AccountId, _ = strconv.Atoi(accountID)
	if d.Get("ref_id") != nil {
		siteV3Request.RefId = d.Get("ref_id").(string)
	}
	siteV3Request.Active = d.Get("active").(bool)
	siteV3Response, diags := client.AddV3Site(&siteV3Request, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to add v3 site to Account ID: %s, %v\n", accountID, diags)
		return diags
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to add v3 site to Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to add v3 site",
			Detail:   fmt.Sprintf("Failed to add v3 site to account%s, %s", accountID, siteV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("account_id", strconv.Itoa(siteV3Response.Data[0].AccountId))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account after add v3 site to Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	siteId := siteV3Response.Data[0].Id
	d.SetId(strconv.Itoa(siteId))

	resourceSiteV3Read(ctx, d, m)
	return diags
}

func resourceSiteV3Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	accountID, _ := d.Get("account_id").(string)

	log.Printf("[INFO] adding v3 site to Account ID: %s to %v", accountID, d)
	siteV3Request := SiteV3Request{}
	siteV3Request.Name = d.Get("name").(string)
	siteV3Request.Id, _ = strconv.Atoi(d.Id())
	siteV3Request.RefId = ""
	if d.Get("ref_id") != nil {
		siteV3Request.RefId = d.Get("ref_id").(string)
	}
	siteV3Request.Active = d.Get("active").(bool)
	siteV3Response, diags := client.UpdateV3Site(&siteV3Request, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to update v3 site to Account ID: %s, %v\n", accountID, diags)
		return diags
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to update v3 site to Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to add v3 site",
			Detail:   fmt.Sprintf("Failed to update v3 site to account%s, %s", accountID, siteV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("account_id", accountID)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account after update v3 site to Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	siteId := siteV3Response.Data[0].Id
	d.SetId(strconv.Itoa(siteId))

	resourceSiteV3Read(ctx, d, m)
	return diags
}

func resourceSiteV3Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountID, _ := d.Get("account_id").(string)
	log.Printf("[DEBUG] ************************* test")
	log.Printf("[INFO] getting v3 site of Account ID: %s to %v", accountID, d)
	siteV3Request := SiteV3Request{}
	siteV3Request.SiteType = d.Get("type").(string)
	siteV3Request.Name = d.Get("name").(string)
	siteV3Request.AccountId, _ = strconv.Atoi(accountID)
	siteV3Request.Id, _ = strconv.Atoi(d.Id())
	siteV3Response, diags := client.GetV3Site(&siteV3Request, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to get v3 site of Account ID: %s, %v\n", accountID, diags)
		return diags
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to get v3 site of Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to get v3 site",
			Detail:   fmt.Sprintf("Failed to get v3 site of account%s, %s", accountID, siteV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("account_id", strconv.Itoa(siteV3Response.Data[0].AccountId))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(siteV3Response.Data[0].Id))

	err = d.Set("name", siteV3Response.Data[0].Name)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula name after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	err = d.Set("creation_time", siteV3Response.Data[0].CreationTime)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula creation time after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	err = d.Set("cname", siteV3Response.Data[0].Cname)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula cname after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	err = d.Set("type", siteV3Response.Data[0].SiteType)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula type after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	err = d.Set("ref_id", siteV3Response.Data[0].RefId)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula ref id after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	err = d.Set("active", siteV3Response.Data[0].Active)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula active mode after get v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(siteV3Response.Data[0].Id))
	return diags
}

func resourceSiteV3Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	accountID, _ := d.Get("account_id").(string)

	log.Printf("[INFO] deleting v3 site of Account ID:%s to %v", accountID, d)
	siteV3Request := SiteV3Request{}
	siteV3Request.SiteType = d.Get("type").(string)
	siteV3Request.Name = d.Get("name").(string)
	siteV3Request.AccountId, _ = strconv.Atoi(accountID)
	siteV3Request.Id, _ = strconv.Atoi(d.Id())
	siteV3Response, diags := client.DeleteV3Site(&siteV3Request, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to delete v3 site of Account ID: %s, %v\n", accountID, diags)
		return diags
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to delete v3 site of Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to delete v3 site",
			Detail:   fmt.Sprintf("Failed to delete v3 site of account%s, %s", accountID, siteV3Response.Errors[0].Detail),
		}}
	}
	err := d.Set("account_id", accountID)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account after delete v3 site of Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}

	siteId := strconv.Itoa(siteV3Response.Data[0].Id)
	d.SetId(siteId)

	return nil
}
