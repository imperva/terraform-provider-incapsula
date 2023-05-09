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
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesBasic(t, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", "4002"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.mitigation_enabled", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesMultipleWebsites(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", "4002"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.mitigation_enabled", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.0.website_id", "11113"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.0.mitigation_enabled", "false"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.1.website_id", "11114"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.1.mitigation_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAbpWebsites_Basic2(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_abp_websites_test.TestAccAbpWebsites_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesBasic(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", "4002"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.mitigation_enabled", "false"),
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

		accountId := rs.Primary.Attributes["account_id"]

		client := testAccProvider.Meta().(*Client)

		response, _ := client.ReadAbpWebsites(accountId)
		if response == nil {
			return fmt.Errorf("Failed to retrieve ABP Websites (id=%d)", accountID)
		}

		*websitesresponse = *response
		return nil
	}
}

func testAccAbpWebsitesMultipleWebsites(t *testing.T) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = 4002
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				website_id = 11112
				mitigation_enabled = true
			}
		}
		website_group {
			name = "sites-2"
			website {
				website_id = 11113
				mitigation_enabled = false
			}
			website {
				website_id = 11114
				mitigation_enabled = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsitesBasic(t *testing.T, mitigationEnabled bool) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = 4002
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				website_id = 11112
				mitigation_enabled = %t
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, mitigationEnabled)
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
