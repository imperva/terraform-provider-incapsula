package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSiemLogConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemLogConfigurationCreate,
		Read:   resourceSiemLogConfigurationRead,
		Update: resourceSiemLogConfigurationUpdate,
		Delete: resourceSiemLogConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Client account id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"configuration_name": {
				Description: "Name of the configuration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provider": {
				Description:  "Type of the provider.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ABP", "NETSEC"}, false),
			},
			"datasets": {
				Description:  "All datasets for the supported providers.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ABP", "ABP_ACCESS", "CONNECTION", "NETFLOW, IP", "ATTACK", "ABP_TEST"}, false),
			},
			"enabled": {
				Description: "True if the connection is enabled, false otherwise.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"connection_id": {
				Description: "The id of the connection for this log configuration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "Version of the log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSiemLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemLogConfigurationInfo := SiemLogConfigurationInfo{
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("provider").(string),
		Datasets:          d.Get("datasets").([]string),
		Enabled:           d.Get("enabled").(bool),
		ConnectionId:      d.Get("connection_id").(string),
	}

	siemLogConfiguration := SiemLogConfiguration{}
	siemLogConfiguration.Data = []SiemLogConfigurationInfo{siemLogConfigurationInfo}

	siemLogConfigurationWithIdAndVersion, responseStatusCode, err := client.CreateSiemLogConfiguration(&siemLogConfiguration)
	if err != nil {
		return err
	}

	if (*responseStatusCode == 201) && (siemLogConfigurationWithIdAndVersion != nil) && (len(siemLogConfigurationWithIdAndVersion.Data) == 1) {
		d.SetId(siemLogConfigurationWithIdAndVersion.Data[0].ID)
		return resourceSiemLogConfigurationRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceSiemLogConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemLogConfigurationWithIdAndVersion, responseStatusCode, err := client.ReadSiemLogConfiguration(d.Id())
	if err != nil {
		return err
	}
	// If the connection is deleted on the server, blow it out locally and run through the normal TF cycle
	if *responseStatusCode == 404 {
		d.SetId("")
		return nil
	} else if (*responseStatusCode == 200) && (siemLogConfigurationWithIdAndVersion != nil) && (len(siemLogConfigurationWithIdAndVersion.Data) == 1) {
		var logConfiguration = siemLogConfigurationWithIdAndVersion.Data[0]
		d.Set("account_id", logConfiguration.AssetID)
		d.Set("configuration_name", logConfiguration.ConfigurationName)
		d.Set("provider", logConfiguration.Provider)
		d.Set("datasets", logConfiguration.Datasets)
		d.Set("enabled", logConfiguration.Enabled)
		d.Set("connection_id", logConfiguration.ConnectionId)
		d.Set("version", logConfiguration.Version)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceSiemLogConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemLogConfigurationWithIdAndVersionInfo := SiemLogConfigurationWithIdAndVersionInfo{
		ID:                d.Id(),
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("provider").(string),
		Datasets:          d.Get("datasets").([]string),
		Enabled:           d.Get("enabled").(bool),
		ConnectionId:      d.Get("connection_id").(string),
		Version:           d.Get("version").(string),
	}

	siemLogConfigurationWithIdAndVersion := SiemLogConfigurationWithIdAndVersion{}
	siemLogConfigurationWithIdAndVersion.Data = []SiemLogConfigurationWithIdAndVersionInfo{siemLogConfigurationWithIdAndVersionInfo}

	_, _, err := client.UpdateSiemLogConfiguration(&siemLogConfigurationWithIdAndVersion)

	if err != nil {
		return err
	}

	return nil
}

func resourceSiemLogConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ID := d.Id()

	_, err := client.DeleteSiemLogConfiguration(ID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")
	return nil
}
