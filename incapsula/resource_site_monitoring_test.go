package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const siteMonitoringResourceName = "incapsula_site_monitoring"
const siteMonitoringResource = siteMonitoringResourceName + "." + siteMonitoringName
const siteMonitoringName = "testacc-terraform-site_monitoring"

func TestAccIncapsulaSiteMonitoring_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_site_monitoring_test.TestAccIncapsulaSiteMonitoring_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSiteMonitoringBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckSiteMonitoringExists(siteMonitoringResource),
					resource.TestCheckResourceAttr(siteMonitoringResource, "monitoring_url", "/users"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "expected_received_string", ""),
					resource.TestCheckResourceAttr(siteMonitoringResource, "up_check_retries", "5"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "up_checks_interval", "1"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "up_checks_interval_units", "MINUTES"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "use_verification_for_down", "false"),

					resource.TestCheckResourceAttr(siteMonitoringResource, "http_request_timeout", "1"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "http_request_timeout_units", "MINUTES"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "http_response_error", "501,503"),

					resource.TestCheckResourceAttr(siteMonitoringResource, "failed_requests_duration", "2"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "failed_requests_duration_units", "MINUTES"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "failed_requests_min_number", "10"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "failed_requests_percentage", "10"),

					resource.TestCheckResourceAttr(siteMonitoringResource, "alarm_on_dc_failover", "false"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "alarm_on_server_failover", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "alarm_on_stands_by_failover", "true"),
					resource.TestCheckResourceAttr(siteMonitoringResource, "required_monitors", "MANY"),
				),
			},
			{
				ResourceName:      siteMonitoringResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiteMonitoringID,
			},
		},
	})
}

func testCheckSiteMonitoringExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Site Monitoring resource not found: %s", name)
		}
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.GetSiteMonitoring(siteId)
		if err != nil {
			fmt.Errorf("Incapsula Site Monitoring doesn't exist")
		}

		return nil
	}
}

func testACCStateSiteMonitoringID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != siteMonitoringResourceName {
			continue
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])

		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int in Site Monitoring resource test", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("Error finding site_id argument in Site Monitoring resource test")
}

func testAccCheckSiteMonitoringBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
  		depends_on = ["%s"]
	
		http_request_timeout       = 1
		http_request_timeout_units = "MINUTES"
		http_response_error        = "501,503"
   
		failed_requests_duration       = 2
		failed_requests_duration_units = "MINUTES"
		failed_requests_min_number     = 10
		failed_requests_percentage     = 10

		alarm_on_dc_failover        = false
		alarm_on_server_failover    = true
		alarm_on_stands_by_failover = true
		required_monitors           = "MANY"

		monitoring_url            = "/users"
		up_check_retries          = 5
		up_checks_interval        = 1
		up_checks_interval_units  = "MINUTES"
		use_verification_for_down = false
	}`,
		siteMonitoringResourceName, siteMonitoringName, siteResourceName,
	)
}
