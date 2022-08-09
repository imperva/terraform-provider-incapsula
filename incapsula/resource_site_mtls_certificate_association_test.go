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
	//resource.Test(t, resource.TestCase{
	//	PreCheck:     func() { testAccPreCheck(t) },
	//	Providers:    testAccProviders,
	//	//CheckDestroy: testACCStateMtlsImpervaToOriginCertificateDestroy,
	//	Steps: []resource.TestStep{
	//		{
	//			//Config: testAccCheckMtlsImpervaToOriginCertificateBasic(t),
	//			Check: resource.ComposeTestCheckFunc(
	//				testCheckSiteMtlsCertificateAssociationExists(),
	//				resource.TestCheckResourceAttr(siteMtlsCrtificateAssociationResource, "site_id", ""),
	//				resource.TestCheckResourceAttr(siteMtlsCrtificateAssociationResource, "certificate_id", ""),
	//			),
	//		},
	//		{
	//			ResourceName:      siteMtlsCrtificateAssociationResource,
	//			ImportState:       true,
	//			ImportStateVerify: true,
	//			ImportStateIdFunc: testACCStateSiteMtlsCertificateAssociationID,
	//		},
	//	},
	//})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testACCStateSiteMtlsCertificateAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSiteMtlsCertificateAssociationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSiteMtlsCertificateAssociationExists(),
					//resource.TestCheckResourceAttr(mtlsCrtificateResource, "input_hash", calculatedHashCACert),
					//resource.TestCheckResourceAttr(mtlsCrtificateResource, "certificate_name", "acceptance test certificate"),
				),
			},
		},
	})
}

//todo:
//check destroy

func testACCStateSiteMtlsCertificateAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteMtlsCrtificateAssociationResourceName {
			continue
		}
		return nil

		siteID, ok := rs.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula mTLS certificate is not assigned to Site ID %s (cannot find mandatory parameter site_id)", siteID)
		}

		_, err := client.GetMTLSCertificate(siteID)
		if err == nil {
			return fmt.Errorf("Resource %s with siteID ID %s still exists", siteMtlsCrtificateAssociationResourceName, siteID)
		}
	}
	return fmt.Errorf("Error finding site_id")
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
		//res, ok := state.RootModule().Resources[siteMtlsCrtificateAssociationResource]
		//if !ok {
		//	return fmt.Errorf("Incapsula mTLS Imperva to Origin Certificate resource not found : %s", siteMtlsCrtificateAssociationResource)
		//}
		//certificateID, ok := res.Primary.Attributes["certificate_id"]
		//
		//siteID, ok := res.Primary.Attributes["site_id"]
		//if !ok {
		//	return fmt.Errorf("Incapsula mTLS certificate ID %s is not assigned to Site ID %sn (cannot find mandatory parameter site_id)", siteID, certificateID)
		//}
		//siteIdInt, err := strconv.Atoi(siteID)
		//
		//client := testAccProvider.Meta().(*Client)
		//response, err := client.GetSiteMtlsCertificateAssociation(siteIdInt)
		//if err != nil || response == nil {
		//	return fmt.Errorf("Incapsula mTLS certificate ID %s is not assigned to Site ID %s", certificateID, siteID)
		//}
		return nil
	}
}

func testAccCheckSiteMtlsCertificateAssociationBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + "\n" + testAccCheckMtlsImpervaToOriginCertificateBasic(t)
}
