# Provider information, like api_id and api_key can be specified here
# or as environment variables: INCAPSULA_API_ID and INCAPSULA_API_KEY
provider "incapsula" {
  api_id = "31228"
  api_key = "35ca97fd-979b-4ba2-b8b6-3be2e7fd3889"
}

# Site information
resource "incapsula_site" "example-site" {
  domain = "foobar.com"
}

//####################################################################
//# Data Center
//####################################################################

# Data Centers
resource "incapsula_data_center" "example-data-center" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center"
  server_address = "8.8.4.4"
  is_standby = "yes"
  is_content = "yes"
  depends_on = ["incapsula_site.example-site"]
}

// todo: review data center servers delete
//# Data Center Servers
//resource "incapsula_data_center_servers" "example-data-center-servers" {
//  dc_id = "${incapsula_data_center.example-data-center.id}"
//  site_id = "${incapsula_site.example-site.id}"
//  server_address = "4.4.4.4"
//  is_standby = "yes"
//  depends_on = ["incapsula_site.example-site", "incapsula_data_center.example-data-center"]
//}

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
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_ips"
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
  enabled = "true"
  priority = "1"
  name = "Example incap rule alert"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule require javascript support 3"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_INTRUSIVE_HTML"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule block ip"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK_IP"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule block request"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["incapsula_site.example-site"]
}

// todo: action block session is unknown
//# Incap Rule: Block Session
//resource "incapsula_incap_rule" "example-incap-rule-block-session" {
//  site_id = "${incapsula_site.example-site.id}"
//  enabled = "true"
//  priority = "1"
//  name = "Example incap rule block session"
//  action = "RULE_ACTION_BLOCK_SESSION"
//  filter = "Full-URL == \"/someurl\""
//  depends_on = ["incapsula_site.example-site"]
//}

// todo: error
//# Incap Rule: Delete Cookie (ADR)
//resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
//  enabled = "true"
//  priority = "1"
//  name = "Example incap rule delete cookie"
//  site_id = "${incapsula_site.example-site.id}"
//  action = "RULE_ACTION_DELETE_COOKIE"
//  filter = "Full-URL == \"/someurl\""
//  depends_on = ["incapsula_site.example-site"]
//}

// todo: error
//# Incap Rule: Delete Header (ADR)
//resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
//  enabled = "true"
//  priority = "1"
//  name = "Example incap rule delete header"
//  site_id = "${incapsula_site.example-site.id}"
//  action = "RULE_ACTION_DELETE_HEADER"
//  filter = "Full-URL == \"/someurl\""
//  depends_on = ["incapsula_site.example-site"]
//}

# Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule forward to data center"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_FORWARD_TO_DC"
  filter = "Full-URL == \"/someurl\""
  dc_id = "${incapsula_data_center.example-data-center.id}"
  allow_caching = "false"
  depends_on = ["incapsula_site.example-site", "incapsula_data_center.example-data-center"]
}

# Incap Rule: Redirect (ADR)
resource "incapsula_incap_rule" "example-incap-rule-redirect" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule redirect"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_REDIRECT"
  filter = "Full-URL == \"/someurl\""
  response_code = "302"
  from = "https://site1.com/url1"
  to = "https://site2.com/url2"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  enabled = "true"
  priority = "1"
  name = "Example incap rule require cookie support"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_RETRY"
  filter = "Full-URL == \"/someurl\""
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Rewrite Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-cookie" {
  enabled = "true"
  priority = "18"
  name = "Example incap rule rewrite cookie"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_REWRITE_COOKIE"
  filter = "Full-URL == \"/someurl\""
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  allow_caching = "false"
  rewrite_name = "my_cookie_name"
  depends_on = ["incapsula_site.example-site"]
}

# Incap Rule: Rewrite Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-header" {
  enabled = "true"
  priority = "17"
  name = "Example incap rule rewrite header"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_REWRITE_HEADER"
  filter = "Full-URL == \"/someurl\""
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  allow_caching = "false"
  rewrite_name = "my_test_header"
  depends_on = ["incapsula_site.example-site"]
}

// todo: error
//# Incap Rule: Rewrite URL (ADR)
//resource "incapsula_incap_rule" "example-incap-rule-rewrite-url" {
//  enabled = "true"
//  priority = "1"
//  name = "Example incap rule rewrite url"
//  action = "RULE_ACTION_REWRITE_URL"
//  filter = "Full-URL == \"/someurl\""
//  depends_on = ["incapsula_site.example-site"]
//}
