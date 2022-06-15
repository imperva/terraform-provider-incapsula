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

const customErrorPage = "<<-EOT\n<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n    <link href=\"https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;700&display=swap\" rel=\"stylesheet\">\n    <title>[Error Title]</title>\n      </head>\n  <body>\n    <div class=\"container\">\n      <div class=\"container-inner\">\n        <div class=\"header\">\n          <div class=\"error-description\">\n            $TITLE$\n          </div>\n        </div>\n        <div class=\"main\">\n          <div class=\"troubleshooting\">\n            <div class=\"content\">\n              $BODY$\n            </div>\n\t    <h1>custom edited error</h1>\n          </div>\n        </div>\n      </div>\n    </div>\n  </body>\n</html>\nEOT"
const customErrorPageRes = "<!DOCTYPE html>\\n<html lang=\\\"en\\\">\\n  <head>\\n    <meta charset=\\\"UTF-8\\\">\\n    <meta name=\\\"viewport\\\" content=\\\"width=device-width, initial-scale=1.0\\\">\\n    <meta http-equiv=\\\"X-UA-Compatible\\\" content=\\\"ie=edge\\\">\\n    <link href=\\\"https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;700&display=swap\\\" rel=\\\"stylesheet\\\">\\n    <title>[Error Title]</title>\\n      </head>\\n  <body>\\n    <div class=\\\"container\\\">\\n      <div class=\\\"container-inner\\\">\\n        <div class=\\\"header\\\">\\n          <div class=\\\"error-description\\\">\\n            $TITLE$\\n          </div>\\n        </div>\\n        <div class=\\\"main\\\">\\n          <div class=\\\"troubleshooting\\\">\\n            <div class=\\\"content\\\">\\n              $BODY$\\n            </div>\\n\\t    <h1>custom edited error</h1>\\n          </div>\\n        </div>\\n      </div>\\n    </div>\\n  </body>\\n</html>"

func TestAccIncapsulaApplicationDelivery_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_application_delivery_test.TestAccIncapsulaApplicationDelivery_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testACCStateApplicationDeliveryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationDeliveryBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationDeliveryExists(applicationDeliveryResource),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "file_compression", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "minify_css", "true"),
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
					resource.TestCheckResourceAttr(applicationDeliveryResource, "ssl_port_to", "555"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "support_non_sni_clients", "true"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "tcp_pre_pooling", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_http_to_https", "false"),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "redirect_naked_to_full", "false"),

					resource.TestCheckResourceAttr(applicationDeliveryResource, "default_error_page_template", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_access_denied", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_connection_failed", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_connection_timeout", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_parse_req_error", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_parse_resp_error", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_ssl_failed", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_deny_and_captcha", customErrorPageRes),
					resource.TestCheckResourceAttr(applicationDeliveryResource, "error_no_ssl_config", customErrorPageRes),
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

func testACCStateApplicationDeliveryDestroy(s *terraform.State) error {
	log.Printf("Destroy:state has resources:\n%v", s.RootModule().Resources)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != applicationDeliveryResourceName {
			continue
		} else {
			return fmt.Errorf("Resource %s for Incapsula Application Delivery: resource still exists\nState:\n%v", applicationDeliveryResourceName, s)
		}
		//return nil
	}
	return nil
}

func testCheckApplicationDeliveryExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Application Delivery resource not found: %s", name)
		}
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Incapsula Application Delivery testCheckApplicationDeliveryExists: Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		_, err = client.GetApplicationDelivery(siteId)
		if err != nil {
			fmt.Errorf("Incapsula Application Delivery doesn't exist")
		}

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

func testAccCheckApplicationDeliveryBasic(t *testing.T) string {
	//site_id = 66107946
	//site_id = incapsula_site.testacc-terraform-site.id
	//depends_on = ["%s"]
	//return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	return fmt.Sprintf(`
resource "%s" "%s" {
  site_id = 53427201
  file_compression = true
  compress_jpeg = true
  minify_js = true
  minify_css = true
  minify_static_html = false
  aggressive_compression = true
  progressive_image_rendering = false
  support_non_sni_clients = true
  enable_http2 = false
  http2_to_origin = false
  origin_connection_reuse = false
  port_to = 225
  tcp_pre_pooling = false
  ssl_port_to = 555
  redirect_naked_to_full = false
  redirect_http_to_https = false
  default_error_page_template = %s
  error_access_denied         = %s
  error_connection_failed     = %s
  error_connection_timeout     = %s
  error_parse_req_error     = %s
  error_parse_resp_error     = %s
  error_ssl_failed     = %s
  error_deny_and_captcha     = %s
  error_no_ssl_config     = %s
}`,
		applicationDeliveryResourceName, applicationDeliveryName, customErrorPage, customErrorPage, customErrorPage,
		customErrorPage, customErrorPage, customErrorPage, customErrorPage, customErrorPage, customErrorPage,
	)
}
