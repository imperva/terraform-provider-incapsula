package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const abpWebsitesResourceName = "incapsula_abp_websites"
const accountConfigName = "testacc-terraform-abp-websites"
const abpWebsitesResource = abpWebsitesResourceName + "." + accountConfigName

const testAccountId = "50201698"

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
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", testAccountId),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesMultipleWebsites(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", testAccountId),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.0.website_id", "11113"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.0.enable_mitigation", "false"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.1.website_id", "11114"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.1.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesBasic2(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", "4002"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11113"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "false"),
				),
			},
		},
	})
}

func TestAccAbpWebsites_DuplicateWebsites(t *testing.T) {

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_abp_websites_test.TestAccAbpWebsites_DuplicateWebsites")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		ErrorCheck: func(err error) error {
			if strings.Contains(err.Error(), "already in use") {
				return nil
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesDuplicate(t),
			},
		},
	})
}

func TestAccAbpWebsites_AutoPublish(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_abp_websites_test.TestAccAbpWebsites_AutoPublish")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesAutoPublish(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", testAccountId),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "false"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesAutoPublish(t, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", testAccountId),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesAutoPublish(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttr(abpWebsitesResource, "account_id", testAccountId),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "false"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.website_id", "11112"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
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

func testAccCheckIncapsulaSiteConfig(name string, domain string) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "%s" {
			domain = "%s"
		}`,
		name,
		domain,
	)
}

func testAccAbpWebsitesMultipleWebsites(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfig("sites-1", GenerateTestDomain(t)) +
		testAccCheckIncapsulaSiteConfig("sites-2", GenerateTestDomain(t)) +
		fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = "%s"
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				website_id = 11112
				enable_mitigation = true
			}
		}
		website_group {
			name = "sites-2"
			website {
				website_id = 11113
				enable_mitigation = false
			}
			website {
				website_id = 11114
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, testAccountId)
}

func testAccAbpWebsitesBasic(t *testing.T, mitigationEnabled bool) string {
	return testAccCheckIncapsulaSiteConfig("sites-1", GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = "%s"
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				website_id = 11112
				enable_mitigation = %t
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, testAccountId, mitigationEnabled)
}

func testAccAbpWebsitesBasic2(t *testing.T, mitigationEnabled bool) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = 4002
		auto_publish = true
		website_group {
			name = "sites-2"
			website {
				website_id = 11113
				enable_mitigation = %t
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, mitigationEnabled)
}

func testAccAbpWebsitesDuplicate(t *testing.T) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = 4002
		auto_publish = true
		website_group {
			name = "sites-2"
			website {
				website_id = 11113
				enable_mitigation = false
			}
			website {
				website_id = 11113
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsitesAutoPublish(t *testing.T, autoPublish bool) string {
	return testAccCheckIncapsulaSiteConfig("sites-1", GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		account_id = "%s"
		auto_publish = %t
		website_group {
			name = "sites-1"
			website {
				website_id = 11112
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, testAccountId, autoPublish)
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

		if len(websitesResponse.WebsiteGroups) != 0 {
			return fmt.Errorf("Found some website groups remaining after delete: %+v", websitesResponse)
		}
		// if websitesResponse.Errors[0].Status != 404 {
		// 	return fmt.Errorf("Incapsula ABP Websites with id %s still exists", accountID)
		// }
	}

	return nil
}
