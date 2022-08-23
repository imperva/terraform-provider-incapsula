package incapsula

import (
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"email": {
				Description: "Email address. For example: joe@example.com.",
				Type:        schema.TypeString,
				Required:    true,
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
			},
			"user_name": {
				Description: "The account owner's name. For example: John Doe.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"plan_id": {
				Description: "An identifier of the plan to assign to the new account. For example, ent100 for the Enterprise 100 plan.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"account_name": {
				Description: "Account name.",
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
				Description:  "The log level. Options are `full`, `security`, and `none`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"full", "security", "none"}, false),
			},
			"error_page_template": {
				Description: "Base64 encoded template for an error page.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"support_all_tls_versions": {
				Description: "Allow sites in the account to support all TLS versions for connectivity between clients (visitors) and the Imperva service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"naked_domain_san_for_new_www_sites": {
				Description:  "Add naked domain SAN to Incapsula SSL certificates for new www sites. Options are `true` and `false`. Defaults to `true`",
				Type:         schema.TypeString,
				Default:      "true",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false", "default"}, false),
			},
			"wildcard_san_for_new_sites": {
				Description:  "Add wildcard SAN to Incapsula SSL certificates for new sites. Options are `true`, `false` and `default`. Defaults to `default`",
				Type:         schema.TypeString,
				Default:      "Default",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"True", "False", "Default"}, true),
			},
			"data_storage_region": {
				Description:  "Default data region of the account for newly created sites. Options are `APAC`, `EU`, `US` and `DEFAULT`. Defaults to `DEFAULT`.",
				Type:         schema.TypeString,
				Default:      "US",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"APAC", "EU", "US", "AU"}, false),
			},

			// Computed Attributes
			"support_level": {
				Description: "Account support level",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"plan_name": {
				Description: "Name of plan associate with account",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"trial_end_date": {
				Description: "End date for trial account",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	email := d.Get("email").(string)

	log.Printf("[INFO] Creating Incapsula account for email: %s\n", email)

	accountAddResponse, err := client.AddAccount(
		email,
		d.Get("ref_id").(string),
		d.Get("user_name").(string),
		d.Get("plan_id").(string),
		d.Get("account_name").(string),
		d.Get("log_level").(string),
		d.Get("logs_account_id").(int),
		d.Get("parent_id").(int),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula account for email: %s, %s\n", email, err)
		return err
	}

	// Set the Account ID
	d.SetId(strconv.Itoa(accountAddResponse.Account.AccountID))
	log.Printf("[INFO] Created Incapsula account for email: %s\n", email)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	err = updateAdditionalAccountProperties(client, d)
	if err != nil {
		return err
	}

	err = updateDefaultDataStorageRegion(client, d)
	if err != nil {
		return err
	}

	// Set the rest of the state from the resource read
	return resourceAccountRead(d, m)
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Reading Incapsula account for Account ID: %d\n", accountID)

	accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)

	// Account object may have been deleted
	if accountStatusResponse != nil && accountStatusResponse.Res.(float64) == 9403 {
		log.Printf("[INFO] Incapsula Account ID %d has already been deleted: %s\n", accountID, err)
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account for Account ID: %d, %s\n", accountID, err)
		return err
	}

	d.Set("parent_id", accountStatusResponse.Account.ParentID)
	d.Set("email", accountStatusResponse.Account.Email)
	d.Set("plan_id", accountStatusResponse.Account.PlanID)
	d.Set("plan_name", accountStatusResponse.Account.PlanName)
	d.Set("trial_end_date", accountStatusResponse.Account.TrialEndDate)
	d.Set("account_id", accountStatusResponse.Account.AccountID)
	d.Set("ref_id", accountStatusResponse.Account.RefID)
	d.Set("user_name", accountStatusResponse.Account.UserName)
	d.Set("account_name", accountStatusResponse.Account.AccountName)
	d.Set("support_level", accountStatusResponse.Account.SupportLevel)
	d.Set("support_all_tls_versions", accountStatusResponse.Account.SupportAllTLSVersions)
	d.Set("wildcard_san_for_new_sites", accountStatusResponse.Account.WildcardSANForNewSites)
	d.Set("naked_domain_san_for_new_www_sites", accountStatusResponse.Account.NakedDomainSANForNewWWWSites)

	// Get the performance settings for the site
	defaultAccountDataStorageRegion, err := client.GetAccountDataStorageRegion(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula default data storage region for account id: %d, %s\n", accountID, err)
		return err
	}
	d.Set("data_storage_region", defaultAccountDataStorageRegion.Region)

	log.Printf("[INFO] Finished reading Incapsula account for account ud: %d\n", accountID)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	updateParams := [7]string{"email", "plan_id", "ref_id", "error_page_template", "support_all_tls_versions", "naked_domain_san_for_new_www_sites", "wildcard_san_for_new_sites"}
	for i := 0; i < len(updateParams); i++ {
		param := updateParams[i]
		if d.HasChange(param) && d.Get(param) != "" {
			log.Printf("[INFO] Updating Incapsula account param (%s) with value (%s) for account_id: %s\n", param, d.Get(param).(string), d.Id())
			_, err := client.UpdateAccount(d.Id(), param, d.Get(param).(string))
			if err != nil {
				log.Printf("[ERROR] Could not update Incapsula account param (%s) with value (%s) for account_id: %s %s\n", param, d.Get(param).(string), d.Id(), err)
				return err
			}
		}
	}

	err := updateAdditionalAccountProperties(client, d)
	if err != nil {
		return err
	}

	err = updateAccountLogLevel(client, d)
	if err != nil {
		return err
	}

	err = updateDefaultDataStorageRegion(client, d)
	if err != nil {
		return err
	}

	// Set the rest of the state from the resource read
	return resourceAccountRead(d, m)
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Deleting Incapsula account id: %d\n", accountID)

	err := client.DeleteAccount(accountID)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula account id: %d, %s\n", accountID, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula account id: %d\n", accountID)

	return nil
}

func updateAdditionalAccountProperties(client *Client, d *schema.ResourceData) error {
	updateParams := [5]string{"name", "error_page_template", "support_all_tls_versions", "naked_domain_san_for_new_www_sites", "wildcard_san_for_new_sites"}
	for i := 0; i < len(updateParams); i++ {
		param := updateParams[i]
		if d.HasChange(param) && d.Get(param) != "" {
			log.Printf("[INFO] Updating Incapsula account param (%s) with value (%s) for account_id: %s\n", param, d.Get(param).(string), d.Id())
			_, err := client.UpdateAccount(d.Id(), param, d.Get(param).(string))
			if err != nil {
				log.Printf("[ERROR] Could not update Incapsula account param (%s) with value (%s) for account_id: %s %s\n", param, d.Get(param).(string), d.Id(), err)
				return err
			}
		}
	}
	return nil
}

func updateAccountLogLevel(client *Client, d *schema.ResourceData) error {
	if d.HasChange("log_level") ||
		d.HasChange("logs_account_id") {
		logLevel := d.Get("log_level").(string)
		logsAccountId := d.Get("logs_account_id").(string)
		err := client.UpdateLogLevel(d.Id(), logLevel, logsAccountId)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula account log level: %s and log account id: %s for account_id: %s %s\n", logLevel, logsAccountId, d.Id(), err)
			return err
		}
	}
	return nil
}

func updateDefaultDataStorageRegion(client *Client, d *schema.ResourceData) error {
	if d.HasChange("data_storage_region") {
		region := d.Get("data_storage_region").(string)
		_, err := client.UpdateAccountDataStorageRegion(d.Id(), region)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula account default data storage region: %s for account_id: %s %s\n", region, d.Id(), err)
			return err
		}
	}
	return nil
}
