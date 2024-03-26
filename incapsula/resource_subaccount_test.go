package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const subAccountResourceType = "incapsula_subaccount"
const subAccountResourceName = "example_subaccount"
const subAccountName = "acceptance-subaccount-test-1"
const dataStorageRegionName = "APAC"

func TestAccIncapsulaSubAccount_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_txt_settings.go.TestAccIncapsulaSubAccount_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSubAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSubAccountConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSubAccountExists(),
					resource.TestCheckResourceAttr(subAccountResourceType+"."+subAccountResourceName, "sub_account_name", subAccountName),
					resource.TestCheckResourceAttr(subAccountResourceType+"."+subAccountResourceName, "data_storage_region", dataStorageRegionName),
				),
			},
			{
				ResourceName:      subAccountResourceType + "." + subAccountResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSubAccountID,
			},
		},
	})
}

func TestAccIncapsulaSubAccount_HTTP2(t *testing.T) {
	testIncapsulaSubAccountHttp2Client(t, true, true)
}

func TestIncapsulaSubAccount_Http2Defaults(t *testing.T) {
	testIncapsulaSubAccountHttp2Client(t, true, false)
}

func TestIncapsulaSubAccount_Http2ClientAndOriginEnabled(t *testing.T) {
	testIncapsulaSubAccountHttp2Client(t, false, false)
}

func testIncapsulaSubAccountHttp2Client(t *testing.T, enableHttp2ForNewSites bool, enableHttp2ToOriginForNewSites bool) {

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_txt_settings.go.TestAccIncapsulaSubAccount_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSubAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSubAccountWithHttp2Config(enableHttp2ForNewSites, enableHttp2ToOriginForNewSites),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSubAccountExists(),
					resource.TestCheckResourceAttr(subAccountResourceType+"."+subAccountResourceName, "sub_account_name", subAccountName),
					resource.TestCheckResourceAttr(subAccountResourceType+"."+subAccountResourceName, "enable_http2_for_new_sites", strconv.FormatBool(enableHttp2ForNewSites)),
					resource.TestCheckResourceAttr(subAccountResourceType+"."+subAccountResourceName, "enable_http2_to_origin_for_new_sites", strconv.FormatBool(enableHttp2ToOriginForNewSites)),
				),
			},
			{
				ResourceName:      subAccountResourceType + "." + subAccountResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSubAccountID,
			},
		},
	})
}

func getAccIncapsulaSubAccountConfigBasic() string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			sub_account_name = "%s"
			data_storage_region = "%s"
		}`,
		subAccountResourceType, subAccountResourceName, subAccountName, dataStorageRegionName,
	)
}

func getAccIncapsulaSubAccountWithHttp2Config(enableHttp2ForNewSites bool, enableHttp2ToOriginForNewSites bool) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			sub_account_name = "%s"
			enable_http2_for_new_sites = "%t"
			enable_http2_to_origin_for_new_sites = "%t"
		}`,
		subAccountResourceType, subAccountResourceName, subAccountName, enableHttp2ForNewSites, enableHttp2ToOriginForNewSites,
	)
}

func testAccIncapsulaSubAccountDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	var subAccountIDStr string
	for _, res := range state.RootModule().Resources {
		if res.Type != subAccountResourceType {
			continue
		}

		subAccountIDStr = res.Primary.ID
		subAccountID, _ := strconv.Atoi(subAccountIDStr)

		found := true
		subAccount, err := client.AccountStatus(subAccountID, ReadSubAccount)
		if err != nil {
			found = false
		} else if subAccount != nil && subAccount.AccountID == subAccountID {
			log.Printf("[INFO] subaccount : %v\n", subAccount)
			found = true
		}

		if found {
			return fmt.Errorf("Incapsula SubAccoint %d still exists", subAccountID)
		}
	}

	//clear subAccount defualt waf policy
	policies, err := client.GetAllPoliciesForAccount(strconv.Itoa(client.accountStatus.AccountID))
	if err != nil {
		return fmt.Errorf("Received an error : %s", err.Error())
	}
	defaultPolicyName := fmt.Sprintf("%s %s", "Generated WAF default Policy For new Account", subAccountIDStr)
	for i := range *policies {
		policy := (*policies)[i]
		if policy.PolicyType == "WAF_RULES" && policy.Name == defaultPolicyName {
			client.DeletePolicy(strconv.Itoa(policy.ID), nil)
		}

	}
	return nil
}

func testACCStateSubAccountID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != subAccountResourceType {
			continue
		}

		subAccountID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", subAccountID), nil
	}

	return "", fmt.Errorf("Error finding SubAccount ID")
}

func testCheckIncapsulaSubAccountExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resource := subAccountResourceType + "." + subAccountResourceName
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Incapsula SubAccount resource not found : %s", subAccountResourceType)
		}

		apiID := res.Primary.ID
		if !ok || apiID == "" {
			return fmt.Errorf("Incapsula API SubAccount API ID does not exist for API SubAccount")
		}

		subAccountIDStr := apiID
		subAccountID, _ := strconv.Atoi(apiID)
		if !ok || subAccountIDStr == "" {
			return fmt.Errorf("Incapsula API SubAccount does not exists")
		}

		client := testAccProvider.Meta().(*Client)
		log.Printf("[INFO] **** subAccountID: %d", subAccountID)
		subAccount, err := client.AccountStatus(subAccountID, ReadSubAccount)
		if err != nil {
			return err
		}

		found := false

		if subAccount != nil && subAccount.AccountID == subAccountID {
			log.Printf("[INFO] subaccount : %v\n", subAccount)
			found = true
		}

		if !found {
			return fmt.Errorf("Incapsula SubAccoint %d does not exist", subAccountID)
		}

		return nil
	}
}
