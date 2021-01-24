package incapsula

import (
	"fmt"
	// "strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testName = "testSubAccount"
const subAccountResourceName = "incapsula_subaccount.test-terraform-subaccount"

func TestIncapsulaSubAccount_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(testEmail),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountExists(subAccountResourceName),
					resource.TestCheckResourceAttr(subAccountResourceName, "testSubAccount", testName),
				),
			},
		},
	})
}

func TestIncapsulaSubAccount_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(testName),
			},
			{
				ResourceName:      "incapsula_subaccount.test-terraform-subaccount",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIncapsulaSubAccountDestroy(state *terraform.State) error {

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_account" {
			continue
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

	}

	return nil
}

func testCheckIncapsulaSubAccountExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula account resource not found: %s", name)
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

		return nil
	}
}

func testCheckIncapsulaSubAccountConfigBasic(email string) string {
	return fmt.Sprintf(`
		resource "incapsula_account" "test-terraform-account" {
			email = "%s"
		}`,
		email,
	)
}
