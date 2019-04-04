package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccSiteID = "42"

func TestAccIncapsulaDataCenter_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaDataCenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfig_basic(testAccDomain),
			},
			{
				Config: testAccCheckIncapsulaDataCenterConfig_basic(testAccSiteID),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDataCenterExists("incapsula_data_center.testacc-terraform-data-center"),
					resource.TestCheckResourceAttr("incapsula_data_center.testacc-terraform-data-center", "site_id", testAccSiteID),
				),
			},
		},
	})
}

func testAccCheckIncapsulaDataCenterDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteIDStr := res.Primary.ID
		if siteIDStr == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			return fmt.Errorf("Site ID conversion error for %s: %s", siteIDStr, err)
		}

		_, err = client.SiteStatus(testAccDomain, siteID)

		if err == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) still exists", testAccDomain, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaDataCenterExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula data center resource not found: %s", name)
		}

		siteIDStr := res.Primary.ID
		if siteIDStr == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			return fmt.Errorf("Site ID conversion error for %s: %s", siteIDStr, err)
		}

		client := testAccProvider.Meta().(*Client)
		siteStatusResponse, err := client.SiteStatus(testAccDomain, siteID)
		if siteStatusResponse == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) does not exist", testAccDomain, siteID)
		}

		return nil
	}
}

func testAccCheckIncapsulaDataCenterConfig_basic(siteID string) string {
	return fmt.Sprintf(`resource "incapsula_data_center" "testacc-terraform-data-center" {
  site_id = "%s"
  name = "Example data center"
  server_address = "8.8.4.4"
  is_standby = "yes"
  is_content = "yes"
  depends_on = ["incapsula_site.example-site"]
}`, siteID,
	)
}
