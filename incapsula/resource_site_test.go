package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDomain = "www.examplesite.com"
const siteResourceName = "incapsula_site.testacc-terraform-site"

func TestAccIncapsulaSite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfig_basic(testAccDomain),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteResourceName),
					resource.TestCheckResourceAttr(siteResourceName, "domain", testAccDomain),
				),
			},
		},
	})
}

func TestAccIncapsulaSite_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfig_basic(testAccDomain),
			},
			{
				ResourceName:      "incapsula_site.testacc-terraform-site",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIncapsulaSiteDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
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

		_, err = client.SiteStatus(testAccDomain, siteID)

		if err == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) still exists", testAccDomain, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaSiteExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula site resource not found: %s", name)
		}

		siteIDStr := res.Primary.ID
		if siteIDStr == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			return fmt.Errorf("Site ID conversion error for %s: %s", siteIDStr, err)
		}

		client := testAccProvider.Meta().(*Client)
		siteStatusResponse, err := client.SiteStatus(testAccDomain, siteID)
		if siteStatusResponse == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) does not exist", testAccDomain, siteID)
		}

		return nil
	}
}

func testAccCheckIncapsulaSiteConfig_basic(domain string) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "testacc-terraform-site" {
			domain = "%s"
		}`,
		domain,
	)
}
