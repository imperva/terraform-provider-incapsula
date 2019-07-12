
variable "id" {
  default = ""
}
variable "key" {
  default = ""
}
provider "incapsula" {
  api_id = "${var.id}"
  api_key = "${var.key}"
}

resource "incapsula_site" "example-site" {
  ####################################################################
  # # The first 6 parameters listed below are designated for initially creating the site in Incapsula.
  ####################################################################
  domain = "mojo.beer.center" //Good for adding site and in response object
  account_id = "2398" //Good for adding site and in response object
  ref_id = "666" //Good for adding and updating site and in response object
  send_site_setup_emails = "true" //Good for adding site, not availble in site onject
  site_ip = "75.80.36.61" //Good for adding and updating site and in response object
  force_ssl = "true" //Good for adding site and in response object

  # # The following 4 parameters are not able to be implemented for TF
  # naked_domain_san = "false" //TODO: Not available in TF site resource
  # wildcard_san = "false" //TODO: Not available in TF site resource
  # log_level = "full" //TODO: Fails with "reason": "Logs are not supported due to feature restrictions"
  # logs_account_id = "1034421" //TODO: Fails with "reason": "The provided logs_account_id invalid"

  ####################################################################
  # # The remaining following parameters below are designated for updating the site after it has been created.
  # # TODO: This entire site update part needs to have logic applied resource; one failure will return out of the resource_incapsula_site.go.resourceSiteUpdate function...
  ####################################################################
  # active = "bypass" # active | bypass //Good for updating site and in response object
  # ignore_ssl = "true" //Good for updating site and in response object statusEnum key from pending-select-approver to pending_ssl_approval
  # acceleration_level = "none"  # off | standard | advanced //Good for updating site and in response object
  # seal_location = "api.seal_location.bottom_right"  //Good for updating site and in response object
  # domain_redirect_to_full = "true" //TODO: Good for updating site but not shown in reponse or in UI
  # remove_ssl = "false"
  # approver = "your@email.com" //TODO: This call fails when trying to apply on update.
      //Error from Incapsula service when updating site for siteID 12608430:
      //{"res":4202,"res_message":"Domain_email invalid","debug_info":{"approver":"approver email value is not a valid email value from the site\u0027s domain_emails","id-info":"13007"}}
  # domain_validation = "dns" //TODO: This call fails after being set by Imperva and cannnot be changed. Default is DNS; call will error on HTML and pass on EMAIL but no change.
      //Error from Incapsula service when updating site for siteID 12608430:
      //{"res":1,"res_message":"Unexpected error","debug_info":{"id-info":"13007","domain_html":"An error occured while trying to modify site\u0027s domain validation"}}
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

Data Center Servers
resource "incapsula_data_center_servers" "example-data-center-servers" {
  dc_id = "${incapsula_data_center.example-data-center.id}"
  site_id = "${incapsula_site.example-site.id}"
  server_address = "4.4.4.4"
  is_standby = "no"
}

# ###################################################################
# Security Rules (WAF)
# ###################################################################

resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url" # (api.threats.action.quarantine_url (default) | api.threats.action.alert | api.threats.action.disabled | api.threats.action.quarantine_url)
}

resource "incapsula_waf_security_rule" "example-waf-cross-site-scripting-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.cross_site_scripting"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

resource "incapsula_waf_security_rule" "example-waf-illegal-resource-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.illegal_resource_access"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

resource "incapsula_waf_security_rule" "example-waf-remote-file-inclusion-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.remote_file_inclusion"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

resource "incapsula_waf_security_rule" "example-waf-sql-injection-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.sql_injection"
  security_rule_action = "api.threats.action.block_ip" # (api.threats.action.disabled | api.threats.action.alert | api.threats.action.block_request | api.threats.action.block_user | api.threats.action.block_ip)
}

resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "true" # true | false (optional, default: true)
  challenge_suspected_bots = "true" # true | false (optional, default: true)
}

resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.ddos"
  activation_mode = "api.threats.ddos.activation_mode.on" # (api.threats.ddos.activation_mode.auto | api.threats.ddos.activation_mode.off | api.threats.ddos.activation_mode.on)
  ddos_traffic_threshold = "5000" # valid values are 10, 20, 50, 100, 200, 500, 750, 1000, 2000, 3000, 4000, 5000
}

# ###################################################################
# Security Rules (ACLs)
# ###################################################################

Security Rule: Country
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
}

Security Rule: Blacklist IP
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_ips"
  ips = "192.168.1.1,192.168.1.2"
}

Security Rule: Blacklist IP Exception
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule_exception" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.1,192.168.1.2"
  urls = "/myurl,/myurl2"
  url_patterns = "EQUALS,CONTAINS"
  countries = "JM,US"
  client_apps= "488,123"
}

Security Rule: URL
resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
  rule_id = "api.acl.blacklisted_urls"
  site_id = "${incapsula_site.example-site.id}"
  url_patterns = "CONTAINS,EQUALS"
  urls = "/alpha,/bravo"
}

Security Rule: Whitelist IP
resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
  rule_id = "api.acl.whitelisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.3,192.168.1.4"
}

# ###################################################################
# Incap Rules
# ###################################################################

Incap Rule: Alert
resource "incapsula_incap_rule" "example-incap-rule-alert" {
  priority = "1"
  name = "Example incap rule alert"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  priority = "1"
  name = "Example incap rule require javascript support 3"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_INTRUSIVE_HTML"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  priority = "1"
  name = "Example incap rule block ip"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK_IP"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  priority = "1"
  name = "Example incap rule block request"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Block Session
resource "incapsula_incap_rule" "example-incap-rule-block-session" {
  priority = "1"
  name = "Example incap rule block session"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_BLOCK_USER"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Delete Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
  priority = "1"
  name = "Example incap rule delete cookie"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_DELETE_COOKIE"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

Incap Rule: Delete Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
  priority = "1"
  name = "Example incap rule delete header"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_DELETE_HEADER"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  priority = "1"
  name = "Example incap rule forward to data center"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_FORWARD_TO_DC"
  filter = "Full-URL == \"/someurl\""
  dc_id = "${incapsula_data_center.example-data-center.id}"
  allow_caching = "false"
}

Incap Rule: Redirect (ADR)
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

Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  priority = "1"
  name = "Example incap rule require cookie support"
  site_id = "${incapsula_site.example-site.id}"
  action = "RULE_ACTION_RETRY"
  filter = "Full-URL == \"/someurl\""
}

Incap Rule: Rewrite Cookie (ADR)
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

Incap Rule: Rewrite Header (ADR)
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

Incap Rule: Rewrite URL (ADR)
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

