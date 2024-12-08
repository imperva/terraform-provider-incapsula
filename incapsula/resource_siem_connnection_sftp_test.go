package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"os"
	"testing"
)

const siemSftpConnectionResourceType = "incapsula_siem_sftp_connection"

const sftpSiemConnectionResourceName = "test_acc_sftp"
const sftpSiemConnectionResource = siemSftpConnectionResourceType + "." + sftpSiemConnectionResourceName

var sftpSiemConnectionName = "SIEMCONNECTIONSFTP" + RandomLetterAndNumberString(10)

var sftpSiemConnectionNameUpdated = "UPDATED" + sftpSiemConnectionName

func isSftpEnvVarExist() (bool, error) {
	skipTest := false
	if v := os.Getenv("SIEM_CONNECTION_SFTP_HOST"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SFTP_HOST environment variable does not exist, if you want to test SIEM connection you must provide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SFTP_PATH"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SFTP_PATH environment variable does not exist, if you want to test SIEM connection you must provide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SFTP_USERNAME"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SFTP_USERNAME environment variable does not exist, if you want to test SIEM connection you must provide it")
	}

	if v := os.Getenv("SIEM_CONNECTION_SFTP_PASSWORD"); v == "" {
		skipTest = true
		log.Printf("[ERROR] SIEM_CONNECTION_SFTP_PASSWORD environment variable does not exist, if you want to test SIEM connection you must provide it")
	}

	return skipTest, nil
}

func TestAccSiemSftpConnection(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection_sftp.go.TestAccSiemSftpConnection")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemSftpConnectionDestroy(siemSftpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isSftpEnvVarExist,
				Config:   getAccIncapsulaSftpSiemConnectionConfig(sftpSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSftpConnectionExists(sftpSiemConnectionResource),
					resource.TestCheckResourceAttr(sftpSiemConnectionResource, "connection_name", sftpSiemConnectionName),
				),
			},
			{
				ResourceName:            sftpSiemConnectionResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateIdFunc:       testACCStateSiemSftpConnectionID(siemSftpConnectionResourceType),
			},
		},
	}

	resource.Test(t, r)
}

func TestAccSiemSftpConnection_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection_sftp.go.TestAccSiemSftpConnection_Update")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemSftpConnectionDestroy(siemSftpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isSftpEnvVarExist,
				Config:   getAccIncapsulaSftpSiemConnectionConfig(sftpSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSftpConnectionExists(sftpSiemConnectionResource),
					resource.TestCheckResourceAttr(sftpSiemConnectionResource, "connection_name", sftpSiemConnectionName),
				),
			},
			{
				SkipFunc: isSftpEnvVarExist,
				Config:   getAccIncapsulaSftpSiemConnectionConfig(sftpSiemConnectionNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemSftpConnectionExists(sftpSiemConnectionResource),
					resource.TestCheckResourceAttr(sftpSiemConnectionResource, "connection_name", sftpSiemConnectionNameUpdated),
				),
			},
		},
	}

	resource.Test(t, r)
}

func getAccIncapsulaSftpSiemConnectionConfig(sftpSiemConnectionName string) string {
	resource := fmt.Sprintf(`
		resource "%s" "%s" {
			connection_name = "%s"
  			host = "%s"
  			path = "%s"
  			username = "%s"
			password = "%s"
		}`,
		siemSftpConnectionResourceType,
		sftpSiemConnectionResourceName,
		sftpSiemConnectionName,
		os.Getenv("SIEM_CONNECTION_SFTP_HOST"),
		os.Getenv("SIEM_CONNECTION_SFTP_PATH"),
		os.Getenv("SIEM_CONNECTION_SFTP_USERNAME"),
		os.Getenv("SIEM_CONNECTION_SFTP_PASSWORD"),
	)

	return resource
}

func testAccReadSiemSftpConnection(client *Client, ID string, accountId string) error {
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

func testAccIncapsulaSiemSftpConnectionDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			err := testAccReadSiemSftpConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"])
			if err != nil {
				return nil
			} else {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}

		}

		return nil
	}
}

func testCheckIncapsulaSiemSftpConnectionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection resource not found : %s", resource)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemSftpConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"])
		}
	}
}

func testACCStateSiemSftpConnectionID(resourceType string) resource.ImportStateIdFunc {
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
