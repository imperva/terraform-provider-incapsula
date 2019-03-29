package incapsula

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"domain": &schema.Schema{
				Description: "The domain name of the site. For example: www.example.com, hello.example.com, example.com.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"account_id": &schema.Schema{
				Description: "Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ref_id": &schema.Schema{
				Description: "Customer specific identifier for this operation.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"send_site_setup_emails": &schema.Schema{
				Description: "If this value is false, end users will not get emails about the add site process such as DNS instructions and SSL setup.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"site_ip": &schema.Schema{
				Description: "Manually set the web server IP/CNAME.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"force_ssl": &schema.Schema{
				Description: "If this value is true, manually set the site to support SSL. This option is only available for sites with manually configured IP/CNAME and for specific accounts.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"log_level": &schema.Schema{
				Description: "Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Sets the log reporting level for the site. Options are full, security, none and default.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"logs_account_id": &schema.Schema{
				Description: "Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"site_id": &schema.Schema{
				Description: "Numeric identifier of the site.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"active": &schema.Schema{
				Description: "active or bypass.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_validation": &schema.Schema{
				Description: "email or html or dns.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"approver": &schema.Schema{
				Description: "my.approver@email.com (some approver email address).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ignore_ssl": &schema.Schema{
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"acceleration_level": &schema.Schema{
				Description: "none | standard | aggressive.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"seal_location": &schema.Schema{
				Description: "api.seal_location.bottom_left | api.seal_location.none | api.seal_location.right_bottom | api.seal_location.right | api.seal_location.left | api.seal_location.bottom_right | api.seal_location.bottom.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_redirect_to_full": &schema.Schema{
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"remove_ssl": &schema.Schema{
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Computed Attributes
			"site_creation_date": &schema.Schema{
				Description: "Numeric representation of the site creation date.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"dns_cname_record_name": &schema.Schema{
				Description: "CNAME record name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_cname_record_value": &schema.Schema{
				Description: "CNAME record value.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_a_record_name": &schema.Schema{
				Description: "A record name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_a_record_value": {
				Description: "A record value.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSiteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	domain := d.Get("domain").(string)

	siteAddResponse, err := client.AddSite(
		domain,
		d.Get("account_id").(string),
		d.Get("ref_id").(string),
		d.Get("send_site_setup_emails").(string),
		d.Get("site_ip").(string),
		d.Get("force_ssl").(string),
		d.Get("log_level").(string),
		d.Get("logs_account_id").(string),
	)

	if err != nil {
		return err
	}

	// Set the Site ID
	d.SetId(strconv.Itoa(siteAddResponse.SiteID))

	// Set the rest of the state from the resource read
	return resourceSiteRead(d, m)
}

func resourceSiteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	domain := d.Get("domain").(string)
	siteID, _ := strconv.Atoi(d.Id())

	siteStatusResponse, err := client.SiteStatus(domain, siteID)

	if err != nil {
		return err
	}

	d.Set("site_creation_date", siteStatusResponse.SiteCreationDate)
	d.Set("domain", siteStatusResponse.Domain)

	// Set the DNS information
	dnsARecordValues := make([]string, 0)
	for _, entry := range siteStatusResponse.DNS {
		if entry.SetTypeTo == "CNAME" && len(entry.SetDataTo) > 0 {
			d.Set("dns_cname_record_name", entry.DNSRecordName)
			d.Set("dns_cname_record_value", entry.SetDataTo[0])
		}
		if entry.SetTypeTo == "A" {
			d.Set("dns_a_record_name", entry.DNSRecordName)
			dnsARecordValues = append(dnsARecordValues, entry.SetDataTo...)
		}
	}
	d.Set("dns_a_record_value", dnsARecordValues)

	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID := d.Get("site_id").(int)

	if active := d.Get("active").(string); active != "" {
		_, err := client.UpdateSiteActive(
			siteID,
			active,
		)
		if err != nil {
			return err
		}
	}

	if siteIP := d.Get("site_ip").(string); siteIP != "" {
		_, err := client.UpdateSiteSiteIP(
			siteID,
			siteIP,
		)
		if err != nil {
			return err
		}
	}

	if domainValidation := d.Get("domain_validation").(string); domainValidation != "" {
		_, err := client.UpdateSiteDomainValidation(
			siteID,
			domainValidation,
		)
		if err != nil {
			return err
		}
	}

	if approver := d.Get("approver").(string); approver != "" {
		_, err := client.UpdateSiteApprover(
			siteID,
			approver,
		)
		if err != nil {
			return err
		}
	}

	if ignoreSSL := d.Get("ignore_ssl").(string); ignoreSSL != "" {
		_, err := client.UpdateSiteIgnoreSSL(
			siteID,
			ignoreSSL,
		)
		if err != nil {
			return err
		}
	}

	if accelerationLevel := d.Get("acceleration_level").(string); accelerationLevel != "" {
		_, err := client.UpdateSiteAccelerationLevel(
			siteID,
			accelerationLevel,
		)
		if err != nil {
			return err
		}
	}

	if sealLocation := d.Get("seal_location").(string); sealLocation != "" {
		_, err := client.UpdateSiteSealLocation(
			siteID,
			sealLocation,
		)
		if err != nil {
			return err
		}
	}

	if domainRedirectToFull := d.Get("domain_redirect_to_full").(string); domainRedirectToFull != "" {
		_, err := client.UpdateSiteDomainRedirectToFull(
			siteID,
			domainRedirectToFull,
		)
		if err != nil {
			return err
		}
	}

	if removeSSL := d.Get("remove_ssl").(string); removeSSL != "" {
		_, err := client.UpdateSiteRemoveSSL(
			siteID,
			removeSSL,
		)
		if err != nil {
			return err
		}
	}

	if refID := d.Get("ref_id").(string); refID != "" {
		_, err := client.UpdateSiteRefID(
			siteID,
			refID,
		)
		if err != nil {
			return err
		}
	}

	// Set the rest of the state from the resource read
	return resourceSiteRead(d, m)
}

func resourceSiteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	domain := d.Get("domain").(string)
	siteID, _ := strconv.Atoi(d.Id())

	err := client.DeleteSite(domain, siteID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
