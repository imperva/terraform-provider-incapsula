package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const dataCenterName = "Example data center"
const dataCenterResourceName = "incapsula_data_center.testacc-terraform-data-center"

func testAccIncapsulaDataCenter_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDataCenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataCenterConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDataCenterExists(dataCenterResourceName),
					resource.TestCheckResourceAttr(dataCenterResourceName, "name", dataCenterName),
				),
			},
			{
				ResourceName:      dataCenterResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateDcID,
			},
		},
	})
}

func testAccStateDcID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_data_center" {
			continue
		}

		dcID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, dcID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaDataCenterDestroy(state *terraform.State) error {
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

func testCheckIncapsulaDataCenterExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error parsing Rule ID %v to int", siteIDString)
		}

		domain := siteRes.Primary.Attributes["domain"]

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula data center resource not found: %s", name)
		}

		dcID := res.Primary.ID
		if dcID == "" {
			return fmt.Errorf("Incapsula data center ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)

		// If the site has already been deleted then return nil
		// Otherwise check the data center list
		_, err = client.SiteStatus(domain, siteID)
		if err != nil {
			return nil
		}

		dataCenterListResponse, err := client.ListDataCenters(siteIDString)
		if dataCenterListResponse == nil {
			return fmt.Errorf("Incapsula data center: %s (Site ID: %d) does not exist\n%s", name, siteID, err)
		}

		for _, dc := range dataCenterListResponse.DCs {
			if dc.Name == dataCenterName && dc.ID != dcID {
				return fmt.Errorf("Incapsula data center: %s (DC ID: %s) has invalid ID", name, dcID)
			}
		}

		return nil
	}
}

func testAccCheckIncapsulaDataCenterConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_data_center" "testacc-terraform-data-center" {
  site_id = "${incapsula_site.testacc-terraform-site.id}"
  name = "%s"
  server_address = "8.8.8.5"
  is_content = "true"
  depends_on = ["%s"]
}`, dataCenterName, siteResourceName,
	)
}
