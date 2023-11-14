---
layout: "incapsula"
page_title: "Incapsula: delivery_rules_configuration"
sidebar_current: "docs-incapsula-resource-delivery_rules_configuration"
description: |-
  Provides a Incapsula delivery_rules_configuration resource.
---

# delivery_rules_configuration

Provides a Incapsula delivery rules configuration resource.
delivery rules configuration include delivery rules.

## Example Usage
```hcl
# delivery rules: REDIRECT
resource "incapsula_delivery_rules_configuration" "redirect-rules" {
  category = "REDIRECT"
  site_id  = incapsula_site.example-site.id
  rule {
    rule_name     = "New delivery rule",
    filter        = "ASN == 1"
    from          = "/1"
    to            = "/2"
    response_code = "302"
    action        = "RULE_ACTION_REDIRECT"
    enabled       = "true"
  }
}
# delivery rules:  SIMPLIFIED_ REDIRECT
resource "incapsula_delivery_rules_configuration" "simplified-redirect-rules" {
  category = "SIMPLIFIED_REDIRECT"
  site_id = incapsula_site.example-site.id
  rule {
    from = "/1"
    to = "$scheme://www.example.com/$city"
    response_code = "302"
    rule_name = "New delivery rule",
    action = "RULE_ACTION_SIMPLIFIED_REDIRECT"
    enabled = "true"
  }
}
```
## Argument Reference - REDIRECT & SIMPLIFIED_ REDIRECT RULE
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, `RULE_ACTION_FORWARD_TO_PORT`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.
* `response_code` - (Optional) For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.
* `from` - (Optional) Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.
* `to` - (Optional) Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.
```hcl

# delivery rules:  REWRITE
resource "incapsula_delivery_rules_configuration" "rewrite-request-rules" {
  category = "REWRITE"
  site_id = incapsula_site.example-site.id
  rule {
    filter = "ASN == 1"
    from = "cookie1"
    to = "cookie2"
    rewrite_existing = "true"
    add_if_missing = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_REWRITE_COOKIE"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    header_name = "abc"
    from = "header1"
    to = "header2"
    rewrite_existing = "true"
    add_if_missing = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_REWRITE_HEADER"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    from = "/1"
    to = "/2"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_REWRITE_URL"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    header_name = "abc"
    multiple_headers_deletion = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_DELETE_HEADER"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    cookie_name = "abc"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_DELETE_COOKIE"
    enabled = "true"
  }
}
```
## Argument Reference - REWRITE RULE

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, `RULE_ACTION_FORWARD_TO_PORT`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.
* `from` - (Optional) Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.
* `to` - (Optional) Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.
* `add_missing` - (Optional) Add cookie or header if it doesn't exist (Rewrite cookie rule only).
* `rewrite_existing` - (Optional) Rewrite cookie or header if it exists.

```hcl

# delivery rules: REWRITE_RESPONSE
resource "incapsula_delivery_rules_configuration" "rewrite-response-rules" {
  category = "REWRITE_RESPONSE"
  site_id = incapsula_site.example-site.id

  rule {
    filter = "ASN == 1"
    header_name = "abc"
    multiple_headers_deletion = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_RESPONSE_DELETE_HEADER"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    header_name = "abc"
    from = "header1"
    to = "header2"
    rewrite_existing="true"
    add_if_missing = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_RESPONSE_REWRITE_HEADER"
    enabled = "true"
  }
  rule {
  filter = "ASN == 1"
  response_code = "302"
  rule_name = "New delivery rule"
  action = "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE"
  enabled = "true"
    
  }
  rule {
  filter = "ASN == 1"
  error_response_format = "[JSON|XML]"
  error_response_data = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
  error_type = "error.type.all"
  response_code = "400"
  rule_name = "New delivery rule"
  action = "RULE_ACTION_CUSTOM_ERROR_RESPONSE"
  enabled = "true"
  }
}
```
## Argument Reference - REWRITE_RESPONSE

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, `RULE_ACTION_FORWARD_TO_PORT`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.
* `response_code` - (Optional) For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.
* `error_type` - (Optional) The error that triggers the rule. `error.type.all` triggers the rule regardless of the error type. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `error.type.all`, `error.type.connection_timeout`, `error.type.access_denied`, `error.type.parse_req_error`, `error.type.parse_resp_error`, `error.type.connection_failed`, `error.type.deny_and_retry`, `error.type.ssl_failed`, `error.type.deny_and_captcha`, `error.type.2fa_required`, `error.type.no_ssl_config`, `error.type.no_ipv6_config`.
* `error_response_format` - (Optional) The format of the given error response in the error_response_data field. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `json`, `xml`.
* `error_response_data` - (Optional) The response returned when the request matches the filter and is blocked. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`.
* `add_missing` - (Optional) Add cookie or header if it doesn't exist (Rewrite cookie rule only).
* `rewrite_existing` - (Optional) Rewrite cookie or header if it exists.

```hcl
# delivery rules: FORWARD
resource "incapsula_delivery_rules_configuration" "rewrite-forward-rules" {
  category = "FORWARD"
  site_id = incapsula_site.example-site.id
  rule {
    filter = "ASN == 1"
    dc_id = 1234
    rule_name = "New delivery rule",
    action = "RULE_ACTION_FORWARD_TO_DC"
    enabled = "true"
  }
  rule {
    filter = "ASN == 1"
    port_forwarding_context = "[Use Header Name/Use Port Value]"
    port_forwarding_value = 1234
    rule_name = "New delivery rule"
    action = "RULE_ACTION_FORWARD_TO_PORT"
    enabled = "true"
  }
}

```

## Argument Reference - FORWARD

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, `RULE_ACTION_FORWARD_TO_PORT`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.
* `dc_id` - (Optional) Data center to forward request to. Applies only for `RULE_ACTION_FORWARD_TO_DC`.
* `port_forwarding_context` - (Optional) Context for port forwarding. \"Use Port Value\" or \"Use Header Name\". Applies only for `RULE_ACTION_FORWARD_TO_PORT`.
* `port_forwarding_value` - (Optional) Port number or header name for port forwarding. Applies only for `RULE_ACTION_FORWARD_TO_PORT`.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.
## Attributes Reference

## Import

Delivery Rule can be imported using the role site_id and category separated by /, e.g.:

```
$ terraform import delivery_rules_configuration.demo site_id/category
```