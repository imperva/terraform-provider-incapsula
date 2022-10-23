---
layout: "incapsula"
page_title: "Incapsula: incap-api-security-api-config"
sidebar_current: "docs-incapsula-resource-api-security-api-config"
description: |- Provides a Incapsula API Security API Config resource.
---

# incapsula_api_security_api_config

Provides an Incapsula API Security API Config resource.

API Security API Config include violation actions set for specific API.

## Example Usage

```hcl
resource "incapsula_api_security_api_config" "demo-terraform-api-security-api-config" {
	site_id = incapsula_site.example-site.id
	api_specification = "${file("path/to/your/swagger/file.yaml")}"
	invalid_url_violation_action = "IGNORE"
	invalid_method_violation_action = "BLOCK_USER"
	missing_param_violation_action = "BLOCK_IP"
	invalid_param_value_violation_action = "BLOCK_REQUEST"
	invalid_param_name_violation_action = "ALERT_ONLY"
	description = "your site API description"
	base_path = "/base/path"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `api_specification` - (Required) The API specification document content. The supported format is OAS2 or OAS3.
* `invalid_url_violation_action` - (Optional) The action taken when an invalid URL Violation occurs. Possible values:
  `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`, `DEFAULT`. Assigning `DEFAULT` will inherit the
  action from parent object.
* `invalid_method_violation_action` - (Optional) The action taken when an invalid method Violation occurs. Possible
  values:
  `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`, `DEFAULT`. Assigning `DEFAULT` will inherit the
  action from parent object.
* `missing_param_violation_action` - (Optional) The action taken when a missing parameter Violation occurs. Possible
  values:
  `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`. Assigning `DEFAULT` will inherit the
  action from parent object.
  
  > **NOTE:** `invalid_param_name_violation_action` parameter is currently not supported. Please do not use/change value.

* `invalid_param_value_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Assigning `DEFAULT` will inherit the action from parent object. Possible values: `ALERT_ONLY`, `BLOCK_REQUEST`
  , `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`. Assigning `DEFAULT` will inherit the action from parent object.
* `invalid_param_name_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Possible values: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`,`DEFAULT`.
* `description` - (Optional) A description that will help recognize the API in the dashboard.
* `base_path` - (Optional) Override the spec basePath / server base path with this value.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the API Security Site Configuration.
* `host_name` - The API's host name
* `last_modified` - (Optional) The last modified timestamp.

## Import

API Security API Configuration can be imported using the site_id and then api_id (id) separated by /, e.g.:

```
$ terraform import incapsula_api_security_api_config.example-terraform-api-security-api-config 1234/100200

```
