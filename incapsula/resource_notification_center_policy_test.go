package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"strings"
	"testing"
)

const notificationCenterPolicyResourceType = "incapsula_notification_center_policy"
const policy1AccountWithoutAssets = "notification-policy-account-without-assets"
const fullResourceNamePolicy1WithoutAsset = notificationCenterPolicyResourceType + "." + policy1AccountWithoutAssets
const policy2AccountWithoutAssets = "notification-policy-account-without-assets"
const fullResourcePolicy2NameWithAsset = notificationCenterPolicyResourceType + "." + policy2AccountWithoutAssets

// ##############################
// TODO: Currently the test is disable
var accountId = 1

//##############################

func testAccNotificationCenterPolicy_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG] Running test TestAccNotificationCenterPolicy_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNotificationCenterPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: getAccPolicyAccountWithoutAssets(),
				Check: resource.ComposeTestCheckFunc(
					testCheckNotificationCenterPolicyExists(fullResourceNamePolicy1WithoutAsset),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "policy_name", "Terraform acceptance test- policy account without assets"),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "account_id", strconv.Itoa(accountId)),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "status", "ENABLE"),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "sub_category", "ACCOUNT_NOTIFICATIONS"),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "emailchannel_external_recipient_list.0", "john.mcclane@externalemail.com"),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "emailchannel_external_recipient_list.1", "another.exernal.email@gmail.com"),
					resource.TestCheckResourceAttr(fullResourceNamePolicy1WithoutAsset, "policy_type", "ACCOUNT"),
				),
			},
			{
				ResourceName:      fullResourceNamePolicy1WithoutAsset,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateNotificationCenterPolicyId,
			},
		},
	})
}

func testAccNotificationCenterPolicy_WithAsst(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG] Running test TestAccNotificationCenterPolicy_WithAsst")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNotificationCenterPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: getAccPolicyAccountWithAssets(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckNotificationCenterPolicyExists(fullResourcePolicy2NameWithAsset),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "policy_name", "Terraform policy account with assets"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "account_id", strconv.Itoa(accountId)),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "status", "ENABLE"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "sub_category", "SITE_NOTIFICATIONS"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "emailchannel_external_recipient_list.0", "john.mcclane@externalemail.com"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "emailchannel_external_recipient_list.1", "another.exernal.email@gmail.com"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "asset.0.asset_type", "SITE"),
					resource.TestCheckResourceAttrSet(fullResourcePolicy2NameWithAsset, "asset.0.asset_id"),
					resource.TestCheckResourceAttr(fullResourcePolicy2NameWithAsset, "policy_type", "ACCOUNT"),
				),
			},
			{
				ResourceName:      fullResourcePolicy2NameWithAsset,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateNotificationCenterPolicyId,
			},
		},
	})
}

func testAccNotificationCenterPolicyDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, res := range state.RootModule().Resources {
		if res.Type != notificationCenterPolicyResourceType {
			continue
		}

		policyIdStr := res.Primary.ID
		policyId, _ := strconv.Atoi(policyIdStr)
		accountIdStr := res.Primary.Attributes["account_id"]
		accountId, _ := strconv.Atoi(accountIdStr)
		notificationCenterPolicy, err := client.GetNotificationCenterPolicy(policyId, accountId)
		policyDeleted := notificationCenterPolicy == nil && err != nil && strings.Contains(err.Error(), "Policy not found")
		if !policyDeleted {
			log.Printf("[INFO] ****Test**** testAccNotificationCenterPolicyDestroy, policy: %+v error: %s \n", notificationCenterPolicy, err)
			return err
		}
	}

	return nil
}

func testCheckNotificationCenterPolicyExists(fullResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] ****Test**** starting testCheckNotificationCenterPolicyExists")
		res, ok := state.RootModule().Resources[fullResourceName]
		if !ok {
			return fmt.Errorf("NotificationCenterPolicy resource not found : %s", fullResourceName)
		}

		policyIdStr := res.Primary.ID
		policyId, _ := strconv.Atoi(policyIdStr)
		if !ok || policyIdStr == "" {
			return fmt.Errorf("NotificationCenterPolicy Id does not exists, policy id string: %s ", policyIdStr)
		}

		client := testAccProvider.Meta().(*Client)
		accountIdStr := res.Primary.Attributes["account_id"]
		accountId, _ := strconv.Atoi(accountIdStr)
		log.Printf("[INFO] ****Test**** policyId: %d accountId:%d", policyId, accountId)
		notificationCenterPolicy, err := client.GetNotificationCenterPolicy(policyId, accountId)
		if err != nil {
			return err
		}

		found := false
		if notificationCenterPolicy != nil && notificationCenterPolicy.Data.PolicyId == policyId {
			log.Printf("[INFO] NotificationCenterPolicy: : %v\n", notificationCenterPolicy)
			found = true
		}

		if !found {
			return fmt.Errorf("NotificationCenterPolicy %d does not exist", policyId)
		}

		return nil
	}
}

func getAccPolicyAccountWithoutAssets() string {
	return fmt.Sprintf(`
resource "%s" "%s" { 
	account_id = %d
	policy_name = "Terraform acceptance test- policy account without assets"
	status = "ENABLE"
	sub_category = "ACCOUNT_NOTIFICATIONS"
	emailchannel_external_recipient_list=["john.mcclane@externalemail.com", "another.exernal.email@gmail.com"]	
	policy_type = "ACCOUNT"
}
`,
		notificationCenterPolicyResourceType, policy1AccountWithoutAssets, accountId,
	)
}

func testACCStateNotificationCenterPolicyId(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != notificationCenterPolicyResourceType {
			continue
		}

		policyId, err := strconv.Atoi(rs.Primary.ID)
		accountIdStr := rs.Primary.Attributes["account_id"]
		accountId, _ := strconv.Atoi(accountIdStr)
		if err != nil {
			return "", fmt.Errorf("Error parsing Id %s to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d/%d", accountId, policyId), nil
	}

	return "", fmt.Errorf("Error finding policyId Id")
}

func getAccPolicyAccountWithAssets(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(` 
resource "%s" "%s" { 
	account_id = 52159558
	policy_name = "Terraform policy account with assets"
	asset {
		asset_type = "SITE"
		asset_id = incapsula_site.testacc-terraform-site.id
	}	
	status = "ENABLE"
	sub_category = "SITE_NOTIFICATIONS"
	emailchannel_external_recipient_list=["john.mcclane@externalemail.com", "another.exernal.email@gmail.com"]		
	policy_type = "ACCOUNT"
}
`,
		notificationCenterPolicyResourceType, policy2AccountWithoutAssets,
	)
}
