package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"os"
	"testing"
)

// for tests without connectivity check put env variable: TESTING_PROFILE=true

const siemSplunkConnectionResourceType = "incapsula_siem_splunk_connection"

const splunkSiemConnectionResourceName = "test_acc_splunk"
const splunkSiemConnectionResource = siemSplunkConnectionResourceType + "." + splunkSiemConnectionResourceName

var splunkSiemConnectionName = "SIEMCONNECTIONSPLUNK" + RandomLetterAndNumberString(10)

var splunkSiemConnectionNameUpdated = "UPDATED" + splunkSiemConnectionName

func isSplunkEnvVarExist() (bool, error) {
	skipTest := false
	if v := os.Getenv("SIEM_CONNECTION_SPLUNK_HOST"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SPLUNK_HOST environment variable does not exist, if you want to test SIEM connection you must prvide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SPLUNK_PORT"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SPLUNK_PORT environment variable does not exist, if you want to test SIEM connection you must prvide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SPLUNK_TOKEN"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SPLUNK_TOKEN environment variable does not exist, if you want to test SIEM connection you must prvide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SPLUNK_DISABLED_CERT_VALIDATION"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SPLUNK_DISABLED_CERT_VALIDATION environment variable does not exist, if you want to test SIEM connection you must prvide it")
	}

	return skipTest, nil
}
func TestAccSiemSplunkConnection(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection_splunk.go.TestAccSiemSplunkConnection")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemSplunkConnectionDestroy(siemSplunkConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isSplunkEnvVarExist,
				Config:   getAccIncapsulaSplunkSiemConnectionConfig(splunkSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSplunkConnectionExists(splunkSiemConnectionResource),
					resource.TestCheckResourceAttr(splunkSiemConnectionResource, "connection_name", splunkSiemConnectionName),
				),
			},
			{
				ResourceName:            splunkSiemConnectionResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
				ImportStateIdFunc:       testACCStateSiemSplunkConnectionID(siemSplunkConnectionResourceType),
			},
		},
	}

	resource.Test(t, r)
}

func TestAccSiemSplunkConnection_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection_splunk.go.TestAccSiemSplunkConnection_Update")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemSplunkConnectionDestroy(siemSplunkConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isSplunkEnvVarExist,
				Config:   getAccIncapsulaSplunkSiemConnectionConfig(splunkSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSplunkConnectionExists(splunkSiemConnectionResource),
					resource.TestCheckResourceAttr(splunkSiemConnectionResource, "connection_name", splunkSiemConnectionName),
				),
			},
			{
				SkipFunc: isSplunkEnvVarExist,
				Config:   getAccIncapsulaSplunkSiemConnectionConfig(splunkSiemConnectionNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSplunkConnectionExists(splunkSiemConnectionResource),
					resource.TestCheckResourceAttr(splunkSiemConnectionResource, "connection_name", splunkSiemConnectionNameUpdated),
				),
			},
		},
	}

	resource.Test(t, r)
}

func getAccIncapsulaSplunkSiemConnectionConfig(splunkSiemConnectionName string) string {
	resource := fmt.Sprintf(`
		resource "%s" "%s" {	
			connection_name = "%s"
  			host = "%s"
  			port = %s
  			token = "%s"
			disable_cert_verification = %s
		}`,
		siemSplunkConnectionResourceType,
		splunkSiemConnectionResourceName,
		splunkSiemConnectionName,
		os.Getenv("SIEM_CONNECTION_SPLUNK_HOST"),
		os.Getenv("SIEM_CONNECTION_SPLUNK_PORT"),
		os.Getenv("SIEM_CONNECTION_SPLUNK_TOKEN"),
		os.Getenv("SIEM_CONNECTION_SPLUNK_DISABLED_CERT_VALIDATION"),
	)

	return resource
}

func testAccReadSiemSplunkConnection(client *Client, ID string, accountId string) error {
	log.Printf("[INFO] SiemConnection ID: %s, accountId: %s", ID, accountId)
	siemConnection, statusCode, err := client.ReadSiemConnection(ID, accountId)
	if err != nil {
		return err
	}

	if (*statusCode == 200) && (siemConnection != nil) && (len(siemConnection.Data) == 1) && (siemConnection.Data[0].ID == ID) {
		log.Printf("[INFO] SiemConnection : %v\n", siemConnection)
		return nil
	} else if *statusCode == 400 {
		return fmt.Errorf("[ERROR] SiemConnection with id %s not found", ID)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func testAccIncapsulaSiemSplunkConnectionDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			err := testAccReadSiemSplunkConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"])
			if err != nil {
				return nil
			} else {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}

		}

		return nil
	}
}

func testCheckIncapsulaSiemSplunkConnectionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection resource not found : %s", resource)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemSplunkConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"])
		}
	}
}

func testACCStateSiemSplunkConnectionID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["account_id"], rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find SiemConnection ID")
	}
}
