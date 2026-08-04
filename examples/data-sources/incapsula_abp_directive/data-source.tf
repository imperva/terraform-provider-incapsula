# Reference a directive by its action. Preferred over indexing into
# `incapsula_abp_policy.<name>.directive[i]` as it is more predictable.
data "incapsula_abp_directive" "std_policy_block" {
  policy_id = incapsula_abp_policy.policy_with_standard_directives.id
  action    = "block"
}

# Reference a directive of the account global policy.
data "incapsula_abp_directive" "account_global_allow" {
  account_id            = var.account_id
  account_global_policy = true
  action                = "allow"
}
