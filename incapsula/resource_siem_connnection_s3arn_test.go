package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"testing"
)

const s3ArnSiemConnectionResourceType = "incapsula_s3arn_siem_connection"
const s3ArnSiemConnectionResourceName = "terraform-siem-connection-config-testacc"
const s3ArnSiemConnectionResource = s3ArnSiemConnectionResourceType + "." + s3ArnSiemConnectionResourceName
const s3ArnSiemConnectionName = "SIEM CONNECTION S3 ARN TESTACC"

func TestAccS3ArnSiemConnection_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_s3arn_siem_connection.go.TestAccS3ArnSiemConnection_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemConnectionDestroy(s3ArnSiemConnectionResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaS3SiemConnectionConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3ArnSiemConnectionResource),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "connection_name", s3ArnSiemConnectionName),
				),
			},
			{
				ResourceName:      s3ArnSiemConnectionResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemConnectionID(s3ArnSiemConnectionResourceType),
			},
		},
	})
}

func getAccIncapsulaS3ArnSiemConnectionConfigBasic() string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "52291885"
			connection_name = "%s"
  			storage_type = "CUSTOMER_S3_ARN"
  			path = "data-platform-access-logs-dev/testacc"
		}`,
		s3ArnSiemConnectionResourceType, s3ArnSiemConnectionResourceName, s3ArnSiemConnectionName,
	)
}
