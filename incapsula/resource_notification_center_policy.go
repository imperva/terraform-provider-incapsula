package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of NotificationCenterPolicy import pair (%q), expected accountId/policyId", d.Id())
				}

				accountId, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("NotificationCenterPolicy- failed to account Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				policyId, err := strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("NotificationCenterPolicy- failed to convert policy Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				d.Set("account_id", accountId)
				d.SetId(idSlice[1])
				log.Printf("[DEBUG] Import NotificationCenterPolicy JSON for account Id: %d, policy Id: %d", accountId, policyId)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Description: "The policy ID. During update must be equal to the updated policy ID.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
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
				ForceNew:    false,
			},
			"status": {
				Description: "Indicates whether policy is enabled or disabled. Default value is enable",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ENABLE",
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
			},

			"apply_to_new_assets": {
				Description: "If value is ‘TRUE’, all newly onboarded assets are automatically added to the " +
					"notification policy's assets list.\nDefault value is no\n",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FALSE",
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FALSE",
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
	notificationCenterPolicyId := data.Get("policy_id").(int)
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
	log.Printf("[INFO] policy_id: %data\n", data.Get("policy_id").(int))
	log.Printf("[INFO] account_id: %d\n", data.Get("account_id").(int))
	log.Printf("[INFO] policy_name: %sata\n", data.Get("policy_name").(string))
	log.Printf("[INFO] status: %s\n", data.Get("status").(string))
	log.Printf("[INFO] sub_category: %s\n", data.Get("sub_category").(string))
	log.Printf("[INFO] emailchannel_user_recipient_list: %s\n", data.Get("emailchannel_user_recipient_list").(interface{}))
	log.Printf("[INFO] emailchannel_external_recipient_list: %s\n", data.Get("emailchannel_external_recipient_list").(interface{}))
	log.Printf("[INFO] asset: %s\n", data.Get("asset").(interface{}))
	log.Printf("[INFO] apply_to_new_assets: %s\n", data.Get("apply_to_new_assets").(string))
	log.Printf("[INFO] policy_type: %s\n", data.Get("policy_type").(string))
	log.Printf("[INFO] apply_to_new_sub_accounts: %s\n", data.Get("apply_to_new_sub_accounts").(string))
	log.Printf("[INFO] sub_account_list: %s\n", data.Get("sub_account_list").(interface{}))

	assetList := getAssetsFromResource(data)
	subAccountsDtoList := getSubAccountsDtoListFromResource(data)
	notificationChannelList := getEmailChannelFromResource(data)
	notificationPolicyFullDto := NotificationPolicyFullDto{
		PolicyId:                data.Get("policy_id").(int),
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

	data.Set("policy_id", notificationCenterPolicy.Data.PolicyId)
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
			for _, recpient := range channel.RecipientToList {
				switch recpient.RecipientType {
				case "External":
					log.Printf("[DEBUG] Adding recipient to external recipients list: %+v", recpient)
					emailChannelExternalRecipientsList = append(emailChannelExternalRecipientsList, recpient.DisplayName)
				case "User":
					log.Printf("[DEBUG] Adding recipient to user recipients list: %+v", recpient)
					emailChannelUserRecipientsList = append(emailChannelUserRecipientsList, recpient.Id)
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
