# Returns the ABP account ID identified by the provider's API credentials.
data "incapsula_abp_account" "current" {}

resource "incapsula_abp_site" "example" {
  account_id = data.incapsula_abp_account.current.account_id
  name       = "example"
}
