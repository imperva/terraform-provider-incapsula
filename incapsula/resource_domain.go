package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSiteSingleDomainConfiguration() *schema.Resource {

	return &schema.Resource{
		Read:   resourceSingleDomainRead,
		Create: resourceDomainCreate,
		Delete: resourceSingleDomainDelete,
		Update: resourceDomainUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

				idSlice := strings.Split(d.Id(), "/")
				log.Printf("[DEBUG] Starting to import domain. Parameters: %s\n", d.Id())

				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id or site_id/domain_id", d.Id())
				}

				err := d.Set("site_id", idSlice[0])

				if err != nil {
					return nil, err
				}

				_, err = strconv.Atoi(idSlice[1])
				if err != nil {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id or site_id/domain_id", d.Id())
				}

				d.SetId(idSlice[1])

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
				Type:        schema.TypeString,
				Description: "Domain name",
				Required:    true,
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
			},
		},
	}
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {

	return fmt.Errorf("resource updates are not allowed")
}
func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	siteID := d.Get("site_id").(string)

	resp, err := client.AddDomainToSite(siteID, d.Get("domain").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(resp.Id))
	return resourceSingleDomainRead(d, m)
}

func resourceSingleDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	id := d.Id()
	siteDomainDetailsDto, err := client.GetDomain(siteID, id)

	if err != nil {
		return err
	}

	if siteDomainDetailsDto.Errors != nil && len(siteDomainDetailsDto.Errors) > 0 {
		out, err := json.Marshal(siteDomainDetailsDto.Errors)
		if err != nil {
			return err
		}

		return fmt.Errorf("error getting domains (%d) for site (%s): %s",
			d.Get("Id").(int),
			d.Get("site_id"),
			out)
	}

	d.Set("domain", siteDomainDetailsDto.Domain)
	d.Set("site_id", strconv.Itoa(siteDomainDetailsDto.SiteId))
	d.Set("status", siteDomainDetailsDto.Status)
	d.SetId(d.Id())
	d.Set("cname_redirection_record", siteDomainDetailsDto.CnameRedirectionRecord)

	if siteDomainDetailsDto.ARecords != nil && len(siteDomainDetailsDto.ARecords) > 0 {
		d.Set("a_records", siteDomainDetailsDto.ARecords)
	} else {
		d.Set("a_records", []string{})
	}

	return nil
}

func resourceSingleDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(string)
	err := client.DeleteDomain(siteID, d.Id())

	if err != nil {
		return err
	}
	return nil
}
