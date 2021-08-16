package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const cacheRuleResourceName = "incapsula_cache_rule.testacc-terraform-cache-rule"
const cacheRuleName = "Example Cache Rule"

func TestAccIncapsulaCacheRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCacheRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCacheRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCacheRuleExists(cacheRuleResourceName),
					resource.TestCheckResourceAttr(cacheRuleResourceName, "name", cacheRuleName),
				),
			},
			{
				ResourceName:      cacheRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateCacheRuleID,
			},
		},
	})
}

func testAccStateCacheRuleID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_cache_rule" {
			continue
		}

		ruleID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing Cache Rule ID %v to int", rs.Primary.ID)
		}

		siteID := rs.Primary.Attributes["site_id"]

		return fmt.Sprintf("%s/%d", siteID, ruleID), nil
	}

	return "", fmt.Errorf("Error finding Site ID")
}

func testAccCheckIncapsulaCacheRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_cache_rule" {
			continue
		}

		ruleID, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula Site ID does not exist for Cache Rule ID %d", ruleID)
		}

		_, statusCode, err := client.ReadCacheRule(siteID, ruleID)
		if statusCode != 404 {
			return fmt.Errorf("Incapsula Cache Rule %d (site id: %s) should have received 404 status code", ruleID, siteID)
		}
		if err == nil {
			return fmt.Errorf("Incapsula Cache Rule %d still exists for Site ID %s", ruleID, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaCacheRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Cache Rule resource not found: %s", name)
		}

		ruleID, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID == "" {
			return fmt.Errorf("Incapsula Site ID does not exist for Cache Rule ID %d", ruleID)
		}

		client := testAccProvider.Meta().(*Client)
		_, statusCode, err := client.ReadCacheRule(siteID, ruleID)
		if statusCode != 200 {
			return fmt.Errorf("Incapsula Cache Rule: %s (site id: %s) should have received 200 status code", name, siteID)
		}
		if err != nil {
			return fmt.Errorf("Incapsula Cache Rule: %s (site id: %s) does not exist", name, siteID)
		}

		return nil
	}
}

func testAccCheckIncapsulaCacheRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_cache_rule" "testacc-terraform-cache-rule" {
	name = "%s"
  	site_id = incapsula_site.testacc-terraform-site.id
	action  = "HTTP_CACHE_MAKE_STATIC"
	enabled = true
	filter  = "ParamExists == \"true\""
	ttl     = 3600
	depends_on = ["%s"]
}`, cacheRuleName, siteResourceName,
	)
}
