# List every Site in the account, including those managed outside of this
# Terraform configuration.
data "incapsula_abp_sites" "all" {
  account_id = var.account_id
}

# `incapsula_abp_account_site_priority` requires every Site of the account
# exactly once, which is awkward when some of them are managed elsewhere. Give
# the Sites managed here the highest priority and append the rest in the order
# the API reports them.
resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
  site_ids = concat(
    [incapsula_abp_site.sample_site.id],
    [for site in data.incapsula_abp_sites.all.sites : site.id if site.id != incapsula_abp_site.sample_site.id],
  )
}
