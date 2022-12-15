package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"os"
	"testing"
)

//for tests without connectivity check put env variable: TESTING_PROFILE=true

const s3SiemConnectionResourceType = "incapsula_siem_connection_s3"
const s3SiemConnectionResourceName = "test_acc"
const s3SiemConnectionResource = s3SiemConnectionResourceType + "." + s3SiemConnectionResourceName

const siemConnectionS3StorageTypeValue = "CUSTOMER_S3"

var s3SiemConnectionName = "SIEMCONNECTIONS3" + RandomLetterAndNumberString(10)

func TestAccS3SiemConnection_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_s3_siem_connection.go.TestAccS3SiemConnection_Basic")

	if v := os.Getenv("TESTING_PROFILE"); v != "" {
		log.Printf("[DEBUG]TESTING_PROFILE environment variable is provided, test is skipped")
		return
	}

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
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "storage_type", siemConnectionS3StorageTypeValue),
				),
			},
			{
				ResourceName:      s3SiemConnectionResource,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testACCStateSiemConnectionID(s3SiemConnectionResourceType),
			},
		},
	})
}

func getAccIncapsulaS3SiemConnectionConfigBasic() string {
	return fmt.Sprintf(`
		resource "%s" "%s" {	
			connection_name = "%s"
  			storage_type = "%s"
  			access_key = "%s"
  			secret_key = "%s"
  			path = "data-platform-access-logs-dev/testacc"
		}`,
		s3SiemConnectionResourceType, s3SiemConnectionResourceName, s3SiemConnectionName,
		siemConnectionS3StorageTypeValue, getSiemConnectionS3AccessKey(), getSiemConnectionS3SecretKey(),
	)
}

func getSiemConnectionS3AccessKey() string {
	k := os.Getenv("SIEM_CONNECTION_S3_ACCESS_KEY")
	if k == "" {
		return RandomCapitalLetterAndNumberString(20)
	}
	return k
}

func getSiemConnectionS3SecretKey() string {
	k := os.Getenv("SIEM_CONNECTION_S3_SECRET_KEY")
	if k == "" {
		return RandomLetterAndNumberString(40)
	}
	return k
}
