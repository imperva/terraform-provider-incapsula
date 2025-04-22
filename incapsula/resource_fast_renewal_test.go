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
var siteV3NameFr = "test-cloudwaf-site-for-fast-renewal" + randomSuffix
var siteV3ResourceNameFr = siteV3NameFr
var manageCertSettingsResourceNameFr = "testacc-terraform-managed_certificate_settings"
var fastRenewalResourceName = "incapsula_fast_renewal" + randomSuffix

func TestAccFastRenewalCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_fast_renewal_test.TestAccFastRenewalCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFastRenewalCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastRenewalCertificateRequestConfig(t),
				Check: resource.ComposeTestCheckFunc(
					printState(),
					testCheckIncapsulaFastRenewalResourceAttributes("incapsula_fast_renewal."+fastRenewalResourceName),
					resource.TestCheckResourceAttrWith(
						"incapsula_fast_renewal."+fastRenewalResourceName,
						"id",
						func(val string) error {
							if _, err := strconv.Atoi(val); err != nil {
								return fmt.Errorf("expected fast_renewal_id to be an integer, got: %s", val)
							}
							return nil
						},
					),
				),
			},
			{
				ResourceName:      siteCertificateResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateFastRenewal,
			},
		},
	})
}

func testAccFastRenewalCertificateDestroy(s *terraform.State) error {
	return nil
}

func testAccFastRenewalCertificateRequestConfig(t *testing.T) string {
	res := fmt.Sprintf(`
	
   resource "incapsula_site_v3" "%s" {
			name = "%s"
	}
	resource incapsula_managed_certificate_settings "%s" {
    site_id = incapsula_site_v3.%s.id
    default_validation_method = "CNAME"
	}

	resource "incapsula_fast_renewal" "%s" {
  	site_id = incapsula_site_v3.%s.id
  	fast_renewal = true
  	depends_on = [ incapsula_managed_certificate_settings.%s ]
}

	`,
		siteV3ResourceNameFr, siteV3NameFr, manageCertSettingsResourceNameFr, siteV3ResourceNameFr, fastRenewalResourceName, siteV3ResourceNameFr, manageCertSettingsResourceNameFr,
	)
	return res
}

func testACCStateFastRenewal(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteCertificateResourceName {
			continue
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("error parsing site ID for import site certificate request. Value %s", rs.Primary.Attributes["site_id"])
		}

		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("error finding an Site certificate request\"")
}

func testCheckIncapsulaFastRenewalResourceAttributes(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("incapsula account resource not found: %s", name)
		}

		var fastRenewal = res.Primary.Attributes["fast_renewal"]
		var id = res.Primary.Attributes["id"]
		var siteId = res.Primary.Attributes["site_id"]

		if strings.Compare(fastRenewal, "true") != 0 {
			return fmt.Errorf("fast_renewal_fast_renewal expected true but got: %s", fastRenewal)
		}

		if _, err := strconv.Atoi(id); err != nil {
			return fmt.Errorf("expected fast_renewal_id to be an integer, got: %s", id)
		}
		if strings.Compare(siteId, id) != 0 {
			return fmt.Errorf("incapsula fast_renewal_site_id does not match fast_renewal_id : %s", siteId)
		}

		return nil
	}
}
