# Defines the order in which a Site's Domains are matched. Domains are matched
# top-down, so list more specific Domains before wildcard/catch-all ones.
#
# The list does not have to name every Domain of the Site. The Domains it names
# are the highest priority ones, in this order; any other Domain attached to the
# Site — including one created in the ABP UI or by another workspace — is
# matched after them. Read `effective_domain_ids` to see the resulting order.
resource "incapsula_abp_site_domain_priority" "sample_site" {
  site_id = incapsula_abp_site.sample_site.id
  domain_ids = [
    incapsula_abp_domain.example_com.id,
    incapsula_abp_domain.test_com.id,
  ]
}
