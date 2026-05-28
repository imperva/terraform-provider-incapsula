terraform {
  required_providers {
    incapsula = {
      source = "imperva/incapsula"
    }
  }
}

provider "incapsula" {
  api_key        = "foo"
  api_id         = "bar"
  base_url       = "http://localhost:8081"
  base_url_rev_2 = "http://localhost:8081"
  base_url_rev_3 = "http://localhost:8081"
  base_url_api   = "http://localhost:8081"
}

locals {
  account_id = "cd3ba503-f034-4912-8f89-a599c8cfbbc6"
}

module "abp" {
  source     = "./abp"
  account_id = local.account_id
}


data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}

resource "incapsula_abp_preflight" "current" {
  account_id   = local.account_id
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
