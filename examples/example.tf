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
# Data Center
####################################################################

# Data Centers
resource "incapsula_data_center" "example-data_center" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center"
  server_address = "192.168.1.10"
  is_standby = "yes"
  is_content = "yes"
  depends_on = ["incapsula_site.example-site"]
}

# Data Center Servers
resource "incapsula_data_center_servers" "example-data_center_servers" {
  dc_id = "${incapsula_data_center.example-data_center.id}"
  server_address = "192.168.2.10"
  is_standby = "yes"
  depends_on = ["incapsula_site.example-site", "incapsula_data_center.example-data_center"]
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

# Security Rule: Blacklist IP Exception
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule_exception" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.1,192.168.1.2"
  urls = "/myurl,/myurl2"
  url_patterns = "EQUALS,CONTAINS"
  countries = "JM,US"
  client_apps= "488,123"
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
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-url-rule"]
}

####################################################################
# Incap Rules
####################################################################

# Incap Rule: Alert
resource "incapsula_incap_rule" "example-incap-rule-alert" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule alert"
  action = "RULE_ACTION_ALERT"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule require javascript support"
  action = "RULE_ACTION_INTRUSIVE_HTML"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule block ip"
  action = "RULE_ACTION_BLOCK_IP"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule block request"
  action = "RULE_ACTION_BLOCK"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site"]
}

# todo: Incap Rule: Block Session
//resource "incapsula_incap_rule" "example-incap-rule-block-session" {
//  site_id = "${incapsula_site.example-site.id}"
//  enabled = "true"
//  priority = "1"
//  name = "Example incap rule block session"
//  action = "RULE_ACTION_BLOCK"
//  filter = "/someurl"
//  depends_on = ["incapsula_site.example-site"]
//}

# Incap Rule: Delete Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule delete cookie"
  action = "RULE_ACTION_DELETE_COOKIE"
  filter = "foo"
  rule_id = "${incapsula_acl_security_rule.example-global-blacklist-country-rule.id}"
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-country-rule"]
}

# Incap Rule: Delete Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule delete header"
  action = "RULE_ACTION_DELETE_HEADER"
  filter = "foo"
  rule_id = "${incapsula_acl_security_rule.example-global-blacklist-country-rule.id}"
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-country-rule"]
}

# Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule forward to data center"
  dc_id = "${incapsula_data_center.example-data_center.id}"
  action = "RULE_ACTION_FORWARD_TO_DC"
  allow_caching = "false"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site", "incapsula_data_center.example-data_center"]
}

# Incap Rule: Redirect (ADR)
resource "incapsula_incap_rule" "example-incap-rule-redirect" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule redirect"
  response_code = "302"
  action = "RULE_ACTION_REDIRECT"
  from = "https://site1.com/url1"
  to = "https://site2.com/url2"
  filter = "testval"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "1"
  name = "Example incap rule require cookie support"
  action = "RULE_ACTION_RETRY"
  filter = "/someurl"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Rewrite Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-cookie" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "18"
  name = "Example incap rule rewrite cookie"
  action = "RULE_ACTION_REWRITE_COOKIE"
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  allow_caching = "false"
  filter = "/someurl"
  rewrite_name = "my_cookie_name"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Rewrite Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-header" {
  site_id = "${incapsula_site.example-site.id}"
  enabled = "true"
  priority = "17"
  name = "Example incap rule rewrite header"
  action = "RULE_ACTION_REWRITE_HEADER"
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  allow_caching = "false"
  filter = "/testurl"
  rewrite_name = "my_test_header"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Rewrite URL (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-url" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule rewrite url"
  action = "RULE_ACTION_REWRITE_URL"
  filter = "/testurl"
  rule_id = "${incapsula_acl_security_rule.example-global-blacklist-country-rule.id}"
  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-country-rule"]
}
