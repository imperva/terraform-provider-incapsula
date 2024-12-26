provider "incapsula" {
  api_id  = "your_api_id"
  api_key = "your_api_key"
}

variable "site_id" {
  default = "your_site_id_here"
}

####################################################################
# Data Centers
####################################################################

resource "incapsula_data_center" "example-data-center-test" {
  site_id        = incapsula_site.example-site.id
  name           = "Example data center test"
  server_address = "8.8.4.8"
  is_content     = "true"
}

resource "incapsula_data_center" "example-data-center" {
  site_id        = incapsula_site.example-site.id
  name           = "Example data center"
  server_address = "8.8.4.4"
  is_content     = "true"
}

resource "incapsula_data_center_server" "example-data-center-server" {
  dc_id          = incapsula_data_center.example-data-center.id
  site_id        = incapsula_site.example-site.id
  server_address = "4.4.4.4"
  is_standby     = "false"
}

####################################################################
# Security Rules (WAF)
####################################################################

resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url"
}

resource "incapsula_waf_security_rule" "example-waf-cross-site-scripting-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.cross_site_scripting"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_waf_security_rule" "example-waf-illegal-resource-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.illegal_resource_access"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_waf_security_rule" "example-waf-remote-file-inclusion-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.remote_file_inclusion"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_waf_security_rule" "example-waf-sql-injection-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.sql_injection"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id                  = incapsula_site.example-site.id
  rule_id                  = "api.threats.bot_access_control"
  block_bad_bots           = "true"
  challenge_suspected_bots = "true"
}

resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id                = incapsula_site.example-site.id
  rule_id                = "api.threats.ddos"
  activation_mode        = "api.threats.ddos.activation_mode.on"
  ddos_traffic_threshold = "5000"
  unknown_clients_challenge = "none"
  block_non_essential_bots  = "false"
}

####################################################################
# Incap Rules
####################################################################

# Incap Rule: Alert
resource "incapsula_incap_rule" "example-incap-rule-alert" {
  name    = "Example incap rule alert"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_ALERT"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  name    = "Example incap rule require javascript support 3"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_INTRUSIVE_HTML"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  name    = "Example incap rule block ip"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_BLOCK_IP"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  name    = "Example incap rule block request"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_BLOCK"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Block Session
resource "incapsula_incap_rule" "example-incap-rule-block-session" {
  name    = "Example incap rule block session"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_BLOCK_USER"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Delete Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
  name         = "Example incap rule delete cookie"
  site_id      = incapsula_site.example-site.id
  action       = "RULE_ACTION_DELETE_COOKIE"
  filter       = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

# Incap Rule: Delete Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
  name         = "Example incap rule delete header"
  site_id      = incapsula_site.example-site.id
  action       = "RULE_ACTION_DELETE_HEADER"
  filter       = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
}

# Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  name    = "Example incap rule forward to data center"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_FORWARD_TO_DC"
  filter  = "Full-URL == \"/someurl\""
  dc_id   = incapsula_data_center.example-data-center.id
}

# Incap Rule: Redirect (ADR)
resource "incapsula_incap_rule" "example-incap-rule-redirect" {
  name          = "Example incap rule redirect"
  site_id       = incapsula_site.example-site.id
  action        = "RULE_ACTION_REDIRECT"
  filter        = "Full-URL == \"/someurl\""
  response_code = "302"
  from          = "https://site1.com/url1"
  to            = "https://site2.com/url2"
}

# Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  name    = "Example incap rule require cookie support"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_RETRY"
  filter  = "Full-URL == \"/someurl\""
}

# Incap Rule: Rewrite Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-cookie" {
  name             = "Example incap rule rewrite cookie"
  site_id          = incapsula_site.example-site.id
  action           = "RULE_ACTION_REWRITE_COOKIE"
  filter           = "Full-URL == \"/someurl\""
  add_missing      = "true"
  rewrite_existing = "true"
  from             = "some_optional_value"
  to               = "some_new_value"
  rewrite_name     = "my_cookie_name"
}

# Incap Rule: Rewrite Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-header" {
  name             = "Example incap rule rewrite header"
  site_id          = incapsula_site.example-site.id
  action           = "RULE_ACTION_REWRITE_HEADER"
  filter           = "Full-URL == \"/someurl\""
  add_missing      = "true"
  rewrite_existing = "true"
  from             = "some_optional_value"
  to               = "some_new_value"
  rewrite_name     = "my_test_header"
}

# Incap Rule: Rewrite URL (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-url" {
  name         = "ExampleRewriteURL"
  site_id      = incapsula_site.example-site.id
  action       = "RULE_ACTION_REWRITE_URL"
  filter       = "Full-URL == \"/someurl\""
  add_missing  = "true"
  from         = "*"
  to           = "/redirect"
  rewrite_name = "my_test_header"
}
