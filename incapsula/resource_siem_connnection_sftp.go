package incapsula

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

const StorageTypeCustomerSftp = "CUSTOMER_SFTP"

func resourceSiemSftpConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemSftpConnectionCreate,
		Read:   resourceSiemSftpConnectionRead,
		Update: resourceSiemSftpConnectionUpdate,
		Delete: resourceSiemSftpConnectionDelete,
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
			"host": {
				Description: "Sftp endpoint host.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
			},
			"username": {
				Description: "Sftp username.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "Sftp password.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					minLen := 1
					actualLen := len(val.(string))
					if actualLen < minLen && val != sensitiveDataPlaceholder {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, minLen, actualLen))
					}
					return
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == sensitiveDataPlaceholder {
						return true
					}
					return false
				},
			},
			"path": {
				Description: "Sftp path.",
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
					newHash := createSftpSiemConnectionHash(d)
					if newHash == old {
						return true
					}
					return false
				},
			},
		},
	}
}

func siemSftpConnectionResourceValidation(d *schema.ResourceData) error {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	path := d.Get("path").(string)
	if host == "" || username == "" || password == "" || path == "" {
		return fmt.Errorf("[ERROR] host, username, password and path should be provided for incapsula_siem_sftp_connection")
	}
	return nil
}
func resourceSiemSftpConnectionCreate(d *schema.ResourceData, m interface{}) error {
	resErr := siemSftpConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	response, statusCode, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    StorageTypeCustomerSftp,
		ConnectionInfo: SftpConnectionInfo{
			Host:     d.Get("host").(string),
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Path:     d.Get("path").(string),
		},
	}}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemSftpConnectionRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemSftpConnectionRead(d *schema.ResourceData, m interface{}) error {
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
		connectionInfo := connection.ConnectionInfo.(SftpConnectionInfo)
		d.Set("host", connectionInfo.Host)
		d.Set("username", connectionInfo.Username)
		d.Set("path", connectionInfo.Path)
		d.Set("input_hash", connectionInfo.Password)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemSftpConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	resErr := siemSftpConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	_, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    StorageTypeCustomerSftp,
		ConnectionInfo: SftpConnectionInfo{
			Host:     d.Get("host").(string),
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Path:     d.Get("path").(string),
		},
	}}})

	if err != nil {
		return err
	}
	return nil
}

func resourceSiemSftpConnectionDelete(d *schema.ResourceData, m interface{}) error {
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
func createSftpSiemConnectionHash(d *schema.ResourceData) string {
	password := d.Get("password").(string)
	result := calculateSftpSiemConnectionHash(password)
	return result
}

func calculateSftpSiemConnectionHash(token string) string {
	h := sha256.New()
	stringForHash := token
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)
	return result
}
