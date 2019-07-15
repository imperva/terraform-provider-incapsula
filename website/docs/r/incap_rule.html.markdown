---
layout: "incapsula"
page_title: "Incapsula: incap-rule"
sidebar_current: "docs-incapsula-resource-incap-rule"
description: |-
  Provides a Incapsula Incap Rule resource.
---

# incapsula_site

Provides a Incapsula Incap Rule resource. 
Incap Rules include security, delivery, and rate rules.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) Rule name.
* `priority` - (Required) Priority for the selected rule.
* `enabled` - (Optional) Enables the rule.
* `action` - (Optional) Rule action. See the possible values in the API documentation.
* `filter` - (Optional) Rule will trigger only a request that matches this filter. The filter may contain up to 400 characters.
* `allow_caching` - (Optional) Allows rule caching.
* `dc_id` - (Optional) Data center to forward request to. Applies only for RULE_ACTION_FORWARD_TO_DC.
* `from` - (Optional) The pattern to rewrite.
* `to` - (Optional) The pattern to change to.
* `response_code` - (Optional) Redirect rule's response code. Valid values are 302, 301, 303, 307, 308.
* `add_missing` - (Optional) Add cookie or header if it doesn't exist (Rewrite cookie rule only).
* `rewrite_name` - (Optional) Name of cookie or header to rewrite. Applies only for RULE_ACTION_REWRITE_COOKIE and RULE_ACTION_REWRITE_HEADER.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Incap Rule.
