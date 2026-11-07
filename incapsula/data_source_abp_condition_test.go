package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for the incapsula_abp_condition data source. It looks a literal
// Condition up by name, and its `managed` argument is the subtle part:
// `true` restricts the lookup to managed Conditions, while `false` or unset
// matches either — the opposite of what the same argument means on
// incapsula_abp_conditions.

const abpConditionDataAddress = "data.incapsula_abp_condition.test"

func abpConditionDataConfig(accountId, name string) string {
	return fmt.Sprintf(`
data "incapsula_abp_condition" "test" {
  account_id = %q
  name       = %q
}
`, accountId, name)
}

func abpConditionDataConfigManaged(accountId, name string, managed bool) string {
	return fmt.Sprintf(`
data "incapsula_abp_condition" "test" {
  account_id = %q
  name       = %q
  managed    = %t
}
`, accountId, name, managed)
}

func TestAbpConditionDataSourceByName(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedCondition(abpFakeAccountId, "okhttp", "matches okhttp", "(matches headers.user_agent \"okhttp\")")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionDataConfig(abpFakeAccountId, "okhttp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionDataAddress, "id", id),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "description", "matches okhttp"),
					// The data source reports the server's expression, not a
					// configured one.
					resource.TestCheckResourceAttr(abpConditionDataAddress, "code",
						abpFakeFormatCode("(matches headers.user_agent \"okhttp\")")),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "template", "custom"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "last_change_by", "fake"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "managed", "false"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "recommended_directives.#", "0"),
					resource.TestCheckResourceAttrSet(abpConditionDataAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpConditionDataAddress, "modified_at"),
				),
			},
		},
	})
}

// A managed Condition is found with `managed = true`, and the API withholds
// its expression, its template and its timestamps.
func TestAbpConditionDataSourceManagedLookup(t *testing.T) {
	f := newAbpFake(t)
	id := abpConditionSeedManaged(f, "Monitoring Tools", "maintained by Imperva", "allow", "captcha_cleared")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionDataConfigManaged(abpFakeAccountId, "Monitoring Tools", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionDataAddress, "id", id),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "managed", "true"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "code", ""),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "template", ""),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "last_change_by", ""),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "created_at", "1970-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "recommended_directives.#", "2"),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "recommended_directives.0", "allow"),
					// `account_id` is a search argument: the matched Condition
					// has none, and it must not be written back over it.
					resource.TestCheckResourceAttr(abpConditionDataAddress, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "name", "Monitoring Tools"),
				),
			},
		},
	})
}

// Without `managed = true` the lookup matches managed Conditions as well.
func TestAbpConditionDataSourceManagedFalseMatchesEither(t *testing.T) {
	f := newAbpFake(t)
	managed := abpConditionSeedManaged(f, "Always true", "maintained by Imperva")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionDataConfig(abpFakeAccountId, "Always true"),
				Check:  resource.TestCheckResourceAttr(abpConditionDataAddress, "id", managed),
			},
			{
				Config: abpConditionDataConfigManaged(abpFakeAccountId, "Always true", false),
				Check:  resource.TestCheckResourceAttr(abpConditionDataAddress, "id", managed),
			},
		},
	})
}

// Because `managed = false` means "either", an account-owned Condition that
// shares its name with a managed one cannot be reached at all. A limitation
// rather than an intention: there is no disambiguator.
func TestAbpConditionDataSourceNameSharedWithManagedIsUnreachable(t *testing.T) {
	f := newAbpFake(t)
	abpConditionSeedManaged(f, "Always true", "maintained by Imperva")
	f.seedCondition(abpFakeAccountId, "Always true", "ours", "(true)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "Always true"),
				ExpectError: regexp.MustCompile(`multiple ABP literal Conditions named "Always true" found`),
			},
			{
				Config:      abpConditionDataConfigManaged(abpFakeAccountId, "Always true", false),
				ExpectError: regexp.MustCompile(`multiple ABP literal Conditions named "Always true" found`),
			},
		},
	})
}

func TestAbpConditionDataSourceDuplicateNameFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedCondition(abpFakeAccountId, "twice", "first", "(true)")
	f.seedCondition(abpFakeAccountId, "twice", "second", "(false)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "twice"),
				ExpectError: regexp.MustCompile(`multiple ABP literal Conditions named "twice" found`),
			},
		},
	})
}

func TestAbpConditionDataSourceNoMatchFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedCondition(abpFakeAccountId, "okhttp", "d", "(true)")
	abpConditionSeedManaged(f, "Monitoring Tools", "d")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "absent"),
				ExpectError: regexp.MustCompile(`no ABP literal Condition named "absent" found`),
			},
			// Names are matched exactly and case-sensitively.
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "OkHttp"),
				ExpectError: regexp.MustCompile(`no ABP literal Condition named "OkHttp" found`),
			},
			// `managed = true` says so, and does not fall back to the
			// account-owned Condition of the same name.
			{
				Config:      abpConditionDataConfigManaged(abpFakeAccountId, "okhttp", true),
				ExpectError: regexp.MustCompile(`no managed ABP literal Condition named "okhttp" found`),
			},
		},
	})
}

// Only the literal variant is a Condition: a Condition List or an entry of the
// same name is skipped.
func TestAbpConditionDataSourceSkipsListsAndEntries(t *testing.T) {
	f := newAbpFake(t)
	list := f.seedConditionList(abpFakeAccountId, "shared name", "a list")
	literal := f.seedCondition(abpFakeAccountId, "shared name", "a literal", "(true)")
	f.seedEntry(list, literal, abpConditionStateActive)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpConditionDataConfig(abpFakeAccountId, "shared name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpConditionDataAddress, "id", literal),
					resource.TestCheckResourceAttr(abpConditionDataAddress, "description", "a literal"),
				),
			},
		},
	})
}

// Another account's Conditions are not visible.
func TestAbpConditionDataSourceOtherAccountIsNotVisible(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedCondition(abpFakeOtherAccountId, "elsewhere", "d", "(true)")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "elsewhere"),
				ExpectError: regexp.MustCompile(`no ABP literal Condition named "elsewhere" found`),
			},
			{
				Config: abpConditionDataConfig(abpFakeOtherAccountId, "elsewhere"),
				Check:  resource.TestCheckResourceAttr(abpConditionDataAddress, "id", id),
			},
		},
	})
}

// A Condition created out of band between two plans is picked up, because the
// lookup runs on every read.
func TestAbpConditionDataSourceRereadsOnEveryPlan(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedCondition(abpFakeAccountId, "okhttp", "d", "(true)")
	config := abpConditionDataConfig(abpFakeAccountId, "okhttp")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpConditionDataAddress, "id", first),
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeConditions, first, func(doc map[string]any) {
						doc["description"] = "changed in the UI"
					})
				},
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpConditionDataAddress)["description"]; got != "changed in the UI" {
						return fmt.Errorf("the data source still reports %q", got)
					}
					return nil
				},
			},
		},
	})
}

func TestAbpConditionDataSourceValidation(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpConditionDataConfig("not-a-uuid", "okhttp"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
			{
				Config:      abpConditionDataConfig(abpFakeAccountId, "  "),
				ExpectError: regexp.MustCompile(`expected "name" to not be an empty string or whitespace`),
			},
		},
	})
}
