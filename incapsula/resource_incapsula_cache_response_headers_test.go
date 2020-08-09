package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const cacheHeaderResponseResourceName = "incapsula_cache_response_headers.testacc-terraform-cache-response-header"
const cacheHeaderResponseName = "server"

func TestAccIncapsulaCacheHeaderResponse_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCacheHeaderResponseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testacccheckincapsulacacheheaderresponseconfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCacheHeaderResponseExists(cacheHeaderResponseResourceName),
					resource.TestCheckResourceAttr(cacheHeaderResponseResourceName, "cache_headers", cacheHeaderResponseName),
				),
			},
			{
				ResourceName:      cacheHeaderResponseResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCheckStatecacheHeaderResponseId,
			},
		},
	})
}

func testAccCheckStatecacheHeaderResponseId(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_cache_response_headers" {
			continue
		}

		cacheResponseId, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, cacheResponseId), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaCacheHeaderResponseDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteIDstr := res.Primary.ID
		if siteIDstr == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		listCacheHeaderResponse, err := client.SiteStatus("cache-response-header", siteID)
		for _, cacheHeader := range listCacheHeaderResponse.PerformanceConfiguration.CacheHeaders {
			if cacheHeader == cacheHeaderResponseName {
				return fmt.Errorf("Incapsula cache header response: %s (site_id: %s) still exists", cacheHeaderResponseName, siteIDstr)
			}
		}
		if err == nil {
			return fmt.Errorf("Incapsula cache header response for: %s (site id: %s) still exists", testAccDomain, siteIDstr)
		}
	}
	return nil
}

func testCheckIncapsulaCacheHeaderResponseExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		siteRes, siteOk := state.RootModule().Resources[siteResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula site resource not found: %s", siteResourceName)
		}

		siteID, err := strconv.Atoi(siteRes.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", siteRes.Primary.ID)
		}

		siteIDstr := siteRes.Primary.ID
		if siteIDstr == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula cache header response resource not found: %s", name)
		}

		cacheResponseId := res.Primary.ID
		if cacheResponseId == "" {
			return fmt.Errorf("Incapsula cache header response ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)
		listCacheHeaderResponse, err := client.SiteStatus("read-cache-responses", siteID)
		if listCacheHeaderResponse == nil {
			return fmt.Errorf("Incapsula cache header response: %s (site id: %s) does not exist\n%s", name, siteIDstr, err)
		}
		for _, cacheHeader := range listCacheHeaderResponse.PerformanceConfiguration.CacheHeaders {
			if cacheHeader == cacheHeaderResponseName {
				return nil
			}
		}
		return nil
	}
}

func testAccCheckIncapsulaCacheHeaderResponseConfigBasic() string {
	return testAccCheckIncapsulaSiteConfigBasic(testAccDomain) + fmt.Sprintf(`
resource "incapsula_cache_response_headers" "testacc-terraform-cache-response-header" {
  site_id           = "${incapsula_site.testacc-terraform-site.id}"
  cache_headers     = "%s"
  depends_on        = ["%s"]
}`, cacheHeaderResponseName, siteResourceName,
	)
}
