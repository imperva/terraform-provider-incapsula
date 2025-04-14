package incapsula

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceSiteDomainConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDomainRead,
		Create: resourceDomainsCreate,
		Delete: resourceDomainDelete,
		Update: resourceDomainsUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"cname_redirection_record": {
				Description: "Cname record for traffic redirection. Point your domain's DNS to this record in order to forward the traffic to Imperva",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:      "A set of domains.",
				Required:         true,
				Type:             schema.TypeSet,
				Set:              getHashFromDomain,
				DiffSuppressFunc: deprecatedFlagDiffSuppress(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Description:      "Domain name",
							Required:         true,
							DiffSuppressFunc: deprecatedFlagDiffSuppress(),
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "Numeric identifier of the domain on Imperva service",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the domain. Indicates if domain DNS is pointed to Imperva's CNAME. Options: BYPASSED, VERIFIED, PROTECTED, MISCONFIGURED",
							Computed:    true,
						},
						"a_records": {
							Description: "A records for traffic redirection. Point your apex domain's DNS to this record in order to forward the traffic to Imperva",
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						}}},
			},
			"deprecated": {
				Description: "Use 'true' to deprecate this resource, any change on the resource will not take effect. Deleting the resource will not delete the domains. Default value: false",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldBool, _ := strconv.ParseBool(old)
					newBool, _ := strconv.ParseBool(new)
					return oldBool == false && newBool == false
				},
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.HasChange("deprecated") {
				oldVal, newVal := d.GetChange("deprecated")
				if oldVal.(bool) && !newVal.(bool) {
					return fmt.Errorf("deprecated flag cannot be changed from true to false")
				}
			}
			return nil
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(3 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
		},
	}
}

func getHashFromDomain(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return PositiveHash(buf.String())
}

func populateFromResourceToDTO(d *schema.ResourceData) []SiteDomainDetails {
	domainsConf := d.Get("domain").(*schema.Set)
	var siteDomainDetails = make([]SiteDomainDetails, len(domainsConf.List()))
	var domainInd = 0
	for _, domain := range domainsConf.List() {
		domainItem := domain.(map[string]interface{})
		siteDomainDetails[domainInd] = SiteDomainDetails{}
		if attr, ok := domainItem["name"]; ok && attr != "" {
			siteDomainDetails[domainInd].Domain = attr.(string)
		}
		domainInd++
	}
	return siteDomainDetails
}

func resourceDomainsCreate(d *schema.ResourceData, m interface{}) error {
	if d.Get("deprecated").(bool) {
		return fmt.Errorf("cannot create deprecated resource")
	}
	return resourceDomainsUpdate(d, m)
}

func resourceDomainsUpdate(d *schema.ResourceData, m interface{}) error {
	if d.Get("deprecated").(bool) {
		return nil
	}
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	siteDomainDetails := populateFromResourceToDTO(d)
	err := client.BulkUpdateDomainsToSite(siteID, siteDomainDetails)
	if err != nil {
		return err
	}
	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	if d.Get("deprecated").(bool) {
		fmt.Printf("[WARN] Resource incapsula_site_domain_configuration is deprecated. Any future changes will be ignored.\n")
		return nil
	}
	client := m.(*Client)
	siteDomainDetailsDto, err := client.GetWebsiteDomains(d.Get("site_id").(string))
	if err != nil {
		return err
	}

	if siteDomainDetailsDto.Errors != nil && len(siteDomainDetailsDto.Errors) > 0 {
		if siteDomainDetailsDto.Errors[0].Status == 404 || siteDomainDetailsDto.Errors[0].Status == 401 {
			log.Printf("[INFO] Operation not allowed: %s\n", siteDomainDetailsDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(siteDomainDetailsDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting domains for site (%s): %s", d.Get("site_id"), string(out))
	}

	domains := &schema.Set{F: getHashFromDomain}
	for _, v := range siteDomainDetailsDto.Data {
		managedCertSettingsID := d.Get("managed_certificate_settings_id")
		hasNoSiteCert := managedCertSettingsID == "" || managedCertSettingsID == nil
		if v.AutoDiscovered == true || (hasNoSiteCert && v.MainDomain == true) { //we ignore the main and auto discovered domains
			continue
		}
		domain := map[string]interface{}{}
		domain["name"] = v.Domain
		domain["id"] = v.Id
		domain["status"] = v.Status
		if v.ARecords != nil && len(v.ARecords) > 0 {
			domain["a_records"] = v.ARecords
		}
		domains.Add(domain)
	}
	d.SetId(d.Get("site_id").(string))
	if len(siteDomainDetailsDto.Data) > 0 {
		d.Set("cname_redirection_record", siteDomainDetailsDto.Data[0].CnameRedirectionRecord)
		d.Set("domain", domains)
	}

	return nil
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	if d.Get("deprecated").(bool) {
		d.SetId("")
		return nil
	}
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	err := client.BulkUpdateDomainsToSite(siteID, []SiteDomainDetails{})
	if err != nil {
		return err
	}

	return resourceDomainRead(d, m)
}
