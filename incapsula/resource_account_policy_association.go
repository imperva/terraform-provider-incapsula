package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"sort"
	"strconv"
	"strings"
)

const noAvailablePoliciesConst = "NO_AVAILABLE_POLICIES"

func resourceAccountPolicyAssociation() *schema.Resource {

	return &schema.Resource{
		Create: resourceAccountPolicyAssociationUpdate,
		Read:   resourceAccountPolicyAssociationRead,
		Update: resourceAccountPolicyAssociationUpdate,
		Delete: resourceAccountPolicyAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account Id.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"default_waf_policy_id": {
				Description: "The WAF policy which is set as default to the account. The account can only have 1 such id." +
					"\n The Default policy will be applied automatically to sites that were create after setting it to default.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"default_non_mandatory_policy_ids": {
				Description: "This list is currently relevant to whitelist and acl policies. More than one policy can be set as default. " +
					"providing an empty list or omitting this argument will clear all the non mandatory default policies.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"available_policy_ids": {
				Description: "Comma separated list of The account’s available policies." +
					" These policies can be applied to the websites in the account." +
					" e.g. available_policy_ids = format(\"%s,%s\", incapsula_policy.acl1-policy.id, incapsula_policy.waf3-policy.id)" +
					" Specify this argument only if you are a parent account trying to update your child account policies availability" +
					" in order to remove availability for all policies please specify \"" + noAvailablePoliciesConst + "\".",
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, oldValue string, newValue string, d *schema.ResourceData) bool {
					if newValue == "" {
						// means that the value was not set and doesn't need to be handled'
						return true
					}
					if oldValue == "" && newValue == noAvailablePoliciesConst {
						//means that we have no available policies
						return true
					}
					if oldValue != "" && newValue != "" {
						return suppressEquivalentStringDiffs(k, oldValue, newValue, d)
					}
					return false
				},
				Optional: true,
			},
		},
	}
}

func resourceAccountPolicyAssociationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountID := d.Get("account_id").(string)
	var err error

	//init default waf policy
	var wafPolicyIDStr string
	if d.Get("default_waf_policy_id") != nil {
		wafPolicyIDStr = d.Get("default_waf_policy_id").(string)
	}

	//init default non mandatory non distinct policies
	defaultNonMandatoryPolicyIds := make([]int, 0)
	defaultNonMandatoryPolicyIdsInter := d.Get("default_non_mandatory_policy_ids").(*schema.Set).List()
	if !IsSetNil(d.Get("default_non_mandatory_policy_ids").(*schema.Set)) {
		defaultNonMandatoryPolicyIds, err = ListToIntSlice(defaultNonMandatoryPolicyIdsInter)
		if err != nil {
			return fmt.Errorf("Cannot convert default_non_mandatory_policy_ids to integer")
		}
	}

	//init available policies
	var availablePolicyIds []int
	availablePolicies := d.Get("available_policy_ids").(string)
	if availablePolicies == noAvailablePoliciesConst {
		availablePolicyIds = make([]int, 0)
	} else if availablePolicies != "" {
		if client.accountStatus.isSubAccount() {
			return fmt.Errorf("sub accounts cannot change thier available_policy_ids")
		}
		splitPoliciesIds := strings.Split(availablePolicies, ",")
		availablePolicyIds, err = ToIntSlice(splitPoliciesIds)
		if err != nil {
			return fmt.Errorf("Cannot convert available_policy_ids to integer")
		}
	}

	_, err = client.PatchAccountPolicyAssociation(accountID, availablePolicyIds, defaultNonMandatoryPolicyIds, wafPolicyIDStr)
	if err != nil {
		return err
	}
	d.SetId(accountID)
	return resourceAccountPolicyAssociationRead(d, m)
}

func IsSetNil(setToCheck *schema.Set) bool {
	emptyStringSet := make(map[string]interface{}, 0)
	if setToCheck.Len() == 0 && !setToCheck.Equal(emptyStringSet) {
		return true
	}
	return false
}

func ListToIntSlice(list []interface{}) ([]int, error) {
	var intArray []int
	for _, val := range list {
		intVar, err := strconv.Atoi(val.(string))
		if err != nil {
			log.Printf("[ERROR] Cannot convert list args to integer. arg: %s\n", val)
			return nil, fmt.Errorf("Cannot convert list args to integer")
		}
		intArray = append(intArray, intVar)
	}
	return intArray, nil
}

func ToIntSlice(list []string) ([]int, error) {
	var intArray []int
	for _, val := range list {
		intVar, err := strconv.Atoi(val)
		if err != nil {
			log.Printf("[ERROR] Cannot convert list of strings to integer. arg: %s\n", val)
			return nil, fmt.Errorf("Cannot convert list of strings to integer")
		}
		intArray = append(intArray, intVar)
	}
	return intArray, nil
}

func ToStringSlice(intArray []int) []string {
	stringArray := make([]string, 0)
	for _, val := range intArray {
		strVar := strconv.Itoa(val)
		stringArray = append(stringArray, strVar)
	}
	return stringArray
}

func resourceAccountPolicyAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountIDStr := d.Get("account_id").(string)
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Printf("[ERROR] Could not convert Account ID. Error: is not numeric: %s", accountIDStr)
		return err
	}
	wafPolicyIdStr := d.Get("default_waf_policy_id").(string)
	defaultNonMandatoryPolicyIds := make([]int, 0)
	var availablePolicyIds []int
	if wafPolicyIdStr != "" && !client.accountStatus.isSubAccount() && client.accountStatus.AccountID != accountID {
		wafPolicyID, err := strconv.Atoi(wafPolicyIdStr)
		if err != nil {
			log.Printf("[ERROR] Could not convert WAF Rule Policy ID. Error: is not numeric: %s", wafPolicyIdStr)
			return err
		}
		availablePolicyIds = append(availablePolicyIds, wafPolicyID)
	}
	_, err = client.PatchAccountPolicyAssociation(accountIDStr, availablePolicyIds, defaultNonMandatoryPolicyIds, wafPolicyIdStr)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAccountPolicyAssociationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := d.Set("account_id", d.Id())
	if err != nil {
		return err
	}

	accountID := d.Id()
	getAccountPolicyAssociation, err := client.GetAccountPolicyAssociation(accountID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula Policies Association for Account ID: %s - %s\n", accountID, err)
		return err
	}

	defaultWafPolicyId := getAccountPolicyAssociation.DefaultWafPolicyId
	nonMandatoryPolicies := getAccountPolicyAssociation.DefaultNonMandatoryNonDistinctPolicyIds
	availablePolicyIds := getAccountPolicyAssociation.AvailablePolicyIds

	if d.Get("default_waf_policy_id").(string) != "" && defaultWafPolicyId != 0 {
		err = d.Set("default_waf_policy_id", strconv.Itoa(defaultWafPolicyId))
		if err != nil {
			return err
		}
	}
	err = d.Set("default_non_mandatory_policy_ids", ToStringSlice(nonMandatoryPolicies))
	if err != nil {
		return err
	}

	sort.Slice(availablePolicyIds, func(i, j int) bool {
		return availablePolicyIds[i] < availablePolicyIds[j]
	})

	err = d.Set("available_policy_ids", strings.Join(ToStringSlice(availablePolicyIds), ","))
	if err != nil {
		return err
	}

	return nil
}
