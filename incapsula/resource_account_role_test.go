package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const accountResourceRoleType = "incapsula_account_role"
const accountResourceRoleName = "test-terraform-account-role"
const accountResourceRoleTypeName = accountResourceRoleType + "." + accountResourceRoleName
const accountRoleName = "role-test-terraform"
const accountRoleDescription = "role-description-test"
const accountRoleDescriptionUpdated = "role-description-test Updated"

func TestIncapsulaAccountRole_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountRoleConfigBasic(t, accountRoleDescription),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountRoleExists(accountResourceRoleTypeName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "name", accountRoleName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "description", accountRoleDescription),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "permissions.#", "2"),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "permissions.0", "canAddSite"),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "permissions.1", "canViewInfraProtectSetting"),
				),
			},
		},
	})
}

func TestIncapsulaAccountRole_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountRoleConfigBasic(t, accountRoleDescription),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountRoleExists(accountResourceRoleTypeName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "name", accountRoleName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "description", accountRoleDescription),
				),
			},
			{
				Config: testCheckIncapsulaAccountRoleConfigBasic(t, accountRoleDescriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountRoleExists(accountResourceRoleTypeName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "name", accountRoleName),
					resource.TestCheckResourceAttr(accountResourceRoleTypeName, "description", accountRoleDescriptionUpdated),
				),
			},
		},
	})
}

func TestIncapsulaAccountRole_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountRoleConfigBasic(t, accountRoleDescription),
			},
			{
				ResourceName: accountResourceRoleTypeName,
				ImportState:  true,
				// TODO - Setting to false - Not supported when state include data sources
				ImportStateVerify: false,
			},
		},
	})
}

func testCheckIncapsulaAccountRoleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != accountResourceRoleType {
			continue
		}

		accountRoleIDStr := res.Primary.ID
		if accountRoleIDStr == "" {
			return fmt.Errorf("Incapsula account role ID does not exist")
		}

		accountRoleID, err := strconv.Atoi(accountRoleIDStr)
		if err != nil {
			return fmt.Errorf("Account Role ID conversion error for %s: %s", accountRoleIDStr, err)
		}

		accountRoleResponse, err := client.GetAccountRole(accountRoleID)

		// Account object may have been deleted
		if accountRoleResponse != nil && accountRoleResponse.ErrorCode != 1047 {
			return fmt.Errorf("Incapsula account role id: %d still exists", accountRoleID)
		}
	}

	return nil
}

func testCheckIncapsulaAccountRoleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula account role resource not found: %s", name)
		}

		accountRoleIDStr := res.Primary.ID
		if accountRoleIDStr == "" {
			return fmt.Errorf("Incapsula account role ID does not exist")
		}

		accountRoleID, err := strconv.Atoi(accountRoleIDStr)
		if err != nil {
			return fmt.Errorf("Account Role ID conversion error for %s: %s", accountRoleIDStr, err)
		}

		client := testAccProvider.Meta().(*Client)
		accountRoleStatusResponse, err := client.GetAccountRole(accountRoleID)
		if accountRoleStatusResponse == nil {
			return fmt.Errorf("Incapsula account role id: %d does not exist", accountRoleID)
		}

		return nil
	}
}

func testCheckIncapsulaAccountRoleConfigBasic(t *testing.T, roleDescription string) string {
	return fmt.Sprintf(`
		data "incapsula_account_data" "account_data" {}

		data "incapsula_account_permissions" "incapsula_account_permissions" {
		  account_id = data.incapsula_account_data.account_data.current_account
		}

		resource "%s" "%s" {
			account_id = data.incapsula_account_data.account_data.current_account
			name = "%s"
			description = "%s"
			permissions = [
				data.incapsula_account_permissions.incapsula_account_permissions.map["View Infra Protect settings"],
				"canAddSite",
			]
		}`,
		accountResourceRoleType, accountResourceRoleName, accountRoleName, roleDescription,
	)
}
