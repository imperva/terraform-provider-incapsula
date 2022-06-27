package incapsula

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceWAFLogSetup() *schema.Resource {
	return &schema.Resource{
		Create: resourceWAFLogSetupCreate,
		Read:   resourceWAFLogSetupRead,
		Update: resourceWAFLogSetupCreate,
		Delete: resourceWAFLogSetupDelete,

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "The Numeric identifier of the account to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			// Optional Arguments
			"enabled": {
				Description: "A boolean flag to enable or disable WAF Logs.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"sftp_host": {
				Description:   "The IP address of your SFTP server.\n.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"sftp_user_name", "sftp_password", "sftp_destination_folder"},
				ConflictsWith: []string{"s3_bucket_name", "s3_access_key", "s3_secret_key"},
			},
			"sftp_user_name": {
				Description:   "A username that will be used to log in to the SFTP server.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"sftp_host", "sftp_password", "sftp_destination_folder"},
				ConflictsWith: []string{"s3_bucket_name", "s3_access_key", "s3_secret_key"},
			},
			"sftp_password": {
				Description:   "A corresponding password for the user account used to log in to the SFTP server.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				RequiredWith:  []string{"sftp_host", "sftp_user_name", "sftp_destination_folder"},
				ConflictsWith: []string{"s3_bucket_name", "s3_access_key", "s3_secret_key"},
			},
			"sftp_destination_folder": {
				Description:   "The path to the directory on the SFTP server.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"sftp_host", "sftp_user_name", "sftp_password"},
				ConflictsWith: []string{"s3_bucket_name", "s3_access_key", "s3_secret_key"},
			},
			"s3_bucket_name": {
				Description:   "S3 bucket name.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"s3_access_key", "s3_secret_key"},
				ConflictsWith: []string{"sftp_host", "sftp_user_name", "sftp_password", "sftp_destination_folder"},
			},
			"s3_access_key": {
				Description:   "S3 access key.",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"s3_bucket_name", "s3_secret_key"},
				ConflictsWith: []string{"sftp_host", "sftp_user_name", "sftp_password", "sftp_destination_folder"},
			},
			"s3_secret_key": {
				Description:   "S3 secret key.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				RequiredWith:  []string{"s3_bucket_name", "s3_secret_key"},
				ConflictsWith: []string{"sftp_host", "sftp_user_name", "sftp_password", "sftp_destination_folder"},
			},
		},
	}
}

func resourceWAFLogSetupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceWAFLogSetupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula WAF Log Setup for account: %d\n", accountID)
	log.Printf("[INFO] enabled: %t\n", d.Get("enabled").(bool))
	log.Printf("[INFO] sftp_host: %s\n", d.Get("sftp_host"))
	log.Printf("[INFO] sftp_user_name: %s\n", d.Get("sftp_user_name"))
	//log.Printf("[INFO] sftp_password: %s\n", d.Get("sftp_password"))
	log.Printf("[INFO] sftp_destination_folder: %s\n", d.Get("sftp_destination_folder"))
	log.Printf("[INFO] s3_bucket_name: %s\n", d.Get("s3_bucket_name"))
	log.Printf("[INFO] s3_access_key: %s\n", d.Get("s3_access_key"))
	//log.Printf("[INFO] s3_secret_key: %s\n", d.Get("s3_secret_key"))

	var wafLogSetupResponse *WAFLogSetupResponse
	var err error

	wafLogSetupPayload := WAFLogSetupPayload{
		accountID,
		d.Get("enabled").(bool),
		d.Get("s3_bucket_name").(string),
		d.Get("s3_access_key").(string),
		d.Get("s3_secret_key").(string),
		d.Get("sftp_host").(string),
		d.Get("sftp_user_name").(string),
		d.Get("sftp_password").(string),
		d.Get("sftp_destination_folder").(string),
	}

	if d.Get("s3_bucket_name") != "" {
		wafLogSetupResponse, err = client.AddWAFLogSetupS3(&wafLogSetupPayload)
	} else if d.Get("sftp_destination_folder") != "" {
		wafLogSetupResponse, err = client.AddWAFLogSetupSFTP(&wafLogSetupPayload)
	} else {
		return errors.New("[ERROR]  Either sftp_* or s3_* arguments are required group")
	}

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula WAF Log Setup for account %d, %s\n", accountID, err)
		return err
	}

	d.SetId(strconv.Itoa(accountID))
	log.Printf("[INFO] red and res_message:  %d, %s\n", wafLogSetupResponse.Res, wafLogSetupResponse.ResMessage)
	log.Printf("[INFO] Created Incapsula WAF Log Setup for account %d\n", accountID)

	return nil
}

func resourceWAFLogSetupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)

	log.Printf("[INFO] Restoring Incapsula WAF Log Setup to default for account: %d\n", accountID)
	wafLogSetupResponse, err := client.DeleteWAFLogSetup(accountID)

	if err != nil {
		log.Printf("[ERROR] Could not restore Incapsula WAF Log Setup to default for account %d, %s\n", accountID, err)
		return err
	}

	d.SetId("")
	log.Printf("[INFO] red and res_message:  %d, %s\n", wafLogSetupResponse.Res, wafLogSetupResponse.ResMessage)
	log.Printf("[INFO] Restored Incapsula WAF Log Setup to default for account %d\n", accountID)

	return nil
}
