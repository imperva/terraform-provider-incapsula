package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const applicationDeliveryResourceName = "incapsula_application_delivery"
const applicationDeliveryResource = applicationDeliveryResourceName + "." + applicationDeliveryName
const applicationDeliveryName = "testacc-terraform-application_delivery"
const customErrorPageBasic = "<!DOCTYPE html>\n<html lang=‘en’>\n  <head>\n    <meta charset=‘UTF-8’>\n    <meta name=‘viewport’ content=‘width=device-width, initial-scale=1.0’>\n    <meta http-equiv=‘X-UA-Compatible’ content=‘ie=edge’>\n    <link href=‘https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;700&display=swap’ rel=‘stylesheet’>\n    <title>[Error Title]</title>\n      </head>\n  <body>\n    <div class=‘container’>\n      <div class=‘container-inner’>\n        <div class=‘header’>\n          <div class=‘error-description’>\n            $TITLE$\n          </div>\n        </div>\n        <div class=‘main’>\n          <div class=‘troubleshooting’>\n            <div class=‘content’>\n              $BODY$\n            </div>\n\t    <h1>custom edited error</h1>\n          </div>\n        </div>\n      </div>\n    </div>\n  </body>\n</html>"
const customErrorPageInput = "<<-EOT\n" + customErrorPageBasic + "\nEOT"

func TestAccIncapsulaApplicationDelivery_basic(t *testing.T) {
	var domainName = GenerateTestDomain(t)
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_application_delivery_test.TestAccIncapsulaApplicationDelivery_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationDeliveryBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationDeliveryExists(applicationDeliveryResource),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "file_compression", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "compression_type", "GZIP"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_css", "true"), //value wasn't set by tf resurce. checking default value from server
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_js", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_static_html", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "aggressive_compression", "true"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "compress_jpeg", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "compress_png", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "aggressive_compression", "true"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "enable_http2", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "http2_to_origin", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "origin_connection_reuse", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "port_to", "225"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "ssl_port_to", "443"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "support_non_sni_clients", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "tcp_pre_pooling", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_http_to_https", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_naked_to_full", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_access_denied", customErrorPageBasic),
				),
			},
			{
				Config: testAccCheckApplicationDeliveryIgnoreHttp2(domainName),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationDeliveryExists(applicationDeliveryResource),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "file_compression", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "compression_type", "GZIP"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_css", "true"), //value wasn't set by tf resurce. checking default value from server
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_js", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_static_html", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "aggressive_compression", "true"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "compress_jpeg", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "compress_png", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "aggressive_compression", "true"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "enable_http2", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "http2_to_origin", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "origin_connection_reuse", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "port_to", "225"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "ssl_port_to", "443"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "support_non_sni_clients", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "tcp_pre_pooling", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_http_to_https", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_naked_to_full", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_access_denied", customErrorPageBasic),
				),
			},
			{
				ResourceName:      applicationDeliveryResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateApplicationDeliveryID,
			},
		},
	})
}

func testCheckApplicationDeliveryExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Application Delivery resource not found: %s", name)
		}
		_, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Incapsula Application Delivery testCheckApplicationDeliveryExists: Error parsing ID %v to int", res.Primary.ID)
		}

		// client := testAccProvider.Meta().(*Client)
		// _, err = client.GetApplicationDelivery(siteId)
		// if err != nil {
		// 	fmt.Errorf("Incapsula Application Delivery doesn't exist")
		// }

		return nil
	}
}

func testACCStateApplicationDeliveryID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != applicationDeliveryResourceName {
			continue
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])

		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int in Application Delivery resource test", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("Error finding site_id argument in Application Delivery resource test")
}

func testAccCheckApplicationDeliveryBasic(domainName string) string {
	return testAccCheckIncapsulaSiteConfigBasic(domainName) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = incapsula_site.testacc-terraform-site.id
  depends_on = ["%s"]
  file_compression = true
  compression_type = "GZIP"
  compress_jpeg = true
  minify_static_html = false
  aggressive_compression = true
  progressive_image_rendering = false
  support_non_sni_clients = true
  enable_http2 = false
  http2_to_origin = false
  origin_connection_reuse = false
  port_to = 225
  tcp_pre_pooling = false
  redirect_naked_to_full = false
  redirect_http_to_https = false
  error_access_denied         = %s
}`,
		applicationDeliveryResourceName, applicationDeliveryName, siteResourceName, customErrorPageInput,
	)
}

func testAccCheckApplicationDeliveryIgnoreHttp2(domainName string) string {
	return testAccCheckIncapsulaSiteConfigBasic(domainName) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = incapsula_site.testacc-terraform-site.id
  depends_on = ["%s"]
  file_compression = true
  compression_type = "GZIP"
  compress_jpeg = true
  minify_static_html = false
  aggressive_compression = true
  progressive_image_rendering = false
  support_non_sni_clients = true
  origin_connection_reuse = false
  port_to = 225
  tcp_pre_pooling = false
  redirect_naked_to_full = false
  redirect_http_to_https = false
  error_access_denied         = %s
}`,
		applicationDeliveryResourceName, applicationDeliveryName, siteResourceName, customErrorPageInput,
	)
}
