package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

const siteDomainConfResourceName = "incapsula_site_domain_configuration"
const siteDomainConfResource = "site_domain_conf"
const rootModuleName = siteDomainConfResourceName + "." + siteDomainConfResource

var siteV3ResourceNameForDomainTest = "test-site-v3-for-domain-resource" + strconv.FormatInt(time.Now().UnixNano()%99999, 10)
var siteV3NameForDomainTests = "test site for domain resource" + strconv.FormatInt(time.Now().UnixNano()%99999, 10)
var domain = strconv.FormatInt(time.Now().UnixNano()%99999, 10) + os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")

func TestAccIncapsulaSiteDomainConfiguration_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				SkipFunc: IsTestDomainEnvVarExist,
				Config:   testAccCheckIncapsulaSiteDomainConfGoodConfig(t, domain),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName),
					resource.TestMatchResourceAttr(rootModuleName, "cname_redirection_record", regexp.MustCompile(".+\\.imperva.+")),
					resource.TestMatchResourceAttr(rootModuleName, "site_id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.name", regexp.MustCompile(domain)),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName, "domain.0.status", regexp.MustCompile("BYPASSED")),
				),
			},
			{
				SkipFunc: IsTestDomainEnvVarExist,
				Config:   testAccCheckIncapsulaSiteV3Domain(t, "b-"+domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(rootModuleName+"-2", "site_id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName+"-2", "domain.0.name", regexp.MustCompile("b-"+domain)),
					resource.TestMatchResourceAttr(rootModuleName+"-2", "domain.0.id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootModuleName+"-2", "domain.0.status", regexp.MustCompile("BYPASSED")),
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

func testAccCheckIncapsulaSiteV3Domain(t *testing.T, domain string) string {
	result := checkIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`

 resource "incapsula_site_v3" "%s" {
			name = "%s"
	}

resource "incapsula_site_certificate_request" "example-site-cert" {
  site_id = incapsula_site_v3.%s.id
}
resource "%s" "%s-2" {
  site_id=incapsula_site_v3.%s.id
  site_certificate_request_id = incapsula_site_certificate_request.example-site-cert.id
  domain {name="%s"}
}`, siteV3ResourceNameForDomainTest, siteV3NameForDomainTests, siteV3ResourceNameForDomainTest, siteDomainConfResourceName, siteDomainConfResource, siteV3ResourceNameForDomainTest, domain)
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
