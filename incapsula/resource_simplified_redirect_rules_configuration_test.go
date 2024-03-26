package incapsula

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const deliverySimplifiedRedirectRulesResourceName = "incapsula_simplified_redirect_rules_configuration"
const simplifiedRedirectRulesConfigurationName = "testacc-terraform-simplified-redirect-rules"
const simplifiedRedirectRulesResourceName = deliverySimplifiedRedirectRulesResourceName + "." + simplifiedRedirectRulesConfigurationName

func TestAccIncapsulaDeliverySimplifiedRedirectRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDeliverySimplifiedRedirectRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSimplifiedRedirectRuleConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSimplifiedRedirectRuleExists(simplifiedRedirectRulesResourceName, 1),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.rule_name", "Simplified Redirect Delivery Rule Test 2"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.enabled", "false"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.response_code", "302"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.from", "/1"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.to", "$scheme://www.example.com/$city"),
				),
			},
			{
				Config: testAccCheckSimplifiedRedirectRulesDiff(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSimplifiedRedirectRuleExists(simplifiedRedirectRulesResourceName, 1),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.rule_name", "Simplified Redirect Delivery Rule Test 2"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.from", "/1"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.to", "$scheme://www.example.com/$city"),
				),
			},
			{
				Config: testAccCheckMultipleSimplifiedRedirectRulesConfig(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSimplifiedRedirectRuleExists(simplifiedRedirectRulesResourceName, 5),

					//Currently we use a list, so result is ordered
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.rule_name", "Simplified Redirect Delivery Rule Test 1"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.from", "/1"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.0.to", "$scheme://www.example.com/$city"),

					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.1.rule_name", "Simplified Redirect Delivery Rule Test 2"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.1.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.1.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.1.from", "/2"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.1.to", "$scheme://www.example.com/$city"),

					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.2.rule_name", "Simplified Redirect Delivery Rule Test 3"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.2.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.2.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.2.from", "/3"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.2.to", "$scheme://www.example.com/$city"),

					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.3.rule_name", "Simplified Redirect Delivery Rule Test 4"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.3.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.3.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.3.from", "/4"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.3.to", "$scheme://www.example.com/$city"),

					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.4.rule_name", "Simplified Redirect Delivery Rule Test 5"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.4.enabled", "true"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.4.response_code", "301"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.4.from", "/5"),
					resource.TestCheckResourceAttr(simplifiedRedirectRulesResourceName, "rule.4.to", "$scheme://www.example.com/$city"),
				),
			},
			{
				ResourceName:      simplifiedRedirectRulesResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetSimplifiedRedirectRulesImportString,
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
			log.Printf("[ERROR] Failed to read simplified redirect rules in category %s for Site ID %s", category, siteID)
			return fmt.Errorf("failed to read simplified redirect rules in category %s for Site ID %s\"", category, siteID)
		}

		if deliveryRulesListDTO == nil || deliveryRulesListDTO.Errors != nil || deliveryRulesListDTO.RulesList == nil || len(deliveryRulesListDTO.RulesList) != 0 {
			log.Printf("The DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
			return fmt.Errorf("the DTO response shouldnt be \"empty\" or \"nil\" or has errors site ID: %s", siteID)
		}
	}
	return nil
}

func testAccCheckSimplifiedRedirectRuleConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
rule {
  	rule_name = "%s"
	response_code = 302
  	enabled = false
    from = "/1"
    to = "$scheme://www.example.com/$city"
  }
}`, deliverySimplifiedRedirectRulesResourceName, simplifiedRedirectRulesConfigurationName, siteResourceName, "Simplified Redirect Delivery Rule Test 2",
	)
}
func testAccCheckSimplifiedRedirectRulesDiff(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
rule {
  	rule_name = "%s"
	response_code = 301
  	enabled = true
    from = "/1"
    to = "$scheme://www.example.com/$city"
  }
}`, deliverySimplifiedRedirectRulesResourceName, simplifiedRedirectRulesConfigurationName, siteResourceName, "Simplified Redirect Delivery Rule Test 2",
	)
}
func testAccCheckMultipleSimplifiedRedirectRulesConfig(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
	rule {
		rule_name = "Simplified Redirect Delivery Rule Test 1"
		response_code = 301
		enabled = true
		from = "/1"
		to = "$scheme://www.example.com/$city"
	}
	rule {
		rule_name = "Simplified Redirect Delivery Rule Test 2"
		response_code = 301
		enabled = true
		from = "/2"
		to = "$scheme://www.example.com/$city"
	}
	rule {
		rule_name = "Simplified Redirect Delivery Rule Test 3"
		response_code = 301
		enabled = true
		from = "/3"
		to = "$scheme://www.example.com/$city"
	}
	rule {
		rule_name = "Simplified Redirect Delivery Rule Test 4"
		response_code = 301
		enabled = true
		from = "/4"
		to = "$scheme://www.example.com/$city"
	}
	rule {
		rule_name = "Simplified Redirect Delivery Rule Test 5"
		response_code = 301
		enabled = true
		from = "/5"
		to = "$scheme://www.example.com/$city"
	}
}`, deliverySimplifiedRedirectRulesResourceName, simplifiedRedirectRulesConfigurationName, siteResourceName,
	)
}
func testAccGetSimplifiedRedirectRulesImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_simplified_redirect_rules_configuration" {
			continue
		}
		return rs.Primary.ID, nil
	}

	return "", fmt.Errorf("error finding simplified redirect Rule Resource")
}

func testCheckSimplifiedRedirectRuleExists(name string, numRules int) resource.TestCheckFunc {
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
		deliveryRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, "SIMPLIFIED_REDIRECT")

		if !ok {
			return fmt.Errorf("Rule category : %s ,does not exist for Site ID is : SIMPLIFIED_REDIRECT ", siteID)
		}
		if diags != nil {
			return fmt.Errorf("Incapsula Delivery Rule: SIMPLIFIED_REDIRECT (site id: %s) returned error", siteID)
		}
		if deliveryRulesListDTO.RulesList == nil || len(deliveryRulesListDTO.RulesList) != numRules {
			return fmt.Errorf("Incapsula Delivery Rule: SIMPLIFIED_REDIRECT (site id: %s) should have %d rules", siteID, numRules)
		}

		return nil
	}
}
