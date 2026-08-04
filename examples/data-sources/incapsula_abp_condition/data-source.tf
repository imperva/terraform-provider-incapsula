# Look up an account-owned condition by name.
data "incapsula_abp_condition" "specific_visitor_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_condition.specific_visitor.name
}

# Look up a managed (Imperva-provided) condition.
data "incapsula_abp_condition" "managed_monitoring_tools" {
  account_id = var.account_id
  name       = "Monitoring Tools"
  managed    = true
}
