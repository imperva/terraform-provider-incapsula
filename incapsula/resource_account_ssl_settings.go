package incapsula

import (
	"bytes"
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
				Optional:    true,
				Computed:    true,
			},
			"allowed_domain_for_cname_validation": {
				Description: "The list of domains that Imperva allow to prove ownership on, on behalf of the customer.",
				Type:        schema.TypeSet,
				Optional:    true,
				Default:     nil,
				Set:         domainUniqueId,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The domain id.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The domain name.",
							Required:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "The domain status.",
							Computed:    true,
						},
						"creation_date": {
							Type:        schema.TypeInt,
							Description: "The domain creation date.",
							Computed:    true,
						},
						"status_since": {
							Type:        schema.TypeInt,
							Description: "The domain status since date.",
							Computed:    true,
						},
						"last_status_check": {
							Type:        schema.TypeInt,
							Description: "The domain last status check date.",
							Computed:    true,
						},
						"cname_record_host": {
							Type:        schema.TypeString,
							Description: "The CNAME record value to use to configure this domain for delegation.",
							Computed:    true,
						},
						"cname_record_value": {
							Type:        schema.TypeString,
							Description: "The CNAME record host to use.",
							Computed:    true,
						},
					},
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
			"allow_support_old_tls_versions": {
				Description: "When true, sites under the account or sub-accounts can allow support of old TLS versions traffic. This can be configured only on the parent account level.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enable_hsts_for_new_sites": {
				Description: "When true, enables HSTS support for newly created websites.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
	if d.Get("allowed_domain_for_cname_validation") != nil {
		domains := d.Get("allowed_domain_for_cname_validation").(*schema.Set).List()
		domainsList := make([]AllowDomainForCnameValidation, len(domains))
		for i, k := range domains {
			updateDomainList(domainsList, i, k, accountID)
		}
		delegationDto.AllowedDomainsForCNAMEValidation = domainsList
	}
	impervaCertificateDto.Delegation = &delegationDto
	accountSSLSettingsDTO.ImpervaCertificate = &impervaCertificateDto
	if d.Get("allow_support_old_tls_versions") != nil {
		fieldVal := d.Get("allow_support_old_tls_versions").(bool)
		accountSSLSettingsDTO.AllowSupportOldTLSVersions = &fieldVal
	}
	if d.Get("enable_hsts_for_new_sites") != nil {
		fieldVal := d.Get("enable_hsts_for_new_sites").(bool)
		accountSSLSettingsDTO.EnableHSTSForNewSites = &fieldVal
	}
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

func updateDomainList(domainList []AllowDomainForCnameValidation, counter int, k interface{}, accountID string) diag.Diagnostics {
	al := k.(map[string]interface{})
	allowDomainForCnameValidation := AllowDomainForCnameValidation{}
	if attr, ok := al["name"]; ok && attr != "" {
		allowDomainForCnameValidation.Name = attr.(string)
	} else {
		log.Printf("[ERROR] Failed to update Incapsula account SSL settings for Account ID: %s, failed to update domain name field for domain %d\n", accountID, allowDomainForCnameValidation.Id)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to update account SSL settings",
			Detail:   fmt.Sprintf("Failed to update account SSL settings for account%s, failed to update domain name field for domain %d", accountID, allowDomainForCnameValidation.Id),
		}}
	}
	domainList[counter] = allowDomainForCnameValidation
	return nil
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
	if accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation != nil {
		numberOfNotInheritedDomains := 0
		for _, k := range accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation {
			if !k.Inherited {
				numberOfNotInheritedDomains++
			}
		}
		domainsList := make([]map[string]interface{}, numberOfNotInheritedDomains)
		counter := 0
		for _, k := range accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation {
			if !k.Inherited {
				domain := map[string]interface{}{}
				domain["id"] = k.Id
				domain["name"] = k.Name
				domain["status"] = k.Status
				domain["creation_date"] = k.CreationDate
				domain["status_since"] = k.StatusSince
				domain["last_status_check"] = k.LastStatusCheck
				domain["cname_record_value"] = k.CnameRecordValue
				domain["cname_record_host"] = k.CnameRecordHost
				domainsList[counter] = domain
				counter++
			}
		}
		if err := d.Set("allowed_domain_for_cname_validation", domainsList); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("allow_cname_validation", accountSSLSettingsDTO.ImpervaCertificate.Delegation.AllowCNAMEValidation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value_for_cname_validation", accountSSLSettingsDTO.ImpervaCertificate.Delegation.ValueForCNAMEValidation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_hsts_for_new_sites", accountSSLSettingsDTO.EnableHSTSForNewSites); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_support_old_tls_versions", accountSSLSettingsDTO.AllowSupportOldTLSVersions); err != nil {
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

func domainUniqueId(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return PositiveHash(buf.String())
}
