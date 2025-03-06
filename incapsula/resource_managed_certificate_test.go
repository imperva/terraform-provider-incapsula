package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
	"time"
)

var siteV3ResourceNameForManagedCert = "test-cloudwaf-site-for-site-cert" + strconv.FormatInt(time.Now().UnixNano()%99999, 10)
var siteV3Name = "test site " + strconv.FormatInt(time.Now().UnixNano()%99999, 10)

const siteCertificateResourceName = "incapsula_managed_certificate_settings"
const siteCertificateResource = siteCertificateResourceName + "." + siteCertificateConfigName
const siteCertificateConfigName = "testacc-terraform-managed_certificate_settings"

func TestAccSiteCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_managed_certificate_test.TestAccSiteCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSiteCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteCertificateRequestConfig(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(siteCertificateResource, "default_validation_method", "DNS"),
				),
			},
			{
				ResourceName:      siteCertificateResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateManagedCertificate,
			},
		},
	})
}

func testAccSiteCertificateDestroy(s *terraform.State) error {
	return nil
}

func testAccSiteCertificateRequestConfig(t *testing.T) string {
	res := fmt.Sprintf(`
	
   resource "imperva_site_v3" "%s" {
			name = "%s"
	}
	resource "%s" "%s" {
    site_id = imperva_site_v3.%s.id
    default_validation_method = "DNS"
	}`,
		siteV3ResourceNameForManagedCert, siteV3Name, siteCertificateResourceName, siteCertificateConfigName, siteV3ResourceNameForManagedCert,
	)
	return res
}

func testACCStateManagedCertificate(s *terraform.State) (string, error) {
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
