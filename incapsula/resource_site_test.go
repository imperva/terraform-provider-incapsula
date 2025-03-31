package incapsula

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	initialDomain := "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000)) + os.Getenv("INCAPSULA_CUSTOM_TEST_DOMAIN")
	generatedDomain = strings.ReplaceAll(initialDomain, " ", "")
	log.Printf("[DEBUG] Generated domain: %s", generatedDomain)
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

func TestAccIncapsulaSite_DeprecationFlag_siteNotCreated(t *testing.T) { //done
	domainName = GenerateTestDomain(t)
	resource_name := "testacc-terraform-site_deprecated_test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIncapsulaSiteConfigDeprecated(resource_name, domainName, true),
				ExpectError: regexp.MustCompile("cannot create deprecated resource"),
			},
		},
	})
}

func TestAccIncapsulaSite_ChangeDeprecatedFlag(t *testing.T) { //done
	domainName = GenerateTestDomain(t)
	resource_name := "testacc-terraform-site_deprecated_test"
	full_resource_name := "incapsula_site." + resource_name
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteNotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfigDeprecated(resource_name, domainName, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(full_resource_name),
					resource.TestCheckResourceAttr(full_resource_name, "domain", domainName),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteConfigDeprecated(resource_name, domainName, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(full_resource_name),
					resource.TestCheckResourceAttr(full_resource_name, "domain", domainName),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteConfigDeprecated(resource_name, domainName, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(full_resource_name),
					resource.TestCheckResourceAttr(full_resource_name, "domain", domainName),
				),
				ExpectError: regexp.MustCompile("deprecated flag cannot be changed from true to false"),
			},
		},
	})
}

func TestAccIncapsulaSite_DeprecationFlagChangeAttributes(t *testing.T) {
	domainName = GenerateTestDomain(t)
	resource_name := "testacc-terraform-site_deprecated_test"
	full_resource_name := "incapsula_site." + resource_name
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteNotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaSiteConfigAllParams(resource_name, domainName, false, "aggressive", "active",
					true, "dns", false, false,
					"ref123", "api.seal_location.bottom_right", true, false, false,
					"EU", true, "salt123", "full",
					false, false, false, true,
					false, "dont_include_html", "smart", 7200,
					false, true, 1200,
					false, false,
					"custom", []string{"Content-Type"},
					false, "custom", 600,
					"X-Cache-Tag-1", false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(full_resource_name),
					resource.TestCheckResourceAttr(full_resource_name, "domain", domainName),
					resource.TestCheckResourceAttr(full_resource_name, "acceleration_level", "aggressive"),
					resource.TestCheckResourceAttr(full_resource_name, "active", "active"),
					resource.TestCheckResourceAttr(full_resource_name, "domain_redirect_to_full", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "domain_validation", "dns"),
					resource.TestCheckResourceAttr(full_resource_name, "ignore_ssl", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "remove_ssl", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "ref_id", "ref123"),
					resource.TestCheckResourceAttr(full_resource_name, "seal_location", "api.seal_location.bottom_right"),
					resource.TestCheckResourceAttr(full_resource_name, "restricted_cname_reuse", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "naked_domain_san", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "wildcard_san", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "data_storage_region", "EU"),
					resource.TestCheckResourceAttr(full_resource_name, "hashing_enabled", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "hash_salt", "salt123"),
					resource.TestCheckResourceAttr(full_resource_name, "log_level", "full"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_comply_no_cache", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_enable_client_side_caching", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_send_age_header", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_key_comply_vary", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_key_unite_naked_full_cache", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_https", "dont_include_html"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_level", "smart"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_time", "7200"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_300x", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_404_enabled", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_404_time", "1200"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_empty_responses", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_http_10_responses", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_response_header_mode", "custom"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_response_headers.0", "Content-Type"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_shield", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_stale_content_mode", "custom"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_stale_content_time", "600"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_tag_response_header", "X-Cache-Tag-1"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_ttl_prefer_last_modified", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_ttl_use_shortest_caching", "false"),
				),
			},
			{
				Config: testAccCheckIncapsulaSiteConfigAllParams(resource_name, domainName, true, "none", "bypass",
					false, "email", true, true,
					"ref456", "api.seal_location.bottom_left", false, true, true,
					"APAC", false, "salt456", "none",
					true, true, true, false,
					true, "disabled", "disable", 3600,
					true, true, 600,
					true, true,
					"all", []string{"Authorization"},
					true, "disabled", 300,
					"X-Cache-Tag-2", true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(full_resource_name),
					resource.TestCheckResourceAttr(full_resource_name, "domain", domainName),
					resource.TestCheckResourceAttr(full_resource_name, "acceleration_level", "aggressive"),
					resource.TestCheckResourceAttr(full_resource_name, "active", "active"),
					resource.TestCheckResourceAttr(full_resource_name, "domain_redirect_to_full", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "domain_validation", "dns"),
					resource.TestCheckResourceAttr(full_resource_name, "ignore_ssl", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "remove_ssl", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "ref_id", "ref123"),
					resource.TestCheckResourceAttr(full_resource_name, "seal_location", "api.seal_location.bottom_right"),
					resource.TestCheckResourceAttr(full_resource_name, "restricted_cname_reuse", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "naked_domain_san", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "wildcard_san", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "data_storage_region", "EU"),
					resource.TestCheckResourceAttr(full_resource_name, "hashing_enabled", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "hash_salt", "salt123"),
					resource.TestCheckResourceAttr(full_resource_name, "log_level", "full"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_comply_no_cache", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_enable_client_side_caching", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_client_send_age_header", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_key_comply_vary", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_key_unite_naked_full_cache", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_https", "dont_include_html"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_level", "smart"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_mode_time", "7200"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_300x", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_404_enabled", "true"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_404_time", "1200"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_empty_responses", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_http_10_responses", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_response_header_mode", "custom"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_response_headers.0", "Content-Type"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_cache_shield", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_stale_content_mode", "custom"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_stale_content_time", "600"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_response_tag_response_header", "X-Cache-Tag-1"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_ttl_prefer_last_modified", "false"),
					resource.TestCheckResourceAttr(full_resource_name, "perf_ttl_use_shortest_caching", "false"),
				),
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

func testAccCheckIncapsulaSiteConfigDeprecated(resource_name string, domain string, deprecated bool) string {
	return fmt.Sprintf(`
		resource "incapsula_site" "%s" {
			domain = "%s"
			deprecated = %t
		}`,
		resource_name,
		domain,
		deprecated,
	)
}

func testAccCheckIncapsulaSiteConfigAllParams(resource_name string, domain string, deprecated bool, accelerationLevel string, active string,
	domainRedirectToFull bool, domainValidation string, ignoreSsl bool, removeSsl bool, refId string, sealLocation string,
	restrictedCnameReuse bool, nakedDomainSan bool, wildcardSan bool, dataStorageRegion string, hashingEnabled bool,
	hashSalt string, logLevel string, perfClientComplyNoCache bool, perfClientEnableClientSideCaching bool,
	perfClientSendAgeHeader bool, perfKeyComplyVary bool, perfKeyUniteNakedFullCache bool, perfModeHttps string,
	perfModeLevel string, perfModeTime int, perfResponseCache300x bool, perfResponseCache404Enabled bool,
	perfResponseCache404Time int, perfResponseCacheEmptyResponses bool, perfResponseCacheHttp10Responses bool,
	perfResponseCacheResponseHeaderMode string, perfResponseCacheResponseHeaders []string, perfResponseCacheShield bool,
	perfResponseStaleContentMode string, perfResponseStaleContentTime int, perfResponseTagResponseHeader string,
	perfTtlPreferLastModified bool, perfTtlUseShortestCaching bool,
) string {
	return fmt.Sprintf(`
resource "incapsula_site"         "%s" {
  domain                          = "%s"
  deprecated                      = %t
  acceleration_level              = "%s"
  active                          = "%s"
  domain_redirect_to_full         = %t
  domain_validation               = "%s"
  ignore_ssl                      = %t
  remove_ssl                      = %t
  ref_id                          = "%s"
  seal_location                   = "%s"
  restricted_cname_reuse          = %t
  naked_domain_san                = %t
  wildcard_san                    = %t
  data_storage_region             = "%s"
  hashing_enabled                 = %t
  hash_salt                       = "%s"
  log_level                       = "%s"
  perf_client_comply_no_cache     = %t
  perf_client_enable_client_side_caching = %t
  perf_client_send_age_header     = %t
  perf_key_comply_vary            = %t
  perf_key_unite_naked_full_cache = %t
  perf_mode_https                 = "%s"
  perf_mode_level                 = "%s"
  perf_mode_time                  = %d
  perf_response_cache_300x        = %t
  perf_response_cache_404_enabled = %t
  perf_response_cache_404_time    = %d
  perf_response_cache_empty_responses = %t
  perf_response_cache_http_10_responses = %t
  perf_response_cache_response_header_mode = "%s"
  perf_response_cache_response_headers = ["%s"]
  perf_response_cache_shield      = %t
  perf_response_stale_content_mode = "%s"
  perf_response_stale_content_time = %d
  perf_response_tag_response_header = "%s"
  perf_ttl_prefer_last_modified   = %t
  perf_ttl_use_shortest_caching   = %t
}`,
		resource_name, domain, deprecated, accelerationLevel, active, domainRedirectToFull, domainValidation, ignoreSsl, removeSsl, refId,
		sealLocation, restrictedCnameReuse, nakedDomainSan, wildcardSan, dataStorageRegion, hashingEnabled, hashSalt, logLevel,
		perfClientComplyNoCache, perfClientEnableClientSideCaching, perfClientSendAgeHeader, perfKeyComplyVary, perfKeyUniteNakedFullCache, perfModeHttps,
		perfModeLevel, perfModeTime, perfResponseCache300x, perfResponseCache404Enabled, perfResponseCache404Time, perfResponseCacheEmptyResponses,
		perfResponseCacheHttp10Responses, perfResponseCacheResponseHeaderMode, strings.Join(perfResponseCacheResponseHeaders, "\", \""),
		perfResponseCacheShield, perfResponseStaleContentMode, perfResponseStaleContentTime, perfResponseTagResponseHeader, perfTtlPreferLastModified,
		perfTtlUseShortestCaching,
	)
}

func testCheckIncapsulaSiteNotDestroy(state *terraform.State) error {
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

		domain := res.Primary.Attributes["domain"]
		siteStatusResponse, err := client.SiteStatus(domain, siteID)

		if err != nil {
			return fmt.Errorf("Site status error for %s: %s", siteIDStr, err)
		}

		if siteStatusResponse == nil || siteStatusResponse.SiteID != siteID {
			response, _ := json.Marshal(siteStatusResponse)
			return fmt.Errorf("Incapsula site status for domain: %s (site id: %d) was not retreived. response: %s", domain, siteID, response)
		}
		err = client.DeleteSite(domain, siteID) // clean the env after checking

		if err != nil {
			return fmt.Errorf("Error deleting site (%d) for domain %s: %s", siteID, domain, err)
		}
	}

	return nil
}
