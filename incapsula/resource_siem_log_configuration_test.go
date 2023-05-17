package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
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
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "producer", "ABP"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "producer", "NETSEC"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "producer", "ATO"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "producer", "AUDIT"),
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
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
				),
			},
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationNameUpdated, "\"ABP\"", "\"IP\", \"ATTACK\"", "\"ATO\"", "\"AUDIT_TRAIL\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationNameUpdated+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationNameUpdated+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationNameUpdated+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationNameUpdated+"audit"),
				),
			},
		},
	})
}

func getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName string, abpDatasets string, netsecDatasets string, atoDatasets string, auditDatasets string) string {
	return getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName) + fmt.Sprintf(`
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
