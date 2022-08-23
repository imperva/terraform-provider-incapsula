package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	//"strconv"
	"testing"
)

const siteMtlsCrtificateAssociationResourceName = "incapsula_site_mtls_certificate_association"
const siteMtlsCrtificateAssociationResource = siteMtlsCrtificateAssociationResourceName + "." + siteMtlsCrtificateAssociationName
const siteMtlsCrtificateAssociationName = "testacc-terraform-site-mtls-certificate-association"

func TestAccIncapsulaSiteMtlsCertificateAssociation_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaSiteMtlsCertificateAssociation_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateSiteMtlsCertificateAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSiteMtlsCertificateAssociationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSiteMtlsCertificateAssociationExists(),
				),
			},
		},
	})
}

//todo KATRIN:
//check destroy!!!!!!

func testACCStateSiteMtlsCertificateAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteMtlsCrtificateAssociationResourceName {
			continue
		}
		siteID, ok := rs.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula mTLS certificate is not assigned to Site ID %s (cannot find mandatory parameter site_id)", siteID)
		}

		_, err := client.GetMTLSCertificate(siteID)
		if err == nil {
			return fmt.Errorf("Resource %s with siteID ID %s still exists", siteMtlsCrtificateAssociationResourceName, siteID)
		}
	}
	return nil
}

//	check import

func testACCStateSiteMtlsCertificateAssociationID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteMtlsCrtificateAssociationResourceName {
			continue
		}
		certificateID, ok := rs.Primary.Attributes["certificate_id"]
		if !ok {
			return "", fmt.Errorf("Incapsula mTLS certificate ID %s to Site association doesn't exist (cannot find mandatory parameter certificate_id)", certificateID)
		}
		siteID, ok := rs.Primary.Attributes["site_id"]
		if !ok {
			return "", fmt.Errorf("Incapsula mTLS certificate ID %s is not assigned to Site ID %sn (cannot find mandatory parameter site_id)", siteID, certificateID)
		}
	}
	return "", fmt.Errorf("Error finding an mTLS Imperva to Origin Certificate\"")
}

func testCheckSiteMtlsCertificateAssociationExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, resource := range state.RootModule().Resources {
			log.Printf("\n%v", resource)
		}
		return nil
	}
}

func testAccCheckSiteMtlsCertificateAssociationBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + "\n" + testAccCheckMtlsImpervaToOriginCertificateBasic(t)
}
