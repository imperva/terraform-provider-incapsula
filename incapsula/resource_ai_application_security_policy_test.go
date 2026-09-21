package incapsula

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const aiApplicationSecurityPolicyResourceName = "incapsula_ai_application_security_policy.test"

// aiApplicationSecurityPolicyMandatoryGuardrailCount is how many guardrails every fixture below
// declares. A policy must carry an EXACT guardian type set per phase, and the required set depends on
// the application's deployment type (ApplicationPolicyValidator.validateExactGuardianSet, AIFW-712,
// driven by deployment-policy-allowance in the service's application.yaml). The fixtures here all
// build an API-type application, whose set is:
//
//	PROMPT   : PROMPT_INJECTION, PII_STATIC, RATE_LIMIT, ZERO_SHOT_CLASSIFICATION
//	RESPONSE : PII_STATIC, RATE_LIMIT, ZERO_SHOT_CLASSIFICATION
//
// SDK and Edge applications require different sets (SDK additionally allows SYSTEM_PROMPT_LEAK and
// MODERATION in RESPONSE, for nine guardrails), so this count is specific to application_type = "API".
// Declaring fewer is rejected with "Invalid guardian set ... Missing guardians", so a fixture cannot
// use a trimmed-down set. The sets are compared by type, so a fixture may add *extra* guardrails of a
// type already present in that phase.
const aiApplicationSecurityPolicyMandatoryGuardrailCount = 7

// TestAccIncapsulaAiApplicationSecurityPolicyBasic exercises the full policy lifecycle:
// create -> read -> update (rename + toggle active via PATCH, guardrail set unchanged) ->
// import (ImportStateVerify) -> destroy. Guardrail mutation is covered by
// TestAccIncapsulaAiApplicationSecurityPolicyGuardrailMutation.
func TestAccIncapsulaAiApplicationSecurityPolicyBasic(t *testing.T) {
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
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", strconv.Itoa(aiApplicationSecurityPolicyMandatoryGuardrailCount)),
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
//
// Every config in the mandatory-set fixture lists its keys in non-sorted order (see the const block
// below), so the plain fixture is what exercises this; the second step is the assertion.
func TestAccIncapsulaAiApplicationSecurityPolicyConfigJSONIdempotent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyConfig("json-policy", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", strconv.Itoa(aiApplicationSecurityPolicyMandatoryGuardrailCount)),
				),
			},
			{
				// Re-applying the identical, non-sorted config must be a no-op.
				Config:             testAccAiApplicationSecurityPolicyConfig("json-policy", true),
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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyGuardrailMutationConfig("mut-policy", "BLOCK", true, `{"message":"first","enabledPatterns":["email_address"]}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", strconv.Itoa(aiApplicationSecurityPolicyMandatoryGuardrailCount)),
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"type":   "PII_STATIC",
						"phase":  "PROMPT",
						"mode":   "BLOCK",
						"active": "true",
						"config": `{"message":"first","enabledPatterns":["email_address"]}`,
					}),
				),
			},
			{
				// Change config value + mode + active in one apply. Every field must round-trip.
				Config: testAccAiApplicationSecurityPolicyGuardrailMutationConfig("mut-policy", "ALERT", false, `{"message":"second","enabledPatterns":["us_social_security_number"]}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", strconv.Itoa(aiApplicationSecurityPolicyMandatoryGuardrailCount)),
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"type":   "PII_STATIC",
						"phase":  "PROMPT",
						"mode":   "ALERT",
						"active": "false",
						"config": `{"message":"second","enabledPatterns":["us_social_security_number"]}`,
					}),
				),
			},
		},
	})
}

// TestAiApplicationSecurityGuardrailSetSize exercises guardrail set cardinality changes: grow the set
// 1 -> 2 and shrink it back 2 -> 1, so expandAiApplicationSecurityGuardrails /
// flattenAiApplicationSecurityGuardrails round-trip as guardrails are added and removed (and the
// request/response phase-split buckets grow and shrink) rather than only being mutated in place.
//
// This is a unit test rather than an acceptance test on purpose. The backend's PATCH is update-only:
// PolicyService.mutateExistingGuardians matches each incoming guardrail to an existing one by
// (phase, guardianType) and "never inserts or deletes guardians" — an incoming type with no existing
// guardian is silently ignored, and an existing type absent from the request is left untouched. A live
// apply that adds or removes a guardrail therefore returns 200 while changing nothing, so cardinality
// could never round-trip through the API. Driving expand/flatten directly still covers the phase-bucket
// bookkeeping this test guards.
func TestAiApplicationSecurityGuardrailSetSize(t *testing.T) {
	guardrail := func(guardrailType, phase, mode, config string) interface{} {
		return map[string]interface{}{
			"type":   guardrailType,
			"phase":  phase,
			"mode":   mode,
			"active": true,
			"config": config,
		}
	}
	newSet := func(elems ...interface{}) *schema.Set {
		return schema.NewSet(aiApplicationSecurityGuardrailHash, elems)
	}

	promptGuardrail := guardrail("PROMPT_INJECTION", "PROMPT", "BLOCK", `{"threshold":0.8}`)
	responseGuardrail := guardrail("MODERATION", "RESPONSE", "ALERT", "{}")

	// One guardrail: a single PROMPT entry and an empty (but non-nil) RESPONSE bucket.
	request, response, err := expandAiApplicationSecurityGuardrails(newSet(promptGuardrail))
	if err != nil {
		t.Fatalf("expand with 1 guardrail: %s", err)
	}
	if len(request) != 1 || len(response) != 0 {
		t.Fatalf("1 guardrail split into %d request / %d response; want 1 / 0", len(request), len(response))
	}

	// Grow to two: the RESPONSE guardrail must land in the response bucket, not the request one.
	request, response, err = expandAiApplicationSecurityGuardrails(newSet(promptGuardrail, responseGuardrail))
	if err != nil {
		t.Fatalf("expand with 2 guardrails: %s", err)
	}
	if len(request) != 1 || len(response) != 1 {
		t.Fatalf("2 guardrails split into %d request / %d response; want 1 / 1", len(request), len(response))
	}

	// Flatten merges both buckets back into a set of 2 that hashes to the configured elements, so a
	// grown set is stable across a refresh instead of showing a spurious remove/add.
	flattened, err := flattenAiApplicationSecurityGuardrails(&AiApplicationSecurityPolicyResponse{Request: request, Response: response})
	if err != nil {
		t.Fatalf("flatten: %s", err)
	}
	roundTripped := schema.NewSet(aiApplicationSecurityGuardrailHash, flattened)
	if roundTripped.Len() != 2 {
		t.Fatalf("flattened set has %d elements; want 2", roundTripped.Len())
	}
	for _, want := range []interface{}{promptGuardrail, responseGuardrail} {
		if !roundTripped.Contains(want) {
			t.Errorf("flattened set is missing configured guardrail %v", want)
		}
	}

	// Shrink back to one: the response bucket empties again.
	request, response, err = expandAiApplicationSecurityGuardrails(newSet(promptGuardrail))
	if err != nil {
		t.Fatalf("expand after shrink: %s", err)
	}
	if len(request) != 1 || len(response) != 0 {
		t.Fatalf("shrunk set split into %d request / %d response; want 1 / 0", len(request), len(response))
	}
}

// TestAccIncapsulaAiApplicationSecurityPolicyGuardrailConfigDistinct declares two guardrails that are
// identical except for their config JSON (same type/phase/mode/active), one carrying nested
// objects with non-sorted keys. It asserts (a) the Set hash keys off the normalized
// config so the two remain distinct elements (the mandatory set plus one, rather than colliding
// into the mandatory set), and (b) re-applying the deeply-nested, non-sorted config is a no-op.
func TestAccIncapsulaAiApplicationSecurityPolicyGuardrailConfigDistinct(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAiApplicationSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig("distinct-policy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAiApplicationSecurityPolicyExists(aiApplicationSecurityPolicyResourceName),
					resource.TestCheckResourceAttr(aiApplicationSecurityPolicyResourceName, "guardrail.#", strconv.Itoa(aiApplicationSecurityPolicyMandatoryGuardrailCount+1)),
					// The nested config round-trips into a distinct set element intact. State
					// keeps the configured (raw) string — the diff against the backend's sorted
					// form is absorbed by DiffSuppressFunc, while aiApplicationSecurityGuardrailHash normalizes
					// before hashing so the two guardrails stay distinct and stable.
					resource.TestCheckTypeSetElemNestedAttrs(aiApplicationSecurityPolicyResourceName, "guardrail.*", map[string]string{
						"config": aiApplicationSecurityPolicyRateLimitNestedConfig,
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

// Guardrail configs used by the fixtures below. Two properties matter:
//
//  1. Keys are deliberately listed in non-sorted order, so every apply exercises
//     normalizeAiApplicationSecurityGuardrailConfig / aiApplicationSecurityGuardrailHash.
//  2. They spell out the fields the backend populates by default in its response — categories for
//     ZERO_SHOT_CLASSIFICATION, enabledPatterns for PII_STATIC. The backend echoes config back with
//     those keys present, so a fixture that omitted them would read back a config it never wrote and
//     every plan would be non-empty. Keys the backend does not recognise are dropped entirely, which
//     is why these use real schema fields rather than arbitrary ones.
const (
	aiApplicationSecurityPolicyPromptInjectionConfig = `{"threshold":0.8,"message":"blocked"}`
	aiApplicationSecurityPolicyZeroShotConfig        = `{"threshold":0.8,"categories":[]}`
	aiApplicationSecurityPolicyPiiStaticConfig       = `{"message":"masked","enabledPatterns":["email_address"]}`
)

// testAccAiApplicationSecurityPolicyMandatoryConfig renders an API-type application plus a policy
// declaring the mandatory guardian set for that type (see
// aiApplicationSecurityPolicyMandatoryGuardrailCount).
//
// The PROMPT-phase PII_STATIC guardrail's mode, active flag and config are caller-supplied: it is the
// guardrail the tests drive, chosen because PII_STATIC appears once per phase in the mandatory set.
// That matters because the backend's PATCH matches guardrails by (phase, type) and keeps only the
// first of any duplicates, so a duplicated type could not be mutated. extra is appended verbatim for
// tests that need an additional guardrail beyond the mandatory set.
func testAccAiApplicationSecurityPolicyMandatoryConfig(name string, active bool, piiMode string, piiActive bool, piiConfig, extra string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"
  active         = %t

  guardrail {
    type   = "PROMPT_INJECTION"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = %q
  }

  guardrail {
    type  = "RATE_LIMIT"
    phase = "PROMPT"
    mode  = "BLOCK"
  }

  guardrail {
    type   = "ZERO_SHOT_CLASSIFICATION"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = %q
  }

  guardrail {
    type   = "PII_STATIC"
    phase  = "PROMPT"
    mode   = "%s"
    active = %t
    config = %q
  }

  guardrail {
    type   = "PII_STATIC"
    phase  = "RESPONSE"
    mode   = "MASK"
    config = %q
  }

  guardrail {
    type  = "RATE_LIMIT"
    phase = "RESPONSE"
    mode  = "BLOCK"
  }

  guardrail {
    type   = "ZERO_SHOT_CLASSIFICATION"
    phase  = "RESPONSE"
    mode   = "BLOCK"
    config = %q
  }
%s
}
`, name, active,
		aiApplicationSecurityPolicyPromptInjectionConfig,
		aiApplicationSecurityPolicyZeroShotConfig,
		piiMode, piiActive, piiConfig,
		aiApplicationSecurityPolicyPiiStaticConfig,
		aiApplicationSecurityPolicyZeroShotConfig,
		extra)
}

// testAccAiApplicationSecurityPolicyConfig is the plain mandatory set, with every guardrail active and
// PII_STATIC masking.
func testAccAiApplicationSecurityPolicyConfig(name string, active bool) string {
	return testAccAiApplicationSecurityPolicyMandatoryConfig(name, active, "MASK", true, aiApplicationSecurityPolicyPiiStaticConfig, "")
}

// testAccAiApplicationSecurityPolicyGuardrailMutationConfig drives the PROMPT-phase PII_STATIC
// guardrail's mode, active flag and config JSON, so a step can mutate exactly one guardrail in place
// and assert the change round-trips.
func testAccAiApplicationSecurityPolicyGuardrailMutationConfig(name, mode string, active bool, config string) string {
	return testAccAiApplicationSecurityPolicyMandatoryConfig(name, true, mode, active, config, "")
}

// testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig adds a second RATE_LIMIT/PROMPT/BLOCK
// guardrail that differs from the mandatory one ONLY by its config JSON, which carries nested objects
// with deliberately non-sorted keys at both levels — exercising config-based set distinctness and deep
// JSON normalization stability. RATE_LIMIT is used because its schema is the only one with nested
// objects the backend round-trips verbatim.
func testAccAiApplicationSecurityPolicyGuardrailConfigDistinctConfig(name string) string {
	return testAccAiApplicationSecurityPolicyMandatoryConfig(name, true, "MASK", true, aiApplicationSecurityPolicyPiiStaticConfig, fmt.Sprintf(`
  guardrail {
    type   = "RATE_LIMIT"
    phase  = "PROMPT"
    mode   = "BLOCK"
    config = %q
  }
`, aiApplicationSecurityPolicyRateLimitNestedConfig))
}

// aiApplicationSecurityPolicyRateLimitNestedConfig is a RATE_LIMIT config with nested objects whose keys
// are non-sorted at both the top level (promptLimitConfig before globalConfig) and inside each nested
// object, so normalization has to recurse.
const aiApplicationSecurityPolicyRateLimitNestedConfig = `{"promptLimitConfig":{"maxCharacters":512,"enabled":true},"globalConfig":{"timeUnitInSec":3600,"maxTokens":1000,"enabled":true}}`

func testAccAiApplicationSecurityPolicyInvalidPhaseConfig(name string) string {
	return fmt.Sprintf(`
resource "incapsula_ai_application_security_application" "policy_app" {
  name             = "policy-app"
  application_type = "API"
  region           = "US"
}

resource "incapsula_ai_application_security_policy" "test" {
  application_id = incapsula_ai_application_security_application.policy_app.id
  name           = "%s"

  guardrail {
    type  = "PROMPT_INJECTION"
    phase = "RESPONSE"
    mode  = "BLOCK"
  }
}
`, name)
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

		accountID, _ := strconv.Atoi(rs.Primary.Attributes["account_id"])
		policy, err := client.GetAiApplicationSecurityPolicy(accountID, rs.Primary.Attributes["application_id"], rs.Primary.ID)
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

		accountID, _ := strconv.Atoi(rs.Primary.Attributes["account_id"])
		policy, err := client.GetAiApplicationSecurityPolicy(accountID, rs.Primary.Attributes["application_id"], rs.Primary.ID)
		if err == nil && policy != nil {
			return fmt.Errorf("AI Application Security policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
