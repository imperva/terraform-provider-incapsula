package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
)

func testAccReadSiemConnection(client *Client, ID string) error {
	log.Printf("[INFO] SiemConnection ID: %s", ID)
	siemConnection, statusCode, err := client.ReadSiemConnection(ID)
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

		siemConnectionID := res.Primary.ID
		if siemConnectionID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemConnection(client, siemConnectionID)
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
