package incapsula

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIncapsulaDataSourceSite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceSiteBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test", "site_id", "incapsula_site_v3.test", "id"),
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test", "name", "incapsula_site_v3.test", "name"),
					resource.TestCheckResourceAttrSet("data.incapsula_site_lookup.test", "cname"),
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test", "type", "incapsula_site_v3.test", "type"),
				),
			},
		},
	})
}

func TestAccIncapsulaDataSourceSite_WithSubAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceSiteWithAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test_account", "site_id", "incapsula_site_v3.test_account", "id"),
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test_account", "name", "incapsula_site_v3.test_account", "name"),
					resource.TestCheckResourceAttrSet("data.incapsula_site_lookup.test_account", "cname"),
					resource.TestCheckResourceAttrPair("data.incapsula_site_lookup.test_account", "type", "incapsula_site_v3.test_account", "type"),
				),
			},
		},
	})
}

func TestAccIncapsulaDataSourceSite_CrossStateReference(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaDataSourceSiteCrossStateReference(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.incapsula_site_lookup.test", "site_id"),
					resource.TestCheckResourceAttrPair("incapsula_cache_rule.example", "site_id", "data.incapsula_site_lookup.test", "site_id"),
				),
			},
		},
	})
}

func testAccCheckIncapsulaDataSourceSiteBasic() string {
	r := rand.New(rand.NewSource(int64(os.Getuid())))
	randomNum := strconv.Itoa(r.Intn(1000))
	customTestDomain := os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")
	if customTestDomain == "" {
		customTestDomain = "testdom.com"
	}

	return fmt.Sprintf(`
resource "incapsula_site_v3" "test" {
  name = "example%s.%s"
  type = "Managed"
}

data "incapsula_site_lookup" "test" {
  name = incapsula_site_v3.test.name

  depends_on = [incapsula_site_v3.test]
}
`, randomNum, customTestDomain)
}

func testAccCheckIncapsulaDataSourceSiteWithAccount() string {
	r := rand.New(rand.NewSource(int64(os.Getuid())))
	randomNum := strconv.Itoa(r.Intn(1000))
	customTestDomain := os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")
	if customTestDomain == "" {
		customTestDomain = "testdom.com"
	}

	return fmt.Sprintf(`
resource "incapsula_site_v3" "test_account" {
  name = "example%s.%s"
  type = "Managed"
}

data "incapsula_site_lookup" "test_account" {
  name = incapsula_site_v3.test_account.name

  depends_on = [incapsula_site_v3.test_account]
}
`, randomNum, customTestDomain)
}

// Optional: Test with explicit account_id (if you have test sub-accounts available)
// Uncomment and configure when testing with multiple accounts
// func testAccCheckIncapsulaDataSourceSiteWithExplicitAccount() string {
//   return fmt.Sprintf(`
// data "incapsula_site_lookup" "test" {
//   name = "existing-site-name"
//   account_id = <your-test-account-id>
// }
// `)
// }

func testAccCheckIncapsulaDataSourceSiteCrossStateReference() string {
	r := rand.New(rand.NewSource(int64(os.Getuid())))
	randomNum := strconv.Itoa(r.Intn(1000))
	customTestDomain := os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")
	if customTestDomain == "" {
		customTestDomain = "testdom.com"
	}

	return fmt.Sprintf(`
resource "incapsula_site_v3" "test" {
  name = "example%s.%s"
  type = "Managed"
}

data "incapsula_site_lookup" "test" {
  name = incapsula_site_v3.test.name

  depends_on = [incapsula_site_v3.test]
}

resource "incapsula_cache_rule" "example" {
  site_id = data.incapsula_site_lookup.test.site_id
  name    = "test-cache-rule-%s"
  action  = "HTTP_CACHE_MAKE_STATIC"
  enabled = true
  filter  = "isMobile == Yes"
  ttl     = 3600
}
`, randomNum, customTestDomain, randomNum)
}
