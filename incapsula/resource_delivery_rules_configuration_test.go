package incapsula

import (
	"fmt"
	"log"
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
			Config : testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t), // working json validation
			Check: resource.ComposeTestCheckFunc(
				testCheckIncapsulaIncapRuleExists(deliveryRuleResourceName),
				resource.TestCheckResourceAttr(redirectRuleName, "name", redirectRuleName),
				resource.TestCheckResourceAttr(incapRuleResourceName, "enabled", "true"),
			),
		},
		}
	}
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

		deliveryRulesListDTO, diags := client.ReadIncapRulePriorities(siteID, category)

		if diags != nil {
			log.Printf("[ERROR] Failed to read delivery rules in category %s for Site ID %s", category, siteID)
			return fmt.Errorf("Failed to read delivery rules in category %s for Site ID %s\", category, siteID")
		}

		if deliveryRulesListDTO != nil && deliveryRulesListDTO.Errors != nil && deliveryRulesListDTO.Errors[0].Status == 404 {
			log.Printf("[INFO] Incapsula Site with ID %s has already been deleted\n", siteID)
			return nil
		}

	}

	return nil
}
func testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_delivery_rules_configuration" "testacc-terraform-redirect-rules" {
  category = %s
  ite_id = "${incapsula_site.testacc-terraform-site.id}"
  rule {
  	rule_name = %s
  	action = "RULE_ACTION_ALERT"
  	filter = "ASN == 1"
	response_code = "302"
	action = "RULE_ACTION_REDIRECT"
  	enabled = true
    from = "/1"
    to = "/2"
	}`,redirectCategory, redirectRuleName,
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
    rule_name = "New delivery rule"
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
    rule_name = "New delivery rule"
    action = "RULE_ACTION_REWRITE_HEADER"
    enabled = "true"
  }
}`,rewriteCategory, rewriteRuleName1, rewriteRuleName2,
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
    rule_name = "New delivery rule",
    action = "RULE_ACTION_FORWARD_TO_DC"
    enabled = "true"
  }
rule {
  	rule_name = %s
	filter = "ASN == 1"
    port_forwarding_context = "[Use Header Name/Use Port Value]"
    port_forwarding_value = 1234
    rule_name = "New delivery rule"
    action = "RULE_ACTION_FORWARD_TO_PORT"
    enabled = "true"
  }
}`,forwardCategory ,forwardRuleName1, forwardRuleName2,
	)
}
