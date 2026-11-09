# Look up a Policy by name.
data "incapsula_abp_policy" "policy_lookup" {
  account_id = var.account_id
  name       = incapsula_abp_policy.policy_with_custom_directives.name
}

# Look up the account global Policy. Note: the global Policy can't be referenced
# by its returned id elsewhere (e.g. in `incapsula_abp_directive.policy_id`);
# use `account_global_policy = true` on `incapsula_abp_directive` instead.
data "incapsula_abp_policy" "account_global" {
  account_id     = var.account_id
  account_global = true
}

# Look up a Site's default Policy via its default Selector. The default Selector
# is indexed with `[0]` even though a Site has exactly one.
data "incapsula_abp_policy" "default_policy" {
  id = incapsula_abp_site.sample_site.default_selector[0].policy_id
}
