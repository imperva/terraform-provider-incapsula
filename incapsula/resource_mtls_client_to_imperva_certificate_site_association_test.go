package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const siteClientToImervaCertificateAssociationResourceName = "incapsula_mtls_site_client_to_imperva_certificate_association"
const siteClientToImervaCertificateAssociationResource = siteMtlsCrtificateAssociationResourceName + "." + siteMtlsCrtificateAssociationName
const siteClientToImervaCertificateAssociationName = "testacc-terraform-site-client-to-imperva-certificate-certificate-association"

//todo KATRIN:
//do we want import?

func TestAccIncapsulaSiteClientToImervaCertificateAssociation_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaSiteMtlsCertificateAssociation_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateSiteClientToImervaCertificateAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckClientToImervaCertificateAssociationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientToImervaCertificateAssociationExists(),
				),
			},
		},
	})
}

//todo KATRIN:
//check destroy!!!!!!
func testACCStateSiteClientToImervaCertificateAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range s.RootModule().Resources {
		if res.Type != siteClientToImervaCertificateAssociationResourceName {
			continue
		}
		return nil

		siteID, certificateID := getResourceDetails(res)

		_, _, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(siteID, certificateID)
		if err == nil {
			return fmt.Errorf("Resource %s with siteID ID %d still exists", siteMtlsCrtificateAssociationResourceName, siteID)
		}
	}
	return fmt.Errorf("Error finding site_id")
}

//	check import
//todo KATRIN:
//func testACCStateClientToImervaCertificateAssociationID(s *terraform.State) (string, error) {
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != siteMtlsCrtificateAssociationResourceName {
//			continue
//		}
//		certificateID, ok := rs.Primary.Attributes["certificate_id"]
//		if !ok {
//			return "", fmt.Errorf("Incapsula mTLS certificate ID %s to Site association doesn't exist (cannot find mandatory parameter certificate_id)", certificateID)
//		}
//		siteID, ok := rs.Primary.Attributes["site_id"]
//		if !ok {
//			return "", fmt.Errorf("Incapsula mTLS certificate ID %s is not assigned to Site ID %sn (cannot find mandatory parameter site_id)", siteID, certificateID)
//		}
//	}
//	return "", fmt.Errorf("Error finding an mTLS Imperva to Origin Certificate\"")
//}

func testCheckClientToImervaCertificateAssociationExists() resource.TestCheckFunc {

	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		log.Print("\nRESOURCES:\n")
		for _, resource := range state.RootModule().Resources {
			log.Printf("\n%v", resource)

			if resource.Type != siteClientToImervaCertificateAssociationResourceName {
				continue
			}
			siteID, certificateID := getResourceDetails(resource)

			response, _, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(siteID, certificateID)
			if err != nil || response == nil {
				return fmt.Errorf("Incapsula mTLS certificate ID %d is not assigned to Site ID %d", certificateID, siteID)
			}
			return nil

		}
		return fmt.Errorf("Error finding %s resource", siteClientToImervaCertificateAssociationResourceName)
	}
}

func testAccCheckClientToImervaCertificateAssociationBasic(t *testing.T) string {
	return testAccCheckIncapsulaCustomCertificateGoodConfig(t) + "\n" + testAccCheckMtlsClientToImervaCertificateBasic(t) + "\n" +
		fmt.Sprintf(`resource "%s" "%s" {
   account_id     = data.incapsula_account_data.account_data.current_account 
   certificate_id = %s.id
   site_id        = incapsula_site.testacc-terraform-site.id
   depends_on     = [%s,%s]
}`, siteClientToImervaCertificateAssociationResourceName,
			siteClientToImervaCertificateAssociationName,
			mtlsClientToImervaCertificateResource,
			mtlsClientToImervaCertificateResource,
			certificateResource)
}

func getResourceDetails(resourceState *terraform.ResourceState) (int, int) {
	siteID, err := strconv.Atoi(resourceState.Primary.Attributes["site_id"])
	if err != nil {
		fmt.Errorf("Error parsing site ID %s to int for %s resource destroy test", resourceState.Primary.Attributes["site_id"], siteClientToImervaCertificateAssociationResourceName)
	}

	certificateID, err := strconv.Atoi(resourceState.Primary.Attributes["certificate_id"])
	if err != nil {
		fmt.Errorf("Error parsing certificate ID %v to int for %s resource destroy test", resourceState.Primary.Attributes["certificate_id"], siteClientToImervaCertificateAssociationResourceName)
	}

	return siteID, certificateID

}
