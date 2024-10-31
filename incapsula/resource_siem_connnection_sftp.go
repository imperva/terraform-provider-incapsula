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
				Type:        schema.TypeInt,
				Required:    true,
			},
			"password": {
				Description: "Sftp password.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"path": {
				Description: "Sftp path.",
				Type:        schema.TypeInt,
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
			Host:                    d.Get("host").(string),
			Username:                d.Get("username").(string),
			Password:                d.Get("password").(string),
			Path:                    d.Get("path").(string),
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
// 		todo should there be also password ?
		d.Set("input_hash", connectionInfo.Token)
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
			Host:                    d.Get("host").(string),
			Username:                d.Get("username").(string),
			Password:                d.Get("password").(string),
			Path:                    d.Get("path").(string),
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