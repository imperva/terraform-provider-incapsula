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

const domainResourceName = "incapsula_domain.test-terraform-domain"

var domainName string

func GenerateTestDomainName(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" && t != nil {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}

	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	domainName = "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000)) + ".incaptest.co"
	return domainName
}

// wont work until we will support delete last domain
func testIncapsulaDomain_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		//	CheckDestroy: testCheckIncapsulaDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaDomainBasic(GenerateTestDomainName(nil), "CLOUD_WAF"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(domainResourceName, "domain", domainName),
				),
			},
			{
				ResourceName:      domainResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDomainID("incapsula_domain"),
			},
		},
	})
}

func testDomainID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["site_id"], rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find domain ID")
	}
}
func testCheckIncapsulaDomainDestroy(state *terraform.State) error {
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

func testCheckIncapsulaDomainBasic(domainName string, siteType string) string {
	result := fmt.Sprintf(`

resource "incapsula_site_v3" "test-terraform-site" {
			name = "%s"
		    type = "%s"
		}

resource "incapsula_domain" "test-terraform-domain" {
	site_id = incapsula_site_v3.test-terraform-site.id
	domain = "%s"

}`, domainName, siteType, domainName)
	return result
}
