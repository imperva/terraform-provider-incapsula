package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for the incapsula_abp_condition_list data source. It looks a Condition
// List up by name within an account; only the list variant matches.

const abpConditionListDataAddress = "data.incapsula_abp_condition_list.test"

func abpConditionListDataConfig(accountId, name string) string {
	return fmt.Sprintf(`
data "incapsula_abp_condition_list" "test" {
  account_id = %q
  name       = %q
}
`, accountId, name)
}

func TestAbpConditionListDataSourceByName(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList(abpFakeAccountId, "Sample condition list", "reusable")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListDataConfig(abpFakeAccountId, "Sample condition list"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionListDataAddress, "id", id),
					resource.TestCheckResourceAttr(abpConditionListDataAddress, "description", "reusable"),
					// The API only ever exposes managed Conditions as literal
					// Conditions, so this is always false.
					resource.TestCheckResourceAttr(abpConditionListDataAddress, "managed", "false"),
					resource.TestCheckResourceAttrSet(abpConditionListDataAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpConditionListDataAddress, "modified_at"),
				),
			},
		},
	})
}

// Literal Conditions and Condition List Entries of the same name are skipped.
func TestAbpConditionListDataSourceSkipsOtherVariants(t *testing.T) {
	f := newAbpFake(t)
	literal := f.seedCondition(abpFakeAccountId, "shared name", "a literal", "(true)")
	list := f.seedConditionList(abpFakeAccountId, "shared name", "a list")
	f.seedEntry(list, literal, abpConditionStateActive)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListDataConfig(abpFakeAccountId, "shared name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionListDataAddress, "id", list),
					resource.TestCheckResourceAttr(abpConditionListDataAddress, "description", "a list"),
				),
			},
		},
	})
}

func TestAbpConditionListDataSourceNoMatchFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedConditionList(abpFakeAccountId, "Sample condition list", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionListDataConfig(abpFakeAccountId, "absent"),
				ExpectError: regexp.MustCompile(`no ABP Condition List named "absent" found`),
			},
			// Names are matched exactly and case-sensitively.
			{
				Config:      abpConditionListDataConfig(abpFakeAccountId, "sample condition list"),
				ExpectError: regexp.MustCompile(`no ABP Condition List named "sample condition list" found`),
			},
		},
	})
}

func TestAbpConditionListDataSourceDuplicateNameFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedConditionList(abpFakeAccountId, "twice", "first")
	f.seedConditionList(abpFakeAccountId, "twice", "second")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionListDataConfig(abpFakeAccountId, "twice"),
				ExpectError: regexp.MustCompile(`multiple ABP Condition Lists named "twice" found`),
			},
		},
	})
}

func TestAbpConditionListDataSourceOtherAccountIsNotVisible(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList(abpFakeOtherAccountId, "elsewhere", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionListDataConfig(abpFakeAccountId, "elsewhere"),
				ExpectError: regexp.MustCompile(`no ABP Condition List named "elsewhere" found`),
			},
			{
				Config: abpConditionListDataConfig(abpFakeOtherAccountId, "elsewhere"),
				Check:  resource.TestCheckResourceAttr(abpConditionListDataAddress, "id", id),
			},
		},
	})
}

// A Policy Directive's Condition List is an ordinary Condition List named
// after the Directive's action, so the lookup finds it. Such a list belongs to
// the Policy, though: incapsula_abp_directive is the way to address it.
func TestAbpConditionListDataSourceFindsPolicyDirectiveList(t *testing.T) {
	f := newAbpFake(t)
	policy := f.seedPolicy("a policy", "d", "block")
	directiveList := abpFakeString(abpFakeList(f.doc(abpFakePolicies, policy)["directives"])[0]["condition_id"])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionListDataConfig(abpFakeAccountId, "block"),
				Check:  resource.TestCheckResourceAttr(abpConditionListDataAddress, "id", directiveList),
			},
		},
	})
}

// A Condition List renamed out of band is found under its new name only.
func TestAbpConditionListDataSourceRereadsOnEveryPlan(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedConditionList(abpFakeAccountId, "a list", "d")
	config := abpConditionListDataConfig(abpFakeAccountId, "a list")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpConditionListDataAddress, "id", id),
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["description"] = "changed in the UI"
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpConditionListDataAddress)["description"]; got != "changed in the UI" {
						return fmt.Errorf("the data source still reports %q", got)
					}
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, id, func(doc map[string]any) {
						doc["name"] = "renamed in the UI"
					})
				},
				Config:      config,
				ExpectError: regexp.MustCompile(`no ABP Condition List named "a list" found`),
			},
		},
	})
}

func TestAbpConditionListDataSourceValidation(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionListDataConfig("not-a-uuid", "a list"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
			{
				Config:      abpConditionListDataConfig(abpFakeAccountId, "  "),
				ExpectError: regexp.MustCompile(`expected "name" to not be an empty string or whitespace`),
			},
		},
	})
}
