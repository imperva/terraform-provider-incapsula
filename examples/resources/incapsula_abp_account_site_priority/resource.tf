# Defines the order in which an account's Sites are matched. Sites are matched
# top-down.
#
# The list has to contain *every* Site of the account, so list the Sites you
# care about explicitly and append all the others — Sites created in the UI, by
# another workspace, or by another module — from the `incapsula_abp_sites` data
# source. A literal list of IDs breaks as soon as anyone adds a Site.
data "incapsula_abp_sites" "all" {
  account_id = var.account_id
}

locals {
  # The Sites whose relative priority matters, highest priority first.
  prioritized_sites = [
    incapsula_abp_site.sample_site.id,
    incapsula_abp_site.site2.id,
  ]
}

resource "incapsula_abp_account_site_priority" "accprio" {
  account_id = var.account_id
  site_ids = concat(
    local.prioritized_sites,
    [for site in data.incapsula_abp_sites.all.sites : site.id if !contains(local.prioritized_sites, site.id)]
  )
}
