package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Tests for the incapsula_abp_site data source: the two ways of looking a Site
// up, the errors when the lookup is ambiguous or matches nothing, and the
// attributes it exposes for a Site this configuration does not manage.

// abpSiteDataConfig looks a Site up by whichever arguments it is given.
func abpSiteDataConfig(arguments string) string {
	return fmt.Sprintf(`
data "incapsula_abp_site" "test" {
  account_id = %q
%s
}
`, abpFakeAccountId, arguments)
}

// A lookup by `site_id` exposes every attribute of the Site, including the
// catch-all Selector and the fields a Selector carries beyond its criteria.
func TestAbpSiteDataSourceLookupBySiteId(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("selector-policy", "", "allow")
	siteId := abpSiteSeedDetailed(f, "looked-up", policyId)
	catchAll := abpSiteCatchAll(t, f.doc(abpFakeSites, siteId))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteDataConfig(fmt.Sprintf("  site_id = %q", siteId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteDataAddress, "id", siteId),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "site_id", siteId),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "name", "looked-up"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "industry", "travel"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_max_requests_per_minute", "11"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_max_requests_per_session", "22"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_max_session_length", "3h"),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "modified_at"),

					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.#", "4"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.purpose", "login"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.policy_id", policyId),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.path_prefix.0.value", "/login"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.path_prefix.0.case_insensitive", "true"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.analysis_settings", abpSitePerSiteSettings),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "selector.0.id"),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "selector.0.derived_id"),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "selector.0.created_at"),
					resource.TestCheckResourceAttrSet(abpSiteDataAddress, "selector.0.modified_at"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.1.path_regex.0.value", `\.png$`),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.2.postback.0.value", "web_interrogation"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.3.api_security_group.0.value", "checkout-api"),

					// The catch-all is reported as the default Selector, not as
					// a fifth user Selector.
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.#", "1"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.0.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.0.id", abpFakeString(catchAll["id"])),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.0.policy_id", abpFakeString(catchAll["policy_id"])),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.0.analysis_settings", `{"rate_limiting":"per_site"}`),
				),
			},
		},
	})
}

// A lookup by `name` finds the same Site. The match is exact and
// case-sensitive.
func TestAbpSiteDataSourceLookupByName(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("looked-up", map[string]any{"path_prefix": "/login"})
	f.seedSite("another", map[string]any{"path_prefix": "/signup"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteDataConfig(`  name = "looked-up"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteDataAddress, "id", siteId),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "site_id", siteId),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.path_prefix.0.value", "/login"),
				),
			},
			{
				Config:      abpSiteDataConfig(`  name = "Looked-Up"`),
				ExpectError: regexp.MustCompile("no ABP Site matching the given criteria found"),
			},
		},
	})
}

// A Site this Terraform configuration creates can be read back by name in the
// same apply.
func TestAbpSiteDataSourceReadsAManagedSite(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", true)) + `
data "incapsula_abp_site" "test" {
  account_id = incapsula_abp_site.test.account_id
  name       = incapsula_abp_site.test.name
}
`
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(abpSiteDataAddress, "site_id", abpSiteAddress, "id"),
					resource.TestCheckResourceAttrPair(abpSiteDataAddress, "selector.0.id", abpSiteAddress, "selector.0.id"),
					resource.TestCheckResourceAttrPair(abpSiteDataAddress, "selector.0.derived_id", abpSiteAddress, "selector.0.derived_id"),
					resource.TestCheckResourceAttrPair(abpSiteDataAddress, "default_selector.0.policy_id", abpSiteAddress, "default_selector.0.policy_id"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.path_prefix.0.case_insensitive", "true"),
				),
			},
		},
	})
}

// Exactly one of `site_id` and `name` is required: neither is refused while
// the configuration is planned, both when it is read.
func TestAbpSiteDataSourceNeedsExactlyOneArgument(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("looked-up", map[string]any{"path_prefix": "/login"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSiteDataConfig(""),
				ExpectError: regexp.MustCompile("one of `name,site_id` must be specified"),
			},
			{
				Config:      abpSiteDataConfig(fmt.Sprintf("  site_id = %q\n  name = \"looked-up\"", siteId)),
				ExpectError: regexp.MustCompile("exactly one of site_id or name should be set, both were set"),
			},
		},
	})
}

// A lookup matching nothing is an error naming the account, whichever argument
// was used.
func TestAbpSiteDataSourceWithoutAMatchFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedSite("looked-up", map[string]any{"path_prefix": "/login"})

	for name, arguments := range map[string]string{
		"by site_id": fmt.Sprintf("  site_id = %q", abpFakeAbsentId),
		"by name":    `  name = "absent"`,
	} {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{
						Config:      abpSiteDataConfig(arguments),
						ExpectError: regexp.MustCompile("no ABP Site matching the given criteria found in account " + abpFakeAccountId),
					},
				},
			})
		})
	}
}

// Site names are not unique, so a name matching two Sites is an error rather
// than an arbitrary pick.
func TestAbpSiteDataSourceWithSeveralMatchesFails(t *testing.T) {
	f := newAbpFake(t)
	f.seedSite("twice", map[string]any{"path_prefix": "/login"})
	f.seedSite("twice", map[string]any{"path_prefix": "/signup"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSiteDataConfig(`  name = "twice"`),
				ExpectError: regexp.MustCompile("multiple ABP Sites match the given criteria in account " + abpFakeAccountId),
			},
		},
	})
}

// The lookup is scoped to one account: a Site of another account is not found,
// not even by its id.
func TestAbpSiteDataSourceDoesNotSeeOtherAccounts(t *testing.T) {
	f := newAbpFake(t)
	elsewhere := abpSiteSeed(f, map[string]any{
		"account_id": abpFakeOtherAccountId,
		"name":       "elsewhere",
		"selectors":  []any{abpSiteSelectorDoc(map[string]any{"path_prefix": "/login"})},
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSiteDataConfig(fmt.Sprintf("  site_id = %q", elsewhere)),
				ExpectError: regexp.MustCompile("no ABP Site matching the given criteria found in account " + abpFakeAccountId),
			},
		},
	})
}

// A Site whose catch-all Selector is not last has no default Selector at all:
// the `/` Selector is reported as a user Selector and `default_selector` is
// empty. This is a limitation (ABP-REVIEW `incapsula_abp_site` issue 16), not
// an intention — the API exposes no flag to tell the two apart, so the
// provider infers the default from the trailing Selector's shape.
func TestAbpSiteDataSourceStrandedCatchAllLeavesNoDefaultSelector(t *testing.T) {
	f := newAbpFake(t)
	siteId := abpSiteSeedVerbatim(f, map[string]any{
		"name": "stranded",
		"selectors": []any{
			abpSiteSelectorDoc(map[string]any{"path_prefix": "/"}),
			abpSiteSelectorDoc(map[string]any{"path_prefix": "/login"}),
		},
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteDataConfig(fmt.Sprintf("  site_id = %q", siteId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSiteDataAddress, "default_selector.#", "0"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.#", "2"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.0.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttr(abpSiteDataAddress, "selector.1.path_prefix.0.value", "/login"),
				),
			},
		},
	})
}

// A Selector matching on something this provider version does not know about
// is an error, rather than a Selector silently reported with no criteria at
// all — which an apply would then rewrite into a different one.
func TestAbpSiteDataSourceUnsupportedCriteriaIsAnError(t *testing.T) {
	f := newAbpFake(t)
	siteId := abpSiteSeedVerbatim(f, map[string]any{
		"name": "from-the-future",
		"selectors": []any{
			abpSiteSelectorDoc(map[string]any{"query_parameter": "token"}),
			abpSiteSelectorDoc(map[string]any{"path_prefix": "/"}),
		},
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSiteDataConfig(fmt.Sprintf("  site_id = %q", siteId)),
				ExpectError: regexp.MustCompile("selector uses a match criteria this provider version does not support"),
			},
		},
	})
}

// The arguments are validated while the configuration is planned.
func TestAbpSiteDataSourceInvalidArgumentsAreRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	for name, test := range map[string]struct{ arguments, message string }{
		"account_id": {fmt.Sprintf("  site_id = %q", abpFakeAbsentId), `expected "account_id" to be a valid UUID`},
		"site_id":    {`  site_id = "not-a-uuid"`, `expected "site_id" to be a valid UUID`},
	} {
		t.Run(name, func(t *testing.T) {
			config := abpSiteDataConfig(test.arguments)
			if name == "account_id" {
				config = fmt.Sprintf(`
data "incapsula_abp_site" "test" {
  account_id = "not-a-uuid"
%s
}
`, test.arguments)
			}
			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{Config: config, ExpectError: regexp.MustCompile(test.message)},
				},
			})
		})
	}
}
