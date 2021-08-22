package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const dataCentersConfigurationResource = "incapsula_data_centers_configuration"
const dataCentersConfigurationName = "testacc-terraform-data-centers-configuration"
const dataCentersConfigurationResourceName = dataCentersConfigurationResource + "." + dataCentersConfigurationName

func TestAccIncapsulaDataCentersConfiguration_Basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_data_centers_configuration_test.TestAccIncapsulaDataCentersConfiguration_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDataCentersConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataCentersConfigurationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDataCentersConfigurationExists(dataCentersConfigurationResourceName),
					resource.TestCheckResourceAttr(dataCentersConfigurationResourceName, "site_topology", "SINGLE_DC"),
				),
			},
			{
				ResourceName:      dataCentersConfigurationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateDataCentersConfigurationID,
			},
		},
	})
}

func testAccStateDataCentersConfigurationID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != dataCentersConfigurationResource {
			continue
		}

		dccID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", dccID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaDataCentersConfigurationDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteID := res.Primary.ID
		if siteID == "" {
			// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
			// Specifically, upgrades/destroys are happening simulatneously and not honoroing
			// dependencies. In this case, it's possible that the site has already been deleted,
			// which means that all of the subresources will have been cleared out.
			// Ordinarily, this should return an error, but until this gets addressed, we're
			// going to simply return nil.
			// return fmt.Errorf("Incapsula site ID does not exist")
			return nil
		}

		listDataCenterResponse, _ := client.GetDataCentersConfiguration(siteID)

		// See comment above - the data center may have already been deleted
		// This workaround will be removed in the future
		if listDataCenterResponse == nil || listDataCenterResponse.Data == nil || len(listDataCenterResponse.Data) == 0 {
			return nil
		}

		// Nothing to check here.
		// Destroying incapsula_data_centers_configuration should not have any backend effect.

	}

	return nil
}

func testCheckIncapsulaDataCentersConfigurationExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		siteRes, siteOk := state.RootModule().Resources[siteResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula site resource not found: %s", siteResourceName)
		}

		siteIDString := siteRes.Primary.ID
		if siteIDString == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(siteIDString)
		if err != nil {
			return fmt.Errorf("Error parsing Site ID %v to int", siteIDString)
		}

		domain := siteRes.Primary.Attributes["domain"]

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula data centers configuration resource not found: %s", name)
		}

		dccID := res.Primary.ID
		if dccID == "" {
			return fmt.Errorf("Incapsula data centers configuration ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)

		// If the site has already been deleted then return nil
		// Otherwise check the data center list
		_, err = client.SiteStatus(domain, siteID)
		if err != nil {
			return nil
		}

		responseDTO, err := client.GetDataCentersConfiguration(siteIDString)
		if responseDTO == nil || responseDTO.Data == nil || len(responseDTO.Data) == 0 {
			return fmt.Errorf("Incapsula data centers configuration: %s (Site ID: %d) does not exist\n%s", name, siteID, err)
		}

		if siteRes.Primary.ID != res.Primary.ID {
			return fmt.Errorf("The ID of Incapsula data centers configuration: %s is invalid (ID: %s). "+
				"It should be identical to site_id (%s)", name, dccID, siteRes.Primary.ID)
		}

		return nil
	}
}

func testAccCheckIncapsulaDataCentersConfigurationBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
  site_topology = "SINGLE_DC"
  data_center {
    name = "%s"
    origin_server {
      address = "4.4.4.4"
    }
  }
  depends_on = ["%s"]
}`, dataCentersConfigurationResource, dataCentersConfigurationName, siteResourceName, dataCenterName, siteResourceName,
	)
}
