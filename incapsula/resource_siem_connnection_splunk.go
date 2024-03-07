package incapsula

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

const StorageTypeCustomerSplunk = "CUSTOMER_SPLUNK"

func resourceSiemSplunkConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemSplunkConnectionCreate,
		Read:   resourceSiemSplunkConnectionRead,
		Update: resourceSiemSplunkConnectionUpdate,
		Delete: resourceSiemSplunkConnectionDelete,
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
			"token": {
				Description: "Splunk access token.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					expectedLen := 36
					actualLen := len(val.(string))
					if actualLen != expectedLen && val != sensitiveDataPlaceholder {
						errs = append(errs, fmt.Errorf("%q length should be %d, got: %d", key, expectedLen, actualLen))
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
			"host": {
				Description: "Splunk endpoint host.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
			},
			"port": {
				Description: "Splunk endpoint port.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"disable_cert_verification": {
				Description: "flag to disable ssl cert verification",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"input_hash": {
				Description: "inputHash",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newHash := createSplunkSiemConnectionHash(d)
					if newHash == old {
						return true
					}
					return false
				},
			},
		},
	}
}

func siemSplunkConnectionResourceValidation(d *schema.ResourceData) error {
	host := d.Get("host").(string)
	portExists := d.Get("port")
	token := d.Get("token").(string)
	if host == "" || portExists == false || token == "" {
		return fmt.Errorf("[ERROR] host, port and token should be provided for incapsula_siem_splunk_connection")
	}
	return nil
}
func resourceSiemSplunkConnectionCreate(d *schema.ResourceData, m interface{}) error {
	resErr := siemSplunkConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	response, statusCode, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    StorageTypeCustomerSplunk,
		ConnectionInfo: SplunkConnectionInfo{
			Host:                    d.Get("host").(string),
			Port:                    d.Get("port").(int),
			Token:                   d.Get("token").(string),
			DisableCertVerification: d.Get("disable_cert_verification").(bool),
		},
	}}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemSplunkConnectionRead(d, m)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemSplunkConnectionRead(d *schema.ResourceData, m interface{}) error {
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
		connectionInfo := connection.ConnectionInfo.(SplunkConnectionInfo)
		d.Set("host", connectionInfo.Host)
		d.Set("port", connectionInfo.Port)
		//d.Set("token", connectionInfo.Token)
		d.Set("disable_cert_verification", connectionInfo.DisableCertVerification)
		d.Set("input_hash", connectionInfo.Token)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func resourceSiemSplunkConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	resErr := siemSplunkConnectionResourceValidation(d)
	if resErr != nil {
		return resErr
	}

	client := m.(*Client)
	_, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    StorageTypeCustomerSplunk,
		ConnectionInfo: SplunkConnectionInfo{
			Host:                    d.Get("host").(string),
			Port:                    d.Get("port").(int),
			Token:                   d.Get("token").(string),
			DisableCertVerification: d.Get("disable_cert_verification").(bool),
		},
	}}})

	if err != nil {
		return err
	}
	return nil
}

func resourceSiemSplunkConnectionDelete(d *schema.ResourceData, m interface{}) error {
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

func createSplunkSiemConnectionHash(d *schema.ResourceData) string {
	token := d.Get("token").(string)
	result := calculateSplunkSiemConnectionHash(token)
	return result
}

func calculateSplunkSiemConnectionHash(token string) string {
	h := sha256.New()
	stringForHash := token
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)
	return result
}
