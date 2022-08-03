package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"

	"testing"
)

const mtlsCrtificateResourceName = "incapsula_mtls_imperva_to_origin_certificate"
const mtlsCrtificateResource = mtlsCrtificateResourceName + "." + mtlsCrtificateName
const mtlsCrtificateName = "testacc-terraform-mtls-certificate"

func TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateMtlsImpervaToOriginCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMtlsImpervaToOriginCertificateBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckMtlsImpervaToOriginCertificateExists(),
					resource.TestCheckResourceAttr(mtlsCrtificateResource, "input_hash", calculatedHash),
					resource.TestCheckResourceAttr(mtlsCrtificateResource, "certificate_name", "acceptance test certificate"),
				),
			},
			{
				ResourceName:      mtlsCrtificateResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateMtlsImpervaToOriginCertificateID,
			},
		},
	})
}

func testACCStateMtlsImpervaToOriginCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != mtlsCrtificateResourceName {
			continue
		}
		return nil

		certificateID := rs.Primary.ID
		if certificateID == "" {
			fmt.Errorf("Parameter id was not found in resource %s", mtlsCrtificateResourceName)
		}

		_, err := client.GetTLSCertificate(certificateID)
		if err == nil {
			return fmt.Errorf("Resource %s with cerificate ID %s still exists", mtlsCrtificateResourceName, certificateID)
		}
	}
	return fmt.Errorf("Error finding site_id")
}

func testCheckMtlsImpervaToOriginCertificateExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[apiSecApiConfigResource]
		if !ok {
			return fmt.Errorf("Incapsula mTLS Imperva to Origin Certificate resource not found : %s", apiSecApiConfigResource)
		}
		certificateID := res.Primary.ID
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetTLSCertificate(certificateID)
		if err != nil {
			return fmt.Errorf("Incapsula mTLS Imperva to Origin Certificate with ID %s does not exist", certificateID)
		}
		return nil
	}
}

func testACCStateMtlsImpervaToOriginCertificateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != mtlsCrtificateResourceName {
			continue
		}
		_, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing mTLS Imperva to Origin Certificate ID %v to int", rs.Primary.Attributes["id"])
		}
		return rs.Primary.Attributes["id"], nil
	}
	return "", fmt.Errorf("Error finding an mTLS Imperva to Origin Certificate\"")
}

func testAccCheckMtlsImpervaToOriginCertificateBasic(t *testing.T) string {
	cert, privateKey := generateKeyPair()
	log.Printf("cert:\n%s\n\nkey:\n%s", cert, privateKey)
	return fmt.Sprintf(`
	resource"%s""%s"{
		certificate_name = "acceptance test certificate"
  		certificate = %s
  		private_key = %s
	}`,
		mtlsCrtificateResourceName, mtlsCrtificateName, cert, privateKey,
	)
}
