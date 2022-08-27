package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const siteTlsSettingsResourceName = "incapsula_site_tls_settings"
const siteTlsSettingsResource = siteTlsSettingsResourceName + "." + siteTlsSettingsgName
const siteTlsSettingsgName = "testacc-terraform-site_tls_settings"

func TestAccIncapsulaSiteTlsSettings_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_site_tls_settings_test.TestAccIncapsulaSiteTlsSettings_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSiteTlsSettingsBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSiteTlsSettingsExists(siteTlsSettingsResource),
					resource.TestCheckResourceAttr(siteTlsSettingsResource, "mandatory", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "ports", "[12,100,305]"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "is_ports_exception", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "hosts", "[\"test.com\", \"secondtest.au\"]"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "is_hosts_exception", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "fingerprints", "[\"fingerprint\"]"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "forward_to_origin", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "header_name", "something"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "header_value", "COMMON_NAME"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "is_disable_session_resumption", "trueE"),
				),
			},
			{
				ResourceName:      siteTlsSettingsResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiteTlsSettinggID,
			},
		},
	})
}

func testACCStateSiteTlsSettinggID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteTlsSettingsResourceName {
			continue
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])

		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int in Site Monitoring resource test", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("Error finding site_id argument in Site TLS Settings resource test")
}

func testCheckSiteTlsSettingsExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Site TLS settings resource not found: %s", name)
		}
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.GetSiteTlsSettings(siteId)
		if err != nil {
			fmt.Errorf("Incapsula Site TLS settings doesn't exist")
		}

		return nil
	}
}

func testAccCheckSiteTlsSettingsBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		depends_on                       = ["%s"]
		mandatory                        = true
		ports                            = [12,100,305]
		is_ports_exception               = true
		hosts                            = ["test.com", "secondtest.au"]
		is_hosts_exception               = true
		fingerprints                     = ["fingerprint"]
		forward_to_origin                = true
		header_name                      = "something"
		header_value                     = "COMMON_NAME"
		is_disable_session_resumption    = true
	}`,
		siteTlsSettingsResourceName, siteTlsSettingsgName, siteResourceName,
	)
}
