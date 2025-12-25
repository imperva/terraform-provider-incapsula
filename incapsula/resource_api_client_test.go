package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

const resourceApiClientType = "incapsula_api_client"
const resourceApiClientName = "test-terraform-api-client"
const resourceApiClientTypeName = resourceApiClientType + "." + resourceApiClientName
const apiClientName = "test-terraform"
const apiClientNameUpdated = "test-terraform updated"
const apiClientDesc = "Test terraform description"
const apiClientDescUpdated = "Test terraform description updated"
const apiClientExpPeriod = "2026-12-31T23:59:59Z"
const apiClientExpPeriodUpdated = "2027-12-31T23:59:59Z"
const apiClientEmail = "test-terraform@www.com"

func TestIncapsulaApiClient_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testCheckIncapsulaApiClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaApiClientConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiClientExists(resourceApiClientTypeName),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "name", apiClientName),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "description", apiClientDesc),
				),
			},
		},
	})
}

func TestIncapsulaApiClient_WithEmail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaApiClientDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testCheckIncapsulaApiClientConfigWithEmail(t),
				ExpectError: regexp.MustCompile(`Can't find user with email`),
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
					testCheckIncapsulaApiClientExists(resourceApiClientTypeName),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "name", apiClientName),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "description", apiClientDesc),
				),
			},
			{
				Config: testCheckIncapsulaApiClientConfigUpdate(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiClientExists(resourceApiClientTypeName),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "name", apiClientNameUpdated),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "description", apiClientDescUpdated),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "enabled", "false"),
				),
			},
		},
	})
}

func TestIncapsulaApiClient_UpdateExpiration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaApiClientDestroy,
		Steps: []resource.TestStep{
			// 1. Create with initial expiration
			{
				Config: testCheckIncapsulaApiClientConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceApiClientTypeName, "api_key"),
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "expiration_date", apiClientExpPeriod),
				),
			},
			// 2. Update to a later expiration date (triggers regeneration)
			{
				Config: testCheckIncapsulaApiClientConfigExpirationUpdate(t), // New config function
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceApiClientTypeName, "expiration_date", apiClientExpPeriodUpdated),
				),
			},
		},
	})
}

// New config function to explicitly test expiration update
func testCheckIncapsulaApiClientConfigExpirationUpdate(t *testing.T) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
			description = "%s"
			enabled = %v
			expiration_date = "%s"
		}`,
		resourceApiClientType, resourceApiClientName, apiClientName, apiClientDesc, true, apiClientExpPeriodUpdated,
	)
}

func testCheckIncapsulaApiClientDestroy(state *terraform.State) error {
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
		apiClientStatusResponse, _ := client.GetAPIClient(0, apiClientIDStr)
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
			expiration_date = "%s"
		}`,
		resourceApiClientType, resourceApiClientName, apiClientName, apiClientDesc, true, apiClientExpPeriod,
	)
}

func testCheckIncapsulaApiClientConfigWithEmail(t *testing.T) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			user_email = "%s"
			name = "%s"
			description = "%s"
			enabled = %v
			expiration_date = "%s"
		}`,
		resourceApiClientType, resourceApiClientName, apiClientEmail, apiClientName, apiClientDesc, true, apiClientExpPeriod,
	)
}

func testCheckIncapsulaApiClientConfigUpdate(t *testing.T) string {
	return fmt.Sprintf(`

		resource "%s" "%s" {
			name = "%s"
			description = "%s"
			enabled = %v
		}`,
		resourceApiClientType, resourceApiClientName, apiClientNameUpdated, apiClientDescUpdated, false,
	)
}
