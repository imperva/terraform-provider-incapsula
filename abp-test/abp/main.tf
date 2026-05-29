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

resource "incapsula_abp_policy" "poltest1" {
  account_id  = var.account_id
  name        = "policy pt 1"
  description = "My cool policy"

  directive {
    action = "allow"
  }
}

resource "incapsula_abp_policy" "poltest2" {
  account_id  = var.account_id
  name        = "policy pt 2"
  description = "My cool policy change desc"

  directive {
    action = "allow"
  }
}

resource "incapsula_abp_policy" "poltest3" {
  account_id  = var.account_id
  name        = "cool name"
  description = "policy 3.2"

  directive {
    action = "allow"
  }
}

resource "incapsula_abp_condition" "cond1" {
  account_id  = var.account_id
  name        = "terraform-0"
  description = "Created through terraform twice"
  code        = "(any true false)"
}

# Attach the literal condition above to the auto-generated condition list of
# poltest1's first directive.
resource "incapsula_abp_condition_list_entry" "poltest1_allow_cond1" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_policy.poltest1.directive[0].condition_id
  condition_id             = incapsula_abp_condition.cond1.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

#
# Proof Of Work
#

resource "incapsula_abp_proof_of_work_configuration" "pow1" {
  account_id = var.account_id
  name       = "terraform-pow-0"
  difficulty = 42
  algorithm  = "bbs"
}

data "incapsula_abp_proof_of_work_configuration" "pow1_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_proof_of_work_configuration.pow1.name
}

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

resource "incapsula_abp_site" "site1" {
  account_id = var.account_id
  name       = "terraform-site-0"

  default_max_requests_per_minute  = 60
  default_max_requests_per_session = 600
  default_max_session_length       = "2h"

  selector {
    path_prefix       = "/login"
    policy_id         = incapsula_abp_policy.poltest1.id
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
