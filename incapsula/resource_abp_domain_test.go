package incapsula

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_domain. Beyond the ordinary create/update/import
// coverage they pin down the two destructive corner cases the resource is built
// around: `criteria` forces replacement, which takes the Domain's encryption
// keys and its slot in the Site's Domain priority order with it, and
// `encryption_key_id` cannot be read back, so neither an import nor a later
// change to it may replace the Domain.

const abpDomainTestAddress = "incapsula_abp_domain.test"

// abpDomainConfig renders a Domain in the fake's Account, with the body given
// verbatim so each test states only what it is about.
func abpDomainConfig(siteId, body string) string {
	return abpDomainConfigIn(abpFakeAccountId, siteId, body)
}

func abpDomainConfigIn(accountId, siteId, body string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_domain" "test" {
  account_id = %q
  site_id    = %q
%s
}
`, accountId, siteId, body)
}

// abpDomainMinimalConfig is the smallest Domain the API accepts.
func abpDomainMinimalConfig(siteId, fqdn string) string {
	return abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = %[1]q

  criteria {
    exact = %[1]q
  }
`, fqdn))
}

// abpDomainStored returns the single Domain the fake holds.
func abpDomainStored(t *testing.T, f *abpFake) map[string]any {
	t.Helper()
	return f.only(t, abpFakeDomains)
}

// --- create -----------------------------------------------------------------

// Every attribute the resource exposes, in one Domain, so that what the
// provider sends is checked against what the API stores and the configuration
// re-plans empty.
func TestAbpDomain_CreateWithEveryAttribute(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("full")
	config := abpDomainConfig(siteId, `
  cookiescope                        = "full.example.com"
  log_region                         = "eu"
  cookie_mode                        = "lax_and_legacy"
  enable_mitigation                  = false
  enable_mobile_sdk_token            = true
  filter_out_static_assets           = true
  obfuscate_path                     = "/obfuscated"
  interstitial_inprogress_iframe_src = "http://www.example.com/iframe"
  divert_host                        = "divert.example.com"
  unmasked_headers                   = ["content-type", "x-request-id"]
  proxy_flags                        = ["enable_referrer_fix", "inject_js_into_body"]

  analysis_ip_lookup_mode {
    header_name   = "X-Forwarded-For"
    reverse_index = 0
  }

  challenge_ip_lookup_mode {
    header_name   = "True-Client-IP"
    reverse_index = 2
  }

  captcha_settings {
    geetest {
      geetest_captcha_id  = "captcha-id"
      geetest_private_key = "captcha-key"
    }
  }

  no_js_injection_path {
    path_prefix = "/static"
  }

  no_js_injection_path {
    incap_rule = "URL == \"/admin\""
  }

  criteria {
    exact = "full.example.com"
  }
`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "criteria.0.exact", "full.example.com"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "log_region", "eu"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "cookie_mode", "lax_and_legacy"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "enable_mitigation", "false"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "captcha_settings.0.geetest.0.geetest_private_key", "captcha-key"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "no_js_injection_path.#", "2"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "analysis_ip_lookup_mode.0.header_name", "X-Forwarded-For"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "challenge_ip_lookup_mode.0.reverse_index", "2"),
					resource.TestCheckResourceAttrSet(abpDomainTestAddress, "created_at"),
					resource.TestCheckResourceAttrSet(abpDomainTestAddress, "modified_at"),
					func(s *terraform.State) error {
						stored := abpDomainStored(t, f)
						if got := abpFakeString(abpFakeMap(stored["criteria"])["exact"]); got != "full.example.com" {
							return fmt.Errorf("stored criteria.exact is %q", got)
						}
						if got := abpFakeString(stored["log_region"]); got != "eu" {
							return fmt.Errorf("stored log_region is %q", got)
						}
						geetest := abpFakeMap(abpFakeMap(stored["captcha_settings"])["geetest"])
						if got := abpFakeString(geetest["geetest_captcha_id"]); got != "captcha-id" {
							return fmt.Errorf("stored geetest_captcha_id is %q", got)
						}
						paths := abpFakeList(stored["no_js_injection_paths"])
						if len(paths) != 2 || abpFakeString(paths[0]["path_prefix"]) != "/static" || abpFakeString(paths[1]["incap_rule"]) != `URL == "/admin"` {
							return fmt.Errorf("stored no_js_injection_paths is %v", paths)
						}
						lookup := abpFakeMap(abpFakeMap(stored["analysis_ip_lookup_mode"])["lookup"])
						if got := abpFakeString(lookup["header_name"]); got != "X-Forwarded-For" {
							return fmt.Errorf("stored analysis_ip_lookup_mode is %v", stored["analysis_ip_lookup_mode"])
						}
						if got := abpFakeStrings(stored["proxy_flags"]); len(got) != 2 {
							return fmt.Errorf("stored proxy_flags is %v", got)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Domain with nothing but the required attributes takes the schema defaults,
// and the API fills in the paths it generates itself.
func TestAbpDomain_CreateMinimal(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("minimal")
	config := abpDomainMinimalConfig(siteId, "minimal.example.com")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "log_region", "usa"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "cookie_mode", "lax"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "enable_mitigation", "true"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "captcha_settings.#", "0"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "analysis_ip_lookup_mode.#", "0"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "no_js_injection_path.#", "0"),
					func(s *terraform.State) error {
						stored := abpDomainStored(t, f)
						// An omitted CAPTCHA and IP lookup mode are the string
						// "none" on the wire, not an absent key.
						if got := abpFakeString(stored["captcha_settings"]); got != "none" {
							return fmt.Errorf("stored captcha_settings is %v, want \"none\"", stored["captcha_settings"])
						}
						if got := abpFakeString(stored["analysis_ip_lookup_mode"]); got != "none" {
							return fmt.Errorf("stored analysis_ip_lookup_mode is %v, want \"none\"", stored["analysis_ip_lookup_mode"])
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Each of the four criteria variants, created and re-planned, because which one
// is present is what the API resolves the Domain's matching on.
func TestAbpDomain_CriteriaVariants(t *testing.T) {
	for _, variant := range []struct {
		name  string
		body  string
		check resource.TestCheckFunc
		key   string
		value any
	}{
		{name: "exact", body: `exact = "one.example.com"`, key: "exact", value: "one.example.com"},
		{name: "prefix", body: `prefix = "two.example.com"`, key: "prefix", value: "two.example.com"},
		{name: "suffix", body: `suffix = "three.example.com"`, key: "suffix", value: "three.example.com"},
		{name: "cloudwaf_id", body: `cloudwaf_id = 12345678`, key: "cloudwaf_id", value: 12345678.0},
	} {
		t.Run(variant.name, func(t *testing.T) {
			f := newAbpFake(t)
			siteId := f.seedSite("criteria-" + variant.name)
			config := abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = "criteria.example.com"

  criteria {
    %s
  }
`, variant.body))

			resource.Test(t, resource.TestCase{
				ProviderFactories: f.providers(t),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: func(s *terraform.State) error {
							criteria := abpFakeMap(abpDomainStored(t, f)["criteria"])
							if len(criteria) != 1 {
								return fmt.Errorf("stored criteria is %v, want exactly one variant", criteria)
							}
							if criteria[variant.key] != variant.value {
								return fmt.Errorf("stored criteria.%s is %v, want %v", variant.key, criteria[variant.key], variant.value)
							}
							return nil
						},
					},
					{Config: config, PlanOnly: true},
				},
			})
		})
	}
}

// The API may report a CloudWAF Domain as `cloudwaf_website_id` plus `fqdn`
// instead of the `cloudwaf_id` a create accepts, so the provider has to read
// both into the one attribute the configuration has.
func TestAbpDomain_CloudwafWebsiteIdIsReadAsCloudwafId(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("cloudwaf")
	seeded := f.seed(abpFakeDomains, map[string]any{
		"site_id":                  siteId,
		"criteria":                 map[string]any{"cloudwaf_website_id": 12345678, "fqdn": "cloudwaf.example.com"},
		"cookiescope":              "cloudwaf.example.com",
		"challenge_ip_lookup_mode": "none",
		"analysis_ip_lookup_mode":  "none",
		"captcha_settings":         "none",
		"log_region":               "usa",
		"cookie_mode":              "lax",
		"no_js_injection_paths":    []any{},
		"unmasked_headers":         []any{},
		"proxy_flags":              []any{},
		"enable_mitigation":        true,
	})
	config := abpDomainConfig(siteId, `
  cookiescope = "cloudwaf.example.com"

  criteria {
    cloudwaf_id = 12345678
  }
`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpDomainTestAddress,
				ImportState:        true,
				ImportStateId:      seeded,
				ImportStatePersist: true,
			},
			{Config: config, PlanOnly: true},
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpDomainTestAddress, "criteria.0.cloudwaf_id", "12345678"),
			},
		},
	})
}

// --- update -----------------------------------------------------------------

// Everything an update can change, changed at once: the Domain keeps its id and
// nothing is deleted.
func TestAbpDomain_UpdateInPlace(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("update")
	before := abpDomainConfig(siteId, `
  cookiescope       = "update.example.com"
  log_region        = "usa"
  cookie_mode       = "lax"
  enable_mitigation = true
  unmasked_headers  = ["content-type"]
  proxy_flags       = ["enable_referrer_fix"]

  analysis_ip_lookup_mode {
    header_name   = "X-Forwarded-For"
    reverse_index = 0
  }

  captcha_settings {
    geetest {
      geetest_captcha_id  = "captcha-id"
      geetest_private_key = "captcha-key"
    }
  }

  no_js_injection_path {
    path_prefix = "/static"
  }

  criteria {
    exact = "update.example.com"
  }
`)
	after := abpDomainConfig(siteId, `
  cookiescope       = "renamed.example.com"
  log_region        = "india"
  cookie_mode       = "none_secure"
  enable_mitigation = false
  unmasked_headers  = ["x-request-id", "x-trace-id"]
  proxy_flags       = ["set_shared_cookie_on_apex"]

  captcha_settings {
    managed_hcaptcha {
      difficulty = "auto"
    }
  }

  no_js_injection_path {
    incap_rule = "URL == \"/admin\""
  }

  criteria {
    exact = "update.example.com"
  }
`)

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					return nil
				},
			},
			{
				Config: after,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "cookiescope", "renamed.example.com"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "analysis_ip_lookup_mode.#", "0"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "captcha_settings.0.managed_hcaptcha.0.difficulty", "auto"),
					func(s *terraform.State) error {
						if id := abpAttrs(t, s, abpDomainTestAddress)["id"]; id != created {
							return fmt.Errorf("the update replaced Domain %s with %s", created, id)
						}
						if deleted := f.deletedIds(); len(deleted) != 0 {
							return fmt.Errorf("the update deleted %v", deleted)
						}
						stored := abpDomainStored(t, f)
						// A oneOf that switched variant must not keep the
						// leftovers of the old one: the API resolves the
						// variant by key presence.
						captcha := abpFakeMap(stored["captcha_settings"])
						if _, ok := captcha["geetest"]; ok {
							return fmt.Errorf("stored captcha_settings still carries the geetest variant: %v", captcha)
						}
						if abpFakeString(stored["analysis_ip_lookup_mode"]) != "none" {
							return fmt.Errorf("stored analysis_ip_lookup_mode is %v, want \"none\"", stored["analysis_ip_lookup_mode"])
						}
						paths := abpFakeList(stored["no_js_injection_paths"])
						if len(paths) != 1 || abpFakeString(paths[0]["incap_rule"]) == "" {
							return fmt.Errorf("stored no_js_injection_paths is %v", paths)
						}
						if got := abpFakeStrings(stored["unmasked_headers"]); len(got) != 2 {
							return fmt.Errorf("stored unmasked_headers is %v", got)
						}
						return nil
					},
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// The update endpoint is a full replace that does not carry `criteria` or
// `mobile_api_obfuscate_path`, so an unrelated update must send them back
// untouched instead of clearing them — along with any key the provider does not
// model at all.
func TestAbpDomain_UpdateKeepsFieldsTheUpdateDoesNotCarry(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("merge")
	before := abpDomainMinimalConfig(siteId, "merge.example.com")
	after := abpDomainConfig(siteId, `
  cookiescope = "merge.example.com"
  log_region  = "australia"

  criteria {
    exact = "merge.example.com"
  }
`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: before,
				Check: func(s *terraform.State) error {
					f.mutate(t, abpFakeDomains, abpAttrs(t, s, abpDomainTestAddress)["id"], func(doc map[string]any) {
						doc["mobile_api_obfuscate_path"] = "/mobile/challenge"
						doc["a_field_the_provider_does_not_model"] = "keep me"
					})
					return nil
				},
			},
			{
				Config: after,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "mobile_api_obfuscate_path", "/mobile/challenge"),
					func(s *terraform.State) error {
						stored := abpDomainStored(t, f)
						if got := abpFakeString(stored["mobile_api_obfuscate_path"]); got != "/mobile/challenge" {
							return fmt.Errorf("the update cleared mobile_api_obfuscate_path, leaving %q", got)
						}
						if got := abpFakeString(stored["a_field_the_provider_does_not_model"]); got != "keep me" {
							return fmt.Errorf("the update cleared an unmodelled field, leaving %q", got)
						}
						if got := abpFakeString(abpFakeMap(stored["criteria"])["exact"]); got != "merge.example.com" {
							return fmt.Errorf("the update lost criteria.exact, leaving %q", got)
						}
						return nil
					},
				),
			},
			{Config: after, PlanOnly: true},
		},
	})
}

// Moving a Domain to another Site of the same Account is an in-place update:
// the Domain keeps its id, its encryption keys and the configuration that
// references it.
func TestAbpDomain_SiteIdChangeIsInPlace(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("first")
	second := f.seedSite("second")

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpDomainMinimalConfig(first, "move.example.com"),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpDomainMinimalConfig(second, "move.example.com"),
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpDomainTestAddress)["id"]; id != created {
						return fmt.Errorf("moving the Domain to another Site replaced %s with %s", created, id)
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("moving the Domain deleted %v", deleted)
					}
					if got := abpFakeString(abpDomainStored(t, f)["site_id"]); got != second {
						return fmt.Errorf("stored site_id is %q, want %q", got, second)
					}
					return nil
				},
			},
			{Config: abpDomainMinimalConfig(second, "move.example.com"), PlanOnly: true},
		},
	})
}

// --- replacement ------------------------------------------------------------

// `criteria` is not settable by an update, so changing it replaces the Domain.
// The replacement is what makes a non-converging criteria value destructive,
// which is why a padded one is rejected at plan time.
func TestAbpDomain_CriteriaChangeForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("replace")

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpDomainMinimalConfig(siteId, "before.example.com"),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpDomainConfig(siteId, `
  cookiescope = "before.example.com"

  criteria {
    suffix = "after.example.com"
  }
`),
				Check: func(s *terraform.State) error {
					id := abpAttrs(t, s, abpDomainTestAddress)["id"]
					if id == created {
						return fmt.Errorf("changing criteria did not replace Domain %s", created)
					}
					if f.exists(abpFakeDomains, created) {
						return fmt.Errorf("the replaced Domain %s is still stored", created)
					}
					if deleted := f.deletedIds(); len(deleted) != 1 || deleted[0] != created {
						return fmt.Errorf("the replacement deleted %v, want just %s", deleted, created)
					}
					return nil
				},
			},
		},
	})
}

// A Domain cannot be moved to another Account, so changing `account_id`
// replaces it.
func TestAbpDomain_AccountIdChangeForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("accounts")
	body := `
  cookiescope = "accounts.example.com"

  criteria {
    exact = "accounts.example.com"
  }
`

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpDomainConfigIn(abpFakeAccountId, siteId, body),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpDomainConfigIn(abpFakeOtherAccountId, siteId, body),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "account_id", abpFakeOtherAccountId),
					func(s *terraform.State) error {
						if id := abpAttrs(t, s, abpDomainTestAddress)["id"]; id == created {
							return fmt.Errorf("changing account_id did not replace Domain %s", created)
						}
						if deleted := f.deletedIds(); len(deleted) != 1 || deleted[0] != created {
							return fmt.Errorf("the replacement deleted %v, want just %s", deleted, created)
						}
						return nil
					},
				),
			},
		},
	})
}

// --- drift ------------------------------------------------------------------

// A Domain edited outside Terraform is brought back to the configuration.
func TestAbpDomain_DriftIsCorrected(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("drift")
	other := f.seedSite("drift-other")
	config := abpDomainConfig(siteId, `
  cookiescope      = "drift.example.com"
  log_region       = "eu"
  unmasked_headers = ["content-type"]

  no_js_injection_path {
    path_prefix = "/static"
  }

  criteria {
    exact = "drift.example.com"
  }
`)

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakeDomains, created, func(doc map[string]any) {
						doc["cookiescope"] = "elsewhere.example.com"
						doc["log_region"] = "apac"
						doc["site_id"] = other
						doc["unmasked_headers"] = []any{"x-smuggled"}
						// An element added to a list the provider models.
						doc["no_js_injection_paths"] = []any{
							map[string]any{"path_prefix": "/static"},
							map[string]any{"path_prefix": "/added-elsewhere"},
						}
					})
				},
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "cookiescope", "drift.example.com"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "log_region", "eu"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "site_id", siteId),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "no_js_injection_path.#", "1"),
					func(s *terraform.State) error {
						stored := abpDomainStored(t, f)
						if got := abpFakeStrings(stored["unmasked_headers"]); len(got) != 1 || got[0] != "content-type" {
							return fmt.Errorf("stored unmasked_headers is %v", got)
						}
						if got := abpFakeList(stored["no_js_injection_paths"]); len(got) != 1 {
							return fmt.Errorf("stored no_js_injection_paths is %v", got)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Domain deleted outside Terraform is recreated, and so is the encryption key
// the Domain took with it — the API deletes a Domain's keys along with it, which
// is why nothing about `criteria` or `encryption_key_id` may replace a Domain.
func TestAbpDomain_DeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("vanished")
	config := abpDomainMinimalConfig(siteId, "vanished.example.com") + fmt.Sprintf(`
resource "incapsula_abp_domain_encryption_key" "test" {
  domain_id = incapsula_abp_domain.test.id
  key       = %q
}
`, abpDomainTestKeyMaterial)

	var created, createdKey string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					createdKey = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.dropDomainWithKeys(created) },
				Config:    config,
				Check: func(s *terraform.State) error {
					id := abpAttrs(t, s, abpDomainTestAddress)["id"]
					if id == created {
						return fmt.Errorf("the vanished Domain %s was not recreated", created)
					}
					if !f.exists(abpFakeDomains, id) {
						return fmt.Errorf("Domain %s is not stored", id)
					}
					if key := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]; key == createdKey {
						return fmt.Errorf("the encryption key %s of the vanished Domain was not recreated", createdKey)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Known limitation: the API omits an optional field that has no value instead
// of returning it as empty, and the provider only writes such a field into
// state when the response carries it. Turning one off outside Terraform is
// therefore invisible — the plan is empty while state and server disagree.
// This pins the behaviour that exists today, not the one that is wanted.
func TestAbpDomain_ClearedOptionalFieldsAreInvisible(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("invisible")
	config := abpDomainConfig(siteId, `
  cookiescope                        = "invisible.example.com"
  enable_mobile_sdk_token            = true
  filter_out_static_assets           = true
  divert_host                        = "divert.example.com"
  interstitial_inprogress_iframe_src = "http://www.example.com/iframe"

  criteria {
    exact = "invisible.example.com"
  }
`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					f.mutate(t, abpFakeDomains, abpAttrs(t, s, abpDomainTestAddress)["id"], func(doc map[string]any) {
						for _, key := range []string{
							"enable_mobile_sdk_token",
							"filter_out_static_assets",
							"divert_host",
							"interstitial_inprogress_iframe_src",
						} {
							delete(doc, key)
						}
					})
					return nil
				},
			},
			// No step asks for a change, because the provider cannot see one.
			{Config: config, PlanOnly: true},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "enable_mobile_sdk_token", "true"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "divert_host", "divert.example.com"),
					func(s *terraform.State) error {
						if _, ok := abpDomainStored(t, f)["divert_host"]; ok {
							return fmt.Errorf("the fake was expected to still hold no divert_host")
						}
						return nil
					},
				),
			},
		},
	})
}

// Known limitation: `enable_mobile_sdk_token` and `filter_out_static_assets`
// are tri-state on the API (no value / true / false), and no value is
// meaningful — for `filter_out_static_assets` it means enabled. Terraform
// cannot express "no value": `false` is sent as `false`, and removing the
// attribute is planned as `-> null` but sent as `false` too.
func TestAbpDomain_OptionalBooleansAreSentAsFalse(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("booleans")
	withFalse := abpDomainConfig(siteId, `
  cookiescope              = "booleans.example.com"
  enable_mobile_sdk_token  = false
  filter_out_static_assets = false

  criteria {
    exact = "booleans.example.com"
  }
`)
	withoutThem := abpDomainMinimalConfig(siteId, "booleans.example.com")

	storedBooleans := func(s *terraform.State) error {
		stored := abpDomainStored(t, f)
		for _, key := range []string{"enable_mobile_sdk_token", "filter_out_static_assets"} {
			value, ok := stored[key]
			if !ok {
				return fmt.Errorf("%s is absent from the stored Domain, want false", key)
			}
			if value != false {
				return fmt.Errorf("stored %s is %v, want false", key, value)
			}
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: withFalse,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "enable_mobile_sdk_token", "false"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "filter_out_static_assets", "false"),
					storedBooleans,
				),
			},
			{Config: withFalse, PlanOnly: true},
			{
				// Removing them is planned as `-> null`, and stored as false.
				Config: withoutThem,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainTestAddress, "enable_mobile_sdk_token", "false"),
					resource.TestCheckResourceAttr(abpDomainTestAddress, "filter_out_static_assets", "false"),
					storedBooleans,
				),
			},
			{Config: withoutThem, PlanOnly: true},
		},
	})
}

// --- delete -----------------------------------------------------------------

// Destroying the configuration deletes the Domain, and the encryption key on it
// is deleted first because it depends on the Domain.
func TestAbpDomain_DestroyDeletesDomainAndItsKey(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("destroy")
	config := abpDomainMinimalConfig(siteId, "destroy.example.com") + fmt.Sprintf(`
resource "incapsula_abp_domain_encryption_key" "test" {
  domain_id = incapsula_abp_domain.test.id
  key       = %q
}
`, abpDomainTestKeyMaterial)

	var created, createdKey string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if f.exists(abpFakeDomains, created) {
				return fmt.Errorf("Domain %s survived the destroy", created)
			}
			if deleted := f.deletedIds(); len(deleted) != 2 || deleted[0] != createdKey || deleted[1] != created {
				return fmt.Errorf("the destroy deleted %v, want the key %s then the Domain %s", deleted, createdKey, created)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					createdKey = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
		},
	})
}

// --- import -----------------------------------------------------------------

func TestAbpDomain_Import(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("import")
	config := abpDomainConfig(siteId, `
  cookiescope             = "import.example.com"
  log_region              = "apac"
  cookie_mode             = "legacy"
  enable_mobile_sdk_token = true
  obfuscate_path          = "/obfuscated"
  unmasked_headers        = ["content-type"]
  proxy_flags             = ["inject_js_into_body"]

  challenge_ip_lookup_mode {
    header_name   = "X-Forwarded-For"
    reverse_index = 1
  }

  captcha_settings {
    managed_geetest {
      difficulty = "normal"
    }
  }

  no_js_injection_path {
    path_prefix = "/static"
  }

  criteria {
    prefix = "import.example.com"
  }
`)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpDomainTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAbpDomain_ImportErrors(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("import-errors")
	config := abpDomainMinimalConfig(siteId, "import-errors.example.com")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:        config,
				ResourceName:  abpDomainTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Domain ` + abpFakeAbsentId + ` not found`),
			},
			{
				Config:        config,
				ResourceName:  abpDomainTestAddress,
				ImportState:   true,
				ImportStateId: " ",
				ExpectError:   regexp.MustCompile(`expected import ID to be '<domain_id>'`),
			},
		},
	})
}

// `encryption_key_id` is consumed on create and never returned, so an imported
// Domain never has one. The first apply after an import must keep the Domain:
// replacing it would destroy its encryption keys and its slot in the Site's
// Domain priority order.
func TestAbpDomain_ImportThenApplyKeepsDomain(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("import-keeps")
	seeded := f.seedDomain(siteId, "seeded.example.com", abpDomainTestKeyId)
	config := abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope       = "seeded.example.com"
  encryption_key_id = %q

  criteria {
    exact = "seeded.example.com"
  }
`, abpDomainTestKeyId))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpDomainTestAddress,
				ImportState:        true,
				ImportStateId:      seeded,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpDomainTestAddress)["id"]; id != seeded {
						return fmt.Errorf("the imported Domain %s was replaced by %s", seeded, id)
					}
					if !f.exists(abpFakeDomains, seeded) {
						return fmt.Errorf("the imported Domain %s was deleted", seeded)
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the first apply after import deleted %v", deleted)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// `encryption_key_id` is honoured on create only, so changing it afterwards
// must neither replace the Domain nor leave a permanent diff.
func TestAbpDomain_EncryptionKeyIdIsCreateOnly(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("create-only")
	config := func(keyId string) string {
		return abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope       = "create-only.example.com"
  encryption_key_id = %q

  criteria {
    exact = "create-only.example.com"
  }
`, keyId))
	}

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config(abpDomainTestKeyId),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainTestAddress)["id"]
					if got := f.encryptionKeyIdOf(created); got != abpDomainTestKeyId {
						return fmt.Errorf("the Domain was created with encryption_key_id %q, want %q", got, abpDomainTestKeyId)
					}
					if _, ok := abpDomainStored(t, f)["encryption_key_id"]; ok {
						return fmt.Errorf("encryption_key_id was expected to be consumed rather than stored")
					}
					return nil
				},
			},
			// The new value is suppressed, so there is nothing to apply.
			{Config: config(abpDomainTestOtherKeyId), PlanOnly: true},
			{
				Config: config(abpDomainTestOtherKeyId),
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpDomainTestAddress)["id"]; id != created {
						return fmt.Errorf("changing encryption_key_id replaced Domain %s with %s", created, id)
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("changing encryption_key_id deleted %v", deleted)
					}
					return nil
				},
			},
		},
	})
}

// --- plan-time validation ---------------------------------------------------

// The API trims the criteria domain name parts, and `criteria` forces
// replacement, so a padded value would make every apply destroy the Domain —
// with its encryption keys and its priority slot — and create a new one,
// forever. It is rejected while the configuration is planned instead.
func TestAbpDomain_PaddedCriteriaIsRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("padded-criteria")
	clean := abpDomainMinimalConfig(siteId, "padded.example.com")
	padded := func(criteria string) string {
		return abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = "padded.example.com"

  criteria {
    %s
  }
`, criteria))
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      padded(`exact = "  padded.example.com  "`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{
				Config:      padded(`prefix = "padded.example.com "`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{
				Config:      padded(`suffix = "\tpadded.example.com"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{
				Config:      padded(`exact = "   "`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not consist only of whitespace`),
			},
			{Config: clean},
			{Config: clean, PlanOnly: true},
		},
	})
}

// `cookiescope` is trimmed by the API too. It is not part of the Domain's
// identity, so a padded value only leaves a diff no apply can resolve, but it
// is rejected at plan time for the same reason.
func TestAbpDomain_PaddedCookiescopeIsRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("padded-cookiescope")
	clean := abpDomainMinimalConfig(siteId, "cookiescope.example.com")
	padded := func(cookiescope string) string {
		return abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = %q

  criteria {
    exact = "cookiescope.example.com"
  }
`, cookiescope))
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      padded("  cookiescope.example.com  "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{
				Config:      padded(" "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not have leading or trailing whitespace`),
			},
			{Config: clean},
			{Config: clean, PlanOnly: true},
		},
	})
}

// The values the schema restricts, and the references it requires to be UUIDs,
// are all rejected before anything is created.
func TestAbpDomain_InvalidValuesAreRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("invalid")

	steps := []struct {
		body  string
		error string
	}{
		{`log_region = "antarctica"`, `expected log_region to be one of`},
		{`cookie_mode = "strict"`, `expected cookie_mode to be one of`},
		{`proxy_flags = ["make_it_fast"]`, `expected proxy_flags.\d+ to be one of`},
		{"captcha_settings {\n    managed_geetest {\n      difficulty = \"impossible\"\n    }\n  }", `expected captcha_settings.0.managed_geetest.0.difficulty to be one of`},
		{"captcha_settings {\n    managed_hcaptcha {\n      difficulty = \"hard\"\n    }\n  }", `expected captcha_settings.0.managed_hcaptcha.0.difficulty to be one of`},
		{"analysis_ip_lookup_mode {\n    header_name   = \"X-Forwarded-For\"\n    reverse_index = -1\n  }", `expected analysis_ip_lookup_mode.0.reverse_index to be at least \(0\)`},
		{"criteria {\n    exact = \"second.example.com\"\n  }", `No more than 1 "criteria" blocks are allowed`},
	}

	var cases []resource.TestStep
	for _, step := range steps {
		cases = append(cases, resource.TestStep{
			Config: abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = "invalid.example.com"
  %s

  criteria {
    exact = "invalid.example.com"
  }
`, step.body)),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(step.error),
		})
	}
	cases = append(cases, resource.TestStep{
		Config:      abpDomainConfigIn("not-a-uuid", siteId, "\n  cookiescope = \"invalid.example.com\"\n\n  criteria {\n    exact = \"invalid.example.com\"\n  }\n"),
		PlanOnly:    true,
		ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
	})
	cases = append(cases, resource.TestStep{
		Config:      abpDomainConfig("not-a-uuid", "\n  cookiescope = \"invalid.example.com\"\n\n  criteria {\n    exact = \"invalid.example.com\"\n  }\n"),
		PlanOnly:    true,
		ExpectError: regexp.MustCompile(`expected "site_id" to be a valid UUID`),
	})

	resource.Test(t, resource.TestCase{ProviderFactories: f.providers(t), Steps: cases})
}

// Known limitation: the "exactly one of" rules live in the create and update
// paths rather than in validation, so these configurations plan cleanly and
// fail during apply. Nothing is created, because the checks run before the
// request is sent.
func TestAbpDomain_OneOfViolationsAreRejectedAtApplyTime(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("one-of")

	bodies := []struct {
		body  string
		error string
	}{
		{"criteria {\n  }", `criteria requires exactly one of exact, prefix, suffix, cloudwaf_id \(got 0\)`},
		{"criteria {\n    exact  = \"one-of.example.com\"\n    suffix = \"example.com\"\n  }", `criteria requires exactly one of exact, prefix, suffix, cloudwaf_id \(got 2\)`},
		{"criteria {\n    exact = \"one-of.example.com\"\n  }\n\n  captcha_settings {\n    managed_geetest {\n      difficulty = \"easy\"\n    }\n    managed_hcaptcha {\n      difficulty = \"auto\"\n    }\n  }", `captcha_settings requires exactly one of geetest, managed_geetest, managed_hcaptcha \(got 2\)`},
		{"criteria {\n    exact = \"one-of.example.com\"\n  }\n\n  no_js_injection_path {\n    path_prefix = \"/static\"\n    incap_rule  = \"URL == \\\"/admin\\\"\"\n  }", `no_js_injection_path\[0\]: exactly one of path_prefix, incap_rule must be set \(got 2\)`},
	}

	var steps []resource.TestStep
	for _, one := range bodies {
		steps = append(steps, resource.TestStep{
			Config: abpDomainConfig(siteId, fmt.Sprintf(`
  cookiescope = "one-of.example.com"

  %s
`, one.body)),
			ExpectError: regexp.MustCompile(one.error),
		})
	}

	resource.Test(t, resource.TestCase{ProviderFactories: f.providers(t), Steps: steps})

	if f.count(abpFakeDomains) != 0 {
		t.Fatalf("a rejected configuration created %d Domains", f.count(abpFakeDomains))
	}
	if f.requestCount(http.MethodPost, "/domain") != 0 {
		t.Fatal("a rejected configuration reached the API")
	}
}

// Provider bug: an empty `no_js_injection_path` block reaches the provider as a
// nil element, and the extraction asserts it is a map, so the provider panics
// instead of reporting the "exactly one of" error the other variants get. The
// crash takes the whole provider process with it, so it is pinned here at the
// function rather than through a configuration.
func TestAbpDomain_EmptyNoJsInjectionPathPanics(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("extractAbpNoJsInjectionPaths no longer panics on an empty block; assert the error instead")
		}
	}()
	//nolint:errcheck // the call is expected to panic before it returns
	_, _ = extractAbpNoJsInjectionPaths([]any{nil})
}

// A Domain whose Site does not exist is rejected by the API, with the message
// the API produced rather than a provider crash.
func TestAbpDomain_UnknownSiteIsRejectedByTheApi(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpDomainMinimalConfig(abpFakeAbsentId, "unknown-site.example.com"),
				ExpectError: regexp.MustCompile(`Unknown site ` + abpFakeAbsentId),
			},
		},
	})
}

// --- helpers ----------------------------------------------------------------

// abpDomainWritesTo returns the methods of the writes whose path contains the
// given fragment, in order.
func abpDomainWritesTo(f *abpFake, fragment string) []string {
	methods := []string{}
	for _, write := range f.writes() {
		if strings.Contains(write.Path, fragment) {
			methods = append(methods, write.Method)
		}
	}
	return methods
}
