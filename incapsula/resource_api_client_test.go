package incapsula

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const accountResourceApiClientType = "incapsula_api_client"
const accountResourceApiClientName = "test-terraform-api-client"
const accountResourceApiClientTypeName = accountResourceApiClientType + "." + accountResourceApiClientName
const apiClientName = "test-terraform"
const apiClientNameUpdated = "test-terraform updated"
const apiClientDesc = "Test terraform description"
const apiClientDescUpdated = "Test terraform description updated"
const apiClientExpPeriod = "2026-01-30T23:59:59Z"
const apiClientExpPeriodUpdated = "2026-02-30T23:59:59Z"

func TestIncapsulaApiClient_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaApiClientDestroy,
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
				Config: testCheckIncapsulaApiClientConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiClientExists(accountResourceApiClientTypeName),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "name", apiClientName),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "description", apiClientDesc),
				),
			},
			{
				Config: testCheckIncapsulaApiClientConfigUpdate(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiClientExists(accountResourceApiClientTypeName),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "name", apiClientNameUpdated),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "description", apiClientDescUpdated),
					resource.TestCheckResourceAttr(accountResourceApiClientTypeName, "enabled", "false"),
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

		// There may be a timing/race condition here
		// Set an arbitrary period to sleep
		time.Sleep(15 * time.Second)

		apiClientResponse, _ := client.GetAPIClient(0, "", apiClientIDStr)

		// Account object may have been deleted
		if apiClientResponse != nil {
			return fmt.Errorf(
				"Incapsula api client with id: %s still exists",
				apiClientIDStr,
			)
		}
	}

	log.Printf("[DEBUG] **** destroy test return nil")

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
			expiration_period = "%s"
		}`,
		accountResourceApiClientType, accountResourceApiClientName, apiClientName, apiClientDesc, true, apiClientExpPeriod,
	)
}

func testCheckIncapsulaApiClientConfigUpdate(t *testing.T) string {
	return fmt.Sprintf(`

		resource "%s" "%s" {
			name = "%s"
			description = "%s"
			enabled = %v
		}`,
		accountResourceApiClientType, accountResourceApiClientName, apiClientNameUpdated, apiClientDescUpdated, false,
	)
}
