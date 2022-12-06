package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSiemLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemLogConfigurationInfo := SiemLogConfigurationInfo{
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("producer").(string),
		Datasets:          d.Get("datasets").([]interface{}),
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
		if len(logConfiguration.AssetID) > 0 {
			d.Set("account_id", logConfiguration.AssetID)
		}
		d.Set("configuration_name", logConfiguration.ConfigurationName)
		d.Set("producer", logConfiguration.Provider)
		d.Set("datasets", logConfiguration.Datasets)
		d.Set("enabled", logConfiguration.Enabled)
		d.Set("connection_id", logConfiguration.ConnectionId)
		if len(logConfiguration.Version) > 0 {
			d.Set("version", logConfiguration.Version)
		}
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
		Provider:          d.Get("producer").(string),
		Datasets:          d.Get("datasets").([]interface{}),
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
