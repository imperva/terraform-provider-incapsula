---
layout: "incapsula"
page_title: "Incapsula: incap-api-security-endpoint-config"
sidebar_current: "docs-incapsula-resource-api-security-endpoint-config"
description: |- Provides a Incapsula API Security API Config resource.
---

# incapsula_api_security_site_config

Provides an Incapsula API Security Endpoint Config resource.

API Security Endpoint Config include violation actions set for specific endpoints.

## Example Usage

```hcl
resource "incapsula_api_security_endpoint_config" "example-api-security-endpoint-config" {
    api_id = 8020
	invalid_param_name_violation_action = "IGNORE"
	invalid_param_value_violation_action = "IGNORE"
	missing_param_violation_action = "IGNORE"
	path = "/endpoint/unit/{id}"
	method = "GET"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) Numeric identifier of the API to operate on.
* `invalid_param_value_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Actions available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`.
* `invalid_param_name_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Actions available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`.
* `missing_param_violation_action` - (Optional) The action taken when a missing parameter Violation occurs. Actions
  available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`.
* `description` - (Optional) A description that will help recognize the API in the dashboard.
* `path` - (Optional) An URL path of specific endpoint.
* `method` - (Optional) HTTP method that describes a specific endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Endpoint for the API Security Site Configuration.

## Import

API Security Endpoint Configuration can be imported using the site_id, and values of "method"  and "path" fields
separated by /. Path should be separated inside by _, not by /, e.g.:

```
$ terraform import incapsula_api_security_endpoint_config.demo-terraform-api-security-api-config site_id/method/your_path_separated_by_underscore

```
