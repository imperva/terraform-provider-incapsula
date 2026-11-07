package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The incapsula_abp_policies data source, which lists an Account's Policies.

const abpPolicyPoliciesDataAddress = "data.incapsula_abp_policies.test"

func abpPolicyPoliciesDataConfig(accountId string) string {
	return fmt.Sprintf(`
data "incapsula_abp_policies" "test" {
  account_id = %q
}
`, accountId)
}

// The list is ordered by name, and it contains the `"<site> Default"` Policy
// the API creates for every Site next to the hand-written ones, with nothing
// to tell them apart. The Account's global Policy is not part of it.
func TestAbpPolicyPoliciesDataSourceListsTheAccount(t *testing.T) {
	f := newAbpFake(t)
	f.seedSite("zsite")
	bravo := f.seedPolicy("bravo", "second", "block")
	alpha := f.seedPolicy("alpha", "first", "allow", "proof_of_work")
	global := f.abpPolicySeedGlobalPolicy("global", "allow")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPoliciesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.#", "3"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.name", "alpha"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.id", alpha),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.description", "first"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.directive.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.directive.1.action", "proof_of_work"),
					resource.TestCheckResourceAttrSet(abpPolicyPoliciesDataAddress, "policies.0.directive.1.skip_condition_list_id"),
					resource.TestCheckResourceAttrSet(abpPolicyPoliciesDataAddress, "policies.0.created_at"),
					resource.TestCheckResourceAttrSet(abpPolicyPoliciesDataAddress, "policies.0.modified_at"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.1.name", "bravo"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.1.id", bravo),
					// The Policy the backend created along with the Site.
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.2.name", "zsite Default"),
					func(s *terraform.State) error {
						for _, attribute := range abpAttrs(t, s, abpPolicyPoliciesDataAddress) {
							if attribute == global {
								return fmt.Errorf("the account's global policy %s is part of the list", global)
							}
						}
						return nil
					},
				),
			},
		},
	})
}

// An Account with no Policies yields an empty list rather than an error.
func TestAbpPolicyPoliciesDataSourceEmptyAccount(t *testing.T) {
	f := newAbpFake(t)
	f.seedPolicy("elsewhere", "", "allow")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPoliciesDataConfig(abpFakeOtherAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "id", abpFakeOtherAccountId),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.#", "0"),
				),
			},
		},
	})
}

// Policies sharing a name are ordered by id, so the list is stable even though
// the API defines no order.
func TestAbpPolicyPoliciesDataSourceSameNameIsOrderedById(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedPolicy("same", "first", "allow")
	second := f.seedPolicy("same", "second", "block")
	if first >= second {
		t.Fatalf("the fake handed out ids out of order: %s then %s", first, second)
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPoliciesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.id", first),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.1.id", second),
				),
			},
		},
	})
}

// The list follows the Policies Terraform manages: one added by a later apply
// shows up in it. `depends_on` is what defers the read until the Policy exists.
func TestAbpPolicyPoliciesDataSourceFollowsTheAccount(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyTestConfig("managed", "desc", "allow") + fmt.Sprintf(`
data "incapsula_abp_policies" "test" {
  account_id = %q
  depends_on = [incapsula_abp_policy.test]
}
`, abpFakeAccountId)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.#", "1"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.name", "managed"),
					resource.TestCheckResourceAttrPair(
						abpPolicyPoliciesDataAddress, "policies.0.id",
						abpPolicyTestAddress, "id"),
				),
			},
			{
				PreConfig: func() { f.seedPolicy("added-elsewhere", "", "block") },
				Config:    config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.#", "2"),
					resource.TestCheckResourceAttr(abpPolicyPoliciesDataAddress, "policies.0.name", "added-elsewhere"),
				),
			},
		},
	})
}

func TestAbpPolicyPoliciesDataSourceInvalidAccountId(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyPoliciesDataConfig("not-a-uuid"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
		},
	})
}

func TestAbpPolicyPoliciesDataSourceAccountIdIsRequired(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: `
data "incapsula_abp_policies" "test" {
}
`,
				ExpectError: regexp.MustCompile(`(?s)"account_id" is required`),
			},
		},
	})
}
