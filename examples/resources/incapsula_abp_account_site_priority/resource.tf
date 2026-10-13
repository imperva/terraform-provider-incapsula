# Defines the order in which an account's Sites are matched. Sites are matched
# top-down. Both Terraform-managed sites and externally-managed sites (looked up
# via the `incapsula_abp_site` data source) can be referenced.
resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
  site_ids = [
    incapsula_abp_site.sample_site.id,
    incapsula_abp_site.site2.id,
    data.incapsula_abp_site.external_site.id,
  ]
}
