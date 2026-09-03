package incapsula

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const aiApplicationSecurityApiKeyResourceName = "incapsula_ai_application_security_api_key.test"

// TestAccIncapsulaAiApplicationSecurityApiKeyBasic exercises the API key lifecycle against the mock
// server: create -> read -> import -> destroy. Because the plaintext api_key is only
// returned on create and is unrecoverable afterward, the test asserts it is non-empty
// after create and is (intentionally) ignored on import (empty after import).
func TestAccIncapsulaAiApplicationSecurityApiKeyBasic(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityApiKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityApiKeyConfig("my-key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityApiKeyExists(aiApplicationSecurityApiKeyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityApiKeyResourceName, "name", "my-key"),
					resource.TestCheckResourceAttrSet(aiApplicationSecurityApiKeyResourceName, "id"),
					resource.TestCheckResourceAttrSet(aiApplicationSecurityApiKeyResourceName, "application_id"),
					// The plaintext key must be present after create.
					resource.TestCheckResourceAttrSet(aiApplicationSecurityApiKeyResourceName, "api_key"),
					resource.TestCheckResourceAttr(aiApplicationSecurityApiKeyResourceName, "active", "true"),
				),
			},
			{
				ResourceName:      aiApplicationSecurityApiKeyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// api_key (plaintext) is only returned on create and cannot be recovered on
				// read/import, so it is empty after import and must be excluded from verify.
				ImportStateVerifyIgnore: []string{"api_key"},
				ImportStateIdFunc:       aiApplicationSecurityApiKeyImportID(aiApplicationSecurityApiKeyResourceName),
			},
		},
	})
}

// TestAiApplicationSecurityApiKeyReadPreservesApplicationIDWhenOmitted is a unit-level regression test
// for the ForceNew application_id guard in Read. The account-level list DTO tags applicationId
// omitempty and the backend can legitimately drop it, so Read must NOT overwrite the configured
// UUID with an empty string — otherwise the next plan would schedule a destroy/recreate whose
// Delete targets an empty application path.
func TestAiApplicationSecurityApiKeyReadPreservesApplicationIDWhenOmitted(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	const configuredAppID = "11111111-1111-1111-1111-111111111111"

	// The list returns the key (id 42) but omits applicationId entirely.
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":{"apiKeys":[{"id":42,"name":"my-key","accountId":55,"active":true}],"totalCount":1}}`))
	}))
	defer server.Close()

	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	d := resourceAiApplicationSecurityApiKey().TestResourceData()
	d.SetId("42")
	d.Set("account_id", 55)
	d.Set("application_id", configuredAppID)

	if diags := resourceAiApplicationSecurityApiKeyRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("Read returned error: %+v", diags)
	}

	if got := d.Get("application_id").(string); got != configuredAppID {
		t.Errorf("application_id was overwritten to %q; the configured UUID %q must be preserved when the list omits applicationId", got, configuredAppID)
	}
	// Sanity: the rest of the response is still applied.
	if got := d.Get("name").(string); got != "my-key" {
		t.Errorf("name not set from read: got %q, want my-key", got)
	}
	if got := d.Get("active").(bool); !got {
		t.Errorf("active not set from read: got %v, want true", got)
	}
}

func testAccAiApplicationSecurityApiKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "api_key_app" {
  account_id       = %d
  name             = "api-key-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_api_key" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.api_key_app.id
  name           = "%s"
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name)
}

// aiApplicationSecurityApiKeyImportID resolves the import key (numeric API key ID) for the named
// resource from Terraform state.
func aiApplicationSecurityApiKeyImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckAiApplicationSecurityApiKeyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No API key ID is set")
		}

		apiKeyID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid API key ID %q: %s", rs.Primary.ID, err)
		}

		client := testAccProvider.Meta().(*Client)

		apiKey, err := client.GetAiApplicationSecurityApiKey(aiApplicationSecurityTestAccountID(), apiKeyID)
		if err != nil {
			return fmt.Errorf("Error getting AI Application Security API key: %s", err)
		}
		if apiKey == nil {
			return fmt.Errorf("AI Application Security API key %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAiApplicationSecurityApiKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_ai_application_security_api_key" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		apiKeyID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			continue
		}

		apiKey, err := client.GetAiApplicationSecurityApiKey(aiApplicationSecurityTestAccountID(), apiKeyID)
		if err == nil && apiKey != nil {
			return fmt.Errorf("AI Application Security API key still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
