terraform {
  required_providers {
    incapsula = {
      source = "imperva/incapsula"
    }
  }
}

variable "account_id" {
  type = string
}

#
# Create a condition with custom MOI
#

resource "incapsula_abp_condition" "specific_visitor" {
  account_id  = var.account_id
  name        = "Specific visitor"
  description = "Match specific UA coming from selected IPs"
  code        = "(all headers.user_agent? (matches headers.user_agent re\"Mozilla\") (in visitor_ip 1.2.3.4 2.3.4.5))"
}

# Demonstrate condition lookup
data "incapsula_abp_condition" "specific_visitor_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_condition.specific_visitor.name
}

resource "incapsula_abp_condition" "okhttp" {
  account_id  = var.account_id
  name        = "OkHttp"
  description = "Matches requests initiated by okhttp"
  code        = "(all headers.user_agent? (matches headers.user_agent \"okhttp/4.12.0\"))"

}

#
# Create a condition from an editor template
#

# `template` makes the server validate `code` against that template's expected
# structure instead of accepting an arbitrary expression. Omitting it is the
# same as `custom`.
#
# NOTE: the template is fixed when the condition is created. Changing it here
# reports a successful apply while nothing changes on the server, so the plan
# never converges; recreate the condition instead
# (`terraform apply -replace=module.abp.incapsula_abp_condition.too_many_requests`).
resource "incapsula_abp_condition" "too_many_requests" {
  account_id  = var.account_id
  name        = "Too many requests"
  description = "Trip when a session exceeds 120 requests per minute"
  template    = "rate_limiting"
  code        = "(all (not flags.no_token) (requests_per_minute > 120))"
}

#
# Lookup a managed condition to subsequently insert into a policy
#

data "incapsula_abp_condition" "managed_monitoring_tools" {
  account_id = var.account_id
  name       = "Monitoring Tools"
  managed    = true
}

#
# Create a condition list and populate it with a condition
#

resource "incapsula_abp_condition_list" "sample_condition_list" {
  account_id  = var.account_id
  name        = "Sample condition list"
  description = "Reusable condition list"
}

resource "incapsula_abp_condition_list_entry" "sample_condition_list_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_condition_list.sample_condition_list.id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

# Demonstrate condition list lookup
data "incapsula_abp_condition_list" "sample_condition_list_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_condition_list.sample_condition_list.name
}

#
# List the conditions of the account
#

# Both the account-owned and the managed conditions, for discovery and for
# building lookups. Note that `managed` means something different here than on
# the singular `incapsula_abp_condition` data source: here it selects *only*
# managed or *only* account-owned conditions, while there `managed = false`
# matches either.
data "incapsula_abp_conditions" "account_owned" {
  account_id = var.account_id
}

# Managed conditions are maintained by Imperva. Deprecated and hidden ones are
# listed too, with no attribute to tell them apart.
data "incapsula_abp_conditions" "managed" {
  account_id = var.account_id
  managed    = true
}

#
# Create a policy and populate it with conditions
# using`incapsula_abp_directive` data source
#
#

# `description` is optional: omitting it leaves the policy without one.
# `policy2` below shows the other case.
resource "incapsula_abp_policy" "policy1" {
  account_id = var.account_id
  name       = "Policy with explicit directives"

  directive {
    action = "allow"
  }

  directive {
    action = "block"
  }

  directive {
    action = "captcha_cleared"
  }

  directive {
    action = "captcha"
  }

  directive {
    action = "identify"
  }

  directive {
    action = "tarpit"
  }

  directive {
    action = "delay"
  }
}

# Demonstrate referencing a directive by action with the `incapsula_abp_directive` data source.
# Always address directives this way, never by index (see the WARNING on `policy2` below).
data "incapsula_abp_directive" "policy1_block" {
  policy_id = incapsula_abp_policy.policy1.id
  action    = "block"
}

resource "incapsula_abp_condition_list_entry" "policy1_block_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.policy1_block.condition_list_id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

#
# Create a policy with a proof_of_work directive and populate its directives
# with conditions
#

resource "incapsula_abp_proof_of_work_configuration" "pow1" {
  account_id = var.account_id
  name       = "terraform-pow-0"
  difficulty = 42
  algorithm  = "bbs"
}

# Demonstrate proof_of_work lookup
data "incapsula_abp_proof_of_work_configuration" "pow1_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_proof_of_work_configuration.pow1.name
}

resource "incapsula_abp_policy" "policy2" {
  account_id  = var.account_id
  name        = "Policy with custom directives"
  description = "Demonstrate how to create policy with custom directives from Terraform"

  directive {
    action = "allow"
  }

  directive {
    action = "block"
  }

  directive {
    action                         = "proof_of_work"
    proof_of_work_configuration_id = incapsula_abp_proof_of_work_configuration.pow1.id
  }
}

# Demonstrate policy lookup
data "incapsula_abp_policy" "policy2" {
  account_id = var.account_id
  name       = incapsula_abp_policy.policy2.name
}

# List the policies of the account.
#
# Note the scope: the account global policy is not part of this list (look it up
# with `incapsula_abp_policy` and `account_global = true`, as below), while the
# `"<site> Default"` policies the backend creates for every site are, with no
# attribute to tell them apart from the ones this configuration owns.
data "incapsula_abp_policies" "all" {
  account_id = var.account_id
}

# Address each directive by its action.
#
# WARNING: never index directives (`incapsula_abp_policy.foo.directive[0]`). The index
# moves when directives are added, removed or reordered, so conditions end up on the
# wrong directive. The nested `condition_list_id` is unreliable for the same reason:
# on the apply that changes the directive list it is still reported at its old
# position. This data source re-reads it from the API and is not affected.
data "incapsula_abp_directive" "policy2_allow" {
  policy_id = incapsula_abp_policy.policy2.id
  action    = "allow"
}

data "incapsula_abp_directive" "policy2_block" {
  policy_id = incapsula_abp_policy.policy2.id
  action    = "block"
}

data "incapsula_abp_directive" "policy2_pow" {
  policy_id = incapsula_abp_policy.policy2.id
  action    = "proof_of_work"
}

# Add condition to the `allow` directive
resource "incapsula_abp_condition_list_entry" "policy2_allow_monitoring_tools" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.policy2_allow.condition_list_id
  condition_id             = data.incapsula_abp_condition.managed_monitoring_tools.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

# Skip the `proof_of_work` directive for the managed monitoring tools condition
resource "incapsula_abp_condition_list_entry" "policy2_pow_skip_monitoring_tools" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.policy2_pow.skip_condition_list_id
  condition_id             = data.incapsula_abp_condition.managed_monitoring_tools.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

# Add a whole condition list to the `block` directive
resource "incapsula_abp_condition_list_entry" "policy2_block_sample_condition_list" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.policy2_block.condition_list_id
  condition_list_id        = incapsula_abp_condition_list.sample_condition_list.id
  state                    = "monitor"
  tags                     = ["terraform_managed"]
}

#
# Lookup the account global policy and populate it with conditions
#

# Lookup the account global policy.
# Note: the account global policy can't be directly referenced by the returned policy id,
# for example in `incapsula_abp_directive.policy_id`. Use `account_global_policy = true`
# on `incapsula_abp_directive` instead (see below).
data "incapsula_abp_policy" "account_global" {
  account_id     = var.account_id
  account_global = true
}

data "incapsula_abp_directive" "account_global_allow" {
  account_id            = var.account_id
  account_global_policy = true
  action                = "allow"
}

# Now add conditions to the automatically created directives
resource "incapsula_abp_condition_list_entry" "account_global_allow_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.account_global_allow.condition_list_id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

#
# Create a site and attach a policy to it
#
# NOTE: a site also gets default policies and conditions generated for it by
# the backend. Destroying the site leaves them behind, so applying and
# destroying this sample repeatedly eventually fails with "Maximum policy count
# will be exceeded"; clean the leftovers up out of band when that happens.
#

data "incapsula_abp_site_analysis_settings" "login" {
  rate_limiting           = "per_site"
  max_requests_per_minute = 100

  # Load-bearing: `use_site_rate_limiting_parameters` defaults to `true`, which
  # makes the runtime use the site's `default_max_*` values and ignore the three
  # per-selector overrides entirely — it is not a per-value fallback. An
  # override next to the default would therefore be dead configuration.
  use_site_rate_limiting_parameters = false
}

data "incapsula_abp_site_analysis_settings" "static" {
  rate_limiting = "none"
}

data "incapsula_abp_site_analysis_settings" "postback" {
  rate_limiting              = "custom_scope"
  rate_limiting_custom_scope = "my scope"
  max_requests_per_minute    = 55
  max_requests_per_session   = 555
  # durations use the largest fitting units: "1h", not "60m"
  max_session_length                = "1h"
  use_site_rate_limiting_parameters = false
}

resource "incapsula_abp_site" "sample_site" {
  account_id = var.account_id
  name       = "Sample site"
  industry   = "retail"

  default_max_requests_per_minute  = 60
  default_max_requests_per_session = 600
  default_max_session_length       = "2h"

  # Each selector matches requests one way, written as exactly one of the
  # path_prefix, path_regex, postback and api_security_group blocks.
  selector {
    purpose = "login"
    path_prefix {
      value            = "/login"
      case_insensitive = true
    }
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  selector {
    path_regex {
      value = "\\.png$"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.static.json
  }

  selector {
    postback {
      value = "web_interrogation"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }

  # The fourth criteria variant: requests belonging to an API Security group.
  selector {
    purpose = "backend_api"
    api_security_group {
      value = "sample-checkout-api"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}

data "incapsula_abp_sites" "all" {
  account_id = var.account_id
}

# Look up a single site. Exactly one of `site_id` or `name` must be set; setting
# both, or neither, is an error, and so is a name that matches no site or more
# than one. Use this for sites this configuration does not manage.
data "incapsula_abp_site" "site2_by_id" {
  account_id = var.account_id
  site_id    = incapsula_abp_site.site2.id
}

# The same site by name. The match is exact and case-sensitive.
data "incapsula_abp_site" "site2_by_name" {
  account_id = var.account_id
  name       = incapsula_abp_site.site2.name
}

locals {
  managed_sites = [incapsula_abp_site.site2.id, incapsula_abp_site.sample_site.id]
}

# The site priority list must name *every* site of the account, so a literal
# list breaks as soon as a site exists that it does not name — one created in
# the UI, by another workspace, or by another module — with the opaque
# `Wrong number of site IDs supplied`. Appending the sites this configuration
# does not own, read from `incapsula_abp_sites`, is the form that survives that.
#
# Two caveats of this form:
#   - `data.incapsula_abp_sites.all` is read before the apply, so the tail only
#     covers pre-existing sites. Every site this module creates must therefore
#     be listed in `local.managed_sites`; adding an `incapsula_abp_site` without
#     adding it there fails the apply with the same error.
#   - in the apply that destroys a managed site the stale data source still
#     returns it, so `site_ids` keeps the dead id for one round and converges on
#     the next apply.
resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
  # The list has to contain every site of the account, so the sites whose order
  # matters come first and everything else - sites created in the UI, by another
  # workspace or by another module - is appended from the data source. A literal
  # list breaks as soon as anyone adds a site.
  site_ids = concat(
    local.managed_sites,
    [for site in data.incapsula_abp_sites.all.sites : site.id if !contains(local.managed_sites, site.id)]
  )
}

#
# Add a condition to the default policy
#

# Lookup a default policy via default selector of the site, demonstrating the
# lookup by `id`.
# Note that default selector is indexed with `[0]` even if it is not techically a list,
# since every site can have only one default selector. This is limitation of
# Terraform Plugin SDK v2
#
# Note also that `name` and `account_id` are treated as search arguments and are
# left empty on a lookup by `id`, so they cannot be read back from here.
data "incapsula_abp_policy" "default_policy" {
  id = incapsula_abp_site.sample_site.default_selector[0].policy_id
}

# Address the `allow` directive of the default policy by action. The directives of
# a default policy are created by the backend in an order this configuration does
# not control, so `data.incapsula_abp_policy.default_policy.directive[0]` is not
# guaranteed to be the `allow` one.
data "incapsula_abp_directive" "sample_site_default_allow" {
  policy_id = data.incapsula_abp_policy.default_policy.id
  action    = "allow"
}

# Add `specific_visitor` to the allow directive of the default policy
resource "incapsula_abp_condition_list_entry" "sample_site_default_allow_specific_visitor" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.sample_site_default_allow.condition_list_id
  condition_id             = incapsula_abp_condition.specific_visitor.id
  tags                     = ["specific_visitor"]
  state                    = "monitor"
}

#
# Create domains in the previously created site
#
resource "incapsula_abp_domain" "test_com" {
  account_id              = var.account_id
  site_id                 = incapsula_abp_site.sample_site.id
  cookiescope             = "test.com"
  log_region              = "apac"
  cookie_mode             = "none_secure"
  enable_mitigation       = false
  enable_mobile_sdk_token = false

  # These three must be given in the exact form the server stores, otherwise the
  # plan never converges, because the server normalises the value and the
  # provider compares the stored form against the configured one:
  #   - `obfuscate_path` needs a leading `/` and no whitespace or fragment;
  #   - `interstitial_inprogress_iframe_src` needs an absolute URL with an
  #     explicit path (`http://www.example.com/` rather than
  #     `http://www.example.com`);
  #   - no two `unmasked_headers` entries may differ only in case.
  obfuscate_path                     = "/spooky-path"
  interstitial_inprogress_iframe_src = "http://www.example.com/iframe-src"
  divert_host                        = "www.example.com"
  unmasked_headers                   = ["content-length", "content-type"]
  proxy_flags                        = ["enable_referrer_fix", "inject_js_into_body"]

  # Keeps common static asset paths out of ABP analysis.
  #
  # NOTE: this value is not applied when the domain is created, only by a later
  # update, and a domain with no stored value behaves as if filtering were
  # enabled. Note also that `false` is indistinguishable from an omitted value,
  # so it is never sent and filtering cannot be turned off from here — which is
  # why this sample sets `true`, the only value that actually reaches the
  # server.
  filter_out_static_assets = true

  no_js_injection_path {
    path_prefix = "/no-js-here"
  }

  captcha_settings {
    // Todo: Could unpack this into a `data`
    geetest {
      geetest_captcha_id  = "abcd"
      geetest_private_key = "my key"
    }
  }

  analysis_ip_lookup_mode {
    header_name   = "X-Forwarded-For"
    reverse_index = 0
  }
  challenge_ip_lookup_mode {
    header_name   = "Origin"
    reverse_index = 0
  }
  criteria {
    exact = "test.com"
  }
}

resource "incapsula_abp_domain" "example_com" {
  account_id  = var.account_id
  site_id     = incapsula_abp_site.sample_site.id
  cookiescope = "example.com"
  log_region  = "usa"
  cookie_mode = "lax"

  // Todo: reference a rule here
  // Commented out because it requires an existing Incapsula rule
  # no_js_injection_path {
  #   incap_rule = "URL == \"/admin\""
  # }

  captcha_settings {
    managed_geetest {
      difficulty = "hard"
    }
  }
  criteria {
    prefix = "example.com"
  }

  // The fourth `criteria` variant matches a CloudWAF site by its numeric id
  // instead of by domain name. Commented out because it needs a CloudWAF site,
  // which is not available locally:
  # criteria {
  #   cloudwaf_id = 12345678
  # }
}

resource "incapsula_abp_domain" "dummy_com" {
  account_id  = var.account_id
  site_id     = incapsula_abp_site.sample_site.id
  cookiescope = "dummy.com"
  log_region  = "eu"
  cookie_mode = "legacy"

  // Clone an existing key's material instead of generating a fresh key. This
  // is honoured on create only: the API never returns it, so changing it later
  // is ignored and does not replace the domain (which would destroy its keys).
  encryption_key_id = incapsula_abp_domain_encryption_key.test_com.id

  captcha_settings {
    managed_hcaptcha {
      difficulty = "auto"
    }
  }
  criteria {
    suffix = "dummy.com"
  }
}

#
# Define domain priority order
#

# The list must contain *every* Domain attached to the site, in priority order.
# A Domain that is attached but missing from the list — including one created
# outside Terraform — fails every apply with
# `Wrong number of domain IDs supplied`, and the removal the plan offers is not
# applicable, so extend this list whenever a Domain is added to the site.
resource "incapsula_abp_site_domain_priority" "sample_site" {
  site_id    = incapsula_abp_site.sample_site.id
  domain_ids = [incapsula_abp_domain.example_com.id, incapsula_abp_domain.test_com.id, incapsula_abp_domain.dummy_com.id]
}

#
# Create encryption key for a domain
#

# NOTE: the key material is immutable, so changing `key` replaces the resource.
# By default Terraform deletes the old key before creating the new one, which
# leaves the domain without a key in between; `create_before_destroy` inverts
# that order. The key is also stored in state in plaintext, so protect the
# state accordingly.
resource "incapsula_abp_domain_encryption_key" "test_com" {
  domain_id = incapsula_abp_domain.test_com.id
  key       = "U2VjcmV0IGtleSB1c2luZyBzdGF0ZS1vZi10aGUtYXJ0IGJhc2U2NCBlbmNyeXB0aW9u"

  lifecycle {
    create_before_destroy = true
  }
}

#------------------------------------------------------------------------------
#
# Others
#

resource "incapsula_abp_site" "site2" {
  account_id = var.account_id
  name       = "terraform-site-2"
  industry   = "financial_services"

  default_max_requests_per_minute  = 30
  default_max_requests_per_session = 300
  default_max_session_length       = "1h"

  # Two selectors of different kinds, so that inserting, removing and
  # reordering them is exercised by this sample. Selectors are matched up with
  # their live counterparts by `id`, which is what keeps the attributes this
  # provider does not model attached to the right selector across an update.
  selector {
    purpose = "login"
    path_prefix {
      value = "/login"
    }
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  selector {
    purpose = "checkout"
    path_regex {
      value = "^/checkout/.*$"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}

resource "incapsula_abp_credential" "my_credential" {
  account_id = var.account_id
  // RSA key to use for encryption, pem encoded
  // decrypt using ex:`terraform output -raw encrypted_secret | base64 -d | openssl pkeyutl -decrypt -inkey <your-private-key-file> -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 -pkeyopt rsa_oaep_label:$(echo -n 'abp_credential' | xxd -p)
  // Changing this key re-encrypts the same secret in place; the credential on
  // the server, and therefore the secret itself, is left untouched.
  rsa_key = <<EOF
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2p9NrmAlOC6ZbXVuKoj4
b4PQhTsPng2VGItWnDzy2VsgsRNz0K1PuPUbRo4ZgbZ5Z5UXi9QEKnf376QjdMGl
NoQRBiBgQAX/fz87ax8bGoDRtHtLgS984M84hIjyhhpZhzwuctnyJ2NnCaiwQIEn
WoH5+hWsxF/YUpP/6DzdvGdJEpDKq0itQl5D4ZpfVbiB/KfU4GOGGXa0bFYZjbT3
xow7+zA4wA29Z+ShKU8fqaTMjwIt8iGF2G7KzYzF1SwTiAgW2qEzNvQP2loFOI6h
yuHpNlEqsQ7r0ov+f3UxSJixcum7H3KEY5BaUdc/i76pgKqVYwi107XKGSBpQFsa
74dCiIy8jQfGNr1usO46swaL7G7WYxKJmOAj32YDF1SaR571N8wEpvV0dU811emR
z2E+I3EIa4FLTLCFKb9SJTGc9kEWA+gndGCInmzzLZSmRCTd5a5GpKn2fBuhSyCP
CfUfcBSpL/iJ4xg3gtx5hgAgELtDYDh8Tv3vDw6AS6/c3hYs4MkAZpbQ8v531emy
gJeOSyLZdk/+ldkl3NcOX0xOqn9JjWKicvTTwpJyO4Gk97lff6GQpxFFDNzs81at
XyQKDg65HAse9wY2TGg8cc/vefRCXpZHoiGv+RlHaF+QpaxwAp2w47fHht39V0VX
ypuTiiPzdbQtr50+N65XJCUCAwEAAQ==
-----END PUBLIC KEY-----
EOF
}

output "my_credential_encrypted_secret" {
  value = incapsula_abp_credential.my_credential.encrypted_secret
}

#
# Server-derived metadata
#

# ABP reports its per-selector request counters against `derived_id`, which is
# recomputed from the selector's fields whenever any of them change, so a
# dashboard needs to be told the current value. Which criteria block a selector
# holds says how it matches, so the matched value is read from the one block
# that is present.
output "sample_site_selector_derived_ids" {
  value = {
    for selector in incapsula_abp_site.sample_site.selector :
    selector.derived_id => one(concat(
      [for criteria in selector.path_prefix : criteria.value],
      [for criteria in selector.path_regex : criteria.value],
      [for criteria in selector.postback : criteria.value],
      [for criteria in selector.api_security_group : criteria.value],
    ))
  }
}

output "policy2_created_at" {
  value = data.incapsula_abp_policy.policy2.created_at
}

# Managed Conditions advertise the directive actions they are meant for.
output "managed_monitoring_tools_recommended_directives" {
  value = data.incapsula_abp_condition.managed_monitoring_tools.recommended_directives
}
