terraform {
  required_providers {
    incapsula = {
      source = "terraform-providers/incapsula"
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

resource "incapsula_abp_condition" "cond1" {
  account_id  = "a9fa7bb9-a36e-40aa-ac81-fe320d634988"
  name        = "terraform-0"
  description = "Created through terraform twice"
  code        = "(any true false)"
}
