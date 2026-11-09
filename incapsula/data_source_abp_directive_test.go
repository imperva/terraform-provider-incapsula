package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The incapsula_abp_directive data source, which addresses a Directive by its
// action — the only stable handle on one — either within a Policy or within the
// Account's global Policy.

const abpPolicyDirectiveDataAddress = "data.incapsula_abp_directive.test"

func TestAbpPolicyDirectiveDataSourceByAction(t *testing.T) {
	f := newAbpFake(t)
	policyId, lists := abpPolicySeedRichPolicy(f, "p1", "desc")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "allow"
}
`, policyId),
				Check: resource.ComposeTestCheckFunc(
					// Directives have no id of their own, so the data source
					// synthesizes one.
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "id", policyId+":allow"),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "policy_id", policyId),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "action", "allow"),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "condition_list_id", lists["allow"]),
					// Absent ids come back as the empty string, never null.
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "skip_condition_list_id", ""),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "proof_of_work_configuration_id", ""),
				),
			},
		},
	})
}

// A proof_of_work Directive is the only one that carries a skip Condition List
// and a Proof Of Work Configuration.
func TestAbpPolicyDirectiveDataSourceProofOfWorkAction(t *testing.T) {
	f := newAbpFake(t)
	policyId, lists := abpPolicySeedRichPolicy(f, "p1", "desc")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "proof_of_work"
}
`, policyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "condition_list_id", lists["proof_of_work"]),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "skip_condition_list_id", lists["skip"]),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "proof_of_work_configuration_id", lists["proof_of_work_pow"]),
				),
			},
		},
	})
}

// A Directive of the Account's global Policy is addressed without a Policy id.
// Known limitation, not an intention: `policy_id` is then not reported, so the
// Policy the Directive belongs to cannot be read off the result.
func TestAbpPolicyDirectiveDataSourceAccountGlobalPolicy(t *testing.T) {
	f := newAbpFake(t)
	globalId := f.abpPolicySeedGlobalPolicy("global", "allow", "block")
	wantList := abpFakeString(abpPolicyStoredDirectives(t, f, globalId)[1]["condition_id"])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  account_id            = %q
  account_global_policy = true
  action                = "block"
}
`, abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "id", "global:"+abpFakeAccountId+":block"),
					resource.TestCheckResourceAttr(abpPolicyDirectiveDataAddress, "condition_list_id", wantList),
					abpPolicyNotReported(abpPolicyDirectiveDataAddress, "policy_id"),
				),
			},
		},
	})
}

func TestAbpPolicyDirectiveDataSourceAccountGlobalPolicyMissing(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  account_id            = %q
  account_global_policy = true
  action                = "allow"
}
`, abpFakeAccountId),
				ExpectError: regexp.MustCompile(`ABP account global policy not found for account ` + abpFakeAccountId),
			},
		},
	})
}

// `account_id` is accepted but ignored when a Policy id is given.
func TestAbpPolicyDirectiveDataSourceIgnoresAccountId(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("p1", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  account_id = %q
  policy_id  = %q
  action     = "allow"
}
`, abpFakeOtherAccountId, policyId),
				Check: resource.TestCheckResourceAttrSet(abpPolicyDirectiveDataAddress, "condition_list_id"),
			},
		},
	})
}

// An action that occurs twice in a Policy — which only a change made outside
// Terraform can produce — is unaddressable, which is why the resource rejects
// duplicate actions in the first place.
func TestAbpPolicyDirectiveDataSourceDuplicateActionIsRejected(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("p1", "", "block", "block")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "block"
}
`, policyId),
				ExpectError: regexp.MustCompile(`multiple directives with action "block" found in ABP Policy ` + policyId),
			},
		},
	})
}

func TestAbpPolicyDirectiveDataSourceUnknownAction(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("p1", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "tarpit"
}
`, policyId),
				ExpectError: regexp.MustCompile(`no directive with action "tarpit" found in ABP Policy ` + policyId),
			},
		},
	})
}

func TestAbpPolicyDirectiveDataSourceUnknownPolicy(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "allow"
}
`, abpFakeAbsentId),
				ExpectError: regexp.MustCompile(`ABP Policy ` + abpFakeAbsentId + ` not found`),
			},
		},
	})
}

func TestAbpPolicyDirectiveDataSourceArgumentCombinations(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("p1", "", "allow")
	providers := f.providers(t)

	for name, step := range map[string]resource.TestStep{
		"no policy to search in": {
			Config: `
data "incapsula_abp_directive" "test" {
  action = "allow"
}
`,
			ExpectError: regexp.MustCompile(`(?s)one of .*account_global_policy.*must be\s+specified`),
		},
		"both a policy and the global one": {
			Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  account_id            = %q
  policy_id             = %q
  account_global_policy = true
  action                = "allow"
}
`, abpFakeAccountId, policyId),
			ExpectError: regexp.MustCompile(`(?s)only one of`),
		},
		// Known limitation, not an intention: SDKv2 counts an attribute as
		// specified whenever it is present, so the default cannot be written
		// out and a parameterised lookup is impossible.
		"account_global_policy written as false": {
			Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id             = %q
  account_global_policy = false
  action                = "allow"
}
`, policyId),
			ExpectError: regexp.MustCompile(`(?s)only one of`),
		},
		"global policy without an account": {
			Config: `
data "incapsula_abp_directive" "test" {
  account_global_policy = true
  action                = "allow"
}
`,
			ExpectError: regexp.MustCompile(`(?s)all of .*account_id.*must be\s+specified`),
		},
		"policy id is not a uuid": {
			Config: `
data "incapsula_abp_directive" "test" {
  policy_id = "not-a-uuid"
  action    = "allow"
}
`,
			ExpectError: regexp.MustCompile(`expected "policy_id" to be a valid UUID`),
		},
		"action is only whitespace": {
			Config: fmt.Sprintf(`
data "incapsula_abp_directive" "test" {
  policy_id = %q
  action    = "  "
}
`, policyId),
			ExpectError: regexp.MustCompile(`expected "action" to not be an empty string or whitespace`),
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

// The supported way to attach Conditions to a Directive: the data source
// re-reads the Condition List from the API, so it stays on the right Directive
// across an insert that moves every list index.
func TestAbpPolicyDirectiveDataSourceSurvivesADirectiveInsert(t *testing.T) {
	f := newAbpFake(t)
	config := func(actions ...string) string {
		return abpPolicyTestConfig("p1", "desc", actions...) + `
data "incapsula_abp_directive" "test" {
  policy_id = incapsula_abp_policy.test.id
  action    = "block"
}

resource "incapsula_abp_condition" "c" {
  account_id  = "` + abpFakeAccountId + `"
  name        = "c"
  description = "d"
  code        = "true"
}

resource "incapsula_abp_condition_list_entry" "e" {
  account_id               = "` + abpFakeAccountId + `"
  parent_condition_list_id = data.incapsula_abp_directive.test.condition_list_id
  condition_id             = incapsula_abp_condition.c.id
  state                    = "active"
}
`
	}
	var blockList string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config("allow", "block"),
				Check: func(s *terraform.State) error {
					blockList = abpAttrs(t, s, abpPolicyDirectiveDataAddress)["condition_list_id"]
					if want := abpAttrs(t, s, abpPolicyTestAddress)["directive.1.condition_list_id"]; blockList != want {
						return fmt.Errorf("the data source reports %q, the resource %q", blockList, want)
					}
					return nil
				},
			},
			{
				// `block` moves from index 1 to index 2.
				Config: config("tarpit", "allow", "block"),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPolicyDirectiveDataAddress)["condition_list_id"]; got != blockList {
						return fmt.Errorf("the `block` directive's condition list changed from %s to %s", blockList, got)
					}
					// The entry itself is replaced rather than kept: while the
					// Policy has changes the data source is read at apply time,
					// so the entry's ForceNew `parent_condition_list_id` is
					// unknown while the plan is made. It lands on the same
					// Condition List, which is what matters here.
					entry := abpAttrs(t, s, "incapsula_abp_condition_list_entry.e")
					if entry["parent_condition_list_id"] != blockList {
						return fmt.Errorf("the entry moved to condition list %s", entry["parent_condition_list_id"])
					}
					return nil
				},
			},
		},
	})
}
