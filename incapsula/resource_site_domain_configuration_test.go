package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"strconv"
	"testing"
	"time"
)

const siteDomainConfResourceName = "incapsula_site_domain_configuration"
const siteDomainConfResource = "site_domain_conf"
const rootModuleName = siteDomainConfResourceName + "." + siteDomainConfResource

var domain = "a-" + strconv.FormatInt(time.Now().UnixNano()%99999, 10) + ".examplewebsite.com"

func TestAccIncapsulaSiteDomainConfiguration_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteDomainConfGoodConfig(t, domain),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName),
					resource.TestMatchResourceAttr(rootModuleName, "cname_redirection_record", regexp.MustCompile(".+\\.imperva.+")),
					resource.TestMatchResourceAttr(rootModuleName, "site_id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.name", regexp.MustCompile(domain)),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.status", regexp.MustCompile("BYPASSED")),
				),
			},
		},
	})
}

func testCheckIncapsulaSiteDomainConfExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[rootModuleName]
		if !ok {
			return fmt.Errorf("incapsula site domain configuration resource not found : %s", siteDomainConfResource)
		}

		siteId, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("incapsula site domain configuration site_id %s does not exist", siteId)
		}

		client := testAccProvider.Meta().(*Client)
		siteDomainDetailsDto, _ := client.GetWebsiteDomains(siteId)
		if siteDomainDetailsDto == nil {
			return fmt.Errorf("incapsula site domain configuration: get domains for siteId %s response returned null", siteId)
		}
		return nil
	}
}

func testAccCheckIncapsulaSiteDomainConfGoodConfig(t *testing.T, domain string) string {
	result := checkIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id=incapsula_site.testacc-terraform-site.id
  domain {name="%s"}
depends_on = ["%s"]
}`, siteDomainConfResourceName, siteDomainConfResource, domain, siteResourceName)
	return result
}

func checkIncapsulaSiteConfigBasic(domain string) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "testacc-terraform-site" {
			domain = "%s"
		}`,
		domain,
	)
}
