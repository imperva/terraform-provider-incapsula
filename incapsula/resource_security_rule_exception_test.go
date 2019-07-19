package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"
)

const securityRuleExceptionName_blacklistedCountries = "Example security rule exception - blacklisted_countries"
const securityRuleExceptionResourceName_blacklistedCountries = "incapsula_security_rule_exception.example-waf-blacklisted-countries-rule-exception"

////////////////////////////////////////////////////////////////
// AccCheckAddSecurityRuleException Tests
////////////////////////////////////////////////////////////////

func testAccCheckSecurityRuleExceptionCreate_validRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionGoodConfig_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceName_blacklistedCountries, "name", securityRuleExceptionName_blacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceName_blacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckSecurityRuleExceptionCreate_invalidRuleId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionInvalidRuleID_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceName_blacklistedCountries, "name", securityRuleExceptionName_blacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceName_blacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateSecurityRuleExceptionID,
			},
		},
	})
}

func testAccCheckSecurityRuleExceptionCreate_invalidParams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleExceptionInvalidParam_blacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityRuleExceptionExists(securityRuleExceptionResourceName_blacklistedCountries),
					resource.TestCheckResourceAttr(securityRuleExceptionResourceName_blacklistedCountries, "name", securityRuleExceptionName_blacklistedCountries),
				),
			},
			{
				ResourceName:      securityRuleExceptionResourceName_blacklistedCountries,
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
			return fmt.Errorf("ListSecurityRuleExceptions Error for site_id (%s) and rule_id (%s) %s\n", siteID, ruleID, err)
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
func testAccCheckACLSecurityRuleExceptionGoodConfig_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}

// Bad Security Rule Exception configs
func testAccCheckACLSecurityRuleExceptionInvalidRuleID_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "bad_rule_id"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}

func testAccCheckACLSecurityRuleExceptionInvalidParam_blacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3."
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,myurl2"
}`, securityRuleExceptionResourceName_blacklistedCountries,
	)
}
