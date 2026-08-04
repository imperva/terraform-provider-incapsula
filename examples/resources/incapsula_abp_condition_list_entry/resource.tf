# Attach a condition to a condition list. The parent list can be a standalone
# `incapsula_abp_condition_list`, or a directive's condition/skip-condition list.
resource "incapsula_abp_condition_list_entry" "sample_condition_list_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_condition_list.sample_condition_list.id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}
