package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"testing"
)

const s3SiemConnectionResourceType = "incapsula_s3_siem_connection"
const s3SiemConnectionResourceName = "terraform-siem-connection-config-testacc"
const s3SiemConnectionResource = s3SiemConnectionResourceType + "." + s3SiemConnectionResourceName
const s3SiemConnectionName = "SIEM CONNECTION S3 TESTACC"

func TestAccS3SiemConnection_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_s3_siem_connection.go.TestAccS3SiemConnection_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemConnectionDestroy(s3SiemConnectionResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaS3SiemConnectionConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3SiemConnectionResource),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "connection_name", s3SiemConnectionName),
				),
			},
			{
				ResourceName:      s3SiemConnectionResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemConnectionID(s3SiemConnectionResourceType),
			},
		},
	})
}

func getAccIncapsulaS3SiemConnectionConfigBasic() string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "52291885"
			connection_name = "%s"
  			storage_type = "CUSTOMER_S3"
  			access_key = "AKIA3TS2JGVQ3VGHMXVG"
  			secret_key = "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"
  			path = "data-platform-access-logs-dev/testacc"
		}`,
		s3SiemConnectionResourceType, s3SiemConnectionResourceName, s3SiemConnectionName,
	)
}
