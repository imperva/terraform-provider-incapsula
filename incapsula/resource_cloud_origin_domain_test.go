package incapsula

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIncapsulaCloudOriginDomainBasic(t *testing.T) {
	testName := "tf-test-cloud-origin-basic"
	domain := fmt.Sprintf("%s.example.com", testName)
	resourceName := "incapsula_cloud_origin_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudOriginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudOriginDomainConfigBasic(testName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudOriginDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttrSet(resourceName, "imperva_origin_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}

func TestAccIncapsulaCloudOriginDomainUpdate(t *testing.T) {
	testName := "tf-test-cloud-origin-update"
	domain := fmt.Sprintf("%s.example.com", testName)
	resourceName := "incapsula_cloud_origin_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudOriginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudOriginDomainConfigBasic(testName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudOriginDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
				),
			},
			{
				Config: testAccCheckCloudOriginDomainConfigUpdate(testName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudOriginDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "port", "8443"),
				),
			},
		},
	})
}

func TestAccIncapsulaCloudOriginDomainImport(t *testing.T) {
	testName := "tf-test-cloud-origin-import"
	domain := fmt.Sprintf("%s.example.com", testName)
	resourceName := "incapsula_cloud_origin_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudOriginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudOriginDomainConfigBasic(testName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudOriginDomainExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					// Import ID format: site_id/origin_id
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccIncapsulaCloudOriginDomainDelete(t *testing.T) {
	testName := "tf-test-cloud-origin-delete"
	domain := fmt.Sprintf("%s.example.com", testName)
	resourceName := "incapsula_cloud_origin_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudOriginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudOriginDomainConfigBasic(testName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudOriginDomainExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckCloudOriginDomainExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Cloud origin domain ID is not set")
		}

		// Parse the composite ID
		if !strings.Contains(rs.Primary.ID, "/") {
			return fmt.Errorf("Invalid cloud origin domain ID format: %s", rs.Primary.ID)
		}

		parts := strings.Split(rs.Primary.ID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid cloud origin domain ID format: %s", rs.Primary.ID)
		}

		siteID, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("Invalid site ID in resource ID: %s", parts[0])
		}

		originID, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Invalid origin ID in resource ID: %s", parts[1])
		}

		client := testAccProvider.Meta().(*Client)

		response, err := client.GetCloudOriginDomain(siteID, 0, originID)
		if err != nil {
			return fmt.Errorf("Error getting cloud origin domain: %s", err)
		}

		if response == nil {
			return fmt.Errorf("Cloud origin domain not found")
		}

		return nil
	}
}

func testAccCheckCloudOriginDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_cloud_origin_domain" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		// Parse the composite ID
		if !strings.Contains(rs.Primary.ID, "/") {
			continue
		}

		parts := strings.Split(rs.Primary.ID, "/")
		if len(parts) != 2 {
			continue
		}

		siteID, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		originID, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		// Try to get the cloud origin domain
		response, err := client.GetCloudOriginDomain(siteID, 0, originID)
		if err == nil && response != nil {
			return fmt.Errorf("Cloud origin domain still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCloudOriginDomainConfigBasic(testName, domain string) string {
	return fmt.Sprintf(`
resource "incapsula_site" "test" {
  domain = "%s"
}

resource "incapsula_cloud_origin_domain" "test" {
  site_id = incapsula_site.test.id
  domain  = "%s"
  region  = "us-east-1"
  port    = 443
}
`, testName+".com", domain)
}

func testAccCheckCloudOriginDomainConfigUpdate(testName, domain string) string {
	return fmt.Sprintf(`
resource "incapsula_site" "test" {
  domain = "%s"
}

resource "incapsula_cloud_origin_domain" "test" {
  site_id = incapsula_site.test.id
  domain  = "%s"
  region  = "eu-west-1"
  port    = 8443
}
`, testName+".com", domain)
}
