package incapsula

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
			"default_website_accounts": {
				Description: "The list of account IDs that current policy is default for. I.e. the policy will be applied for all future added assets in these accounts.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"available_for_accounts": {
				Description: "The list of account IDs that current policy is available for. If parameter equals empty list (\"[]\") then current policy is available for all subaccounts.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policySettingsString := d.Get("policy_settings").(string)
	var policySettings []PolicySetting
	err := json.Unmarshal([]byte(policySettingsString), &policySettings)

	defaultPolicyConfig, err := getDefaultPolicyConfigForRequest(d, true)
	if err != nil {
		return fmt.Errorf("Failed to create policy %s, error: %s", d.Get("name").(string), err.Error())
	}
	policySubmitted := PolicySubmitted{
		Name:                d.Get("name").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyType:          d.Get("policy_type").(string),
		Description:         d.Get("description").(string),
		AccountID:           d.Get("account_id").(int),
		PolicySettings:      policySettings,
		DefaultPolicyConfig: defaultPolicyConfig,
	}

	policyAddResponse, err := client.AddPolicy(&policySubmitted)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}

	policyID := strconv.Itoa(policyAddResponse.Value.ID)
	// Set the policyID
	// We set ID here since if policy was created then we will start managing the resource in TF, even if Update Policy Account Assoociation action will fail.

	d.SetId(policyID)

	associatedAccountsList, err := getAccountAssociationListForRequest(d.Get("available_for_accounts").(*schema.Set).List())
	policyAccountAssociation, err := client.UpdatePolicyAccountAssociation(policyID, associatedAccountsList)
	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}
	d.Set("available_for_accounts", getAccountAssociationListForSchema(policyAccountAssociation))

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

	policyAccountAssociation, err := client.GetPolicyAccountAssociation(policyID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", policyID, err)
		return err
	}
	d.Set("available_for_accounts", getAccountAssociationListForSchema(policyAccountAssociation))

	// Set computed values
	d.Set("name", policyGetResponse.Value.Name)
	d.Set("enabled", policyGetResponse.Value.Enabled)
	d.Set("policy_type", policyGetResponse.Value.PolicyType)
	d.Set("description", policyGetResponse.Value.Description)
	d.Set("account_id", policyGetResponse.Value.AccountID)
	d.Set("default_website_accounts", getDefaultPolicyConfigAccountListForSchema(policyGetResponse.Value.DefaultPolicyConfig))

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

	defaultPolicyConfig, err := getDefaultPolicyConfigForRequest(d, false)
	if err != nil {
		return fmt.Errorf("Failed to update policy %s, error: %s", d.Get("name").(string), err.Error())
	}

	policySubmitted := PolicySubmitted{
		Name:                d.Get("name").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyType:          d.Get("policy_type").(string),
		AccountID:           d.Get("account_id").(int),
		Description:         d.Get("description").(string),
		PolicySettings:      policySettings,
		DefaultPolicyConfig: defaultPolicyConfig,
	}

	_, err = client.UpdatePolicy(id, &policySubmitted)

	if err != nil {
		return err
	}

	associatedAccountsList, err := getAccountAssociationListForRequest(d.Get("available_for_accounts").(*schema.Set).List())
	policyAccountAssociation, err := client.UpdatePolicyAccountAssociation(d.Id(), associatedAccountsList)
	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula policy: %s - %s\n", policySubmitted.Name, err)
		return err
	}
	d.Set("available_for_accounts", getAccountAssociationListForSchema(policyAccountAssociation))
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

func getAccountAssociationListForSchema(policyAccountAssociation *PolicyAccountAssociation) *schema.Set {
	associatedAccounts := &schema.Set{F: schema.HashString}
	for i := range policyAccountAssociation.Value {
		associatedAccounts.Add(strconv.Itoa(policyAccountAssociation.Value[i]))
	}
	return associatedAccounts
}

func getDefaultPolicyConfigAccountListForSchema(defaultPolicies []DefaultPolicyConfig) *schema.Set {
	accountIdSet := &schema.Set{F: schema.HashString}
	for i := range defaultPolicies {
		log.Printf("Adding %s - %d to schema,", strconv.Itoa(defaultPolicies[i].AccountID), defaultPolicies[i].AccountID)
		accountIdSet.Add(strconv.Itoa(defaultPolicies[i].AccountID))
	}
	return accountIdSet
}

//convert account IDs from string format ti int
func getAccountAssociationListForRequest(policyAccountAssoiationSchema []interface{}) ([]int, error) {
	var accountIdList []int
	for _, account := range policyAccountAssoiationSchema {
		accountIdInt, err := strconv.Atoi(fmt.Sprint(account))
		if err != nil {
			return nil, err
		}
		accountIdList = append(accountIdList, accountIdInt)
	}
	return accountIdList, nil
}

func getDefaultPolicyConfigForRequest(d *schema.ResourceData, isCreate bool) ([]DefaultPolicyConfig, error) {
	defaultAccountIdsFromSchema := d.Get("default_website_accounts").(*schema.Set).List()

	defaultWebsiteAccountsSForRequest := make([]string, len(defaultAccountIdsFromSchema))
	for i, v := range defaultAccountIdsFromSchema {
		defaultWebsiteAccountsSForRequest[i] = fmt.Sprint(v)
	}

	var defaultWebsiteAccountList []DefaultPolicyConfig

	for _, entry := range defaultWebsiteAccountsSForRequest {
		accountID, err := strconv.Atoi(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert default policy Website Account ID %s. Reason: is not numeric", entry)
		}

		defaultWebsiteAccount := DefaultPolicyConfig{
			AccountID: accountID,
			AssetType: "WEBSITE",
		}
		//if !isCreate {
		//	policyId, err := strconv.Atoi(d.Id())
		//	if err != nil {
		//		return nil, fmt.Errorf("Failed to convert Policy ID, error: %s", err.Error())
		//	}
		//	defaultWebsiteAccount.PolicyID = policyId
		//}
		defaultWebsiteAccountList = append(defaultWebsiteAccountList, defaultWebsiteAccount)
	}
	return defaultWebsiteAccountList, nil
}
