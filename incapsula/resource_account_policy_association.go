package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
)

const WAF_RULES = "WAF_RULES"
const WEBSITE = "WEBSITE"

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
				Description: "The policy name.",
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
				Description: "This list is currently relevant to whitelist and acl policies. More than one policy can be set as default.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccountPolicyAssociationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountIDStr := d.Get("account_id").(string)
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Printf("[ERROR] Could not convert Account ID. Error: is not numeric: %s", accountIDStr)
		return err
	}

	_, ok := d.GetOk("default_waf_policy_id")
	if ok {
		//add WAF
		wafPolicyIDStr := d.Get("default_waf_policy_id").(string)

		if err != nil {
			log.Printf("[ERROR] Could not convert Waf Policy ID. Error: is not numeric: %s", accountIDStr)
			return err
		}

		//to update WAF policy
		policyGetResponse, err := client.GetPolicy(wafPolicyIDStr)
		if err != nil {
			log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", wafPolicyIDStr, err)
			return err
		}

		if policyGetResponse.Value.PolicyType != WAF_RULES {
			log.Printf("[ERROR] Cannot set a policy of type %s as a default WAF Policy. Policy ID: %d", policyGetResponse.Value.PolicyType, policyGetResponse.Value.ID)
			return fmt.Errorf("Cannot set a policy of type %s as a default WAF Policy. Policy ID: %d", policyGetResponse.Value.PolicyType, policyGetResponse.Value.ID)
		}

		err = updatePolicy(policyGetResponse.Value, accountID, *client)
		if err != nil {
			log.Printf("[ERROR] Could not update Default WAF Policy ID %s for Account ID %d. Reason: %s", wafPolicyIDStr, accountID, err.Error())
			resourceAccountPolicyAssociationRead(d, m)
			return err
		}

	}

	err = updateNonMandatoryPolicies(d.Get("default_non_mandatory_policy_ids").(*schema.Set).List(), accountID, *client)
	if err != nil {
		log.Printf("[ERROR] Could not update Default Non Mandatory Policies. Reason: %s", err.Error())
		resourceAccountPolicyAssociationRead(d, m)
		return err
	}
	d.SetId(accountIDStr)
	return resourceAccountPolicyAssociationRead(d, m)
}

func resourceAccountPolicyAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountIDStr := d.Get("account_id").(string)
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Printf("[ERROR] Could not convert Account ID. Error: is not numeric: %s", accountIDStr)
		return err
	}
	//for policyId of non default
	nonMandatoryPolicyIdList := (d.Get("default_non_mandatory_policy_ids").(*schema.Set).List())
	for _, policy := range nonMandatoryPolicyIdList {
		policyIdStr := fmt.Sprint(policy)
		policyGetResponse, err := client.GetPolicy(policyIdStr)
		if err != nil {
			log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", policyIdStr, err)
			return err
		}
		err = removePolicy(policyGetResponse.Value, accountID, *client)
		if err != nil {
			log.Printf("[ERROR] Could not remove Default Non Mandatory Policy ID %d for Account ID %d. Reason: %s", policyGetResponse.Value.ID, accountID, err.Error())
			return err
		}
	}
	d.SetId("")
	return nil
}

func resourceAccountPolicyAssociationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	d.Set("account_id", d.Id())
	accountID := d.Id()
	getAllPoliciesResponse, err := client.GetAllPoliciesForAccount(accountID)

	if err != nil {
		log.Printf("[ERROR] Could not get All Incapsula Policies for Account ID: %s - %s\n", accountID, err)
		return err
	}

	defaultWafPolicyId, nonMandatoryPolicies := filterPoliciesByTypeForAccount(accountID, getAllPoliciesResponse)

	if defaultWafPolicyId != "0" {
		d.Set("default_waf_policy_id", defaultWafPolicyId)
	}
	d.Set("default_non_mandatory_policy_ids", nonMandatoryPolicies)
	return nil
}

func filterPoliciesByTypeForAccount(accountId string, getAllPoliciesResponse *[]Policy) (string, *schema.Set) {
	var wafPolicy Policy
	nonMandatoryPoliciesSet := &schema.Set{F: schema.HashString}
	for _, policy := range *getAllPoliciesResponse {
		if containsAccountIdInDefaultPolicyConfig(policy.DefaultPolicyConfig, accountId) {
			if policy.PolicyType == WAF_RULES {
				wafPolicy = policy
			} else {
				nonMandatoryPoliciesSet.Add(strconv.Itoa(policy.ID))
			}
		}
	}
	return strconv.Itoa(wafPolicy.ID), nonMandatoryPoliciesSet
}

func updateNonMandatoryPolicies(policyIds []interface{}, accountID int, client Client) error {
	var policyIdsCopy = make([]interface{}, len(policyIds))
	copy(policyIdsCopy, policyIds)
	log.Printf("%v", policyIdsCopy)

	getAllPoliciesResponse, err := client.GetAllPoliciesForAccount(strconv.Itoa(accountID))
	if err != nil {
		log.Printf("[ERROR] Could not get All Incapsula Policies for Account ID: %d - %s\n", accountID, err)
		return err
	}

	// remove  default policies the are not present in default resource, but present in DB.
	// It means that user wants to remove default status from Policy
	for _, policyFromResponse := range *getAllPoliciesResponse {
		if contains(policyIds, strconv.Itoa(policyFromResponse.ID)) {
			policyIdsCopy = removeElement(policyIdsCopy, strconv.Itoa(policyFromResponse.ID))
		}

		if containsAccountIdInDefaultPolicyConfig(policyFromResponse.DefaultPolicyConfig, strconv.Itoa(accountID)) && policyFromResponse.PolicyType != WAF_RULES {
			if !contains(policyIds, strconv.Itoa(policyFromResponse.ID)) {
				//update policy
				//then policy was ommitted and we need to update policy default list
				//and remove this policy from list of defaults
				err = removePolicy(policyFromResponse, accountID, client)
				if err != nil {
					log.Printf("[ERROR] Could not remove Default Non Mandatory Policy ID %d for Account ID %d. Reason: %s", policyFromResponse.ID, accountID, err.Error())
					return err
				}
			}
		} else {

			if contains(policyIds, strconv.Itoa(policyFromResponse.ID)) {
				if policyFromResponse.PolicyType == WAF_RULES {
					log.Printf("[ERROR] Cannot set a policy of type %s as a default non mandatory Policy. Policy ID: %d", policyFromResponse.PolicyType, policyFromResponse.ID)
					return fmt.Errorf("Cannot set a policy of type %s as a default non mandatory Policy. Policy ID: %d", policyFromResponse.PolicyType, policyFromResponse.ID)
				}

				//update policy, add status default, since for now in DB it is not default
				//added associationns
				policyIdStr := strconv.Itoa(policyFromResponse.ID)
				if err != nil {
					log.Printf("[ERROR] Could not convert Non Mandatory Policy ID. Error: is not numeric: %s", policyIdStr)
					return err
				}
				err = updatePolicy(policyFromResponse, accountID, client)
				if err != nil {
					log.Printf("[ERROR] Could not update Default Non Mandatory Policy ID %s for Account ID %d. Reason: %s", policyIdStr, accountID, err.Error())
					return err
				}
			}
		}
	}

	if len(policyIdsCopy) > 0 {
		log.Printf("[ERROR] Non mandatory default policies list contains policies that are not available for this account: %v", policyIdsCopy)
		return fmt.Errorf("[ERROR] Non mandatory default policies list contains policies that are not available for this account: %v", policyIdsCopy)
	}
	return nil
}

func removePolicy(policy Policy, accountId int, client Client) error {
	currentDefaultPolicyConfigList := policy.DefaultPolicyConfig
	var updatedDefaultPolicyConfigList []DefaultPolicyConfig
	for _, defaultPolicyConfig := range currentDefaultPolicyConfigList {
		if defaultPolicyConfig.AccountID != accountId {
			updatedDefaultPolicyConfigList = append(updatedDefaultPolicyConfigList, defaultPolicyConfig)
		}
	}
	return upsertPolicy(policy, updatedDefaultPolicyConfigList, client)
}

func updatePolicy(policy Policy, accountId int, client Client) error {
	currentDefaultPolicyConfigList := policy.DefaultPolicyConfig
	log.Printf("[DEBUG] updatePolicy ID %v for accountID %v currentDefaultPolicyConfigList\n%v", policy.ID, accountId, currentDefaultPolicyConfigList)
	for _, defaultPolicyConfig := range currentDefaultPolicyConfigList {
		if defaultPolicyConfig.AccountID == accountId {
			log.Print("don't need to update policy")
			return nil
		}

	}

	newDefaultConfig := DefaultPolicyConfig{AccountID: accountId, AssetType: WEBSITE}
	updatedDefaultPolicyConfigList := append(policy.DefaultPolicyConfig, newDefaultConfig)
	return upsertPolicy(policy, updatedDefaultPolicyConfigList, client)
}

func upsertPolicy(policy Policy, updatedDefaultPolicyConfigList []DefaultPolicyConfig, client Client) error {
	policyUpserted := PolicySubmitted{
		Name:                policy.Name,
		Description:         policy.Description,
		Enabled:             policy.Enabled,
		AccountID:           policy.AccountID,
		PolicyType:          policy.PolicyType,
		PolicySettings:      policy.PolicySettings,
		DefaultPolicyConfig: updatedDefaultPolicyConfigList,
	}
	_, err := client.UpdatePolicy(policy.ID, &policyUpserted)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula policy: %s - %s\n", policyUpserted.Name, err)
		return err
	}
	return nil
}

func contains(s []interface{}, str string) bool {
	for _, v := range s {
		vToStr := fmt.Sprint(v)
		if vToStr == str {
			return true
		}
	}

	return false
}

func findIndex(slice []interface{}, elementToRemove string) int {
	for index, element := range slice {
		elementString := fmt.Sprint(element)
		if elementString == elementToRemove {
			return index
		}
	}
	return -1 // not found
}

func removeElement(slice []interface{}, elementStr string) []interface{} {
	index := findIndex(slice, elementStr)
	return append(slice[:index], slice[index+1:]...)
}

func containsAccountIdInDefaultPolicyConfig(defaultPolicyConfigs []DefaultPolicyConfig, accountID string) bool {
	for _, defaultPolicy := range defaultPolicyConfigs {
		if strconv.Itoa(defaultPolicy.AccountID) == accountID {
			return true
		}
	}
	return false
}
