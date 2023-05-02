package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

const AbpProvider = "ABP"

var AbpDatasets = []string{"ABP"}

const NetsecProvider = "NETSEC"

var NetsecDatasets = []string{"CONNECTION", "IP", "NETFLOW", "ATTACK"}

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
				Description:  "Type of the producer.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{AbpProvider, NetsecProvider}, false),
			},
			"datasets": {
				Description: "All datasets for the supported producers.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{AbpDatasets[0], NetsecDatasets[0], NetsecDatasets[1], NetsecDatasets[2], NetsecDatasets[3]}, false),
				},
				Required: true,
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
		},
	}
}

func resourceValidation(d *schema.ResourceData) error {
	producer := d.Get("producer").(string)
	datasets := d.Get("datasets").([]interface{})

	var providerDatasets []string
	if producer == AbpProvider {
		providerDatasets = AbpDatasets
	} else if producer == NetsecProvider {
		providerDatasets = NetsecDatasets
	}

	for _, s := range datasets {
		found := false
		for _, k := range providerDatasets {
			if s.(string) == k {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("[ERROR] Unsupported dataset %v for producesr %s", datasets, producer)
		}
	}
	return nil
}

func resourceSiemLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	resErr := resourceValidation(d)
	if resErr != nil {
		return resErr
	}

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
	reponse, statusCode, err := client.ReadSiemLogConfiguration(d.Id(), d.Get("account_id").(string))
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
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemLogConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	resErr := resourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	_, _, err := client.UpdateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		ID:                d.Id(),
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
