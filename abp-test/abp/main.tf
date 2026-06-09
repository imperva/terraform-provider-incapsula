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
# Create a policy with standard directives and populate it with conditions
#

resource "incapsula_abp_policy" "policy_with_standard_directives" {
  account_id              = var.account_id
  name                    = "Policy with standard directives"
  description             = "Terraform-managed policy with standard directives"
  use_standard_directives = true
}

# Now add conditions to the automatically created directives
resource "incapsula_abp_condition_list_entry" "std_policy_allow_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.policy_with_standard_directives.directive[0].condition_list_id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

#
# Create a policy with custom directives and populate it with conditions
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

resource "incapsula_abp_condition_list_entry" "policy2_allow_monitoring_tools" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.policy2.directive[0].condition_list_id
  condition_id             = data.incapsula_abp_condition.managed_monitoring_tools.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

# Skip the proof_of_work directive for the managed monitoring tools condition
resource "incapsula_abp_condition_list_entry" "policy2_pow_skip_monitoring_tools" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.policy2.directive[2].skip_condition_list_id
  condition_id             = data.incapsula_abp_condition.managed_monitoring_tools.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

resource "incapsula_abp_condition_list_entry" "policy2_block_sample_condition_list" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.policy2.directive[1].condition_list_id
  condition_list_id        = incapsula_abp_condition_list.sample_condition_list.id
  state                    = "monitor"
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

  default_max_requests_per_minute  = 60
  default_max_requests_per_session = 600
  default_max_session_length       = "2h"

  selector {
    path_prefix       = "/login"
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  selector {
    path_regex        = "\\.png$"
    analysis_settings = data.incapsula_abp_site_analysis_settings.static.json
  }

  selector {
    postback          = "web_interrogation"
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}

#
# Add a condition to the default policy
#

# Lookup a default policy via default selector of the site. This is required
# to get access to the policy directives
# Note that default selector is indexed with `[0]` even if it is not techically a list,
# since every site can have only one default selector. This is limitation of
# Terraform Plugin SDK v2
data "incapsula_abp_policy" "default_policy" {
  id = incapsula_abp_site.sample_site.default_selector[0].policy_id
}


# Add `specific_visitor` to the allow directive of the default policy
resource "incapsula_abp_condition_list_entry" "sample_site_default_allow_specific_visitor" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_policy.default_policy.directive[0].condition_list_id
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

  default_max_requests_per_minute  = 30
  default_max_requests_per_session = 300
  default_max_session_length       = "1h"

  selector {
    path_prefix       = "/login"
    policy_id         = incapsula_abp_policy.policy2.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }
}

# TODO
# resource "incapsula_abp_account_site_priority" "site_priority" {
#   account_id = var.account_id
#   site_ids = [
#     incapsula_abp_site.sample_site.id,
#     incapsula_abp_site.site2.id,
#   ]
# }

resource "incapsula_abp_credential" "my_credential" {
  account_id = var.account_id
  // RSA key to use for encryption, pem encoded
  // decrypt using ex:`terraform output -raw encrypted_secret | base64 -d | openssl pkeyutl -decrypt -inkey <your-private-key-file> -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 -pkeyopt rsa_oaep_label:$(echo -n 'abp_credential' | xxd -p)
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
