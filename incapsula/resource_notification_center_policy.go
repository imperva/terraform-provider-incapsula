package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

var assetResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		"asset_type": {
			Type:        schema.TypeString,
			Description: "The asset type",
			Required:    true,
		},

		"asset_id": {
			Type:        schema.TypeInt,
			Description: "The asset id",
			Optional:    true,
		},
	},
}

func resourceNotificationCenterPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationCenterPolicyCreate,
		Read:   resourceNotificationCenterPolicyRead,
		Update: resourceNotificationCenterPolicyUpdate,
		Delete: resourceNotificationCenterPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				isOnlyPolicyIdFormat := len(idSlice) == 1
				isAccountIdSlashPolicyIdFormat := !(len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "")

				if !isOnlyPolicyIdFormat && !isAccountIdSlashPolicyIdFormat {
					return nil, fmt.Errorf("unexpected format of NotificationCenterPolicy import input (%q), expected accountId/policyId or policyId only", d.Id())
				}

				policyIdIndex := 0
				if isAccountIdSlashPolicyIdFormat {
					policyIdIndex = 1
				}

				policyIdString := idSlice[policyIdIndex]
				policyId, err := strconv.Atoi(policyIdString)
				if err != nil {
					fmt.Errorf("NotificationCenterPolicy- failed to convert policy Id from import command, actual value: %s, expected numeric id", policyIdString)
				}

				d.SetId(policyIdString)
				if isAccountIdSlashPolicyIdFormat {
					accountId, err := strconv.Atoi(idSlice[0])
					if err != nil {
						fmt.Errorf("NotificationCenterPolicy- failed to account Id from import command, actual value: %s, expected numeric id", idSlice[0])
					}
					d.Set("account_id", accountId)
					log.Printf("[DEBUG] Import NotificationCenterPolicy JSON for account Id: %d, policy Id: %d", accountId, policyId)
				} else {
					log.Printf("[DEBUG] Import NotificationCenterPolicy JSON for policy Id: %d", policyId)
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account ID",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"policy_name": {
				Description: "The name of the policy",
				Type:        schema.TypeString,
				Required:    true,
			},
			"status": {
				Description:  "Indicates whether policy is enabled or disabled. Default value is ENABLE",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLE",
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
			},
			"sub_category": {
				Description: "Subtype of notification policy. Example values include: ‘account_notifications’; " +
					"‘website_notifications’; ‘certificate_management_notifications’",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"emailchannel_user_recipient_list": {
				Description: "List of Imperva users id to get the notifications",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"emailchannel_external_recipient_list": {
				Description: "List of external email to get the notifications (not Imperva users)",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},

			"asset": {
				Description: "Assets to receive notifications (if assets are relevant to the sub category type). " +
					"\nObject struct:\nassetType: the asset type. Example: websites, router connections, network prefixes, " +
					"individual IPs, Flow exporters\nassetId: the asset id.\n",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &assetResource,
				Set:      schema.HashResource(&assetResource),
			},

			"apply_to_new_assets": {
				Description: "If value is ‘TRUE’, all newly onboarded assets are automatically added to the " +
					"notification policy's assets list.\nDefault value is no\n",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TRUE", "FALSE"}, false),
				Default:      "FALSE",
			},
			"policy_type": {
				Description: "If value is ‘ACCOUNT’, the policy will apply only to the current account. \nIf the value" +
					" is 'SUB_ACCOUNT' the policy applies to the sub accounts only. \n The parent account will receive " +
					"notifications for activity in the sub accounts that are specified in the subAccountList parameter." +
					"\nThis parameter is available only in accounts that can contain sub accounts.\n",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ACCOUNT",
			},
			"apply_to_new_sub_accounts": {
				Description: "If value is ‘TRUE’, all newly onboarded sub accounts are automatically added to the " +
					"notification policy's sub account list.\nDefault value is no\n",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TRUE", "FALSE"}, false),
				Default:      "FALSE",
			},
			"sub_account_list": {
				Description: "The policy ID. During update must be equal to the updated policy ID.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
		},
	}
}

func resourceNotificationCenterPolicyUpdate(data *schema.ResourceData, i interface{}) error {
	client := i.(*Client)
	notificationCenterPolicyName := data.Get("policy_name").(string)
	notificationCenterPolicyId, _ := getPolicyId(data)
	accountId := data.Get("account_id").(int)
	log.Printf("[INFO] Updateding NotificationCenterPolicy with policyId:%d accountId:%d and name: %s\n",
		notificationCenterPolicyId, accountId, notificationCenterPolicyName)
	notificationPolicyFullDto := getNotificationCenterPolicyFromResource(data)
	notificationCenterPolicyUpdateResponse, err := client.UpdateNotificationCenterPolicy(&notificationPolicyFullDto)
	if err != nil {
		log.Printf("[ERROR] Could not update NotificationCenterPolicy id:%d. \nThe policy: %+v  \nThe response:%+v \nThe error: %s\n",
			notificationCenterPolicyId, notificationPolicyFullDto, notificationCenterPolicyUpdateResponse, err)
		return err
	} else {
		log.Printf("[DEBUG] NotificationCenter update policy with json reponse: %+v ", notificationCenterPolicyUpdateResponse)
	}

	return resourceNotificationCenterPolicyRead(data, client)
}

func resourceNotificationCenterPolicyCreate(data *schema.ResourceData, i interface{}) error {
	client := i.(*Client)
	notificationCenterPolicyName := data.Get("policy_name").(string)
	log.Printf("[INFO] Creating NotificationCenterPolicy: %s\n", notificationCenterPolicyName)
	notificationPolicyFullDto := getNotificationCenterPolicyFromResource(data)
	notificationCenterPolicyAddResponse, err := client.AddNotificationCenterPolicy(&notificationPolicyFullDto)

	if err != nil {
		log.Printf("[ERROR] Could not create NotificationCenterPolicy. \nThe policy: %+v  \nThe response:%+v \nThe error: %s\n",
			notificationPolicyFullDto, notificationCenterPolicyAddResponse, err)
		return err
	} else {
		log.Printf("[DEBUG] NotificationCenter create policy with json response: %+v ", notificationCenterPolicyAddResponse)
	}

	policyID := strconv.Itoa(notificationCenterPolicyAddResponse.Data.PolicyId)
	data.SetId(policyID)
	log.Printf("[DEBUG] NotificationCenter create policy with id %s ", policyID)

	return resourceNotificationCenterPolicyRead(data, client)
}

//This function get all the properties of NotificationCenterPolicy from the resource,
//so we can share it with create & update function
func getNotificationCenterPolicyFromResource(data *schema.ResourceData) NotificationPolicyFullDto {
	policyId, _ := getPolicyId(data)
	log.Printf(
		"[INFO] policyId: %d\n"+
			"[INFO] account_id: %d\n"+
			"[INFO] policy_name: %s\n"+
			"[INFO] status: %s\n"+
			"[INFO] sub_category: %s\n"+
			"[INFO] emailchannel_user_recipient_list: %s\n"+
			"[INFO] emailchannel_external_recipient_list: %s\n"+
			"[INFO] asset: %s\n"+
			"[INFO] apply_to_new_assets: %s\n"+
			"[INFO] policy_type: %s\n"+
			"[INFO] apply_to_new_sub_accounts: %s\n"+
			"[INFO] sub_account_list: %s\n",
		policyId,
		data.Get("account_id").(int),
		data.Get("policy_name").(string),
		data.Get("status").(string),
		data.Get("sub_category").(string),
		data.Get("emailchannel_user_recipient_list").(interface{}),
		data.Get("emailchannel_external_recipient_list").(interface{}),
		data.Get("asset").(interface{}),
		data.Get("apply_to_new_assets").(string),
		data.Get("policy_type").(string),
		data.Get("apply_to_new_sub_accounts").(string),
		data.Get("sub_account_list").(interface{}))

	assetList := getAssetsFromResource(data)
	subAccountsDtoList := getSubAccountsDtoListFromResource(data)
	notificationChannelList := getEmailChannelFromResource(data)
	notificationPolicyFullDto := NotificationPolicyFullDto{
		PolicyId:                policyId,
		AccountId:               data.Get("account_id").(int),
		PolicyName:              data.Get("policy_name").(string),
		Status:                  data.Get("status").(string),
		SubCategory:             data.Get("sub_category").(string),
		NotificationChannelList: []NotificationChannelEmailDto{notificationChannelList},
		AssetList:               assetList,
		ApplyToNewAssets:        data.Get("apply_to_new_assets").(string),
		PolicyType:              data.Get("policy_type").(string),
		SubAccountPolicyInfo: SubAccountPolicyInfo{
			ApplyToNewSubAccounts: data.Get("apply_to_new_sub_accounts").(string),
			SubAccountList:        subAccountsDtoList,
		},
	}
	log.Printf("[DEBUG] getNotificationCenterPolicyFromResource build a NotificationPolicyFullDto from the resource file: %+v", notificationPolicyFullDto)

	return notificationPolicyFullDto
}

func getEmailChannelFromResource(data *schema.ResourceData) NotificationChannelEmailDto {
	var userRecipientDto []RecipientDto
	usersIds := data.Get("emailchannel_user_recipient_list").([]interface{})
	for _, userId := range usersIds {
		recipientDto := RecipientDto{
			RecipientType: "User",
			Id:            userId.(int),
		}
		userRecipientDto = append(userRecipientDto, recipientDto)
	}

	externalUsersEmail := data.Get("emailchannel_external_recipient_list").([]interface{})
	for _, userEmail := range externalUsersEmail {
		recipientDto := RecipientDto{
			RecipientType: "External",
			DisplayName:   userEmail.(string),
		}
		userRecipientDto = append(userRecipientDto, recipientDto)
	}

	notificationChannelList := NotificationChannelEmailDto{
		ChannelType:     "email",
		RecipientToList: userRecipientDto,
	}
	return notificationChannelList
}

func getSubAccountsDtoListFromResource(d *schema.ResourceData) []SubAccountDTO {
	subAccountsIds := d.Get("sub_account_list").([]interface{})
	var subAccountsDtoList []SubAccountDTO
	for _, subAccountId := range subAccountsIds {
		subAccountDTO := SubAccountDTO{SubAccountId: subAccountId.(int)}
		subAccountsDtoList = append(subAccountsDtoList, subAccountDTO)
	}
	return subAccountsDtoList
}

func getAssetsFromResource(d *schema.ResourceData) []AssetDto {
	var assetList []AssetDto
	assets := d.Get("asset").(*schema.Set)
	for _, asset := range assets.List() {
		assetResource := asset.(map[string]interface{})
		assetDto := AssetDto{
			AssetType: assetResource["asset_type"].(string),
			AssetId:   assetResource["asset_id"].(int),
		}
		assetList = append(assetList, assetDto)
	}
	return assetList
}

func resourceNotificationCenterPolicyRead(data *schema.ResourceData, i interface{}) error {
	client := i.(*Client)
	policyID, _ := getPolicyId(data)
	accountId := data.Get("account_id").(int)
	notificationCenterPolicy, err := client.GetNotificationCenterPolicy(policyID, accountId)
	log.Printf("[INFO] Reading NotificationCenterPolicy with id %d \nThe policy: %+v", policyID, notificationCenterPolicy)
	if err != nil {
		return err
	}

	if notificationCenterPolicy == nil {
		log.Printf("[INFO] notificationCenterPolicy %s has already been deleted: %s\n", data.Id(), err)
		return nil
	}

	data.Set("account_id", notificationCenterPolicy.Data.AccountId)
	data.Set("policy_name", notificationCenterPolicy.Data.PolicyName)
	data.Set("status", notificationCenterPolicy.Data.Status)
	data.Set("sub_category", notificationCenterPolicy.Data.SubCategory)
	handleEmailChannelRead(data, notificationCenterPolicy)
	handleAssetsRead(data, notificationCenterPolicy)
	data.Set("apply_to_new_assets", notificationCenterPolicy.Data.ApplyToNewAssets)
	data.Set("policy_type", notificationCenterPolicy.Data.PolicyType)
	data.Set("apply_to_new_sub_accounts", notificationCenterPolicy.Data.SubAccountPolicyInfo.ApplyToNewSubAccounts)

	subAccountList := make([]int, 0)
	for _, subAccount := range notificationCenterPolicy.Data.SubAccountPolicyInfo.SubAccountList {
		subAccountList = append(subAccountList, subAccount.SubAccountId)
	}

	data.Set("sub_account_list", subAccountList)
	log.Printf("[INFO] Finished reading notificationCenterPolicy: %s\n", data.Id())

	return nil
}

func handleAssetsRead(data *schema.ResourceData, notificationCenterPolicy *NotificationPolicy) {
	var assets []interface{}
	for _, assetFromServer := range notificationCenterPolicy.Data.AssetList {
		asset := map[string]interface{}{}
		asset["asset_type"] = assetFromServer.AssetType
		asset["asset_id"] = assetFromServer.AssetId
		log.Printf("[DEBUG] Adding asset to assets set: %+v", asset)
		assets = append(assets, asset)
	}
	log.Printf("[DEBUG] Assets set to save: %+v", assets)
	assetSet := schema.NewSet(schema.HashResource(&assetResource), assets)
	data.Set("asset", assetSet)
}

func handleEmailChannelRead(data *schema.ResourceData, notificationCenterPolicy *NotificationPolicy) {
	var emailChannelUserRecipientsList []int
	var emailChannelExternalRecipientsList []string
	for _, channel := range notificationCenterPolicy.Data.NotificationChannelList {
		if channel.ChannelType == "email" {
			for _, recipient := range channel.RecipientToList {
				switch recipient.RecipientType {
				case "External":
					log.Printf("[DEBUG] Adding recipient to external recipients list: %+v", recipient)
					emailChannelExternalRecipientsList = append(emailChannelExternalRecipientsList, recipient.DisplayName)
				case "User":
					log.Printf("[DEBUG] Adding recipient to user recipients list: %+v", recipient)
					emailChannelUserRecipientsList = append(emailChannelUserRecipientsList, recipient.Id)
				}
			}
		}
	}
	log.Printf("[DEBUG] External recipients list to save: %+v", emailChannelUserRecipientsList)
	data.Set("emailchannel_user_recipient_list", emailChannelUserRecipientsList)
	log.Printf("[DEBUG] User recipients list to save: %+v", emailChannelExternalRecipientsList)
	data.Set("emailchannel_external_recipient_list", emailChannelExternalRecipientsList)
}

func getPolicyId(data *schema.ResourceData) (int, error) {
	policyID, err := strconv.Atoi(data.Id())
	if err != nil {
		log.Printf("[ERROR] failed to convert NotificationCenter policy id to int while reading policy, the id is %s ", data.Id())
	}

	return policyID, err
}

func resourceNotificationCenterPolicyDelete(data *schema.ResourceData, i interface{}) error {
	client := i.(*Client)
	policyID, _ := getPolicyId(data)
	accountId := data.Get("account_id").(int)
	log.Printf("[INFO] Deleting NotificationCenterPolicy policyId: %d and accountId: %d", policyID, accountId)
	err := client.DeleteNotificationCenterPolicy(policyID, accountId)

	if err != nil {
		log.Printf("[ERROR] Could not delete NotificationCenterPolicy id: %d, %s", policyID, err)
		return err
	}

	log.Printf("[INFO] Deleted NotificationCenterPolicy id: %d ", policyID)

	return nil
}
