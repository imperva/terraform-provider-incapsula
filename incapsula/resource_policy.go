package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"name": {
				Description: "The policy name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Enables the policy.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"policy_type": {
				Description: "The policy type. Possible values: ACL, WHITELIST, WAF_RULES",
				Type:        schema.TypeString,
				Required:    true,
			},
			"policy_settings": {
				Description: "The policy settings as JSON string. See Imperva documentation for help with constructing a correct value.",
				Type:        schema.TypeString,
				Required:    true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					oldValue = strings.ReplaceAll(oldValue, " ", "")
					oldValue = strings.ReplaceAll(oldValue, "\n", "")
					oldValue = strings.ReplaceAll(oldValue, "\"data\":{}", "")
					oldValue = strings.ReplaceAll(oldValue, "\"policyDataExceptions\":[]", "")
					oldValue = strings.ReplaceAll(oldValue, ",,", ",")
					oldValue = strings.ReplaceAll(oldValue, ",}", "}")
					oldValue = strings.ReplaceAll(oldValue, "{,", "{")

					newValue = strings.ReplaceAll(newValue, " ", "")
					newValue = strings.ReplaceAll(newValue, "\n", "")
					newValue = strings.ReplaceAll(newValue, "\"data\":{}", "")
					newValue = strings.ReplaceAll(newValue, "\"policyDataExceptions\":[]", "")
					newValue = strings.ReplaceAll(newValue, ",,", ",")
					newValue = strings.ReplaceAll(newValue, ",}", "}")
					newValue = strings.ReplaceAll(newValue, "{,", "{")

					return suppressEquivalentJSONStringDiffs(k, oldValue, newValue, d)
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					// Check if valid JSON
					d := val.(string)
					var js interface{}
					unMarshalErr := json.Unmarshal([]byte(d), &js)
					if unMarshalErr != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid JSON policy, please check your syntax, got: %s, message: %s", key, d, unMarshalErr))
					}
					return
				},
			},
			// Optional Arguments
			"account_id": {
				Description: "The Account ID of the policy.",
				Type:        schema.TypeInt,
				Deprecated:  "Use the incapsula_account_policy_association resource instead",
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "The policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func getCurrentAccountId(d *schema.ResourceData, accountStatus *AccountStatusResponse) *int {
	caid := d.Get("account_id").(int)
	if accountStatus.isSubAccount() || caid == 0 {
		//in case of sub account we do not want to send the caid since the policy owner is the sub account's parent
		return nil
	}
	return &caid
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policySettingsString := d.Get("policy_settings").(string)
	var policySettings []PolicySetting
	err := json.Unmarshal([]byte(policySettingsString), &policySettings)

	policySubmitted := PolicySubmitted{
		Name:           d.Get("name").(string),
		Enabled:        d.Get("enabled").(bool),
		PolicyType:     d.Get("policy_type").(string),
		Description:    d.Get("description").(string),
		AccountID:      d.Get("account_id").(int),
		PolicySettings: policySettings,
	}

	policyAddResponse, err := client.AddPolicy(&policySubmitted)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}

	policyID := strconv.Itoa(policyAddResponse.Value.ID)

	d.SetId(policyID)
	log.Printf("[INFO] Created Incapsula policy with ID: %s\n", policyID)
	return resourcePolicyRead(d, m)
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Id()

	currentAccountId := getCurrentAccountId(d, client.accountStatus)

	policyGetResponse, err := client.GetPolicy(policyID, currentAccountId)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", policyID, err)
		return err
	}

	// Set computed values
	d.Set("name", policyGetResponse.Value.Name)
	d.Set("enabled", policyGetResponse.Value.Enabled)
	d.Set("policy_type", policyGetResponse.Value.PolicyType)
	d.Set("description", policyGetResponse.Value.Description)
	if d.Get("account_id") != nil {
		d.Set("account_id", d.Get("account_id"))
	} else {
		d.Set("account_id", policyGetResponse.Value.AccountID)
	}

	// JSON encode policy settings
	policySettingsJSONBytes, err := json.MarshalIndent(policyGetResponse.Value.PolicySettings, "", "    ")
	if err != nil {
		log.Printf("[ERROR] Could not get marshal Incapsula policy settings: %s - %s - %s\n", policyID, err, policySettingsJSONBytes)
		return err
	}
	d.Set("policy_settings", string(policySettingsJSONBytes))

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if d.Get("account_id") != nil {
		log.Printf("[WARN] Incapsula policy account id attribute is deprecated - please remove it\n")
	}
	policySettingsString := d.Get("policy_settings").(string)
	var policySettings []PolicySetting
	err = json.Unmarshal([]byte(policySettingsString), &policySettings)

	currentAccountId := getCurrentAccountId(d, client.accountStatus)
	policyGetResponse, err := client.GetPolicy(d.Id(), currentAccountId)
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula policy: %d - %s\n", id, err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Incapsula policy ID %d has already been deleted: %s\n", id, err)
			d.SetId("")
			return nil
		}
		return err
	}

	policySubmitted := PolicySubmitted{
		Name:                d.Get("name").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyType:          d.Get("policy_type").(string),
		Description:         d.Get("description").(string),
		DefaultPolicyConfig: policyGetResponse.Value.DefaultPolicyConfig,
		PolicySettings:      policySettings,
	}

	_, err = client.UpdatePolicy(id, &policySubmitted, currentAccountId)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}

	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	currentAccountId := getCurrentAccountId(d, client.accountStatus)
	err := client.DeletePolicy(d.Id(), currentAccountId)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
