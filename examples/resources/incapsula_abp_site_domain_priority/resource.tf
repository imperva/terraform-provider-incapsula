# Defines the order in which a Site's Domains are matched. Domains are matched
# top-down, so list more specific Domains before wildcard/catch-all ones.
resource "incapsula_abp_site_domain_priority" "sample_site" {
  site_id = incapsula_abp_site.sample_site.id
  domain_ids = [
    incapsula_abp_domain.example_com.id,
    incapsula_abp_domain.test_com.id,
    incapsula_abp_domain.dummy_com.id,
  ]
}
