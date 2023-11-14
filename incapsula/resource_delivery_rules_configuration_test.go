package incapsula

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const deliveryRuleResourceName = "incapsula_delivery_rules_configuration"

const redirectRulesConfigurationName = "testacc-terraform-redirect-rules"
const redirectRulesResourceName = deliveryRuleResourceName + "." + redirectRulesConfigurationName

const rewriteRulesConfigurationName = "testacc-terraform-rewrite-rules"
const rewriteRulesResourceName = deliveryRuleResourceName + "." + rewriteRulesConfigurationName

const forwardRulesConfigurationName = "testacc-terraform-forward-rules"
const forwardRulesResourceName = deliveryRuleResourceName + "." + forwardRulesConfigurationName

func TestAccIncapsulaDeliveryRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDeliveryRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(redirectRulesResourceName, 1),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "category", "REDIRECT"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.rule_name", "Redirect Delivery Rule Test"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.action", "RULE_ACTION_REDIRECT"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.filter", "ASN == 1"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.response_code", "302"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.from", "*/1"),
					resource.TestCheckResourceAttr(redirectRulesResourceName, "rule.0.to", "$1/2"),
				),
			},
			{
				ResourceName:      redirectRulesResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetDeliveryRulesConfigurationImportString,
			},
			{
				Config: testAccCheckIncapsulaDeliveryRewriteRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(rewriteRulesResourceName, 2),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "category", "REWRITE"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.rule_name", "rewrite Delivery Rule Test 1"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.action", "RULE_ACTION_REWRITE_COOKIE"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.filter", "ASN == 2"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.rewrite_existing", "true"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.cookie_name", "cookie_1"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.from", "cookie1"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.to", "cookie2"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.0.add_if_missing", "false"),

					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.rule_name", "rewrite Delivery Rule Test 2"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.enabled", "true"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.action", "RULE_ACTION_REWRITE_HEADER"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.filter", "ASN == 3"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.rewrite_existing", "false"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.from", "header1"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.to", "header2"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.add_if_missing", "false"),
					resource.TestCheckResourceAttr(rewriteRulesResourceName, "rule.1.header_name", "abc"),
				),
			},
			{
				ResourceName:      rewriteRulesResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetDeliveryRulesConfigurationImportString,
			},
			{
				Config: testAccCheckIncapsulaDeliveryForwardRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(forwardRulesResourceName, 1),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "category", "FORWARD"),
					// resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.rule_name", "forward Delivery Rule Test 1"),
					// resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.enabled", "true"),
					// resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.action", "RULE_ACTION_FORWARD_TO_DC"),
					// resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.dc_id", "1234"),
					// resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.filter", "ASN == 1"),

					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.rule_name", "forward Delivery Rule Test 1"),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.action", "RULE_ACTION_FORWARD_TO_PORT"),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.filter", "ASN == 1"),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.port_forwarding_context", "Use Port Value"),
					resource.TestCheckResourceAttr(forwardRulesResourceName, "rule.0.port_forwarding_value", "1234"),
				),
			},
			{
				ResourceName:      forwardRulesResourceName,
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

func testCheckIncapsulaDeliveryRuleExists(name string, numRules int) resource.TestCheckFunc {
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
		if diags != nil {
			return fmt.Errorf("Incapsula Delivery Rule: %s (site id: %s) returned error", category, siteID)
		}
		if deliveryRulesListDTO.RulesList == nil || len(deliveryRulesListDTO.RulesList) != numRules {
			return fmt.Errorf("Incapsula Delivery Rule: %s (site id: %s) should have %d rules", category, siteID, numRules)
		}

		return nil
	}
}
func testAccCheckIncapsulaDeliveryRedirectRuleConfigBasic(t *testing.T) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  category = "%s"
  site_id = "1657653"
  rule {
  	rule_name = "%s"
  	filter = "ASN == 1"
	response_code = "302"
	action = "RULE_ACTION_REDIRECT"
  	enabled = true
    from = "*/1"
    to = "$1/2"
  }
}`, deliveryRuleResourceName, redirectRulesConfigurationName, "REDIRECT", "Redirect Delivery Rule Test",
	)
}
func testAccCheckIncapsulaDeliveryRewriteRuleConfigBasic(t *testing.T) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  category = "%s"
  site_id = "1657653"
  rule {
   	rule_name = "%s"
    filter = "ASN == 2"
	cookie_name = "cookie_1"
    from = "cookie1"
    to = "cookie2"
    add_if_missing = "false"
    action = "RULE_ACTION_REWRITE_COOKIE"
    enabled = "true"
  }
  rule {
  	rule_name = "%s"
    filter = "ASN == 3"
    header_name = "abc"
    from = "header1"
    to = "header2"
    rewrite_existing = "false"
    add_if_missing = "false"
    action = "RULE_ACTION_REWRITE_HEADER"
    enabled = "true"
  }
}`, deliveryRuleResourceName, rewriteRulesConfigurationName, "REWRITE", "rewrite Delivery Rule Test 1", "rewrite Delivery Rule Test 2",
	)
}
func testAccCheckIncapsulaDeliveryForwardRuleConfigBasic(t *testing.T) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  site_id = "1657653"
  category = "%s"
  rule {
  	rule_name = "%s"
	filter = "ASN == 1"
    port_forwarding_context = "Use Port Value"
    port_forwarding_value = "1234"
    action = "RULE_ACTION_FORWARD_TO_PORT"
    enabled = "true"
  }
}`, deliveryRuleResourceName, forwardRulesConfigurationName, "FORWARD", "forward Delivery Rule Test 1",
	)
}

func testAccGetDeliveryRulesConfigurationImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_delivery_rules_configuration" {
			continue
		}
		return rs.Primary.ID, nil
	}

	return "", fmt.Errorf("error finding Delivery Rule Resource")
}
