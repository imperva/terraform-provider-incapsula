package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/resource"
)

const dataCenterServerAddress = "4.4.4.4"
const dataCenterServerResourceName = "incapsula_data_center_server.testacc-terraform-data-center-server"

func TestAccIncapsulaDataCenterServer_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDataCenterServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataCenterServerConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDataCenterServerExists(dataCenterServerResourceName),
					resource.TestCheckResourceAttr(dataCenterServerResourceName, "server_address", dataCenterServerAddress),
				),
			},
			{
				ResourceName:      dataCenterServerResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateDCID,
			},
		},
	})
}

func testAccStateDCID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_data_center_server" {
			continue
		}

		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing data center ID %v to int", rs.Primary.ID)
		}
		dcID, err := strconv.Atoi(rs.Primary.Attributes["dc_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing dc_id %v to int", rs.Primary.Attributes["dc_id"])
		}
		return fmt.Sprintf("%d/%d", dcID, serverID), nil
	}

	return "", fmt.Errorf("Error finding dc_id")
}

func testAccCheckIncapsulaDataCenterServerDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteID := res.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		listDataCenterResponse, err := client.ListDataCenters(siteID)
		for _, dc := range listDataCenterResponse.DCs {
			if dc.Name == dataCenterName {
				return fmt.Errorf("Incapsula data center: %s (site_id: %s) still exists", dataCenterName, siteID)
			}
		}
		if err == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %s) still exists", testAccDomain, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaDataCenterServerExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		siteRes, siteOk := state.RootModule().Resources[siteResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula site resource not found: %s", siteResourceName)
		}

		siteID := siteRes.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula data center resource not found: %s", name)
		}

		serverID := res.Primary.ID
		if serverID == "" {
			return fmt.Errorf("Incapsula data center ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)
		dataCenterListResponse, err := client.ListDataCenters(siteID)
		if dataCenterListResponse == nil {
			return fmt.Errorf("Incapsula data center: %s (site id: %s) does not exist\n%s", name, siteID, err)
		}

		for _, dc := range dataCenterListResponse.DCs {
			if dc.Name == dataCenterName {
				for _, server := range dc.Servers {
					if server.ID == serverID {
						return nil
					}
				}

			}
		}

		return fmt.Errorf("Incapsula data center: %s serverID: %s is not found in data center servers", name, serverID)
	}
}

func testAccCheckIncapsulaDataCenterServerConfig_basic() string {
	return testAccCheckIncapsulaDataCenterConfig_basic() + fmt.Sprintf(`
resource "incapsula_data_center_server" "testacc-terraform-data-center-server" {
  dc_id = "${incapsula_data_center.testacc-terraform-data-center.id}"
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  server_address = "4.4.4.4"
  is_standby = "yes"
  depends_on = ["%s", "%s"]
}`, siteResourceName, dataCenterResourceName,
	)
}
