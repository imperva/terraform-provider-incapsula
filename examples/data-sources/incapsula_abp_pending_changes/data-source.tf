# Returns the hash of the account's pending (unpublished) configuration. Make it
# depend on the ABP resources you want included in the snapshot so the hash is
# recomputed whenever any of them change.
data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}
