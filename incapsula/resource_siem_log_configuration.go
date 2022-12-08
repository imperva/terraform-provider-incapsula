package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSiemLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	response, statusCode, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("producer").(string),
		Datasets:          d.Get("datasets").([]interface{}),
		Enabled:           d.Get("enabled").(bool),
		ConnectionId:      d.Get("connection_id").(string),
	}}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemLogConfigurationRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemLogConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	reponse, statusCode, err := client.ReadSiemLogConfiguration(d.Id())
	if err != nil {
		return err
	}
	// If the connection is deleted on the server, blow it out locally and run through the normal TF cycle
	if *statusCode == 404 {
		d.SetId("")
		return nil
	} else if (*statusCode == 200) && (reponse != nil) && (len(reponse.Data) == 1) {
		var logConfiguration = reponse.Data[0]
		d.Set("account_id", logConfiguration.AssetID)
		d.Set("configuration_name", logConfiguration.ConfigurationName)
		d.Set("producer", logConfiguration.Provider)
		d.Set("datasets", logConfiguration.Datasets)
		d.Set("enabled", logConfiguration.Enabled)
		d.Set("connection_id", logConfiguration.ConnectionId)
		d.Set("version", logConfiguration.Version)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemLogConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, _, err := client.UpdateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		ID:                d.Id(),
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("producer").(string),
		Datasets:          d.Get("datasets").([]interface{}),
		Enabled:           d.Get("enabled").(bool),
		ConnectionId:      d.Get("connection_id").(string),
		Version:           d.Get("version").(string),
	}}})

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
