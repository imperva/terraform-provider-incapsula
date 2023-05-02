package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const abpWebsitesResourceName = "incapsula_abp_websites"
const accountConfigName = "testacc-terraform-abp-websites"
const abpWebsitesResource = abpWebsitesResourceName + "." + accountConfigName

func TestAccAbpWebsites_Basic(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_abp_websites_test.TestAccAbpWebsites_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckAbpWebsitesExists(websitesresponse *AbpTerraformAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[abpWebsitesResource]
		if !ok {
			return fmt.Errorf("Not found: %s", abpWebsitesResource)
		}

		accountID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Incapsula ABP Websites with id %d doesn't have site ID", accountID)
		}

		accountId := rs.Primary.Attributes["account_id"]

		client := testAccProvider.Meta().(*Client)

		response, _ := client.ReadAbpWebsites(accountId)
		if response == nil {
			return fmt.Errorf("Failed to retrieve ABP Websites (id=%d)", accountID)
		}

		*websitesresponse = *response
		return fmt.Errorf("resource %s was not updated correctly after full update", abpWebsitesResource)
	}
}

func testAccAbpWebsitesBasic(t *testing.T) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = 4001
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				website_id = 1
				mitigation_enabled = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccCheckAbpWebsitesDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != abpWebsitesResourceName {
			continue
		}

		accountID := res.Primary.ID
		if accountID == "" {
			return fmt.Errorf("Incapsula ABP Websites ID does not exist")
		}

		siteID := res.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Incapsula ABP Websites with id %s doesn't have site ID", accountID)
		}

		accountId := res.Primary.Attributes["account_id"]

		websitesResponse, _ := client.ReadAbpWebsites(accountId)
		if websitesResponse == nil {
			return fmt.Errorf("Failed to check ABP Websites status (id=%s)", accountID)
		}
		// if websitesResponse.Errors[0].Status != 404 {
		// 	return fmt.Errorf("Incapsula ABP Websites with id %s still exists", accountID)
		// }
	}

	return nil
}

func testAccGetAbpWebsitesImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_abp_websites" {
			continue
		}

		accountID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %s to int", rs.Primary.Attributes["site_id"])
		}
		accountId := rs.Primary.Attributes["account_id"]

		return fmt.Sprintf("%s/%d/%d", accountId, siteID, accountID), nil
	}

	return "", fmt.Errorf("Error finding ABP Websites")
}
