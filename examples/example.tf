
provider "incapsula" {
  api_id = "your_api_id"
  api_key = "your_api_key"
}

resource "incapsula_site" "example-site" {
  ####################################################################
  # The first 6 parameters listed below are designated for initially creating the site in Incapsula.
  ####################################################################
  domain = ""
  account_id = "1234"
  ref_id = "123456"
  send_site_setup_emails = "true"
  site_ip = "10.10.10.11"
  force_ssl = "true"

  ####################################################################
  # The remaining following parameters below are designated for updating the site after it has been created.
  ####################################################################
  active = "bypass" # active | bypass
  ignore_ssl = "true"
  acceleration_level = "none"  # off | standard | advanced
  seal_location = "api.seal_location.bottom_right"
  domain_redirect_to_full = "true"
  remove_ssl = "false"
}

####################################################################
# Custom Certificates
####################################################################
resource "incapsula_custom_certificate" "custom-certificate" {
   site_id = "${incapsula_site.example-site.id}"
   certificate = "${file("path/to/your/cert.crt")}"
   private_key = "${file("path/to/your/private_key.key")}"
   passphrase = "yourpassphrase"
}

####################################################################
# Data Center
####################################################################

resource "incapsula_data_center" "example-data-center-test" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center test"
  server_address = "8.8.4.8"
  is_content = "yes"
}

resource "incapsula_data_center" "example-data-center" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center"
  server_address = "8.8.4.4"
  is_content = "yes"
}

# Data Center Servers
resource "incapsula_data_center_server" "example-data-center-server" {
  dc_id = "${incapsula_data_center.example-data-center.id}"
  site_id = "${incapsula_site.example-site.id}"
  server_address = "4.4.4.4"
  is_standby = "no"
}

###################################################################
# Security Rules (ACLs) and Security Rule (ACLs) Exceptions/Whitelists
###################################################################

# api.acl.blacklisted_countries Security Rule (one instance per site)
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
}

# api.threats.blacklisted_countries Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-blacklisted-countries-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}

# api.acl.blacklisted_ips Security Rule (one instance per site)
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_ips"
  ips = "192.168.1.1,192.168.1.2"
}

# api.threats.blacklisted_ips Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-blacklisted-ips-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_ips"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}

# api.acl.blacklisted_urls Security Rule (one instance per site)
resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
  rule_id = "api.acl.blacklisted_urls"
  site_id = "${incapsula_site.example-site.id}"
  url_patterns = "CONTAINS,EQUALS"
  urls = "/alpha,/bravo"
}

# api.acl.blacklisted_urls Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-blacklisted-urls-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_urls"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}

# api.acl.whitelisted_ips Security Rule (one instance per site)
resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
  rule_id = "api.acl.whitelisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.3,192.168.1.4"
}

####################################################################
# Security Rules (WAF) and Security Rule (WAF) Exceptions/Whitelists
####################################################################

# api.threats.backdoor Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url" # (api.threats.action.quarantine_url (default) | api.threats.action.alert | api.threats.action.disabled | api.threats.action.quarantine_url)
}

# api.threats.backdoor Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-backdoor-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  user_agents="myUserAgent"
  parameters="myparam"
}

# api.acl.bot_access_control Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "true" # true | false (optional, default: true)
  challenge_suspected_bots = "true" # true | false (optional, default: true)
}

# api.threats.bot_access_control Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-bot_access-control-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  user_agents="myUserAgent"
}

# api.threats.cross_site_scripting Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-cross-site-scripting-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.cross_site_scripting"
  security_rule_action = "api.threats.action.block_ip"
}

# api.threats.cross_site_scripting Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-cross-site-scripting-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.cross_site_scripting"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  parameters="myparam"
}

# api.acl.ddos Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "api.threats.ddos.activation_mode.on" # (api.threats.ddos.activation_mode.auto | api.threats.ddos.activation_mode.off | api.threats.ddos.activation_mode.on)
  ddos_traffic_threshold = "5000" # valid values are 10, 20, 50, 100, 200, 500, 750, 1000, 2000, 3000, 4000, 5000
}

# api.threats.ddos Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-ddos-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}

# api.acl.illegal_resource_access Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-illegal-resource-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.illegal_resource_access"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

# api.threats.illegal_resource_access Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-illegal-resource-access-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.illegal_resource_access"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  parameters="myparam"
}

# api.acl.remote_file_inclusion Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-remote-file-inclusion-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.remote_file_inclusion"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

# api.threats.remote_file_inclusion Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-remote-file-inclusion-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.remote_file_inclusion"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  user_agents="myUserAgent"
  parameters="myparam"
}

# api.acl.sql_injection Security Rule (one instance per site)
resource "incapsula_waf_security_rule" "example-waf-sql-injection-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.sql_injection"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

# api.threats.sql_injection Security Rule Sample Exception
resource "incapsula_security_rule_exception" "example-waf-sql-injection-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.sql_injection"
  client_apps="488,123"
  countries="JM,US"
  continents="NA,AF"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
}

###################################################################
# Incap Rules
###################################################################

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
