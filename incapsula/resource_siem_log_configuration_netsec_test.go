package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"testing"
)

const siemLogConfigurationNetsecResourceType = "incapsula_siem_log_configuration_netsec"
const siemLogConfigurationNetsecResourceName = "test_acc"
const siemLogConfigurationNetsecResource = siemLogConfigurationNetsecResourceType + "." + siemLogConfigurationNetsecResourceName

const siemLogConfigurationNetsecValue = "NETSEC"

const siemLogConfigurationNetsecAccountIdValue = "52291885"

var siemLogConfigurationNetsecName = "SIEMLOGCONFIGURATIONNETSEC" + RandomLetterAndNumberString(10)

func TestAccNetsecSiemLogConfiguration_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccNetsecSiemLogConfiguration_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationNetsecResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaNetsecSiemLogConfigurationConfigBasic(siemLogConfigurationNetsecAccountIdValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationNetsecResource),
					resource.TestCheckResourceAttr(siemLogConfigurationNetsecResource, "configuration_name", siemLogConfigurationNetsecName),
					resource.TestCheckResourceAttr(siemLogConfigurationNetsecResource, "producer", siemLogConfigurationNetsecValue),
					resource.TestCheckResourceAttr(siemLogConfigurationNetsecResource, "account_id", siemLogConfigurationNetsecAccountIdValue),
				),
			},
			{
				ResourceName:      siemLogConfigurationNetsecResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationNetsecResourceType),
			},
		},
	})
}

func getAccIncapsulaNetsecSiemLogConfigurationConfigBasic(accountId string) string {
	return getAccIncapsulaS3SiemConnectionConfigBasic(accountId) + fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%s"
			configuration_name = "%s"
  			producer = "%s"
			datasets = ["CONNECTION", "NETFLOW", "IP", "ATTACK"]
			enabled = true
			version = "1.0"
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationNetsecResourceType, siemLogConfigurationNetsecResourceName, accountId, siemLogConfigurationNetsecName,
		siemLogConfigurationNetsecValue,
		s3SiemConnectionResourceType, s3SiemConnectionResourceName,
	)
}
