package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const applicationPerformanceResourceName = "incapsula_application_performance"
const applicationPerformanceResource = applicationPerformanceResourceName + "." + applicationPerformanceName
const applicationPerformanceName = "testacc-terraform-application_performance"

func TestAccIncapsulaApplicationPerformance_basic(t *testing.T) {
	var domainName = GenerateTestDomain(t)
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_application_performance_test.TestAccIncapsulaApplicationPerformance_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAllCachingEnabled(domainName),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationPerformanceExists(applicationPerformanceResource),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_comply_no_cache", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_enable_client_side_caching", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_send_age_header", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "key_comply_vary", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "key_unite_naked_full_cache", "true"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_level", "all_resources"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_https", "disabled"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_time", "120"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_300x", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_404_enabled", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_404_time", "120"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_empty_responses", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_http_10_responses", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_response_header_mode", "custom"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_response_headers.#", "1"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_response_headers.0", "cache2"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_stale_content_mode", "custom"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_stale_content_time", "120"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_tag_response_header", "myHeader3"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "ttl_prefer_last_modified", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "ttl_use_shortest_caching", "true"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_shield", "false"),
				),
			},
			{
				Config: testAllCachingDisabled(domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_comply_no_cache", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_enable_client_side_caching", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "client_send_age_header", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "key_comply_vary", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "key_unite_naked_full_cache", "false"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_level", "disabled"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_https", "disabled"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "mode_time", "0"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_300x", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_404_enabled", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_404_time", "0"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_empty_responses", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_http_10_responses", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_response_header_mode", "disabled"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_response_headers.#", "0"),

					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_stale_content_mode", "disabled"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_stale_content_time", "0"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_tag_response_header", ""),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "ttl_prefer_last_modified", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "ttl_use_shortest_caching", "false"),
					resource.TestCheckResourceAttr(applicationPerformanceResource, "response_cache_shield", "false"),
				),
			},
			{
				ResourceName:      applicationPerformanceResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateApplicationPerformanceID,
			},
		},
	})
}

func testCheckApplicationPerformanceExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Application Performance resource not found: %s", name)
		}

		client := testAccProvider.Meta().(*Client)
		_, err2 := client.GetPerformanceSettings(siteId)
		if err2 != nil {
			fmt.Errorf("Incapsula Application Performance doesn't exist")
		}

		return nil
	}
}

func testAccStateApplicationPerformanceID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != applicationPerformanceResourceName {
			continue
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])

		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int in Application Performance resource test", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("Error finding site_id argument in Application Performance resource test")
}

func testAllCachingEnabled(domainName string) string {
	return testCheckIncapsulaSiteV3ConfigBasic(domainName, "CLOUD_WAF") + fmt.Sprintf(`
resource "%s" "%s" {
	site_id = incapsula_site_v3.test-terraform-site-v3.id
	depends_on = ["incapsula_site_v3.test-terraform-site-v3"]
	client_comply_no_cache = true
	client_enable_client_side_caching = true
	client_send_age_header = true
	key_comply_vary = true
	key_unite_naked_full_cache = true
	mode_level = "all_resources"
	mode_https = "disabled"
	mode_time = 120
	response_cache_300x = true
	response_cache_404_enabled = true
	response_cache_404_time = 120
	response_cache_empty_responses = true
	response_cache_http_10_responses = true
	response_cache_response_header_mode = "custom"
	response_cache_response_headers = ["cache2"]
	response_stale_content_mode = "custom"
	response_tag_response_header = "myHeader3"
	response_stale_content_time = 120
	ttl_prefer_last_modified = true
	ttl_use_shortest_caching = true
}`,
		applicationPerformanceResourceName, applicationPerformanceName,
	)
}

func testAllCachingDisabled(domainName string) string {
	return testCheckIncapsulaSiteV3ConfigBasic(domainName, "CLOUD_WAF") + fmt.Sprintf(`
resource "%s" "%s" {
	site_id = incapsula_site_v3.test-terraform-site-v3.id
	depends_on = ["incapsula_site_v3.test-terraform-site-v3"]
	client_comply_no_cache = false
	client_enable_client_side_caching = false
	client_send_age_header = false
	key_comply_vary = false
	key_unite_naked_full_cache = false
	mode_level = "disabled"
	mode_https = "disabled"
	mode_time = 0
	response_cache_300x = false
	response_cache_404_enabled = false
	response_cache_404_time = 0
	response_cache_empty_responses = false
	response_cache_http_10_responses = false
	response_cache_response_header_mode = "disabled"
	response_stale_content_mode = "disabled"
	response_stale_content_time = 0
	ttl_prefer_last_modified = false
	ttl_use_shortest_caching = false
}`,
		applicationPerformanceResourceName, applicationPerformanceName,
	)
}
