# Provider information, like api_id and api_key can be specified here
# or as environment variables: INCAPSULA_API_ID and INCAPSULA_API_KEY
provider "incapsula" {
  api_id = "foo"
  api_key = "bar"
}

# Site information
resource "incapsula_site" "example-site" {
  domain = "www.examplesite.com"
}

####################################################################
# Data Center
####################################################################

# Data Centers
resource "incapsula_data_center" "example-data-center" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center"
  server_address = "8.8.4.4"
  is_content = "yes"
}

# Data Center Servers
resource "incapsula_data_center_servers" "example-data-center-servers" {
  dc_id = "${incapsula_data_center.example-data-center.id}"
  site_id = "${incapsula_site.example-site.id}"
  server_address = "4.4.4.4"
  is_standby = "no"
}

####################################################################
# Security Rules
####################################################################

# Security Rule: Country
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = ["AI","AN"]
}

# Security Rule: Blacklist IP
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_ips"
  ips = ["192.168.1.1","192.168.1.2"]
}

# Security Rule: Blacklist IP Exception
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule_exception" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = ["192.168.1.1","192.168.1.2"]
  urls = ["/myurl","/myurl2"]
  url_patterns = ["EQUALS","CONTAINS"]
  countries = ["JM","US"]
  client_apps= ["488","123"]
}

# Security Rule: URL
resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
  rule_id = "api.acl.blacklisted_urls"
  site_id = "${incapsula_site.example-site.id}"
  url_patterns = ["CONTAINS","EQUALS"]
  urls = ["/alpha","/bravo"]
}

# Security Rule: Whitelist IP
resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
  rule_id = "api.acl.whitelisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = ["192.168.1.3","192.168.1.4"]
}

####################################################################
# Incap Rules
####################################################################

# Incap Rule: Alert
resource "incapsula_incap_rule" "example-incap-rule-alert" {
  priority = "1"
  name = "Example incap rule alert"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  priority = "1"
  name = "Example incap rule require javascript support 3"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_INTRUSIVE_HTML"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  priority = "1"
  name = "Example incap rule block ip"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK_IP"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  priority = "1"
  name = "Example incap rule block request"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Block Session
resource "incapsula_incap_rule" "example-incap-rule-block-session" {
  priority = "1"
  name = "Example incap rule block session"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK_USER"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Delete Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
  priority = "1"
  name = "Example incap rule delete cookie"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_DELETE_COOKIE"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

# Incap Rule: Delete Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
  priority = "1"
  name = "Example incap rule delete header"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_DELETE_HEADER"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

# Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  priority = "1"
  name = "Example incap rule forward to data center"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_FORWARD_TO_DC"
  filter = "Full-URL == \"/someurl\""
  dc_id = "${incapsula_data_center.example-data-center.id}"
  allow_caching = "false"
 }

# Incap Rule: Redirect (ADR)
resource "incapsula_incap_rule" "example-incap-rule-redirect" {
  priority = "1"
  name = "Example incap rule redirect"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_REDIRECT"
  filter = "Full-URL == \"/someurl\""
  response_code = "302"
  from = "https://site1.com/url1"
  to = "https://site2.com/url2"
}

# Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  priority = "1"
  name = "Example incap rule require cookie support"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_RETRY"
  filter = "Full-URL == \"/someurl\""
}

# Incap Rule: Rewrite Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-cookie" {
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
}

# Incap Rule: Rewrite Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-header" {
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
}

# Incap Rule: Rewrite URL (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-url" {
  priority = "1"
  name = "ExampleRewriteURL"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_REWRITE_URL"
  filter = "Full-URL == \"/someurl\""
  add_missing = "true"
  from = "*"
  to = "/redirect"
  allow_caching = "false"
  rewrite_name = "my_test_header"
}
