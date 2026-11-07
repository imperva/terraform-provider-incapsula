package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_condition_list_entry. An entry is just the link
// between a parent Condition List and either a literal Condition or a nested
// Condition List, so the interesting parts are which kind it may reference,
// what makes it be replaced rather than updated, and that tags and state
// round-trip.

const abpConditionEntryTestAddress = "incapsula_abp_condition_list_entry.test"

// abpConditionEntryTestConfig takes the reference and the tags block
// verbatim, so a test can state whichever variant it needs.
func abpConditionEntryTestConfig(accountId, parent, reference, state, tags string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition_list_entry" "test" {
  account_id               = %q
  parent_condition_list_id = %q
  %s
  state = %q
  %s
}
`, accountId, parent, reference, state, tags)
}

// abpConditionEntryTestFixture seeds a parent Condition List, a literal
// Condition and a nested Condition List to point entries at.
type abpConditionEntryTestFixture struct {
	parent     string
	literal    string
	other      string
	nestedList string
}

func abpConditionEntryTestSeed(f *abpFake) abpConditionEntryTestFixture {
	return abpConditionEntryTestFixture{
		parent:     f.seedConditionList(abpFakeAccountId, "parent list", "d"),
		literal:    f.seedCondition(abpFakeAccountId, "okhttp", "d", "(matches headers.user_agent \"okhttp\")"),
		other:      f.seedCondition(abpFakeAccountId, "curl", "d", "(matches headers.user_agent \"curl\")"),
		nestedList: f.seedConditionList(abpFakeAccountId, "nested list", "d"),
	}
}

func TestAbpConditionListEntryCreate(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive,
		`tags = ["terraform_managed", "okhttp"]`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_id", fx.literal),
					// The unselected reference attribute is cleared, not left
					// unknown.
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_list_id", ""),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "state", abpConditionStateActive),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr(abpConditionEntryTestAddress, "tags.*", "terraform_managed"),
					resource.TestCheckTypeSetElemAttr(abpConditionEntryTestAddress, "tags.*", "okhttp"),
					func(s *terraform.State) error {
						attrs := abpAttrs(t, s, abpConditionEntryTestAddress)
						doc := f.doc(abpFakeConditions, attrs["id"])
						if doc == nil {
							return fmt.Errorf("the entry was not stored")
						}
						if abpFakeConditionKind(doc) != AbpConditionKindReference {
							return fmt.Errorf("the stored Condition is a %s", abpFakeConditionKind(doc))
						}
						if abpFakeString(doc["parent"]) != fx.parent {
							return fmt.Errorf("the entry is parented by %q", doc["parent"])
						}
						if abpFakeString(doc["reference"]) != fx.literal {
							return fmt.Errorf("the entry references %q", doc["reference"])
						}
						for _, key := range []string{"created_at", "modified_at"} {
							if attrs[key] == "" {
								return fmt.Errorf("%s is empty", key)
							}
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// An entry may reference a whole Condition List instead of a single Condition.
func TestAbpConditionListEntryCreateReferencingList(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_list_id = %q", fx.nestedList), abpConditionStateMonitor, "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_list_id", fx.nestedList),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_id", ""),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "state", abpConditionStateMonitor),
					// An omitted tags argument is an empty tag set, which the
					// API requires to be sent explicitly.
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "tags.#", "0"),
					func(s *terraform.State) error {
						tags, ok := f.lastWrite(t)["tags"].([]any)
						if !ok || len(tags) != 0 {
							return fmt.Errorf("the create sent tags %v, want an empty array", f.lastWrite(t)["tags"])
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Only tags and state can be changed without replacing the entry.
func TestAbpConditionListEntryUpdateTagsAndStateInPlace(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	reference := fmt.Sprintf("condition_id = %q", fx.literal)
	final := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference,
		abpConditionStateInactive, `tags = ["kept", "added"]`)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference,
					abpConditionStateActive, `tags = ["kept", "removed"]`),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: final,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionEntryTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the entry was replaced by %s", attrs["id"])
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the update deleted %v", deleted)
					}
					doc := f.doc(abpFakeConditions, id)
					if abpFakeString(doc["state"]) != abpConditionStateInactive {
						return fmt.Errorf("stored state is %q", doc["state"])
					}
					got := abpFakeStrings(doc["tags"])
					want := map[string]bool{"kept": true, "added": true}
					if len(got) != len(want) {
						return fmt.Errorf("stored tags are %v", got)
					}
					for _, tag := range got {
						if !want[tag] {
							return fmt.Errorf("stored tags are %v", got)
						}
					}
					return nil
				},
			},
			{Config: final, PlanOnly: true},
		},
	})
}

// Clearing the tag set must reach the API: an entry with no tags is a valid
// state, not a missing value.
func TestAbpConditionListEntryTagsCanBeCleared(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	reference := fmt.Sprintf("condition_id = %q", fx.literal)
	cleared := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference,
		abpConditionStateActive, `tags = []`)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference,
					abpConditionStateActive, `tags = ["terraform_managed"]`),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: cleared,
				Check: func(s *terraform.State) error {
					if got := abpFakeStrings(f.doc(abpFakeConditions, id)["tags"]); len(got) != 0 {
						return fmt.Errorf("stored tags are %v, want none", got)
					}
					return nil
				},
			},
			{Config: cleared, PlanOnly: true},
		},
	})
}

// Re-pointing an entry replaces it: the API could update the reference in
// place, but the provider forces a new entry, so the Condition briefly does
// not apply. A limitation rather than an intention.
func TestAbpConditionListEntryRepointingReplacesIt(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	second := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.other), abpConditionStateActive, "")
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, ""),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: second,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionEntryTestAddress)
					if attrs["id"] == first {
						return fmt.Errorf("the entry was updated in place")
					}
					if f.exists(abpFakeConditions, first) {
						return fmt.Errorf("the old entry %s survived", first)
					}
					// Only the link is recreated.
					if !f.exists(abpFakeConditions, fx.literal) {
						return fmt.Errorf("the previously referenced Condition was deleted")
					}
					if !f.exists(abpFakeConditions, fx.parent) {
						return fmt.Errorf("the parent Condition List was deleted")
					}
					return nil
				},
			},
			{Config: second, PlanOnly: true},
		},
	})
}

// Swapping a literal reference for a Condition List reference replaces the
// entry too, and the unselected attribute is cleared.
func TestAbpConditionListEntrySwitchingReferenceKindReplacesIt(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	asList := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_list_id = %q", fx.nestedList), abpConditionStateActive, "")
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, ""),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: asList,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_id", ""),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_list_id", fx.nestedList),
					func(s *terraform.State) error {
						if abpAttrs(t, s, abpConditionEntryTestAddress)["id"] == first {
							return fmt.Errorf("the entry was updated in place")
						}
						return nil
					},
				),
			},
			{Config: asList, PlanOnly: true},
		},
	})
}

func TestAbpConditionListEntryParentChangeReplacesIt(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	otherParent := f.seedConditionList(abpFakeAccountId, "another parent", "d")
	reference := fmt.Sprintf("condition_id = %q", fx.literal)
	moved := abpConditionEntryTestConfig(abpFakeAccountId, otherParent, reference, abpConditionStateActive, "")
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference, abpConditionStateActive, ""),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: moved,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionEntryTestAddress)
					if attrs["id"] == first {
						return fmt.Errorf("the entry was updated in place")
					}
					if got := abpFakeString(f.doc(abpFakeConditions, attrs["id"])["parent"]); got != otherParent {
						return fmt.Errorf("the new entry is parented by %q", got)
					}
					return nil
				},
			},
			{Config: moved, PlanOnly: true},
		},
	})
}

func TestAbpConditionListEntryAccountIdChangeReplacesIt(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	reference := fmt.Sprintf("condition_id = %q", fx.literal)
	moved := abpConditionEntryTestConfig(abpFakeOtherAccountId, fx.parent, reference, abpConditionStateActive, "")
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference, abpConditionStateActive, ""),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				Config: moved,
				Check: func(s *terraform.State) error {
					if abpAttrs(t, s, abpConditionEntryTestAddress)["id"] == first {
						return fmt.Errorf("the entry was updated in place")
					}
					if f.exists(abpFakeConditions, first) {
						return fmt.Errorf("the old entry %s survived", first)
					}
					return nil
				},
			},
			{Config: moved, PlanOnly: true},
		},
	})
}

// The provider checks the kind of the referenced Condition before writing,
// because the API does not.
func TestAbpConditionListEntryReferencedKindIsChecked(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	entry := f.seedEntry(fx.parent, fx.literal, abpConditionStateActive)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q", fx.nestedList), abpConditionStateActive, ""),
				ExpectError: regexp.MustCompile(`condition_id .* must reference a literal condition, but it is a list`),
			},
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_list_id = %q", fx.literal), abpConditionStateActive, ""),
				ExpectError: regexp.MustCompile(`condition_list_id .* must reference a list condition, but it is a literal`),
			},
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q", entry), abpConditionStateActive, ""),
				ExpectError: regexp.MustCompile(`condition_id .* must reference a literal condition, but it is a reference`),
			},
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q", abpFakeAbsentId), abpConditionStateActive, ""),
				ExpectError: regexp.MustCompile(`condition_id .* does not exist`),
			},
		},
	})
}

// The parent is not checked client-side, so an unusable parent is only
// rejected by the API — and the message names the generic "ABP Condition"
// rather than the entry. A limitation rather than an intention.
func TestAbpConditionListEntryUnknownParentIsRejectedByTheApi(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, abpFakeAbsentId,
					fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, ""),
				ExpectError: regexp.MustCompile(`creating ABP Condition in ABP account .*Unknown parent condition list`),
			},
		},
	})
}

func TestAbpConditionListEntryValidation(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	reference := fmt.Sprintf("condition_id = %q", fx.literal)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, "", abpConditionStateActive, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`one of .condition_id,condition_list_id. must be specified`),
			},
			{
				Config: abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
					fmt.Sprintf("condition_id = %q\n  condition_list_id = %q", fx.literal, fx.nestedList),
					abpConditionStateActive, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only one of .condition_id,condition_list_id. can be`),
			},
			{
				Config:      abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference, "bogus", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected state to be one of`),
			},
			{
				Config:      abpConditionEntryTestConfig(abpFakeAccountId, fx.parent, reference, abpConditionStateActive, `tags = ["Not-Snake-Case"]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`tag must be snake_case`),
			},
			{
				Config:      abpConditionEntryTestConfig(abpFakeAccountId, "not-a-uuid", reference, abpConditionStateActive, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "parent_condition_list_id" to be a valid UUID`),
			},
			{
				Config:      abpConditionEntryTestConfig("not-a-uuid", fx.parent, reference, abpConditionStateActive, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
		},
	})
}

func TestAbpConditionListEntryStateAndTagDriftAreCorrected(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive,
		`tags = ["terraform_managed"]`)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["state"] = abpConditionStateInactive
						doc["tags"] = []any{"added_in_the_ui"}
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					doc := f.doc(abpFakeConditions, id)
					if abpFakeString(doc["state"]) != abpConditionStateActive {
						return fmt.Errorf("the state was not restored, it is %q", doc["state"])
					}
					got := abpFakeStrings(doc["tags"])
					if len(got) != 1 || got[0] != "terraform_managed" {
						return fmt.Errorf("the tags were not restored, they are %v", got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListEntryDeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, "")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakeConditions, id) },
				Config:    config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionEntryTestAddress)
					if attrs["id"] == id {
						return fmt.Errorf("the vanished entry %s is still in state", id)
					}
					if !f.exists(abpFakeConditions, attrs["id"]) {
						return fmt.Errorf("the entry was not recreated")
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// When the referenced Condition disappears, the refresh fails instead of
// proposing to recreate the entry, so the configuration has to be unblocked by
// hand. A limitation rather than an intention.
func TestAbpConditionListEntryReferencedConditionDeletedFailsRefresh(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig:   func() { f.drop(abpFakeConditions, fx.literal) },
				Config:      config,
				ExpectError: regexp.MustCompile(`referenced condition .* no longer exists`),
			},
		},
	})
}

// Destroying an entry removes only the link.
func TestAbpConditionListEntryDestroyKeepsParentAndReference(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, "")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if f.exists(abpFakeConditions, id) {
				return fmt.Errorf("the entry survived the destroy")
			}
			if !f.exists(abpFakeConditions, fx.parent) || !f.exists(abpFakeConditions, fx.literal) {
				return fmt.Errorf("the destroy took the parent or the referenced Condition with it")
			}
			if got := f.deletedIds(); len(got) != 1 || got[0] != id {
				return fmt.Errorf("the destroy deleted %v, want just the entry %s", got, id)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
					return nil
				},
			},
		},
	})
}

func TestAbpConditionListEntryImport(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive,
		`tags = ["terraform_managed"]`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpConditionEntryTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// An entry created elsewhere imports completely, so the first apply after an
// import plans nothing.
func TestAbpConditionListEntryImportThenApplyPlansNothing(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	id := f.seedEntry(fx.parent, fx.nestedList, abpConditionStateMonitor, "added_elsewhere")
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_list_id = %q", fx.nestedList), abpConditionStateMonitor,
		`tags = ["added_elsewhere"]`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpConditionEntryTestAddress,
				ImportState:        true,
				ImportStateId:      id,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					attrs := states[0].Attributes
					if attrs["condition_list_id"] != fx.nestedList {
						return fmt.Errorf("the imported condition_list_id is %q", attrs["condition_list_id"])
					}
					if attrs["condition_id"] != "" {
						return fmt.Errorf("the imported condition_id is %q, want empty", attrs["condition_id"])
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListEntryImportRejectsNonEntry(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	config := abpConditionEntryTestConfig(abpFakeAccountId, fx.parent,
		fmt.Sprintf("condition_id = %q", fx.literal), abpConditionStateActive, "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionEntryTestAddress,
				ImportState:   true,
				ImportStateId: fx.literal,
				ExpectError:   regexp.MustCompile(`is not a reference variant \(it is a literal\)`),
			},
			{
				Config:        config,
				ResourceName:  abpConditionEntryTestAddress,
				ImportState:   true,
				ImportStateId: fx.parent,
				ExpectError:   regexp.MustCompile(`is not a reference variant \(it is a list\)`),
			},
			{
				Config:        config,
				ResourceName:  abpConditionEntryTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Condition List Entry .* not found`),
			},
		},
	})
}

// Attaching a Condition to a Policy Directive is done by parenting an entry on
// the Directive's Condition List, which the Directive data source exposes.
func TestAbpConditionListEntryInPolicyDirective(t *testing.T) {
	f := newAbpFake(t)
	fx := abpConditionEntryTestSeed(f)
	policy := f.seedPolicy("a policy", "d", "block", "proof_of_work")
	directives := abpFakeList(f.doc(abpFakePolicies, policy)["directives"])
	blockList := abpFakeString(directives[0]["condition_id"])
	skipList := abpFakeString(directives[1]["skip_condition_id"])

	config := fmt.Sprintf(`
data "incapsula_abp_directive" "block" {
  policy_id = %[1]q
  action    = "block"
}

data "incapsula_abp_directive" "pow" {
  policy_id = %[1]q
  action    = "proof_of_work"
}

resource "incapsula_abp_condition_list_entry" "test" {
  account_id               = %[2]q
  parent_condition_list_id = data.incapsula_abp_directive.block.condition_list_id
  condition_id             = %[3]q
  state                    = "active"
  tags                     = ["terraform_managed"]
}

resource "incapsula_abp_condition_list_entry" "skip" {
  account_id               = %[2]q
  parent_condition_list_id = data.incapsula_abp_directive.pow.skip_condition_list_id
  condition_id             = %[3]q
  state                    = "active"
}
`, policy, abpFakeAccountId, fx.literal)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress,
						"parent_condition_list_id", blockList),
					resource.TestCheckResourceAttr("incapsula_abp_condition_list_entry.skip",
						"parent_condition_list_id", skipList),
					func(s *terraform.State) error {
						id := abpAttrs(t, s, abpConditionEntryTestAddress)["id"]
						if got := abpFakeString(f.doc(abpFakeConditions, id)["parent"]); got != blockList {
							return fmt.Errorf("the entry is parented by %q, want the Directive's list %q", got, blockList)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A cycle — an entry inside a Condition List that references that very list —
// is accepted by both the provider and the API. A limitation rather than an
// intention; the effect downstream is unspecified.
func TestAbpConditionListEntrySelfReferenceIsAccepted(t *testing.T) {
	f := newAbpFake(t)
	list := f.seedConditionList(abpFakeAccountId, "self referencing", "d")
	config := abpConditionEntryTestConfig(abpFakeAccountId, list,
		fmt.Sprintf("condition_list_id = %q", list), abpConditionStateActive, "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "parent_condition_list_id", list),
					resource.TestCheckResourceAttr(abpConditionEntryTestAddress, "condition_list_id", list),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}
