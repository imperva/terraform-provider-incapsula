package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The incapsula_abp_policy data source: lookup by id, by name within an
// Account, and the Account's global Policy.

const abpPolicyDataAddress = "data.incapsula_abp_policy.test"

// abpPolicySeedRichPolicy seeds a Policy that exercises every attribute the
// data sources report, including the two ids only a proof_of_work Directive
// carries.
func abpPolicySeedRichPolicy(f *abpFake, name, description string) (policyId string, lists map[string]string) {
	lists = map[string]string{
		"allow":             f.seedConditionList(abpFakeAccountId, "allow", ""),
		"proof_of_work":     f.seedConditionList(abpFakeAccountId, "proof_of_work", ""),
		"proof_of_work_pow": f.seedPowConfig("pow-for-"+name, 42),
		"skip":              f.seedConditionList(abpFakeAccountId, "proof_of_work_skip", ""),
	}
	policyId = f.abpPolicySeedWithDirectives(name, description,
		map[string]any{"action": "allow", "condition_id": lists["allow"]},
		map[string]any{
			"action":                         "proof_of_work",
			"condition_id":                   lists["proof_of_work"],
			"skip_condition_id":              lists["skip"],
			"proof_of_work_configuration_id": lists["proof_of_work_pow"],
		},
	)
	return policyId, lists
}

// --- lookup by name ---------------------------------------------------------

func TestAbpPolicyDataSourceByName(t *testing.T) {
	f := newAbpFake(t)
	policyId, lists := abpPolicySeedRichPolicy(f, "lookup-me", "a description")
	config := fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "lookup-me"
}
`, abpFakeAccountId)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "id", policyId),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "description", "a description"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.action", "allow"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.condition_list_id", lists["allow"]),
					// Only a proof_of_work Directive carries these two.
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.skip_condition_list_id", ""),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.proof_of_work_configuration_id", ""),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.1.action", "proof_of_work"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.1.condition_list_id", lists["proof_of_work"]),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.1.skip_condition_list_id", lists["skip"]),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.1.proof_of_work_configuration_id", lists["proof_of_work_pow"]),
					resource.TestCheckResourceAttrSet(abpPolicyDataAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpPolicyDataAddress, "modified_at"),
				),
			},
		},
	})
}

// Names are matched exactly and are case-sensitive.
func TestAbpPolicyDataSourceNameIsCaseSensitive(t *testing.T) {
	f := newAbpFake(t)
	f.seedPolicy("Lookup-Me", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "lookup-me"
}
`, abpFakeAccountId),
				ExpectError: regexp.MustCompile(`no ABP Policy named "lookup-me" found in account`),
			},
		},
	})
}

// A Policy of another Account is not visible.
func TestAbpPolicyDataSourceNameIsScopedToTheAccount(t *testing.T) {
	f := newAbpFake(t)
	f.seedPolicy("elsewhere", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "elsewhere"
}
`, abpFakeOtherAccountId),
				ExpectError: regexp.MustCompile(`no ABP Policy named "elsewhere" found in account ` + abpFakeOtherAccountId),
			},
		},
	})
}

// Names are not unique, so an ambiguous lookup has to fail rather than pick one.
func TestAbpPolicyDataSourceDuplicateNameIsRejected(t *testing.T) {
	f := newAbpFake(t)
	f.seedPolicy("twice", "first", "allow")
	f.seedPolicy("twice", "second", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "twice"
}
`, abpFakeAccountId),
				ExpectError: regexp.MustCompile(`multiple ABP Policies named "twice" found in account`),
			},
		},
	})
}

// --- lookup by id -----------------------------------------------------------

// A lookup by id returns the Policy's contents. Known limitation, not an
// intention: the search arguments `name` and `account_id` are dropped from the
// result, so neither reports the Policy's own value.
func TestAbpPolicyDataSourceById(t *testing.T) {
	f := newAbpFake(t)
	policyId, lists := abpPolicySeedRichPolicy(f, "by-id", "a description")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  id = %q
}
`, policyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "id", policyId),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "description", "a description"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.1.condition_list_id", lists["proof_of_work"]),
					abpPolicyNotReported(abpPolicyDataAddress, "name", "account_id"),
				),
			},
		},
	})
}

// `account_id` is accepted but ignored by a lookup by id, so it cannot be used
// to scope one.
func TestAbpPolicyDataSourceByIdIgnoresAccountId(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("by-id", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  id         = %q
}
`, abpFakeOtherAccountId, policyId),
				Check: resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.action", "allow"),
			},
		},
	})
}

func TestAbpPolicyDataSourceUnknownIdIsRejected(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  id = %q
}
`, abpFakeAbsentId),
				ExpectError: regexp.MustCompile(`no ABP Policy with id "` + abpFakeAbsentId + `" found`),
			},
		},
	})
}

// --- the Account's global Policy -------------------------------------------

func TestAbpPolicyDataSourceAccountGlobal(t *testing.T) {
	f := newAbpFake(t)
	globalId := f.abpPolicySeedGlobalPolicy("global", "allow", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id     = %q
  account_global = true
}
`, abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "id", globalId),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.0.action", "allow"),
					resource.TestCheckResourceAttrSet(abpPolicyDataAddress, "directive.0.condition_list_id"),
					// The global Policy has no name of its own.
					abpPolicyNotReported(abpPolicyDataAddress, "name"),
				),
			},
		},
	})
}

func TestAbpPolicyDataSourceAccountGlobalMissing(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id     = %q
  account_global = true
}
`, abpFakeAccountId),
				ExpectError: regexp.MustCompile(`no global ABP Policy found for account ` + abpFakeAccountId),
			},
		},
	})
}

// --- argument combinations --------------------------------------------------

func TestAbpPolicyDataSourceArgumentCombinations(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("p1", "", "allow")
	providers := f.providers(t)

	for name, step := range map[string]resource.TestStep{
		"nothing to search by": {
			Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
}
`, abpFakeAccountId),
			ExpectError: regexp.MustCompile(`(?s)one of .*account_global.*must be\s+specified`),
		},
		"name and id together": {
			Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "p1"
  id         = %q
}
`, abpFakeAccountId, policyId),
			ExpectError: regexp.MustCompile(`(?s)only one of`),
		},
		// Known limitation, not an intention: SDKv2 counts an attribute as
		// specified whenever it is present, so stating the default explicitly
		// breaks the data source and a parameterised lookup is impossible.
		"account_global written as false": {
			Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id     = %q
  name           = "p1"
  account_global = false
}
`, abpFakeAccountId),
			ExpectError: regexp.MustCompile(`(?s)only one of`),
		},
		"name without an account": {
			Config: `
data "incapsula_abp_policy" "test" {
  name = "p1"
}
`,
			ExpectError: regexp.MustCompile(`(?s)all of .*account_id.*must be\s+specified`),
		},
		"account id is not a uuid": {
			Config: `
data "incapsula_abp_policy" "test" {
  account_id = "not-a-uuid"
  name       = "p1"
}
`,
			ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
		},
		"name is only whitespace": {
			Config: fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "   "
}
`, abpFakeAccountId),
			ExpectError: regexp.MustCompile(`(?s)expected "name" to not be an empty string or whitespace`),
		},
	} {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: providers,
				Steps:             []resource.TestStep{step},
			})
		})
	}
}

// --- following a resource --------------------------------------------------

// The data source re-reads the Policy from the API, so it reports what an
// update left behind rather than what state happened to hold.
func TestAbpPolicyDataSourceFollowsTheResource(t *testing.T) {
	f := newAbpFake(t)
	config := func(actions ...string) string {
		return abpPolicyTestConfig("tracked", "desc", actions...) + fmt.Sprintf(`
data "incapsula_abp_policy" "test" {
  account_id = %q
  name       = incapsula_abp_policy.test.name
}
`, abpFakeAccountId)
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config("allow"),
				Check:  resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.#", "1"),
			},
			{
				Config: config("allow", "block"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDataAddress, "directive.#", "2"),
					func(s *terraform.State) error {
						data := abpAttrs(t, s, abpPolicyDataAddress)
						managed := abpAttrs(t, s, abpPolicyTestAddress)
						for _, key := range []string{"directive.0.condition_list_id", "directive.1.condition_list_id"} {
							if data[key] != managed[key] {
								return fmt.Errorf("%s: the data source reports %q, the resource %q", key, data[key], managed[key])
							}
						}
						return nil
					},
				),
			},
		},
	})
}
