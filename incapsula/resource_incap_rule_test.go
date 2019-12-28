package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, ruleID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaIncapRuleDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteID := res.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		incapRes, ok := state.RootModule().Resources[incapRuleResourceName]
		if !ok {
			return fmt.Errorf("Incapsula incap rule: %s resource not found: (site_id: %s)", incapRuleName, siteID)
		}

		ruleID := incapRes.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula incap rule ID does not exist")
		}

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.ReadIncapRule(siteID, ruleIDInt)
		if err == nil {
			return fmt.Errorf("Incapsula incap rule: %s (site id: %s) still exists\n%s", incapRuleName, siteID, err)
		}

		return nil
	}

	return nil
}

func testCheckIncapsulaIncapRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		siteRes, siteOk := state.RootModule().Resources[siteResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula site resource not found: %s", siteResourceName)
		}

		siteID := siteRes.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula incap rule resource not found: %s", name)
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula incap rule ID does not exist")
		}

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.ReadIncapRule(siteID, ruleIDInt)
		if err != nil {
			return fmt.Errorf("Incapsula incap rule: %s (site id: %s) does not exist\n%s", name, siteID, err)
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
