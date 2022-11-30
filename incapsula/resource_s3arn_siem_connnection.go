package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceS3ArnSiemConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3ArnSiemConnectionCreate,
		Read:   resourceS3ArnSiemConnectionRead,
		Update: resourceS3ArnSiemConnectionUpdate,
		Delete: resourceS3ArnSiemConnectionDelete,
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

func resourceS3ArnSiemConnectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemConnectionInfo := SiemConnectionInfo{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			Path: d.Get("path").(string),
		},
	}

	siemConnection := SiemConnection{}
	siemConnection.Data = []SiemConnectionInfo{siemConnectionInfo}

	siemConnectionWithIdAndVersion, responseStatusCode, err := client.CreateSiemConnection(&siemConnection)
	if err != nil {
		return err
	}

	if (*responseStatusCode == 201) && (siemConnectionWithIdAndVersion != nil) && (len(siemConnectionWithIdAndVersion.Data) == 1) {
		d.SetId(siemConnectionWithIdAndVersion.Data[0].ID)
		return resourceS3SiemConnectionRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceS3ArnSiemConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemConnectionWithIdAndVersion, responseStatusCode, err := client.ReadSiemConnection(d.Id())
	if err != nil {
		return err
	}
	// If the connection is deleted on the server, blow it out locally and run through the normal TF cycle
	if *responseStatusCode == 404 {
		d.SetId("")
		return nil
	} else if (*responseStatusCode == 200) && (siemConnectionWithIdAndVersion != nil) && (len(siemConnectionWithIdAndVersion.Data) == 1) {
		var connection = siemConnectionWithIdAndVersion.Data[0]
		d.Set("account_id", connection.AssetID)
		d.Set("connection_name", connection.ConnectionName)
		d.Set("storage_type", connection.StorageType)
		d.Set("path", connection.ConnectionInfo.Path)
		d.Set("version", connection.Version)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceS3ArnSiemConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemConnectionWithIdAndVersionInfo := SiemConnectionWithIdAndVersionInfo{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		Version:        d.Get("version").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			Path: d.Get("path").(string),
		},
	}

	siemConnectionWithIdAndVersionUpdate := SiemConnectionWithIdAndVersion{}
	siemConnectionWithIdAndVersionUpdate.Data = []SiemConnectionWithIdAndVersionInfo{siemConnectionWithIdAndVersionInfo}

	_, _, err := client.UpdateSiemConnection(&siemConnectionWithIdAndVersionUpdate)

	if err != nil {
		return err
	}

	return nil
}

func resourceS3ArnSiemConnectionDelete(d *schema.ResourceData, m interface{}) error {
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
