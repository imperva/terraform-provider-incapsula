# A Site (a.k.a. Website Group) groups Domains together and maps incoming
# requests to Policies via an ordered list of Selectors.
resource "incapsula_abp_site" "sample_site" {
  account_id = var.account_id
  name       = "Sample site"

  default_max_requests_per_minute  = 60
  default_max_requests_per_session = 600
  default_max_session_length       = "2h"

  # Match by path prefix and apply a specific policy.
  selector {
    path_prefix       = "/login"
    policy_id         = incapsula_abp_policy.policy_with_custom_directives.id
    analysis_settings = data.incapsula_abp_site_analysis_settings.login.json
  }

  # Match by regex.
  selector {
    path_regex        = "\\.png$"
    analysis_settings = data.incapsula_abp_site_analysis_settings.static.json
  }

  # A postback selector.
  selector {
    postback          = "web_interrogation"
    analysis_settings = data.incapsula_abp_site_analysis_settings.postback.json
  }
}
