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

resource "incapsula_abp_condition_list_entry" "sample_condition_list_specific_visitor" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_condition_list.sample_condition_list.id
  condition_id             = incapsula_abp_condition.specific_visitor.id
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

# TODO

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
    action = "proof_of_work"
    # TODO: proof_of_work should attach the configuration
    # TODO: skip conditions for proof_of_work
  }
}

resource "incapsula_abp_condition_list_entry" "policy2_allow_monitoring_tools" {
  account_id = var.account_id
  # TODO: index by action?
  parent_condition_list_id = incapsula_abp_policy.policy2.directive[0].condition_id
  condition_id             = data.incapsula_abp_condition.managed_monitoring_tools.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

resource "incapsula_abp_condition_list_entry" "policy2_block_sample_condition_list" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.policy2.directive[1].condition_id
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

  # Temporary commented out as it doesn't work locally due to MY dependency
  # no_js_injection_path {
  #   path_prefix = "/no-js-here"
  # }

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

  # Temporary commented out as it doesn't work locally due to MY dependency
  // Todo: reference a rule here
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

/*resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
  site_ids = [
    incapsula_abp_site.site2.id,
    incapsula_abp_site.site1.id,
    // .. fill out complete list
  ]
}*/
