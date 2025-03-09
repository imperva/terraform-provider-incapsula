package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const securityRuleExceptionNameBlacklistedCountries = "Example security rule exception - blacklisted_countries"
const securityRuleExceptionResourceNameBlacklistedCountries = "incapsula_security_rule_exception.example-waf-blacklisted-countries-rule-exception"

////////////////////////////////////////////////////////////////
// AccCheckAddSecurityRuleException Tests
////////////////////////////////////////////////////////////////

func testAccCheckSecurityRuleExceptionCreateValidRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionGoodConfigBlacklistedCountries(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceNameBlacklistedCountries, "name", securityRuleExceptionNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckSecurityRuleExceptionCreateInvalidRuleID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionInvalidRuleIDBlacklistedCountries(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceNameBlacklistedCountries, "name", securityRuleExceptionNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckSecurityRuleExceptionCreateInvalidParams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionInvalidParamBlacklistedCountries(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceNameBlacklistedCountries, "name", securityRuleExceptionNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

////////////////////////////////////////////////////////////////
// testAccCheckSecurityRuleExceptionDestroy Tests
////////////////////////////////////////////////////////////////

func testAccCheckSecurityRuleExceptionDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_security_rule_exception" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula security rule exception - rule ID (%s) does not exist", ruleID)
		}

		err := "nil"
		if err == "nil" {
			return fmt.Errorf("Incapsula security rule exception for site site_id (%s) still exists", ruleID)
		}
	}

	return nil
}

func testCheckSecurityRuleExceptionExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula security rule exception resource not found: %s", name)
		}

		ruleID := "api.acl.blacklisted_countries"
		siteID := res.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula security rule exception ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)
		siteStatusResponse, err := client.ListSecurityRuleExceptions(siteID, ruleID)
		if err != nil {
			return fmt.Errorf("ListSecurityRuleExceptions Error for site_id (%s) and rule_id (%s) %s", siteID, ruleID, err)
		}

		if siteStatusResponse == nil {
			return fmt.Errorf("Incapsula security rule exception for site id (%s) and rule_id (%s) does not exist", siteID, ruleID)
		}

		return nil
	}
}

func testAccStateSecurityRuleExceptionID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_waf_security_rule" {
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

// Good Security Rule Exception configs
func testAccCheckACLSecurityRuleExceptionGoodConfigBlacklistedCountries(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  urls="/myurl,/myurl2"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}

// Bad Security Rule Exception configs
func testAccCheckACLSecurityRuleExceptionInvalidRuleIDBlacklistedCountries(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "bad_rule_id"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  urls="/myurl,/myurl2"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}

func testAccCheckACLSecurityRuleExceptionInvalidParamBlacklistedCountries(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3."
  urls="/myurl,myurl2"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}
