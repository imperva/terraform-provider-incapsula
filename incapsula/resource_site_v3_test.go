package incapsula

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"math/rand"
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
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF"),
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

func testCheckIncapsulaSiteV3ConfigBasic(name string, siteType string) string {
	return fmt.Sprintf(`
		resource "incapsula_site_v3" "test-terraform-site-v3" {
			name = "%s"
		    type = "%s"

		}`,
		name,
		siteType,
	)
}
