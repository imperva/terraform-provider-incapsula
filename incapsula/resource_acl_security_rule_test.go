package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strconv"
	"testing"
)

const aclSecurityRuleNameBlacklistedCountries = "Example ACL security rule - blacklisted_countries"
const aclSecurityRuleResourceNameBlacklistedCountries = "incapsula_acl_security_rule.example-global-blacklist-country-rule"

////////////////////////////////////////////////////////////////
// testAccCheckACLSecurityRuleCreate Tests
////////////////////////////////////////////////////////////////

func testAccCheckACLSecurityRuleCreateValidRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleGoodConfigBlacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceNameBlacklistedCountries, "name", aclSecurityRuleNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateACLSecurityRuleID,
			},
		},
	})
}

func testAccCheckACLSecurityRuleCreateInvalidRuleID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleInvalidRuleIDBlacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceNameBlacklistedCountries, "name", aclSecurityRuleNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateACLSecurityRuleID,
			},
		},
	})
}

func testAccCheckACLSecurityRuleCreateInvalidParams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityRuleExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckACLSecurityRuleInvalidParamBlacklistedCountries(),
				Check: resource.ComposeTestCheckFunc(
					testCheckACLSecurityRuleExists(aclSecurityRuleResourceNameBlacklistedCountries),
					resource.TestCheckResourceAttr(aclSecurityRuleResourceNameBlacklistedCountries, "name", aclSecurityRuleNameBlacklistedCountries),
				),
			},
			{
				ResourceName:      aclSecurityRuleResourceNameBlacklistedCountries,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateACLSecurityRuleID,
			},
		},
	})
}

////////////////////////////////////////////////////////////////
// testAccCheckSecurityRuleDestroy Tests
////////////////////////////////////////////////////////////////

func testAccCheckACLSecurityRuleDestroy(state *terraform.State) error {
	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_acl_security_rule" {
			continue
		}

		ruleID := res.Primary.ID
		if ruleID == "" {
			return fmt.Errorf("Incapsula acl security rule - rule ID (%s) does not exist", ruleID)
		}

		return fmt.Errorf("Incapsula acl security rule for site site_id (%s) still exists", ruleID)
	}

	return nil
}

func testCheckACLSecurityRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// each security rule will always exist as a part of the site.  Returning nil.
		return nil
	}
}

func testAccStateACLSecurityRuleID(s *terraform.State) (string, error) {
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
func testAccCheckACLSecurityRuleGoodConfigBlacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}

// Bad Security Rule Exception configs
func testAccCheckACLSecurityRuleInvalidRuleIDBlacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "bad_rule_id"
  countries = "AI,AN"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}

func testAccCheckACLSecurityRuleInvalidParamBlacklistedCountries() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf("%s%s", `
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "Bad_Value"
}`, securityRuleExceptionResourceNameBlacklistedCountries,
	)
}
