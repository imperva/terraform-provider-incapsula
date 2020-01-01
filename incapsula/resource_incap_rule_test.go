package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const incapRuleResourceName = "incapsula_incap_rule.testacc-terraform-incap-rule"
const incapRuleName = "Example Incap Rule Alert"

func TestAccIncapsulaIncapRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaIncapRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaIncapRuleConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapRuleExists(incapRuleResourceName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "name", incapRuleName),
				),
			},
			{
				ResourceName:      incapRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateRuleID,
			},
		},
	})
}

func testAccStateRuleID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_incap_rule" {
			continue
		}

		ruleID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing Rule ID %v to int", rs.Primary.ID)
		}

		siteID := rs.Primary.Attributes["site_id"]

		return fmt.Sprintf("%s/%d", siteID, ruleID), nil
	}

	return "", fmt.Errorf("Error finding Site ID")
}

func testAccCheckIncapsulaIncapRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_incap_rule" {
			continue
		}

		ruleID, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula Site ID does not exist for Rule ID %d", ruleID)
		}

		_, err = client.ReadIncapRule(siteID, ruleID)

		if err == nil {
			return fmt.Errorf("Incapsula Incap Rule %d still exists for Site ID %s", ruleID, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaIncapRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Incap Rule resource not found: %s", name)
		}

		ruleID, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID == "" {
			return fmt.Errorf("Incapsula Site ID does not exist for Rule ID %d", ruleID)
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.ReadIncapRule(siteID, ruleID)
		if err != nil {
			return fmt.Errorf("Incapsula Incap Rule: %s (site id: %s) does not exist\n%s", name, siteID, err)
		}

		return nil
	}
}

func testAccCheckIncapsulaIncapRuleConfig_basic() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf(`
resource "incapsula_incap_rule" "testacc-terraform-incap-rule" {
  name = "%s"
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["%s"]
}`, incapRuleName, siteResourceName,
	)
}
