package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"

	"testing"
)

const POLICY_NON_MANDATORY = "450257"
const POLICY_WAF = "1093549"

const accountPolicyAssociationResourceName = "incapsula_account_policy_association"
const accountPolicyAssociationResource = accountPolicyAssociationResourceName + "." + accountPolicyAssociationName
const accountPolicyAssociationName = "testacc-account-policy-association-parent"

func TestAccIncapsulaAccountPolicyAssociation_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_account_policy_association.TestAccIncapsulaAccountPolicyAssociation_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckAccountPolicyAssociationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountPolicyAssociationExists(),
					resource.TestCheckResourceAttr(accountPolicyAssociationResource, "default_non_mandatory_policy_ids.0", POLICY_NON_MANDATORY),
					resource.TestCheckResourceAttr(accountPolicyAssociationResource, "default_waf_policy_id", POLICY_WAF),
				),
			},
			{
				ResourceName:      accountPolicyAssociationResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateAccountPolicyAssociationID,
			},
		},
	})
}

func testACCStateAccountPolicyAssociationID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != accountPolicyAssociationResourceName {
			continue
		}
		accountID, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing account ID for import Account Policy Association resource command. Value %s", rs.Primary.Attributes["account_id"])
		}

		return fmt.Sprintf("%d", accountID), nil
	}
	return "", fmt.Errorf("Error finding an Account Policy Association\"")
}

func testCheckIncapsulaAccountPolicyAssociationExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[accountPolicyAssociationResource]
		if !ok {
			return fmt.Errorf("Incapsula Account Policy Association resource not found : %s", accountPolicyAssociationResource)
		}

		accountID := res.Primary.ID
		if !ok || accountID == "" {
			return fmt.Errorf("Incapsula Account Policy Association does not exist")
		}
		return nil
	}
}

func testAccCheckAccountPolicyAssociationBasic(t *testing.T) string {

	return fmt.Sprintf(`
data "incapsula_account_data" "account_data" {
}

resource "%s" "%s" {
    account_id                       = data.incapsula_account_data.account_data.current_account
    default_non_mandatory_policy_ids = [
       %s
    ]
    default_waf_policy_id = "%s"
}
`, accountPolicyAssociationResourceName, accountPolicyAssociationName, POLICY_NON_MANDATORY, POLICY_WAF)
}
