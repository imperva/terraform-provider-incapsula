package incapsula

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Tests for the incapsula_abp_sites data source, which lists every Site of an
// account. The documented order is "most to least significant priority"; the
// shared fake's list endpoint returns the Sites in creation order, so the
// ordering itself is not what these tests assert.

// abpSitesDataConfig lists the Sites of an account.
func abpSitesDataConfig(accountId string) string {
	return fmt.Sprintf(`
data "incapsula_abp_sites" "test" {
  account_id = %q
}
`, accountId)
}

// The list carries every Site of the account, with the same attributes the
// single-Site data source exposes.
func TestAbpSitesDataSourceListsEveryAttribute(t *testing.T) {
	f := newAbpFake(t)
	policyId := f.seedPolicy("selector-policy", "", "allow")
	first := abpSiteSeedDetailed(f, "first", policyId)
	second := f.seedSite("second", map[string]any{"postback": "web_automation"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSitesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "2"),

					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.id", first),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.account_id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.name", "first"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.industry", "travel"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.default_max_requests_per_minute", "11"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.default_max_requests_per_session", "22"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.default_max_session_length", "3h"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.#", "4"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.0.purpose", "login"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.0.policy_id", policyId),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.0.path_prefix.0.value", "/login"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.0.path_prefix.0.case_insensitive", "true"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.1.path_regex.0.value", `\.png$`),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.2.postback.0.value", "web_interrogation"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.selector.3.api_security_group.0.value", "checkout-api"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.default_selector.#", "1"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.default_selector.0.path_prefix.0.value", "/"),
					resource.TestCheckResourceAttrSet(abpSitesDataAddress, "sites.0.default_selector.0.policy_id"),
					resource.TestCheckResourceAttrSet(abpSitesDataAddress, "sites.0.created_at"),

					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.id", second),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.name", "second"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.industry", ""),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.selector.#", "1"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.selector.0.postback.0.value", "web_automation"),
				),
			},
		},
	})
}

// An account with no Sites yields an empty list rather than an error.
func TestAbpSitesDataSourceEmptyAccountIsEmpty(t *testing.T) {
	f := newAbpFake(t)
	f.seedSite("elsewhere", map[string]any{"path_prefix": "/login"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSitesDataConfig(abpFakeOtherAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "id", abpFakeOtherAccountId),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "0"),
				),
			},
		},
	})
}

// The list is scoped to one account.
func TestAbpSitesDataSourceExcludesOtherAccounts(t *testing.T) {
	f := newAbpFake(t)
	mine := f.seedSite("mine", map[string]any{"path_prefix": "/login"})
	abpSiteSeed(f, map[string]any{
		"account_id": abpFakeOtherAccountId,
		"name":       "theirs",
		"selectors":  []any{abpSiteSelectorDoc(map[string]any{"path_prefix": "/login"})},
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSitesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "1"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.0.id", mine),
				),
			},
		},
	})
}

// The list is read afresh every time, so a Site created or deleted out of band
// appears in, and disappears from, the next plan.
func TestAbpSitesDataSourceFollowsSitesCreatedAndDeletedOutOfBand(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("first", map[string]any{"path_prefix": "/login"})
	var second string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSitesDataConfig(abpFakeAccountId),
				Check:  resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "1"),
			},
			{
				PreConfig: func() {
					second = f.seedSite("second", map[string]any{"path_prefix": "/signup"})
				},
				Config: abpSitesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "2"),
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.1.name", "second"),
				),
			},
			{
				PreConfig: func() { f.drop(abpFakeSites, first) },
				Config:    abpSitesDataConfig(abpFakeAccountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "1"),
					resource.TestCheckResourceAttrPtr(abpSitesDataAddress, "sites.0.id", &second),
				),
			},
		},
	})
}

// Every Site the configuration manages is in the list, which is what makes the
// list usable for the priority resources.
func TestAbpSitesDataSourceIncludesManagedSites(t *testing.T) {
	f := newAbpFake(t)
	config := abpSiteConfig(abpSitePathPrefix("/login", false)) + `
data "incapsula_abp_sites" "test" {
  account_id = incapsula_abp_site.test.account_id
}
`
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpSitesDataAddress, "sites.#", "1"),
					resource.TestCheckResourceAttrPair(abpSitesDataAddress, "sites.0.id", abpSiteAddress, "id"),
					resource.TestCheckResourceAttrPair(abpSitesDataAddress, "sites.0.selector.0.id", abpSiteAddress, "selector.0.id"),
				),
			},
		},
	})
}

// A Site with a Selector this provider version cannot represent fails the
// listing, naming the Site, rather than being listed with an empty criteria.
func TestAbpSitesDataSourceUnsupportedCriteriaNamesTheSite(t *testing.T) {
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
				Config:      abpSitesDataConfig(abpFakeAccountId),
				ExpectError: regexp.MustCompile("ABP Site " + siteId + `: selector\[0\]: selector uses a match criteria this provider version does not support`),
			},
		},
	})
}

// The account is validated while the configuration is planned.
func TestAbpSitesDataSourceInvalidAccountIdIsRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpSitesDataConfig("not-a-uuid"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
		},
	})
}
