package incapsula

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_condition. The theme throughout is the expression:
// the API stores a reformatted `code`, which the provider records in
// `code_normalized` instead of `code`, so applying must be stable, an
// out-of-band expression change cannot be reported as drift, and an update
// that does not touch `code` must not overwrite the stored expression.

const abpConditionTestAddress = "incapsula_abp_condition.test"

func abpConditionTestConfig(name, description, code string) string {
	return abpConditionTestConfigIn(abpFakeAccountId, name, description, code)
}

func abpConditionTestConfigIn(accountId, name, description, code string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition" "test" {
  account_id  = %q
  name        = %q
  description = %q
  code        = %q
}
`, accountId, name, description, code)
}

func abpConditionTestConfigWithTemplate(name, template, code string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition" "test" {
  account_id  = %q
  name        = %q
  description = "d"
  template    = %q
  code        = %q
}
`, abpFakeAccountId, name, template, code)
}

// abpConditionTestStored returns one field of the Condition the fake holds.
func abpConditionTestStored(t *testing.T, f *abpFake, id, field string) string {
	t.Helper()
	doc := f.doc(abpFakeConditions, id)
	if doc == nil {
		t.Fatalf("Condition %s is not stored", id)
	}
	return abpFakeString(doc[field])
}

// The expression the most recent write carried.
func abpConditionTestSentCode(t *testing.T, f *abpFake) string {
	t.Helper()
	return abpFakeString(f.lastWrite(t)["code"])
}

func TestAbpConditionCreate(t *testing.T) {
	f := newAbpFake(t)
	const code = "(all flags.no_token)"
	config := abpConditionTestConfig("cond", "a description", code)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					id := attrs["id"]
					if got, want := abpConditionTestStored(t, f, id, "code"), abpFakeFormatCode(code); got != want {
						return fmt.Errorf("stored code is %q, want %q", got, want)
					}
					// `code` keeps what the configuration said; the server's
					// form is exposed separately.
					if attrs["code"] != code {
						return fmt.Errorf("code is %q, want the configured %q", attrs["code"], code)
					}
					if got, want := attrs["code_normalized"], abpFakeFormatCode(code); got != want {
						return fmt.Errorf("code_normalized is %q, want %q", got, want)
					}
					// An omitted template is the same as `custom`.
					if attrs["template"] != "custom" {
						return fmt.Errorf("template is %q, want custom", attrs["template"])
					}
					for _, key := range []string{"created_at", "modified_at", "last_change_by"} {
						if attrs[key] == "" {
							return fmt.Errorf("%s is empty", key)
						}
					}
					if attrs["account_id"] != abpFakeAccountId {
						return fmt.Errorf("account_id is %q", attrs["account_id"])
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Reformatting is idempotent, so a configuration that already carries the
// server's form converges just as well.
func TestAbpConditionCreateWithAlreadyFormattedCode(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(requests_per_minute > 7)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttr(abpConditionTestAddress,
					"code_normalized", "(requests_per_minute > 7)\n"),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionUpdateInPlace(t *testing.T) {
	f := newAbpFake(t)
	var id, modifiedAt string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfig("before", "d1", "(all flags.no_token)"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					id, modifiedAt = attrs["id"], attrs["modified_at"]
					return nil
				},
			},
			{
				Config: abpConditionTestConfig("after", "d2", "(all flags.no_token)"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the Condition was replaced by %s", attrs["id"])
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the update deleted %v", deleted)
					}
					if got := abpConditionTestStored(t, f, id, "name"); got != "after" {
						return fmt.Errorf("stored name is %q", got)
					}
					if got := abpConditionTestStored(t, f, id, "description"); got != "d2" {
						return fmt.Errorf("stored description is %q", got)
					}
					if attrs["modified_at"] == modifiedAt {
						return fmt.Errorf("modified_at did not change")
					}
					return nil
				},
			},
			{Config: abpConditionTestConfig("after", "d2", "(all flags.no_token)"), PlanOnly: true},
		},
	})
}

// The API cannot move a Condition between accounts, so account_id is ForceNew.
func TestAbpConditionAccountIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	var first string
	config := abpConditionTestConfigIn(abpFakeOtherAccountId, "cond", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfig("cond", "d", "(all flags.no_token)"),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] == first {
						return fmt.Errorf("the Condition was updated in place")
					}
					if f.exists(abpFakeConditions, first) {
						return fmt.Errorf("the old Condition %s survived", first)
					}
					if got := abpConditionTestStored(t, f, attrs["id"], "account_id"); got != abpFakeOtherAccountId {
						return fmt.Errorf("the new Condition belongs to account %q", got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionNameDriftIsCorrected(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(all flags.no_token)")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["name"] = "renamed in the UI"
						doc["description"] = "changed in the UI"
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpConditionTestStored(t, f, id, "name"); got != "cond" {
						return fmt.Errorf("the name was not restored, it is %q", got)
					}
					if got := abpConditionTestStored(t, f, id, "description"); got != "d" {
						return fmt.Errorf("the description was not restored, it is %q", got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// An out-of-band expression change produces no diff at all: the provider only
// ever writes the server's expression into `code_normalized`, which is
// Computed. This is a documented limitation, not an intention — comparing
// `code` with `code_normalized` is the only way to notice the divergence.
func TestAbpConditionCodeDriftIsInvisible(t *testing.T) {
	f := newAbpFake(t)
	const configured = "(requests_per_minute > 7)"
	const outOfBand = "(in visitor_ip 9.9.9.9)"
	config := abpConditionTestConfig("cond", "d", configured)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["code"] = abpFakeFormatCode(outOfBand)
					})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if got, want := attrs["code_normalized"], abpFakeFormatCode(outOfBand); got != want {
						return fmt.Errorf("code_normalized is %q, want the server's %q", got, want)
					}
					if attrs["code"] != configured {
						return fmt.Errorf("code is %q, want the configured %q", attrs["code"], configured)
					}
					if got, want := abpConditionTestStored(t, f, id, "code"), abpFakeFormatCode(outOfBand); got != want {
						return fmt.Errorf("the stored expression became %q, want %q", got, want)
					}
					return nil
				},
			},
		},
	})
}

// An update that does not change `code` must leave the stored expression
// alone: the plan does not mention the expression, so the apply must not
// change it. Before this was fixed, changing only the description replaced the
// live expression with the configured one.
func TestAbpConditionUnrelatedUpdateKeepsStoredCode(t *testing.T) {
	f := newAbpFake(t)
	const outOfBand = "(in visitor_ip 9.9.9.9)"
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfig("cond", "d1", "(requests_per_minute > 7)"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["code"] = abpFakeFormatCode(outOfBand)
					})
				},
				Config: abpConditionTestConfig("cond", "d2", "(requests_per_minute > 7)"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the Condition was replaced by %s", attrs["id"])
					}
					want := abpFakeFormatCode(outOfBand)
					if got := abpConditionTestStored(t, f, id, "code"); got != want {
						return fmt.Errorf("a description-only update changed the stored expression to %q, want the untouched %q", got, want)
					}
					if got := abpConditionTestSentCode(t, f); got != want {
						return fmt.Errorf("the update sent %q, want the stored expression %q", got, want)
					}
					if got := abpConditionTestStored(t, f, id, "description"); got != "d2" {
						return fmt.Errorf("the description was not updated, it is %q", got)
					}
					if attrs["code_normalized"] != want {
						return fmt.Errorf("code_normalized is %q, want %q", attrs["code_normalized"], want)
					}
					return nil
				},
			},
			{Config: abpConditionTestConfig("cond", "d2", "(requests_per_minute > 7)"), PlanOnly: true},
		},
	})
}

// Changing `code` in the configuration is the one thing that replaces the
// stored expression, including one that was changed outside Terraform.
func TestAbpConditionConfiguredCodeChangeIsApplied(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(requests_per_minute > 42)")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfig("cond", "d", "(requests_per_minute > 7)"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["code"] = abpFakeFormatCode("(in visitor_ip 9.9.9.9)")
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the Condition was replaced by %s", attrs["id"])
					}
					want := abpFakeFormatCode("(requests_per_minute > 42)")
					if got := abpConditionTestStored(t, f, id, "code"); got != want {
						return fmt.Errorf("the configured expression was not applied: stored %q, want %q", got, want)
					}
					if attrs["code_normalized"] != want {
						return fmt.Errorf("code_normalized is %q, want %q", attrs["code_normalized"], want)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// `template` is fixed when the Condition is created: the API ignores a new
// value, so the apply reports success while nothing changes and the plan keeps
// showing the same change. A limitation — the attribute is not ForceNew, and
// the only way out is to replace the Condition.
func TestAbpConditionTemplateIsFixedAtCreate(t *testing.T) {
	f := newAbpFake(t)
	var id string
	changed := abpConditionTestConfigWithTemplate("cond", "ip_set", "(in visitor_ip 1.2.3.4)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfigWithTemplate("cond", "rate_limiting", "(in visitor_ip 1.2.3.4)"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					id = attrs["id"]
					if attrs["template"] != "rate_limiting" {
						return fmt.Errorf("template is %q at create", attrs["template"])
					}
					return nil
				},
			},
			{
				Config:             changed,
				ExpectNonEmptyPlan: true,
				Check: func(s *terraform.State) error {
					if attrs := abpAttrs(t, s, abpConditionTestAddress); attrs["id"] != id {
						return fmt.Errorf("the Condition was replaced by %s", attrs["id"])
					}
					if got := abpConditionTestStored(t, f, id, "template"); got != "rate_limiting" {
						return fmt.Errorf("the stored template became %q", got)
					}
					return nil
				},
			},
			// The plan never converges.
			{Config: changed, PlanOnly: true, ExpectNonEmptyPlan: true},
		},
	})
}

// The update endpoint is a full replace, so the provider overlays its payload
// onto the document the API holds. Fields it does not model must ride along
// untouched instead of being cleared.
func TestAbpConditionUpdatePreservesUnmodelledFields(t *testing.T) {
	f := newAbpFake(t)
	const code = "(all flags.no_token)"
	final := abpConditionTestConfig("cond", "d2", code)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionTestConfig("cond", "d1", code),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["experimental_setting"] = "keep me"
					})
				},
				Config: final,
				Check: func(s *terraform.State) error {
					if got := abpFakeString(f.doc(abpFakeConditions, id)["experimental_setting"]); got != "keep me" {
						return fmt.Errorf("the unmodelled field is %q after the update", got)
					}
					if got := abpConditionTestStored(t, f, id, "description"); got != "d2" {
						return fmt.Errorf("the description was not updated, it is %q", got)
					}
					return nil
				},
			},
			{Config: final, PlanOnly: true},
		},
	})
}

func TestAbpConditionDeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(all flags.no_token)")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakeConditions, id) },
				Config:    config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] == id {
						return fmt.Errorf("the vanished Condition %s is still in state", id)
					}
					if !f.exists(abpFakeConditions, attrs["id"]) {
						return fmt.Errorf("the Condition was not recreated")
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionDestroyDeletesIt(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if f.count(abpFakeConditions) != 0 {
				return fmt.Errorf("the Condition survived the destroy")
			}
			if len(f.deletedIds()) != 1 {
				return fmt.Errorf("the destroy deleted %v", f.deletedIds())
			}
			return nil
		},
		Steps: []resource.TestStep{{Config: config}},
	})
}

func TestAbpConditionImport(t *testing.T) {
	f := newAbpFake(t)
	const code = "(requests_per_minute > 100)"
	config := abpConditionTestConfig("imported", "d", code)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:       config,
				ResourceName: abpConditionTestAddress,
				ImportState:  true,
				// `code` is seeded from the server's reformatted expression on
				// import, so it does not match the configured one.
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("the import produced %d resources", len(states))
					}
					if got, want := states[0].Attributes["code"], abpFakeFormatCode(code); got != want {
						return fmt.Errorf("the imported code is %q, want the server's %q", got, want)
					}
					return nil
				},
			},
		},
	})
}

// Because `code` is seeded from the server's reformatted expression, the first
// plan after an import is never empty. The apply must keep the Condition and
// leave the expression intact, and the plan must converge afterwards.
func TestAbpConditionImportThenApplyKeepsCondition(t *testing.T) {
	f := newAbpFake(t)
	const code = "(requests_per_minute > 100)"
	id := f.seedCondition(abpFakeAccountId, "seeded", "d", code)
	config := abpConditionTestConfig("seeded", "d", code)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpConditionTestAddress,
				ImportState:        true,
				ImportStateId:      id,
				ImportStatePersist: true,
			},
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the imported Condition %s was replaced by %s", id, attrs["id"])
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the first apply after import deleted %v", deleted)
					}
					if got, want := abpConditionTestStored(t, f, id, "code"), abpFakeFormatCode(code); got != want {
						return fmt.Errorf("the stored expression became %q, want %q", got, want)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A managed (Imperva-owned) Condition carries no account and cannot be managed
// as a resource.
func TestAbpConditionManagedCannotBeImported(t *testing.T) {
	f := newAbpFake(t)
	id := abpConditionSeedManaged(f, "Monitoring Tools", "maintained by Imperva")
	config := abpConditionTestConfig("Monitoring Tools", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionTestAddress,
				ImportState:   true,
				ImportStateId: id,
				ExpectError:   regexp.MustCompile(`is a managed condition and cannot be imported`),
			},
		},
	})
}

// Only the literal variant is an incapsula_abp_condition. Importing a
// Condition List as one fails on the read that follows the import.
func TestAbpConditionImportRejectsConditionList(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList(abpFakeAccountId, "a list", "d")
	config := abpConditionTestConfig("a list", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionTestAddress,
				ImportState:   true,
				ImportStateId: id,
				ExpectError:   regexp.MustCompile(`is not a literal variant \(it is a list\)`),
			},
		},
	})
}

func TestAbpConditionImportOfAbsentIdFails(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionTestConfig("cond", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Condition .* not found`),
			},
		},
	})
}

func TestAbpConditionValidation(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionTestConfigIn("not-a-uuid", "cond", "d", "(all flags.no_token)"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
			{
				Config:      abpConditionTestConfig(strings.Repeat("x", 101), "d", "(all flags.no_token)"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 100\)`),
			},
			{
				Config:      abpConditionTestConfig("cond", "d", "   "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "code" to not be an empty string or whitespace`),
			},
		},
	})
}

// The 1..100 name limit is the provider's own, so a longer Condition created
// elsewhere imports fine but can never be expressed in a configuration. A
// limitation rather than an intention.
func TestAbpConditionLongNameImportsButCannotBePlanned(t *testing.T) {
	f := newAbpFake(t)
	longName := strings.Repeat("n", 120)
	id := f.seedCondition(abpFakeAccountId, longName, "d", "(all flags.no_token)")
	placeholder := abpConditionTestConfig("short", "d", "(all flags.no_token)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        placeholder,
				ResourceName:  abpConditionTestAddress,
				ImportState:   true,
				ImportStateId: id,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if got := states[0].Attributes["name"]; got != longName {
						return fmt.Errorf("the imported name is %q", got)
					}
					return nil
				},
			},
			{
				Config:      abpConditionTestConfig(longName, "d", "(all flags.no_token)"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 100\)`),
			},
		},
	})
}
