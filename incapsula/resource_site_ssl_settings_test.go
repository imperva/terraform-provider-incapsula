package incapsula

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// NOTE: SSL settings endpoint requires the site to have SSL configured with a custom certificate.
// This test uploads a self-signed certificate before configuring SSL settings.
const sslSettingsResourceType = "incapsula_site_ssl_settings"
const sslSettingsResourceName = "testacc-terraform-site-ssl-settings"
const sslSettingsFullResourceName = sslSettingsResourceType + "." + sslSettingsResourceName

func TestAccSiteSSLSettings_Basic(t *testing.T) {
	domainName = GenerateTestDomain(t)
	cert, pkey := generateKeyPairBase64(domainName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSiteCheckIncapsulaSiteSSLSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSSLSettingsConfig(domainName, cert, pkey, 31536000, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapSiteSSLSettingsExists(),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.max_age", "31536000"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.sub_domains_included", "true"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.pre_loaded", "true"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "inbound_tls_settings.0.configuration_profile", "CUSTOM"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "inbound_tls_settings.0.tls_configuration.0.tls_version", "TLS_1_3"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "inbound_tls_settings.0.tls_configuration.0.ciphers_support.0", "TLS_AES_128_GCM_SHA256"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "inbound_tls_settings.0.tls_configuration.0.ciphers_support.1", "TLS_AES_256_GCM_SHA384"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "disable_pqc_support", "true"),
				),
			},
		},
	})
}

func TestAccSiteSSLSettings_Update(t *testing.T) {
	domainName = GenerateTestDomain(t)
	cert, pkey := generateKeyPairBase64(domainName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSiteCheckIncapsulaSiteSSLSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSSLSettingsConfig(domainName, cert, pkey, 31536000, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapSiteSSLSettingsExists(),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.max_age", "31536000"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.sub_domains_included", "true"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.pre_loaded", "true"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "disable_pqc_support", "true"),
				),
			},
			{
				Config: testSSLSettingsConfig(domainName, cert, pkey, 86400, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaIncapSiteSSLSettingsExists(),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.max_age", "86400"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.sub_domains_included", "false"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "hsts.0.pre_loaded", "false"),
					resource.TestCheckResourceAttr(sslSettingsFullResourceName, "disable_pqc_support", "false"),
				),
			},
		},
	})
}

func testAccSiteCheckIncapsulaSiteSSLSettingsDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site_ssl_settings" {
			continue
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID == "" {
			return fmt.Errorf("incapsula Site ID does not exist for SSL settings")
		}

		accountID, ok := res.Primary.Attributes["account_id"]
		if !ok || accountID == "" {
			return fmt.Errorf("incapsula Account ID does not exist for ssl settings site id: %s", siteID)
		}

		siteIDToInt, err := strconv.Atoi(siteID)
		if err != nil {
			return fmt.Errorf("failed to parse site_id %s: %v", siteID, err)
		}
		accountIDToInt, err := strconv.Atoi(accountID)
		if err != nil {
			return fmt.Errorf("failed to parse account_id %s: %v", accountID, err)
		}

		_, statusCode, err := client.ReadSiteSSLSettings(siteIDToInt, accountIDToInt)
		if statusCode != 404 {
			return fmt.Errorf("incapsula Incap Site ssl settings (site id: %s) should have received 404 status code", siteID)
		}
		if err == nil {
			return fmt.Errorf("incapsula Incap Site ssl settings still exists for Site ID %s", siteID)
		}
	}

	return nil
}

func testCheckIncapsulaIncapSiteSSLSettingsExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[sslSettingsFullResourceName]
		if !ok {
			return fmt.Errorf("incapsula Site SSL settings resource not found")
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok || siteID == "" {
			return fmt.Errorf("incapsula Site ID does not exist for ssl settings")
		}

		accountID, ok := res.Primary.Attributes["account_id"]
		if !ok || accountID == "" {
			return fmt.Errorf("incapsula Account ID does not exist for ssl settings, site id: %s", siteID)
		}

		siteIDToInt, err := strconv.Atoi(siteID)
		if err != nil {
			return fmt.Errorf("failed to parse site_id %s: %v", siteID, err)
		}
		accountIDToInt, err := strconv.Atoi(accountID)
		if err != nil {
			return fmt.Errorf("failed to parse account_id %s: %v", accountID, err)
		}

		client := testAccProvider.Meta().(*Client)
		_, statusCode, err := client.ReadSiteSSLSettings(siteIDToInt, accountIDToInt)
		if statusCode != 200 {
			return fmt.Errorf("incapsula site ssl settings (site id: %s) should have received 200 status code", siteID)
		}
		if err != nil {
			return fmt.Errorf("incapsula site ssl settings (site id: %s) does not exist", siteID)
		}

		return nil
	}
}

func testSSLSettingsConfig(domain, cert, pkey string, maxAge int, subDomainsIncluded, preLoaded, disablePQC bool) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "testacc-terraform-site" {
			domain = "%s"
		}

		resource "incapsula_custom_certificate" "ssl-settings-test-certificate" {
			site_id = incapsula_site.testacc-terraform-site.id
			certificate = "%s"
			private_key = "%s"
			depends_on = ["incapsula_site.testacc-terraform-site"]
		}

		resource "incapsula_site_ssl_settings" "testacc-terraform-site-ssl-settings" {
		  site_id = incapsula_site.testacc-terraform-site.id
		  account_id = incapsula_site.testacc-terraform-site.account_id
		  hsts {
			is_enabled           = true
			max_age              = %d
			sub_domains_included = %t
			pre_loaded           = %t
		  }
		  inbound_tls_settings {
			configuration_profile = "CUSTOM"
			tls_configuration {
				tls_version = "TLS_1_3"
				ciphers_support = [
					"TLS_AES_128_GCM_SHA256",
					"TLS_AES_256_GCM_SHA384"
				]
			}
		  }
		  disable_pqc_support = %t
		  depends_on = ["incapsula_custom_certificate.ssl-settings-test-certificate"]
		}
`, domain, cert, pkey, maxAge, subDomainsIncluded, preLoaded, disablePQC)
}
