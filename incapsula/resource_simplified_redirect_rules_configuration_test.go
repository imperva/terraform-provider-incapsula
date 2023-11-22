package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
)

const deliverySimplifiedRedirectRulesResourceName = "incapsula_simplified_redirect_rules_configuration"
const simplifiedRedirectRulesConfigurationName = "testacc-terraform-simplified-redirect-rules"
const simplifiedRedirectRulesResourceName = simplifiedRedirectRulesConfigurationName + "." + deliverySimplifiedRedirectRulesResourceName

func TestAccIncapsulaDeliverySimplifiedRedirectRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDeliverySimplifiedRedirectRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDeliverySimplifiedRedirectRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDeliveryRuleExists(deliverySimplifiedRedirectRulesResourceName, 2),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.rule_name", "Simplified Redirect Delivery Rule Test 1"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.action", "RULE_ACTION_SIMPLIFIED_REDIRECT"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.response_code", "302"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.from", "*/1"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.0.to", "$1/2"),

					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.rule_name", "Simplified Redirect Delivery Rule Test 2"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.action", "RULE_ACTION_SIMPLIFIED_REDIRECT"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.response_code", "302"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.from", "*/1"),
					resource.TestCheckResourceAttr(deliverySimplifiedRedirectRulesResourceName, "rule.1.to", "$1/2"),
				),
			},
			{
				ResourceName:      simplifiedRedirectRulesResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetDeliverySimplifiedRedirectRulesConfigurationImportString,
			},
		},
	})
}

func testAccCheckIncapsulaDeliverySimplifiedRedirectRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	category := "SIMPLIFIED_REDIRECT"

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_simplified_redirect_rules_configuration" {
			continue
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula Site ID does not exist")
		}

		deliveryRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, category)

		if diags != nil {
			log.Printf("[ERROR] Failed to read delivery simplified redirect rules in category %s for Site ID %s", category, siteID)
			return fmt.Errorf("failed to read delivery simplified redirect rules in category %s for Site ID %s\"", category, siteID)
		}

		if deliveryRulesListDTO == nil || deliveryRulesListDTO.Errors != nil || deliveryRulesListDTO.RulesList == nil || len(deliveryRulesListDTO.RulesList) != 0 {
			log.Printf("The DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
			return fmt.Errorf("the DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
		}
	}
	return nil
}

func testAccCheckIncapsulaDeliverySimplifiedRedirectRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  rule {
  	rule_name = "%s"
	response_code = "302"
	action = "RULE_ACTION_SIMPLIFIED_REDIRECT"
  	enabled = true
    from = "*/1"
    to = "$1/2"
  }
rule {
  	rule_name = "%s"
	response_code = "302"
	action = "RULE_ACTION_SIMPLIFIED_REDIRECT"
  	enabled = false
    from = "*/1"
    to = "$1/2"
  }
}`, deliverySimplifiedRedirectRulesResourceName, simplifiedRedirectRulesConfigurationName, "Simplified Redirect Delivery Rule Test 1", "Simplified Redirect Delivery Rule Test 2",
	)
}
func testAccGetDeliverySimplifiedRedirectRulesConfigurationImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_simplified_redirect_rules_configuration" {
			continue
		}
		return rs.Primary.ID, nil
	}

	return "", fmt.Errorf("error finding Delivery simplified redirect Rule Resource")
}
