# Attach a Condition to a Condition List. The parent list can be a standalone
# `incapsula_abp_condition_list`, or a Directive's Condition List or skip
# Condition List.
resource "incapsula_abp_condition_list_entry" "sample_condition_list_okhttp" {
  account_id               = var.account_id
  parent_condition_list_id = incapsula_abp_condition_list.sample_condition_list.id
  condition_id             = incapsula_abp_condition.okhttp.id
  state                    = "active"
  tags                     = ["terraform_managed"]
}
