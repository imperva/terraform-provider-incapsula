package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSiemConnectionS3Arn() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemConnectionS3ArnCreate,
		Read:   resourceSiemConnectionS3ArnRead,
		Update: resourceSiemConnectionS3ArnUpdate,
		Delete: resourceSiemConnectionS3ArnDelete,
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
			"connection_name": {
				Description: "Name of the connection.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"storage_type": {
				Description:  "Type of the storage.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOMER_S3_ARN"}, false),
			},
			"path": {
				Description: "Store data from the specified connection under this path.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"version": {
				Description: "Version of the connection.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSiemConnectionS3ArnCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	response, statusCode, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			Path: d.Get("path").(string),
		},
	}}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemConnectionS3ArnRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemConnectionS3ArnRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	response, statusCode, err := client.ReadSiemConnection(d.Id())
	if err != nil {
		return err
	}
	// If the connection is deleted on the server, blow it out locally and run through the normal TF cycle
	if *statusCode == 404 {
		d.SetId("")
		return nil
	} else if (*statusCode == 200) && (response != nil) && (len(response.Data) == 1) {
		var connection = response.Data[0]
		d.Set("account_id", connection.AssetID)
		d.Set("connection_name", connection.ConnectionName)
		d.Set("storage_type", connection.StorageType)
		d.Set("path", connection.ConnectionInfo.Path)
		d.Set("version", connection.Version)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemConnectionS3ArnUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		Version:        d.Get("version").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			Path: d.Get("path").(string),
		},
	}}})

	if err != nil {
		return err
	}
	return nil
}

func resourceSiemConnectionS3ArnDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ID := d.Id()

	_, err := client.DeleteSiemConnection(ID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")
	return nil
}
