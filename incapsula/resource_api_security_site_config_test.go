package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const apiSiteConfigResourceName = "incapsula_api_security_site_config"
const apiSiteConfigResource = apiSiteConfigResourceName + "." + apiSiteConfigName
const apiSiteConfigName = "testacc-terraform-api-security-site-config"

func TestAccIncapsulaApiSecuritySiteConfig_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_api_security_site_config_test.TestAccIncapsulaApiSecuritySiteConfig_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiSiteConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckApiSecuritySiteConfigExists(apiSiteConfigResource),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "is_api_only_site", "true"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "non_api_request_violation_action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_url_violation_action", "BLOCK_REQUEST"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_method_violation_action", "BLOCK_IP"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "missing_param_violation_action", "ALERT_ONLY"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_param_value_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_param_name_violation_action", "BLOCK_IP"),
				),
			},
			{
				ResourceName:      apiSiteConfigResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateApiSiteConfigID,
			},
		},
	})
}

func testCheckApiSecuritySiteConfigExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Api Security Site Config resource not found: %s", name)
		}
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.ReadApiSecuritySiteConfig(siteId)
		if err != nil {
			fmt.Errorf("Incapsula Api Security Site Config doesn't exist")
		}

		return nil
	}
}

func testACCStateApiSiteConfigID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != apiSiteConfigResourceName {
			continue
		}

		//we don't need thi block because site_config doesn't have an ID
		ruleID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing Api Security Site Config ID %v to int", rs.Primary.ID)
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		iD, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.Attributes["site_id"])
		}
		if siteID != iD {
			return "", fmt.Errorf("Incapsula API Security Site Config does not exist")
		}
		return fmt.Sprintf("%d", ruleID), nil
	}
	return "", fmt.Errorf("Error finding site_id in Api Security Site Config resource")
}

func testAccCheckApiSiteConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "incapsula_api_security_site_config" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		is_automatic_discovery_api_integration_enabled = false
		is_api_only_site = true
		non_api_request_violation_action = "BLOCK_USER"
		invalid_url_violation_action = "BLOCK_REQUEST"
		invalid_method_violation_action = "BLOCK_IP"
		invalid_param_value_violation_action = "IGNORE"
		invalid_param_name_violation_action = "BLOCK_IP"
  		depends_on = ["%s"]
	}`,
		apiSiteConfigName, siteResourceName,
	)
}
