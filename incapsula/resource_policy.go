package incapsula

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description:      "The policy settings as JSON string. See Imperva documentation for help with constructing a correct value.",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentJSONStringDiffs,
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
	policyGetResponse, err := client.GetPolicy(policyID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", policyID, err)
		return err
	}

	// Set computed values
	d.Set("name", policyGetResponse.Value.Name)
	d.Set("enabled", policyGetResponse.Value.Enabled)
	d.Set("policy_type", policyGetResponse.Value.PolicyType)
	d.Set("description", policyGetResponse.Value.Description)
	d.Set("account_id", policyGetResponse.Value.AccountID)

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

	policySettingsString := d.Get("policy_settings").(string)
	var policySettings []PolicySetting
	err = json.Unmarshal([]byte(policySettingsString), &policySettings)

	policyGetResponse, err := client.GetPolicy(d.Id())
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
		AccountID:           d.Get("account_id").(int),
		Description:         d.Get("description").(string),
		DefaultPolicyConfig: policyGetResponse.Value.DefaultPolicyConfig,
		PolicySettings:      policySettings,
	}

	_, err = client.UpdatePolicy(id, &policySubmitted)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}

	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.DeletePolicy(d.Id())

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
