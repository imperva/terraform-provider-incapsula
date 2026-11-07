package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAbpSmoke(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "incapsula_abp_condition" "c" {
  account_id  = %[1]q
  name        = "cond"
  description = "d"
  code        = "(all flags.no_token)"
}

resource "incapsula_abp_condition_list" "l" {
  account_id  = %[1]q
  name        = "list"
  description = "d"
}

resource "incapsula_abp_condition_list_entry" "e" {
  account_id               = %[1]q
  parent_condition_list_id = incapsula_abp_condition_list.l.id
  condition_id             = incapsula_abp_condition.c.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}

resource "incapsula_abp_proof_of_work_configuration" "p" {
  account_id = %[1]q
  name       = "pow"
  difficulty = 42
  algorithm  = "bbs"
}

resource "incapsula_abp_policy" "pol" {
  account_id = %[1]q
  name       = "pol"

  directive {
    action = "allow"
  }

  directive {
    action                         = "proof_of_work"
    proof_of_work_configuration_id = incapsula_abp_proof_of_work_configuration.p.id
  }
}

data "incapsula_abp_directive" "d" {
  policy_id = incapsula_abp_policy.pol.id
  action    = "allow"
}

data "incapsula_abp_site_analysis_settings" "s" {
  rate_limiting = "per_site"
}

resource "incapsula_abp_site" "site" {
  account_id = %[1]q
  name       = "site"
  industry   = "retail"

  selector {
    purpose = "login"
    path_prefix {
      value = "/login"
    }
    policy_id         = incapsula_abp_policy.pol.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.s.json
  }
}

resource "incapsula_abp_domain" "dom" {
  account_id  = %[1]q
  site_id     = incapsula_abp_site.site.id
  cookiescope = "test.com"
  log_region  = "eu"
  cookie_mode = "lax"

  captcha_settings {
    managed_geetest {
      difficulty = "hard"
    }
  }

  criteria {
    exact = "test.com"
  }
}

resource "incapsula_abp_domain_encryption_key" "k" {
  domain_id = incapsula_abp_domain.dom.id
  key       = "U2VjcmV0IGtleSB1c2luZyBzdGF0ZS1vZi10aGUtYXJ0IGJhc2U2NCBlbmNyeXB0aW9u"
}

resource "incapsula_abp_site_domain_priority" "sdp" {
  site_id    = incapsula_abp_site.site.id
  domain_ids = [incapsula_abp_domain.dom.id]
}

resource "incapsula_abp_account_site_priority" "asp" {
  account_id = %[1]q
  site_ids   = [incapsula_abp_site.site.id]
}

data "incapsula_abp_account" "a" {}

data "incapsula_abp_conditions" "all" {
  account_id = %[1]q
}

data "incapsula_abp_policies" "all" {
  account_id = %[1]q
}

data "incapsula_abp_sites" "all" {
  account_id = %[1]q
}

data "incapsula_abp_site" "one" {
  account_id = %[1]q
  site_id    = incapsula_abp_site.site.id
}

data "incapsula_abp_condition" "one" {
  account_id = %[1]q
  name       = incapsula_abp_condition.c.name
}

data "incapsula_abp_condition_list" "one" {
  account_id = %[1]q
  name       = incapsula_abp_condition_list.l.name
}

data "incapsula_abp_policy" "one" {
  account_id = %[1]q
  name       = incapsula_abp_policy.pol.name
}

data "incapsula_abp_proof_of_work_configuration" "one" {
  account_id = %[1]q
  name       = incapsula_abp_proof_of_work_configuration.p.name
}

data "incapsula_abp_pending_changes" "pc" {}

resource "incapsula_abp_preflight" "pf" {
  account_id = %[1]q
  pending_hash = data.incapsula_abp_pending_changes.pc.hash
}

resource "incapsula_abp_publish" "pub" {
  preflight_id = incapsula_abp_preflight.pf.id
}

resource "incapsula_abp_credential" "cred" {
  account_id = %[1]q
  rsa_key    = <<-EOT
%[2]s
  EOT
}
`, abpFakeAccountId, abpTestRsaPublicPem(t, 0))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{Config: config, PlanOnly: true},
		},
	})
}
