---
layout: "incapsula"
page_title: "Incapsula: incap-api-security-site-config"
sidebar_current: "docs-incapsula-resource-api-security-site-config"
description: |- Provides a Incapsula API Security Site Config resource.
---

# incapsula_api_security_site_config

Provides an Incapsula API Security Site Config resource. //todo API Security Site Config include violation actions.

## Example Usage

```hcl
resource "incapsula_api_security_site_config" "example-terraform-api-security-site-config" {
  	site_id = 123
  	api_only_site = "true"
  	non_api_request_violation_action = "ALERT_ONLY"
  	invalid_url_violation_action = "ALERT_ONLY"
  	invalid_method_violation_action = "ALERT_ONLY"
  	missing_param_violation_action = "ALERT_ONLY"
  	invalid_param_value_violation_action = "ALERT_ONLY"
  	invalid_param_name_violation_action = "ALERT_ONLY"
	is_automatic_discovery_api_integration_enabled = "false"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `account_id` - (Optional) Numeric identifier of the account that owns the site.
* `site_name` - (Optional) Alphabetic identifier of the site to operate on .
* `api_only_site` - (Optional)
* `discovery_enabled` - (Optional) Parameter activates/deactivates automatic API discovery function.
* `non_api_request_violation_action` - (Optional) .
* `invalid_url_violation_action` - (Optional) The action taken when an invalid URL Violation occurs. Actions
  available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_method_violation_action` - (Optional) The action taken when an invalid method Violation occurs. Actions
  available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `missing_param_violation_action` - (Optional) The action taken when a missing parameter Violation occurs. Actions
  available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_param_value_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Actions available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_param_name_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Actions available: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `last_modified` - (Optional) The latest date when the resource was updated.
* `is_automatic_discovery_api_integration_enabled` - (Optional) Parameter shows whether automatic API discovery is
  enabled.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the API Security Site Configuration.

## Import

API Security Site Configuration can be imported using the site_id

```
$ terraform import incapsula_api_security_site_config.demo_site_config site_id
```
