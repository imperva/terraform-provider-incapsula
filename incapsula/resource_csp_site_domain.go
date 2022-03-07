package incapsula

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

const (
	cspDomainStatusAllowed = "Allowed"
	cspDomainStatusBlocked = "Blocked"
)

func resourceCspSiteDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceCspSiteDomainCreate,
		Read:   resourceCspSiteDomainRead,
		Update: resourceCspSiteDomainUpdate,
		Delete: resourceCspSiteDomainDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				keyParts := strings.Split(d.Id(), ".")
				if len(keyParts) != 2 {
					return nil, fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '.'\n", d.Id())
				}
				siteID, err := strconv.Atoi(keyParts[0])
				if err != nil {
					return nil, fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", keyParts[0])
				}
				domain, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(keyParts[1])
				if err != nil {
					return nil, fmt.Errorf("failed to convert domain reference ID from import command, actual value: %s, expected Base64 id", keyParts[1])
				}

				d.Set("site_id", siteID)
				d.Set("domain", string(domain))
				log.Printf("[DEBUG] Import CSP Domain %s for site ID %d", domain, siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"domain": {
				Description: "The fully qualified domain name of the site. For example: www.example.com, hello.example.com.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"status": {
				Description: "Defines whether the domain should be Blocked or Allowed once the site's mode changes to the Enforcement. Values: Blocked, Allowed",
				Type:        schema.TypeString,
				Required:    true,
			},
			"include_subdomains": {
				Description: "Defines Whether or not subdomains will inherit the allowance of the parent domain. Values: true, false",
				Type:        schema.TypeBool,
				Required:    true,
			},
			//Optional
			"notes": {
				Description: "Add a quick note to a domain to help in future analysis and investigation. You can add as many notes as you like.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceCspSiteDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	log.Printf("[DEBUG] Reading CSP domain for site ID: %d , domain reference: %s , domain: %s", siteId, domainRef, domain)

	/*
		domainData, err := client.getCspDomainData(siteId, domainRef)
		if err != nil {
			log.Printf("[ERROR] Could not get CSP domain data: %s - %s\n", d.Id(), err)
		}
		log.Printf("[DEBUG] Reading CSP domain %s (ref %s) configuration for site ID: %d , response: %v.", domain, domainRef, siteId, domainData)
	*/

	preApprovedDomains, err := client.getCspPreApprovedDomains(siteId)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP pre-approved domains list: %s - %s\n", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Reading CSP pre-approved domains list for site ID: %d , response: %v.", siteId, preApprovedDomains)

	dom, ok := preApprovedDomains[domainRef]
	if !ok {
		d.SetId("")
		fmt.Errorf("Error reading any CSP domain data for domain %s from site ID %d\n",
			domain, siteId)
		return nil
	}

	log.Printf("[DEBUG] Reading CSP domain, found matching pre-approved: %v", dom)
	d.Set("domain", dom.Domain)
	d.Set("include_subdomains", dom.Subdomains)
	d.Set("status", cspDomainStatusAllowed)

	/*
		if domainData != nil {
			log.Printf("[DEBUG] Reading CSP domain, found domain data, using for allowance status: %v", domainData)
			if domainData.Status.Blocked == true {
				d.Set("status", cspDomainStatusBlocked)
			} else {
				d.Set("status", cspDomainStatusAllowed)
			}
		}

		d.Set("notes", []string{})
	*/

	return nil
}

func resourceCspSiteDomainCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceCspSiteDomainUpdate(d, m)
	if err != nil {
		return err
	}
	domRef := base64.RawURLEncoding.EncodeToString([]byte(d.Get("domain").(string)))

	newID := fmt.Sprintf("%d.%s", d.Get("site_id").(int), domRef)
	log.Printf("[DEBUG] Create CSP Domain, changing key %s to: %s", d.Id(), newID)
	d.SetId(newID)

	return nil
}

func resourceCspSiteDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	st := d.Get("status").(string)
	log.Printf("[DEBUG] Updating CSP domain, read site ID %d and status=\"%s\"\n", siteID, st)
	if strings.Compare(st, cspDomainStatusAllowed) != 0 {
		log.Printf("[DEBUG] Updating CSP domain, skipping Blocked domains for now. %s\n", d.Get("domain").(string))
		return nil
	}

	dom := CspPreApprovedDomain{
		Domain:      d.Get("domain").(string),
		Subdomains:  d.Get("include_subdomains").(bool),
		ReferenceID: base64.RawURLEncoding.EncodeToString([]byte(d.Get("domain").(string))),
	}
	log.Printf("[DEBUG] Updating CSP domain for site ID: %d , domain: %v\n", siteID, dom)
	updatedDom, err := client.updateCspPreApprovedDomain(siteID, &dom)
	if err != nil {
		log.Printf("[ERROR] Could not update CSP pre-approved domain: %v - %s\n", dom, err)
		return err
	}
	log.Printf("[DEBUG] Updating CSP domain %v for site ID: %d , got response: %v.", dom, siteID, updatedDom)

	d.SetId(fmt.Sprintf("%d.%s", siteID, dom.ReferenceID))

	return nil
}

func resourceCspSiteDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	log.Printf("[DEBUG] Deleting CSP domain %s from site ID %d\n", domain, siteID)

	err := client.deleteCspPreApprovedDomains(siteID, base64.RawURLEncoding.EncodeToString([]byte(domain)))
	if err != nil {
		log.Printf("[ERROR] Could not delete CSP pre-approved domain %s for site ID %d: %s\n", domain, siteID, err)
		return err
	}
	log.Printf("[DEBUG] Deleted CSP domain %s for site ID: %d successfully", domain, siteID)

	return nil
}
