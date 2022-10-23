package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

const policyResourceType = "incapsula_policy"
const policyResourceName = "testacc-terraform-acl-policy"
const policyResourceTypeAndName = policyResourceType + "." + policyResourceName

const aclPolicyName = "acl-policy-test"

var createdPoliciesNames = []string{aclPolicyName}

const aclPolicySettingsUrlExceptions = "[\n" +
	"    {\n" +
	"        \"settingsAction\": \"BLOCK\",\n" +
	"        \"policySettingType\": \"IP\",\n" +
	"        \"data\": {\n" +
	"            \"ips\": [\n" +
	"                \"10.10.10.10\",\n" +
	"                \"10.10.10.11\",\n" +
	"                \"10.10.10.24\"\n" +
	"            ]\n" +
	"        },\n" +
	"        \"policyDataExceptions\": [\n" +
	"            {\n" +
	"                \"data\": [\n" +
	"                    {\n" +
	"                        \"exceptionType\": \"CLIENT_ID\",\n" +
	"                        \"values\": [\n" +
	"                            \"144\"\n" +
	"                        ]\n" +
	"                    },\n" +
	"                    {\n" +
	"                        \"exceptionType\": \"IP\",\n" +
	"                        \"values\": [\n" +
	"                            \"10.10.192.10\"\n" +
	"                        ]\n" +
	"                    }\n" +
	"                ],\n" +
	"                \"comment\": \"Adding first exception to policy settings\"\n" +
	"            }\n" +
	"        ]\n" +
	"    },\n" +
	"    {\n" +
	"        \"settingsAction\": \"BLOCK\",\n" +
	"        \"policySettingType\": \"GEO\",\n" +
	"        \"data\": {\n" +
	"            \"geo\": {\n" +
	"                \"countries\": [\n" +
	"                    \"WF\",\n" +
	"                    \"AD\"\n" +
	"                ],\n" +
	"                \"continents\": [\n" +
	"                    \"AS\",\n" +
	"                    \"EU\"\n" +
	"                ]\n" +
	"            }\n" +
	"        }\n" +
	"    }\n" +
	"]"

func TestAccIncapsulaPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, aclPolicyName, true, "ACL", aclPolicySettingsUrlExceptions),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName, "name", aclPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName, "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(policyResourceTypeAndName, "policy_settings", aclPolicySettingsUrlExceptions),
				),
			},
			{
				ResourceName:      policyResourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStatePolicyID,
			},
		},
	})
}

func testAccStatePolicyID(state *terraform.State) (string, error) {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != policyResourceType {
			continue
		}

		policyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", policyID), nil
	}

	return "", fmt.Errorf("Error finding Policy ID")
}

func testCheckIncapsulaPolicyExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Policy resource not found: %s", name)
		}
		policyIDStr := res.Primary.ID
		if policyIDStr == "" {
			return fmt.Errorf("Incapsula policy ID does not exist")
		}

		_, err := strconv.Atoi(policyIDStr)
		if err != nil {
			return fmt.Errorf("Policy ID conversion error for %s: %s", policyIDStr, err)
		}
		client := testAccProvider.Meta().(*Client)
		_, err = client.GetPolicy(policyIDStr)
		if err != nil {
			return fmt.Errorf("Get Incapsula Policy return error %s", err)
		}
		return nil
	}
}

func testAccCheckIncapsulaPolicyConfigBasic(t *testing.T, policyName string, enabled bool, policyType string, policySettings string) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
    name        = "%s"
    enabled     = %s
    policy_type = "%s"
    policy_settings = <<POLICY
%s 
POLICY
}`, policyResourceType, policyResourceName, policyName, strconv.FormatBool(enabled), policyType, policySettings,
	)
}

func testAccCheckIncapsulaPolicyDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_account" {
			continue
		}

		accountID := res.Primary.ID
		if accountID == "" {
			// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
			// Specifically, upgrades/destroys are happening simultaneously and not honoring
			// dependencies. In this case, it's possible that the site has already been deleted,
			// which means that all the sub resources will have been cleared out.
			// Ordinarily, this should return an error, but until this gets addressed, we're
			// going to simply return nil.
			// return fmt.Errorf("Incapsula site ID does not exist")
			return nil
		}

		getAllPoliciesResponse, _ := client.GetAllPoliciesForAccount(accountID)

		for _, policyFromResponse := range *getAllPoliciesResponse {

			for _, val := range createdPoliciesNames {
				if val == policyFromResponse.Name {
					return fmt.Errorf("Incapsula Policy with name: %s (id: %d) still exists", policyFromResponse.Name, policyFromResponse.ID)
				}
			}
		}
	}
	return nil
}
