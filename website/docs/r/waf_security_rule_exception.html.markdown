---
layout: "incapsula"
page_title: "Incapsula: waf-security-rule-exception"
sidebar_current: "docs-incapsula-resource-waf-security-rule-exception"
description: |-
  Provides a Incapsula WAF Security Rule Exception resource.
---

# incapsula_site

Provides a Incapsula WAF Security Rule Exception resource.  Important to note that based on the rule_id the exception is being created for, that there are rule specific parameters that apply to each.  The example resources listed below include all of the supported resources for each rule_id or rule type, although it is not required to use all listed parameters when creating an exception. Exception parameters are optional but at least one is required.

## Example Usage

```hcl
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

resource "incapsula_security_rule_exception" "example-waf-bot_access-control-rule-exception" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.threats.bot_access_control"
  client_app_types="DataScraper,"
  ips="1.2.3.6,1.2.3.7"
  url_patterns="EQUALS,CONTAINS"
  urls="/myurl,/myurl2"
  user_agents="myUserAgent"
}

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
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_id` - (Required) The identifier of the WAF rule, e.g api.threats.cross_site_scripting.
* `client_app_types` - (Optional) A comma separated list of client application types.
* `client_apps` - (Optional) A comma separated list of client application IDs.
* `countries` - (Optional) A comma separated list of country codes.
* `continents` - (Optional) A comma separated list of continent codes.
* `ips=` - (Optional) A comma separated list of IPs or IP ranges, e.g: 192.168.1.1, 192.168.1.1-192.168.1.100 or 192.168.1.1/24
* `urls=` - (Optional) A comma separated list of resource paths. For example, /home and /admin/index.html are resource paths, while http://www.example.com/home is not. Each URL should be encoded separately using percent encoding as specified by RFC 3986 (http://tools.ietf.org/html/rfc3986#section-2.1).  An empty URL list will remove all URLs. urls="/someurl1,/path/to/my/resource/2.html,/some/url/3"
* `url_patterns` - (Optional) A comma separated list of patters that correlate to the list of urls.  url_patterns are required if you have urls specified, and patters are applied in the order specified and map literally to the list of urls. Supported values are: contains,equals,prefix,suffix,not_equals,not_contain,not_prefix,not_suffix.  Example of how to apply url_patters to the three urls listed above in order: url_patters="prefix,equals,prefix".  
* `user_agents` - (Optional) A comma separated list of encoded user agents.
* `parameters` - (Optional) A comma separated list of encoded parameters.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Incap Rule.
