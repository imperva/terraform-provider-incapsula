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

# To avoid conflicts:
# - create file `vars.tfvars` containing
# ```
# account_id = "<your account id>"
# ```
# - specify it as `terraform apply -var-file=vars.tfvars`
# (environment can also be used, choose your preferred approach)
variable "account_id" {
  type = string
}

module "abp" {
  source     = "./abp"
  account_id = var.account_id
}


data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}

resource "incapsula_abp_preflight" "current" {
  account_id   = var.account_id
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
