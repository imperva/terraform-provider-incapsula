# Builds the opaque `json` analysis-settings value consumed by an
# `incapsula_abp_site` selector's `analysis_settings` field.

# Per-site rate limiting, using the Site's default limits.
data "incapsula_abp_site_analysis_settings" "login" {
  rate_limiting = "per_site"
}

# Per-site rate limiting with its own limits. `use_site_rate_limiting_parameters`
# must be false, otherwise the Site defaults win and `max_requests_per_minute` is
# ignored.
data "incapsula_abp_site_analysis_settings" "checkout" {
  rate_limiting                     = "per_site"
  max_requests_per_minute           = 100
  use_site_rate_limiting_parameters = false
}

# No rate limiting.
data "incapsula_abp_site_analysis_settings" "static" {
  rate_limiting = "none"
}

# Custom-scope rate limiting.
data "incapsula_abp_site_analysis_settings" "postback" {
  rate_limiting                     = "custom_scope"
  rate_limiting_custom_scope        = "my scope"
  max_requests_per_minute           = 55
  max_requests_per_session          = 555
  max_session_length                = "1h"
  use_site_rate_limiting_parameters = false
}
