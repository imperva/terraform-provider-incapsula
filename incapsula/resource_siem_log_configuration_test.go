package incapsula

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const siemLogConfigurationResourceType = "incapsula_siem_log_configuration"
const siemLogConfigurationResourceName = "test_acc"
const siemLogConfigurationResource = siemLogConfigurationResourceType + "." + siemLogConfigurationResourceName

var siemLogConfigurationName = "SIEMLOGCONFIGURATION" + RandomLetterAndNumberString(10)

var siemLogConfigurationNameUpdated = "UPDATE" + siemLogConfigurationName

func TestSiemLogConfiguration_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccSiemLogConfiguration_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\"", "\"DNSMS_SECURITY_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_dnsms"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "producer", "ABP"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "producer", "NETSEC"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "producer", "ATO"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "producer", "AUDIT"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationName+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "producer", "CSP"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationName+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "producer", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "format", "CEF"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationName+"attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "producer", "ATTACK_ANALYTICS"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "format", "CEF"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_dnsms", "configuration_name", siemLogConfigurationName+"dnsms"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_dnsms", "producer", "DNSMS"),
				),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_abp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_netsec",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_ato",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_audit",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_csp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_cloudwaf",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_attackanalytics",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_dnsms",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
		},
	})
}

func TestSiemLogConfiguration_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccSiemLogConfiguration_Update")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\"", "\"DNSMS_SECURITY_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_dnsms"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationName+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationName+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationName+"attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_dnsms", "configuration_name", siemLogConfigurationName+"dnsms"),
				),
			},
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationNameUpdated, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\"", "\"DNSMS_SECURITY_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_dnsms"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationNameUpdated+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationNameUpdated+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationNameUpdated+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationNameUpdated+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationNameUpdated+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationNameUpdated+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationNameUpdated+"attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_dnsms", "configuration_name", siemLogConfigurationNameUpdated+"dnsms"),
				),
			},
		},
	})
}

func getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName string, abpDatasets string, netsecDatasets string, atoDatasets string, auditDatasets string, cspDatasets string, cloudWafDatasets string, attackAnalyticsDatasets string, dnsMsDatasets string) string {
	return getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName, "data-platform-access-logs-dev/test/cwaf/51319839") + fmt.Sprintf(`
		resource "%s" "%s" {
			configuration_name = "%s"
  			producer = "ABP"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_abp", siemLogConfigurationName+"abp",
		abpDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "NETSEC"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_netsec", siemLogConfigurationName+"netsec",
		netsecDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "ATO"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_ato", siemLogConfigurationName+"ato",
		atoDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "AUDIT"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_audit", siemLogConfigurationName+"audit",
		auditDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {
			configuration_name = "%s"
			producer = "CSP"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_csp", siemLogConfigurationName+"csp",
		cspDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
	    resource "%s" "%s" {
	       configuration_name = "%s"
	       producer = "CLOUD_WAF"
	       datasets = [%s]
	       enabled = true
	       connection_id = %s.%s.id
           format = "CEF"
        }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_cloudwaf", siemLogConfigurationName+"cloudwaf",
		cloudWafDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
      	 resource "%s" "%s" {
      	    configuration_name = "%s"
      	    producer = "ATTACK_ANALYTICS"
      	    datasets = [%s]
      	    enabled = true
      	    connection_id = %s.%s.id
			format = "CEF"
         }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_attackanalytics", siemLogConfigurationName+"attackanalytics",
		attackAnalyticsDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
        resource "%s" "%s" {
            configuration_name = "%s"
            producer = "DNSMS"
            datasets = [%s]
            enabled = true
            connection_id = %s.%s.id
        }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_dnsms", siemLogConfigurationName+"dnsms",
		dnsMsDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	)
}

func testAccReadSiemLogConfiguration(client *Client, ID string, accountId string) error {
	log.Printf("[INFO] SiemLogConfiguration ID: %s", ID)
	siemLogConfiguration, statusCode, err := client.ReadSiemLogConfiguration(ID, accountId)
	if err != nil {
		return err
	}

	if (*statusCode == 200) && (siemLogConfiguration != nil) && (len(siemLogConfiguration.Data) == 1) && (siemLogConfiguration.Data[0].ID == ID) {
		log.Printf("[INFO] SiemLogConfiguration : %v\n", siemLogConfiguration)
		return nil
	} else if *statusCode == 400 {
		return fmt.Errorf("[ERROR] SiemLogConfiguration with id %s not found", ID)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func testAccIncapsulaSiemLogConfigurationDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			err := testAccReadSiemLogConfiguration(client, res.Primary.ID, res.Primary.Attributes["account_id"])
			if err != nil {
				return nil
			} else {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}

		}

		return nil
	}
}

func testCheckIncapsulaSiemLogConfigurationExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemLogConfiguration resource not found : %s", resource)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemLogConfiguration does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemLogConfiguration(client, res.Primary.ID, res.Primary.Attributes["account_id"])
		}
	}
}

func testACCStateSiemLogConfigurationID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["account_id"], rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find SiemLogConfiguration ID")
	}
}
