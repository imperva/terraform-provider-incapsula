package incapsula

import (
	"encoding/json"
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
const rootSiteDomainConfigurationModuleName = siteDomainConfResourceName + "." + siteDomainConfResource

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
					testCheckIncapsulaSiteDomainConfExists(rootSiteDomainConfigurationModuleName),
					resource.TestMatchResourceAttr(rootSiteDomainConfigurationModuleName, "cname_redirection_record", regexp.MustCompile(".+\\.imperva.+")),
					resource.TestMatchResourceAttr(rootSiteDomainConfigurationModuleName, "site_id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootSiteDomainConfigurationModuleName, "domain.0.name", regexp.MustCompile(domain)),
					resource.TestMatchResourceAttr(rootSiteDomainConfigurationModuleName, "domain.0.id", regexp.MustCompile("\\d+")),
					resource.TestMatchResourceAttr(rootSiteDomainConfigurationModuleName, "domain.0.status", regexp.MustCompile("BYPASSED")),
				),
			},
		},
	})
}

func testCheckIncapsulaSiteDomainConfExists(fullSiteDomainResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[fullSiteDomainResourceName]
		if !ok {
			return fmt.Errorf("incapsula site domain configuration resource not found : %s", fullSiteDomainResourceName)
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
	result := checkIncapsulaSiteConfigBasic(GenerateTestDomain(t), "testacc-terraform-site", false) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id=incapsula_site.testacc-terraform-site.id
  domain {name="%s"}
depends_on = ["%s"]
}`, siteDomainConfResourceName, siteDomainConfResource, domain, siteResourceName)
	return result
}

func checkIncapsulaSiteConfigBasic(domain string, siteResourceName string, deprecated bool) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "%s" {
			domain = "%s"
            deprecated = %t
		}`,
		siteResourceName, domain, deprecated,
	)
}

func testAccCheckIncapsulaSiteDomainConfigDeprecated(siteDomain string, siteResourceName string, siteDomainConfigResourceName string, domain string, deprecatedDomain bool, deprecatedSite bool) string {
	result := checkIncapsulaSiteConfigBasic(siteDomain, siteResourceName, deprecatedSite) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id=incapsula_site.%s.id
  domain {name="%s"}
  deprecated = %t
depends_on = ["incapsula_site.%s"]
}`, siteDomainConfResourceName, siteDomainConfigResourceName, siteResourceName, domain, deprecatedDomain, siteResourceName)
	return result
}

func testAccCheckIncapsulaSiteDomainConfigDeprecatedMultiDomains(siteDomain string, siteResourceName string, siteDomainConfigResourceName string, domain1 string, domain2 string, deprecatedDomain bool, deprecatedSite bool) string {
	result := checkIncapsulaSiteConfigBasic(siteDomain, siteResourceName, deprecatedSite) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id=incapsula_site.%s.id
  domain {name="%s"}
  domain {name="%s"}
  deprecated = %t
depends_on = ["incapsula_site.%s"]
}`, siteDomainConfResourceName, siteDomainConfigResourceName, siteResourceName, domain1, domain2, deprecatedDomain, siteResourceName)
	return result
}

func TestAccIncapsulaSiteDomainConfig_DeprecationFlag_siteDomainConfigNotCreated(t *testing.T) { //done
	testDomain := GenerateTestDomain(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIncapsulaSiteDomainConfigDeprecated(GenerateTestDomain(t), "TestAccIncapsulaSiteDomainConfig-site1", "TestAccIncapsulaSiteDomainConfig-domain1", testDomain, true, true),
				ExpectError: regexp.MustCompile("cannot create deprecated resource"),
			},
		},
	})
}

func TestAccIncapsulaSiteDomainConfig_ChangeDeprecatedFlag(t *testing.T) {
	siteTestDomain := GenerateTestDomain(t) //done
	testDomain := GenerateTestDomain(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteDomainConfNotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteDomainConfigDeprecated(siteTestDomain, "TestAccIncapsulaSiteDomainConfig-site2", "TestAccIncapsulaSiteDomainConfig-domain2", testDomain, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2"),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2", "domain.0.name", testDomain),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteDomainConfigDeprecated(siteTestDomain, "TestAccIncapsulaSiteDomainConfig-site2", "TestAccIncapsulaSiteDomainConfig-domain2", testDomain, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2"),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2", "domain.0.name", testDomain),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteDomainConfigDeprecated(siteTestDomain, "TestAccIncapsulaSiteDomainConfig-site2", "TestAccIncapsulaSiteDomainConfig-domain2", testDomain, false, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2"),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain2", "domain.0.name", testDomain),
				),
				ExpectError: regexp.MustCompile("deprecated flag cannot be changed from true to false"),
			},
		},
	})
}

func TestAccIncapsulaSiteDomainConfig_DeprecationFlagChangeAttributes(t *testing.T) {
	siteTestDomain := GenerateTestDomain(t)
	domain1 := GenerateTestDomain(t)
	domain2 := GenerateTestDomain(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckIncapsulaSiteDomainConfNotDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteDomainConfigDeprecatedMultiDomains(siteTestDomain, "TestAccIncapsulaSiteDomainConfig-site3", "TestAccIncapsulaSiteDomainConfig-domain3", domain1, domain2, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3"),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3", "domain.0.name", domain2),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3", "domain.1.name", domain1),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteDomainConfigDeprecatedMultiDomains(siteTestDomain, "TestAccIncapsulaSiteDomainConfig-site3", "TestAccIncapsulaSiteDomainConfig-domain3", "c.example.com", "d.example.com", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteDomainConfExists(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3"),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3", "domain.0.name", domain2),
					resource.TestCheckResourceAttr(siteDomainConfResourceName+".TestAccIncapsulaSiteDomainConfig-domain3", "domain.1.name", domain1),
				),
			},
		},
	})
}

func testCheckIncapsulaSiteDomainConfNotDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, res := range state.RootModule().Resources {
		if res.Type != siteDomainConfResourceName && res.Primary.Attributes["name"] != "TestAccIncapsulaSiteDomainConfig-domain2" {
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

		domain := res.Primary.Attributes["domain"]
		siteStatusResponse, err := client.SiteStatus(domain, siteID)

		if err != nil {
			return fmt.Errorf("Site status error for %s: %s", siteIDStr, err)
		}

		if siteStatusResponse == nil || siteStatusResponse.SiteID != siteID {
			response, _ := json.Marshal(siteStatusResponse)
			return fmt.Errorf("Incapsula site status for domain: %s (site id: %d) was not retreived. response: %s", domain, siteID, response)
		}
		err = client.DeleteSite(domain, siteID) // clean the env after checking

		if err != nil {
			return fmt.Errorf("Error deleting site (%d) for domain %s: %s", siteID, domain, err)
		}
	}

	return nil
}
