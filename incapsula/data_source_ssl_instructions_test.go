package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
	"time"
)

const incapsulaSslInstructionsDataSource = "incapsula_ssl_instructions_example"
const dataSourceName = "data.incapsula_ssl_instructions." + incapsulaSslInstructionsDataSource

var siteV3ResourceNameForDomainTest = "test-site-v3-for-domain-resource" + strconv.FormatInt(time.Now().UnixNano()%99999, 10)

func TestIncapsulaSSL_Instructions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteV3Domain(t, "b-"+domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "domain_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "instructions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "instructions.0.type", "CNAME"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instructions.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instructions.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instructions.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instructions.0.san_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "managed_certificate_settings_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instructions.0.type"),
				),
			},
		},
	})
}

func testAccCheckIncapsulaSiteV3Domain(t *testing.T, domain string) string {
	result := checkIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`

 resource "incapsula_site_v3" "%s" {
			name = "%s"
	}

resource "incapsula_managed_certificate_settings" "example-site-cert" {
  site_id = incapsula_site_v3.%s.id
}


resource "incapsula_domain" "domain1" {
    site_id = incapsula_site_v3.%s.id
    domain="%s"
}

data "incapsula_ssl_instructions" "%s" {
domain_ids = toset([incapsula_domain.domain1.id])
managed_certificate_settings_id = incapsula_managed_certificate_settings.example-site-cert.id
site_id=incapsula_site_v3.%s.id

}
`, siteV3ResourceNameForDomainTest, siteV3NameForDomainTests, siteV3ResourceNameForDomainTest, siteV3ResourceNameForDomainTest, domain,
		incapsulaSslInstructionsDataSource, siteV3ResourceNameForDomainTest)
	return result
}
