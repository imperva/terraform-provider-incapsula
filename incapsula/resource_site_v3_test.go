package incapsula

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const siteV3ResourceName = "incapsula_site_v3.test-terraform-site-v3"

var siteName string

func GenerateTestSiteName(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" && t != nil {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}

	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	siteName = "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000))
	return siteName
}

func TestIncapsulaSiteV3_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_refId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", "ref_id = \"123456\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siteV3ResourceName, "ref_id", "123456"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_isActive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", "active = false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siteV3ResourceName, "active", "false"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_AccountIdUpdateFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
				),
			},
			{
				Config:      testCheckIncapsulaSiteV3ConfigWithAccountId(siteName, "CLOUD_WAF", "999999"),
				ExpectError: regexp.MustCompile("account_id cannot be updated for an existing site"),
			},
		},
	})
}

func testCheckIncapsulaSiteV3Destroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site_v3" {
			continue
		}

		siteIDStr := res.Primary.ID
		if siteIDStr == "" {
			return fmt.Errorf("incapsula site v3 ID does not exist")
		}
		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			return fmt.Errorf("Site ID conversion error for %s: %s", siteIDStr, err)
		}

		siteV3Request := SiteV3Request{}
		siteV3Request.Name = siteName

		_, diags := client.GetV3Site(&siteV3Request, "123")

		if diags == nil {
			return fmt.Errorf("incapsula site for domain: %s (site id: %d) still exists", siteName, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaSiteV3ConfigBasic(name string, siteType string, extraAttr string) string {
	return fmt.Sprintf(`
		resource "incapsula_site_v3" "test-terraform-site-v3" {
			name = "%s"
		    type = "%s"
			%s
		}`,
		name,
		siteType,
		extraAttr,
	)
}

func testCheckIncapsulaSiteV3ConfigWithAccountId(name string, siteType string, accountId string) string {
	return fmt.Sprintf(`
		resource "incapsula_site_v3" "test-terraform-site-v3" {
			name = "%s"
		    type = "%s"
			account_id = "%s"
		}`,
		name,
		siteType,
		accountId,
	)
}

func testSiteV3Importer(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {

		accountId1, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing API ID %v to int", rs.Primary.Attributes["id"])
		}

		return fmt.Sprintf("%d/%s", accountId1, rs.Primary.ID), nil
	}
	return "", fmt.Errorf("Error finding an Site V3")
}
