---
layout: "incapsula"
page_title: "Incapsula: acl-security-rule"
sidebar_current: "docs-incapsula-resource-acl-security-rule"
description: |-
  Provides a Incapsula ACL Security Rule resource.
---

# incapsula_acl_security_rule

Provides a Incapsula ACL Security Rule resource. 
ACL Security Rules allow for blacklisting or whitelisting countries, IP addresses, and URLs.

## Example Usage

```hcl
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url"
}

resource "incapsula_security_rule_exception" "example-waf-backdoor-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.backdoor"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  ips          = "1.2.3.6,1.2.3.7"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
  user_agents  = "myUserAgent"
  parameters   = "myparam"
}

resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id                  = incapsula_site.example-site.id
  rule_id                  = "api.threats.bot_access_control"
  block_bad_bots           = "true"
  challenge_suspected_bots = "true"
}

resource "incapsula_security_rule_exception" "example-waf-bot_access-control-rule-exception" {
  site_id          = incapsula_site.example-site.id
  rule_id          = "api.threats.bot_access_control"
  client_app_types = "DataScraper,"
  ips              = "1.2.3.6,1.2.3.7"
  url_patterns     = "EQUALS,CONTAINS"
  urls             = "/myurl,/myurl2"
  user_agents      = "myUserAgent"
}

resource "incapsula_waf_security_rule" "example-waf-cross-site-scripting-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.cross_site_scripting"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_security_rule_exception" "example-waf-cross-site-scripting-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.cross_site_scripting"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
  parameters   = "myparam"
}

resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id                = incapsula_site.example-site.id
  rule_id                = "api.threats.ddos"
  activation_mode        = "api.threats.ddos.activation_mode.on"
  ddos_traffic_threshold = "5000"
}

resource "incapsula_security_rule_exception" "example-waf-ddos-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.ddos"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  ips          = "1.2.3.6,1.2.3.7"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
}

resource "incapsula_waf_security_rule" "example-waf-illegal-resource-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.illegal_resource_access"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_security_rule_exception" "example-waf-illegal-resource-access-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.illegal_resource_access"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  ips          = "1.2.3.6,1.2.3.7"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
  parameters   = "myparam"
}

resource "incapsula_waf_security_rule" "example-waf-remote-file-inclusion-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.remote_file_inclusion"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_security_rule_exception" "example-waf-remote-file-inclusion-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.remote_file_inclusion"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  ips          = "1.2.3.6,1.2.3.7"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
  user_agents  = "myUserAgent"
  parameters   = "myparam"
}

resource "incapsula_waf_security_rule" "example-waf-sql-injection-rule" {
  site_id              = incapsula_site.example-site.id
  rule_id              = "api.threats.sql_injection"
  security_rule_action = "api.threats.action.block_ip"
}

resource "incapsula_security_rule_exception" "example-waf-sql-injection-rule-exception" {
  site_id      = incapsula_site.example-site.id
  rule_id      = "api.threats.sql_injection"
  client_apps  = "488,123"
  countries    = "JM,US"
  continents   = "NA,AF"
  ips          = "1.2.3.6,1.2.3.7"
  url_patterns = "EQUALS,CONTAINS"
  urls         = "/myurl,/myurl2"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_id` - (Required) The id of the acl, e.g api.acl.blacklisted_ips. Options are `api.acl.blacklisted_countries`, `api.acl.blacklisted_urls`, `api.acl.blacklisted_ips`, and `api.acl.whitelisted_ips`.
* `continents` - (Optional) A comma separated list of continent codes.
* `countries` - (Optional) A comma separated list of country codes.
* `ips` - (Optional) A comma separated list of IPs or IP ranges, e.g: `192.168.1.1`, `192.168.1.1-192.168.1.100` or `192.168.1.1/24`.
* `urls` - (Optional) A comma separated list of resource paths.
* `url_patterns` - (Optional) The patterns should be in accordance with the matching urls sent by the urls parameter. Options are `CONTAINS`, `EQUALS`, `PREFIX`, `SUFFIX`, `NOT_EQUALS`, `NOT_CONTAIN`, `NOT_PREFIX`,and `NOT_SUFFIX`.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the ACL security rule.
