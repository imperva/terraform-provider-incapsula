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
				keyParts := strings.Split(d.Id(), "/")
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
			"include_subdomains": {
				Description: "Defines Whether or not subdomains will inherit the allowance of the parent domain. Values: true, false",
				Type:        schema.TypeBool,
				Required:    true,
			},
			//Optional
			"status": {
				Description: "Defines whether the domain should be Blocked or Allowed once the site's mode changes to the Enforcement. Values: Blocked, Allowed",
				Type:        schema.TypeString,
				Default:     cspDomainStatusAllowed,
				Optional:    true,
			},
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
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	log.Printf("[DEBUG] Reading CSP domain for site ID: %d , domain reference: %s , domain: %s", siteID, domainRef, domain)

	cspNotes, err := client.getCspDomainNotes(siteID, domain)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP domain notes: %s - %s\n", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Reading CSP domain notes for domain %s from site ID: %d , response: %v.", domain, siteID, cspNotes)

		notes := make([]string, len(cspNotes))
		for i := range cspNotes {
			notes = append(notes, cspNotes[i].Text)
		}
		d.Set("notes", notes)
	}

	// First check if it's a pre-approved domain, and update resource according to that
	preApprovedDomain, err := client.getCspPreApprovedDomain(siteID, domain)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP pre-approved domain : %s - %s\n", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Reading CSP pre-approved domain %s for site ID: %d , response: %v.", domain, siteID, preApprovedDomain)

		d.Set("include_subdomains", preApprovedDomain.Subdomains)
		d.Set("status", cspDomainStatusAllowed)

		return nil
	}

	// If domain wasn't found as pre-approved domain, check if status set directly and update accordingly
	status, err := client.getCspDomainStatus(siteID, domain)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP domain status: %s - %s\n", d.Id(), err)
	} else if status.Blocked != nil {
		log.Printf("[DEBUG] Reading CSP domain status for domain %s from site ID: %d , response: %v.", domain, siteID, status)
		d.Set("include_subdomains", strings.HasPrefix(domain, "*."))
		if !*(status.Blocked) {
			d.Set("status", cspDomainStatusAllowed)
		} else {
			d.Set("status", cspDomainStatusBlocked)
		}

		return nil
	}

	// In case we couldn't find data of pre-approved/status for the domain, remove it as a resource
	d.SetId("")
	fmt.Errorf("Error no CSP domain data found for domain %s from site ID %d\n",
		domain, siteID)
	return nil
}

func resourceCspSiteDomainCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceCspSiteDomainUpdate(d, m)
	if err != nil {
		return err
	}
	domRef := base64.RawURLEncoding.EncodeToString([]byte(d.Get("domain").(string)))

	newID := fmt.Sprintf("%d/%s", d.Get("site_id").(int), domRef)
	log.Printf("[DEBUG] Create CSP Domain, changing key %s to: %s", d.Id(), newID)
	d.SetId(newID)

	return nil
}

func resourceCspSiteDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	status := d.Get("status").(string)
	notes := d.Get("notes").([]interface{})

	log.Printf("[DEBUG] Updating CSP domain %s site ID %d status=\"%s\"\n", domain, siteID, status)

	if strings.Compare(status, cspDomainStatusAllowed) == 0 {
		// If the domain is allowed just put it in the pre-approved list
		dom := CspPreApprovedDomain{
			Domain:      domain,
			Subdomains:  d.Get("include_subdomains").(bool),
			ReferenceID: base64.RawURLEncoding.EncodeToString([]byte(domain)),
		}
		log.Printf("[DEBUG] Updating CSP domain for site ID: %d , domain: %v\n", siteID, dom)
		updatedDom, err := client.updateCspPreApprovedDomain(siteID, &dom)
		if err != nil {
			log.Printf("[ERROR] Could not update CSP pre-approved domain: %v - %s\n", dom, err)
			return err
		}
		log.Printf("[DEBUG] Updating CSP domain %v for site ID: %d , got response: %v.", dom, siteID, updatedDom)
	} else if strings.Compare(status, cspDomainStatusBlocked) == 0 {
		// Otherwise update the status directly to blocked
		st := CspDomainStatus{
			Blocked:  new(bool),
			Reviewed: new(bool),
		}
		*(st.Blocked) = true
		*(st.Reviewed) = true

		domainStatus, err := client.updateCspDomainStatus(siteID, domain, &st)
		if err != nil || domainStatus.Blocked == nil || domainStatus.Reviewed == nil {
			e := fmt.Errorf("[ERROR] Could not update CSP domain %s status: %v - %s\n", domain, status, err)
			return e
		}
	}

	// Remove all existing notes and add them freshly
	client.deleteCspDomainNotes(siteID, domain)
	for i := range notes {
		client.addCspDomainNote(siteID, domain, notes[i].(string))
	}

	return resourceCspSiteDomainRead(d, m)
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
