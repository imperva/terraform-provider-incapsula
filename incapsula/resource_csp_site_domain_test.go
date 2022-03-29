package incapsula

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"strings"
	"testing"
)

const cspDomainResourceName = "incapsula_csp_site_domain"
const cspDomainResource = cspDomainResourceName + "." + cspDomainName
const cspDomainName = "testacc-terraform-csp-domain"

func TestAccIncapsulaCSPDomain_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_csp_site_domain_test.TestAccIncapsulaCSPDomain_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateCSPDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCSPDomainBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSPDomainExists(cspDomainResource),
					resource.TestCheckResourceAttr(cspDomainResource, "domain", "stam-domain.com"),
					resource.TestCheckResourceAttr(cspDomainResource, "include_subdomains", "true"),
					resource.TestCheckResourceAttr(cspDomainResource, "notes.#", "1"),
					resource.TestCheckResourceAttr(cspDomainResource, "notes.0", "new note"),
					resource.TestCheckResourceAttr(cspDomainResource, "status", "Allowed"),
				),
			},
			{
				ResourceName:      cspDomainResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateCSPDomainID,
			},
		},
	})
}

func testCheckCSPDomainExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula CSP domain resource not found: %s", name)
		}
		siteId, err := strconv.Atoi(res.Primary.Attributes["site_id"])
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.Attributes["site_id"])
		}

		client := testAccProvider.Meta().(*Client)
		cspDomain, err := client.getCSPPreApprovedDomain(siteId, res.Primary.Attributes["domain"])
		if err != nil || cspDomain == nil {
			return fmt.Errorf("Incapsula CSP domain %s doesn't exist for site ID %d", res.Primary.Attributes["domain"], siteId)
		}

		return nil
	}
}

func testACCStateCSPDomainID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		fmt.Errorf("Resource: %v", rs)
		if rs.Type != cspDomainResourceName {
			continue
		}

		keyParts := strings.Split(rs.Primary.ID, "/")
		if len(keyParts) != 2 {
			return "", fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", rs.Primary.ID)
		}
		keySiteID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric id", keyParts[0])
		}
		keyDomain, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(keyParts[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert domain reference ID from import command, actual value: %s, expected Base64 id", keyParts[1])
		}

		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		domain := rs.Primary.Attributes["domain"]

		if siteID != keySiteID || strings.Compare(domain, string(keyDomain)) != 0 {
			return "", fmt.Errorf("Incapsula CSP domain does not exist")
		}
		return fmt.Sprintf("%d/%s", siteID, keyParts[1]), nil
	}
	return "", fmt.Errorf("Error finding correct resource %s", cspDomainResourceName)
}

func testACCStateCSPDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != cspDomainResourceName {
			continue
		}

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Parameter site_id was not found in resource %s", cspDomainResourceName)
		}
		siteIDInt, err := strconv.Atoi(siteID)
		if err != nil {
			return fmt.Errorf("failed to convert site ID from import command, actual value : %s, expected numeric id", siteID)
		}
		domain := rs.Primary.Attributes["domain"]
		if domain == "" {
			return fmt.Errorf("Parameter domain was not found in resource %s", cspDomainResourceName)
		}

		cspDomain, err := client.getCSPPreApprovedDomain(siteIDInt, domain)

		fmt.Sprintf("Got CSP domain for site ID %d: %v", siteIDInt, cspDomain)
		if err != nil && cspDomain != nil {
			return fmt.Errorf("Resource %s for CSP domain: Api Id %s, site ID %d still exists", cspDomainResourceName, rs.Primary.ID, siteIDInt)
		}
		return nil
	}
	return fmt.Errorf("Error finding the correct resource: %s", cspDomainResourceName)
}

func testAccCheckCSPDomainBasic(t *testing.T) string {
	return testAccCheckCSPSiteConfigBasic(t) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id				= %s.id
		domain				= "stam-domain.com"
		include_subdomains	= true
		notes				= ["new note"]
		status				= "Allowed"
		depends_on			= ["%s"]
	}`,
		cspDomainResourceName, cspDomainName, cspSiteConfigResource, cspSiteConfigResource,
	)
}
