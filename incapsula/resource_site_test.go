package incapsula

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"math/rand"
)

const siteResourceName = "incapsula_site.testacc-terraform-site"

var generatedDomain string

func GenerateTestDomain(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" && t != nil {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}
	if v := os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN"); v == "" && t != nil {
		t.Fatal("INCAPSULA_CUSTOM_TEST_DOMAIN must be set for acceptance tests which require onboarding a domain")
	}
	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	generatedDomain = "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000)) + os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")
	return generatedDomain
}

func TestAccIncapsulaSite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaSiteDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: IsTestDomainEnvVarExist,
				Config:   testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(nil)),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteResourceName),
					resource.TestCheckResourceAttr(siteResourceName, "domain", generatedDomain),
				),
			},
			{
				ResourceName:            "incapsula_site.testacc-terraform-site",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site_ip", "domain_validation"},
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

		_, err = client.SiteStatus(generatedDomain, siteID)

		if err == nil {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %d) still exists", GenerateTestDomain(nil), siteID)
		}
	}

	return nil
}

func IsTestDomainEnvVarExist() (bool, error) {
	skipTest := false
	if v := os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN"); v == "" {
		skipTest = true
		log.Printf("[ERROR] INCAPSULA_CUSTOM_TEST_DOMAIN environment variable does not exist, if you want to if you want to test features which require site onboarding, you must prvide it")
	}

	return skipTest, nil
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
