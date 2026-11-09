# Produces a marker that changes whenever Terraform has a planned change inside
# `depends_on`, which is what makes `incapsula_abp_preflight` take a fresh
# snapshot. List every ABP resource here (most simply: the module that holds
# them) — a change outside `depends_on` does not trigger a publish.
data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}
