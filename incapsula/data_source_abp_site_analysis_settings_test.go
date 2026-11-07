package incapsula

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for the incapsula_abp_site_analysis_settings data source. It builds
// the JSON document a Selector's `analysis_settings` holds, and because that
// value is compared as a string, the exact output is part of the contract.

const abpSiteSettingsAddress = "data.incapsula_abp_site_analysis_settings.test"

// abpSiteSettingsConfig builds the data source from its arguments.
func abpSiteSettingsConfig(arguments string) string {
	return fmt.Sprintf(`
data "incapsula_abp_site_analysis_settings" "test" {
%s
}
`, arguments)
}

// Each rate limiting mode produces its own shape: `none` and `per_site` are
// written as the mode name, `custom_scope` as the object the API's oneOf
// expects.
func TestAbpSiteAnalysisSettingsDataSourceEveryRateLimitingMode(t *testing.T) {
	f := newAbpFake(t)
	for _, test := range []struct {
		arguments string
		json      string
	}{
		{
			arguments: `  rate_limiting = "none"`,
			json:      `{"rate_limiting":"none","use_site_rate_limiting_parameters":true}`,
		},
		{
			arguments: `  rate_limiting = "per_site"`,
			json:      `{"rate_limiting":"per_site","use_site_rate_limiting_parameters":true}`,
		},
		{
			arguments: `  rate_limiting              = "custom_scope"
  rate_limiting_custom_scope = "checkout"`,
			json: `{"rate_limiting":{"custom_scope":"checkout"},"use_site_rate_limiting_parameters":true}`,
		},
	} {
		resource.Test(t, resource.TestCase{
			ProviderFactories: f.providers(t),
			Steps: []resource.TestStep{
				{
					Config: abpSiteSettingsConfig(test.arguments),
					Check:  resource.TestCheckResourceAttr(abpSiteSettingsAddress, "json", test.json),
				},
			},
		})
	}
}

// The per-Selector maxima are only written next to
// `use_site_rate_limiting_parameters = false`, which is the setting that makes
// them take effect at all.
func TestAbpSiteAnalysisSettingsDataSourceCarriesTheMaxima(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteSettingsConfig(`  rate_limiting = "per_site"

  max_requests_per_minute  = 55
  max_requests_per_session = 555
  max_session_length       = "2d1h"

  use_site_rate_limiting_parameters = false`),
				Check: resource.TestCheckResourceAttr(abpSiteSettingsAddress, "json",
					`{"rate_limiting":"per_site","max_requests_per_minute":55,"max_requests_per_session":555,"max_session_length":"2d1h","use_site_rate_limiting_parameters":false}`),
			},
			{
				// The maxima are written even when the flag that ignores them
				// keeps its default, which is how dead configuration gets
				// written in the first place.
				Config: abpSiteSettingsConfig(`  rate_limiting           = "per_site"
  max_requests_per_minute = 55`),
				Check: resource.TestCheckResourceAttr(abpSiteSettingsAddress, "json",
					`{"rate_limiting":"per_site","max_requests_per_minute":55,"use_site_rate_limiting_parameters":true}`),
			},
		},
	})
}

// A custom scope is required by, and only by, `rate_limiting =
// "custom_scope"`.
func TestAbpSiteAnalysisSettingsDataSourceValidatesTheCustomScope(t *testing.T) {
	f := newAbpFake(t)
	for name, test := range map[string]struct{ arguments, message string }{
		"custom_scope without a scope": {
			`  rate_limiting = "custom_scope"`,
			`rate_limiting_custom_scope must be set when rate_limiting = "custom_scope"`,
		},
		"a scope with none": {
			`  rate_limiting              = "none"
  rate_limiting_custom_scope = "checkout"`,
			`rate_limiting_custom_scope must not be set when rate_limiting = "none"`,
		},
		"a scope with per_site": {
			`  rate_limiting              = "per_site"
  rate_limiting_custom_scope = "checkout"`,
			`rate_limiting_custom_scope must not be set when rate_limiting = "per_site"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{
						Config:      abpSiteSettingsConfig(test.arguments),
						ExpectError: regexp.MustCompile(test.message),
					},
				},
			})
		})
	}
}

// The arguments are checked while the configuration is planned.
func TestAbpSiteAnalysisSettingsDataSourceInvalidArgumentsAreRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	for name, test := range map[string]struct{ arguments, message string }{
		"missing rate_limiting": {
			`  max_requests_per_minute = 10`,
			`The argument "rate_limiting" is required`,
		},
		"unknown rate_limiting": {
			`  rate_limiting = "per_selector"`,
			`expected rate_limiting to be one of`,
		},
		"negative max_requests_per_minute": {
			`  rate_limiting           = "per_site"
  max_requests_per_minute = -1`,
			`expected max_requests_per_minute to be at least \(0\)`,
		},
		"negative max_requests_per_session": {
			`  rate_limiting            = "per_site"
  max_requests_per_session = -1`,
			`expected max_requests_per_session to be at least \(0\)`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{
						Config:      abpSiteSettingsConfig(test.arguments),
						ExpectError: regexp.MustCompile(test.message),
					},
				},
			})
		})
	}
}

// A maximum of `0` passes validation but is left out of the document, so it
// means "no override" rather than "zero". This is a limitation (ABP-REVIEW
// `incapsula_abp_site` issue 5), not an intention.
func TestAbpSiteAnalysisSettingsDataSourceDropsAZeroMaximum(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteSettingsConfig(`  rate_limiting = "per_site"

  max_requests_per_minute  = 0
  max_requests_per_session = 0

  use_site_rate_limiting_parameters = false`),
				Check: resource.TestCheckResourceAttr(abpSiteSettingsAddress, "json",
					`{"rate_limiting":"per_site","use_site_rate_limiting_parameters":false}`),
			},
		},
	})
}

// A duration is passed through verbatim: the data source does not normalise
// it, and the API stores it in its own canonical form, so a
// non-canonical duration leaves a Selector that never converges. A limitation
// (ABP-REVIEW `incapsula_abp_site` issue 3) worth pinning down, because the
// value reaching the Site is the one written here.
func TestAbpSiteAnalysisSettingsDataSourceDoesNotNormaliseDurations(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteSettingsConfig(`  rate_limiting      = "per_site"
  max_session_length = "120m"

  use_site_rate_limiting_parameters = false`),
				Check: resource.TestCheckResourceAttr(abpSiteSettingsAddress, "json",
					`{"rate_limiting":"per_site","max_session_length":"120m","use_site_rate_limiting_parameters":false}`),
			},
		},
	})
}

// The id is the digest of the document, so settings that differ are different
// data sources and identical settings are the same one.
func TestAbpSiteAnalysisSettingsDataSourceIdIsTheDigestOfTheDocument(t *testing.T) {
	f := newAbpFake(t)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpSiteSettingsConfig(`  rate_limiting = "per_site"`) + `
data "incapsula_abp_site_analysis_settings" "same" {
  rate_limiting = "per_site"
}

data "incapsula_abp_site_analysis_settings" "other" {
  rate_limiting = "none"
}
`,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						attributes := abpAttrs(t, s, abpSiteSettingsAddress)
						digest := sha256.Sum256([]byte(attributes["json"]))
						if want := hex.EncodeToString(digest[:]); attributes["id"] != want {
							return fmt.Errorf("id is %q, want the digest %q of %q", attributes["id"], want, attributes["json"])
						}
						return nil
					},
					resource.TestCheckResourceAttrPair(abpSiteSettingsAddress, "id",
						"data.incapsula_abp_site_analysis_settings.same", "id"),
					func(s *terraform.State) error {
						mine := abpAttrs(t, s, abpSiteSettingsAddress)["id"]
						other := abpAttrs(t, s, "data.incapsula_abp_site_analysis_settings.other")["id"]
						if mine == other {
							return fmt.Errorf("different settings share the id %q", mine)
						}
						return nil
					},
				),
			},
		},
	})
}
