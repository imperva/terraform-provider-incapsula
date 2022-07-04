package incapsula

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"math/rand"
	"time"
)

const siteResourceName = "incapsula_site.testacc-terraform-site"

func GenerateTestDomain(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" && t != nil {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}
	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	return "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000)) + ".examplesite.com"
}

func TestAccIncapsulaSite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteResourceName),
					resource.TestCheckResourceAttr(siteResourceName, "domain", GenerateTestDomain(t)),
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
				Config: testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)),
			},
			{
				ResourceName:            "incapsula_site.testacc-terraform-site",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site_ip"},
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

		_, err = client.SiteStatus(GenerateTestDomain(nil), siteID)

		if err == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) still exists", GenerateTestDomain(nil), siteID)
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
		siteStatusResponse, err := client.SiteStatus(GenerateTestDomain(nil), siteID)
		if siteStatusResponse == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) does not exist", GenerateTestDomain(nil), siteID)
		}

		return nil
	}
}

func testAccCheckIncapsulaSiteConfigBasic(domain string) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "testacc-terraform-site" {
			domain = "%s"
		}`,
		domain,
	)
}
