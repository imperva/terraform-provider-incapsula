package incapsula

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const gexResourceName = "incapsula_ai_application_security_guardrail_exception.test"

// Lifecycle against the mock: create (heredoc pattern - its trailing newline is trimmed by the service and is no
// change) -> import by application/policy/id -> in-place edit of every field, including a move to another guardrail
// -> destroy. The framework plans again after each step and fails on any difference.
func TestAccIncapsulaAiApplicationSecurityGuardrailExceptionBasic(t *testing.T) {
	var createdID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuardrailExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailExceptionConfig(`
  guardrail_type    = "PROMPT_INJECTION"
  phase             = "PROMPT"
  exception_pattern = <<EOT
ignore previous instructions in test prompts
EOT
  exception_reason  = "FALSE_POSITIVE"
`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuardrailExceptionExists(gexResourceName),
					func(s *terraform.State) error {
						createdID = s.RootModule().Resources[gexResourceName].Primary.ID
						return nil
					},
					resource.TestCheckResourceAttr(gexResourceName, "guardrail_type", "PROMPT_INJECTION"),
					resource.TestCheckResourceAttr(gexResourceName, "phase", "PROMPT"),
					resource.TestCheckResourceAttr(gexResourceName, "exception_pattern", "ignore previous instructions in test prompts"),
					resource.TestCheckResourceAttr(gexResourceName, "active", "true"),
					resource.TestCheckResourceAttrSet(gexResourceName, "guardrail_id"),
					resource.TestCheckResourceAttrSet(gexResourceName, "author"),
				),
			},
			{
				ResourceName:      gexResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					a := s.RootModule().Resources[gexResourceName].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s", a["application_id"], a["policy_id"], a["id"]), nil
				},
			},
			{
				Config: testAccGuardrailExceptionConfig(`
  guardrail_type    = "PII_STATIC"
  phase             = "RESPONSE"
  exception_pattern = "a different allowed sentence"
  exception_reason  = "OTHER"
  note              = "  checked by security  "
  active            = false
`),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						if id := s.RootModule().Resources[gexResourceName].Primary.ID; id != createdID {
							return fmt.Errorf("replaced (%s -> %s), expected an in-place update", createdID, id)
						}
						return nil
					},
					resource.TestCheckResourceAttr(gexResourceName, "guardrail_type", "PII_STATIC"),
					resource.TestCheckResourceAttr(gexResourceName, "phase", "RESPONSE"),
					resource.TestCheckResourceAttr(gexResourceName, "note", "checked by security"),
					// the output reads guardrail_id in the same apply - it must already hold the new guardrail's id
					func(s *terraform.State) error {
						want := s.RootModule().Resources[gexResourceName].Primary.Attributes["guardrail_id"]
						if got := fmt.Sprint(s.RootModule().Outputs["gex_guardrail_id"].Value); got != want {
							return fmt.Errorf("output gex_guardrail_id = %s, want the new guardrail_id %s", got, want)
						}
						return nil
					},
					resource.TestCheckResourceAttr(gexResourceName, "active", "false"),
					resource.TestCheckResourceAttrSet(gexResourceName, "updated_by"),
				),
			},
		},
	})
}

// A guardrail the policy doesn't have fails with a clear error, before anything is written.
func TestAccIncapsulaAiApplicationSecurityGuardrailExceptionGuardrailNotInPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuardrailExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailExceptionConfig(`
  guardrail_type    = "MODERATION"
  phase             = "RESPONSE"
  exception_pattern = "this is a valid sentence"
  exception_reason  = "OTHER"
`),
				ExpectError: regexp.MustCompile(`has no MODERATION guardrail in phase RESPONSE`),
			},
		},
	})
}

// RATE_LIMIT has no text to except - rejected at plan time.
func TestAccIncapsulaAiApplicationSecurityGuardrailExceptionRateLimitRejected(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailExceptionConfig(`
  guardrail_type    = "RATE_LIMIT"
  phase             = "PROMPT"
  exception_pattern = "this is a valid sentence"
  exception_reason  = "OTHER"
`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected guardrail_type to be one of`),
			},
		},
	})
}

func TestGuardrailIdOf(t *testing.T) {
	policy := &AiApplicationSecurityPolicyResponse{
		Id:       "p1",
		Request:  []AiApplicationSecurityGuardrail{{Id: 1, GuardrailType: "PROMPT_INJECTION"}, {Id: 2, GuardrailType: "MODERATION"}, {Id: 3, GuardrailType: "MODERATION"}},
		Response: []AiApplicationSecurityGuardrail{{Id: 4, GuardrailType: "PII_STATIC"}},
	}

	if id, err := guardrailIdOf(policy, "PROMPT_INJECTION", "PROMPT"); err != nil || id != 1 {
		t.Errorf("PROMPT_INJECTION/PROMPT: got %d, %v", id, err)
	}
	if id, err := guardrailIdOf(policy, "PII_STATIC", "RESPONSE"); err != nil || id != 4 {
		t.Errorf("PII_STATIC/RESPONSE: got %d, %v", id, err)
	}
	if _, err := guardrailIdOf(policy, "PII_STATIC", "PROMPT"); err == nil || !strings.Contains(err.Error(), "has no PII_STATIC guardrail in phase PROMPT") {
		t.Errorf("only in the other phase: expected not found, got %v", err)
	}
	if _, err := guardrailIdOf(policy, "MODERATION", "PROMPT"); err == nil || !strings.Contains(err.Error(), "has 2 MODERATION guardrails") {
		t.Errorf("two matches: expected an error, got %v", err)
	}
	if gt, phase, ok := guardrailTypeAndPhaseOf(policy, 4); !ok || gt != "PII_STATIC" || phase != "RESPONSE" {
		t.Errorf("id 4: got %s %s %v", gt, phase, ok)
	}
	if _, _, ok := guardrailTypeAndPhaseOf(policy, 99); ok {
		t.Errorf("id 99 is not in the policy")
	}
}

func TestSuppressOuterWhitespaceDiff(t *testing.T) {
	if !suppressOuterWhitespaceDiff("", "allowed text", "allowed text\n", nil) {
		t.Errorf("a trailing newline must not be a change")
	}
	if suppressOuterWhitespaceDiff("", "allowed text", "allowed  text", nil) {
		t.Errorf("inner whitespace is a real change")
	}
}

func TestValidateGuardrailExceptionPatternLength(t *testing.T) {
	atMax := "one two " + strings.Repeat("a", aiApplicationSecurityGuardrailExceptionPatternMaxLength-len("one two "))
	if _, errs := validateGuardrailExceptionPatternLength(atMax, "exception_pattern"); len(errs) != 0 {
		t.Errorf("a pattern at the max length must be accepted, got %v", errs)
	}
	if _, errs := validateGuardrailExceptionPatternLength("  "+atMax+"\n", "exception_pattern"); len(errs) != 0 {
		t.Errorf("outer whitespace is trimmed before counting, got %v", errs)
	}
	if _, errs := validateGuardrailExceptionPatternLength(atMax+"a", "exception_pattern"); len(errs) != 1 {
		t.Errorf("a pattern one over the max length must be rejected")
	}
	// Hebrew letters are 2 bytes each in UTF-8 but 1 character for the service
	hebrew := strings.Repeat("ש", aiApplicationSecurityGuardrailExceptionPatternMaxLength)
	if _, errs := validateGuardrailExceptionPatternLength(hebrew, "exception_pattern"); len(errs) != 0 {
		t.Errorf("characters are counted, not bytes, got %v", errs)
	}
	// an emoji is 2 UTF-16 code units, like Java counts it
	emojis := strings.Repeat("😀", aiApplicationSecurityGuardrailExceptionPatternMaxLength/2+1)
	if _, errs := validateGuardrailExceptionPatternLength(emojis, "exception_pattern"); len(errs) != 1 {
		t.Errorf("an emoji counts as 2 characters, like the service")
	}
}

// Read keeps type and phase when the exception's guardrail is no longer in the policy, instead of blanking the two
// required fields - the next apply attaches the exception to the guardrail of that type and phase again.
func TestGuardrailExceptionReadKeepsTypeAndPhaseWhenGuardrailIsGone(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "/guardian-exceptions/") {
			rw.Write([]byte(`{"data":{"id":"` + testGexID + `","guardianId":99,"exceptionPattern":"some allowed text","exceptionReason":"OTHER","active":true}}`))
			return
		}
		rw.Write([]byte(`{"data":{"id":"` + testGexPolicyID + `","accountId":55,"request":[{"id":1,"guardrailType":"MODERATION","guardrailPhase":"PROMPT"}],"response":[]}}`))
	}))
	defer server.Close()

	d := resourceAiApplicationSecurityGuardrailException().TestResourceData()
	d.SetId(testGexID)
	d.Set("account_id", 55)
	d.Set("application_id", testGexAppID)
	d.Set("policy_id", testGexPolicyID)
	d.Set("guardrail_type", "PROMPT_INJECTION")
	d.Set("phase", "PROMPT")

	if diags := resourceAiApplicationSecurityGuardrailExceptionRead(context.Background(), d, testGexClient(server.URL)); diags.HasError() {
		t.Fatalf("Read returned error: %+v", diags)
	}
	if got := d.Get("guardrail_type").(string); got != "PROMPT_INJECTION" {
		t.Errorf("guardrail_type must keep its value, got %q", got)
	}
	if got := d.Get("phase").(string); got != "PROMPT" {
		t.Errorf("phase must keep its value, got %q", got)
	}
	if got := d.Get("guardrail_id").(int); got != 99 {
		t.Errorf("guardrail_id must show what the service returned, got %d", got)
	}
}

// updated_by is planned as "known after apply" when the update changes something, and not at all when nothing
// changes (otherwise every plan would show a change). Outer whitespace is no change: the service trims it.
func TestGuardrailExceptionPlanMarksUpdatedByOnlyOnUpdate(t *testing.T) {
	resourceSchema := resourceAiApplicationSecurityGuardrailException()
	state := &terraform.InstanceState{ID: testGexID, Attributes: map[string]string{
		"id": testGexID, "account_id": "55", "application_id": testGexAppID, "policy_id": testGexPolicyID,
		"guardrail_type": "PROMPT_INJECTION", "phase": "PROMPT", "exception_pattern": "some allowed text",
		"exception_reason": "OTHER", "note": "a note", "active": "true", "guardrail_id": "1", "author": "alice", "updated_by": "alice",
	}}
	plan := func(overrides map[string]interface{}) *terraform.InstanceDiff {
		config := map[string]interface{}{
			"account_id": 55, "application_id": testGexAppID, "policy_id": testGexPolicyID,
			"guardrail_type": "PROMPT_INJECTION", "phase": "PROMPT", "exception_pattern": "some allowed text",
			"exception_reason": "OTHER", "note": "a note", "active": true,
		}
		for k, v := range overrides {
			config[k] = v
		}
		diff, err := resourceSchema.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatalf("Diff returned error: %v", err)
		}
		return diff
	}
	updatedByUnknown := func(diff *terraform.InstanceDiff) bool {
		return diff != nil && diff.Attributes["updated_by"] != nil && diff.Attributes["updated_by"].NewComputed
	}

	if diff := plan(nil); diff != nil && !diff.Empty() {
		t.Errorf("no change must give an empty plan, got %+v", diff.Attributes)
	}
	if diff := plan(map[string]interface{}{"exception_pattern": "some allowed text\n", "note": "  a note "}); diff != nil && !diff.Empty() {
		t.Errorf("outer whitespace must give an empty plan, got %+v", diff.Attributes)
	}
	for name, overrides := range map[string]map[string]interface{}{
		"note":              {"note": "another note"},
		"exception_reason":  {"exception_reason": "FALSE_POSITIVE"},
		"active":            {"active": false},
		"exception_pattern": {"exception_pattern": "other allowed text"},
	} {
		if diff := plan(overrides); !updatedByUnknown(diff) {
			t.Errorf("changing %s must plan updated_by as known after apply", name)
		}
	}
	if diff := plan(map[string]interface{}{"phase": "RESPONSE"}); !updatedByUnknown(diff) || diff.Attributes["guardrail_id"] == nil || !diff.Attributes["guardrail_id"].NewComputed {
		t.Errorf("moving the guardrail must plan updated_by and guardrail_id as known after apply")
	}
}

// A trailing slash on the import ID is ignored.
func TestGuardrailExceptionImportIgnoresTrailingSlash(t *testing.T) {
	for _, id := range []string{
		testGexAppID + "/" + testGexPolicyID + "/" + testGexID,
		testGexAppID + "/" + testGexPolicyID + "/" + testGexID + "/",
		"55/" + testGexAppID + "/" + testGexPolicyID + "/" + testGexID + "/",
	} {
		d := resourceAiApplicationSecurityGuardrailException().TestResourceData()
		d.SetId(id)
		if _, err := resourceAiApplicationSecurityGuardrailExceptionImport(context.Background(), d, nil); err != nil {
			t.Errorf("import ID %q: unexpected error %v", id, err)
			continue
		}
		if d.Id() != testGexID || d.Get("application_id") != testGexAppID || d.Get("policy_id") != testGexPolicyID {
			t.Errorf("import ID %q: got id %q, application %q, policy %q", id, d.Id(), d.Get("application_id"), d.Get("policy_id"))
		}
	}
	d := resourceAiApplicationSecurityGuardrailException().TestResourceData()
	d.SetId(testGexAppID + "/" + testGexID)
	if _, err := resourceAiApplicationSecurityGuardrailExceptionImport(context.Background(), d, nil); err == nil || !strings.Contains(err.Error(), "invalid import ID") {
		t.Errorf("two parts must give the format error, got %v", err)
	}
}

// testAccGuardrailExceptionConfig is the policy test's mandatory policy plus one guardrail exception on it.
func testAccGuardrailExceptionConfig(attributes string) string {
	return testAccAiApplicationSecurityPolicyConfig("gex-policy", true) + fmt.Sprintf(`
resource "incapsula_ai_application_security_guardrail_exception" "test" {
  application_id = incapsula_ai_application_security_application.policy_app.id
  policy_id      = incapsula_ai_application_security_policy.test.id
%s}

output "gex_guardrail_id" {
  value = incapsula_ai_application_security_guardrail_exception.test.guardrail_id
}
`, attributes)
}

func getGuardrailExceptionOf(rs *terraform.ResourceState) (*AiApplicationSecurityGuardrailException, error) {
	accountID, _ := strconv.Atoi(rs.Primary.Attributes["account_id"])
	return testAccProvider.Meta().(*Client).GetAiApplicationSecurityGuardrailException(accountID,
		rs.Primary.Attributes["application_id"], rs.Primary.Attributes["policy_id"], rs.Primary.ID)
}

func testAccCheckGuardrailExceptionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		gex, err := getGuardrailExceptionOf(rs)
		if err != nil {
			return err
		}
		if gex == nil {
			return fmt.Errorf("guardrail exception %s not found", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckGuardrailExceptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_ai_application_security_guardrail_exception" || rs.Primary.ID == "" {
			continue
		}
		if gex, err := getGuardrailExceptionOf(rs); err == nil && gex != nil {
			return fmt.Errorf("guardrail exception still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}
