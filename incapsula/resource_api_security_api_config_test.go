package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	//	"strings"
	"testing"
)

const apiSecApiConfigResourceName = "incapsula_api_security_api_config"
const apiSecApiConfigResource = apiSecApiConfigResourceName + "." + apiSecApiConfigName
const apiSecApiConfigName = "testacc-terraform-api-security-api-config"
const swaggerFileContent = "<<-EOT\nswagger: \"2.0\"\ninfo:\n  title: Sample API\n  description: API description in Markdown.\n  version: 1.0.0\nhost: api.example.com\nbasePath: /v1\nschemes:\n  - https\npaths:\n  /users:\n    get:\n      summary: Returns a list of users.\n      description: Optional extended description in Markdown.\n      produces:\n        - application/json\n      responses:\n        200:\n          description: OK\nEOT"

func TestAccIncapsulaApiSecurityApiConfig_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_api_security_api_config_test.TestAccIncapsulaApiSecurityApiConfig_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateApiSecurityApiConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiConfigExists(),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "validate_host", "false"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "description", "first-api-security-collection"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "invalid_url_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "invalid_method_violation_action", "BLOCK_IP"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "missing_param_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "invalid_param_value_violation_action", "BLOCK_IP"),
					resource.TestCheckResourceAttr(apiSecApiConfigResource, "invalid_param_name_violation_action", "IGNORE"),
				),
			},
			{
				ResourceName:      apiSecApiConfigResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateApiSecurityApiConfigID,
			},
		},
	})
}

func testACCStateApiSecurityApiConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != apiSecApiConfigResourceName {
			continue
		}
		return nil

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			fmt.Errorf("Parameter site_id was not found in resource %s", apiSecApiConfigResourceName)
		}
		siteIDInt, err := strconv.Atoi(siteID)
		if err != nil {
			fmt.Errorf("failed to convert Site Id from import command, actual value : %s, expected numeric id", siteID)
		}

		apiID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			fmt.Errorf("failed to convert API Id from import command, actual value : %s, expected numeric id", rs.Primary.ID)
		}

		_, err = client.GetApiSecurityApiConfig(siteIDInt, apiID)
		if err == nil {
			return fmt.Errorf("Resource%sforIncapsulaApiSecurityApi:ApiId%d,siteID%dstillexists", apiSecApiConfigResourceName, apiID, siteIDInt)
		}
	}
	return fmt.Errorf("Error finding site_id")
}

func testACCStateApiSecurityApiConfigID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != apiSecApiConfigResourceName {
			continue
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site ID %v to int", rs.Primary.Attributes["site_id"])
		}
		apiID, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing API ID %v to int", rs.Primary.Attributes["id"])
		}

		return fmt.Sprintf("%d/%d", siteID, apiID), nil
	}
	return "", fmt.Errorf("Error finding an API Security API Config\"")
}

func testCheckIncapsulaApiConfigExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[apiSecApiConfigResource]
		if !ok {
			return fmt.Errorf("Incapsula Api Securiy Config resource not found : %s", apiSecApiConfigResource)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID != siteID {
			return fmt.Errorf("Incapsula API Security Site ID %s does not exist for API Site Config", siteID)
		}
		siteIdInt, err := strconv.Atoi(siteID)

		apiID := res.Primary.ID
		if !ok || apiID == "" {
			return fmt.Errorf("Incapsula API Security API ID%s does not exist for API Config", siteID)
		}
		apiIDInt, err := strconv.Atoi(apiID)

		client := testAccProvider.Meta().(*Client)
		_, err = client.GetApiSecurityApiConfig(siteIdInt, apiIDInt)
		if err != nil {
			return fmt.Errorf("Incapsula API Security API Config : %s (SiteId : %d, API Id %d) does not exist", apiSecApiConfigResource, siteIdInt, apiIDInt)
		}

		return nil
	}
}

func testAccCheckApiConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource"%s""%s"{
site_id=incapsula_site.testacc-terraform-site.id
validate_host="false"
description="first-api-security-collection"
invalid_url_violation_action="IGNORE"
invalid_method_violation_action="BLOCK_IP"
missing_param_violation_action="IGNORE"
invalid_param_value_violation_action="BLOCK_IP"
invalid_param_name_violation_action="IGNORE"
depends_on=["%s"]
api_specification = %s
}`,
		apiSecApiConfigResourceName, apiSecApiConfigName, siteResourceName, swaggerFileContent,
	)
}
