package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"testing"
)

const s3ArnSiemConnectionResourceType = "incapsula_siem_connection_s3arn"
const s3ArnSiemConnectionResourceName = "test_acc"
const s3ArnSiemConnectionResource = s3ArnSiemConnectionResourceType + "." + s3ArnSiemConnectionResourceName

const siemConnectionS3ArnStorageTypeValue = "CUSTOMER_S3_ARN"

var s3ArnSiemConnectionName = "SIEMCONNECTIONS3ARN" + RandomLetterAndNumberString(10)

func TestAccS3ArnSiemConnection_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_s3arn_siem_connection.go.TestAccS3ArnSiemConnection_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemConnectionDestroy(s3ArnSiemConnectionResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaS3ArnSiemConnectionConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3ArnSiemConnectionResource),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "connection_name", s3ArnSiemConnectionName),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "storage_type", siemConnectionS3ArnStorageTypeValue),
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
			connection_name = "%s"
  			storage_type = "%s"	
  			path = "data-platform-access-logs-dev/testacc"
		}`,
		s3ArnSiemConnectionResourceType, s3ArnSiemConnectionResourceName,
		s3ArnSiemConnectionName, siemConnectionS3ArnStorageTypeValue,
	)
}
