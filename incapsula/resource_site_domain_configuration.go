package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func resourceSiteDomainConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDomainRead,
		Create: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Update: resourceDomainUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cname_redirection_record": {
				Description: "Cname record for traffic redirection. Point your domain's DNS to this record in order to forward the traffic to Imperva",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description: "A set of domains.",
				Required:    true,
				Type:        schema.TypeSet,
				Set:         getHashFromDomain,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Description: "Domain name",
							Required:    true,
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
					}},
			},
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

	if v, ok := m["domain_name"]; ok {
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
		if attr, ok := domainItem["domain_name"]; ok && attr != "" {
			siteDomainDetails[domainInd].Domain = attr.(string)
		}
		domainInd++
	}
	return siteDomainDetails
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	siteDomainDetails := populateFromResourceToDTO(d)
	siteDomainDetailsDto, err := client.BulkUpdateDomainsToSite(siteID, siteDomainDetails)
	if err != nil {
		resourceDomainRead(d, m)
		return err
	}

	if siteDomainDetailsDto.Errors != nil && len(siteDomainDetailsDto.Errors) > 0 {
		if siteDomainDetailsDto.Errors[0].Status == "404" || siteDomainDetailsDto.Errors[0].Status == "400" {
			log.Printf("[INFO] Operation not allowed: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		if siteDomainDetailsDto.Errors[0].Status == "500" {
			log.Printf("[INFO] Internal Server Error: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(siteDomainDetailsDto.Errors)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("error updating domains for site (%s): %s", d.Get("site_id"), string(out))
	}

	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteDomainDetailsDto, err := client.GetWebsiteDomains(d.Get("site_id").(string))
	if err != nil {
		return err
	}

	if siteDomainDetailsDto.Errors != nil && len(siteDomainDetailsDto.Errors) > 0 {
		if siteDomainDetailsDto.Errors[0].Status == "404" || siteDomainDetailsDto.Errors[0].Status == "400" {
			log.Printf("[INFO] Operation not allowed: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		if siteDomainDetailsDto.Errors[0].Status == "500" {
			log.Printf("[INFO] Internal Server Error: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(siteDomainDetailsDto.Errors)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("error getting domains for site (%s): %s", d.Get("site_id"), string(out))
	}

	domains := &schema.Set{F: getHashFromDomain}
	for _, v := range siteDomainDetailsDto.Data {
		if v.MainDomain == true || v.AutoDiscovered == true { //we ignore the main and auto discovered domains
			continue
		}
		domain := map[string]interface{}{}
		domain["domain_name"] = v.Domain
		domain["id"] = v.Id
		domain["status"] = v.Status
		domains.Add(domain)
	}
	d.SetId(d.Get("site_id").(string))
	d.Set("site_id", d.Get("site_id"))
	d.Set("cname_redirection_record", siteDomainDetailsDto.Data[0].CnameRedirectionRecord)
	d.Set("domain", domains)
	return nil
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	siteDomainDetailsDto, err := client.BulkUpdateDomainsToSite(siteID, []SiteDomainDetails{})
	if err != nil {
		return err
	}

	if siteDomainDetailsDto.Errors != nil && len(siteDomainDetailsDto.Errors) > 0 {
		if siteDomainDetailsDto.Errors[0].Status == "404" || siteDomainDetailsDto.Errors[0].Status == "400" {
			log.Printf("[INFO] Operation not allowed: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		if siteDomainDetailsDto.Errors[0].Status == "500" {
			log.Printf("[INFO] Internal Server Error: %s\n", siteDomainDetailsDto.Errors)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(siteDomainDetailsDto.Errors)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("error deleting domains for site (%s): %s", d.Get("site_id"), string(out))
	}

	return resourceDomainRead(d, m)
}
