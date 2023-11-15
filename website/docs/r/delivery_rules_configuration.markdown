---
layout: "incapsula"
page_title: "Incapsula: delivery_rules_configuration"
sidebar_current: "docs-incapsula-resource-delivery_rules_configuration"
description: |-
  Provides a Incapsula delivery_rules_configuration resource.
---

# incapsula_delivery_rules_configuration

Provides the delivery rules configuration for a specific site. The order of rules execution (a.k.a. priority) is the same as the order they are defined in the resource configuration. 

Currently there are 5 possible types of delivery rule:
* **REDIRECT** - Redirect requests with 30X response.
* **SIMPLIFIED_REDIRECT** - Redirect requests with 30X response. (this category doesn't support condition triggers, and needs to be enbled at the account level before being used)
* **REWRITE** - Modify, add, and remove different request attributes such as URL, headers and cookies.
* **REWRITE_RESPONSE** - Modify, add, and remove different response attributes such as headers, statuc code and error responses.
* **FORWARD** - Forward the request to a specific data-center or port.

**Important Note:** When using this resource, the rule names within each category must be unique. When multiple rules have the same name, the update would fail with an error message specifying the index of the offending rules.


## Example Usage


## `REDIRECT` RULES

```hcl
resource "incapsula_delivery_rules_configuration" "redirect-rules" {
  category = "REDIRECT"
  site_id  = incapsula_site.example-site.id

  rule {
    rule_name     = "New delivery rule",
    filter        = "ASN == 1"
    from          = "*/url"
    to            = "http://www.example.com"
    response_code = "302"
    action        = "RULE_ACTION_REDIRECT"
    enabled       = "true"
  }

  rule {
    ...
  }
}
```

### Argument Reference
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `category` - (Required) Category of rules - `REDIRECT`.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. Possible value:
  * `RULE_ACTION_REDIRECT` - Redirects incoming requests.
* `filter` - (Required) The filter defines the conditions that trigger the rule action.
* `response_code` - (Required) Redirect status code. Valid values are `302`, `301`, `303`, `307`, `308`.
* `from` - (Required) URL pattern to rewrite.
* `to` - (Required) URL pattern to change to.
* `enabled` - (Optional) Boolean that enables the rule. Default value is true.


## `SIMPLIFIED_REDIRECT` RULES

```hcl
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

  rule {
    ...
  }
}
```

### Argument Reference
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `category` - (Required) Category of rules - `SIMPLIFIED_REDIRECT`.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. Possible value:
  * `RULE_ACTION_SIMPLIFIED_REDIRECT` - Redirects incoming requests.
* `response_code` - (Required) Redirect status code. Valid values are `302`, `301`, `303`, `307`, `308`.
* `from` - (Required) URL pattern to rewrite. **Note**: this field must be unique among other rules of the same category.
* `to` - (Required) URL pattern to change to.
* `enabled` - (Optional) Boolean that enables the rule. Default value is true.


## `REWRITE` RULES

```hcl
resource "incapsula_delivery_rules_configuration" "rewrite-request-rules" {
  category = "REWRITE"
  site_id = incapsula_site.example-site.id

  rule {
    filter = "ASN == 1"
    cookie_name = "cookieName",
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
    header_name = "headerName"
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
    from = "/folder1"
    to = "/folder2"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_REWRITE_URL"
    enabled = "true"
  }

  rule {
    filter = "ASN == 1"
    header_name = "headerName"
    multiple_headers_deletion = "false"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_DELETE_HEADER"
    enabled = "true"
  }

  rule {
    filter = "ASN == 1"
    cookie_name = "cookieName"
    rule_name = "New delivery rule"
    action = "RULE_ACTION_DELETE_COOKIE"
    enabled = "true"
  }
}
```

### Argument Reference

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `category` - (Required) Category of rules - `REWRITE`.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. Possible values:
  * `RULE_ACTION_REWRITE_HEADER` - Modify header of incoming request
  * `RULE_ACTION_REWRITE_COOKIE` - Modify cookie of incoming request
  * `RULE_ACTION_REWRITE_URL` - Modify URL of incoming request
  * `RULE_ACTION_DELETE_HEADER` - delete header of incoming request
  * `RULE_ACTION_DELETE_COOKIE` - delete cookie of incoming request
* `filter` - (Optional) The filter defines the conditions that trigger the rule action.
* `cookie_name` - (Required) The cookie name that the rules applies to.
* `header_name` - (Required) The header name that the rules applies to.
* `from` - (Optional) Header/Cookie/URL pattern to rewrite.
* `to` - (Required) Header/Cookie/URL pattern to change to.
* `add_missing` - (Optional) When rewriting cookie or header, add it if it doesn't exist.
* `rewrite_existing` - (Optional) Rewrite cookie or header even if it exists already.
* `multiple_headers_deletion` - (Optional) Delete all header occurrences.
* `enabled` - (Optional) Boolean that enables the rule. Default value is true.


## `REWRITE_RESPONSE` RULES

```hcl
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

### Argument Reference

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `category` - (Required) Category of rules - `REWRITE`.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. Possible values:
  * `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Modify header of outgoing response
  * `RULE_ACTION_RESPONSE_DELETE_HEADER` - Remove header from outgoing response
  * `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` - Modify HTTP status code of outgoing response
  * `RULE_ACTION_CUSTOM_ERROR_RESPONSE` - Set custom template for various error responses
* `filter` - (Optional) The filter defines the conditions that trigger the rule action.
* `header_name` - (Required) The header name that the rules applies to.
* `from` - (Optional) Header pattern to rewrite.
* `to` - (Required) Header pattern to change to.
* `add_missing` - (Optional) When rewriting a header, add it if it doesn't exist.
* `rewrite_existing` - (Optional) Rewrite a header even it if it exists already.
* `multiple_headers_deletion` - (Optional) Delete all header occurrences.
* `response_code` - (Required) HTTP status code. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.
* `error_type` - (Optional) The error that triggers the rule. `error.type.all` triggers the rule regardless of the error type. Possible values: `error.type.all`, `error.type.connection_timeout`, `error.type.access_denied`, `error.type.parse_req_error`, `error.type.parse_resp_error`, `error.type.connection_failed`, `error.type.deny_and_retry`, `error.type.ssl_failed`, `error.type.deny_and_captcha`, `error.type.2fa_required`, `error.type.no_ssl_config`, `error.type.no_ipv6_config`.
* `error_response_format` - (Optional) The format of the given error response in the error_response_data field. Possible values: `json`, `xml`.
* `error_response_data` - (Optional) The response returned when the request matches the filter and is blocked.
* `enabled` - (Optional) Boolean that enables the rule. Default value is true.


## `FORWARD` RULES

```hcl
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

### Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `category` - (Required) Category of rules - `FORWARD`.
* `rule_name` - (Required) Rule name.
* `action` - (Required) Rule action. Possible values:
  * `RULE_ACTION_FORWARD_TO_DC` - Forward requests to a specific data-center
  * `RULE_ACTION_FORWARD_TO_PORT` - Forward requests to a specific port
* `filter` - (Optional) The filter defines the conditions that trigger the rule action.
* `dc_id` - (Required) ID of the data center to forward the request to.
* `port_forwarding_context` - (Required) Context for port forwarding. Possible values: `Use Port Value` or `Use Header Name`.
* `port_forwarding_value` - (Required) Port number or header name for port forwarding. When using a header, its value should be of format IP:PORT.
* `enabled` - (Optional) Boolean that enables the rule. Possible values: true, false. Default value is true.


## Import

Delivery rules configuration can be imported using the role site_id and category separated by /, e.g.:

```
$ terraform import delivery_rules_configuration.demo site_id/category
```