# List every policy of the account.
data "incapsula_abp_policies" "all" {
  account_id = var.account_id
}

output "policy_names" {
  value = data.incapsula_abp_policies.all.policies[*].name
}
