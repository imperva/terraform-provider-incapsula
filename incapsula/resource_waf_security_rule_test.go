package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const wafSecurityRuleNameBackdoor = "Example waf backdoor rule"
const wafSecurityRuleNameDDoS = "Example waf bot rule"
const wafSecurityRuleNameBots = "Example waf ddos rule"
const wafSecurityRuleResourceNameBackdoor = "incapsula_waf_security_rule.example-waf-backdoor-rule"
const wafSecurityRuleResourceNameBotAccessControl = "incapsula_waf_security_rule.example-waf-bot-access-control-rule"
const wafSecurityRuleResourceNameDDoS = "incapsula_waf_security_rule.example-waf-ddos-rule"

// Test all WAF security rule good configurations, one for ruleID that uses security_rule_action, one for ddos and bots (the three variations of param combinations)
func testAccCheckWAFSecurityRuleCreate_goodConfig_backdoor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_backdoor,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfig_backdoor(),
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

func testAccCheckWAFSecurityRuleCreate_goodConfig_ddos(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_ddos,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfig_ddos(),
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

func testAccCheckWAFSecurityRuleCreate_goodConfig_bots(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_bots,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleGoodConfig_bots(),
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

func testAccCheckWAFSecurityRuleCreate_badConfig_backdoor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_backdoor,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfig_backdoor(),
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

func testAccCheckWAFSecurityRuleCreate_badConfig_ddos(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_ddos,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfig_ddos(),
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

func testAccCheckWAFSecurityRuleCreate_badConfig_bots(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWAFSecurityRuleDestroy_bots,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWAFSecurityRuleBadConfig_bots(),
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

func testAccCheckWAFSecurityRuleDestroy_backdoor(state *terraform.State) error {
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

func testAccCheckWAFSecurityRuleDestroy_ddos(state *terraform.State) error {
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

func testAccCheckWAFSecurityRuleDestroy_bots(state *terraform.State) error {
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
func testAccCheckWAFSecurityRuleGoodConfig_backdoor() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleGoodConfig_ddos() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "api.threats.ddos.activation_mode.on"
  ddos_traffic_threshold = "5000"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleGoodConfig_bots() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "true"
  challenge_suspected_bots = "true"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfig_backdoor() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "bad_action"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfig_ddos() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "bad_activation_mode"
  ddos_traffic_threshold = "1234"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckWAFSecurityRuleBadConfig_bots() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "abc"
  challenge_suspected_bots = "abc"
}`, certificateName, siteResourceName,
	)
}
