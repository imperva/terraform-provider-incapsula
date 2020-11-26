package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testEmail = "example@example.com"
const accountResourceName = "incapsula_account.test-terraform-account"

func TestIncapsulaAccount_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(testEmail),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountExists(accountResourceName),
					resource.TestCheckResourceAttr(accountResourceName, "email", testEmail),
				),
			},
		},
	})
}

func TestIncapsulaAccount_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(testEmail),
			},
			{
				ResourceName:      "incapsula_account.test-terraform-account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIncapsulaAccountDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_account" {
			continue
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("Account ID conversion error for %s: %s", accountIDStr, err)
		}

		_, err = client.AccountStatus(accountID)

		if err == nil {
			return fmt.Errorf("Incapsula account id: %d still exists", accountID)
		}
	}

	return nil
}

func testCheckIncapsulaAccountExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula account resource not found: %s", name)
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("Account ID conversion error for %s: %s", accountIDStr, err)
		}

		client := testAccProvider.Meta().(*Client)
		accountStatusResponse, err := client.AccountStatus(accountID)
		if accountStatusResponse == nil {
			return fmt.Errorf("Incapsula account id: %d does not exist", accountID)
		}

		return nil
	}
}

func testCheckIncapsulaAccountConfigBasic(email string) string {
	return fmt.Sprintf(`
		resource "incapsula_account" "test-terraform-account" {
			email = "%s"
		}`,
		email,
	)
}
