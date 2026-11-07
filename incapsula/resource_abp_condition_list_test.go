package incapsula

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_condition_list. Two things shape the resource: the
// `description` attribute cannot express "no description", and destroying a
// list takes every entry inside it with it, including entries this
// configuration does not manage.

const abpConditionListTestAddress = "incapsula_abp_condition_list.test"

func abpConditionListTestConfig(name, description string) string {
	return abpConditionListTestConfigIn(abpFakeAccountId, name, description)
}

func abpConditionListTestConfigIn(accountId, name, description string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition_list" "test" {
  account_id  = %q
  name        = %q
  description = %q
}
`, accountId, name, description)
}

// abpConditionListTestConfigNoDescription omits `description` entirely.
func abpConditionListTestConfigNoDescription(name string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition_list" "test" {
  account_id = %q
  name       = %q
}
`, abpFakeAccountId, name)
}

func TestAbpConditionListCreate(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfig("a list", "a description")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionListTestAddress)
					doc := f.doc(abpFakeConditions, attrs["id"])
					if doc == nil {
						return fmt.Errorf("the Condition List was not stored")
					}
					// A Condition List is the list variant: no code, no
					// reference.
					if abpFakeConditionKind(doc) != AbpConditionKindList {
						return fmt.Errorf("the stored Condition is a %s", abpFakeConditionKind(doc))
					}
					if abpFakeString(doc["name"]) != "a list" {
						return fmt.Errorf("stored name is %q", doc["name"])
					}
					if attrs["description"] != "a description" {
						return fmt.Errorf("description is %q", attrs["description"])
					}
					for _, key := range []string{"created_at", "modified_at"} {
						if attrs[key] == "" {
							return fmt.Errorf("%s is empty", key)
						}
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Condition List without a description is the normal case; it stays empty.
func TestAbpConditionListCreateWithoutDescription(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfigNoDescription("a list")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpConditionListTestAddress, "description", ""),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListUpdateInPlace(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfig("after", "d2")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListTestConfig("before", "d1"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionListTestAddress)
					if attrs["id"] != id {
						return fmt.Errorf("the Condition List was replaced by %s", attrs["id"])
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the update deleted %v", deleted)
					}
					doc := f.doc(abpFakeConditions, id)
					if abpFakeString(doc["name"]) != "after" || abpFakeString(doc["description"]) != "d2" {
						return fmt.Errorf("stored name/description are %q/%q", doc["name"], doc["description"])
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListAccountIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfigIn(abpFakeOtherAccountId, "a list", "d")
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListTestConfig("a list", "d"),
				Check: func(s *terraform.State) error {
					first = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionListTestAddress)
					if attrs["id"] == first {
						return fmt.Errorf("the Condition List was updated in place")
					}
					if f.exists(abpFakeConditions, first) {
						return fmt.Errorf("the old Condition List %s survived", first)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// `description` is Optional + Computed, so an empty string cannot be told from
// an omitted value: a description that has already been set cannot be cleared.
// A limitation rather than an intention.
func TestAbpConditionListEmptyDescriptionIsIgnored(t *testing.T) {
	f := newAbpFake(t)
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListTestConfig("a list", "d"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			// Asking for no description at all plans nothing.
			{Config: abpConditionListTestConfig("a list", ""), PlanOnly: true},
			{Config: abpConditionListTestConfigNoDescription("a list"), PlanOnly: true},
			{
				Config: abpConditionListTestConfig("a list", ""),
				Check: func(s *terraform.State) error {
					if got := abpFakeString(f.doc(abpFakeConditions, id)["description"]); got != "d" {
						return fmt.Errorf("the stored description became %q, want the kept %q", got, "d")
					}
					if got := abpAttrs(t, s, abpConditionListTestAddress)["description"]; got != "d" {
						return fmt.Errorf("state description is %q, want the kept %q", got, "d")
					}
					return nil
				},
			},
		},
	})
}

// A description changed outside Terraform is adopted into state rather than
// restored when the configuration does not state one, because `description` is
// Optional + Computed. A limitation rather than an intention.
func TestAbpConditionListDescriptionDriftIsAdopted(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfigNoDescription("a list")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["description"] = "changed in the UI"
					})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpConditionListTestAddress)["description"]; got != "changed in the UI" {
						return fmt.Errorf("state description is %q, want the adopted one", got)
					}
					return nil
				},
			},
		},
	})
}

// A name changed outside Terraform is restored, unlike the description.
func TestAbpConditionListNameDriftIsCorrected(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfig("a list", "d")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["name"] = "renamed in the UI"
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpFakeString(f.doc(abpFakeConditions, id)["name"]); got != "a list" {
						return fmt.Errorf("the name was not restored, it is %q", got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListDeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfig("a list", "d")
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakeConditions, id) },
				Config:    config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionListTestAddress)
					if attrs["id"] == id {
						return fmt.Errorf("the vanished Condition List %s is still in state", id)
					}
					if !f.exists(abpFakeConditions, attrs["id"]) {
						return fmt.Errorf("the Condition List was not recreated")
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Destroying a Condition List asks the API to delete the list itself and
// nothing else: neither the Conditions its entries point at nor the entries
// are deleted by Terraform. The service deletes the entries parented by the
// list on its own — including entries created elsewhere, which is what the
// resource's warning is about — but that cascade happens server side, so
// Terraform never learns of it.
func TestAbpConditionListDestroyDeletesOnlyTheList(t *testing.T) {
	f := newAbpFake(t)
	referenced := f.seedCondition(abpFakeAccountId, "referenced", "d", "(all flags.no_token)")
	config := abpConditionListTestConfig("a list", "d")
	var list string
	var entry string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if f.exists(abpFakeConditions, list) {
				return fmt.Errorf("the Condition List survived the destroy")
			}
			if !f.exists(abpFakeConditions, referenced) {
				return fmt.Errorf("the referenced Condition was deleted too")
			}
			if f.exists(abpFakeConditions, entry) {
				return fmt.Errorf("the service was expected to delete entry %s with its list", entry)
			}
			if got := f.deletedIds(); len(got) != 1 || got[0] != list {
				return fmt.Errorf("the destroy deleted %v, want just the Condition List %s", got, list)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					list = abpAttrs(t, s, abpConditionListTestAddress)["id"]
					// An entry added through the UI, which this configuration
					// knows nothing about.
					entry = f.seedEntry(list, referenced, abpConditionStateActive, "added_in_the_ui")
					return nil
				},
			},
		},
	})
}

func TestAbpConditionListImport(t *testing.T) {
	f := newAbpFake(t)
	config := abpConditionListTestConfig("a list", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpConditionListTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Importing a Condition List created elsewhere and applying it must keep it:
// a replacement would delete every entry inside it.
func TestAbpConditionListImportThenApplyKeepsIt(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList(abpFakeAccountId, "seeded", "d")
	config := abpConditionListTestConfig("seeded", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpConditionListTestAddress,
				ImportState:        true,
				ImportStateId:      id,
				ImportStatePersist: true,
			},
			{Config: config, PlanOnly: true},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpConditionListTestAddress)["id"]; got != id {
						return fmt.Errorf("the imported Condition List %s was replaced by %s", id, got)
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the first apply after import deleted %v", deleted)
					}
					return nil
				},
			},
		},
	})
}

func TestAbpConditionListImportRejectsLiteral(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedCondition(abpFakeAccountId, "a literal", "d", "(all flags.no_token)")
	config := abpConditionListTestConfig("a literal", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionListTestAddress,
				ImportState:   true,
				ImportStateId: id,
				ExpectError:   regexp.MustCompile(`is not a list variant \(it is a literal\)`),
			},
		},
	})
}

func TestAbpConditionListImportRejectsEntry(t *testing.T) {
	f := newAbpFake(t)
	referenced := f.seedCondition(abpFakeAccountId, "referenced", "d", "(all flags.no_token)")
	list := f.seedConditionList(abpFakeAccountId, "a list", "d")
	entry := f.seedEntry(list, referenced, abpConditionStateActive)
	config := abpConditionListTestConfig("a list", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionListTestAddress,
				ImportState:   true,
				ImportStateId: entry,
				ExpectError:   regexp.MustCompile(`is not a list variant \(it is a reference\)`),
			},
		},
	})
}

func TestAbpConditionListManagedCannotBeImported(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList("", "managed list", "d")
	config := abpConditionListTestConfig("managed list", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpConditionListTestAddress,
				ImportState:   true,
				ImportStateId: id,
				ExpectError:   regexp.MustCompile(`is a managed condition list and cannot be imported`),
			},
		},
	})
}

// A Policy Directive's Condition List is an ordinary Condition of the list
// variant, so the importer accepts it — and the resulting resource renames a
// list a live Policy owns. A limitation rather than an intention: such a list
// is not meant to be managed on its own.
func TestAbpConditionListImportOfPolicyDirectiveListRenamesIt(t *testing.T) {
	f := newAbpFake(t)
	policy := f.seedPolicy("a policy", "d", "block")
	directiveList := abpFakeString(abpFakeList(f.doc(abpFakePolicies, policy)["directives"])[0]["condition_id"])
	config := abpConditionListTestConfig("stolen from the policy", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpConditionListTestAddress,
				ImportState:        true,
				ImportStateId:      directiveList,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					// The backend names a Directive's list after its action.
					if got := states[0].Attributes["name"]; got != "block" {
						return fmt.Errorf("the imported name is %q, want the Directive's action", got)
					}
					return nil
				},
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpFakeString(f.doc(abpFakeConditions, directiveList)["name"]); got != "stolen from the policy" {
						return fmt.Errorf("the Policy's Condition List is named %q", got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpConditionListValidation(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionListTestConfigIn("not-a-uuid", "a list", "d"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
			{
				Config:      abpConditionListTestConfig(strings.Repeat("x", 101), "d"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 100\)`),
			},
			{
				Config:      abpConditionListTestConfig("", "d"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 100\)`),
			},
		},
	})
}
