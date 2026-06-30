# Look up a Site that is not managed by this Terraform configuration, e.g. when
# building an `incapsula_abp_account_site_priority` list. Provide exactly one of
# `site_id` or `name`.
data "incapsula_abp_site" "external_site" {
  account_id = var.account_id
  name       = "Externally managed site"
}
