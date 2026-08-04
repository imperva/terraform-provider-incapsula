# Publishes a specific preflight, making the corresponding configuration active.
resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
