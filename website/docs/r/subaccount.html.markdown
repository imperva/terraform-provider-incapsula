---
layout: "incapsula"
page_title: "Incapsula: subaccount"
sidebar_current: "docs-incapsula-resource-subaccount"
description: |-
  Provides a Incapsula SubAccount resource.
---

# incapsula_subaccount

Provides a Incapsula SubAccount resource. 

## Example Usage

```hcl
resource "incapsula_subaccount" "example-subaccount" {
  sub_account_name                   = "Example SubAccount"
  logs_account_id                    = "789"
  log_level                          = "full"
}
```

## Argument Reference

The following arguments are supported:

* `parent_id` - (Optional) The newly created sub-account's parent id. If not specified, the invoking account will be assigned as the parent.
* `ref_id` - (Optional) Customer specific identifier for this operation.
* `sub_account_name` - (Mandatory) SubAccount name.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, `none`, `default`.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the SubAccount.
* `sub_account_id` - SubAccount ID
* `is_for_special_ssl_configuration` - Is using special SSL configuration.
* `support_level` - The CNAME record name.

## Import

SubAccount can be imported using the `id`, e.g.:

```
$ terraform import incapsula_subaccount.demo 1234
```
