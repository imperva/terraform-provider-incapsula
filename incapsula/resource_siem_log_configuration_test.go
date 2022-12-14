package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
)

func testAccReadSiemLogConfiguration(client *Client, ID string) error {
	log.Printf("[INFO] SiemLogConfiguration ID: %s", ID)
	siemLogConfiguration, statusCode, err := client.ReadSiemLogConfiguration(ID)
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

			err := testAccReadSiemLogConfiguration(client, res.Primary.ID)
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
			return testAccReadSiemLogConfiguration(client, res.Primary.ID)
		}
	}
}

func testACCStateSiemLogConfigurationID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s", rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find SiemLogConfiguration ID")
	}
}
