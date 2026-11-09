package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Helpers shared by the Site tests: configuration builders, seeders for the
// Site shapes the shared fake's seedSite cannot express, and checks on what
// the fake ended up holding and on what the provider sent to get there.

// abpSiteAddress is the Site every resource test manages.
const abpSiteAddress = "incapsula_abp_site.test"

// abpSiteDataAddress and abpSitesDataAddress are the data sources under test.
const (
	abpSiteDataAddress  = "data.incapsula_abp_site.test"
	abpSitesDataAddress = "data.incapsula_abp_sites.test"
)

// abpSitePerSiteSettings is the analysis settings document the
// `incapsula_abp_site_analysis_settings` data source produces for
// `rate_limiting = "per_site"` and nothing else. Selectors are compared as
// strings, so this is also the only value a configuration can write by hand
// without diffing forever.
const abpSitePerSiteSettings = `{"rate_limiting":"per_site","use_site_rate_limiting_parameters":true}`

// abpSiteSettingsData declares the analysis settings every selector gets when
// a test does not ask for its own.
const abpSiteSettingsData = `
data "incapsula_abp_site_analysis_settings" "s" {
  rate_limiting = "per_site"
}
`

// abpSiteRawConfig assembles a Site configuration. `extra` is HCL placed
// beside the Site (data sources, Policies), `attributes` are written inside
// it, and each selector is the body of one `selector` block; a body that does
// not set `analysis_settings` itself gets the shared `per_site` one.
func abpSiteRawConfig(extra, accountId, name, attributes string, selectors ...string) string {
	var b strings.Builder
	b.WriteString(abpSiteSettingsData)
	b.WriteString(extra)
	fmt.Fprintf(&b, `
resource "incapsula_abp_site" "test" {
  account_id = %q
  name       = %q
%s
`, accountId, name, attributes)
	for _, body := range selectors {
		fmt.Fprintf(&b, "\n  selector {\n%s\n", body)
		if !strings.Contains(body, "analysis_settings") {
			b.WriteString("    analysis_settings = data.incapsula_abp_site_analysis_settings.s.json\n")
		}
		b.WriteString("  }\n")
	}
	b.WriteString("}\n")
	return b.String()
}

// abpSiteFullConfig is abpSiteRawConfig for the account and name the tests
// use unless they are about those two.
func abpSiteFullConfig(extra, attributes string, selectors ...string) string {
	return abpSiteRawConfig(extra, abpFakeAccountId, "abp-site", attributes, selectors...)
}

// abpSiteConfig builds a Site with one selector per criteria block given.
func abpSiteConfig(selectors ...string) string {
	return abpSiteFullConfig("", "", selectors...)
}

func abpSitePathPrefix(value string, caseInsensitive bool) string {
	return fmt.Sprintf("    path_prefix {\n      value            = %q\n      case_insensitive = %t\n    }", value, caseInsensitive)
}

func abpSitePathRegex(value string) string {
	return fmt.Sprintf("    path_regex {\n      value = %q\n    }", value)
}

func abpSitePostback(value string) string {
	return fmt.Sprintf("    postback {\n      value = %q\n    }", value)
}

func abpSiteApiSecurityGroup(value string) string {
	return fmt.Sprintf("    api_security_group {\n      value = %q\n    }", value)
}

// --- seeding -----------------------------------------------------------------

// abpSiteSeed stores a Site created outside Terraform, taking the attributes
// the shared seedSite does not: `industry`, the `default_max_*` values and
// whatever a Selector carries beyond its criteria. The catch-all Selector the
// backend maintains is appended, as it is for every real Site.
func abpSiteSeed(f *abpFake, site map[string]any) string {
	id := f.seed(abpFakeSites, site)

	f.mu.Lock()
	defer f.mu.Unlock()
	stored := f.store[abpFakeSites][id]
	f.normalizeSelectors(stored)
	f.addDefaultSelector(stored)
	return id
}

// abpSiteSeedVerbatim stores a Site exactly as given, without appending the
// catch-all. It is the only way to get the Selector order the backend would
// never produce by itself: the `/` Selector somewhere other than last.
func abpSiteSeedVerbatim(f *abpFake, site map[string]any) string {
	id := f.seed(abpFakeSites, site)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.normalizeSelectors(f.store[abpFakeSites][id])
	return id
}

// abpSiteSelectorDoc is one stored Selector, for seeding and for out-of-band
// changes. Only the criteria varies between the tests that use it.
func abpSiteSelectorDoc(criteria map[string]any) map[string]any {
	return map[string]any{
		"criteria":          criteria,
		"analysis_settings": map[string]any{"rate_limiting": "per_site", "use_site_rate_limiting_parameters": true},
	}
}

// abpSiteSeedDetailed stores a Site carrying every attribute the data sources
// expose: an industry, the three default maxima, and one Selector per criteria
// variant, the first of them with a purpose and a Policy of its own.
func abpSiteSeedDetailed(f *abpFake, name, policyId string) string {
	login := abpSiteSelectorDoc(map[string]any{"path_prefix": "/login", "case_insensitive": true})
	login["purpose"] = "login"
	login["policy_id"] = policyId

	return abpSiteSeed(f, map[string]any{
		"name":                             name,
		"industry":                         "travel",
		"default_max_requests_per_minute":  11,
		"default_max_requests_per_session": 22,
		"default_max_session_length":       "3h",
		"selectors": []any{
			login,
			abpSiteSelectorDoc(map[string]any{"path_regex": `\.png$`}),
			abpSiteSelectorDoc(map[string]any{"postback": "web_interrogation"}),
			abpSiteSelectorDoc(map[string]any{"api_security_group": "checkout-api"}),
		},
	})
}

// --- inspecting the fake -----------------------------------------------------

// abpSiteOnly returns the single stored Site.
func abpSiteOnly(t *testing.T, f *abpFake) map[string]any {
	t.Helper()
	return f.only(t, abpFakeSites)
}

// abpSiteUserSelectors returns a Site's user-defined Selectors, in order,
// without the catch-all the backend keeps last.
func abpSiteUserSelectors(t *testing.T, site map[string]any) []map[string]any {
	t.Helper()
	selectors := abpFakeList(site["selectors"])
	if len(selectors) == 0 {
		t.Fatalf("site %s has no selectors at all, not even the catch-all", abpFakeString(site["id"]))
	}
	last := abpFakeResolveCriteria(selectors[len(selectors)-1])
	if abpFakeString(last["path_prefix"]) != "/" {
		t.Fatalf("the last selector of site %s is %v, not the catch-all", abpFakeString(site["id"]), last)
	}
	return selectors[:len(selectors)-1]
}

// abpSiteCatchAll returns the catch-all Selector the backend maintains.
func abpSiteCatchAll(t *testing.T, site map[string]any) map[string]any {
	t.Helper()
	selectors := abpFakeList(site["selectors"])
	return selectors[len(abpSiteUserSelectors(t, site))]
}

// abpSiteCriteria returns the criteria of the only Site's user-defined
// Selectors, in order.
func abpSiteCriteria(t *testing.T, f *abpFake) []map[string]any {
	t.Helper()
	out := []map[string]any{}
	for _, selector := range abpSiteUserSelectors(t, abpSiteOnly(t, f)) {
		out = append(out, abpFakeMap(selector["criteria"]))
	}
	return out
}

// abpSiteSelectorIds returns the ids of the only Site's user-defined
// Selectors, in order.
func abpSiteSelectorIds(t *testing.T, f *abpFake) []string {
	t.Helper()
	out := []string{}
	for _, selector := range abpSiteUserSelectors(t, abpSiteOnly(t, f)) {
		out = append(out, abpFakeString(selector["id"]))
	}
	return out
}

// abpSiteWrittenSelectors returns, for every Site write the provider made, the
// selectors it submitted.
func abpSiteWrittenSelectors(t *testing.T, f *abpFake) [][]map[string]any {
	t.Helper()
	var out [][]map[string]any
	for _, request := range f.writes() {
		if request.Method == http.MethodDelete || request.Body == "" || !strings.Contains(request.Path, "/site") {
			continue
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
			t.Fatalf("recorded request is not JSON: %s", err)
		}
		if _, ok := body["selectors"]; !ok {
			continue
		}
		out = append(out, abpFakeList(body["selectors"]))
	}
	return out
}

// abpSiteLastCreateBody returns the body of the Site create request, for tests
// about what the provider leaves out of it.
func abpSiteLastCreateBody(t *testing.T, f *abpFake) map[string]any {
	t.Helper()
	var found map[string]any
	for _, request := range f.writes() {
		if request.Method != http.MethodPost || !strings.HasSuffix(request.Path, "/site") {
			continue
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
			t.Fatalf("recorded request is not JSON: %s", err)
		}
		found = body
	}
	if found == nil {
		t.Fatal("the provider made no Site create request")
	}
	return found
}

// --- checks ------------------------------------------------------------------

// abpSiteCheckCriteria asserts the criteria the fake holds for the only Site's
// user-defined Selectors, in order.
func abpSiteCheckCriteria(t *testing.T, f *abpFake, want ...map[string]any) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := abpSiteCriteria(t, f)
		if len(got) != len(want) {
			return fmt.Errorf("expected %d user selectors, got %d: %v", len(want), len(got), got)
		}
		for i := range want {
			if len(got[i]) != len(want[i]) {
				return fmt.Errorf("selector[%d] criteria %v, want %v", i, got[i], want[i])
			}
			for key, value := range want[i] {
				if fmt.Sprint(got[i][key]) != fmt.Sprint(value) {
					return fmt.Errorf("selector[%d] criteria %v, want %v", i, got[i], want[i])
				}
			}
		}
		return nil
	}
}

// abpSiteCheckOneCriteriaPerWrite asserts that every selector of every Site
// write carried exactly one criteria key. The API resolves the criteria oneOf
// by key presence, so a leftover key from the previous variant would silently
// win instead of failing.
func abpSiteCheckOneCriteriaPerWrite(t *testing.T, f *abpFake) resource.TestCheckFunc {
	return func(*terraform.State) error {
		for request, selectors := range abpSiteWrittenSelectors(t, f) {
			for i, selector := range selectors {
				var keys []string
				criteria := abpFakeMap(selector["criteria"])
				for _, key := range abpFakeCriteriaKeys {
					if _, ok := criteria[key]; ok {
						keys = append(keys, key)
					}
				}
				if len(keys) != 1 {
					return fmt.Errorf("write %d selector %d carried criteria keys %v, want exactly one",
						request, i, keys)
				}
			}
		}
		return nil
	}
}

// abpSiteCheckDoc runs an assertion against the only stored Site.
func abpSiteCheckDoc(t *testing.T, f *abpFake, check func(site map[string]any) error) resource.TestCheckFunc {
	return func(*terraform.State) error {
		return check(abpSiteOnly(t, f))
	}
}

// abpSiteCapture records a state attribute, so a later step can assert whether
// it changed.
func abpSiteCapture(t *testing.T, address, key string, into *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*into = abpAttrs(t, s, address)[key]
		return nil
	}
}

// abpSiteCheckUnchanged asserts a state attribute still has the captured
// value.
func abpSiteCheckUnchanged(t *testing.T, address, key string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got := abpAttrs(t, s, address)[key]; got != *previous {
			return fmt.Errorf("%s changed from %q to %q", key, *previous, got)
		}
		return nil
	}
}

// abpSiteCheckChanged asserts a state attribute no longer has the captured
// value.
func abpSiteCheckChanged(t *testing.T, address, key string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got := abpAttrs(t, s, address)[key]; got == *previous {
			return fmt.Errorf("%s is still %q", key, got)
		}
		return nil
	}
}

// abpSiteMustJson encodes a value, for comparing two of them as strings.
func abpSiteMustJson(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encoding %+v: %s", value, err)
	}
	return string(encoded)
}
