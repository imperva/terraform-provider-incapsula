data "incapsula_abp_condition_list" "sample_condition_list_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_condition_list.sample_condition_list.name
}
