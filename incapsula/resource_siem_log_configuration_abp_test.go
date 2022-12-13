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

const siemLogConfigurationAbpValue = "ABP"

const siemLogConfigurationAbpAccountIdValue = "52291885"

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
				Config: getAccIncapsulaAbpSiemLogConfigurationConfigBasic(siemLogConfigurationAbpAccountIdValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationAbpResource),
					resource.TestCheckResourceAttr(siemLogConfigurationAbpResource, "configuration_name", siemLogConfigurationAbpName),
					resource.TestCheckResourceAttr(siemLogConfigurationAbpResource, "producer", siemLogConfigurationAbpValue),
					resource.TestCheckResourceAttr(siemLogConfigurationAbpResource, "account_id", siemLogConfigurationAbpAccountIdValue),
				),
			},
			{
				ResourceName:      siemLogConfigurationAbpResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationAbpResourceType),
			},
		},
	})
}

func getAccIncapsulaAbpSiemLogConfigurationConfigBasic(accountId string) string {
	return getAccIncapsulaS3SiemConnectionConfigBasic(accountId) + fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%s"
			configuration_name = "%s"
  			producer = "%s"
			datasets = ["%s"]
			enabled = true
			version = "1.0"
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationAbpResourceType, siemLogConfigurationAbpResourceName, accountId, siemLogConfigurationAbpName,
		siemLogConfigurationAbpValue, siemLogConfigurationAbpValue,
		s3SiemConnectionResourceType, s3SiemConnectionResourceName,
	)
}
