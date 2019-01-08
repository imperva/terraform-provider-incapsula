package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

const incapAddSiteURL = "https://my.incapsula.com/api/prov/v1/sites/add"

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"api_id": &schema.Schema{
				Description: "API authentication identifier.",
				Type:        schema.TypeString,
				Required:    true,
				StateFunc: func(val interface{}) string {
					return "redacted"
				},
			},
			"api_key": &schema.Schema{
				Description: "API authentication identifier.",
				Type:        schema.TypeString,
				Required:    true,
				StateFunc: func(val interface{}) string {
					return "redacted"
				},
			},
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

			// Computed Attributes
			// NOTE: Site ID will also be used for the ID
			"site_id": &schema.Schema{
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
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
	type DNSRecord struct {
		DNSRecordName string   `json:"dns_record_name"`
		SetTypeTo     string   `json:"set_type_to"`
		SetDataTo     []string `json:"set_data_to"`
	}

	type AddSiteResponse struct {
		SiteID           int         `json:"site_id"`
		SiteCreationDate int         `json:"site_creation_date"`
		DNS              []DNSRecord `json:"dns"`
		Res              int         `json:"res"`
	}

	domain := d.Get("domain").(string)
	log.Printf("[INFO] Adding Incapsula site for domain: %s\n", domain)

	// Post form to Incapsula
	resp, err := http.PostForm(incapAddSiteURL, url.Values{
		"api_id":                 {d.Get("api_id").(string)},
		"api_key":                {d.Get("api_key").(string)},
		"domain":                 {domain},
		"account_id":             {d.Get("account_id").(string)},
		"ref_id":                 {d.Get("ref_id").(string)},
		"send_site_setup_emails": {d.Get("send_site_setup_emails").(string)},
		"site_ip":                {d.Get("site_ip").(string)},
		"force_ssl":              {d.Get("force_ssl").(string)},
		"log_level":              {d.Get("log_level").(string)},
		"logs_account_id":        {d.Get("logs_account_id").(string)},
	})
	if err != nil {
		return fmt.Errorf("Error adding site for domain %q: %s", domain, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var addSiteResponse AddSiteResponse
	err = json.Unmarshal([]byte(responseBody), &addSiteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing add site JSON response for domain %q: %s", domain, err)
	}

	// Look at the response status code from Incapsula
	if addSiteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when creating site for domain %q: %s", domain, string(responseBody))
	}

	// Set the Site ID
	d.SetId(strconv.Itoa(addSiteResponse.SiteID))
	d.Set("site_id", addSiteResponse.SiteID)

	// Set the Site Creation Date
	d.Set("site_creation_date", addSiteResponse.SiteCreationDate)

	// Set the DNS information
	for _, entry := range addSiteResponse.DNS {
		if entry.SetTypeTo == "CNAME" && len(entry.SetDataTo) > 0 {
			d.Set("dns_cname_record_name", entry.DNSRecordName)
			d.Set("dns_cname_record_value", entry.SetDataTo[0])
		} else {
			d.Set("dns_a_record_name", entry.DNSRecordName)
			d.Set("dns_a_record_value", entry.SetDataTo)
		}
	}

	return resourceSiteRead(d, m)
}

func resourceSiteRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
