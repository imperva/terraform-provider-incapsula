# A Domain (a.k.a. Website) lives inside a Site and is matched against incoming
# request hosts via its `criteria` block.
resource "incapsula_abp_domain" "test_com" {
  account_id  = var.account_id
  site_id     = incapsula_abp_site.sample_site.id
  cookiescope = "test.com"
  log_region  = "apac"
  cookie_mode = "none_secure"

  enable_mitigation       = false
  enable_mobile_sdk_token = false

  obfuscate_path                     = "/spooky-path"
  interstitial_inprogress_iframe_src = "http://www.example.com/iframe-src"
  divert_host                        = "www.example.com"
  unmasked_headers                   = ["content-length", "content-type"]
  proxy_flags                        = ["enable_referrer_fix", "inject_js_into_body"]

  no_js_injection_path {
    path_prefix = "/no-js-here"
  }

  captcha_settings {
    geetest {
      geetest_captcha_id  = "abcd"
      geetest_private_key = "my key"
    }
  }

  analysis_ip_lookup_mode {
    header_name   = "X-Forwarded-For"
    reverse_index = 0
  }

  challenge_ip_lookup_mode {
    header_name   = "Origin"
    reverse_index = 0
  }

  criteria {
    exact = "test.com"
  }
}
