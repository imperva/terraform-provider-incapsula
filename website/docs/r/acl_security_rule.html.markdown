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
resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
  site_id = "${incapsula_site.example-site.id}"
  rule_id = "api.acl.blacklisted_countries"
  countries = "AI,AN"
}

resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.1,192.168.1.2"
}

resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
  rule_id = "api.acl.blacklisted_urls"
  site_id = "${incapsula_site.example-site.id}"
  url_patterns = "CONTAINS,EQUALS"
  urls = "/alpha,/bravo"
}

resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
  rule_id = "api.acl.whitelisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.3,192.168.1.4"
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
