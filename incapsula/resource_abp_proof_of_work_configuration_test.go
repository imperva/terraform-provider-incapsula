package incapsula

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_proof_of_work_configuration, against the shared fake ABP
// backend.

const abpPolicyPowTestAddress = "incapsula_abp_proof_of_work_configuration.test"

// abpPolicyPowTestConfig builds a Proof Of Work Configuration. An empty
// algorithm leaves the attribute out of the configuration altogether.
func abpPolicyPowTestConfig(accountId, name string, difficulty int, algorithm string) string {
	var algorithmLine string
	if algorithm != "" {
		algorithmLine = fmt.Sprintf("\n  algorithm  = %q", algorithm)
	}
	return fmt.Sprintf(`
resource "incapsula_abp_proof_of_work_configuration" "test" {
  account_id = %q
  name       = %q
  difficulty = %d%s
}
`, accountId, name, difficulty, algorithmLine)
}

// --- baseline ---------------------------------------------------------------

func TestAbpPolicyPowConfigurationCreateIsIdempotent(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "name", "pow"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "difficulty", "42"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "algorithm", "bbs"),
					resource.TestCheckResourceAttrSet(abpPolicyPowTestAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpPolicyPowTestAddress, "modified_at"),
					func(s *terraform.State) error {
						stored := f.only(t, abpFakePowConfigs)
						if abpFakeString(stored["algorithm"]) != "bbs" {
							return fmt.Errorf("the server stored %+v", stored)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Difficulty and algorithm are updated in place, keeping the configuration's id.
func TestAbpPolicyPowConfigurationUpdateIsInPlace(t *testing.T) {
	f := newAbpFake(t)
	var id string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpPolicyPowTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow-renamed", 7, "sha1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "name", "pow-renamed"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "difficulty", "7"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "algorithm", "sha1"),
					func(s *terraform.State) error {
						if got := abpAttrs(t, s, abpPolicyPowTestAddress)["id"]; got != id {
							return fmt.Errorf("the configuration was replaced: %s became %s", id, got)
						}
						if len(f.deletedIds()) != 0 {
							return fmt.Errorf("an in-place update deleted %v", f.deletedIds())
						}
						return nil
					},
				),
			},
			{Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow-renamed", 7, "sha1"), PlanOnly: true},
		},
	})
}

// account_id is ForceNew, so moving the configuration to another Account
// replaces it.
func TestAbpPolicyPowConfigurationAccountIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	var id string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs"),
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpPolicyPowTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpPolicyPowTestConfig(abpFakeOtherAccountId, "pow", 42, "bbs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "account_id", abpFakeOtherAccountId),
					func(s *terraform.State) error {
						if abpAttrs(t, s, abpPolicyPowTestAddress)["id"] == id {
							return fmt.Errorf("the configuration was updated in place instead of replaced")
						}
						if f.exists(abpFakePowConfigs, id) {
							return fmt.Errorf("the replaced configuration %s still exists", id)
						}
						return nil
					},
				),
			},
		},
	})
}

// --- the algorithm, which the backend derives nothing from ------------------

// With `algorithm` omitted the attribute is filled in from whatever the API
// reports, and the configuration still re-plans empty.
func TestAbpPolicyPowConfigurationAlgorithmCanBeOmitted(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					want := abpFakeString(f.only(t, abpFakePowConfigs)["algorithm"])
					if got := abpAttrs(t, s, abpPolicyPowTestAddress)["algorithm"]; got != want {
						return fmt.Errorf("state holds algorithm %q, the server holds %q", got, want)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Known limitation, not an intention: `algorithm` is Optional + Computed, so
// dropping it from the configuration is ignored rather than clearing it. There
// is no way to reach "no algorithm" once one has been set.
func TestAbpPolicyPowConfigurationAlgorithmCannotBeCleared(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")},
			{
				Config:   abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, ""),
				PlanOnly: true,
			},
			{
				Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "algorithm", "bbs"),
					func(s *terraform.State) error {
						if got := abpFakeString(f.only(t, abpFakePowConfigs)["algorithm"]); got != "bbs" {
							return fmt.Errorf("the server holds algorithm %q", got)
						}
						return nil
					},
				),
			},
		},
	})
}

// Known limitation, not an intention: because `algorithm` is Optional +
// Computed, an algorithm changed outside Terraform is adopted into state
// instead of being restored — as long as the configuration does not name one.
func TestAbpPolicyPowConfigurationOmittedAlgorithmDriftIsAdopted(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePowConfigs)["id"])
					f.mutate(t, abpFakePowConfigs, id, func(doc map[string]any) {
						doc["algorithm"] = "sha1"
					})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "algorithm", "sha1"),
			},
		},
	})
}

// --- drift ------------------------------------------------------------------

// Name, difficulty and an explicitly configured algorithm are all restored
// after an out-of-band change.
func TestAbpPolicyPowConfigurationDriftIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePowConfigs)["id"])
					f.mutate(t, abpFakePowConfigs, id, func(doc map[string]any) {
						doc["name"] = "renamed-out-of-band"
						doc["difficulty"] = 1
						doc["algorithm"] = "sha1"
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "name", "pow"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "difficulty", "42"),
					resource.TestCheckResourceAttr(abpPolicyPowTestAddress, "algorithm", "bbs"),
					func(s *terraform.State) error {
						stored := f.only(t, abpFakePowConfigs)
						if abpFakeString(stored["name"]) != "pow" ||
							abpFakeString(stored["algorithm"]) != "bbs" ||
							fmt.Sprint(stored["difficulty"]) != "42" {
							return fmt.Errorf("server state not restored: %+v", stored)
						}
						return nil
					},
				),
			},
		},
	})
}

// Out-of-band deletion must lead to a re-create, not an error.
func TestAbpPolicyPowConfigurationDeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	var id string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = abpAttrs(t, s, abpPolicyPowTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakePowConfigs, id) },
				Config:    config,
				Check: func(s *terraform.State) error {
					recreated := abpFakeString(f.only(t, abpFakePowConfigs)["id"])
					if recreated == id {
						return fmt.Errorf("the configuration was not recreated")
					}
					if got := abpAttrs(t, s, abpPolicyPowTestAddress)["id"]; got != recreated {
						return fmt.Errorf("state holds %s, the server holds %s", got, recreated)
					}
					return nil
				},
			},
		},
	})
}

// An update is a read-modify-write, so a field the provider does not model
// survives an unrelated change.
func TestAbpPolicyPowConfigurationUnmodelledFieldSurvivesUpdate(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")},
			{
				PreConfig: func() {
					id := abpFakeString(f.only(t, abpFakePowConfigs)["id"])
					f.mutate(t, abpFakePowConfigs, id, func(doc map[string]any) {
						doc["invented_by_the_service"] = "keep me"
					})
				},
				Config: abpPolicyPowTestConfig(abpFakeAccountId, "pow", 7, "bbs"),
				Check: func(s *terraform.State) error {
					stored := f.only(t, abpFakePowConfigs)
					if abpFakeString(stored["invented_by_the_service"]) != "keep me" {
						return fmt.Errorf("the update cleared a field the provider does not model: %+v", stored)
					}
					return nil
				},
			},
		},
	})
}

// --- import -----------------------------------------------------------------

func TestAbpPolicyPowConfigurationImportStateVerify(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpPolicyPowTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Importing a configuration created elsewhere, then applying, must neither
// replace it nor change it.
func TestAbpPolicyPowConfigurationImportThenApply(t *testing.T) {
	f := newAbpFake(t)
	seeded := f.seedPowConfig("pow", 42)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpPolicyPowTestAddress,
				ImportState:        true,
				ImportStateId:      seeded,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPolicyPowTestAddress)["id"]; got != seeded {
						return fmt.Errorf("the imported configuration %s was replaced by %s", seeded, got)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestAbpPolicyPowConfigurationImportMissing(t *testing.T) {
	f := newAbpFake(t)
	config := abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs")
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:        config,
				ResourceName:  abpPolicyPowTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Proof Of Work Configuration ` + abpFakeAbsentId + ` not found`),
			},
		},
	})
}

// --- validation -------------------------------------------------------------

func TestAbpPolicyPowConfigurationInvalidConfigurations(t *testing.T) {
	f := newAbpFake(t)
	providers := f.providers(t)
	for name, step := range map[string]resource.TestStep{
		"account id is not a uuid": {
			Config:      abpPolicyPowTestConfig("not-a-uuid", "pow", 42, "bbs"),
			ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
		},
		"difficulty below one": {
			Config:      abpPolicyPowTestConfig(abpFakeAccountId, "pow", 0, "bbs"),
			ExpectError: regexp.MustCompile(`expected difficulty to be at least \(1\)`),
		},
		"unknown algorithm": {
			Config:      abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "md5"),
			ExpectError: regexp.MustCompile(`expected algorithm to be one of`),
		},
		"padded name": {
			Config:      abpPolicyPowTestConfig(abpFakeAccountId, "  pow  ", 42, "bbs"),
			ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
		},
		"whitespace only name": {
			Config:      abpPolicyPowTestConfig(abpFakeAccountId, "   ", 42, "bbs"),
			ExpectError: regexp.MustCompile(`must not be empty or consist only of whitespace`),
		},
		"name too long": {
			Config:      abpPolicyPowTestConfig(abpFakeAccountId, strings.Repeat("n", 101), 42, "bbs"),
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

// --- use from a Policy Directive -------------------------------------------

// A Directive may reference a Proof Of Work Configuration only when its action
// is `proof_of_work`; on any other action the API refuses the whole Policy.
func TestAbpPolicyPowConfigurationIsOnlyValidOnAProofOfWorkDirective(t *testing.T) {
	f := newAbpFake(t)
	policy := func(action string) string {
		return abpPolicyPowTestConfig(abpFakeAccountId, "pow", 42, "bbs") + fmt.Sprintf(`
resource "incapsula_abp_policy" "test" {
  account_id = %q
  name       = "p1"

  directive {
    action                         = %q
    proof_of_work_configuration_id = incapsula_abp_proof_of_work_configuration.test.id
  }
}
`, abpFakeAccountId, action)
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: policy("proof_of_work"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						abpPolicyTestAddress, "directive.0.proof_of_work_configuration_id",
						abpPolicyPowTestAddress, "id"),
					resource.TestCheckResourceAttrSet(abpPolicyTestAddress, "directive.0.skip_condition_list_id"),
				),
			},
			{Config: policy("proof_of_work"), PlanOnly: true},
			{
				Config:      policy("captcha"),
				ExpectError: regexp.MustCompile(`proof_of_work_configuration_id is only allowed for proof_of_work action`),
			},
		},
	})
}
