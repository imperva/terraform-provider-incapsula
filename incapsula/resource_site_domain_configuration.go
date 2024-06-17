package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"slices"
	"strconv"
	"time"
)

func resourceSiteDomainConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDomainRead,
		Create: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Update: resourceDomainUpdate,
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
			"managed_certificate_id": {
				Description: "Numeric identifier of the site certificate id.",
				Type:        schema.TypeString,
				Optional:    true,
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
						"name": {
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

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	siteDomainDetails := populateFromResourceToDTO(d)
	err := client.BulkUpdateDomainsToSite(siteID, siteDomainDetails)
	if err != nil {
		return err
	}
	siteDomainDetailsDto, err := client.GetWebsiteDomains(d.Get("site_id").(string))
	id, _ := strconv.Atoi(siteId)
	if d.Get("managed_certificate_id") != "" {
		for i := 0; i < 20; i++ {
			siteCertificateV3Response, _ := client.GetSiteCertificateRequestStatus(id)
			b := siteCertificateV3Response != nil && siteCertificateV3Response.Data != nil && len(siteCertificateV3Response.Data) > 0
			b = b && siteCertificateV3Response.Data[0].CertificatesDetails != nil && validateSans(siteCertificateV3Response.Data[0], siteDomainDetailsDto.Data)
			if b {
				return resourceDomainRead(d, m)
			}
			time.Sleep(10 * time.Second)
		}
	} else {
		return resourceDomainRead(d, m)
	}
	panic("managed certificate was not created as expected")
}

func validateSans(siteCertificateDTO SiteCertificateDTO, siteDomainDetails []SiteDomainDetails) bool {
	if len(siteDomainDetails) == 0 { //todo handle multiple certificates
		return len(siteCertificateDTO.CertificatesDetails) == 0 || validateEmptyCertificates(siteCertificateDTO.CertificatesDetails)
	} else {
		return len(siteCertificateDTO.CertificatesDetails) > 0 && validateDomainsSans(siteDomainDetails, siteCertificateDTO.CertificatesDetails)
	}
}

func validateEmptyCertificates(certificatesDetails []CertificateDTO) bool {
	for _, t := range certificatesDetails {
		if t.Sans != nil && len(t.Sans) > 0 && !allSansAreDeleted(t.Sans) {
			return false
		}
	}
	return true
}

func allSansAreDeleted(sans []SanDTO) bool {
	deleteStatuses := []string{"REMOVED", "DELETED_PENDING_PUBLICATION", "DELETED_LOCALLY", "CANCELED_LOCALLY", "PENDING_CANCELATION", "PENDING_DELETION"}
	for _, t := range sans {
		if !slices.Contains(deleteStatuses, t.Status) {
			return false
		}
	}
	return true
}

func validateDomainsSans(siteDomainDetails []SiteDomainDetails, certificateDTO []CertificateDTO) bool {
	return validateSanForEachDomain(siteDomainDetails, certificateDTO) && validateRedundantSansAreDeleted(siteDomainDetails, certificateDTO)
}

func validateSanForEachDomain(siteDomainDetails []SiteDomainDetails, certificateDTO []CertificateDTO) bool {
	for _, t := range siteDomainDetails {
		if !validateSanForDomain(t, certificateDTO) {
			return false
		}
	}
	return true
}

func validateSanForDomain(domain SiteDomainDetails, certificates []CertificateDTO) bool {
	for _, t := range certificates {
		if t.Sans != nil && len(t.Sans) > 0 {
			for _, s := range t.Sans {
				if slices.Contains(s.DomainIds, domain.Id) && s.SanValue == domain.Domain {
					return true
				}
			}
		}
	}
	return false
}

func validateRedundantSansAreDeleted(domains []SiteDomainDetails, certificates []CertificateDTO) bool {
	res := true
	for _, t := range certificates {
		if t.Sans != nil && len(t.Sans) > 0 {
			for _, s := range t.Sans {
				res = res && validateExistDomainOrDeletedSan(s, domains)
			}
		}
	}
	return res
}

func validateExistDomainOrDeletedSan(s SanDTO, domains []SiteDomainDetails) bool {
	deleteStatuses := []string{"REMOVED", "DELETED_PENDING_PUBLICATION", "DELETED_LOCALLY", "CANCELED_LOCALLY", "PENDING_CANCELATION", "PENDING_DELETION"}
	for _, t := range domains {
		if t.Domain == s.SanValue && slices.Contains(s.DomainIds, t.Id) {
			return true
		}
	}
	return !slices.Contains(deleteStatuses, s.Status)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
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
		if v.AutoDiscovered == true || (d.Get("managed_certificate_id") == "" && v.MainDomain == true) { //we ignore the main and auto discovered domains
			continue
		}
		domain := map[string]interface{}{}
		domain["name"] = v.Domain
		domain["id"] = v.Id
		domain["status"] = v.Status
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
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	err := client.BulkUpdateDomainsToSite(siteID, []SiteDomainDetails{})
	if err != nil {
		return err
	}

	return resourceDomainRead(d, m)
}
