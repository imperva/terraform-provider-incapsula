# Policy with explicit Directives.
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

# Conditions are attached to a Directive's Condition List. Address the Directive
# by its action with the `incapsula_abp_directive` data source -- never by its
# index in the `directive` list, which moves when Directives are added, removed
# or reordered.
data "incapsula_abp_directive" "block" {
  policy_id = incapsula_abp_policy.policy_with_custom_directives.id
  action    = "block"
}

resource "incapsula_abp_condition_list_entry" "block_bad_bots" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.block.condition_list_id
  condition_id             = incapsula_abp_condition.bad_bots.id
  state                    = "active"
}

# A `proof_of_work` Directive also has a skip Condition List: matches there
# bypass the Directive.
data "incapsula_abp_directive" "pow" {
  policy_id = incapsula_abp_policy.policy_with_custom_directives.id
  action    = "proof_of_work"
}

resource "incapsula_abp_condition_list_entry" "pow_skip_monitoring" {
  account_id               = var.account_id
  parent_condition_list_id = data.incapsula_abp_directive.pow.skip_condition_list_id
  condition_id             = data.incapsula_abp_condition.monitoring_tools.id
  state                    = "active"
}
