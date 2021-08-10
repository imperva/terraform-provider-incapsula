package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const multipleDataCentersConfigurationName = "testacc-terraform-multiple-data-centers-configuration"
const multipleDataCentersConfigurationResourceName = dataCentersConfigurationResource + "." + multipleDataCentersConfigurationName
const dataSourceDataCenterResource = "incapsula_data_center"
const dataSourceDcName1 = "ds_dc_rest_of_the_world"
const dataSourceDcResourceName1 = "data." + dataSourceDataCenterResource + "." + dataSourceDcName1
const dataSourceDcName2 = "ds_dc_geo_europe"
const dataSourceDcResourceName2 = "data." + dataSourceDataCenterResource + "." + dataSourceDcName2
const dataSourceDcName3 = "ds_dc_standby"
const dataSourceDcResourceName3 = "data." + dataSourceDataCenterResource + "." + dataSourceDcName3
const dataSourceDcName4 = "ds_dc_name_americas_and_is_content"
const dataSourceDcResourceName4 = "data." + dataSourceDataCenterResource + "." + dataSourceDcName4

func TestAccIncapsulaDataSourceDataCenter_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceDataCenterConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaDataSourceDataCenterExists(dataSourceDcResourceName1),
					resource.TestCheckResourceAttr(dataSourceDcResourceName1, "origin_pop", "hkg"),
					testCheckIncapsulaDataSourceDataCenterExists(dataSourceDcResourceName2),
					resource.TestCheckResourceAttr(dataSourceDcResourceName2, "origin_pop", "lon"),
					testCheckIncapsulaDataSourceDataCenterExists(dataSourceDcResourceName3),
					resource.TestCheckResourceAttr(dataSourceDcResourceName3, "ip_mode", "SINGLE_IP"),
					testCheckIncapsulaDataSourceDataCenterExists(dataSourceDcResourceName4),
					resource.TestCheckResourceAttr(dataSourceDcResourceName4, "name", "Americas DC"),
				),
			},
			{
				ResourceName:      multipleDataCentersConfigurationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateDataCentersConfigurationID,
			},
		},
	})
}

func testCheckIncapsulaDataSourceDataCenterExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		dcsConfigurationRes, siteOk := state.RootModule().Resources[multipleDataCentersConfigurationResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula data centers configuration resource not found: %s", multipleDataCentersConfigurationResourceName)
		}

		siteID := dcsConfigurationRes.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula data centers configuration ID does not exist")
		}

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula data source data center resource not found: %s", name)
		}

		dcID := res.Primary.ID
		if dcID == "" {
			return fmt.Errorf("Incapsula data source data center ID does not exist")
		}

		return nil
	}
}

func testAccCheckIncapsulaDataSourceDataCenterConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
  site_lb_algorithm = "GEO_PREFERRED"
  site_topology = "MULTIPLE_DC"

  data_center {
    is_rest_of_the_world = true
    name = "Rest of the world DC"
    origin_pop = "hkg"

    origin_server {
      address = "55.66.77.123"
    }

  }

  data_center  {
    ip_mode = "SINGLE_IP"
    is_active = false
    name = "Standby DC"
    web_servers_per_server = 2

    origin_server {
      address = "12.69.14.73"
    }

  }

  data_center  {
    geo_locations = "AFRICA,EUROPE,ASIA"
    name = "EMEA-DC1"
    origin_pop = "lon"

    origin_server {
      address = "54.74.193.120"
    }

  }

  data_center  {
    geo_locations = "US_EAST,US_WEST"
    name = "Americas DC"
    origin_pop = "iad"
	is_content = true

    origin_server {
      address = "54.90.145.67"
    }

  }

  depends_on = [%s]
}

data "incapsula_data_center" "%s" {
  site_id = %s.id
  filter_by_is_rest_of_the_world = true
}

data "incapsula_data_center" "%s" {
  site_id = %s.id
  filter_by_geo_location = "EUROPE"
}

data "incapsula_data_center" "%s" {
  site_id = %s.id
  filter_by_is_standby = true
}

data "incapsula_data_center" "%s" {
  site_id = %s.id
  filter_by_name = "Americas DC"
  filter_by_is_content = true
}`, dataCentersConfigurationResource, multipleDataCentersConfigurationName, siteResourceName, siteResourceName,
		dataSourceDcName1, multipleDataCentersConfigurationResourceName, dataSourceDcName2, multipleDataCentersConfigurationResourceName,
		dataSourceDcName3, multipleDataCentersConfigurationResourceName, dataSourceDcName4, multipleDataCentersConfigurationResourceName)
}
