package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
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
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		cspSite, err := client.GetCSPSite(siteId)
		if err != nil {
			return fmt.Errorf("Incapsula CSP Site Config doesn't exist for site ID %d", siteId)
		}
		if cspSite == nil || cspSite.Discovery != CSPDiscoveryOn {
			return fmt.Errorf("Incapsula CSP Site Config isn't on for site ID %d", siteId)
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

		resourceID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing CSP site config ID %v to int", rs.Primary.ID)
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])

		if siteID != resourceID {
			return "", fmt.Errorf("Incapsula CSP Site Config does not exist")
		}
		return fmt.Sprintf("%d", resourceID), nil
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

		apiID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to convert API ID from import command, actual value : %s, expected numeric id", rs.Primary.ID)
		}

		if apiID != siteIDInt {
			return fmt.Errorf("Resource %s for CSP site configuration has mismatch with IDs: Api Id %d, site ID %d", cspSiteConfigResourceType, apiID, siteIDInt)
		}

		cspSite, err := client.GetCSPSite(siteIDInt)
		fmt.Sprintf("Got CSP site config for site ID %d: %v", siteIDInt, cspSite)
		if err != nil && cspSite != nil && cspSite.Discovery != CSPDiscoveryOff {
			return fmt.Errorf("Resource %s for CSP site configuration: Api Id %d, site ID %d still exists", cspSiteConfigResourceType, apiID, siteIDInt)
		}
		return nil
	}
	return fmt.Errorf("Error finding the correct resource: %s", cspSiteConfigResourceType)
}

func testAccCheckCSPSiteConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id             = %s.id
		mode                = "monitor"
		email_addresses     = ["amiranc@imperva.com"]
		depends_on = ["%s"]
	}`,
		cspSiteConfigResourceType, cspSiteConfigName, siteResourceName, siteResourceName,
	)
}
