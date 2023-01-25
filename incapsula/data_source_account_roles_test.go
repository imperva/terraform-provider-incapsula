package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const incapsulaAccountRole = "incapsula_account_role"
const terraformAccTestRole = "terraform_acc_test_role"
const roleName = "terraform acc test role"
const incapsulaAccountRoles = "incapsula_account_roles"
const terraformAccTestRoles = "terraform_acc_test_roles"
const accountRolesDataSourceName = "data." + incapsulaAccountRoles + "." + terraformAccTestRoles
const defaultRoleAdministrator = "Administrator"
const defaultRoleReader = "Reader"

func TestAccIncapsulaDataSourceAccountRoles_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceAccountRolesConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					// Check the map have at least 2 elements (Administrator, Reader)
					resource.TestMatchResourceAttr(accountRolesDataSourceName, "map.%", regexp.MustCompile("^([2-9]\\d*|\\d{2,})$")),
					resource.TestCheckResourceAttrSet(accountRolesDataSourceName, "map."+defaultRoleAdministrator),
					resource.TestCheckResourceAttrSet(accountRolesDataSourceName, "map."+defaultRoleReader),
					// Negative Test, ensure the role doesn't exist before
					resource.TestCheckNoResourceAttr(accountRolesDataSourceName, "map."+roleName),
				),
			},
			{
				Config: testAccCheckIncapsulaDataSourceAccountRolesConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					// Check the map have at least 3 elements (Administrator, Reader and new created role)
					resource.TestMatchResourceAttr(accountRolesDataSourceName, "map.%", regexp.MustCompile("^([3-9]\\d*|\\d{2,})$")),
					resource.TestCheckResourceAttrSet(accountRolesDataSourceName, "map."+defaultRoleAdministrator),
					resource.TestCheckResourceAttrSet(accountRolesDataSourceName, "map."+defaultRoleReader),
					resource.TestCheckResourceAttrSet(accountRolesDataSourceName, "map."+roleName),
				),
			},
			{
				ResourceName:      incapsulaAccountRole + "." + terraformAccTestRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIncapsulaDataSourceAccountRolesConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
data "%s" "%s" {}
resource "%s" "%s" {
  account_id      = data.%s.%s.current_account
  name            = "%s"
}
data "%s" "%s" { 
  account_id          = data.%s.%s.current_account
}`,
		incapsulaAccountData, accountData,
		incapsulaAccountRole, terraformAccTestRole, incapsulaAccountData, accountData, roleName,
		incapsulaAccountRoles, terraformAccTestRoles, incapsulaAccountData, accountData,
	)
}
