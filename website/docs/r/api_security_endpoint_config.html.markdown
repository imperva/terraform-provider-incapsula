---
layout: "incapsula"
page_title: "Incapsula: incap-api-security-endpoint-config"
sidebar_current: "docs-incapsula-resource-api-security-endpoint-config"
description: |- Provides a Incapsula API Security Endpoint Config resource.
---

# incapsula_api_security_endpoint_config

Provides an Incapsula API Security Endpoint Config resource.

API Security Endpoint Config include violation actions set for specific endpoints.

## Example Usage

```hcl
resource "incapsula_api_security_endpoint_config" "demo-api-security-endpoint-config" {
    api_id = incapsula_api_security_api_config.demo_api_security_api_config.id
    path = "/endpoint/unit/{id}"
	method = "GET"
	invalid_param_value_violation_action = "ALERT_ONLY"
	missing_param_violation_action = "BLOCK_IP"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) Numeric identifier of the API Security API Configuration to operate on.
* `path` - (Required) An URL path of specific Endpoint.
* `method` - (Required) HTTP method that describes a specific Endpoint.
* `invalid_param_value_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Possible values: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`, `DEFAULT`. Assigning `DEFAULT`
  will inherit the action from parent object.
* `invalid_param_name_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Possible values: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`, `DEFAULT`. Assigning `DEFAULT`
  will inherit the action from parent object.
* `missing_param_violation_action` - (Optional) The action taken when a missing parameter Violation occurs. Possible
  values: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`, `DEFAULT`. Assigning `DEFAULT` will inherit
  the action from parent object.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Endpoint for the API Security Endpoint Configuration.

## Import

API Security Endpoint Configuration can be imported using api_id and then endpoint_id (id) separated by /, e.g.
```
$ terraform import incapsula_api_security_endpoint_config.example-terraform-api-security-endpoint-config 100200/1122

```
