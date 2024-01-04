---
layout: "incapsula"
page_title: "Incapsula: simplified_redirect_rules_configuration"
sidebar_current: "docs-incapsula-resource-simplified_redirect_rules_configuration"
description: |-
  Provides a Incapsula simplified_redirect_rules_configuration resource.
---

# incapsula_simplified_redirect_rules_configuration

Provides simplified redirect rules. The functionality is similar to REDIRECT rules in `incapsula_delivery_rules_configuration` but with the following limitations:
* They cannot specify a `filter` (logical predicate tested before executing the rule)
* There cannot be 2 rules with the same origin (`from`) argument
* Wildcards cannot be used in the origin (`from`) argument
* Consequently, they cannot be assigned a priority value, since there can be at most one rule for any origin URL

Due to their simplicity, the limits of simplified redirect rules per sites is much higher than normal rules (currently 20,000, compared to only 500 for other delivery rules).

**Note:** Simplified redirect rules should be enabled in the account settings before being able to use them. 

## Example Usage

```hcl
resource "incapsula_simplified_redirect_rules_configuration" "simplified-redirect-rules" {
  site_id = incapsula_site.example-site.id
  rule {
    from = "/url/1"
    to = "$scheme://www.example.com/$city"
    response_code = "302"
    rule_name = "rule 1",
    enabled = "true"
  }
  rule {
    from = "/url/2"
    to = "http://www.google.com"
    response_code = "302"
    rule_name = "rule 2",
    enabled = "true"
  }
}
```

### Argument Reference
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `response_code` - (Required) Redirect status code. Valid values are `302`, `301`, `303`, `307`, `308`.
* `from` - (Required) URL to rewrite.
* `to` - (Required) URL to change to.
* `enabled` - (Optional) Boolean that enables the rule. Default value is `true`.

## Import

Simplified redirect rules configuration can be imported using the site_id, e.g.:

```
$ terraform import incapsula_simplified_redirect_rules_configuration.demo site_id
```
