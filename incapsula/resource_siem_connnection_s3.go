package incapsula

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSiemConnectionS3() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemConnectionS3Create,
		Read:   resourceSiemConnectionS3Read,
		Update: resourceSiemConnectionS3Update,
		Delete: resourceSiemConnectionS3Delete,
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
				ValidateFunc: validation.StringInSlice([]string{"CUSTOMER_S3"}, false),
			},
			"access_key": {
				Description: "Access key in AWS.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					expectedLen := 20
					actualLen := len(val.(string))
					if actualLen != expectedLen {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, expectedLen, actualLen))
					}
					return
				},
			},
			"secret_key": {
				Description: "Secret key in AWS.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					expectedLen := 40
					actualLen := len(val.(string))
					if actualLen != expectedLen {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, expectedLen, actualLen))
					}
					return
				},
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
			"input_hash": {
				Description: "inputHash",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newHash := createS3SiemConnectionHash(d)
					if newHash == old {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceSiemConnectionS3Create(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemConnectionInfo := SiemConnectionInfo{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			AccessKey: d.Get("access_key").(string),
			SecretKey: d.Get("secret_key").(string),
			Path:      d.Get("path").(string),
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
		return resourceSiemConnectionS3Read(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceSiemConnectionS3Read(d *schema.ResourceData, m interface{}) error {
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
		d.Set("access_key", connection.ConnectionInfo.AccessKey)
		d.Set("path", connection.ConnectionInfo.Path)
		d.Set("version", connection.Version)
		d.Set("input_hash", connection.ConnectionInfo.SecretKey)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *responseStatusCode)
	}
}

func resourceSiemConnectionS3Update(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siemConnectionWithIdAndVersionInfo := SiemConnectionWithIdAndVersionInfo{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		Version:        d.Get("version").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			AccessKey: d.Get("access_key").(string),
			SecretKey: d.Get("secret_key").(string),
			Path:      d.Get("path").(string),
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

func resourceSiemConnectionS3Delete(d *schema.ResourceData, m interface{}) error {
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

func createS3SiemConnectionHash(d *schema.ResourceData) string {
	secretKey := d.Get("secret_key").(string)
	result := calculateS3SiemConnectionHash(secretKey)
	return result
}

func calculateS3SiemConnectionHash(secretKey string) string {
	h := sha256.New()
	stringForHash := secretKey
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)
	return result
}
