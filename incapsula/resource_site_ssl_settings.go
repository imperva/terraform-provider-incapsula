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

var inboundTLSSettingsResource = schema.Resource{
	Schema: map[string]*schema.Schema{

		"configuration_profile": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tls_configuration": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"tls_version": {
						Type:     schema.TypeString,
						Required: true,
					},
					"ciphers_support": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
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
			"inbound_tls_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &inboundTLSSettingsResource,
				Set:      schema.HashResource(&inboundTLSSettingsResource),
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

	mapHSTSResponseToHSTSResource(d, settingsData)
	mapInboundTLSSettingsResponseToResource(d, settingsData)
	// map other settings here

	return nil
}

func resourceSiteSSLSettingsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	// currently only disables HSTS and set default InboundTLSSettings
	// If more settings are implemented in the endpoint, add delete logic for them here.
	setting := prepareDisableHSTSStructure()
	prepareDefaultTLSStructure(&setting)
	var _, err = client.UpdateSiteSSLSettings(d.Get("site_id").(int), setting)

	if err != nil {
		return err
	}

	return nil
}

func prepareDisableHSTSStructure() SSLSettingsResponse {
	disableHSTSSetting := HSTSConfiguration{
		IsEnabled: false,
	}

	return SSLSettingsResponse{
		[]SSLSettingsDTO{
			{
				HstsConfiguration: disableHSTSSetting,
				// add more setting types here
			},
		},
	}
}

func prepareDefaultTLSStructure(settingsData *SSLSettingsResponse) {
	defaultTLSSettings := &InboundTLSSettingsConfiguration{
		ConfigurationProfile: "DEFAULT",
		TLSConfigurations:    []TLSConfiguration{},
	}

	settingsData.Data[0].InboundTLSSettingsConfiguration = defaultTLSSettings
}

func mapHSTSResponseToHSTSResource(d *schema.ResourceData, settingsData *SSLSettingsResponse) {
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

func mapInboundTLSSettingsResponseToResource(d *schema.ResourceData, settingsData *SSLSettingsResponse) {

	var inboundTLSSettingsFromServer *InboundTLSSettingsConfiguration
	inboundTLSSettingsFromServer = settingsData.Data[0].InboundTLSSettingsConfiguration

	if inboundTLSSettingsFromServer == nil {
		d.Set("inbound_tls_settings", nil)
		return
	}
	inboundTLSSettingsMap := make(map[string]interface{})
	inboundTLSSettingsMap["configuration_profile"] = inboundTLSSettingsFromServer.ConfigurationProfile

	tlsConfigurations := make([]map[string]interface{}, 0)
	for _, tlsConfig := range inboundTLSSettingsFromServer.TLSConfigurations {
		tlsConfigMap := make(map[string]interface{})
		tlsConfigMap["tls_version"] = tlsConfig.TLSVersion
		tlsConfigMap["ciphers_support"] = toStringInterfaceSlice(tlsConfig.CiphersSupport)

		tlsConfigurations = append(tlsConfigurations, tlsConfigMap)
	}

	inboundTLSSettingsMap["tls_configuration"] = tlsConfigurations

	d.Set("inbound_tls_settings", []map[string]interface{}{inboundTLSSettingsMap})

}

func mapInboundTLSSettingsResourceToDTO(resourceData *schema.ResourceData) *InboundTLSSettingsConfiguration {
	inboundTLSSettings, ok := resourceData.Get("inbound_tls_settings").(*schema.Set)
	if !ok || inboundTLSSettings.Len() == 0 {
		return nil
	}
	inboundTLSSettingsList := inboundTLSSettings.List()
	inboundTLSSettingsMap := inboundTLSSettingsList[0].(map[string]interface{})

	dto := &InboundTLSSettingsConfiguration{
		ConfigurationProfile: inboundTLSSettingsMap["configuration_profile"].(string),
		TLSConfigurations:    make([]TLSConfiguration, 0),
	}

	if tlsConfigurations, ok := inboundTLSSettingsMap["tls_configuration"].([]interface{}); ok {
		for _, tlsConfig := range tlsConfigurations {
			tlsConfigMap := tlsConfig.(map[string]interface{})
			tlsVersion := tlsConfigMap["tls_version"].(string)
			ciphersSupport := tlsConfigMap["ciphers_support"].([]interface{})

			tlsConfigDTO := TLSConfiguration{
				TLSVersion:     tlsVersion,
				CiphersSupport: toStringSlice(ciphersSupport),
			}

			dto.TLSConfigurations = append(dto.TLSConfigurations, tlsConfigDTO)
		}
	}

	return dto
}

// Helper function to convert []interface{} to []string
func toStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.(string)
	}
	return result
}

// Helper function to convert []string to []interface{}
func toStringInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

func getSSLSettingsDTO(d *schema.ResourceData) SSLSettingsResponse {

	// setup hsts config structure
	hstsSettings := mapHSTSResourceToHSTSDTO(d)
	// setup inbound TLS settings structure
	inboundTLSSettings := mapInboundTLSSettingsResourceToDTO(d)
	// scale - add other structures here...

	return SSLSettingsResponse{
		[]SSLSettingsDTO{
			{
				HstsConfiguration:               hstsSettings,
				InboundTLSSettingsConfiguration: inboundTLSSettings,
				// add more setting types here
			},
		},
	}
}
