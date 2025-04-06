package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

const siteLogConfigResourceType = "incapsula_site_log_configuration"
const siteLogConfigResourceName = "example_site_log_configuration"

const logLevel = "full"
const dataStorageRegion = "US"
const hashingEnabled = true
const hashSalt = "EJKHRT48375N4TKE7956NG"

func TestAccIncapsulaSiteLogConfiguration_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiteLogConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "log_level", logLevel),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "data_storage_region", dataStorageRegion),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hashing_enabled", strconv.FormatBool(hashingEnabled)),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hash_salt", hashSalt),
				),
			},
		},
	})
}

func TestAccIncapsulaSiteLogConfiguration_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiteLogConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "log_level", logLevel),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hashing_enabled", strconv.FormatBool(hashingEnabled)),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hash_salt", hashSalt),
				),
			},
			{
				Config: getAccIncapsulaSiteLogConfigUpdated(GenerateTestSiteName(nil), "CLOUD_WAF"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "log_level", "security"),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "data_storage_region", "EU"),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hashing_enabled", "false"),
					resource.TestCheckResourceAttr(siteLogConfigResourceType+"."+siteLogConfigResourceName, "hash_salt", "EJKHRT48375N777777TE"),
				),
			},
		},
	})
}

func getAccIncapsulaSiteLogConfigBasic(name string, siteType string) string {
	result := getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName, "data-platform-access-logs-dev/test/cwaf/51319839") + fmt.Sprintf(`
	    resource "%s" "%s" {
	       configuration_name = "%s"
	       producer = "CLOUD_WAF"
	       datasets = [%s]
	       enabled = true
	       connection_id = %s.%s.id
           format = "CEF"
        }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_cloudwaf", siemLogConfigurationName+"cloudwaf",
		"\"CLOUD_WAF_ACCESS\"", siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "incapsula_site_v3" "test_log_config_site" {
			name = "%s"
		    type = "%s"
		}
		resource "%s" "%s" {
			site_id = incapsula_site_v3.test_log_config_site.id
			log_level = "%s"
			data_storage_region = "%s"
			hashing_enabled = %t
			hash_salt = "%s"
		}`,
		name, siteType, siteLogConfigResourceType, siteLogConfigResourceName, logLevel, dataStorageRegion, hashingEnabled, hashSalt,
	)
	return result
}

func getAccIncapsulaSiteLogConfigUpdated(name string, siteType string) string {
	result := getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName, "data-platform-access-logs-dev/test/cwaf/51319839") + fmt.Sprintf(`
	    resource "%s" "%s" {
	       configuration_name = "%s"
	       producer = "CLOUD_WAF"
	       datasets = [%s]
	       enabled = true
	       connection_id = %s.%s.id
           format = "CEF"
        }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_cloudwaf", siemLogConfigurationName+"cloudwaf",
		"\"CLOUD_WAF_ACCESS\"", siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "incapsula_site_v3" "test_log_config_site" {
			name = "%s"
		    type = "%s"
		}
		resource "%s" "%s" {
			site_id = incapsula_site_v3.test_log_config_site.id
			log_level = "security"
			data_storage_region = "EU"
			hashing_enabled = false
			hash_salt = "EJKHRT48375N777777TE"
		}`,
		name, siteType, siteLogConfigResourceType, siteLogConfigResourceName,
	)
	return result
}
