package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

var hstsConfigResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
			Optional: true,
		},
		"max_age": {
			Type:     schema.TypeInt,
			Default:  31536000,
			Optional: true,
		},
		"sub_domains_included": {
			Type:     schema.TypeBool,
			Computed: true,
			Optional: true,
		},
		"pre_loaded": {
			Type:     schema.TypeBool,
			Computed: true,
			Optional: true,
		},
	},
}

func resourceSiteSSLSettings() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSiteSSLSettingsRead,
		Update: resourceSiteSSLSettingsUpdate,
		Create: resourceSiteSSLSettingsUpdate,
		Delete: resourceSiteSSLSettingsDelete,
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
			// Add all types of configurations here that are related to TSL configuration endpoint
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"hsts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &hstsConfigResource,
				Set:      schema.HashResource(&hstsConfigResource),
			},
		},
	}
}

func resourceSiteSSLSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	setting := getSSLSettingsDTO(d)

	_, err := client.UpdateSiteSSLSettings(d.Get("site_id").(int), setting)

	if err != nil {
		return err
	}

	return resourceSiteSSLSettingsRead(d, m)
}

func resourceSiteSSLSettingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	settingsData, statusCode, err := client.ReadSiteSSLSettings(d.Get("site_id").(int))
	if statusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	if len(settingsData.Data) == 0 {
		return nil
	}

	d.SetId(fmt.Sprintf("site_ssl_settings_%d", d.Get("site_id").(int)))

	mapHSTSDTOtoHSTSResource(d, settingsData)
	// map other settings here

	return nil
}

func resourceSiteSSLSettingsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	// currently only disables HSTS
	// If more settings are implemented in the endpoint, add delete logic for them here.
	setting := prepareDisableHSTSStructure()
	var _, err = client.UpdateSiteSSLSettings(d.Get("site_id").(int), setting)

	if err != nil {
		return err
	}

	return nil
}

func prepareDisableHSTSStructure() SSLSettingsDTO {
	disableHSTSSetting := HSTSConfiguration{
		IsEnabled: false,
	}

	return SSLSettingsDTO{
		[]Data{
			{
				HstsConfiguration: disableHSTSSetting,
			},
		},
	}
}

func mapHSTSDTOtoHSTSResource(d *schema.ResourceData, settingsData *SSLSettingsDTO) {
	// handle HSTS remote configuration mapping
	var hstsSettingsFromServer HSTSConfiguration
	hstsSettingsFromServer = settingsData.Data[0].HstsConfiguration
	// Get the "hsts" attribute from the resource data
	// Create a map to hold the values for the "hsts" nested object
	hstsMap := make(map[string]interface{})

	// Set the values for each property of the "hsts" object
	hstsMap["is_enabled"] = hstsSettingsFromServer.IsEnabled
	hstsMap["max_age"] = hstsSettingsFromServer.MaxAge
	hstsMap["sub_domains_included"] = hstsSettingsFromServer.SubDomainsIncluded
	hstsMap["pre_loaded"] = hstsSettingsFromServer.PreLoaded

	// Set the "hsts" object in the resource data
	d.Set("hsts", []map[string]interface{}{hstsMap})
	// END HSTS mapping
}

func mapHSTSResourceToHSTSDTO(d *schema.ResourceData) HSTSConfiguration {
	hsts := d.Get("hsts").(*schema.Set)
	hstsList := hsts.List()
	hstsMap := hstsList[0].(map[string]interface{})

	return HSTSConfiguration{
		IsEnabled:          hstsMap["is_enabled"].(bool),
		MaxAge:             hstsMap["max_age"].(int),
		PreLoaded:          hstsMap["pre_loaded"].(bool),
		SubDomainsIncluded: hstsMap["sub_domains_included"].(bool),
	}
}

func getSSLSettingsDTO(d *schema.ResourceData) SSLSettingsDTO {
	// setup hsts config structure
	hstsSettings := mapHSTSResourceToHSTSDTO(d)
	// scale - add other structures here...

	return SSLSettingsDTO{
		[]Data{
			{
				HstsConfiguration: hstsSettings,
				// add more setting types here
			},
		},
	}
}
