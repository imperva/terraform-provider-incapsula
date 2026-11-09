package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The incapsula_abp_proof_of_work_configuration data source, which looks a
// configuration up by name within an Account.

const abpPolicyPowDataAddress = "data.incapsula_abp_proof_of_work_configuration.test"

func abpPolicyPowDataConfig(accountId, name string) string {
	return fmt.Sprintf(`
data "incapsula_abp_proof_of_work_configuration" "test" {
  account_id = %q
  name       = %q
}
`, accountId, name)
}

func TestAbpPolicyPowConfigurationDataSourceByName(t *testing.T) {
	f := newAbpFake(t)
	id := f.seedPowConfig("lookup-me", 42)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPowDataConfig(abpFakeAccountId, "lookup-me"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "id", id),
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "name", "lookup-me"),
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "difficulty", "42"),
					resource.TestCheckResourceAttrSet(abpPolicyPowDataAddress, "algorithm"),
					resource.TestCheckResourceAttrSet(abpPolicyPowDataAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpPolicyPowDataAddress, "modified_at"),
				),
			},
		},
	})
}

// The data source follows the resource: it reports what an update left behind.
func TestAbpPolicyPowConfigurationDataSourceFollowsTheResource(t *testing.T) {
	f := newAbpFake(t)
	config := func(difficulty int, algorithm string) string {
		return abpPolicyPowTestConfig(abpFakeAccountId, "tracked", difficulty, algorithm) + fmt.Sprintf(`
data "incapsula_abp_proof_of_work_configuration" "test" {
  account_id = %q
  name       = incapsula_abp_proof_of_work_configuration.test.name
}
`, abpFakeAccountId)
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config(42, "bbs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "difficulty", "42"),
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "algorithm", "bbs"),
				),
			},
			{
				Config: config(7, "sha1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "difficulty", "7"),
					resource.TestCheckResourceAttr(abpPolicyPowDataAddress, "algorithm", "sha1"),
					resource.TestCheckResourceAttrPair(
						abpPolicyPowDataAddress, "id",
						abpPolicyPowTestAddress, "id"),
				),
			},
		},
	})
}

func TestAbpPolicyPowConfigurationDataSourceNoMatch(t *testing.T) {
	f := newAbpFake(t)
	f.seedPowConfig("something-else", 42)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyPowDataConfig(abpFakeAccountId, "lookup-me"),
				ExpectError: regexp.MustCompile(`no ABP Proof Of Work Configuration named "lookup-me" found in account ` + abpFakeAccountId),
			},
		},
	})
}

// The lookup is scoped to the Account.
func TestAbpPolicyPowConfigurationDataSourceIsScopedToTheAccount(t *testing.T) {
	f := newAbpFake(t)
	f.seedPowConfig("elsewhere", 42)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyPowDataConfig(abpFakeOtherAccountId, "elsewhere"),
				ExpectError: regexp.MustCompile(`no ABP Proof Of Work Configuration named "elsewhere" found in account ` + abpFakeOtherAccountId),
			},
		},
	})
}

// Names are not unique, so an ambiguous lookup has to fail rather than pick one.
func TestAbpPolicyPowConfigurationDataSourceDuplicateName(t *testing.T) {
	f := newAbpFake(t)
	f.seedPowConfig("twice", 42)
	f.seedPowConfig("twice", 7)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPolicyPowDataConfig(abpFakeAccountId, "twice"),
				ExpectError: regexp.MustCompile(`multiple ABP Proof Of Work Configurations named "twice" found in account`),
			},
		},
	})
}

func TestAbpPolicyPowConfigurationDataSourceInvalidArguments(t *testing.T) {
	f := newAbpFake(t)
	providers := f.providers(t)
	for name, step := range map[string]resource.TestStep{
		"account id is not a uuid": {
			Config:      abpPolicyPowDataConfig("not-a-uuid", "pow"),
			ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
		},
		"name is only whitespace": {
			Config:      abpPolicyPowDataConfig(abpFakeAccountId, "   "),
			ExpectError: regexp.MustCompile(`expected "name" to not be an empty string or whitespace`),
		},
		"account id is required": {
			Config: `
data "incapsula_abp_proof_of_work_configuration" "test" {
  name = "pow"
}
`,
			ExpectError: regexp.MustCompile(`(?s)"account_id" is required`),
		},
		"name is required": {
			Config: fmt.Sprintf(`
data "incapsula_abp_proof_of_work_configuration" "test" {
  account_id = %q
}
`, abpFakeAccountId),
			ExpectError: regexp.MustCompile(`(?s)"name" is required`),
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
