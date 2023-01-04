package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const incapsulaAccountPermissions = "incapsula_account_permissions"
const accountPermissionsNoFilter = "account_permissions_no_filter"
const accountPermissionsWithFilter = "account_permissions_with_filter"
const filterText = "site"
const accountPermissionsDataSourceNoFilterName = "data." + incapsulaAccountPermissions + "." + accountPermissionsNoFilter
const accountPermissionsDataSourceWithFilterName = "data." + incapsulaAccountPermissions + "." + accountPermissionsWithFilter

func TestAccIncapsulaDataSourceAccountPermissions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceAccountPermissionsConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					// Check the map is not empty
					resource.TestMatchResourceAttr(accountPermissionsDataSourceNoFilterName, "map.%", regexp.MustCompile("^[1-9]\\d*$")),
					// Check the list (keys) is empty (null)
					resource.TestMatchResourceAttr(accountPermissionsDataSourceNoFilterName, "keys.%", regexp.MustCompile("")),
					// Check the map is not empty
					resource.TestMatchResourceAttr(accountPermissionsDataSourceWithFilterName, "map.%", regexp.MustCompile("^[1-9]\\d*$")),
					// Check the list (keys) exists (cannot check a specific number as it depends on the account and the filter)
					resource.TestMatchResourceAttr(accountPermissionsDataSourceWithFilterName, "keys.#", regexp.MustCompile("^\\d+$")),
				),
			},
		},
	})
}

func testAccCheckIncapsulaDataSourceAccountPermissionsConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
data "%s" "%s" {}
data "%s" "%s" { 
  account_id = data.%s.%s.current_account
}
data "%s" "%s" { 
  account_id = data.%s.%s.current_account
  filter_by_text="%s"
}`,
		incapsulaAccountData, accountData,
		incapsulaAccountPermissions, accountPermissionsNoFilter, incapsulaAccountData, accountData,
		incapsulaAccountPermissions, accountPermissionsWithFilter, incapsulaAccountData, accountData, filterText,
	)
}
