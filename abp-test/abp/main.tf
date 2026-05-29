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

#
# Proof Of Work
#

resource "incapsula_abp_proof_of_work_configuration" "pow1" {
  account_id = var.account_id
  name       = "terraform-pow-0"
  difficulty = 42
  algorithm  = "bbs"
}

# TODO: explore how data sources impact publishing
data "incapsula_abp_proof_of_work_configuration" "pow1_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_proof_of_work_configuration.pow1.name
}
