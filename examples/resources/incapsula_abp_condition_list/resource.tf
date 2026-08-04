resource "incapsula_abp_condition_list" "sample_condition_list" {
  account_id  = var.account_id
  name        = "Sample condition list"
  description = "Reusable condition list"
}
