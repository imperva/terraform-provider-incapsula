package incapsula

import (
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description:  "The log level. Options are `full`, `security`, `none` and `default`. Defaults to `default`.",
				Type:         schema.TypeString,
				Default:      "default",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"full", "security", "none", "default"}, false),
			},

			// Computed Attributes
			"sub_account_id": {
				Description: "SubAccount ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"is_for_special_ssl_configuration": {
				Description: "Is using SSL configuration",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"support_levels": {
				Description: "Support level",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	sub_account_name := d.Get("sub_account_name").(string)

	log.Printf("[INFO] Creating Incapsula subaccount: %s\n", sub_account_name)

	log.Printf("[INFO] logs_account_id: %d\n", d.Get("logs_account_id").(int))
	log.Printf("[INFO] log_level: %s\n", d.Get("log_level").(string))
	log.Printf("[INFO] parent_id: %d\n", d.Get("parent_id").(int))

	SubAccountAddResponse, err := client.AddSubAccount(
		sub_account_name,
		d.Get("ref_id").(string),
		d.Get("log_level").(string),
		d.Get("logs_account_id").(int),
		d.Get("parent_id").(int),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula subaccount %s, %s\n", sub_account_name, err)
		return err
	}

	// Set the SubAccount ID
	d.SetId(strconv.Itoa(SubAccountAddResponse.SubAccount.SubAccountID))
	log.Printf("[DEBUG] Account id for new sub account : %d", SubAccountAddResponse.SubAccount.SubAccountID)
	log.Printf("[INFO] Created Incapsula subaccount %s\n", sub_account_name)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	return resourceSubAccountRead(d, m)
}

func resourceSubAccountRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data center
	client := m.(*Client)
	subAccountID, _ := strconv.Atoi(d.Id())
	listSubAccountsResponse, err := client.ListSubAccounts(d.Get("parent_id").(int))

	if err != nil {
		return err
	}

	found := false

	for _, subAccount := range listSubAccountsResponse.SubAccounts {
		if subAccount.SubAccountID == subAccountID {
			log.Printf("[INFO] subaccount : %v\n", subAccount)
			d.Set("sub_account_id", subAccount.SubAccountID)
			d.Set("sub_account_name", subAccount.SubAccountName)
			d.Set("is_for_special_ssl_configuration", subAccount.SpeicalSSL)
			d.Set("support_levels", subAccount.SupportLevel)
			found = true
			break
		}
	}

	if !found {
		log.Printf("[INFO] Incapsula subaccount %s has already been deleted: %s\n", d.Id(), err)
		d.SetId("")
		return nil
	}

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
