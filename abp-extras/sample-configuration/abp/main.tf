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

data "incapsula_abp_site_analysis_settings" "login" {
  rate_limiting           = "per_site"
  max_requests_per_minute = 100
}

data "incapsula_abp_site_analysis_settings" "static" {
  rate_limiting = "none"
}

data "incapsula_abp_site_analysis_settings" "postback" {
  rate_limiting                     = "custom_scope"
  rate_limiting_custom_scope        = "my scope"
  max_requests_per_minute           = 55
  max_requests_per_session          = 555
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

  selector {
    purpose = "login"
    kind {
      path_prefix      = "/login"
      case_insensitive = true
    }
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  selector {
    kind {
      path_regex = "\\.png$"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.static.json
  }

  selector {
    kind {
      postback = "web_interrogation"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }

  # The fourth criteria variant: requests belonging to an API Security group.
  selector {
    purpose = "backend_api"
    kind {
      api_security_group = "sample-checkout-api"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}

data "incapsula_abp_sites" "all" {
  account_id = var.account_id
}

locals {
  managed_sites = [incapsula_abp_site.site2.id, incapsula_abp_site.sample_site.id]
}

resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
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
  // Todo: backend auto-prefixes with `/` causing a perpetual change-detection if omitted on this field
  // Other paths are validated enforcing path prefixing, we could do that in the tf-layer, backend, or not at all
  obfuscate_path                     = "/spooky-path"
  interstitial_inprogress_iframe_src = "http://www.example.com/iframe-src"
  divert_host                        = "www.example.com"
  unmasked_headers                   = ["content-length", "content-type"]
  proxy_flags                        = ["enable_referrer_fix", "inject_js_into_body"]

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
  // Commented out due to no MY available locally
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

resource "incapsula_abp_site_domain_priority" "sample_site" {
  site_id    = incapsula_abp_site.sample_site.id
  domain_ids = [incapsula_abp_domain.example_com.id, incapsula_abp_domain.test_com.id, incapsula_abp_domain.dummy_com.id]
}

#
# Create encryption key for a domain
#

resource "incapsula_abp_domain_encryption_key" "test_com" {
  domain_id = incapsula_abp_domain.test_com.id
  key       = "U2VjcmV0IGtleSB1c2luZyBzdGF0ZS1vZi10aGUtYXJ0IGJhc2U2NCBlbmNyeXB0aW9u"
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
    kind {
      path_prefix = "/login"
    }
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  selector {
    purpose = "checkout"
    kind {
      path_regex = "^/checkout/.*$"
    }
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}

# TODO
#resource "incapsula_abp_account_site_priority" "accprio" {
#  account_id = var.account_id
#  site_ids = [
#    incapsula_abp_site.sample_site.id,
#    incapsula_abp_site.site2.id,
#    data.incapsula_abp_site.ext_site1.id,
#    data.incapsula_abp_site.ext_site2.id,
#    data.incapsula_abp_site.ext_site3.id,
#    data.incapsula_abp_site.ext_site4.id,
#    data.incapsula_abp_site.ext_site5.id,
#    data.incapsula_abp_site.ext_site6.id,
#    data.incapsula_abp_site.ext_site7.id,
#    data.incapsula_abp_site.ext_site8.id,
#    data.incapsula_abp_site.ext_site9.id,
#    data.incapsula_abp_site.ext_site10.id,
#    data.incapsula_abp_site.ext_site11.id,
#    data.incapsula_abp_site.ext_site12.id,
#    data.incapsula_abp_site.ext_site13.id,
#  ]
#}

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
# dashboard needs to be told the current value.
output "sample_site_selector_derived_ids" {
  value = {
    for selector in incapsula_abp_site.sample_site.selector :
    selector.derived_id => one([
      for criteria in [
        selector.kind[0].path_prefix,
        selector.kind[0].path_regex,
        selector.kind[0].postback,
        selector.kind[0].api_security_group,
      ] : criteria if criteria != ""
    ])
  }
}

output "policy2_created_at" {
  value = data.incapsula_abp_policy.policy2.created_at
}

# Managed Conditions advertise the directive actions they are meant for.
output "managed_monitoring_tools_recommended_directives" {
  value = data.incapsula_abp_condition.managed_monitoring_tools.recommended_directives
}
