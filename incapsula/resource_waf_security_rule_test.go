package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const wafSecurityRuleNameBackdoor = "Example waf backdoor rule"
const wafSecurityRuleNameDDoS = "Example waf bot rule"
const wafSecurityRuleNameBots = "Example waf ddos rule"
const wafSecurityRuleResourceNameBackdoor = "incapsula_waf_security_rule.example-waf-backdoor-rule"
const wafSecurityRuleResourceNameBotAccessControl = "incapsula_waf_security_rule.example-waf-bot-access-control-rule"
const wafSecurityRuleResourceNameDDoS = "incapsula_waf_security_rule.example-waf-ddos-rule"

// Test all WAF security rule good configurations, one for ruleID that uses security_rule_action, one for ddos and bots (the three variations of param combinations)
func testAccCheckWAFSecurityRuleCreateGoodConfigBackdoor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyBackdoor,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfigBackdoor(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameBackdoor),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameBackdoor, "name", wafSecurityRuleNameBackdoor),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameBackdoor,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleCreateGoodConfigDDoS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyDDoS,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfigDDoS(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameDDoS),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameDDoS, "name", wafSecurityRuleResourceNameDDoS),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameDDoS,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleCreateGoodConfigBots(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyBots,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfigBots(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameBotAccessControl),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameBotAccessControl, "name", wafSecurityRuleResourceNameBotAccessControl),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameBotAccessControl,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleCreateBadConfigBackdoor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyBackdoor,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfigBackdoor(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameBackdoor),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameBackdoor, "name", wafSecurityRuleResourceNameBackdoor),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameBackdoor,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleCreateBadConfigDDoS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyDDoS,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfigDDoS(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameDDoS),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameDDoS, "name", wafSecurityRuleResourceNameDDoS),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameDDoS,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleCreateBadConfigBots(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroyBots,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfigBots(),
				Check: resource.ComposeTestCheckFunc(
					testCheckWAAFSecurityRuleExists(wafSecurityRuleResourceNameBotAccessControl),
					resource.TestCheckResourceAttr(wafSecurityRuleResourceNameBotAccessControl, "name", wafSecurityRuleNameBots),
				),
			},
			{
				ResourceName:      wafSecurityRuleResourceNameBotAccessControl,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWAFSecurityRuleID,
			},
		},
	})
}

func testAccCheckWAFSecurityRuleDestroyBackdoor(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_waf_security_rule" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula WAF rule ID " + backdoorRuleID + " does not exist")
		}

		err := "nil"
		if err == "nil" {
			return fmt.Errorf("Incapsula WAF security rule for site for domain: %s (site id: %s) still exists", backdoorRuleID, ruleID)
		}
	}

	return nil
}

func testAccCheckWAFSecurityRuleDestroyDDoS(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_waf_security_rule" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula WAF rule ID " + backdoorRuleID + " does not exist")
		}

		err := "nil"
		if err == "nil" {
			return fmt.Errorf("Incapsula WAF security rule for site for domain: %s (site id: %s) still exists", backdoorRuleID, ruleID)
		}
	}

	return nil
}

func testAccCheckWAFSecurityRuleDestroyBots(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_waf_security_rule" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula WAF rule ID " + backdoorRuleID + " does not exist")
		}

		err := "nil"
		if err == "nil" {
			return fmt.Errorf("Incapsula WAF security rule for site for domain: %s (site id: %s) still exists", backdoorRuleID, ruleID)
		}
	}

	return nil
}

func testCheckWAAFSecurityRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// each security rule will always exist as a part of the site.  Returning nil.
		return nil
	}
}

// Good and bad WAF Security Rule configs
func testAccCheckWAFSecurityRuleGoodConfigBackdoor() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleGoodConfigDDoS() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "api.threats.ddos.activation_mode.on"
  ddos_traffic_threshold = "5000"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleGoodConfigBots() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "true"
  challenge_suspected_bots = "true"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfigBackdoor() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "bad_action"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfigDDoS() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "bad_activation_mode"
  ddos_traffic_threshold = "1234"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfigBots() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "abc"
  challenge_suspected_bots = "abc"
}`, certificateName, siteResourceName,
	)
}
