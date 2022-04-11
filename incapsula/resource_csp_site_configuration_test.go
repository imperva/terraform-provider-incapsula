package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"strings"
	"testing"
)

const cspSiteConfigResourceType = "incapsula_csp_site_configuration"
const cspSiteConfigResource = cspSiteConfigResourceType + "." + cspSiteConfigName
const cspSiteConfigName = "testacc-terraform-csp-site-config"

func TestAccIncapsulaCSPSiteConfig_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_csp_site_configuration_test.TestAccIncapsulaCSPSiteConfig_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateCSPSiteConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCSPSiteConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSPSiteConfigExists(cspSiteConfigResource),
					resource.TestCheckResourceAttr(cspSiteConfigResource, "mode", "monitor"),
					resource.TestCheckResourceAttr(cspSiteConfigResource, "email_addresses.#", "1"),
					resource.TestCheckResourceAttr(cspSiteConfigResource, "email_addresses.0", "amiranc@imperva.com"),
				),
			},
			{
				ResourceName:      cspSiteConfigResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateCSPSiteConfigID,
			},
		},
	})
}

func testCheckCSPSiteConfigExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula CSP Site Config resource not found: %s", name)
		}

		keyParts := strings.Split(res.Primary.ID, "/")
		if len(keyParts) != 2 {
			return fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", res.Primary.ID)
		}
		accountID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric ID", res.Primary.ID)
		}
		siteID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			fmt.Errorf("failed to convert account ID from import command, actual value: %s, expected numeric ID", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		cspSite, err := client.GetCSPSite(accountID, siteID)
		if err != nil {
			return fmt.Errorf("Incapsula CSP Site Config doesn't exist for site ID %d", siteID)
		}
		if cspSite == nil || cspSite.Discovery != CSPDiscoveryOn {
			return fmt.Errorf("Incapsula CSP Site Config isn't on for site ID %d", siteID)
		}

		return nil
	}
}

func testACCStateCSPSiteConfigID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		fmt.Errorf("Resource: %v", rs)
		if rs.Type != cspSiteConfigResourceType {
			continue
		}

		keyParts := strings.Split(rs.Primary.ID, "/")
		if len(keyParts) != 2 {
			return "", fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", rs.Primary.ID)
		}
		accountID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric ID", rs.Primary.ID)
		}

		siteID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert account ID from import command, actual value: %s, expected numeric ID", rs.Primary.ID)
		}
		resourceID := fmt.Sprintf("%d/%d", accountID, siteID)

		schemaSiteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		schemaAccountID, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		newID := fmt.Sprintf("%d/%d", schemaAccountID, schemaSiteID)

		if strings.Compare(newID, resourceID) != 0 {
			// if newID != resourceID {
			return "", fmt.Errorf("Incapsula CSP Site Config does not exist")
		}
		return resourceID, nil
	}
	return "", fmt.Errorf("Error finding correct resource %s", cspSiteConfigResourceType)
}

func testACCStateCSPSiteConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != cspSiteConfigResourceType {
			continue
		}

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Parameter site_id was not found in resource %s", cspSiteConfigResourceType)
		}
		siteIDInt, err := strconv.Atoi(siteID)
		if err != nil {
			return fmt.Errorf("failed to convert site ID from import command, actual value : %s, expected numeric id", siteID)
		}
		accountID := rs.Primary.Attributes["account_id"]
		accountIDInt, err := strconv.Atoi(accountID)
		if err != nil {
			return fmt.Errorf("failed to convert site ID from import command, actual value : %s, expected numeric id", siteID)
		}

		cspSite, err := client.GetCSPSite(accountIDInt, siteIDInt)
		fmt.Sprintf("Got CSP site config for site ID %d: %v", siteIDInt, cspSite)
		if err != nil && cspSite != nil && cspSite.Discovery != CSPDiscoveryOff {
			return fmt.Errorf("Resource %s for CSP site configuration: Api Id %s, site ID %d still exists", cspSiteConfigResourceType, rs.Primary.ID, siteIDInt)
		}
		return nil
	}
	return fmt.Errorf("Error finding the correct resource: %s", cspSiteConfigResourceType)
}

func testAccCheckCSPSiteConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		account_id          = %s.account_id
		site_id             = %s.id
		mode                = "monitor"
		email_addresses     = ["amiranc@imperva.com"]
		depends_on = ["%s"]
	}`,
		cspSiteConfigResourceType, cspSiteConfigName, siteResourceName, siteResourceName, siteResourceName,
	)
}
