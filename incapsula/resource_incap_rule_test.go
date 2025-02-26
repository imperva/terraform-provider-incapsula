package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				Config: testAccCheckIncapsulaIncapRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapRuleExists(incapRuleResourceName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "name", incapRuleName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(incapRuleResourceName, "send_notifications", "true"),
				),
			},
			{
				Config: testAccCheckIncapsulaIncapRuleConfigEnabledFlagNotProvided(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapRuleExists(incapRuleResourceName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "name", incapRuleName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(incapRuleResourceName, "send_notifications", "false"),
				),
			},
			{
				Config: testAccCheckIncapsulaIncapRuleConfigDisabled(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapRuleExists(incapRuleResourceName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "name", incapRuleName),
					resource.TestCheckResourceAttr(incapRuleResourceName, "enabled", "false"),
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

		_, statusCode, err := client.ReadIncapRule(siteID, ruleID)
		if statusCode != 404 {
			return fmt.Errorf("Incapsula Incap Rule %d (site id: %s) should have received 404 status code", ruleID, siteID)
		}
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
		_, statusCode, err := client.ReadIncapRule(siteID, ruleID)
		if statusCode != 200 {
			return fmt.Errorf("Incapsula Incap Rule: %s (site id: %s) should have received 200 status code", name, siteID)
		}
		if err != nil {
			return fmt.Errorf("Incapsula Incap Rule: %s (site id: %s) does not exist", name, siteID)
		}

		return nil
	}
}

func testAccCheckIncapsulaIncapRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_incap_rule" "testacc-terraform-incap-rule" {
  name = "%s"
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["%s"]
  enabled = true
  send_notifications = "true"
}`, incapRuleName, siteResourceName,
	)
}

func testAccCheckIncapsulaIncapRuleConfigEnabledFlagNotProvided(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_incap_rule" "testacc-terraform-incap-rule" {
  name = "%s"
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["%s"]
  send_notifications = "false"
}`, incapRuleName, siteResourceName,
	)
}

func testAccCheckIncapsulaIncapRuleConfigDisabled(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_incap_rule" "testacc-terraform-incap-rule" {
  name = "%s"
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["%s"]
  enabled = false
  send_notifications = "false"
}`, incapRuleName, siteResourceName,
	)
}
