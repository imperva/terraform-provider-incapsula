package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const incapsulaAccountData = "incapsula_account_data"
const accountData = "account_data"
const accountDataSourceName = "data." + incapsulaAccountData + "." + accountData

func TestAccIncapsulaDataSourceAccount_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceAccountConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(accountDataSourceName, "current_account"),
					resource.TestCheckResourceAttrSet(accountDataSourceName, "plan_name"),
				),
			},
		},
	})
}

func testAccCheckIncapsulaDataSourceAccountConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
data "%s" "%s" {}`,
		incapsulaAccountData, accountData,
	)
}
