---
layout: "incapsula"
page_title: "Incapsula: incap-rule"
sidebar_current: "docs-incapsula-resource-incap-rule"
description: |-
  Provides a Incapsula Incap Rule resource.
---

# incapsula_incap_rule

Provides a Incapsula Incap Rule resource. 
Incap Rules include security, delivery, and rate rules.

## Example Usage

```hcl
resource "incapsula_incap_rule" "example-incap-rule-alert" {
  name = "Example incap rule alert"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_ALERT"
  filter = "Full-URL == \"/someurl\""
  enabled = true
}

# Incap Rule: Require javascript support
resource "incapsula_incap_rule" "example-incap-rule-require-js-support" {
  name = "Example incap rule require javascript support 3"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_INTRUSIVE_HTML"
  filter = "Full-URL == \"/someurl\""
  enabled = true  
}

# Incap Rule: Block IP
resource "incapsula_incap_rule" "example-incap-rule-block-ip" {
  name = "Example incap rule block ip"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_BLOCK_IP"
  filter = "Full-URL == \"/someurl\""
  enabled = true  
}

# Incap Rule: Block Request
resource "incapsula_incap_rule" "example-incap-rule-block-request" {
  name = "Example incap rule block request"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_BLOCK"
  filter = "Full-URL == \"/someurl\""
  enabled = true
}

# Incap Rule: Block Session
resource "incapsula_incap_rule" "example-incap-rule-block-session" {
  name = "Example incap rule block session"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_BLOCK_USER"
  filter = "Full-URL == \"/someurl\""
  enabled = true
}

# Incap Rule: Delete Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-cookie" {
  name = "Example incap rule delete cookie"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_DELETE_COOKIE"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
  enabled = true
}

# Incap Rule: Delete Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-delete-header" {
  name = "Example incap rule delete header"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_DELETE_HEADER"
  filter = "Full-URL == \"/someurl\""
  rewrite_name = "my_test_header"
  enabled = true
}

# Incap Rule: Forward to Data Center (ADR)
# For a more detailed example of how to reference a data center id, look at datasource incapsula_data_center  
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  name = "Example incap rule forward to data center"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_FORWARD_TO_DC"
  filter = "Full-URL == \"/someurl\""
  dc_id = data.incapsula_data_center.example_content_dc.id
  enabled = true  
}

# Incap Rule: Redirect (ADR)
resource "incapsula_incap_rule" "example-incap-rule-redirect" {
  name = "Example incap rule redirect"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_REDIRECT"
  filter = "Full-URL == \"/someurl\""
  response_code = "302"
  from = "https://site1.com/url1"
  to = "https://site2.com/url2"
  enabled = true
}

# Incap Rule: Require Cookie Support (IncapRule)
resource "incapsula_incap_rule" "example-incap-rule-require-cookie-support" {
  name = "Example incap rule require cookie support"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_RETRY"
  filter = "Full-URL == \"/someurl\""
  enabled = true
}

# Incap Rule: Rewrite Cookie (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-cookie" {
  name = "Example incap rule rewrite cookie"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_REWRITE_COOKIE"
  filter = "Full-URL == \"/someurl\""
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  rewrite_name = "my_cookie_name"
  enabled = true  
}

# Incap Rule: Rewrite Header (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-header" {
  name = "Example incap rule rewrite header"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_REWRITE_HEADER"
  filter = "Full-URL == \"/someurl\""
  add_missing = "true"
  from = "some_optional_value"
  to = "some_new_value"
  rewrite_name = "my_test_header"
  enabled = true  
}

# Incap Rule: Rewrite URL (ADR)
resource "incapsula_incap_rule" "example-incap-rule-rewrite-url" {
  name = "ExampleRewriteURL"
  site_id = incapsula_site.example-site.id
  action = "RULE_ACTION_REWRITE_URL"
  filter = "Full-URL == \"/someurl\""
  from = "*"
  to = "/redirect"
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, `RULE_ACTION_FORWARD_TO_PORT`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.
* `response_code` - (Optional) For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.
* `add_missing` - (Optional) Add cookie or header if it doesn't exist (Rewrite cookie rule only).
* `from` - (Optional) Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.
* `to` - (Optional) Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.
* `rewrite_name` - (Optional) Name of cookie or header to rewrite. Applies only for `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER`.
* `dc_id` - (Optional) Data center to forward request to. Applies only for `RULE_ACTION_FORWARD_TO_DC`.
* `port_forwarding_context` - (Optional) Context for port forwarding. \"Use Port Value\" or \"Use Header Name\". Applies only for `RULE_ACTION_FORWARD_TO_PORT`.
* `port_forwarding_value` - (Optional) Port number or header name for port forwarding. Applies only for `RULE_ACTION_FORWARD_TO_PORT`.
* `rate_context` - (Optional) The context of the rate counter. Possible values `IP` or `Session`. Applies only to rules using `RULE_ACTION_RATE`.
* `rate_interval` - (Optional) The interval in seconds of the rate counter. Possible values is a multiple of `10`; minimum `10` and maximum `300`. Applies only to rules using `RULE_ACTION_RATE`.
* `error_type` - (Optional) The error that triggers the rule. `error.type.all` triggers the rule regardless of the error type. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `error.type.all`, `error.type.connection_timeout`, `error.type.access_denied`, `error.type.parse_req_error`, `error.type.parse_resp_error`, `error.type.connection_failed`, `error.type.deny_and_retry`, `error.type.ssl_failed`, `error.type.deny_and_captcha`, `error.type.2fa_required`, `error.type.no_ssl_config`, `error.type.no_ipv6_config`.
* `error_response_format` - (Optional) The format of the given error response in the error_response_data field. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `json`, `xml`.
* `error_response_data` - (Optional) The response returned when the request matches the filter and is blocked. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`.
* `multiple_deletions` - (Optional) Delete multiple header occurrences. Applies only to rules using `RULE_ACTION_DELETE_HEADER` and `RULE_ACTION_RESPONSE_DELETE_HEADER`.
* `overrideWafAction` - (Optional) The response returned when the request matches the filter and is blocked. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`.
* `overrideWafRule` - (Optional) The action for the override rule. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Incap Rule.

## Import

Incap Rule can be imported using the role site_id and rule_id separated by /, e.g.:

```
$ terraform import incapsula_incap_rule.demo site_id/rule_id
```