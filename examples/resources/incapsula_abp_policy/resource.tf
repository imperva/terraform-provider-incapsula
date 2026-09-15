# Policy with standard (auto-generated) directives.
resource "incapsula_abp_policy" "policy_with_standard_directives" {
  account_id              = var.account_id
  name                    = "Policy with standard directives"
  description             = "Terraform-managed policy with standard directives"
  use_standard_directives = true
}

# Policy with custom directives. Conditions are attached to a directive's
# condition list via `incapsula_abp_condition_list_entry`.
resource "incapsula_abp_policy" "policy_with_custom_directives" {
  account_id  = var.account_id
  name        = "Policy with custom directives"
  description = "Demonstrate how to create a policy with custom directives"

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
