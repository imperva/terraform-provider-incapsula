---
layout: "incapsula"
page_title: "Incapsula: incap-api-security-site-config"
sidebar_current: "docs-incapsula-resource-api-security-site-config"
description: |- Provides a Incapsula API Security Site Config resource.
---

# incapsula_api_security_site_config

Provides an Incapsula API Security Site Config resource.

## Example Usage

```hcl
resource "incapsula_api_security_site_config" "demo-terraform-api-security-site-config" {
  	site_id = incapsula_site.example-site.id
  	is_automatic_discovery_api_integration_enabled = false
  	is_api_only_site = true
  	non_api_request_violation_action = "ALERT_ONLY"
  	invalid_url_violation_action = "BLOCK_IP"
  	invalid_method_violation_action = "BLOCK_REQUEST"
  	missing_param_violation_action = "BLOCK_IP"
  	invalid_param_value_violation_action = "IGNORE"
  	invalid_param_name_violation_action = "ALERT_ONLY"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `is_automatic_discovery_api_integration_enabled` - (Optional) Parameter shows whether automatic API discovery API
  Integration is enabled.
* `invalid_url_violation_action` - (Optional) The action taken when an invalid URL Violation occurs. Possible
  values: `ALERT_ONLY` (default value), `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_method_violation_action` - (Optional) The action taken when an invalid method Violation occurs. Possible
  values: `ALERT_ONLY` (default value), `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `missing_param_violation_action` - (Optional) The action taken when a missing parameter Violation occurs. Possible
  values: `ALERT_ONLY` (default value), `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_param_value_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Possible values: `ALERT_ONLY` (default value), `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `invalid_param_name_violation_action` - (Optional) The action taken when an invalid parameter value Violation occurs.
  Possible values: `ALERT_ONLY` (default value), `BLOCK_REQUEST`, `BLOCK_USER`, `BLOCK_IP`, `IGNORE`.
* `is_api_only_site` - (Optional) Apply positive security model for all traffic on the site. Applying the positive
  security model for all traffic on the site may lead to undesired request blocking.
* `non_api_request_violation_action` - (Optional) Action to be taken for traffic on the site that does not target the
  uploaded APIs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE. This parameter is required
  when `is_api_only_site` is set true. Possible values: `ALERT_ONLY`, `BLOCK_REQUEST`, `BLOCK_USER`
  , `BLOCK_IP`, `IGNORE`.
## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the API Security Site Configuration.
* `last_modified` - The last modified timestamp.

## Import

API Security Site Configuration can be imported using the site_id

```
$ terraform import incapsula_api_security_site_config.demo_site_config 1234
```
