package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"time"
)

func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubAccountCreate,
		Read:   resourceSubAccountRead,
		Delete: resourceSubAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"sub_account_name": {
				Description: "The name of the new sub-account.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			// Optional Arguments
			"parent_id": {
				Description: "The newly created account's parent id. If not specified, the invoking account will be assigned as the parent.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"ref_id": {
				Description: "Customer specific identifier for this operation.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"logs_account_id": {
				Description: "Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"log_level": {
				Description:  "The log level. Options are `full`, `security`, `none` and `default`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"full", "security", "none", "default"}, false),
			},
		},
	}
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	subAccountName := d.Get("sub_account_name").(string)

	log.Printf("[INFO] Creating Incapsula subaccount: %s\n", subAccountName)
	log.Printf("[INFO] logs_account_id: %d\n", d.Get("logs_account_id").(int))
	log.Printf("[INFO] log_level: %s\n", d.Get("log_level").(string))
	log.Printf("[INFO] parent_id: %d\n", d.Get("parent_id").(int))
	log.Printf("[INFO] ref_id: %s\n", d.Get("ref_id").(string))

	subAccountPayload := SubAccountPayload{subAccountName,
		d.Get("ref_id").(string),
		d.Get("log_level").(string),
		d.Get("parent_id").(int),
		d.Get("logs_account_id").(int)}

	SubAccountAddResponse, err := client.AddSubAccount(&subAccountPayload)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula subaccount %s, %s\n", subAccountName, err)
		return err
	}

	// Set the SubAccount ID
	d.SetId(strconv.Itoa(SubAccountAddResponse.SubAccount.SubAccountID))
	log.Printf("[DEBUG] Account id for new sub account : %d", SubAccountAddResponse.SubAccount.SubAccountID)
	log.Printf("[INFO] Created Incapsula subaccount %s\n", subAccountName)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	return resourceSubAccountRead(d, m)
}

func resourceSubAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	subAccountID, _ := strconv.Atoi(d.Id())
	subAccount, err := client.GetSubAccount(d.Get("parent_id").(int), subAccountID)

	if err != nil {
		return err
	}

	if subAccount == nil {
		log.Printf("[INFO] Incapsula subaccount %s has already been deleted: %s\n", d.Id(), err)
		d.SetId("")
		return nil
	}

	d.Set("sub_account_name", subAccount.SubAccountName)
	d.Set("ref_id", subAccount.RefID)
	d.Set("log_level", subAccount.LogLevel)
	d.Set("parent_id", subAccount.ParentID)
	d.Set("logs_account_id", subAccount.LogsAccountID)

	log.Printf("[INFO] Finished reading Incapsula subaccount: %s\n", d.Id())

	return nil
}

func resourceSubAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	subAccountID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Deleting Incapsula subaccount id: %d\n", subAccountID)

	err := client.DeleteSubAccount(subAccountID)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula subaccount id: %d, %s\n", subAccountID, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula subaccount id: %d\n", subAccountID)

	return nil
}
