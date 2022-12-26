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

const siemConnectionResourceType = "incapsula_siem_connection"
const s3ArnSiemConnectionResourceName = "test_acc_s3arn"
const s3ArnSiemConnectionResource = siemConnectionResourceType + "." + s3ArnSiemConnectionResourceName

var s3ArnSiemConnectionName = "SIEMCONNECTIONS3ARN" + RandomLetterAndNumberString(10)

var s3ArnSiemConnectionNameUpdated = "UPDATED" + s3ArnSiemConnectionName

const s3SiemConnectionResourceName = "test_acc_s3"
const s3SiemConnectionResource = siemConnectionResourceType + "." + s3SiemConnectionResourceName

var s3SiemConnectionName = "SIEMCONNECTIONS3" + RandomLetterAndNumberString(10)

var s3SiemConnectionNameUpdated = "UPDATED" + s3SiemConnectionName

func TestAccSiemConnection_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection.go.TestAccSiemConnection_Basic")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemConnectionDestroy(siemConnectionResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3ArnSiemConnectionResource),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "connection_name", s3ArnSiemConnectionName),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "storage_type", StorageTypeCustomerS3Arn),
				),
			},
			{
				ResourceName:      s3ArnSiemConnectionResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemConnectionID(siemConnectionResourceType),
			},
		},
	}

	if v := os.Getenv("TESTING_PROFILE"); v == "" {
		r.Steps = append(r.Steps,
			resource.TestStep{
				Config: getAccIncapsulaS3SiemConnectionConfigBasic(s3SiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3SiemConnectionResource),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "connection_name", s3SiemConnectionName),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "storage_type", StorageTypeCustomerS3),
				),
			},
			resource.TestStep{
				ResourceName:            s3SiemConnectionResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
				ImportStateIdFunc:       testACCStateSiemConnectionID(siemConnectionResourceType),
			})
	}
	resource.Test(t, r)
}

func TestAccSiemConnection_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connection.go.TestAccSiemConnection_Update")

	r := resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemConnectionDestroy(siemConnectionResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3ArnSiemConnectionResource),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "connection_name", s3ArnSiemConnectionName),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "storage_type", StorageTypeCustomerS3Arn),
				),
			},
			{
				Config: getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3ArnSiemConnectionResource),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "connection_name", s3ArnSiemConnectionNameUpdated),
					resource.TestCheckResourceAttr(s3ArnSiemConnectionResource, "storage_type", StorageTypeCustomerS3Arn),
				),
			},
		},
	}

	if v := os.Getenv("TESTING_PROFILE"); v == "" {
		r.Steps = append(r.Steps,
			resource.TestStep{
				Config: getAccIncapsulaS3SiemConnectionConfigBasic(s3SiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3SiemConnectionResource),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "connection_name", s3SiemConnectionName),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "storage_type", StorageTypeCustomerS3),
				),
			},
			resource.TestStep{
				Config: getAccIncapsulaS3SiemConnectionConfigBasic(s3SiemConnectionNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemConnectionExists(s3SiemConnectionResource),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "connection_name", s3SiemConnectionNameUpdated),
					resource.TestCheckResourceAttr(s3SiemConnectionResource, "storage_type", StorageTypeCustomerS3),
				),
			})
	}

	resource.Test(t, r)
}

func getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {	
			connection_name = "%s"
  			storage_type = "%s"	
  			path = "%s"
		}`,
		siemConnectionResourceType, s3ArnSiemConnectionResourceName,
		s3ArnSiemConnectionName, StorageTypeCustomerS3Arn,
		os.Getenv("SIEM_CONNECTION_S3_PATH"),
	)
}

func getAccIncapsulaS3SiemConnectionConfigBasic(s3SiemConnectionName string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {	
			connection_name = "%s"
  			storage_type = "%s"
  			access_key = "%s"
  			secret_key = "%s"
  			path = "%s"
		}`,
		siemConnectionResourceType, s3SiemConnectionResourceName, s3SiemConnectionName,
		StorageTypeCustomerS3, os.Getenv("SIEM_CONNECTION_S3_ACCESS_KEY"),
		os.Getenv("SIEM_CONNECTION_S3_SECRET_KEY"),
		os.Getenv("SIEM_CONNECTION_S3_PATH"),
	)
}

func testAccReadSiemConnection(client *Client, ID string) error {
	log.Printf("[INFO] SiemConnection ID: %s", ID)
	siemConnection, statusCode, err := client.ReadSiemConnection(ID, "")
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

func testAccIncapsulaSiemConnectionDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			err := testAccReadSiemConnection(client, res.Primary.ID)
			if err != nil {
				return nil
			} else {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}

		}

		return nil
	}
}

func testCheckIncapsulaSiemConnectionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection resource not found : %s", resource)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemConnection(client, res.Primary.ID)
		}
	}
}

func testACCStateSiemConnectionID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s", rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find SiemConnection ID")
	}
}
