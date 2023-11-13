package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const deliveryRuleResourceName = "incapsula_delivery_rules_configuration.testacc-terraform-delivery_rules_configuration"
const redirectRuleName = "Redirect Delivery Rule Test"
const rewriteRuleName1 = "rewrite Delivery Rule Test 1"
const rewriteRuleName2 = "rewrite Delivery Rule Test 2"
const forwardRuleName1 = "forward Delivery Rule Test 1"
const forwardRuleName2 = "forward Delivery Rule Test 2"
const redirectCategory = "REDIRECT"
const rewriteCategory = "REWRITE"
const forwardCategory = "FORWARD"

func TestAccIncapsulaDeliveryRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDeliveryRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(deliveryRuleResourceName),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "category", redirectCategory),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.rule_name", redirectRuleName),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.action", "RULE_ACTION_REDIRECT"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.filter", "ASN == 1"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.response_code", "302"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.from", "/1"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.to", "/2"),
				),
			},
			{
				Config: testAccCheckIncapsulaDeliveryRewriteRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(deliveryRuleResourceName),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "category", rewriteCategory),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.rule_name", rewriteRuleName1),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.action", "RULE_ACTION_REWRITE_COOKIE"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.filter", "ASN == 2"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.rewrite_existing", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.from", "cookie1"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.to", "cookie2"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.add_if_missing", "false"),

					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.rule_name", rewriteRuleName2),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.enabled", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.action", "RULE_ACTION_REWRITE_HEADER"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.filter", "ASN == 3"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.rewrite_existing", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.from", "header1"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.to", "header2"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.add_if_missing", "false"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.header_name", "abc"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.response_code", "302"),
				),
			},
			{
				Config: testAccCheckIncapsulaDeliveryForwardRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(deliveryRuleResourceName),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "category", forwardCategory),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.rule_name", forwardRuleName1),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.action", "RULE_ACTION_FORWARD_TO_DC"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.dc_id", "1234"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.0.filter", "ASN == 1"),

					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.rule_name", forwardRuleName2),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.enabled", "true"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.action", "RULE_ACTION_FORWARD_TO_PORT"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.filter", "ASN == 1"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.port_forwarding_context", "[Use Header Name/Use Port Value]"),
					resource.TestCheckResourceAttr(deliveryRuleResourceName, "rule.1.port_forwarding_value", "1234"),
				),
			},
			{
				ResourceName:      deliveryRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetDeliveryRulesConfigurationImportString,
			},
		},
	})
}

func testAccCheckIncapsulaDeliveryRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_delivery_rules_configuration" {
			continue
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula Site ID does not exist")
		}
		category, ok := res.Primary.Attributes["category"]
		if !ok {
			return fmt.Errorf("Rule category does not exist for Site ID is : %s ", siteID)
		}

		deliveryRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, category)

		if diags != nil {
			log.Printf("[ERROR] Failed to read delivery rules in category %s for Site ID %s", category, siteID)
			return fmt.Errorf("failed to read delivery rules in category %s for Site ID %s\"", category, siteID)
		}

		if deliveryRulesListDTO == nil || deliveryRulesListDTO.Errors != nil || deliveryRulesListDTO.RulesList == nil || len(deliveryRulesListDTO.RulesList) != 0 {
			log.Printf("The DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
			return fmt.Errorf("the DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
		}
	}
	return nil
}

func testCheckIncapsulaDeliveryRuleExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Delivery Rule resource not found: %s", name)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID == "" {
			return fmt.Errorf("Incapsula Site ID does not exist ")
		}

		client := testAccProvider.Meta().(*Client)
		category, ok := res.Primary.Attributes["category"]
		deliveryRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, category)

		if !ok {
			return fmt.Errorf("Rule category : %s ,does not exist for Site ID is : %s ", siteID, category)
		}
		if deliveryRulesListDTO.Errors[0].Status != 200 {
			return fmt.Errorf("Incapsula Delivery Rule: %s (site id: %s) should have received 200 status code", category, siteID)
		}
		if diags != nil {
			return fmt.Errorf("Incapsula Delivery Rule: %s (site id: %s) does not exist", category, siteID)
		}

		return nil
	}
}
func testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_delivery_rules_configuration" "testacc-terraform-redirect-rules" {
  category = %s
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  rule {
  	rule_name = %s
  	filter = "ASN == 1"
	response_code = "302"
	action = "RULE_ACTION_REDIRECT"
  	enabled = true
    from = "/1"
    to = "/2"
	}`, redirectCategory, redirectRuleName,
	)
}
func testAccCheckIncapsulaDeliveryRewriteRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_delivery_rules_configuration" "testacc-terraform-ewrite-rules" {
  category = %s
  site_id = "${incapsula_site.testacc-terraform-site.id}"
 rule {
  	rule_name = %s
    filter = "ASN == 2"
    from = "cookie1"
    to = "cookie2"
    rewrite_existing = "true"
    add_if_missing = "false"
    action = "RULE_ACTION_REWRITE_COOKIE"
    enabled = "true"
 }
rule {
  	rule_name = %s
    filter = "ASN == 3"
    header_name = "abc"
    from = "header1"
    to = "header2"
    response_code = 302
    rewrite_existing = "true"
    add_if_missing = "false"
    action = "RULE_ACTION_REWRITE_HEADER"
    enabled = "true"
  }
}`, rewriteCategory, rewriteRuleName1, rewriteRuleName2,
	)
}
func testAccCheckIncapsulaDeliveryForwardRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_delivery_rules_configuration" "testacc-terraform-ewrite-rules" {
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  category = %s
  rule {
  	rule_name = %s
    filter = "ASN == 1"
    dc_id = 1234
    action = "RULE_ACTION_FORWARD_TO_DC"
    enabled = "true"
  }
rule {
  	rule_name = %s
	filter = "ASN == 1"
    port_forwarding_context = "[Use Header Name/Use Port Value]"
    port_forwarding_value = 1234
    action = "RULE_ACTION_FORWARD_TO_PORT"
    enabled = "true"
  }
}`, forwardCategory, forwardRuleName1, forwardRuleName2,
	)
}

func testAccGetDeliveryRulesConfigurationImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_delivery_rules_configuration" {
			continue
		}
		category, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("error parsing site_id %s to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%s/%d", category, siteID), nil
	}

	return "", fmt.Errorf("error finding Delivery Rule Resource")
}
