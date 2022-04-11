package incapsula

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

const (
	cspDomainStatusAllowed = "allowed"
	cspDomainStatusBlocked = "blocked"
)

func resourceCSPSiteDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSPSiteDomainUpdate,
		Read:   resourceCSPSiteDomainRead,
		Update: resourceCSPSiteDomainUpdate,
		Delete: resourceCSPSiteDomainDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				keyParts := strings.Split(d.Id(), "/")
				if len(keyParts) != 3 {
					return nil, fmt.Errorf("Error parsing ID, actual value: %s, expected two numeric IDs and string seperated by '/'\n", d.Id())
				}
				accountID, err := strconv.Atoi(keyParts[0])
				if err != nil {
					return nil, fmt.Errorf("failed to convert account ID from import command, actual value: %s, expected numeric id", keyParts[0])
				}
				siteID, err := strconv.Atoi(keyParts[1])
				if err != nil {
					return nil, fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric id", keyParts[1])
				}
				domain, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(keyParts[2])
				if err != nil {
					return nil, fmt.Errorf("failed to convert domain reference ID from import command, actual value: %s, expected Base64 id", keyParts[2])
				}

				d.Set("account_id", accountID)
				d.Set("site_id", siteID)
				d.Set("domain", string(domain))
				log.Printf("[DEBUG] Import CSP Domain %s for site ID %d", domain, siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
			},
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
			//Optional
			"include_subdomains": {
				Description: "Defines Whether or not subdomains will inherit the allowance of the parent domain. Values: true, false",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"status": {
				Description:  "Defines whether the domain should be Blocked or Allowed once the site's mode changes to the Enforcement. Values: Blocked, Allowed",
				Type:         schema.TypeString,
				Default:      cspDomainStatusAllowed,
				ValidateFunc: validation.StringInSlice([]string{cspDomainStatusAllowed, cspDomainStatusBlocked}, false),
				Optional:     true,
			},
			"notes": {
				Description: "Add a quick note to a domain to help in future analysis and investigation. You can add as many notes as you like.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceCSPSiteDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))

	log.Printf("[DEBUG] Reading CSP domain for site ID: %d , domain reference: %s , domain: %s", siteID, domainRef, domain)

	cspNotes, err := client.getCSPDomainNotes(accountID, siteID, domain)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP domain notes: %s - %s\n", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Reading CSP domain notes for domain %s from site ID: %d , response: %v.", domain, siteID, cspNotes)

		notes := &schema.Set{F: schema.HashString}
		for i := range cspNotes {
			notes.Add(cspNotes[i].Text)
		}
		log.Printf("[DEBUG] Reading CSP domain notes for domain %s from site ID: %d , updating notes to: %v.", domain, siteID, notes)

		d.Set("notes", notes)
	}

	// First check if it's a pre-approved domain, and update resource according to that
	preApprovedDomain, err := client.getCSPPreApprovedDomain(accountID, siteID, domain)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP pre-approved domain : %s - %s\n", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Reading CSP pre-approved domain %s for site ID: %d , response: %v.", domain, siteID, preApprovedDomain)

		d.Set("include_subdomains", preApprovedDomain.Subdomains)
		d.Set("status", cspDomainStatusAllowed)

		return nil
	}

	// If domain wasn't found as pre-approved domain, check if status set directly and update accordingly
	status, err := client.getCSPDomainStatus(accountID, siteID, domain)
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

func resourceCSPSiteDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	domRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	status := d.Get("status").(string)
	notes := d.Get("notes").(*schema.Set)

	log.Printf("[DEBUG] Updating CSP domain %s site ID %d status=\"%s\"\n", domain, siteID, status)

	if strings.Compare(status, cspDomainStatusAllowed) == 0 {
		// If the domain is allowed just put it in the pre-approved list
		dom := CSPPreApprovedDomain{
			Domain:      domain,
			Subdomains:  d.Get("include_subdomains").(bool),
			ReferenceID: base64.RawURLEncoding.EncodeToString([]byte(domain)),
		}
		log.Printf("[DEBUG] Updating CSP domain for site ID: %d , domain: %v\n", siteID, dom)
		updatedDom, err := client.updateCSPPreApprovedDomain(accountID, siteID, &dom)
		if err != nil {
			log.Printf("[ERROR] Could not update CSP pre-approved domain: %v - %s\n", dom, err)
			return err
		}
		log.Printf("[DEBUG] Updating CSP domain %v for site ID: %d , got response: %v.", dom, siteID, updatedDom)
	} else if strings.Compare(status, cspDomainStatusBlocked) == 0 {
		// Otherwise update the status directly to blocked
		st := CSPDomainStatus{
			Blocked:  new(bool),
			Reviewed: new(bool),
		}
		*(st.Blocked) = true
		*(st.Reviewed) = true

		domainStatus, err := client.updateCSPDomainStatus(accountID, siteID, domain, &st)
		if err != nil || domainStatus.Blocked == nil || domainStatus.Reviewed == nil {
			e := fmt.Errorf("[ERROR] Could not update CSP domain %s status: %v - %s\n", domain, status, err)
			return e
		}
	}

	// Remove all existing notes and add them freshly
	client.deleteCSPDomainNotes(accountID, siteID, domain)
	for _, note := range notes.List() {
		client.addCSPDomainNote(accountID, siteID, domain, note.(string))
	}

	newID := fmt.Sprintf("%d/%d/%s", accountID, siteID, domRef)
	log.Printf("[DEBUG] Update CSP Domain, setting key %s to: %s", d.Id(), newID)
	d.SetId(newID)

	return resourceCSPSiteDomainRead(d, m)
}

func resourceCSPSiteDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	status := d.Get("status").(string)
	log.Printf("[DEBUG] Deleting CSP domain %s from site ID %d\n", domain, siteID)

	if strings.Compare(status, cspDomainStatusAllowed) == 0 {
		err := client.deleteCSPPreApprovedDomains(accountID, siteID, base64.RawURLEncoding.EncodeToString([]byte(domain)))
		if err != nil {
			log.Printf("[ERROR] Could not delete CSP pre-approved domain %s for site ID %d: %s\n", domain, siteID, err)
			return err
		}
	} else if strings.Compare(status, cspDomainStatusBlocked) == 0 {
		newStatus := CSPDomainStatus{
			Blocked:  new(bool),
			Reviewed: new(bool),
		}
		*newStatus.Blocked = false
		*newStatus.Reviewed = false
		ret, err := client.updateCSPDomainStatus(accountID, siteID, domain, &newStatus)
		if err != nil {
			log.Printf("[ERROR] Could not delete CSP domain status %s for site ID %d: %s\n", domain, siteID, err)
			return err
		}
		if ret.Blocked == nil || ret.Reviewed == nil {
			return fmt.Errorf("[ERROR] Could not update CSP domain %s status to: %v got: %v\n", domain, newStatus, ret)
		}
	}
	log.Printf("[DEBUG] Deleted CSP domain %s for site ID: %d successfully", domain, siteID)

	return nil
}
