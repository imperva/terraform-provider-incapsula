package incapsula

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const aiApplicationSecurityApplicationApiResourceName = "incapsula_ai_application_security_application.api_test"
const aiApplicationSecurityApplicationEdgeResourceName = "incapsula_ai_application_security_application.edge_test"
const aiApplicationSecurityApplicationEdgePartialResourceName = "incapsula_ai_application_security_application.edge_partial"

// TestAccIncapsulaAiApplicationSecurityApplicationBasic exercises the full API-type resource
// lifecycle against the mock server: create -> read -> update (name + region via
// PATCH) -> import (ImportStateVerify) -> destroy.
// TestAiApplicationSecurityApplicationReadDefaultsContentTypeWhenEmpty is a unit-level regression test
// for the content_type fallback in Read. The backend applies no server-side default for contentType
// and tags it omitempty, so an EDGE app can read back with it absent. Without the fallback, the
// schema Default ("application/json") would re-apply on the next plan and cause a perpetual in-place
// update; Read must resolve the empty value to the default so state matches config.
func TestAiApplicationSecurityApplicationReadDefaultsContentTypeWhenEmpty(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	const appID = "app-123"

	// EDGE app whose configuration omits contentType.
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":[{"applicationId":"app-123","name":"edge-app","accountId":55,"region":"US","applicationType":"EDGE","configuration":{"siteId":100,"path":"/api"}}]}`))
	}))
	defer server.Close()

	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	d := resourceAiApplicationSecurityApplication().TestResourceData()
	d.SetId(appID)
	d.Set("account_id", 55)

	if diags := resourceAiApplicationSecurityApplicationRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("Read returned error: %+v", diags)
	}

	if got := d.Get("configuration.0.content_type").(string); got != "application/json" {
		t.Errorf("content_type not defaulted: got %q, want application/json when the backend omits contentType", got)
	}
}

func TestAccIncapsulaAiApplicationSecurityApplicationBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityApplicationConfig("my-app", "US"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApplicationExists(aiApplicationSecurityApplicationApiResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "name", "my-app"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "application_type", "API"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "region", "US"),
					resource.TestCheckResourceAttrSet(aiApplicationSecurityApplicationApiResourceName, "id"),
					// API applications carry no configuration block.
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "configuration.#", "0"),
				),
			},
			{
				Config: testAccAiApplicationSecurityApplicationConfig("renamed-app", "EU"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApplicationExists(aiApplicationSecurityApplicationApiResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "name", "renamed-app"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationApiResourceName, "region", "EU"),
				),
			},
			{
				ResourceName:      aiApplicationSecurityApplicationApiResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: aiApplicationSecurityApplicationImportID(aiApplicationSecurityApplicationApiResourceName),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityApplicationEdge exercises the EDGE-type resource, whose
// nested configuration block (site_id/path/prompt_location/is_streaming plus the
// nested request{} and response{} blocks) drives expandAiApplicationSecurityApplicationConfig
// and flattenAiApplicationSecurityApplicationConfig. The ImportStateVerify step compares every
// attribute after import against live state, so it is the strongest guard against a
// flatten mismatch in the nested blocks.
func TestAccIncapsulaAiApplicationSecurityApplicationEdge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityApplicationEdgeConfig(t, "edge-app", "/v1/chat/completions"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApplicationExists(aiApplicationSecurityApplicationEdgeResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "name", "edge-app"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "application_type", "EDGE"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.#", "1"),
					// site_id is created inline (incapsula_site) and referenced, so compare against it.
					resource.TestCheckResourceAttrPair(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.site_id", siteResourceName, "id"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.path", "/v1/chat/completions"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.is_streaming", "true"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.request.#", "1"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.request.0.message_path", "$.messages"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.response.#", "1"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.response.0.end_of_stream_marker", "[DONE]"),
				),
			},
			{
				// Mutate a nested configuration field -> PATCH carrying the whole config block.
				Config: testAccAiApplicationSecurityApplicationEdgeConfig(t, "edge-app", "/v2/chat/completions"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApplicationExists(aiApplicationSecurityApplicationEdgeResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgeResourceName, "configuration.0.path", "/v2/chat/completions"),
				),
			},
			{
				ResourceName:      aiApplicationSecurityApplicationEdgeResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: aiApplicationSecurityApplicationImportID(aiApplicationSecurityApplicationEdgeResourceName),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityApplicationEdgePartialConfig covers a configuration block that
// declares a request{} sub-block but NO response{}. flattenAiApplicationSecurityApplicationConfig only
// emits request/response when the backend returns them, so this guards against a phantom
// empty response block (configuration.0.response.# must be 0) and, via ImportStateVerify,
// against any expand/flatten asymmetry when one nested sub-block is absent.
func TestAccIncapsulaAiApplicationSecurityApplicationEdgePartialConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityApplicationEdgePartialConfig(t, "edge-partial"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApplicationExists(aiApplicationSecurityApplicationEdgePartialResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgePartialResourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgePartialResourceName, "configuration.0.request.#", "1"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgePartialResourceName, "configuration.0.request.0.message_path", "$.messages"),
					// No response{} was declared — flatten must not synthesize a phantom one.
					resource.TestCheckResourceAttr(aiApplicationSecurityApplicationEdgePartialResourceName, "configuration.0.response.#", "0"),
				),
			},
			{
				ResourceName:      aiApplicationSecurityApplicationEdgePartialResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: aiApplicationSecurityApplicationImportID(aiApplicationSecurityApplicationEdgePartialResourceName),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityApplicationConfigValidation asserts the type<->configuration
// rules enforced by resourceAiApplicationSecurityApplicationCustomizeDiff: a configuration block is
// required for EDGE and rejected for SDK/API. Both errors surface at plan time, so no
// apply/backend call is needed.
func TestAccIncapsulaAiApplicationSecurityApplicationConfigValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// EDGE without a configuration block must be rejected.
				Config:      testAccAiApplicationSecurityApplicationEdgeMissingConfig("edge-missing-config"),
				ExpectError: regexp.MustCompile("configuration is required for EDGE application type"),
			},
			{
				// SDK with a configuration block must be rejected.
				Config:      testAccAiApplicationSecurityApplicationSdkWithConfig("sdk-with-config"),
				ExpectError: regexp.MustCompile("configuration is not supported for SDK application type"),
			},
		},
	})
}

func testAccAiApplicationSecurityApplicationConfig(name, region string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "api_test" {
  name             = "%s"
  application_type = "API"
  region           = "%s"
}
`, name, region)
}

func testAccAiApplicationSecurityApplicationEdgeConfig(t *testing.T, name, path string) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "edge_test" {
  name             = "%s"
  application_type = "EDGE"
  region           = "US"

  configuration {
    site_id                    = incapsula_site.testacc-terraform-site.id
    path                       = "%s"
    content_type               = "application/json"
    prompt_location            = "$.messages[-1].content"
    blocked_response_structure = "{\"error\": \"$BLOCKED_MESSAGE\"}"
    is_streaming               = true

    request {
      message_path = "$.messages"
      content_path = "$.content"
      role_path    = "$.role"
    }

    response {
      role_path            = "$.choices[0].delta.role"
      content_path         = "$.choices[0].delta.content"
      finish_reason_path   = "$.choices[0].finish_reason"
      finish_reason_value  = "stop"
      end_of_stream_marker = "[DONE]"
    }
  }
}
`, name, path)
}

// testAccAiApplicationSecurityApplicationEdgePartialConfig builds an EDGE configuration that contains a
// request{} block but no response{} block.
func testAccAiApplicationSecurityApplicationEdgePartialConfig(t *testing.T, name string) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "edge_partial" {
  name             = "%s"
  application_type = "EDGE"
  region           = "US"

  configuration {
    site_id                    = incapsula_site.testacc-terraform-site.id
    path                       = "/v1/chat/completions"
    prompt_location            = "$.messages[-1].content"
    blocked_response_structure = "{\"error\": \"$BLOCKED_MESSAGE\"}"

    request {
      message_path = "$.messages"
      content_path = "$.content"
      role_path    = "$.role"
    }
  }
}
`, name)
}

func testAccAiApplicationSecurityApplicationEdgeMissingConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "edge_test" {
  name             = "%s"
  application_type = "EDGE"
  region           = "US"
}
`, name)
}

func testAccAiApplicationSecurityApplicationSdkWithConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "sdk_test" {
  name             = "%s"
  application_type = "SDK"
  region           = "US"

  configuration {
    path = "/v1/chat/completions"
  }
}
`, name)
}

// aiApplicationSecurityApplicationImportID resolves the import key (application UUID) for the
// named resource from Terraform state.
func aiApplicationSecurityApplicationImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckAiApplicationSecurityApplicationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No application ID is set")
		}

		client := testAccProvider.Meta().(*Client)

		accountID, _ := strconv.Atoi(rs.Primary.Attributes["account_id"])
		app, err := client.GetAiApplicationSecurityApplication(accountID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting AI Application Security application: %s", err)
		}
		if app == nil {
			return fmt.Errorf("AI Application Security application %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAiApplicationSecurityApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_ai_application_security_application" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		accountID, _ := strconv.Atoi(rs.Primary.Attributes["account_id"])
		app, err := client.GetAiApplicationSecurityApplication(accountID, rs.Primary.ID)
		if err == nil && app != nil {
			return fmt.Errorf("AI Application Security application still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
