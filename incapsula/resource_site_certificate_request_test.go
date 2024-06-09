package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const siteCertificateResourceName = "incapsula_site_certificate_request"
const siteCertificateResource = siteCertificateResourceName + "." + siteCertificateConfigName
const siteCertificateConfigName = "testacc-terraform-site_certificate_request"

func TestAccSiteCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_site_certificate_request_test.TestAccSiteCertificate_Basic")
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
				ImportStateIdFunc: testACCStateSiteCertificateRequest,
			},
		},
	})
}

func testAccSiteCertificateDestroy(s *terraform.State) error {
	return nil
}

func testAccSiteCertificateRequestConfig(t *testing.T) string {
	return fmt.Sprintf(`
	
	resource"%s""%s"{
    site_id = 722883409
    default_validation_method = "DNS"
	}`,
		siteCertificateResourceName, siteCertificateConfigName,
	)
}

func testACCStateSiteCertificateRequest(s *terraform.State) (string, error) {
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
