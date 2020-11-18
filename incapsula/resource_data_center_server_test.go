package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				Config: testAccCheckIncapsulaDataCenterServerConfigBasic(),
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

		serverID := rs.Primary.ID
		siteID := rs.Primary.Attributes["site_id"]
		dcID := rs.Primary.Attributes["dc_id"]

		return fmt.Sprintf("%s/%s/%s", siteID, dcID, serverID), nil
	}

	return "", fmt.Errorf("Error finding a data center server")
}

func testAccCheckIncapsulaDataCenterServerDestroy(state *terraform.State) error {
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

		listDataCenterResponse, _ := client.ListDataCenters(siteID)

		// See comment above - the data center may have already been deleted
		// This workaround will be removed in the future
		if listDataCenterResponse == nil || listDataCenterResponse.DCs == nil || len(listDataCenterResponse.DCs) == 0 {
			return nil
		}

		for _, dc := range listDataCenterResponse.DCs {
			if dc.Name == dataCenterName {
				return fmt.Errorf("Incapsula data center: %s (site_id: %s) still exists", dataCenterName, siteID)
			}
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

func testAccCheckIncapsulaDataCenterServerConfigBasic() string {
	return testAccCheckIncapsulaDataCenterConfigBasic() + fmt.Sprintf(`
resource "incapsula_data_center_server" "testacc-terraform-data-center-server" {
  dc_id = incapsula_data_center.testacc-terraform-data-center.id
  site_id = incapsula_site.testacc-terraform-site.id
  server_address = "4.4.4.4"
	is_enabled = "true"
  depends_on = [%s, %s]
}`, dataCenterResourceName, siteResourceName,
	)
}
