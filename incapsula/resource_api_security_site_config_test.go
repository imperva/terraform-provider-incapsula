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

func TestAccIncapsulaApiSecuritySiteConfig_Basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_api_security_site_config_test.TestAccIncapsulaApiSecuritySiteConfig_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckIncapsulaIncapRuleDestroy, //todo - don't do destroy because there's no option to delete site
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiSiteConfigBasic(t), //todo  change
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteConfigAttributeCorrect(apiSiteConfigResource, "site_id"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "site_id", "78865045"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "api_only_site", "true"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "non_api_request_violation_action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_url_violation_action", "BLOCK_REQUEST"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_method_violation_action", "BLOCK_IP"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "missing_param_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_param_value_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "invalid_param_name_violation_action", "BLOCK_IP"),
					resource.TestCheckResourceAttr(apiSiteConfigResource, "is_automatic_discovery_api_integration_enabled", "false"),
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

func testACCStateApiSiteConfigID(s *terraform.State) (string, error) {
	//return "", fmt.Errorf("Resources: %v", s.RootModule().Resources)
	for _, rs := range s.RootModule().Resources {
		fmt.Errorf("Resource: %v", rs)
		if rs.Type != apiSiteConfigResourceName {
			continue
		}

		//we don't need thi block because site_config doesn't have an ID
		ruleID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing Cache Rule ID %v to int", rs.Primary.ID)
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		siteID, err := strconv.Atoi(rs.Primary.ID)
		//return "", fmt.Errorf("Extracting ID %v to int", rs.Primary.ID)
		//todo compare ruleId to SiteID
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.Attributes["site_id"])
		}
		fmt.Errorf("%d", ruleID)
		return fmt.Sprintf("%d", ruleID), nil
	}
	return "", fmt.Errorf("Error finding site_id")
}

func testCheckIncapsulaSiteConfigAttributeCorrect(resourceName, attrName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Incapsula Api Securiy Cinfig resource not found: %s", resourceName)
		}

		siteID, ok := res.Primary.Attributes[attrName]
		if !ok || siteID == "" {
			return fmt.Errorf("Incapsula Site ID %s does not exist for API Site Config", siteID)
		}
		client := testAccProvider.Meta().(*Client)
		siteIdInt, err := strconv.Atoi(siteID)
		_, err = client.ReadApiSecuritySiteConfig(siteIdInt)
		if err != nil {
			return fmt.Errorf("Incapsula API Site Config: %s (site id: %s) does not exist", resourceName, siteID)
		}

		return nil
	}
}

//		site_id = 78865045		todo - move id to variable
func testAccCheckApiSiteConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "incapsula_api_security_site_config" "%s" {
		site_id = 78865045		
		api_only_site = "true"
		discovery_enabled = "true"
		non_api_request_violation_action = "BLOCK_USER"
		invalid_url_violation_action = "BLOCK_REQUEST"
		invalid_method_violation_action = "BLOCK_IP"
		missing_param_violation_action = "IGNORE"
		invalid_param_value_violation_action = "IGNORE"
		invalid_param_name_violation_action = "BLOCK_IP"
		is_automatic_discovery_api_integration_enabled = "false"
	}`,
		apiSiteConfigName,
	)
}
