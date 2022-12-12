package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"testing"
)

const siemLogConfigurationAbpResourceType = "incapsula_siem_log_configuration_abp"
const siemLogConfigurationAbpResourceName = "test_acc"
const siemLogConfigurationAbpResource = siemLogConfigurationAbpResourceType + "." + siemLogConfigurationAbpResourceName

var siemLogConfigurationAbpName = "SIEMLOGCONFIGURATIONABP" + RandomLetterAndNumberString(10)

func TestAccAbpSiemLogConfiguration_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccAbpSiemLogConfiguration_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationAbpResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaAbpSiemLogConfigurationConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationAbpResource),
					resource.TestCheckResourceAttr(siemLogConfigurationAbpResource, "connection_name", siemLogConfigurationAbpName),
				),
			},
			{
				ResourceName:      siemLogConfigurationAbpResource,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationAbpResourceType),
			},
		},
	})
}

func getAccIncapsulaAbpSiemLogConfigurationConfigBasic() string {
	return getAccIncapsulaS3SiemConnectionConfigBasic() + fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "52291885"
			configuration_name = "%s"
  			producer = "ABP"
			datasets = ["ABP"]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationAbpResourceType, siemLogConfigurationAbpResourceName, siemLogConfigurationAbpName,
		s3SiemConnectionResourceType, s3SiemConnectionResourceName,
	)
}
