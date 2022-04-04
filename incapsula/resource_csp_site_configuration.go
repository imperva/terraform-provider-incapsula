package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

const (
	cspSiteModeMonitor = "monitor"
	cspSiteModeEnforce = "enforce"
	cspSiteModeOff     = "off"
)

func resourceCSPSiteConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSPSiteConfigurationUpdate,
		Read:   resourceCSPSiteConfigurationRead,
		Update: resourceCSPSiteConfigurationUpdate,
		Delete: resourceCSPSiteConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import CSP Site Config JSON for Site ID %d", siteID)
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
			//Optional
			"mode": {
				Description:  "Website Protection Mode. When in \"enforce\" mode, blocked resources will not be available in the application and new resources will be automatically blocked. When in \"monitor\" mode, all resources are available in the application and the system keeps track of all new domains that are discovered.\nValues: monitor\\enforce\\off\n",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      cspSiteModeMonitor,
				ValidateFunc: validation.StringInSlice([]string{cspSiteModeMonitor, cspSiteModeEnforce, cspSiteModeOff}, false),
			},
			"email_addresses": {
				Description: "Email address for the event notification recipient list of a specific website. Notifications are reasonably small and limited in frequency",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceCSPSiteConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)

	log.Printf("[DEBUG] Reading CSP site configuration for site ID:  %d.", siteId)

	cspSite, err := client.GetCSPSite(siteId)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP site config: %s - %s\n", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Reading CSP site configuration for site ID: %d , response: %v.", siteId, cspSite)

	emails := &schema.Set{F: schema.HashString}
	for i := range cspSite.Settings.Emails {
		emails.Add(cspSite.Settings.Emails[i].Email)
	}
	d.Set("email_addresses", emails)

	switch {
	case strings.Compare(cspSite.Discovery, CSPDiscoveryOff) == 0:
		d.Set("mode", cspSiteModeOff)
	case strings.Compare(cspSite.Discovery, CSPDiscoveryOn) == 0 && strings.Compare(cspSite.Mode, cspSiteModeMonitor) == 0:
		d.Set("mode", cspSiteModeMonitor)
	case strings.Compare(cspSite.Discovery, CSPDiscoveryOn) == 0 && strings.Compare(cspSite.Mode, cspSiteModeEnforce) == 0:
		d.Set("mode", cspSiteModeEnforce)
	default:
		d.Set("mode", cspSiteModeOff)
	}

	return nil
}

func resourceCSPSiteConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	emails := d.Get("email_addresses").(*schema.Set)
	siteId := d.Get("site_id").(int)

	cspSiteConfig := CSPSiteConfig{
		Name:      "",
		Mode:      "",
		Discovery: "",
		Settings: struct {
			Emails []CSPSiteConfigEmail `json:"emails"`
		}{},
		TrackingIDs: nil,
	}

	cspSiteConfig.Settings.Emails = []CSPSiteConfigEmail{}
	for _, email := range emails.List() {
		cspSiteConfig.Settings.Emails = append(cspSiteConfig.Settings.Emails, CSPSiteConfigEmail{Email: email.(string)})
	}

	switch d.Get("mode").(string) {
	case cspSiteModeOff:
		cspSiteConfig.Discovery = CSPDiscoveryOff
		cspSiteConfig.Mode = cspSiteModeMonitor
	case cspSiteModeMonitor:
		cspSiteConfig.Discovery = CSPDiscoveryOn
		cspSiteConfig.Mode = cspSiteModeMonitor
	case cspSiteModeEnforce:
		cspSiteConfig.Discovery = CSPDiscoveryOn
		cspSiteConfig.Mode = cspSiteModeEnforce
	}

	log.Printf("[DEBUG] Updating CSP site configuration for site ID: %d , values: %v.", siteId, cspSiteConfig)
	updatedSite, err := client.UpdateCSPSite(siteId, &cspSiteConfig)
	if err != nil {
		log.Printf("[ERROR] Could not update CSP site config: %s - %s\n", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Updating CSP site configuration for site ID: %d , got response: %v.", siteId, updatedSite)
	id := strconv.Itoa(siteId)
	d.SetId(id)

	return resourceCSPSiteConfigurationRead(d, m)
}

func resourceCSPSiteConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	// Deleting the CSP settings is just setting the site mode to off
	d.Set("mode", cspSiteModeOff)
	if err := resourceCSPSiteConfigurationUpdate(d, m); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
