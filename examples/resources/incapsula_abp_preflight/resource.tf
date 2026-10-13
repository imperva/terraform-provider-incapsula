# A preflight captures a snapshot (hash) of the account's pending configuration.
# Publishing the preflight makes that configuration active. `pending_hash` comes
# from the `incapsula_abp_pending_changes` data source, which should depend on
# all the ABP resources you want included in the snapshot.
data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}

resource "incapsula_abp_preflight" "current" {
  account_id   = var.account_id
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}
