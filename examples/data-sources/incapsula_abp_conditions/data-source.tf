# List the account-owned conditions.
data "incapsula_abp_conditions" "account" {
  account_id = var.account_id
}

# List the managed (Imperva-provided) conditions.
data "incapsula_abp_conditions" "managed" {
  account_id = var.account_id
  managed    = true
}

output "managed_condition_names" {
  value = data.incapsula_abp_conditions.managed.conditions[*].name
}
