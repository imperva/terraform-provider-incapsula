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
const fullResourceName = notificationCenterPolicyResourceType + "." + policy1AccountWithoutAssets

var accountId int

func TestAccNotificationCenterPolicy_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG] Running test TestAccNotificationCenterPolicy_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNotificationCenterPolicyDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*Client)
					bla, _ := client.Verify()
					log.Printf("zzzzzz %d", bla.AccountID)
				},
				Config: getAccPolicyAccountWithoutAssets(),
				Check: resource.ComposeTestCheckFunc(
					testCheckNotificationCenterPolicyExists(),
					resource.TestCheckResourceAttr(fullResourceName, "policy_name", "Terraform acceptance test- policy account without assets"),
					resource.TestCheckResourceAttr(fullResourceName, "account_id", strconv.Itoa(accountId)),
					resource.TestCheckResourceAttr(fullResourceName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(fullResourceName, "sub_category", "ACCOUNT_NOTIFICATIONS"),
					resource.TestCheckResourceAttr(fullResourceName, "emailchannel_external_recipient_list.0", "john.mcclane@externalemail.com"),
					resource.TestCheckResourceAttr(fullResourceName, "emailchannel_external_recipient_list.1", "another.exernal.email@gmail.com"),
					resource.TestCheckResourceAttr(fullResourceName, "policy_type", "ACCOUNT"),
				),
			},
			{
				ResourceName:      fullResourceName,
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
		notificationCenterPolicy, err := client.GetNotificationCenterPolicy(policyId)
		policyDeleted := notificationCenterPolicy == nil && err != nil && strings.Contains(err.Error(), "Policy not found")
		if !policyDeleted {
			log.Printf("[INFO] ****Test**** testAccNotificationCenterPolicyDestroy, policy: %+v error: %s \n", notificationCenterPolicy, err)
			return err
		}
	}

	return nil
}

func testCheckNotificationCenterPolicyExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] ****Test**** starting testCheckNotificationCenterPolicyExists")
		//resource := notificationCenterPolicyResourceType + "." + policy1AccountWithoutAssets
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
		log.Printf("[INFO] ****Test**** policyId: %d", policyId)
		notificationCenterPolicy, err := client.GetNotificationCenterPolicy(policyId)
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
	//client := testAccProvider.Meta().(*Client)
	//bla, _ := client.Verify()
	//log.Printf("zzzzzz %d", bla.AccountID)
	accountId = 52159558
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
		if err != nil {
			return "", fmt.Errorf("Error parsing Id %s to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", policyId), nil
	}

	return "", fmt.Errorf("Error finding policyId Id")
}

//TODO: for another test
func getAccPolicyAccountWithAssets() string {
	return fmt.Sprintf(`
resource "%s" "%s" { 
	account_id = 52159558
	policy_name = "Terraform policy account with assets"
	asset {
		asset_type = "SITE"
		asset_id = incapsula_site.tmp-site.id
	}
	asset {
		asset_type = "SITE"
		asset_id = 7999203
	}	
	status = "ENABLE"
	sub_category = "SITE_NOTIFICATIONS"
	emailchannel_external_recipient_list=["john.mcclane@externalemail.com", "another.exernal.email@gmail.com"]		
	policy_type = "ACCOUNT"
}
`,
		notificationCenterPolicyResourceType, subAccountResourceName,
	)
}
