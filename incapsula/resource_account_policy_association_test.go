package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"sort"
	"strconv"
	"time"

	"testing"
)

const accountPolicyAssociationResourceName = "incapsula_account_policy_association"
const accountPolicyAssociationResource = accountPolicyAssociationResourceName + "." + accountPolicyAssociationName
const accountPolicyAssociationName = "testacc-account-policy-association-parent"

func TestAccIncapsulaAccountPolicyAssociation_basic(t *testing.T) {
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
					validateResourceSlice("default_non_mandatory_policy_ids", 1, func(apav3 AccountPolicyAssociationV3) []string {
						return ToStringSlice(apav3.DefaultNonMandatoryNonDistinctPolicyIds)
					}),
				),
			},
			{
				ResourceName:      accountPolicyAssociationResource,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testACCStateAccountPolicyAssociationID,
			},
		},
	})
}

func getValueFromResource(state *terraform.State, name string, key string) (string, error) {
	res, ok := state.RootModule().Resources[name]
	if !ok {
		log.Printf("[ERROR] resource not found: %s\n", name)
		return "", fmt.Errorf("Failed to find resource: %s ", name)
	}

	val, ok := res.Primary.Attributes[key]

	if !ok {
		log.Printf("[ERROR] Failed to extract value for: %s\n", key)
		return "", fmt.Errorf("Failed to extract value for: %s", key)
	}
	return val, nil
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func compareClientResourceSlices(valFromResource []string, valFromClient []string, resourceKey string) error {
	sort.Slice(valFromResource, func(i, j int) bool {
		return valFromResource[i] < valFromResource[j]
	})

	sort.Slice(valFromClient, func(i, j int) bool {
		return valFromClient[i] < valFromClient[j]
	})

	if !stringSlicesEqual(valFromResource, valFromClient) {
		log.Printf("[ERROR] Failed to match values for key: %s. expected: %v, got: %s\n", resourceKey, valFromClient, valFromResource)
		return fmt.Errorf("Failed to match values for key: %s. expected: %v, got: %s", resourceKey, valFromClient, valFromResource)
	}

	return nil
}

func getAccountPolicyAssociationFromClient(state *terraform.State) (*AccountPolicyAssociationV3, error) {
	accountIDStr, err := getValueFromResource(state, accountPolicyAssociationResource, "account_id")
	if err != nil {
		return nil, err
	}
	_, err = strconv.Atoi(accountIDStr)
	if err != nil {
		log.Printf("[ERROR] Could not convert to int ID: %s - %s\n", accountIDStr, err)
		return nil, err
	}

	client := testAccProvider.Meta().(*Client)
	return client.GetAccountPolicyAssociation(accountIDStr)
}

func validateResourceSlice(resourceKey string, expectedSize int, f func(AccountPolicyAssociationV3) []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		getAccountPolicyAssociation, err := getAccountPolicyAssociationFromClient(state)
		if err != nil {
			return err
		}

		valFromClient := f(*getAccountPolicyAssociation)
		var valFromResource []string
		for i := 0; i < expectedSize; i++ {
			val, err := getValueFromResource(state, accountPolicyAssociationResource, resourceKey+"."+strconv.Itoa(i))
			if err != nil {
				return err
			}
			valFromResource = append(valFromResource, val)
		}

		return compareClientResourceSlices(valFromResource, valFromClient, resourceKey)
	}
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

	aclName := "My-acl-created-on-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	return createPolicyResourceString(aclName, true, "ACL", aclPolicySettingsUrlExceptions) +
		fmt.Sprintf(`
			data "incapsula_account_data" "account_data" {
			}
			
			resource "%s" "%s" {
				account_id                       = data.incapsula_account_data.account_data.current_account
				default_non_mandatory_policy_ids = [
				   "${%s.id}"
				]
			}`, accountPolicyAssociationResourceName, accountPolicyAssociationName, policyResourceTypeAndName+aclName)
}
