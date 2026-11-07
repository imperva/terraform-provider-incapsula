package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for the incapsula_abp_conditions data source. It lists the literal
// Conditions of an account ordered by name, and its `managed` argument
// selects *either* the managed or the account-owned ones — the opposite of
// what the same argument means on the singular incapsula_abp_condition.

const abpConditionsDataAddress = "data.incapsula_abp_conditions.test"

func abpConditionsDataConfig(accountId string) string {
	return fmt.Sprintf(`
data "incapsula_abp_conditions" "test" {
  account_id = %q
}
`, accountId)
}

func abpConditionsDataConfigManaged(accountId string, managed bool) string {
	return fmt.Sprintf(`
data "incapsula_abp_conditions" "test" {
  account_id = %q
  managed    = %t
}
`, accountId, managed)
}

// The list holds only the account's literal Conditions, sorted by name.
// Condition Lists, Condition List Entries, managed Conditions and other
// accounts' Conditions are all left out.
func TestAbpConditionsDataSourceListsAccountOwnedSortedByName(t *testing.T) {
	f := newAbpFake(t)
	zulu := f.seedCondition(abpFakeAccountId, "zulu", "d", "(true)")
	alpha := f.seedCondition(abpFakeAccountId, "alpha", "d", "(false)")
	list := f.seedConditionList(abpFakeAccountId, "a list", "d")
	f.seedEntry(list, zulu, abpConditionStateActive)
	abpConditionSeedManaged(f, "Monitoring Tools", "d")
	f.seedCondition(abpFakeOtherAccountId, "elsewhere", "d", "(true)")
	// A Policy brings Condition Lists of its own, which are not Conditions.
	f.seedPolicy("a policy", "d", "block", "allow")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionsDataConfig(abpFakeAccountId),
				Check: resource.ComposeAggregateTestCheckFunc(
					// The data source id is the account it listed.
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.#", "2"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.id", alpha),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.name", "alpha"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.1.id", zulu),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.1.name", "zulu"),
				),
			},
		},
	})
}

func TestAbpConditionsDataSourceElementAttributes(t *testing.T) {
	f := newAbpFake(t)
	const code = "(matches headers.user_agent \"okhttp\")"
	id := f.seedCondition(abpFakeAccountId, "okhttp", "matches okhttp", code)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionsDataConfig(abpFakeAccountId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.id", id),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.managed", "false"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.description", "matches okhttp"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.code", abpFakeFormatCode(code)),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.template", "custom"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.last_change_by", "fake"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.recommended_directives.#", "0"),
					resource.TestCheckResourceAttrSet(abpConditionsDataAddress, "conditions.0.created_at"),
					resource.TestCheckResourceAttrSet(abpConditionsDataAddress, "conditions.0.modified_at"),
				),
			},
		},
	})
}

// With `managed = true` only the managed Conditions are listed, including the
// deprecated ones the UI hides — there is no attribute to tell those apart.
func TestAbpConditionsDataSourceManagedTrueListsOnlyManaged(t *testing.T) {
	f := newAbpFake(t)
	f.seedCondition(abpFakeAccountId, "ours", "d", "(true)")
	deprecated := abpConditionSeedManaged(f, "AI Tools - AI Assistant - Deprecated", "d")
	monitoring := abpConditionSeedManaged(f, "Monitoring Tools", "d", "allow")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionsDataConfigManaged(abpFakeAccountId, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.#", "2"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.id", deprecated),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.managed", "true"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.account_id", ""),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.1.id", monitoring),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.1.recommended_directives.0", "allow"),
				),
			},
			// `false` is the default and the complement of `true`.
			{
				Config: abpConditionsDataConfigManaged(abpFakeAccountId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.#", "1"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.name", "ours"),
				),
			},
		},
	})
}

func TestAbpConditionsDataSourceEmptyAccount(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionsDataConfig(abpFakeAccountId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.#", "0"),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "id", abpFakeAccountId),
				),
			},
		},
	})
}

// Conditions sharing a name are ordered by id, so the list is stable even
// then.
func TestAbpConditionsDataSourceOrdersEqualNamesById(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedCondition(abpFakeAccountId, "twice", "first", "(true)")
	second := f.seedCondition(abpFakeAccountId, "twice", "second", "(false)")
	config := abpConditionsDataConfig(abpFakeAccountId)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.0.id", first),
					resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.1.id", second),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Conditions added and removed out of band show up on the next read.
func TestAbpConditionsDataSourceFollowsOutOfBandChanges(t *testing.T) {
	f := newAbpFake(t)
	kept := f.seedCondition(abpFakeAccountId, "kept", "d", "(true)")
	removed := f.seedCondition(abpFakeAccountId, "removed", "d", "(true)")
	config := abpConditionsDataConfig(abpFakeAccountId)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpConditionsDataAddress, "conditions.#", "2"),
			},
			{
				PreConfig: func() {
					f.drop(abpFakeConditions, removed)
					f.seedCondition(abpFakeAccountId, "added", "d", "(true)")
				},
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpConditionsDataAddress)
					if attrs["conditions.#"] != "2" {
						return fmt.Errorf("the list holds %s Conditions", attrs["conditions.#"])
					}
					if attrs["conditions.0.name"] != "added" || attrs["conditions.1.id"] != kept {
						return fmt.Errorf("the list is %q, %q", attrs["conditions.0.name"], attrs["conditions.1.name"])
					}
					return nil
				},
			},
		},
	})
}

func TestAbpConditionsDataSourceValidation(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionsDataConfig("not-a-uuid"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
		},
	})
}
