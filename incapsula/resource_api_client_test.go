package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const accountResourceApiClientType = "incapsula_api_client"
const accountResourceApiClientName = "test-terraform-api-client"
const accountResourceApiClientTypeName = accountResourceUserType + "." + accountResourceUserName
const apiClientName = "test-terraform"
const apiClientDesc = "Test terraform description"

func TestIncapsulaApiClient_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaApiClientConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiClientExists(accountResourceApiClientTypeName),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "name", apiClientName),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "description", apiClientDesc),
				),
			},
		},
	})
}

func TestIncapsulaApiClient_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaApiClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountUserConfigBasic(t, accountUserEmail),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmail),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "role_ids.#", "0"),
				),
			},
			{
				Config: testCheckIncapsulaAccountUserConfigUpdate(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmail),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "role_ids.#", "1"),
				),
			},
		},
	})
}

func TestIncapsulaApiClient_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaApiClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaApiClientConfigBasic(t),
			},
			{
				ResourceName:      accountResourceApiClientTypeName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIncapsulaApiClientDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != accountResourceApiClientType {
			continue
		}

		apiClientIDStr := res.Primary.ID
		if apiClientIDStr == "" {
			return fmt.Errorf("Incapsula api client ID does not exist")
		}

		accountUserStatusResponse, _ := client.GetAccountUser(0, apiClientIDStr)

		// Account object may have been deleted
		if accountUserStatusResponse != nil {
			return fmt.Errorf(
				"Incapsula api client with id: %s still exists",
				apiClientIDStr,
			)
		}
	}

	return nil
}

func testCheckIncapsulaApiClientExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula api client resource not found: %s", name)
		}

		apiClientIDStr := res.Primary.ID
		if apiClientIDStr == "" {
			return fmt.Errorf("Incapsula api client ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)
		apiClientStatusResponse, _ := client.GetAPIClient(0, "", apiClientIDStr)
		if apiClientStatusResponse == nil {
			return fmt.Errorf(
				"Incapsula api client with id: %s does not exist", apiClientIDStr)
		}

		return nil
	}
}

func testCheckIncapsulaApiClientConfigBasic(t *testing.T) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
			description = "%s"
			enabled = %v
		}`,
		accountResourceApiClientType, accountResourceApiClientName, apiClientName, apiClientDesc, false,
	)
}

func testCheckIncapsulaApiClientConfigUpdate(t *testing.T) string {
	return fmt.Sprintf(`

		resource "%s" "%s" {
			name = "%s"
			description = "%s"
			enabled = %v
		}`,
		accountResourceApiClientType, accountResourceApiClientName, apiClientName, apiClientDesc, false,
	)
}
