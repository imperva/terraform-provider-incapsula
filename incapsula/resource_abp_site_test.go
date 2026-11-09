package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_site. The bulk of them is about Selectors: which
// criteria block a Selector holds is the way it matches, so switching,
// reordering and inserting Selectors has to converge and must never leave the
// previous criteria behind. The rest covers the Site's own attributes, the
// catch-all Selector the backend appends, and the analysis settings a Selector
// carries as a JSON string.

// --- create ------------------------------------------------------------------

// A Site with one Selector per criteria variant applies, re-plans empty and
// reads back every attribute the backend assigned.
func TestAbpSiteCreateWithEveryCriteriaVariant(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("selector-policy", "", "allow")
	config := abpSiteFullConfig("", `  industry = "retail"

  default_max_requests_per_minute  = 60
  default_max_requests_per_session = 600
  default_max_session_length       = "2h"`,
		fmt.Sprintf("    purpose   = \"login\"\n    policy_id = %q\n%s", policyId, abpSitePathPrefix("/login", true)),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
		fmt.Sprintf("    purpose = \"backend_api\"\n%s", abpSiteApiSecurityGroup("checkout-api")),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpSiteAddress, "name", "abp-site"),
					resource.TestCheckResourceAttr(abpSiteAddress, "industry", "retail"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_requests_per_minute", "60"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_requests_per_session", "600"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_session_length", "2h"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "modified_at"),

					// The catch-all the backend appends is not a user Selector.
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "4"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.purpose", "login"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.policy_id", policyId),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.value", "/login"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.case_insensitive", "true"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_regex.#", "0"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.analysis_settings", abpSitePerSiteSettings),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "selector.0.id"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "selector.0.derived_id"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "selector.0.created_at"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "selector.0.modified_at"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.1.path_regex.0.value", `\.png$`),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.1.purpose", ""),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.1.policy_id", ""),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.2.postback.0.value", "web_interrogation"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.3.api_security_group.0.value", "checkout-api"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.3.purpose", "backend_api"),

					abpSiteCheckCriteria(t, f,
						map[string]any{"path_prefix": "/login", "case_insensitive": true},
						map[string]any{"path_regex": `\.png$`},
						map[string]any{"postback": "web_interrogation"},
						map[string]any{"api_security_group": "checkout-api"},
					),
					abpSiteCheckOneCriteriaPerWrite(t, f),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Site with no Selectors of its own is still served by the catch-all the
// backend appends, which the provider exposes as `default_selector` — with the
// `"<site> Default"` Policy the backend created along with the Site.
func TestAbpSiteCreateWithoutSelectorsExposesTheDefaultSelector(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "0"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.#", "1"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.path_prefix.0.case_insensitive", "false"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.purpose", ""),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.analysis_settings", `{"rate_limiting":"per_site"}`),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "default_selector.0.id"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "default_selector.0.derived_id"),
					resource.TestCheckResourceAttrSet(abpSiteAddress, "default_selector.0.policy_id"),
					func(s *terraform.State) error {
						policy := f.doc(abpFakePolicies, abpAttrs(t, s, abpSiteAddress)["default_selector.0.policy_id"])
						if policy == nil {
							return fmt.Errorf("default_selector.0.policy_id names no Policy")
						}
						if name := abpFakeString(policy["name"]); name != "abp-site Default" {
							return fmt.Errorf("the default Selector points at Policy %q", name)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// --- update ------------------------------------------------------------------

// Renaming a Site and changing its defaults is an in-place update, and leaves
// the catch-all Selector — and the Policy it points at — exactly as it was.
func TestAbpSiteUpdateInPlaceKeepsTheDefaultSelector(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteFullConfig("", `  industry                        = "retail"
  default_max_requests_per_minute = 60`, abpSitePathPrefix("/login", false))
	after := abpSiteFullConfig("", `  industry                        = "travel"
  default_max_requests_per_minute = 90
  default_max_session_length      = "1h"`, abpSitePathPrefix("/login", false))

	var siteId, selectorId, defaultId, defaultPolicyId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCapture(t, abpSiteAddress, "id", &siteId),
					abpSiteCapture(t, abpSiteAddress, "selector.0.id", &selectorId),
					abpSiteCapture(t, abpSiteAddress, "default_selector.0.id", &defaultId),
					abpSiteCapture(t, abpSiteAddress, "default_selector.0.policy_id", &defaultPolicyId),
				),
			},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckUnchanged(t, abpSiteAddress, "id", &siteId),
					abpSiteCheckUnchanged(t, abpSiteAddress, "selector.0.id", &selectorId),
					abpSiteCheckUnchanged(t, abpSiteAddress, "default_selector.0.id", &defaultId),
					abpSiteCheckUnchanged(t, abpSiteAddress, "default_selector.0.policy_id", &defaultPolicyId),
					resource.TestCheckResourceAttr(abpSiteAddress, "industry", "travel"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_requests_per_minute", "90"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_session_length", "1h"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "1"),
					// The update must not have the backend append a second
					// catch-all, nor create a second default Policy.
					abpSiteCheckDoc(t, f, func(site map[string]any) error {
						if got := len(abpFakeList(site["selectors"])); got != 2 {
							return fmt.Errorf("the site holds %d selectors, want the user's one plus the catch-all", got)
						}
						return nil
					}),
					func(s *terraform.State) error {
						if got := f.count(abpFakePolicies); got != 1 {
							return fmt.Errorf("the account holds %d policies, want only the site's default one", got)
						}
						return nil
					},
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// Removing `industry` clears it, because the Site's update is a full replace
// and `industry` is not Computed. The `default_max_*` attributes are Computed,
// so removing one keeps whatever the backend holds instead.
func TestAbpSiteUpdateRemovingIndustryClearsItButKeepsComputedDefaults(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteFullConfig("", `  industry                        = "retail"
  default_max_requests_per_minute = 60`, abpSitePathPrefix("/login", false))
	after := abpSiteConfig(abpSitePathPrefix("/login", false))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: before},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "industry", ""),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_max_requests_per_minute", "60"),
					abpSiteCheckDoc(t, f, func(site map[string]any) error {
						if _, ok := site["industry"]; ok {
							return fmt.Errorf("the site still holds industry %v", site["industry"])
						}
						return nil
					}),
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// Changing the account replaces the Site, because the API cannot move one.
func TestAbpSiteUpdateAccountIdReplacesTheSite(t *testing.T) {
	f := newAbpFake(t)
	selector := abpSitePathPrefix("/login", false)
	before := abpSiteConfig(selector)
	after := fmt.Sprintf(`%s
resource "incapsula_abp_site" "test" {
  account_id = %q
  name       = "abp-site"

  selector {
%s
    analysis_settings = data.incapsula_abp_site_analysis_settings.s.json
  }
}
`, abpSiteSettingsData, abpFakeOtherAccountId, selector)

	var siteId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check:  abpSiteCapture(t, abpSiteAddress, "id", &siteId),
			},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckChanged(t, abpSiteAddress, "id", &siteId),
					resource.TestCheckResourceAttr(abpSiteAddress, "account_id", abpFakeOtherAccountId),
					func(s *terraform.State) error {
						if f.exists(abpFakeSites, siteId) {
							return fmt.Errorf("the replaced site %s was not deleted", siteId)
						}
						return nil
					},
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// --- the criteria a Selector matches by -------------------------------------

// Switching a Selector from one match criteria to another must apply, keep the
// Selector's identity, and leave nothing of the previous criteria behind.
func TestAbpSiteMatchCriteriaSwitchConverges(t *testing.T) {
	f := newAbpFake(t)
	var idBefore string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteConfig(abpSitePathPrefix("/login", true)),
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login", "case_insensitive": true}),
					abpSiteCapture(t, abpSiteAddress, "selector.0.id", &idBefore),
				),
			},
			{
				Config: abpSiteConfig(abpSitePathRegex("^/login")),
				Check: resource.ComposeTestCheckFunc(
					// The path prefix and its case_insensitive are gone.
					abpSiteCheckCriteria(t, f, map[string]any{"path_regex": "^/login"}),
					abpSiteCheckOneCriteriaPerWrite(t, f),
					abpSiteCheckUnchanged(t, abpSiteAddress, "selector.0.id", &idBefore),
				),
			},
			{Config: abpSiteConfig(abpSitePathRegex("^/login")), PlanOnly: true},
		},
	})
}

// Every pair of criteria variants converges when one is switched for the
// other, so no variant can strand a key the API would resolve in preference to
// the new one.
func TestAbpSiteMatchCriteriaSwitchConvergesForEveryPair(t *testing.T) {
	variants := []struct {
		name  string
		block string
		want  map[string]any
	}{
		{"path_prefix", abpSitePathPrefix("/login", true), map[string]any{"path_prefix": "/login", "case_insensitive": true}},
		{"path_regex", abpSitePathRegex("^/login"), map[string]any{"path_regex": "^/login"}},
		{"postback", abpSitePostback("web_automation"), map[string]any{"postback": "web_automation"}},
		{"api_security_group", abpSiteApiSecurityGroup("checkout-api"), map[string]any{"api_security_group": "checkout-api"}},
	}

	for _, from := range variants {
		for _, to := range variants {
			if from.name == to.name {
				continue
			}
			t.Run(from.name+"-to-"+to.name, func(t *testing.T) {
				f := newAbpFake(t)
				resource.Test(t, resource.TestCase{
					ProviderFactories: f.providers(t),
					Steps: []resource.TestStep{
						{Config: abpSiteConfig(from.block)},
						{
							Config: abpSiteConfig(to.block),
							Check: resource.ComposeTestCheckFunc(
								abpSiteCheckCriteria(t, f, to.want),
								abpSiteCheckOneCriteriaPerWrite(t, f),
							),
						},
						{Config: abpSiteConfig(to.block), PlanOnly: true},
					},
				})
			})
		}
	}
}

// Turning `case_insensitive` off again converges: the API omits the flag when
// it is false, so a merge that kept the stored value would never let it go.
func TestAbpSiteCaseInsensitiveTogglesBothWays(t *testing.T) {
	f := newAbpFake(t)
	insensitive := abpSiteConfig(abpSitePathPrefix("/login", true))
	sensitive := abpSiteConfig(abpSitePathPrefix("/login", false))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: insensitive,
				Check:  abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login", "case_insensitive": true}),
			},
			{
				Config: sensitive,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login"}),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.case_insensitive", "false"),
				),
			},
			{Config: sensitive, PlanOnly: true},
			{
				Config: insensitive,
				Check:  abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login", "case_insensitive": true}),
			},
			{Config: insensitive, PlanOnly: true},
		},
	})
}

// Reordering Selectors of different kinds must apply and converge.
func TestAbpSiteReorderingMixedKindSelectorsConverges(t *testing.T) {
	f := newAbpFake(t)
	forward := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
	)
	reversed := abpSiteConfig(
		abpSitePostback("web_interrogation"),
		abpSitePathRegex(`\.png$`),
		abpSitePathPrefix("/login", false),
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: forward},
			{
				Config: reversed,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckCriteria(t, f,
						map[string]any{"postback": "web_interrogation"},
						map[string]any{"path_regex": `\.png$`},
						map[string]any{"path_prefix": "/login"},
					),
					abpSiteCheckOneCriteriaPerWrite(t, f),
				),
			},
			{Config: reversed, PlanOnly: true},
		},
	})
}

// Inserting a Selector above Selectors of other kinds must apply and
// converge.
func TestAbpSiteInsertingIntoMixedKindListConverges(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
	)
	after := abpSiteConfig(
		abpSiteApiSecurityGroup("checkout-api"),
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: before},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckCriteria(t, f,
						map[string]any{"api_security_group": "checkout-api"},
						map[string]any{"path_prefix": "/login"},
						map[string]any{"path_regex": `\.png$`},
						map[string]any{"postback": "web_interrogation"},
					),
					abpSiteCheckOneCriteriaPerWrite(t, f),
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// Removing a Selector from the middle of the list converges, and the catch-all
// is not touched.
func TestAbpSiteRemovingASelectorConverges(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
	)
	after := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePostback("web_interrogation"),
	)
	var defaultId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check:  abpSiteCapture(t, abpSiteAddress, "default_selector.0.id", &defaultId),
			},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "2"),
					abpSiteCheckCriteria(t, f,
						map[string]any{"path_prefix": "/login"},
						map[string]any{"postback": "web_interrogation"},
					),
					abpSiteCheckUnchanged(t, abpSiteAddress, "default_selector.0.id", &defaultId),
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// `derived_id` is recomputed from a Selector's criteria: ABP counts the
// requests that hit a Selector against it, so it has to move when the criteria
// does and stay put when nothing else about the Selector matters to it.
func TestAbpSiteDerivedIdFollowsTheCriteria(t *testing.T) {
	f := newAbpFake(t)
	prefix := abpSiteConfig(abpSitePathPrefix("/login", false))
	insensitive := abpSiteConfig(abpSitePathPrefix("/login", true))
	renamed := abpSiteFullConfig("", `  industry = "retail"`, abpSitePathPrefix("/login", true))

	var derived, insensitiveDerived string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: prefix,
				Check:  abpSiteCapture(t, abpSiteAddress, "selector.0.derived_id", &derived),
			},
			{
				Config: insensitive,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckChanged(t, abpSiteAddress, "selector.0.derived_id", &derived),
					abpSiteCapture(t, abpSiteAddress, "selector.0.derived_id", &insensitiveDerived),
				),
			},
			{
				// A change elsewhere on the Site leaves it alone.
				Config: renamed,
				Check:  abpSiteCheckUnchanged(t, abpSiteAddress, "selector.0.derived_id", &insensitiveDerived),
			},
			{Config: renamed, PlanOnly: true},
		},
	})
}

// --- plan-time validation ----------------------------------------------------

// A Selector holding no criteria block is refused at plan time.
func TestAbpSiteSelectorWithoutCriteriaIsRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSiteConfig(""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`selector\[0\]: a selector requires exactly one of .* but holds none`),
			},
		},
	})
}

// A Selector holding two criteria blocks is refused at plan time.
func TestAbpSiteTwoCriteriaBlocksAreRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteConfig(
					abpSitePathPrefix("/login", false) + "\n" + abpSitePathRegex("^/login"),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`selector\[0\]: a selector requires exactly one of .* but holds path_prefix and path_regex`),
			},
		},
	})
}

// The Site's and its Selectors' enumerations, UUIDs and non-empty strings are
// all checked while the configuration is planned.
func TestAbpSiteInvalidValuesAreRejectedAtPlanTime(t *testing.T) {
	cases := []struct {
		name    string
		config  string
		message string
	}{
		{
			name:    "account_id",
			config:  abpSiteRawConfig("", "not-a-uuid", "abp-site", "", abpSitePathPrefix("/login", false)),
			message: `expected "account_id" to be a valid UUID`,
		},
		{
			name:    "padded name",
			config:  abpSiteRawConfig("", abpFakeAccountId, " abp-site ", "", abpSitePathPrefix("/login", false)),
			message: `must not have leading or trailing whitespace`,
		},
		{
			name:    "empty name",
			config:  abpSiteRawConfig("", abpFakeAccountId, "", "", abpSitePathPrefix("/login", false)),
			message: `expected length of name to be in the range \(1 - 100\)`,
		},
		{
			name:    "industry",
			config:  abpSiteFullConfig("", `  industry = "fishing"`, abpSitePathPrefix("/login", false)),
			message: `expected industry to be one of`,
		},
		{
			name:    "purpose",
			config:  abpSiteConfig("    purpose = \"shopping\"\n" + abpSitePathPrefix("/login", false)),
			message: `expected selector.0.purpose to be one of`,
		},
		{
			name:    "policy_id",
			config:  abpSiteConfig("    policy_id = \"not-a-uuid\"\n" + abpSitePathPrefix("/login", false)),
			message: `expected "selector.0.policy_id" to be a valid UUID`,
		},
		{
			name:    "postback",
			config:  abpSiteConfig(abpSitePostback("carrier_pigeon")),
			message: `expected selector.0.postback.0.value to be one of`,
		},
		{
			name:    "empty path_prefix",
			config:  abpSiteConfig(abpSitePathPrefix("", false)),
			message: `expected "selector.0.path_prefix.0.value" to not be an empty string`,
		},
		{
			name: "analysis_settings that is not JSON",
			config: abpSiteConfig(abpSitePathPrefix("/login", false) +
				"\n    analysis_settings = \"per_site\""),
			message: `contains an invalid JSON`,
		},
		{
			name: "two path_prefix blocks",
			config: abpSiteConfig(abpSitePathPrefix("/login", false) + "\n" +
				abpSitePathPrefix("/signup", false)),
			message: `No more than 1 "path_prefix" blocks are allowed`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			f := newAbpFake(t)
			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{
						Config:      test.config,
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(test.message),
					},
				},
			})
		})
	}
}

// --- drift -------------------------------------------------------------------

// A criteria changed out of band is reported and restored, rather than
// adopted.
func TestAbpSiteCriteriaDriftIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", true))
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						abpFakeList(doc["selectors"])[0]["criteria"] = map[string]any{"path_regex": "^/out-of-band"}
					})
				},
				Config: config,
				Check: abpSiteCheckCriteria(t, f,
					map[string]any{"path_prefix": "/login", "case_insensitive": true}),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Site renamed out of band is renamed back.
func TestAbpSiteNameDriftIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false))
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						doc["name"] = "renamed-elsewhere"
					})
				},
				Config: config,
				Check: abpSiteCheckDoc(t, f, func(site map[string]any) error {
					if got := abpFakeString(site["name"]); got != "abp-site" {
						return fmt.Errorf("the site is still named %q", got)
					}
					return nil
				}),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Selector added out of band is removed again, and the Selectors the
// configuration does declare keep their identity.
func TestAbpSiteSelectorAddedOutOfBandIsRemoved(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false))
	var selectorId, defaultId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCapture(t, abpSiteAddress, "selector.0.id", &selectorId),
					abpSiteCapture(t, abpSiteAddress, "default_selector.0.id", &defaultId),
				),
			},
			{
				PreConfig: func() {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						selectors, _ := doc["selectors"].([]any)
						extra := abpSiteSelectorDoc(map[string]any{"path_regex": `\.png$`})
						extra["id"] = abpFakeAbsentId
						// before the catch-all, which stays last
						doc["selectors"] = append([]any{selectors[0], extra}, selectors[1:]...)
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "1"),
					abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login"}),
					abpSiteCheckUnchanged(t, abpSiteAddress, "selector.0.id", &selectorId),
					abpSiteCheckUnchanged(t, abpSiteAddress, "default_selector.0.id", &defaultId),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Selector removed out of band is added back.
func TestAbpSiteSelectorRemovedOutOfBandIsRestored(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						selectors, _ := doc["selectors"].([]any)
						doc["selectors"] = append([]any{selectors[0]}, selectors[2:]...)
					})
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "2"),
					abpSiteCheckCriteria(t, f,
						map[string]any{"path_prefix": "/login"},
						map[string]any{"path_regex": `\.png$`},
					),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Site deleted out of band is created again.
func TestAbpSiteDeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false))
	var siteId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  abpSiteCapture(t, abpSiteAddress, "id", &siteId),
			},
			{
				PreConfig: func() { f.drop(abpFakeSites, siteId) },
				Config:    config,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckChanged(t, abpSiteAddress, "id", &siteId),
					abpSiteCheckCriteria(t, f, map[string]any{"path_prefix": "/login"}),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// --- fields this provider does not model ------------------------------------

// An update preserves the fields the provider does not model, and keeps each
// Selector's with the Selector it belongs to: Selectors are reconciled with
// their live counterparts by `id`.
func TestAbpSiteUnmodelledFieldsSurviveAnUpdate(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteConfig(
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
	)
	after := abpSiteFullConfig("", `  industry = "retail"`,
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check: func(*terraform.State) error {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						doc["something_newer_than_this_provider"] = "site-wide"
						for i, selector := range abpFakeList(doc["selectors"]) {
							selector["something_newer_than_this_provider"] = fmt.Sprintf("selector-%d", i)
						}
					})
					return nil
				},
			},
			{
				Config: after,
				Check: abpSiteCheckDoc(t, f, func(site map[string]any) error {
					if got := abpFakeString(site["something_newer_than_this_provider"]); got != "site-wide" {
						return fmt.Errorf("the site's unmodelled field is %q", got)
					}
					for i, selector := range abpFakeList(site["selectors"]) {
						want := fmt.Sprintf("selector-%d", i)
						if got := abpFakeString(selector["something_newer_than_this_provider"]); got != want {
							return fmt.Errorf("selector[%d] carries %q, want %q", i, got, want)
						}
					}
					return nil
				}),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// Inserting a Selector shifts the ids the provider sends, because a Selector's
// identity is its position in the list. The fields the provider does not model
// follow the id, so they end up on the neighbouring Selector. This is a
// limitation (ABP-REVIEW `incapsula_abp_site` issue 7), not an intention; the
// test pins down what happens today.
func TestAbpSiteInsertingASelectorMovesUnmodelledFields(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteConfig(abpSitePathPrefix("/login", false))
	after := abpSiteConfig(
		abpSitePathPrefix("/first", false),
		abpSitePathPrefix("/login", false),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check: func(*terraform.State) error {
					site := abpSiteOnly(t, f)
					f.mutate(t, abpFakeSites, abpFakeString(site["id"]), func(doc map[string]any) {
						abpFakeList(doc["selectors"])[0]["counted_by_something_else"] = "/login"
					})
					return nil
				},
			},
			{
				Config: after,
				Check: abpSiteCheckDoc(t, f, func(site map[string]any) error {
					selectors := abpSiteUserSelectors(t, site)
					if got := abpFakeString(selectors[0]["counted_by_something_else"]); got != "/login" {
						return fmt.Errorf("the inserted selector carries %q, want the value the shifted id brought with it", got)
					}
					if _, ok := selectors[1]["counted_by_something_else"]; ok {
						return fmt.Errorf("/login kept the field although its id moved")
					}
					return nil
				}),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// --- analysis settings -------------------------------------------------------

// Every rate limiting mode, the per-Selector maxima and
// `use_site_rate_limiting_parameters` round-trip through a Selector's
// `analysis_settings`.
func TestAbpSiteAnalysisSettingsCoverEveryRateLimitingMode(t *testing.T) {
	f := newAbpFake(t)
	extra := `
data "incapsula_abp_site_analysis_settings" "none" {
  rate_limiting = "none"
}

data "incapsula_abp_site_analysis_settings" "custom" {
  rate_limiting              = "custom_scope"
  rate_limiting_custom_scope = "checkout"

  max_requests_per_minute  = 55
  max_requests_per_session = 555
  max_session_length       = "1h"

  use_site_rate_limiting_parameters = false
}
`
	config := abpSiteFullConfig(extra, "",
		abpSitePathPrefix("/login", false),
		abpSitePathRegex(`\.png$`)+"\n    analysis_settings = data.incapsula_abp_site_analysis_settings.none.json",
		abpSitePostback("web_interrogation")+"\n    analysis_settings = data.incapsula_abp_site_analysis_settings.custom.json",
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.analysis_settings", abpSitePerSiteSettings),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.1.analysis_settings",
						`{"rate_limiting":"none","use_site_rate_limiting_parameters":true}`),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.2.analysis_settings",
						`{"rate_limiting":{"custom_scope":"checkout"},"max_requests_per_minute":55,"max_requests_per_session":555,"max_session_length":"1h","use_site_rate_limiting_parameters":false}`),
					// The custom scope reaches the backend as the object the
					// API's oneOf expects, not as the mode name.
					abpSiteCheckDoc(t, f, func(site map[string]any) error {
						settings := abpFakeMap(abpSiteUserSelectors(t, site)[2]["analysis_settings"])
						scope := abpFakeMap(settings["rate_limiting"])
						if got := abpFakeString(scope["custom_scope"]); got != "checkout" {
							return fmt.Errorf("the stored rate_limiting is %v", settings["rate_limiting"])
						}
						return nil
					}),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Changing a Selector's analysis settings is an in-place update that keeps the
// Selector.
func TestAbpSiteAnalysisSettingsUpdateKeepsTheSelector(t *testing.T) {
	f := newAbpFake(t)
	before := abpSiteConfig(abpSitePathPrefix("/login", false))
	after := abpSiteFullConfig(`
data "incapsula_abp_site_analysis_settings" "capped" {
  rate_limiting           = "per_site"
  max_requests_per_minute = 100

  use_site_rate_limiting_parameters = false
}
`, "", abpSitePathPrefix("/login", false)+
		"\n    analysis_settings = data.incapsula_abp_site_analysis_settings.capped.json")

	var selectorId string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check:  abpSiteCapture(t, abpSiteAddress, "selector.0.id", &selectorId),
			},
			{
				Config: after,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCheckUnchanged(t, abpSiteAddress, "selector.0.id", &selectorId),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.analysis_settings",
						`{"rate_limiting":"per_site","max_requests_per_minute":100,"use_site_rate_limiting_parameters":false}`),
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// `analysis_settings` is compared as a string, so JSON that is not byte for
// byte what the provider writes back never converges. This is a limitation
// (ABP-REVIEW `incapsula_abp_site` issue 8), not an intention — the
// `incapsula_abp_site_analysis_settings` data source is what keeps the value
// stable.
func TestAbpSiteHandWrittenAnalysisSettingsNeverConverge(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false) +
		"\n    analysis_settings = \"{ \\\"use_site_rate_limiting_parameters\\\": true, \\\"rate_limiting\\\": \\\"per_site\\\" }\"")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.analysis_settings",
					abpSitePerSiteSettings),
				ExpectNonEmptyPlan: true,
			},
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
		},
	})
}

// --- the zero that reads as "unset" ------------------------------------------

// `default_max_requests_per_minute = 0` passes validation but is left out of
// the request, because the provider reads it with `GetOk`. The API then
// substitutes its own default and the plan asks for `0` again on every run.
// This is a limitation (ABP-REVIEW `incapsula_abp_site` issue 5), not an
// intention; the fake does not substitute anything, so what the test can pin
// down is the request.
func TestAbpSiteZeroDefaultMaximaAreLeftOutOfTheRequest(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteFullConfig("", `  default_max_requests_per_minute  = 0
  default_max_requests_per_session = 0`, abpSitePathPrefix("/login", false))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(*terraform.State) error {
					body := abpSiteLastCreateBody(t, f)
					for _, key := range []string{"default_max_requests_per_minute", "default_max_requests_per_session"} {
						if value, ok := body[key]; ok {
							return fmt.Errorf("the create request carried %s = %v; a zero reaching the API would be a fix, and this test with it", key, value)
						}
					}
					return nil
				},
			},
		},
	})
}

// --- delete ------------------------------------------------------------------

// Deleting a Site deletes the Site, and leaves behind the default Policy and
// the Condition Lists the backend created along with it. This is a limitation
// (ABP-REVIEW `incapsula_abp_site` issue 10), not an intention: the leftovers
// eventually exhaust the account's Policy quota.
func TestAbpSiteDeleteLeavesTheGeneratedPolicyBehind(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false))
	var siteId, defaultPolicyId string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					abpSiteCapture(t, abpSiteAddress, "id", &siteId),
					abpSiteCapture(t, abpSiteAddress, "default_selector.0.policy_id", &defaultPolicyId),
				),
			},
			{Config: abpSiteSettingsData},
		},
		CheckDestroy: func(*terraform.State) error {
			if f.exists(abpFakeSites, siteId) {
				return fmt.Errorf("site %s was not deleted", siteId)
			}
			return nil
		},
	})

	if f.exists(abpFakeSites, siteId) {
		t.Errorf("site %s outlived the configuration", siteId)
	}
	if !f.exists(abpFakePolicies, defaultPolicyId) {
		t.Errorf("the generated default policy %s was deleted too; the leak is fixed and so should this test be", defaultPolicyId)
	}
	if f.count(abpFakeConditions) == 0 {
		t.Error("the generated condition lists were deleted too; the leak is fixed and so should this test be")
	}
}

// --- import ------------------------------------------------------------------

// Import populates every Selector's criteria block, the catch-all and the
// Site's own attributes, and re-plans empty.
func TestAbpSiteImportReplansEmpty(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteFullConfig("", `  industry = "retail"

  default_max_requests_per_minute = 60`,
		"    purpose = \"login\"\n"+abpSitePathPrefix("/login", true),
		abpSitePathRegex(`\.png$`),
		abpSitePostback("web_interrogation"),
		abpSiteApiSecurityGroup("checkout-api"),
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpSiteAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// A Site created outside Terraform can be imported and then managed: the first
// plan after the import is empty.
func TestAbpSiteImportThenApplyPlansNothing(t *testing.T) {
	f := newAbpFake(t)
	siteId := abpSiteSeed(f, map[string]any{
		"name":     "abp-site",
		"industry": "retail",
		"selectors": []any{
			abpSiteSelectorDoc(map[string]any{"path_prefix": "/login", "case_insensitive": true}),
			abpSiteSelectorDoc(map[string]any{"postback": "web_interrogation"}),
		},
	})
	config := abpSiteFullConfig("", `  industry = "retail"`,
		abpSitePathPrefix("/login", true),
		abpSitePostback("web_interrogation"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpSiteAddress,
				ImportState:        true,
				ImportStateId:      siteId,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("imported %d objects, want one", len(states))
					}
					attributes := states[0].Attributes
					for key, want := range map[string]string{
						"account_id":                             abpFakeAccountId,
						"name":                                   "abp-site",
						"industry":                               "retail",
						"selector.#":                             "2",
						"selector.0.path_prefix.0.value":         "/login",
						"selector.1.postback.0.value":            "web_interrogation",
						"default_selector.#":                     "1",
						"default_selector.0.path_prefix.0.value": "/",
					} {
						if attributes[key] != want {
							return fmt.Errorf("imported %s = %q, want %q", key, attributes[key], want)
						}
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
			{
				// And it is the same Site afterwards, not a replacement.
				Config: abpSiteFullConfig("", `  industry = "travel"`,
					abpSitePathPrefix("/login", true),
					abpSitePostback("web_interrogation"),
				),
				Check: resource.TestCheckResourceAttr(abpSiteAddress, "id", siteId),
			},
		},
	})
}

// Importing a Site id that does not exist fails with a message naming it.
func TestAbpSiteImportOfAnUnknownIdFails(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        abpSiteConfig(abpSitePathPrefix("/login", false)),
				ResourceName:  abpSiteAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile("ABP Site " + abpFakeAbsentId + " not found"),
			},
		},
	})
}

// A Selector the configuration writes with a `/` path prefix matches every
// path, just as the catch-all does. It stays the configuration's own: the API
// appends its catch-all behind it, the provider tells the two apart by id, and
// the plan empties. Recognising the catch-all by its criteria instead used to
// swallow such a Selector into `default_selector`, leaving the configuration's
// Selector list one short and a plan that could never converge.
func TestAbpSiteUserSelectorMatchingTheCatchAllConverges(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/", false))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "1"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.#", "1"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.path_prefix.0.value", "/"),
					// The configuration's Selector and the appended catch-all
					// are two different Selectors of the Site.
					func(s *terraform.State) error {
						attributes := abpAttrs(t, s, abpSiteAddress)
						if own, appended := attributes["selector.0.id"], attributes["default_selector.0.id"]; own == appended {
							return fmt.Errorf("the configuration's Selector and the catch-all are both %s", own)
						}
						return nil
					},
					abpSiteCheckDoc(t, f, func(site map[string]any) error {
						if got := len(abpFakeList(site["selectors"])); got != 2 {
							return fmt.Errorf("the site holds %d selectors, want the configuration's one plus the catch-all", got)
						}
						return nil
					}),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Site that has no catch-all — one left that way through the API — can be
// given one by the configuration: the trailing Selector matching every path
// serves the paths no other Selector claims, so `default_selector` stays
// empty and the Selector itself stays in `selector`. The first apply is
// enough, and it adds exactly the one Selector the plan asked for. Reading the
// catch-all off the trailing Selector's criteria used to adopt that Selector
// as the catch-all instead, so the plan asked for it again and the second
// apply left the Site with a duplicate.
func TestAbpSiteConfiguredCatchAllOnASiteWithoutOne(t *testing.T) {
	f := newAbpFake(t)
	siteId := abpSiteSeedVerbatim(f, map[string]any{
		"name":      "no-catch-all",
		"selectors": []any{abpSiteSelectorDoc(map[string]any{"path_prefix": "/login"})},
	})
	config := abpSiteRawConfig("", abpFakeAccountId, "no-catch-all", "",
		abpSitePathPrefix("/login", false),
		abpSitePathPrefix("/", false),
	)

	selectorCount := func(want int) resource.TestCheckFunc {
		return abpSiteCheckDoc(t, f, func(site map[string]any) error {
			if got := len(abpFakeList(site["selectors"])); got != want {
				return fmt.Errorf("the site holds %d selectors, want %d", got, want)
			}
			return nil
		})
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpSiteAddress,
				ImportState:        true,
				ImportStateId:      siteId,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "2"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.value", "/login"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.1.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.#", "0"),
					selectorCount(2),
				),
			},
			{Config: config, PlanOnly: true},
			// And the Site is still those two Selectors after another apply.
			{Config: config, Check: selectorCount(2)},
		},
	})
}

// The catch-all is recognised by the id state recorded for it, so a Selector
// added behind it out of band is reported as drift on the Site's own Selector
// list rather than mistaken for the catch-all.
func TestAbpSiteSelectorAppendedBehindTheCatchAllIsDrift(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					f.mutate(t, abpFakeSites, abpAttrs(t, s, abpSiteAddress)["id"], func(site map[string]any) {
						selectors, _ := site["selectors"].([]any)
						site["selectors"] = append(selectors,
							abpSiteSelectorDoc(map[string]any{"path_regex": `\.png$`}))
						f.normalizeSelectors(site)
					})
					return nil
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.#", "1"),
					resource.TestCheckResourceAttr(abpSiteAddress, "selector.0.path_prefix.0.value", "/login"),
					resource.TestCheckResourceAttr(abpSiteAddress, "default_selector.0.path_prefix.0.value", "/"),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// `default_selector` belongs to the backend, so a configuration cannot write
// it.
func TestAbpSiteDefaultSelectorCannotBeConfigured(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteFullConfig("", `  default_selector {
    path_prefix {
      value = "/"
    }
  }`, abpSitePathPrefix("/login", false)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Can't configure a value for "default_selector"`),
			},
		},
	})
}

// --- schema and round trip ---------------------------------------------------

// A Selector's criteria is only ever what the configuration says, so no
// criteria block may be Computed: a Computed criteria would survive its own
// removal from the configuration, which is what made switching a Selector's
// match criteria impossible.
func TestAbpSiteSelectorCriteriaIsNotComputed(t *testing.T) {
	selector := Provider().ResourcesMap["incapsula_abp_site"].
		Schema["selector"].Elem.(*schema.Resource).Schema
	for _, name := range abpSelectorCriteriaBlocks {
		block := selector[name]
		if block.Computed {
			t.Errorf("selector.%s is Computed", name)
		}
		for attr, s := range block.Elem.(*schema.Resource).Schema {
			if s.Computed {
				t.Errorf("selector.%s.%s is Computed", name, attr)
			}
		}
	}
}

// Every criteria round-trips through its block unchanged.
func TestAbpSiteSelectorCriteriaRoundTrip(t *testing.T) {
	prefix, regex, postback, group := "/login", `\.png$`, "web_interrogation", "checkout-api"
	for _, criteria := range []AbpSelectorCriteria{
		{PathPrefix: &prefix},
		{PathPrefix: &prefix, CaseInsensitive: true},
		{PathRegex: &regex},
		{Postback: &postback},
		{ApiSecurityGroup: &group},
	} {
		block, err := flattenAbpSelectorCriteria(criteria)
		if err != nil {
			t.Fatalf("flattening %+v: %s", criteria, err)
		}
		got, err := extractAbpSelectorCriteria(block)
		if err != nil {
			t.Fatalf("extracting %+v: %s", block, err)
		}
		if before, after := abpSiteMustJson(t, criteria), abpSiteMustJson(t, got); before != after {
			t.Errorf("round-tripped %s into %s", before, after)
		}
	}
}

// A criteria this provider version does not know about is an error rather than
// no block at all, which would be silently rewritten into another criteria.
func TestAbpSiteSelectorCriteriaRejectsAnUnknownCriteria(t *testing.T) {
	if _, err := flattenAbpSelectorCriteria(AbpSelectorCriteria{}); err == nil {
		t.Error("an unrecognised criteria was flattened without an error")
	}
}
