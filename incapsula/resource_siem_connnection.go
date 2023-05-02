package incapsula

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

const StorageTypeCustomerS3 = "CUSTOMER_S3"
const StorageTypeCustomerS3Arn = "CUSTOMER_S3_ARN"
const sensitiveDataPlaceholder = "Sensitive data placeholder"

func resourceSiemConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemConnectionCreate,
		Read:   resourceSiemConnectionRead,
		Update: resourceSiemConnectionUpdate,
		Delete: resourceSiemConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected account_id/connection_id", d.Id())
				}

				accountID := idSlice[0]
				d.Set("account_id", accountID)

				connectionID := idSlice[1]
				d.SetId(connectionID)

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
			"connection_name": {
				Description: "Name of the connection.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"storage_type": {
				Description:  "Type of the storage.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{StorageTypeCustomerS3, StorageTypeCustomerS3Arn}, false),
			},
			"access_key": {
				Description: "Access key in AWS.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					expectedLen := 20
					actualLen := len(val.(string))
					if actualLen != expectedLen && val != sensitiveDataPlaceholder {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, expectedLen, actualLen))
					}
					return
				},
				RequiredWith: []string{"secret_key"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == sensitiveDataPlaceholder {
						return true
					}
					return false
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
					if actualLen != expectedLen && val != sensitiveDataPlaceholder {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, expectedLen, actualLen))
					}
					return
				},
				RequiredWith: []string{"access_key"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == sensitiveDataPlaceholder {
						return true
					}
					return false
				},
			},
			"path": {
				Description: "Store data from the specified connection under this path.",
				Type:        schema.TypeString,
				Required:    true,
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

func siemConnectionResourceValidation(d *schema.ResourceData) error {
	storageType := d.Get("storage_type").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	if storageType == StorageTypeCustomerS3 && (accessKey == "" || secretKey == "") {
		return fmt.Errorf("[ERROR] access_key and secret_key should be provided for storage_type=%s", storageType)
	} else if storageType == StorageTypeCustomerS3Arn && (accessKey != "" || secretKey != "") {
		return fmt.Errorf("[ERROR] access_key and secret_key should not be provided for storage_type=%s", storageType)
	}
	return nil
}
func resourceSiemConnectionCreate(d *schema.ResourceData, m interface{}) error {
	resErr := siemConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	response, statusCode, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			AccessKey: d.Get("access_key").(string),
			SecretKey: d.Get("secret_key").(string),
			Path:      d.Get("path").(string),
		},
	}}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemConnectionRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	response, statusCode, err := client.ReadSiemConnection(d.Id(), d.Get("account_id").(string))
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
		if connection.StorageType == StorageTypeCustomerS3 {
			d.Set("access_key", connection.ConnectionInfo.AccessKey)
			d.Set("input_hash", connection.ConnectionInfo.SecretKey)
		}
		d.Set("path", connection.ConnectionInfo.Path)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	resErr := siemConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	_, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    d.Get("storage_type").(string),
		ConnectionInfo: ConnectionInfo{
			AccessKey: d.Get("access_key").(string),
			SecretKey: d.Get("secret_key").(string),
			Path:      d.Get("path").(string),
		},
	}}})

	if err != nil {
		return err
	}
	return nil
}

func resourceSiemConnectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ID := d.Id()
	accountId := d.Get("account_id").(string)

	_, err := client.DeleteSiemConnection(ID, accountId)

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
