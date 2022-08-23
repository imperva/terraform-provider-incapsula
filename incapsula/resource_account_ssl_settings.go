package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceAccountSSLSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountSSLSettingsUpdate,
		ReadContext:   resourceAccountSSLSettingsRead,
		UpdateContext: resourceAccountSSLSettingsUpdate,
		DeleteContext: resourceAccountSSLSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				accountId := d.Id()

				d.Set("account_id", accountId)
				log.Printf("[DEBUG] account ssl settings resource: Import  Account Config JSON for Account ID %s", accountId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"allow_cname_validation": {
				Description: "Allow Imperva to prove ownership on the domains under the allowed_domains_for_cname_validation list on behalf of the customer.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"value_for_cname_validation": {
				Description: "The value of the CNAME records that need to create for each domain under the allowed_domains_for_cname_validation list to allow delegation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"allowed_domains_for_cname_validation": {
				Description: "The list of domains that Imperva allow to prove ownership on, on behalf of the customer.",
				Type:        schema.TypeSet,
				Optional:    true,
				Default:     nil,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"use_wild_card_san_instead_of_fqdn": {
				Description: "Add wildcard SAN instead of FQDN SAN on the Imperva generated certificate for newly created sites.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"add_naked_domain_san_for_www_sites": {
				Description: "Add naked domain SAN on the Imperva generated certificate for newly created WWW sites.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceAccountSSLSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	accountID := d.Id()
	if d.Get("account_id") != nil {
		accountID, _ = d.Get("account_id").(string)
	}
	accountSSLSettingsDTO := AccountSSLSettingsDTO{}
	log.Printf("[INFO] Updating Incapsula account SSL settings for Account ID: %s to %v", accountID, d)
	impervaCertificateDto := ImpervaCertificate{}
	if d.Get("add_naked_domain_san_for_www_sites") != nil {
		fieldVal := d.Get("add_naked_domain_san_for_www_sites").(bool)
		impervaCertificateDto.AddNakedDomainSanForWWWSites = &fieldVal
	}
	if d.Get("use_wild_card_san_instead_of_fqdn") != nil {
		fieldVal := d.Get("use_wild_card_san_instead_of_fqdn").(bool)
		impervaCertificateDto.UseWildCardSanInsteadOfFQDN = &fieldVal
	}
	delegationDto := Delegation{}
	if d.Get("allow_cname_validation") != nil {
		fieldVal := d.Get("allow_cname_validation").(bool)
		delegationDto.AllowCNAMEValidation = &fieldVal
	}
	if d.Get("allowed_domains_for_cname_validation") != nil {
		domains := d.Get("allowed_domains_for_cname_validation").(*schema.Set).List()
		domainsList := make([]string, len(domains))
		for i, k := range domains {
			domainsList[i] = k.(string)
		}
		delegationDto.AllowedDomainsForCNAMEValidation = domainsList
	}
	impervaCertificateDto.Delegation = &delegationDto
	accountSSLSettingsDTO.ImpervaCertificate = &impervaCertificateDto
	accountSSLSettingsDTOResponse, diags := client.UpdateAccountSSLSettings(&accountSSLSettingsDTO, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Could not update Incapsula account SSL settings for Account ID: %s, %v\n", accountID, diags)
		return diags
	} else if accountSSLSettingsDTOResponse.Errors != nil {
		log.Printf("[ERROR] Failed to update Incapsula account SSL settings for Account ID: %s, %v\n", accountID, accountSSLSettingsDTOResponse.Errors[0].Detail)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to update account SSL settings",
			Detail:   fmt.Sprintf("Failed to update account SSL settings for account%s, %s", accountID, accountSSLSettingsDTOResponse.Errors[0].Detail),
		}}
	}
	err := d.Set("account_id", accountID)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account SSL settings after update for Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	if err != nil {
		log.Printf("[ERROR] Could not update last_update field of Incapsula account SSL settings resource for Account ID: %s, %s\n", accountID, err)
		return diag.FromErr(err)
	}
	resourceAccountSSLSettingsRead(ctx, d, m)

	return diags
}

func resourceAccountSSLSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	accountID := d.Id()
	if d.Get("account_id") != nil {
		accountID, _ = d.Get("account_id").(string)
	}

	log.Printf("[INFO] Reading Incapsula account SSL settings for Account ID: %s\n", accountID)

	accountSSLSettingsDTOResponse, diagFromClient := client.GetAccountSSLSettings(accountID)

	if diagFromClient != nil && diagFromClient.HasError() {
		log.Printf("[ERROR] Could not read Incapsula account SSL settings for Account ID: %s, %v\n", accountID, diagFromClient)
		return diagFromClient
	}
	accountSSLSettingsDTO := accountSSLSettingsDTOResponse.Data[0]
	if err := d.Set("use_wild_card_san_instead_of_fqdn", accountSSLSettingsDTO.ImpervaCertificate.UseWildCardSanInsteadOfFQDN); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("add_naked_domain_san_for_www_sites", accountSSLSettingsDTO.ImpervaCertificate.AddNakedDomainSanForWWWSites); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_domains_for_cname_validation", accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_cname_validation", accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowCNAMEValidation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value_for_cname_validation", accountSSLSettingsDTO.ImpervaCertificate.Delegation.ValueForCNAMEValidation); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(accountID)
	return diags
}

func resourceAccountSSLSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diag diag.Diagnostics
	accountID := d.Id()
	if d.Get("account_id") != nil {
		accountID, _ = d.Get("account_id").(string)
	}

	log.Printf("[INFO] Reseting Incapsula account SSL settings for Account ID: %s\n", accountID)

	diag = client.DeleteAccountSSLSettings(accountID)

	if diag != nil && diag.HasError() {
		log.Printf("[ERROR] Could not delete Incapsula account SSL settings for Account ID: %s, %v\n", accountID, diag[0].Detail)
		return diag
	}

	d.SetId("")
	return diag
}
