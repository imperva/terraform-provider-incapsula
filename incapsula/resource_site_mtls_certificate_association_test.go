package incapsula

//
//import (
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"log"
//	"strconv"
//
//	"testing"
//)
//
//const mtlsCrtificateResourceName = "incapsula_mtls_imperva_to_origin_certificate"
//const mtlsCrtificateResource = mtlsCrtificateResourceName + "." + mtlsCrtificateName
//const mtlsCrtificateName = "testacc-terraform-mtls-certificate"
//
//func TestAccIncapsulaSiteCertificateAssociation_Basic(t *testing.T) {
//	log.Printf("========================BEGIN TEST========================")
//	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic")
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testACCStateMtlsImpervaToOriginCertificateDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckMtlsImpervaToOriginCertificateBasic(t),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckMtlsImpervaToOriginCertificateExists(),
//					resource.TestCheckResourceAttr(mtlsCrtificateResource, "input_hash", calculatedHash),
//					resource.TestCheckResourceAttr(mtlsCrtificateResource, "certificate_name", "acceptance test certificate"),
//				),
//			},
//			{
//				ResourceName:      mtlsCrtificateResource,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: testACCStateMtlsImpervaToOriginCertificateID,
//			},
//		},
//	})
//}
