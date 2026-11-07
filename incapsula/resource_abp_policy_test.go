package incapsula

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_policy, against the shared fake ABP backend. The Directive
// identity tests are the important ones: a Directive is addressed by its
// action, and the Condition Lists the backend owns have to stay with the
// action they belong to when Directives are inserted, removed or reordered.

const abpPolicyTestAddress = "incapsula_abp_policy.test"

// abpPolicyTestConfig builds a Policy with one Directive per action.
func abpPolicyTestConfig(name, description string, actions ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `
resource "incapsula_abp_policy" "test" {
  account_id  = %q
  name        = %q
  description = %q
`, abpFakeAccountId, name, description)
	for _, action := range actions {
		fmt.Fprintf(&b, "\n  directive {\n    action = %q\n  }\n", action)
	}
	b.WriteString("}\n")
	return b.String()
}

// --- baseline ---------------------------------------------------------------

// Creating a Policy and re-running the same configuration must not produce a
// plan.
func TestAbpPolicyCreateIsIdempotent(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "name", "p1"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", "desc"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.0.action", "allow"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.1.action", "block"),
					resource.TestCheckResourceAttrSet(abpPolicyTestAddress, "directive.0.condition_list_id"),
					resource.TestCheckResourceAttrSet(abpPolicyTestAddress, "directive.1.condition_list_id"),
					// Only a proof_of_work Directive gets a skip list.
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.0.skip_condition_list_id", ""),
					func(s *terraform.State) error {
						// The create request must not carry condition_id: the
						// Condition Lists are allocated by the backend.
						for _, write := range f.writes() {
							if !strings.Contains(write.Body, "\"directives\"") {
								continue
							}
							if strings.Contains(write.Body, "condition_id") {
								return fmt.Errorf("the create request sent a condition list id: %s", write.Body)
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

// A Policy resource exposes no timestamps, although both Policy data sources
// do. A limitation rather than an intention: the attributes are simply not in
// the resource's schema.
func TestAbpPolicyTimestampsAreNotExposed(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(abpPolicyTestAddress, "created_at"),
					resource.TestCheckNoResourceAttr(abpPolicyTestAddress, "modified_at"),
					func(s *terraform.State) error {
						// The API does report them.
						if abpFakeString(f.only(t, abpFakePolicies)["modified_at"]) == "" {
							return fmt.Errorf("the backend stored no modified_at, so the attribute is untestable")
						}
						return nil
					},
				),
			},
		},
	})
}

// A proof_of_work Directive gets a skip Condition List from the backend, which
// must stay stable across runs.
func TestAbpPolicyProofOfWorkDirectiveIsStable(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "incapsula_abp_policy" "test" {
  account_id  = %q
  name        = "p1"
  description = "desc"

  directive {
    action = "allow"
  }

  directive {
    action                         = "proof_of_work"
    proof_of_work_configuration_id = "eeeeeeee-0000-4000-8000-000000000001"
  }
}
`, abpFakeAccountId)
	var skipList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(abpPolicyTestAddress, "directive.1.skip_condition_list_id"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.1.proof_of_work_configuration_id", "eeeeeeee-0000-4000-8000-000000000001"),
					func(s *terraform.State) error {
						skipList = abpAttrs(t, s, abpPolicyTestAddress)["directive.1.skip_condition_list_id"]
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
			{
				// Renaming the Policy must not reallocate the skip list.
				Config: strings.Replace(config, `"p1"`, `"p2"`, 1),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPolicyTestAddress)["directive.1.skip_condition_list_id"]; got != skipList {
						return fmt.Errorf("skip condition list changed from %s to %s", skipList, got)
					}
					if !f.exists(abpFakeConditions, skipList) {
						return fmt.Errorf("skip condition list %s was deleted", skipList)
					}
					return nil
				},
			},
		},
	})
}

// Renaming a Policy and rewriting its description is an in-place update that
// keeps the Policy and its Condition Lists.
func TestAbpPolicyUpdateIsInPlace(t *testing.T) {
	f := newAbpFake(t)
	var id, allowList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					id, allowList = attrs["id"], attrs["directive.0.condition_list_id"]
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p2", "other", "allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "name", "p2"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", "other"),
					func(s *terraform.State) error {
						attrs := abpAttrs(t, s, abpPolicyTestAddress)
						if attrs["id"] != id {
							return fmt.Errorf("the policy was replaced: %s became %s", id, attrs["id"])
						}
						if attrs["directive.0.condition_list_id"] != allowList {
							return fmt.Errorf("the `allow` directive's condition list changed from %s to %s",
								allowList, attrs["directive.0.condition_list_id"])
						}
						if len(f.deletedIds()) != 0 {
							return fmt.Errorf("an in-place update deleted %v", f.deletedIds())
						}
						return nil
					},
				),
			},
		},
	})
}

// account_id is ForceNew, so moving a Policy to another Account replaces it:
// the old Policy is deleted, and with it every Condition List it owned.
func TestAbpPolicyAccountIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	first := abpPolicyTestConfig("p1", "desc", "allow")
	second := strings.Replace(first, abpFakeAccountId, abpFakeOtherAccountId, 1)

	var id, allowList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: first,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					id, allowList = attrs["id"], attrs["directive.0.condition_list_id"]
					return nil
				},
			},
			{
				Config: second,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "account_id", abpFakeOtherAccountId),
					func(s *terraform.State) error {
						if abpAttrs(t, s, abpPolicyTestAddress)["id"] == id {
							return fmt.Errorf("the policy was updated in place instead of replaced")
						}
						if f.exists(abpFakePolicies, id) {
							return fmt.Errorf("the replaced policy %s still exists", id)
						}
						if f.exists(abpFakeConditions, allowList) {
							return fmt.Errorf("condition list %s outlived the policy that owned it", allowList)
						}
						return nil
					},
				),
			},
		},
	})
}

// Deleting a Policy deletes its Directives' Condition Lists, and with them
// every entry inside those lists — including entries created elsewhere.
func TestAbpPolicyDeleteTakesConditionListsWithIt(t *testing.T) {
	f := newAbpFake(t)
	var lists []string
	var entry string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			for _, list := range lists {
				if f.exists(abpFakeConditions, list) {
					return fmt.Errorf("condition list %s survived the policy", list)
				}
			}
			if f.exists(abpFakeConditions, entry) {
				return fmt.Errorf("condition list entry %s survived the policy", entry)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "proof_of_work"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					lists = []string{
						attrs["directive.0.condition_list_id"],
						attrs["directive.1.condition_list_id"],
						attrs["directive.1.skip_condition_list_id"],
					}
					for _, list := range lists {
						if list == "" {
							return fmt.Errorf("a condition list id is missing from %v", attrs)
						}
					}
					condition := f.seedCondition(abpFakeAccountId, "c", "", "true")
					entry = f.seedEntry(lists[0], condition, abpConditionStateActive)
					return nil
				},
			},
		},
	})
}

// --- drift restoration ------------------------------------------------------

// External changes to name, description and a Directive action must be
// detected and reverted.
func TestAbpPolicyDriftIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePolicies)["id"])
					f.mutate(t, abpFakePolicies, id, func(doc map[string]any) {
						doc["name"] = "renamed-out-of-band"
						doc["description"] = "changed out of band"
						abpFakeList(doc["directives"])[0]["action"] = "tarpit"
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "name", "p1"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", "desc"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.0.action", "allow"),
					func(s *terraform.State) error {
						policy := f.only(t, abpFakePolicies)
						directives := abpFakeList(policy["directives"])
						if abpFakeString(policy["name"]) != "p1" ||
							abpFakeString(policy["description"]) != "desc" ||
							abpFakeString(directives[0]["action"]) != "allow" {
							return fmt.Errorf("server state not restored: %+v", policy)
						}
						return nil
					},
				),
			},
		},
	})
}

// A Directive appended out of band must be removed again, and the Condition
// List it brought with it goes with it.
func TestAbpPolicyDriftExtraDirectiveIsRemoved(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow", "block")
	var extra string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePolicies)["id"])
					extra = f.seedConditionList(abpFakeAccountId, "delay", "")
					f.mutate(t, abpFakePolicies, id, func(doc map[string]any) {
						raw, _ := doc["directives"].([]any)
						doc["directives"] = append(raw, map[string]any{
							"action":       "delay",
							"condition_id": extra,
						})
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.#", "2"),
					func(s *terraform.State) error {
						if got := len(abpFakeList(f.only(t, abpFakePolicies)["directives"])); got != 2 {
							return fmt.Errorf("server still has %d directives", got)
						}
						if f.exists(abpFakeConditions, extra) {
							return fmt.Errorf("the removed directive's condition list %s was left behind", extra)
						}
						return nil
					},
				),
			},
		},
	})
}

// A Directive removed out of band must be restored.
func TestAbpPolicyDriftMissingDirectiveIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePolicies)["id"])
					f.mutate(t, abpFakePolicies, id, func(doc map[string]any) {
						raw, _ := doc["directives"].([]any)
						doc["directives"] = raw[:1]
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.1.action", "block"),
					resource.TestCheckResourceAttrSet(abpPolicyTestAddress, "directive.1.condition_list_id"),
					func(s *terraform.State) error {
						directives := abpFakeList(f.only(t, abpFakePolicies)["directives"])
						if len(directives) != 2 || abpFakeString(directives[1]["action"]) != "block" {
							return fmt.Errorf("server state not restored: %+v", directives)
						}
						return nil
					},
				),
			},
		},
	})
}

// Out-of-band deletion must lead to a re-create, not an error.
func TestAbpPolicyDriftDeletedIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow")
	var id string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpPolicyTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakePolicies, id) },
				Config:    config,
				Check: func(s *terraform.State) error {
					recreated := abpFakeString(f.only(t, abpFakePolicies)["id"])
					if recreated == id {
						return fmt.Errorf("the policy was not recreated")
					}
					if got := abpAttrs(t, s, abpPolicyTestAddress)["id"]; got != recreated {
						return fmt.Errorf("state holds %s, the server holds %s", got, recreated)
					}
					return nil
				},
			},
		},
	})
}

// --- import -----------------------------------------------------------------

// Import must produce state that matches the configuration.
func TestAbpPolicyImportStateVerify(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow", "proof_of_work")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpPolicyTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// The first apply after an import must not replace the Policy, and must not
// take its Condition List entries with it.
func TestAbpPolicyImportThenApplyKeepsPolicy(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow")
	seeded := f.seedPolicy("p1", "desc", "allow")
	list := abpFakeString(abpPolicyStoredDirectives(t, f, seeded)[0]["condition_id"])
	condition := f.seedCondition(abpFakeAccountId, "c", "", "true")
	entry := f.seedEntry(list, condition, abpConditionStateActive)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpPolicyTestAddress,
				ImportState:        true,
				ImportStateId:      seeded,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpPolicyTestAddress)["id"]; id != seeded {
						return fmt.Errorf("the imported policy %s was replaced by a new policy %s", seeded, id)
					}
					if !f.exists(abpFakeConditions, entry) {
						return fmt.Errorf("condition list entry %s was deleted by the first apply after import", entry)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Importing a Policy that does not exist fails with a message naming it.
func TestAbpPolicyImportMissing(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("p1", "desc", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:        config,
				ResourceName:  abpPolicyTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Policy ` + abpFakeAbsentId + ` not found`),
			},
		},
	})
}

// --- directive identity ----------------------------------------------------

// Inserting a Directive in front of existing ones must not move the existing
// Directives' Condition Lists onto different actions.
func TestAbpPolicyInsertDirectiveKeepsConditionLists(t *testing.T) {
	f := newAbpFake(t)
	var allowList, blockList, allowEntry string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "block"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					allowList = attrs["directive.0.condition_list_id"]
					blockList = attrs["directive.1.condition_list_id"]
					condition := f.seedCondition(abpFakeAccountId, "c", "", "true")
					allowEntry = f.seedEntry(allowList, condition, abpConditionStateActive)
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p1", "desc", "tarpit", "allow", "block"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					if attrs["directive.0.condition_list_id"] == allowList {
						return fmt.Errorf("the new `tarpit` directive took over the `allow` directive's condition list %s (and with it entry %s)",
							allowList, allowEntry)
					}
					if attrs["directive.1.condition_list_id"] != allowList {
						return fmt.Errorf("the `allow` directive's condition list changed from %s to %s",
							allowList, attrs["directive.1.condition_list_id"])
					}
					if attrs["directive.2.condition_list_id"] != blockList {
						return fmt.Errorf("the `block` directive's condition list changed from %s to %s",
							blockList, attrs["directive.2.condition_list_id"])
					}
					if !f.exists(abpFakeConditions, allowEntry) {
						return fmt.Errorf("condition list entry %s was deleted", allowEntry)
					}
					return nil
				},
			},
		},
	})
}

// Removing the first of two Directives must not move the surviving
// Directive's Condition List, and must not delete its entries.
func TestAbpPolicyRemoveDirectiveKeepsConditionLists(t *testing.T) {
	f := newAbpFake(t)
	var allowList, blockList, blockEntry string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "block"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					allowList = attrs["directive.0.condition_list_id"]
					blockList = attrs["directive.1.condition_list_id"]
					condition := f.seedCondition(abpFakeAccountId, "c", "", "true")
					blockEntry = f.seedEntry(blockList, condition, abpConditionStateActive)
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p1", "desc", "block"),
				Check: func(s *terraform.State) error {
					got := abpAttrs(t, s, abpPolicyTestAddress)["directive.0.condition_list_id"]
					if got != blockList {
						return fmt.Errorf("the surviving `block` directive's condition list changed from %s to %s", blockList, got)
					}
					if !f.exists(abpFakeConditions, blockList) {
						return fmt.Errorf("the `block` directive's condition list %s was deleted", blockList)
					}
					if !f.exists(abpFakeConditions, blockEntry) {
						return fmt.Errorf("condition list entry %s was deleted", blockEntry)
					}
					// The dropped Directive's own list is deleted, by design.
					if f.exists(abpFakeConditions, allowList) {
						return fmt.Errorf("the dropped `allow` directive's condition list %s was left behind", allowList)
					}
					return nil
				},
			},
		},
	})
}

// Reordering Directives must only change their priority, not which Condition
// List belongs to which action.
func TestAbpPolicyReorderDirectivesKeepsConditionLists(t *testing.T) {
	f := newAbpFake(t)
	var allowList, blockList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "block"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					allowList = attrs["directive.0.condition_list_id"]
					blockList = attrs["directive.1.condition_list_id"]
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p1", "desc", "block", "allow"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					if attrs["directive.0.condition_list_id"] != blockList || attrs["directive.1.condition_list_id"] != allowList {
						return fmt.Errorf("condition lists did not follow their actions: block=%s (was %s), allow=%s (was %s)",
							attrs["directive.0.condition_list_id"], blockList,
							attrs["directive.1.condition_list_id"], allowList)
					}
					if len(f.deletedIds()) != 0 {
						return fmt.Errorf("reordering deleted %v", f.deletedIds())
					}
					return nil
				},
			},
		},
	})
}

// Replacing a Directive's action in place must not hand the old action's
// Condition List to the new action: the Conditions attached to it were chosen
// for the action that is gone.
func TestAbpPolicyChangedActionGetsAFreshConditionList(t *testing.T) {
	f := newAbpFake(t)
	var allowList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow"),
				Check: func(s *terraform.State) error {
					allowList = abpAttrs(t, s, abpPolicyTestAddress)["directive.0.condition_list_id"]
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p1", "desc", "block"),
				Check: func(s *terraform.State) error {
					got := abpAttrs(t, s, abpPolicyTestAddress)["directive.0.condition_list_id"]
					if got == allowList {
						return fmt.Errorf("the `allow` directive's condition list %s was reused for the new `block` directive", got)
					}
					return nil
				},
			},
		},
	})
}

// Inserting a Directive in front of a proof_of_work Directive shifts the list
// positions, which must not shift the skip Condition List onto a non
// proof_of_work action — the API rejects that outright.
func TestAbpPolicyInsertBeforeProofOfWorkDirective(t *testing.T) {
	f := newAbpFake(t)
	var powList, skipList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "proof_of_work", "allow"),
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPolicyTestAddress)
					powList = attrs["directive.0.condition_list_id"]
					skipList = attrs["directive.0.skip_condition_list_id"]
					if skipList == "" {
						return fmt.Errorf("the proof_of_work directive got no skip condition list")
					}
					return nil
				},
			},
			{
				Config: abpPolicyTestConfig("p1", "desc", "block", "proof_of_work", "allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.0.skip_condition_list_id", ""),
					func(s *terraform.State) error {
						attrs := abpAttrs(t, s, abpPolicyTestAddress)
						if attrs["directive.1.condition_list_id"] != powList {
							return fmt.Errorf("the proof_of_work directive's condition list changed from %s to %s",
								powList, attrs["directive.1.condition_list_id"])
						}
						if attrs["directive.1.skip_condition_list_id"] != skipList {
							return fmt.Errorf("the proof_of_work directive's skip condition list changed from %s to %s",
								skipList, attrs["directive.1.skip_condition_list_id"])
						}
						return nil
					},
				),
			},
		},
	})
}

// On the initial create the whole resource is unknown, so a consumer of a
// Directive's condition_list_id works. This scopes the test below.
func TestAbpPolicyCreateWithConditionListConsumer(t *testing.T) {
	f := newAbpFake(t)
	const consumer = `
resource "terraform_data" "consumer" {
  input = incapsula_abp_policy.test.directive[1].condition_list_id
}
`
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "block") + consumer,
				Check: func(s *terraform.State) error {
					want := abpAttrs(t, s, abpPolicyTestAddress)["directive.1.condition_list_id"]
					got := abpAttrs(t, s, "terraform_data.consumer")["output"]
					if got != want || got == "" {
						return fmt.Errorf("consumer got %q, policy has %q", got, want)
					}
					return nil
				},
			},
		},
	})
}

// Known limitation, not an intention: a Condition List belonging to a Directive
// added by an *update* cannot be consumed in the same apply. SDKv2 keeps a
// computed attribute nested in a TypeList attached to its list index, so the
// new Directive's id is planned at its old position instead of as unknown and
// the apply dies. `incapsula_abp_directive` is the supported way to reach these
// ids; see the schema comment in resource_abp_policy.go for why there is no
// escape hatch.
func TestAbpPolicyNewDirectiveConditionListCannotBeConsumed(t *testing.T) {
	f := newAbpFake(t)
	const consumer = `
resource "terraform_data" "consumer" {
  input = incapsula_abp_policy.test.directive[1].condition_list_id
}
`
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPolicyTestConfig("p1", "desc", "allow")},
			{
				Config:      abpPolicyTestConfig("p1", "desc", "allow", "block") + consumer,
				ExpectError: regexp.MustCompile(`(?s)Provider produced inconsistent final plan`),
			},
		},
	})
}

// The documented alternative: the data source re-reads the ids from the API, so
// a Directive added by an update is immediately usable.
func TestAbpPolicyNewDirectiveConditionListViaDataSource(t *testing.T) {
	f := newAbpFake(t)
	consumer := `
data "incapsula_abp_directive" "block" {
  policy_id = incapsula_abp_policy.test.id
  action    = "block"
}

resource "terraform_data" "consumer" {
  input = data.incapsula_abp_directive.block.condition_list_id
}
`
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPolicyTestConfig("p1", "desc", "allow")},
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow", "block") + consumer,
				Check: func(s *terraform.State) error {
					want := abpAttrs(t, s, abpPolicyTestAddress)["directive.1.condition_list_id"]
					got := abpAttrs(t, s, "terraform_data.consumer")["output"]
					if got != want || got == "" {
						return fmt.Errorf("consumer got %q, policy has %q", got, want)
					}
					return nil
				},
			},
		},
	})
}

// --- unmodelled fields ------------------------------------------------------

// An update is a read-modify-write, so a field the provider does not model
// survives an unrelated change instead of being cleared by the full-replace
// endpoint.
func TestAbpPolicyUnmodelledFieldSurvivesUpdate(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPolicyTestConfig("p1", "desc", "allow")},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePolicies)["id"])
					f.mutate(t, abpFakePolicies, id, func(doc map[string]any) {
						doc["invented_by_the_service"] = "keep me"
					})
				},
				Config: abpPolicyTestConfig("p2", "desc", "allow"),
				Check: func(s *terraform.State) error {
					policy := f.only(t, abpFakePolicies)
					if abpFakeString(policy["invented_by_the_service"]) != "keep me" {
						return fmt.Errorf("the rename cleared a field the provider does not model: %+v", policy)
					}
					return nil
				},
			},
		},
	})
}

// --- validation and plan-time errors ---------------------------------------

// Two Directives with the same action must be rejected: Directive identity
// comes from the action, and `incapsula_abp_directive` cannot address an
// action that occurs twice.
func TestAbpPolicyDuplicateActions(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyTestConfig("p1", "desc", "block", "block"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`more than one directive with action "block"`),
			},
		},
	})
}

// An action that is only known after apply must not be mistaken for a missing
// or duplicate one at plan time.
func TestAbpPolicyUnknownActionIsPlannable(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "terraform_data" "action" {
  input = "block"
}

resource "incapsula_abp_policy" "test" {
  account_id  = %q
  name        = "p1"
  description = "desc"

  directive {
    action = "allow"
  }

  directive {
    action = terraform_data.action.output
  }
}
`, abpFakeAccountId)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "directive.1.action", "block"),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A duplicate hidden behind an unknown action at plan time is still rejected,
// because Terraform re-plans the instance with the resolved configuration
// before applying it.
func TestAbpPolicyUnknownDuplicateActionIsRejectedOnApply(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "terraform_data" "action" {
  input = "allow"
}

resource "incapsula_abp_policy" "test" {
  account_id  = %q
  name        = "p1"
  description = "desc"

  directive {
    action = "allow"
  }

  directive {
    action = terraform_data.action.output
  }
}
`, abpFakeAccountId)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`more than one directive with action "allow"`),
			},
		},
	})
}

// `description` is optional, so it must be omittable. An omitted description is
// left out of the request (the API then stores none) and read back as the empty
// string, and removing it from a configuration that had one clears it.
func TestAbpPolicyDescriptionCanBeOmitted(t *testing.T) {
	f := newAbpFake(t)
	withoutDescription := fmt.Sprintf(`
resource "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "p1"

  directive {
    action = "allow"
  }
}
`, abpFakeAccountId)
	checkServerDescription := func(want any) resource.TestCheckFunc {
		return func(*terraform.State) error {
			got, present := f.only(t, abpFakePolicies)["description"]
			switch {
			case want == nil && !present:
				return nil
			case want == nil:
				return fmt.Errorf("expected no description on the server, got %q", got)
			case !present:
				return fmt.Errorf("expected description %q on the server, got none", want)
			case want != got:
				return fmt.Errorf("expected description %q on the server, got %q", want, got)
			}
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: withoutDescription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", ""),
					checkServerDescription(nil),
				),
			},
			{Config: withoutDescription, PlanOnly: true},
			{
				Config: abpPolicyTestConfig("p1", "desc", "allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", "desc"),
					checkServerDescription("desc"),
				),
			},
			{
				Config: withoutDescription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyTestAddress, "description", ""),
					checkServerDescription(nil),
				),
			},
			{Config: withoutDescription, PlanOnly: true},
		},
	})
}

// The API trims `name`, so a padded name could never converge. It is rejected
// at plan time instead, as is a name that is only whitespace (which the API
// refuses once trimmed).
func TestAbpPolicyPaddedNameIsRejected(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyTestConfig("  p1  ", "desc", "allow"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{
				Config:      abpPolicyTestConfig("   ", "desc", "allow"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not be empty or consist only of whitespace`),
			},
			{Config: abpPolicyTestConfig("p1", "desc", "allow")},
			{Config: abpPolicyTestConfig("p1", "desc", "allow"), PlanOnly: true},
		},
	})
}

// The arguments the schema rejects before anything is sent.
func TestAbpPolicyInvalidConfigurations(t *testing.T) {
	f := newAbpFake(t)
	providers := f.providers(t)
	for name, step := range map[string]resource.TestStep{
		"not a uuid": {
			Config:      strings.Replace(abpPolicyTestConfig("p1", "desc", "allow"), abpFakeAccountId, "not-a-uuid", 1),
			ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
		},
		"action not snake case": {
			Config:      abpPolicyTestConfig("p1", "desc", "Block"),
			ExpectError: regexp.MustCompile(`action must be snake_case`),
		},
		"no directive": {
			Config:      abpPolicyTestConfig("p1", "desc"),
			ExpectError: regexp.MustCompile(`(?s)[Ii]nsufficient directive blocks`),
		},
		"name too long": {
			Config:      abpPolicyTestConfig(strings.Repeat("n", 101), "desc", "allow"),
			ExpectError: regexp.MustCompile(`expected length of name to be in the range`),
		},
	} {
		t.Run(name, func(t *testing.T) {
			step.PlanOnly = true
			resource.Test(t, resource.TestCase{
				ProviderFactories: providers,
				Steps:             []resource.TestStep{step},
			})
		})
	}
}

// Known limitation, not an intention: `proof_of_work_configuration_id` on a non
// proof_of_work Directive is rejected only by the API, so the configuration
// plans cleanly and the apply fails.
func TestAbpPolicyPowConfigOnNonPowDirectiveIsRejectedByTheApi(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "incapsula_abp_policy" "test" {
  account_id  = %q
  name        = "p1"
  description = "desc"

  directive {
    action                         = "allow"
    proof_of_work_configuration_id = "eeeeeeee-0000-4000-8000-000000000001"
  }
}
`, abpFakeAccountId)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`proof_of_work_configuration_id is only allowed for proof_of_work action`),
			},
		},
	})
}
