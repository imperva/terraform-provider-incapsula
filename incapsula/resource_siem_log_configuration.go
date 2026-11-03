package incapsula

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var BatchPackagingTypes = []string{"NONE", "IMPERVA_METADATA_PREFIX"}

func resourceSiemLogConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemLogConfigurationCreate,
		Read:   resourceSiemLogConfigurationRead,
		Update: resourceSiemLogConfigurationUpdate,
		Delete: resourceSiemLogConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected account_id/logConfiguration_id", d.Id())
				}

				accountID := idSlice[0]
				d.Set("account_id", accountID)

				confID := idSlice[1]
				d.SetId(confID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Client account id.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"configuration_name": {
				Description: "Name of the configuration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"producer": {
				Description: "Type of the producer. Validated by the API.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"datasets": {
				Description: "Datasets to send. Validated against the producer by the API.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
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
			"format": {
				Description: "The format of the logs of this log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
			},
			"logs_level": {
				Description: "The default logs level of this log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
			},
			"compress_logs": {
				Description: "True if the logs are compressed, false otherwise.",
				Type:        schema.TypeBool,
				Optional:    true,
				Required:    false,
			},
			"public_key": {
				Description: "The public key of the log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
			},
			"public_key_file_name": {
				Description: "The public key file name of the log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
			},
			"batch_packaging_type": {
				Description:  "How log batches are packaged in the payload sent to the connection.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(BatchPackagingTypes, false),
			},
		},
	}
}

func buildSiemLogConfigurationPayload(d *schema.ResourceData) *Payload {
	batchPackagingType := d.Get("batch_packaging_type").(string)
	if batchPackagingType == "" {
		return nil
	}
	return &Payload{BatchPackaging: BatchPackaging{Type: batchPackagingType}}
}

func batchPackagingTypeFromPayload(payload *Payload) string {
	if payload == nil {
		return ""
	}
	return payload.BatchPackaging.Type
}

func resourceSiemLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	response, statusCode, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		AssetID:           d.Get("account_id").(string),
		ConfigurationName: d.Get("configuration_name").(string),
		Provider:          d.Get("producer").(string),
		Datasets:          d.Get("datasets").([]interface{}),
		Enabled:           d.Get("enabled").(bool),
		ConnectionId:      d.Get("connection_id").(string),
		CompressLogs:      d.Get("compress_logs").(bool),
		Format:            d.Get("format").(string),
		LogsLevel:         d.Get("logs_level").(string),
		PublicKey:         d.Get("public_key").(string),
		PublicKeyFileNAme: d.Get("public_key_file_name").(string),
		Payload:           buildSiemLogConfigurationPayload(d),
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
	reponse, statusCode, err := client.ReadSiemLogConfiguration(d.Id(), d.Get("account_id").(string))
	// If the resource is deleted on the server, blow it out locally and run through the normal TF cycle.
	// The client reports any non-200 as an error, so 404 is checked first.
	if statusCode != nil && *statusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if (*statusCode == 200) && (reponse != nil) && (len(reponse.Data) == 1) {
		var logConfiguration = reponse.Data[0]
		d.Set("account_id", logConfiguration.AssetID)
		d.Set("configuration_name", logConfiguration.ConfigurationName)
		d.Set("producer", logConfiguration.Provider)
		d.Set("datasets", logConfiguration.Datasets)
		d.Set("enabled", logConfiguration.Enabled)
		d.Set("connection_id", logConfiguration.ConnectionId)
		d.Set("compress_logs", logConfiguration.CompressLogs)
		d.Set("format", logConfiguration.Format)
		d.Set("logs_level", logConfiguration.LogsLevel)
		d.Set("public_key", logConfiguration.PublicKey)
		d.Set("public_key_file_name", logConfiguration.PublicKeyFileNAme)
		d.Set("batch_packaging_type", batchPackagingTypeFromPayload(logConfiguration.Payload))
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
		CompressLogs:      d.Get("compress_logs").(bool),
		Format:            d.Get("format").(string),
		LogsLevel:         d.Get("logs_level").(string),
		PublicKey:         d.Get("public_key").(string),
		PublicKeyFileNAme: d.Get("public_key_file_name").(string),
		Payload:           buildSiemLogConfigurationPayload(d),
	}}})

	if err != nil {
		return err
	}

	return nil
}

func resourceSiemLogConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ID := d.Id()
	accountId := d.Get("account_id").(string)
	_, err := client.DeleteSiemLogConfiguration(ID, accountId)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")
	return nil
}
