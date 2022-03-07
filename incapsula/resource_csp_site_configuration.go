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

func resourceCspSiteConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCspSiteConfigurationUpdate,
		Read:   resourceCspSiteConfigurationRead,
		Update: resourceCspSiteConfigurationUpdate,
		Delete: resourceCspSiteConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import  Site Config JSON for Site ID %d", siteID)
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
				Default:      "monitor",
				ValidateFunc: validation.StringInSlice([]string{"monitor", "enforce", "off"}, true),
			},
			"email_addresses": {
				Description: "Email address for the event notification recipient list of a specific website. Notifications are reasonably small and limited in frequency",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceCspSiteConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)

	log.Printf("[DEBUG] Reading CSP site configuration for site ID:  %d.", siteId)

	cspSite, err := client.GetCspSite(siteId)
	if err != nil {
		log.Printf("[ERROR] Could not get CSP site config: %s - %s\n", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Reading CSP site configuration for site ID: %d , response: %v.", siteId, cspSite)

	var emails []string
	for i := range cspSite.Settings.Emails {
		emails = append(emails, cspSite.Settings.Emails[i].Email)
	}
	d.Set("email_addresses", emails)

	switch {
	case strings.Compare(cspSite.Discovery, "pause") == 0:
		d.Set("mode", cspSiteModeOff) // Do I need to set ID to blank if off?
	case strings.Compare(cspSite.Discovery, "start") == 0 && strings.Compare(cspSite.Mode, "monitor") == 0:
		d.Set("mode", cspSiteModeMonitor)
	case strings.Compare(cspSite.Discovery, "start") == 0 && strings.Compare(cspSite.Mode, "enforce") == 0:
		d.Set("mode", cspSiteModeEnforce)
	default:
		d.Set("mode", cspSiteModeOff)
	}

	return nil
}

func resourceCspSiteConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	emails := d.Get("email_addresses").([]interface{})
	siteId := d.Get("site_id").(int)

	cspSiteConfig := CspSiteConfig{
		Name:      "",
		Mode:      "",
		Discovery: "",
		Settings: struct {
			Emails []CspSiteConfigEmail `json:"emails"`
		}{},
		TrackingIDs: nil,
	}

	cspSiteConfig.Settings.Emails = []CspSiteConfigEmail{}
	for i := range emails {
		a := emails[i]
		switch a.(type) {
		case string:
			cspSiteConfig.Settings.Emails = append(cspSiteConfig.Settings.Emails, CspSiteConfigEmail{Email: a.(string)})
		}
	}

	switch d.Get("mode").(string) {
	case cspSiteModeOff:
		cspSiteConfig.Discovery = "pause"
		cspSiteConfig.Mode = "monitor"
	case cspSiteModeMonitor:
		cspSiteConfig.Discovery = "start"
		cspSiteConfig.Mode = "monitor"
	case cspSiteModeEnforce:
		cspSiteConfig.Discovery = "start"
		cspSiteConfig.Mode = "enforce"
	}

	log.Printf("[DEBUG] Updating CSP site configuration for site ID: %d , values: %v.", siteId, cspSiteConfig)
	updatedSite, err := client.UpdateCspSite(siteId, &cspSiteConfig)
	if err != nil {
		log.Printf("[ERROR] Could not update CSP site config: %s - %s\n", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Updating CSP site configuration for site ID: %d , got response: %v.", siteId, updatedSite)

	return nil
}

func resourceCspSiteConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	// Deleting the CSP settings is just setting the site mode to off
	d.Set("mode", cspSiteModeOff)
	d.SetId("") // do I need to do this here?
	return resourceCspSiteConfigurationUpdate(d, m)
}
