---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_subaccount"
description: |- 
  Provides a Incapsula SubAccount resource.
---

# incapsula_subaccount

Provides an Incapsula SubAccount resource. 
Please note any change on this resource will create a new SubAccount instance, 
while non-supported terraform dependent resources, such as users, 
will not be automatically created.

## Example Usage

```hcl
resource "incapsula_subaccount" "example-subaccount" {
  sub_account_name                     = "Example SubAccount"
  logs_account_id                      = "789"
  log_level                            = "full"
  data_storage_region                  = "US"
  enable_http2_for_new_sites           = true
  enable_http2_to_origin_for_new_sites = true
}
```

## Argument Reference

The following arguments are supported:

* `sub_account_name` - (Mandatory) SubAccount name.
* `parent_id` - (Optional) The newly created subaccount's parent id. If not specified, the account identified by the authentication parameters will be assigned as the parent.
* `ref_id` - (Optional) Customer specific identifier for this operation.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, `none`, `default`.
* `data_storage_region` - (Optional) Default data region of the subaccount for newly created sites. Options are `APAC`, `EU`, `US` and `AU`. Defaults to `US`.
* `enable_http2_for_new_sites` - (Optional) Use this option to enable HTTP/2 support for traffic between end-users (visitors) and Imperva for newly created SSL sites. Options are `true` and `false`. Defaults to `true`.
* `enable_http2_to_origin_for_new_sites` - (Optional) Use this option to enable HTTP/2 support for traffic between Imperva and your origin server for newly created SSL sites. This option can only be 'true' once 'enable_http2_for_new_sites' is enabled for newly created sites. Options are `true` and `false`. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Sub Account ID.
```
$ terraform import incapsula_subaccount.demo 1234
```
