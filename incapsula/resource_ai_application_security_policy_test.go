package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const aiApplicationSecurityPolicyResourceName = "incapsula_ai_application_security_policy.test"

// TestAccIncapsulaAiApplicationSecurityPolicyBasic exercises the full policy lifecycle against the
// mock server: create -> read -> update (rename + toggle active + mutate guardrails via
// PATCH) -> import (ImportStateVerify) -> destroy.
func TestAccIncapsulaAiApplicationSecurityPolicyBasic(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyConfig("my-policy", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "name", "my-policy"),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "active", "true"),
					resource.TestCheckResourceAttrSet(aiApplicationSecurityPolicyResourceName, "id"),
					resource.TestCheckResourceAttrSet(aiApplicationSecurityPolicyResourceName, "application_id"),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "2"),
				),
			},
			{
				Config: testAccAiApplicationSecurityPolicyConfig("renamed-policy", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "name", "renamed-policy"),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "active", "false"),
				),
			},
			{
				ResourceName:      aiApplicationSecurityPolicyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: aiApplicationSecurityPolicyImportID(aiApplicationSecurityPolicyResourceName),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityPolicyInvalidPhase asserts the plan-time type/phase validation
// in resourceAiApplicationSecurityPolicyCustomizeDiff: PROMPT_INJECTION is only valid in the PROMPT
// phase, so declaring it in RESPONSE must fail before any backend call.
func TestAccIncapsulaAiApplicationSecurityPolicyInvalidPhase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAiApplicationSecurityPolicyInvalidPhaseConfig("bad-policy"),
				ExpectError: regexp.MustCompile(`guardrail type "PROMPT_INJECTION" is not valid in phase "RESPONSE"`),
			},
		},
	})
}

// TestAiApplicationSecurityGuardrailHashIgnoresEmbeddedType is a unit-level regression test for the
// guardrail set hash. The backend requires a "type" discriminator inside the config JSON on write and
// always returns it on read; flattenAiApplicationSecurityGuardrail strips it so state matches config.
// A user who redundantly embeds "type" in their config JSON must therefore hash identically to the
// stripped read-back form — otherwise the set element shows a perpetual remove/add. This is enforced
// by normalizeAiApplicationSecurityGuardrailConfig dropping "type".
func TestAiApplicationSecurityGuardrailHashIgnoresEmbeddedType(t *testing.T) {
	withEmbeddedType := map[string]interface{}{
		"type":   "RATE_LIMIT",
		"phase":  "PROMPT",
		"mode":   "BLOCK",
		"active": true,
		"config": `{"type":"RATE_LIMIT","threshold":5}`,
	}
	strippedType := map[string]interface{}{
		"type":   "RATE_LIMIT",
		"phase":  "PROMPT",
		"mode":   "BLOCK",
		"active": true,
		"config": `{"threshold":5}`,
	}

	if h1, h2 := aiApplicationSecurityGuardrailHash(withEmbeddedType), aiApplicationSecurityGuardrailHash(strippedType); h1 != h2 {
		t.Errorf("guardrail hash differs with embedded type (%d) vs stripped (%d); must be equal to avoid a perpetual diff", h1, h2)
	}

	// The normalizer must drop the embedded discriminator directly.
	if got := normalizeAiApplicationSecurityGuardrailConfig(`{"type":"RATE_LIMIT","threshold":5}`); got != `{"threshold":5}` {
		t.Errorf("normalize did not strip embedded type: got %q, want {\"threshold\":5}", got)
	}
}

// TestAccIncapsulaAiApplicationSecurityPolicyConfigJSONIdempotent guards against the TypeSet + JSON
// config hashing trap: a guardrail config written with keys in non-sorted order must not
// produce a perpetual diff. The backend (and flatten) round-trip config with sorted keys,
// so without the guardrail set's normalized Set hash (aiApplicationSecurityGuardrailHash) the refreshed
// set element would hash differently than the configured one and every plan would show a
// spurious guardrail remove/add. The framework's automatic post-apply plan check fails if
// that regresses.
func TestAccIncapsulaAiApplicationSecurityPolicyConfigJSONIdempotent(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyConfigJSONConfig("json-policy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "1"),
				),
			},
			{
				// Re-applying the identical, non-sorted config must be a no-op.
				Config:             testAccAiApplicationSecurityPolicyConfigJSONConfig("json-policy"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityPolicyGuardrailMutation is the positive control for the guardrail
// set hashing / DiffSuppressFunc: it mutates a single guardrail's config value, mode and
// active flag across a second apply and asserts each change actually lands in state. If the
// normalized-config Set hash (aiApplicationSecurityGuardrailHash) or suppressor ever over-suppressed, a
// real edit would be silently dropped and these post-apply checks would fail.
func TestAccIncapsulaAiApplicationSecurityPolicyGuardrailMutation(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyGuardrailMutationConfig("mut-policy", "BLOCK", true, `{"threshold":0.5}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"type":   "PII_STATIC",
						"phase":  "PROMPT",
						"mode":   "BLOCK",
						"active": "true",
						"config": `{"threshold":0.5}`,
					}),
				),
			},
			{
				// Change config value + mode + active in one apply. Every field must round-trip.
				Config: testAccAiApplicationSecurityPolicyGuardrailMutationConfig("mut-policy", "ALERT", false, `{"threshold":0.9}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"type":   "PII_STATIC",
						"phase":  "PROMPT",
						"mode":   "ALERT",
						"active": "false",
						"config": `{"threshold":0.9}`,
					}),
				),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityPolicyGuardrailSetSize exercises guardrail set cardinality changes:
// grow the set 1 -> 2 and shrink it back 2 -> 1, so expandAiApplicationSecurityGuardrails /
// flattenAiApplicationSecurityGuardrails round-trip as guardrails are added and removed (and the request/
// response phase-split buckets grow and shrink) rather than only being mutated in place.
func TestAccIncapsulaAiApplicationSecurityPolicyGuardrailSetSize(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyOneGuardrailConfig("size-policy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "1"),
				),
			},
			{
				// Add a second guardrail (RESPONSE phase) -> set grows to 2.
				Config: testAccAiApplicationSecurityPolicyConfig("size-policy", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "2"),
				),
			},
			{
				// Remove it again -> set shrinks back to 1.
				Config: testAccAiApplicationSecurityPolicyOneGuardrailConfig("size-policy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "1"),
				),
			},
		},
	})
}

// TestAccIncapsulaAiApplicationSecurityPolicyGuardrailConfigDistinct declares two guardrails that are
// identical except for their config JSON (same type/phase/mode/active), one carrying nested
// objects/arrays with non-sorted keys. It asserts (a) the Set hash keys off the normalized
// config so the two remain distinct elements (guardrail.# == 2 rather than colliding to 1),
// and (b) re-applying the deeply-nested, non-sorted config is a no-op.
func TestAccIncapsulaAiApplicationSecurityPolicyGuardrailConfigDistinct(t *testing.T) {
	skipAiApplicationSecurityLiveAccTest(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig("distinct-policy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", "2"),
					// The nested config round-trips into a distinct set element intact. State
					// keeps the configured (raw) string — the diff against the backend's sorted
					// form is absorbed by DiffSuppressFunc, while aiApplicationSecurityGuardrailHash normalizes
					// before hashing so the two guardrails stay distinct and stable.
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"config": `{"nested":{"b":2,"a":1},"arr":[3,1,2]}`,
					}),
				),
			},
			{
				// Re-applying the identical, non-sorted, nested config must be a no-op.
				Config:             testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig("distinct-policy"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccAiApplicationSecurityPolicyConfig(name string, active bool) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"
  active         = %t

  guardrail {
    type  = "PROMPT_INJECTION"
    phase = "PROMPT"
    mode  = "BLOCK"
  }

  guardrail {
    type  = "MODERATION"
    phase = "RESPONSE"
    mode  = "ALERT"
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name, active)
}

// testAccAiApplicationSecurityPolicyConfigJSONConfig declares a guardrail whose config JSON has multiple
// keys in deliberately non-sorted order, to exercise config normalization/hash stability.
func testAccAiApplicationSecurityPolicyConfigJSONConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type   = "PROMPT_INJECTION"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = "{\"z_threshold\":0.9,\"a_enabled\":true}"
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name)
}

// testAccAiApplicationSecurityPolicyGuardrailMutationConfig declares a single PII_STATIC/PROMPT guardrail
// whose mode, active flag and config JSON are caller-controlled, so a step can mutate exactly
// one guardrail in place and assert the change round-trips.
func testAccAiApplicationSecurityPolicyGuardrailMutationConfig(name, mode string, active bool, config string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type   = "PII_STATIC"
    phase  = "PROMPT"
    mode   = "%s"
    active = %t
    config = %q
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name, mode, active, config)
}

// testAccAiApplicationSecurityPolicyOneGuardrailConfig declares a policy with a single guardrail, used to
// grow/shrink the guardrail set against the two-guardrail testAccAiApplicationSecurityPolicyConfig.
func testAccAiApplicationSecurityPolicyOneGuardrailConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type  = "PROMPT_INJECTION"
    phase = "PROMPT"
    mode  = "BLOCK"
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name)
}

// testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig declares two guardrails that differ ONLY
// by their config JSON (same type/phase/mode), one of which carries nested objects/arrays with
// deliberately non-sorted keys — exercising both config-based set distinctness and deep JSON
// normalization stability.
func testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type   = "PII_STATIC"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = "{\"z_threshold\":0.9,\"a_enabled\":true}"
  }

  guardrail {
    type   = "PII_STATIC"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = "{\"nested\":{\"b\":2,\"a\":1},\"arr\":[3,1,2]}"
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name)
}

func testAccAiApplicationSecurityPolicyInvalidPhaseConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  account_id       = %d
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  account_id     = %d
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type  = "PROMPT_INJECTION"
    phase = "RESPONSE"
    mode  = "BLOCK"
  }
}
`, aiApplicationSecurityTestAccountID(), aiApplicationSecurityTestAccountID(), name)
}

// aiApplicationSecurityPolicyImportID resolves the import key (policy UUID) for the named resource
// from Terraform state.
func aiApplicationSecurityPolicyImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckAiApplicationSecurityPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No policy ID is set")
		}

		client := testAccProvider.Meta().(*Client)

		policy, err := client.GetAiApplicationSecurityPolicy(aiApplicationSecurityTestAccountID(), rs.Primary.Attributes["application_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting AI Application Security policy: %s", err)
		}
		if policy == nil {
			return fmt.Errorf("AI Application Security policy %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAiApplicationSecurityPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_ai_application_security_policy" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		policy, err := client.GetAiApplicationSecurityPolicy(aiApplicationSecurityTestAccountID(), rs.Primary.Attributes["application_id"], rs.Primary.ID)
		if err == nil && policy != nil {
			return fmt.Errorf("AI Application Security policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
