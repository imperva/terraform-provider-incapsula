# Provider information, like api_id and api_key can be specified here
# or as environment variables: INCAPSULA_API_ID and INCAPSULA_API_KEY
provider "incapsula" {
  api_id = "foo"
  api_key = "bar"
}

# Site information
resource "incapsula_site" "example-site" {
  domain = "examplesite.com"
}

####################################################################
# Security Rules
#################################################################### 

# Security Rule: Country
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
  depends_on = ["incapsula_site.example-site"]
}

# Security Rule: Blacklist IP
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.1,192.168.1.2"
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-country-rule"]
}

# Security Rule: URL
resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
  rule_id = "api.acl.blacklisted_urls"
  site_id = "${incapsula_site.example-site.id}"
  url_patterns = "CONTAINS,EQUALS"
  urls = "/alpha,/bravo"
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-ip-rule"]
}

# Security Rule: Whitelist IP
resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
  rule_id = "api.acl.whitelisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.3,192.168.1.4"
  depends_on = ["incapsula_site.example-site","incapsula_acl_security_rule.example-global-blacklist-url-rule"]
}
