---
layout: "incapsula"
page_title: "Incapsula: subaccount"
sidebar_current: "docs-incapsula-resource-subaccount"
description: |-
  Provides a Incapsula SubAccount resource.
---

# incapsula_subaccount

Provides a Incapsula SubAccount resource. 
Please note any change on this resource will force create a new SubAccount instance, 
while non-supported terraform dependent resources won't auto create 
(Users for example) 

## Example Usage

```hcl
resource "incapsula_subaccount" "example-subaccount" {
  sub_account_name                   = "Example SubAccount"
  logs_account_id                    = "789"
  log_level                          = "full"
  data_storage_region                = "US"
}
```

## Argument Reference

The following arguments are supported:

* `sub_account_name` - (Mandatory) SubAccount name.
* `parent_id` - (Optional) The newly created sub-account's parent id. If not specified, the invoking account will be assigned as the parent.
* `ref_id` - (Optional) Customer specific identifier for this operation.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, `none`, `default`.
* `data_storage_region` - (Optional) Default data region of the sub-account for newly created sites. Options are `APAC`, `EU`, `US` and `AU`. Defaults to `US`.

SubAccount can be imported using the `id`, e.g.:

```
$ terraform import incapsula_subaccount.demo 1234
```
