package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

var randomSuffix = strconv.FormatInt(time.Now().UnixNano()%99999, 10)
var siteV3NameFr = "v3site-short-renewal-cycle-test" + randomSuffix
var siteV3ResourceNameFr = siteV3NameFr
var manageCertSettingsResourceNameFr = "cert-settings-short-renewal-cycle-test" + randomSuffix
var shortRenewalCycleResourceName = "short-renewal-cycle-test" + randomSuffix

func TestAccShortRenewalCycleCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_short_renewal_cycle_test.TestAccShortRenewalCycleCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccShortRenewalCycleCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccShortRenewalCycleCertificateRequestConfig(t),
				Check: resource.ComposeTestCheckFunc(
					printState(),
					testCheckIncapsulaShortRenewalCycleResourceAttributes("incapsula_short_renewal_cycle."+shortRenewalCycleResourceName),
					resource.TestCheckResourceAttrWith(
						"incapsula_short_renewal_cycle."+shortRenewalCycleResourceName,
						"id",
						func(val string) error {
							if _, err := strconv.Atoi(val); err != nil {
								return fmt.Errorf("expected short_renewal_cycle_id to be an integer, got: %s", val)
							}
							return nil
						},
					),
				),
			},
			{
				ResourceName:      "incapsula_short_renewal_cycle." + shortRenewalCycleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateShortRenewalCycle,
			},
		},
	})
}

func testAccShortRenewalCycleCertificateDestroy(s *terraform.State) error {
	return nil
}

func testAccShortRenewalCycleCertificateRequestConfig(t *testing.T) string {
	res := fmt.Sprintf(`
   resource "incapsula_site_v3" "%s" {
			name = "%s"
	}
	resource incapsula_managed_certificate_settings "%s" {
    site_id = incapsula_site_v3.%s.id
    default_validation_method = "CNAME"
	}

	resource "incapsula_short_renewal_cycle" "%s" {
  	site_id = incapsula_site_v3.%s.id
  	short_renewal_cycle = true
	managed_certificate_settings_id = incapsula_managed_certificate_settings.%s.id
}
	`,
		siteV3ResourceNameFr, siteV3NameFr, manageCertSettingsResourceNameFr, siteV3ResourceNameFr, shortRenewalCycleResourceName, siteV3ResourceNameFr, manageCertSettingsResourceNameFr,
	)
	return res
}

func testACCStateShortRenewalCycle(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_short_renewal_cycle" {
			continue
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("error parsing site id for import short renewal cycle. Value %s", rs.Primary.Attributes["site_id"])
		}

		managedCertSettingsId, err := strconv.Atoi(rs.Primary.Attributes["managed_certificate_settings_id"])
		if err != nil {
			return "", fmt.Errorf("error parsing site id for import short renewal cycle. Value %s", rs.Primary.Attributes["managed_certificate_settings_id"])
		}

		return fmt.Sprintf("%d/%d", siteID, managedCertSettingsId), nil
	}
	return "", fmt.Errorf("error finding a short renewal cycle resource\"")
}

func testCheckIncapsulaShortRenewalCycleResourceAttributes(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("incapsula account resource not found: %s", name)
		}

		var shortRenewalCycle = res.Primary.Attributes["short_renewal_cycle"]
		var id = res.Primary.Attributes["id"]
		var siteId = res.Primary.Attributes["site_id"]
		var managedCertSettingsId = res.Primary.Attributes["managed_certificate_settings_id"]

		if strings.Compare(shortRenewalCycle, "true") != 0 {
			return fmt.Errorf("short_renewal_cycle expected true but got: %s", shortRenewalCycle)
		}

		if _, err := strconv.Atoi(id); err != nil {
			return fmt.Errorf("expected short_renewal_cycle_id to be an integer, got: %s", id)
		}
		if strings.Compare(siteId, id) != 0 {
			return fmt.Errorf("incapsula short_renewal_cycle.site_id does not match short_renewal_cycle_id : %s", siteId)
		}
		if strings.Compare(managedCertSettingsId, id) != 0 {
			return fmt.Errorf("incapsula short_renewal_cycle.managed_certificate_settings_id does not match short_renewal_cycle_id : %s", managedCertSettingsId)
		}
		return nil
	}
}
