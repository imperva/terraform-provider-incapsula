package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
)

const mtlsClientToImervaCertificateResourceName = "incapsula_mtls_client_to_imperva_ca_certificate"
const mtlsClientToImervaCertificateResource = mtlsClientToImervaCertificateResourceName + "." + mtlsClientToImervaCertificateName
const mtlsClientToImervaCertificateName = "testacc-terraform-client_to_imperva_certificate"

var calculatedHashForClientCACert = ""

// todo KATRIn - do I want to add import test??   do I have import? Probably not
func TestAccIncapsulaMtlsClientToImervaCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateMtlsClientToImervaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMtlsClientToImervaCertificateBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckMtlsClientToImervaCertificateExists(),
					resource.TestCheckResourceAttr(mtlsClientToImervaCertificateResource, "certificate_name", "acceptance test CA certificate"),
				),
			},
		},
	})
}

func testACCStateMtlsClientToImervaCertificateDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	res, ok := state.RootModule().Resources[mtlsClientToImervaCertificateResource]
	if !ok {
		return fmt.Errorf("Incapsula mTLS Client CA To Imperva Certificate resource not found : %s", apiSecApiConfigResource)
	}
	certificateID := res.Primary.ID
	if certificateID == "" {
		fmt.Errorf("Parameter id was not found in resource %s", mtlsClientToImervaCertificateResourceName)
	}

	accountID, ok := res.Primary.Attributes["account_id"]
	if !ok {
		return fmt.Errorf("Mandatory parameter account_id doesn't exist for Incapsula mTLS Client CA To Imperva Certificate with ID %s", certificateID)
	}

	_, exists, _ := client.GetClientCaCertificate(accountID, certificateID)
	if exists == true {
		return fmt.Errorf("Resource %s with cerificate ID %s still exists", mtlsClientToImervaCertificateResourceName, certificateID)
	}
	return nil
}

func testCheckMtlsClientToImervaCertificateExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[mtlsClientToImervaCertificateResource]
		if !ok {
			return fmt.Errorf("Incapsula mTLS Client CA To Imperva Certificate resource not found : %s", apiSecApiConfigResource)
		}
		certificateID := res.Primary.ID
		accountID, ok := res.Primary.Attributes["account_id"]
		if !ok {
			return fmt.Errorf("Incapsula mTLS Client CA To Imperva Certificate withID %s does not exist for Account ID %s", certificateID, accountID)
		}

		client := testAccProvider.Meta().(*Client)
		_, _, err := client.GetClientCaCertificate(accountID, certificateID)
		if err != nil {
			return fmt.Errorf("Incapsula mTLS Client CA To Imperva Certificate with ID %s does not exist for Account ID %s", certificateID, accountID)
		}
		return nil
	}
}

//
//func testACCStateMtlsImpervaToOriginCertificateID(s *terraform.State) (string, error) {
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != mtlsCrtificateResourceName {
//			continue
//		}
//		_, err := strconv.Atoi(rs.Primary.Attributes["id"])
//		if err != nil {
//			return "", fmt.Errorf("Error parsing mTLS Imperva to Origin Certificate ID %v to int", rs.Primary.Attributes["id"])
//		}
//		return rs.Primary.Attributes["id"], nil
//	}
//	return "", fmt.Errorf("Error finding an mTLS Imperva to Origin Certificate\"")
//}

func testAccCheckMtlsClientToImervaCertificateBasic(t *testing.T) string {
	cert, _ := generateKeyPairBase64("dash.beer.center")
	calculatedHashForClientCACert = calculateHash(cert+"\n", "", "")
	res := fmt.Sprintf(`
data "incapsula_account_data" "account_data" {
}

	resource"%s""%s"{
	    account_id       = data.incapsula_account_data.account_data.current_account
		certificate_name = "acceptance test CA certificate"
  		certificate      = "%s"
	}`,
		mtlsClientToImervaCertificateResourceName, mtlsClientToImervaCertificateName, cert,
	)
	return res
}
